package dir

import (
	"github.com/google/btree"
	"github.com/pkg/errors"
)

// DiffStatus represents the status of a node in the diff.
type DiffStatus string

const (
	Added     DiffStatus = "added"
	Removed   DiffStatus = "removed"
	Modified  DiffStatus = "modified"
	Unchanged DiffStatus = "unchanged"
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

	if node.Type == Directory {
		diffNode.Entries = btree.NewG(2, func(a, b *DiffNode) bool {
			return a.Node.Name < b.Node.Name
		})
	}

	return diffNode
}

// Diff compares two directories and returns a DiffNode tree with the differences.
func Diff(current, next *FsNode, recursive bool) (*DiffNode, error) {
	if current.Type != Directory || next.Type != Directory {
		return nil, errors.New("diff is only supported for directories")
	}

	return diff(current, next, recursive), nil
}

// diff is a recursive function that computes the differences between two directory nodes.
func diff(current, next *FsNode, recursive bool) *DiffNode {
	root := NewDiffNode(current, Unchanged)

	// processes entries from the current directory.
	for _, currentEntry := range current.Entries {
		nextEntry, found := next.Search(currentEntry.Name)
		if !found {
			root.Entries.ReplaceOrInsert(NewDiffNode(currentEntry, Removed))
			root.Status = Modified
			continue
		}

		if currentEntry.Hash == nextEntry.Hash {
			root.Entries.ReplaceOrInsert(NewDiffNode(currentEntry, Unchanged))
			continue
		}

		root.Status = Modified
		if recursive && currentEntry.Type == Directory && nextEntry.Type == Directory {
			subDiff := diff(currentEntry, nextEntry, recursive)
			root.Entries.ReplaceOrInsert(subDiff)
		} else {
			root.Entries.ReplaceOrInsert(NewDiffNode(currentEntry, Modified))
		}
	}

	// processes entries from the next directory that were not found in the current directory.
	for _, nextEntry := range next.Entries {
		if _, found := current.Search(nextEntry.Name); !found {
			root.Status = Modified
			root.Entries.ReplaceOrInsert(NewDiffNode(nextEntry, Added))
		}
	}

	return root
}
