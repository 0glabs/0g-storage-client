package unixfs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/0glabs/0g-storage-client/transfer/unixfs"
	"github.com/stretchr/testify/assert"
)

func setupTestFixtures(t *testing.T) string {
	/*
		Create a temporary directory to hold the test fixtures:
		fixtures/
		├── file1.txt
		├── file2.txt
		└── subdir/
			└── file3.txt
	*/
	tempDir, err := os.MkdirTemp("", "buildfiletree_test")
	assert.NoError(t, err)

	assert.NoError(t, os.Mkdir(filepath.Join(tempDir, "subdir"), 0755))
	assert.NoError(t, os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("This is file 1"), 0644))
	assert.NoError(t, os.WriteFile(filepath.Join(tempDir, "file2.txt"), []byte("This is file 2"), 0644))
	assert.NoError(t, os.WriteFile(filepath.Join(tempDir, "subdir/file3.txt"), []byte("This is file 3"), 0644))

	return tempDir
}

func cleanupTestFixtures(t *testing.T, tempDir string) {
	assert.NoError(t, os.RemoveAll(tempDir))
}

func TestBuildFileTree(t *testing.T) {
	// Setup the test fixtures
	tempDir := setupTestFixtures(t)
	defer cleanupTestFixtures(t, tempDir)

	// Build the file tree from the fixtures directory
	tree, err := unixfs.BuildFileTree(tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, tree)

	// Validate the root node
	assert.Equal(t, filepath.Base(tempDir), tree.Name)
	assert.Equal(t, unixfs.Directory, tree.Type)
	assert.Empty(t, tree.Hash)

	// Validate the root node has 3 entries
	assert.Len(t, tree.Entries, 3)

	// Validate the files and subdirectory
	var file1, file2, subdir *unixfs.FsNode
	for _, entry := range tree.Entries {
		switch entry.Name {
		case "file1.txt":
			file1 = entry
		case "file2.txt":
			file2 = entry
		case "subdir":
			subdir = entry
		}
	}

	assert.NotNil(t, file1)
	assert.Equal(t, unixfs.File, file1.Type)
	assert.Equal(t, int64(14), file1.Size) // "This is file 1" length is 14
	assert.NotEmpty(t, file1.Hash)

	assert.NotNil(t, file2)
	assert.Equal(t, unixfs.File, file2.Type)
	assert.Equal(t, int64(14), file2.Size) // "This is file 2" length is 14
	assert.NotEmpty(t, file2.Hash)

	assert.NotNil(t, subdir)
	assert.Equal(t, unixfs.Directory, subdir.Type)
	assert.Len(t, subdir.Entries, 1) // subdir should contain 1 file
	assert.Empty(t, subdir.Hash)

	// Validate the file in subdir
	subfile := subdir.Entries[0]
	assert.Equal(t, "file3.txt", subfile.Name)
	assert.Equal(t, unixfs.File, subfile.Type)
	assert.Equal(t, int64(14), subfile.Size) // "This is file 3" length is 14
	assert.NotEmpty(t, subfile.Hash)
}

func TestAddFile(t *testing.T) {
	dir := &unixfs.FsNode{
		Name:    "root",
		Type:    unixfs.Directory,
		Entries: []*unixfs.FsNode{},
	}

	// Add a file
	err := dir.Add(&unixfs.FsNode{Name: "file1.txt", Type: unixfs.File, Size: 1234})
	assert.NoError(t, err)
	assert.Len(t, dir.Entries, 1)
	assert.Equal(t, "file1.txt", dir.Entries[0].Name)

	// Add another file and check if it's sorted
	err = dir.Add(&unixfs.FsNode{Name: "file2.txt", Type: unixfs.File, Size: 5678})
	assert.NoError(t, err)
	assert.Len(t, dir.Entries, 2)
	assert.Equal(t, "file1.txt", dir.Entries[0].Name)
	assert.Equal(t, "file2.txt", dir.Entries[1].Name)

	// Add a file with the same name to replace it
	err = dir.Add(&unixfs.FsNode{Name: "file1.txt", Type: unixfs.File, Size: 4321})
	assert.NoError(t, err)
	assert.Len(t, dir.Entries, 2)
	assert.Equal(t, "file1.txt", dir.Entries[0].Name)
	assert.Equal(t, int64(4321), dir.Entries[0].Size)
}

func TestBatchAddFiles(t *testing.T) {
	dir := &unixfs.FsNode{
		Name:    "root",
		Type:    unixfs.Directory,
		Entries: []*unixfs.FsNode{},
	}

	// Batch add files
	err := dir.BatchAdd([]*unixfs.FsNode{
		{Name: "file2.txt", Type: unixfs.File, Size: 5678},
		{Name: "file1.txt", Type: unixfs.File, Size: 1234},
		{Name: "file3.txt", Type: unixfs.File, Size: 91011},
	})
	assert.NoError(t, err)
	assert.Len(t, dir.Entries, 3)

	// Check if files are sorted
	assert.Equal(t, "file1.txt", dir.Entries[0].Name)
	assert.Equal(t, "file2.txt", dir.Entries[1].Name)
	assert.Equal(t, "file3.txt", dir.Entries[2].Name)
}

