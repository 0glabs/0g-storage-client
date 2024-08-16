package unixfs

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

// FileType represents the type of a file in the FsNode structure.
type FileType string

const (
	File      FileType = "file"
	Directory FileType = "directory"
)

var (
	ErrDirectoryRequired = errors.New("operation supported only on directories")
	ErrEntryNotFound     = errors.New("directory entry not found")
	ErrFileNotFound      = errors.New("file not found")
)

// FsNode represents a node in the filesystem hierarchy.
type FsNode struct {
	Name    string    `json:"name"`              // File or directory name
	Type    FileType  `json:"type"`              // Type of the node (file or directory)
	Size    int64     `json:"size,omitempty"`    // File size in bytes (omitted for directories)
	Hash    string    `json:"hash,omitempty"`    // Merkle root hash (only for files)
	Entries []*FsNode `json:"entries,omitempty"` // Directory entries (only for directories)
}

// Add adds a new file to the directory and keeps the entries sorted by name.
func (node *FsNode) Add(file *FsNode) error {
	if node.Type != Directory {
		return ErrDirectoryRequired
	}

	// Find the index where the file should be inserted
	index, found := node.Search(file.Name)

	// If the file already exists, replace it
	if found {
		node.Entries[index] = file
		return nil
	}

	// Insert and sort the new file at the correct position
	node.Entries = append(node.Entries, file)
	node.sort()

	return nil
}

// BatchAdd adds multiple files or directories to the current directory and sorts once.
func (node *FsNode) BatchAdd(entries []*FsNode) error {
	if node.Type != Directory {
		return ErrDirectoryRequired
	}

	// Append all new entries
	node.Entries = append(node.Entries, entries...)

	// Sort the entries after all additions
	node.sort()
	return nil
}

// Update updates an existing file's metadata within the directory.
func (node *FsNode) Update(file *FsNode) error {
	if node.Type != Directory {
		return ErrDirectoryRequired
	}

	index, found := node.Search(file.Name)
	if !found {
		return ErrFileNotFound
	}

	// Update the existing entry
	node.Entries[index] = file
	return nil
}

// Rename renames a file or directory within the directory and keeps the entries sorted.
func (node *FsNode) Rename(oldName, newName string) error {
	if node.Type != Directory {
		return ErrDirectoryRequired
	}

	// Find the index of the old file/directory
	index, found := node.Search(oldName)
	if !found {
		return ErrFileNotFound
	}

	// Retrieve the existing entry
	entry := node.Entries[index]

	// If the name is the same, there's nothing to rename
	if entry.Name == newName {
		return nil
	}

	// Just in case new name already exists
	if _, found := node.Search(newName); found {
		return errors.New("file with the new name already exists")
	}

	// Remove the old entry
	node.Entries = append(node.Entries[:index], node.Entries[index+1:]...)

	// Update the entry's name
	entry.Name = newName

	// Re-insert the renamed entry at the correct position
	return node.Add(entry)
}

// Delete removes a file from the directory.
func (node *FsNode) Delete(fileName string) error {
	if node.Type != Directory {
		return ErrDirectoryRequired
	}

	index, found := node.Search(fileName)
	if !found {
		return ErrFileNotFound
	}

	// Remove the file
	node.Entries = append(node.Entries[:index], node.Entries[index+1:]...)
	return nil
}

// Locate locates a child node by its path within the FsNode.
func (node *FsNode) Locate(path string) (parent, current *FsNode, found bool) {
	// Normalize the path and split it into parts
	parts := strings.Split(filepath.Clean(path), string(os.PathSeparator))

	// Start from the current node
	current = node

	// Traverse the path
	for _, part := range parts {
		if current.Type != Directory {
			return nil, nil, false // Cannot traverse a non-directory node
		}

		index, found := current.Search(part)
		if !found {
			return nil, nil, false // Part not found in current directory
		}

		parent = current                 // Save the parent node
		current = current.Entries[index] // Move to the next level
	}

	return parent, current, true
}

// Search finds the index of an entry by name using binary search.
// It returns the index and a boolean indicating if the entry was found.
func (node *FsNode) Search(name string) (int, bool) {
	index := sort.Search(len(node.Entries), func(i int) bool {
		return node.Entries[i].Name >= name
	})
	if index < len(node.Entries) && node.Entries[index].Name == name {
		return index, true
	}
	return index, false
}

// sort sorts the entries slice by the name field in alphabetical order.
func (node *FsNode) sort() {
	sort.Slice(node.Entries, func(i, j int) bool {
		return node.Entries[i].Name < node.Entries[j].Name
	})
}

// BuildFileTree recursively builds a file tree from the given file path.
func BuildFileTree(path string) (*FsNode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// Handle file case
	if !info.IsDir() {
		merkleRoot, err := GetMerkleTreeRootOfFile(path)
		if err != nil {
			return nil, err
		}

		return &FsNode{
			Name: info.Name(),
			Type: File,
			Size: info.Size(),
			Hash: merkleRoot,
		}, nil
	}

	// Handle directory case
	node := &FsNode{
		Name:    info.Name(),
		Type:    Directory,
		Entries: []*FsNode{},
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var children []*FsNode
	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		entryNode, err := BuildFileTree(entryPath)
		if err != nil {
			return nil, err
		}
		children = append(children, entryNode)
	}

	if len(children) > 0 {
		return node, node.BatchAdd(children)
	}

	return node, nil
}
