package dir

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/google/btree"
	"github.com/pkg/errors"
)

// DiffStatus represents the status of a node in the diff.
type DiffStatus string

const (
	DiffStatusAdded     DiffStatus = "added"
	DiffStatusRemoved   DiffStatus = "removed"
	DiffStatusModified  DiffStatus = "modified"
	DiffStatusUnchanged DiffStatus = "unchanged"
)

// DiffNode represents a node in the diff structure with its status.
type DiffNode struct {
	Node    *FsNode                  // The original node
	Status  DiffStatus               // Diff status of the node
	Entries *btree.BTreeG[*DiffNode] // Directory entries as a B-Tree
}

// NewDiffNode creates a new DiffNode.
func NewDiffNode(node *FsNode, status DiffStatus) *DiffNode {
	diffNode := &DiffNode{
		Node:   node,
		Status: status,
	}

	if node.Type == FileTypeDirectory {
		diffNode.Entries = btree.NewG(2, func(a, b *DiffNode) bool {
			return a.Node.Name < b.Node.Name
		})
	}

	return diffNode
}

// Diff compares two directories and returns a DiffNode tree with the differences.
func Diff(current, next *FsNode) (*DiffNode, error) {
	if current.Type != FileTypeDirectory || next.Type != FileTypeDirectory {
		return nil, errors.New("diff is only supported for directories")
	}

	return diff(current, next), nil
}

// diff is a recursive function that computes the differences between two directory nodes.
func diff(current, next *FsNode) *DiffNode {
	root := NewDiffNode(current, DiffStatusUnchanged)

	// processes entries from the current directory.
	for _, currentEntry := range current.Entries {
		nextEntry, found := next.Search(currentEntry.Name)
		if !found {
			root.Entries.ReplaceOrInsert(NewDiffNode(currentEntry, DiffStatusRemoved))
			root.Status = DiffStatusModified
			continue
		}

		if currentEntry.Type == FileTypeDirectory && nextEntry.Type == FileTypeDirectory {
			subDiff := diff(currentEntry, nextEntry)
			root.Entries.ReplaceOrInsert(subDiff)
			if subDiff.Status != DiffStatusUnchanged {
				root.Status = DiffStatusModified
			}
		} else if currentEntry.Equal(nextEntry) {
			root.Entries.ReplaceOrInsert(NewDiffNode(currentEntry, DiffStatusUnchanged))
		} else {
			root.Entries.ReplaceOrInsert(NewDiffNode(currentEntry, DiffStatusModified))
			root.Status = DiffStatusModified
		}
	}

	// processes entries from the next directory that were not found in the current directory.
	for _, nextEntry := range next.Entries {
		if _, found := current.Search(nextEntry.Name); !found {
			root.Status = DiffStatusModified
			root.Entries.ReplaceOrInsert(NewDiffNode(nextEntry, DiffStatusAdded))
		}
	}

	return root
}

// PrettyPrint prints the DiffNode tree in a human-readable format with a tree skeleton structure.
func PrettyPrint(root *DiffNode) {
	prettyPrint(root, 0, false, nil)
}

func prettyPrint(node *DiffNode, depth int, isLast bool, prefixes []string) {
	if depth == 0 { // Root directory case
		fmt.Println(node.Node.Name)
		printChildEntries(node, 1, nil)
		return
	}

	// Create prefix and branch for current node
	prefix := strings.Join(prefixes, "")
	branch := "├─"
	if isLast {
		branch = "└─"
	}

	// Print node with optional status
	coloredStatus := colorStringForDiffStatus(node.Status)
	fmt.Printf("%s%s %s %s\n", prefix, branch, node.Node.Name, coloredStatus)

	// Update prefixes for children
	newPrefixes := append(prefixes, "│  ")
	if isLast {
		newPrefixes[len(newPrefixes)-1] = "    "
	}

	// Recursively print children if it's a directory
	if node.Node.Type == FileTypeDirectory {
		printChildEntries(node, depth+1, newPrefixes)
	}
}

func printChildEntries(node *DiffNode, depth int, prefixes []string) {
	if node.Entries == nil || node.Entries.Len() == 0 {
		return
	}

	// Get the last node in a BTree efficiently by descending to the last element.
	var last *DiffNode
	node.Entries.Descend(func(item *DiffNode) bool {
		last = item
		return false // Stop the iteration as soon as we get the last node.
	})

	// Print all child entries of the current node
	node.Entries.Ascend(func(item *DiffNode) bool {
		// Check if this is the last entry in the directory
		isLast := item == last
		prettyPrint(item, depth, isLast, prefixes)
		return true
	})
}

// colorStringForDiffStatus returns a color string based on the diff status.
func colorStringForDiffStatus(status DiffStatus) string {
	switch status {
	case DiffStatusAdded:
		return color.GreenString("[+]")
	case DiffStatusRemoved:
		return color.RedString("[-]")
	case DiffStatusModified:
		return color.YellowString("[*]")
	default:
		return ""
	}
}