func TestUpdateFile(t *testing.T) {
	dir := &unixfs.FsNode{
		Name:    "root",
		Type:    unixfs.Directory,
		Entries: []*unixfs.FsNode{},
	}

	// Add a file
	err := dir.Add(&unixfs.FsNode{Name: "file1.txt", Type: unixfs.File, Size: 1234})
	assert.NoError(t, err)

	// Update the file
	err = dir.Update(&unixfs.FsNode{Name: "file1.txt", Type: unixfs.File, Size: 5678})
	assert.NoError(t, err)
	assert.Len(t, dir.Entries, 1)
	assert.Equal(t, "file1.txt", dir.Entries[0].Name)
	assert.Equal(t, int64(5678), dir.Entries[0].Size)

	// Try to update a non-existent file
	err = dir.Update(&unixfs.FsNode{Name: "file2.txt", Type: unixfs.File, Size: 91011})
	assert.ErrorIs(t, err, unixfs.ErrFileNotFound)
}

func TestRenameFile(t *testing.T) {
	dir := &unixfs.FsNode{
		Name:    "root",
		Type:    unixfs.Directory,
		Entries: []*unixfs.FsNode{},
	}

	// Add a file
	err := dir.Add(&unixfs.FsNode{Name: "file1.txt", Type: unixfs.File, Size: 1234})
	assert.NoError(t, err)

	// Rename the file
	err = dir.Rename("file1.txt", "file2.txt")
	assert.NoError(t, err)
	assert.Len(t, dir.Entries, 1)
	assert.Equal(t, "file2.txt", dir.Entries[0].Name)

	// Rename to an existing name (should replace)
	err = dir.Add(&unixfs.FsNode{Name: "file3.txt", Type: unixfs.File, Size: 5678})
	assert.NoError(t, err)
	err = dir.Rename("file2.txt", "file3.txt")
	assert.Error(t, err)
}

func TestDeleteFile(t *testing.T) {
	dir := &unixfs.FsNode{
		Name:    "root",
		Type:    unixfs.Directory,
		Entries: []*unixfs.FsNode{},
	}

	// Add a file
	err := dir.Add(&unixfs.FsNode{Name: "file1.txt", Type: unixfs.File, Size: 1234})
	assert.NoError(t, err)

	// Delete the file
	err = dir.Delete("file1.txt")
	assert.NoError(t, err)
	assert.Len(t, dir.Entries, 0)

	// Try to delete a non-existent file
	err = dir.Delete("file2.txt")
	assert.ErrorIs(t, err, unixfs.ErrFileNotFound)
}

func TestLocate(t *testing.T) {
	root := &unixfs.FsNode{
		Name:    "root",
		Type:    unixfs.Directory,
		Entries: []*unixfs.FsNode{},
	}

	// Build a directory structure
	root.BatchAdd([]*unixfs.FsNode{
		{Name: "dir1", Type: unixfs.Directory, Entries: []*unixfs.FsNode{
			{Name: "file1.txt", Type: unixfs.File, Size: 1234},
			{Name: "file2.txt", Type: unixfs.File, Size: 5678},
		}},
		{Name: "dir2", Type: unixfs.Directory, Entries: []*unixfs.FsNode{
			{Name: "subdir1", Type: unixfs.Directory, Entries: []*unixfs.FsNode{
				{Name: "file3.txt", Type: unixfs.File, Size: 91011},
			}},
		}},
		{Name: "file4.txt", Type: unixfs.File, Size: 4321},
	})

	// Locate a file in the root directory
	parent, node, found := root.Locate("file4.txt")
	assert.True(t, found)
	assert.NotNil(t, node)
	assert.Equal(t, "file4.txt", node.Name)
	assert.Equal(t, root, parent)

	// Locate a file in a subdirectory
	parent, node, found = root.Locate("dir1/file1.txt")
	assert.True(t, found)
	assert.NotNil(t, node)
	assert.Equal(t, "file1.txt", node.Name)
	assert.Equal(t, root.Entries[0], parent)

	// Locate a file in a nested subdirectory
	parent, node, found = root.Locate("dir2/subdir1/file3.txt")
	assert.True(t, found)
	assert.NotNil(t, node)
	assert.Equal(t, "file3.txt", node.Name)
	assert.Equal(t, root.Entries[1].Entries[0], parent)

	// Try to locate a non-existent file
	parent, node, found = root.Locate("dir2/subdir1/fileX.txt")
	assert.False(t, found)
	assert.Nil(t, node)
	assert.Nil(t, parent)

	// Try to locate a file through a non-directory path
	parent, node, found = root.Locate("file4.txt/fileX.txt")
	assert.False(t, found)
	assert.Nil(t, node)
	assert.Nil(t, parent)
}

func TestSearchFile(t *testing.T) {
	dir := &unixfs.FsNode{
		Name:    "root",
		Type:    unixfs.Directory,
		Entries: []*unixfs.FsNode{},
	}

	// Add files
	err := dir.BatchAdd([]*unixfs.FsNode{
		{Name: "file1.txt", Type: unixfs.File, Size: 1234},
		{Name: "file2.txt", Type: unixfs.File, Size: 5678},
		{Name: "file3.txt", Type: unixfs.File, Size: 91011},
	})
	assert.NoError(t, err)

	// Search for an existing file
	index, found := dir.Search("file2.txt")
	assert.True(t, found)
	assert.Equal(t, 1, index)

	// Search for a non-existent file
	index, found = dir.Search("file4.txt")
	assert.False(t, found)
	assert.Equal(t, 3, index) // index is where the file would be inserted
}
