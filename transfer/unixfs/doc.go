// Package unixfs implements a data format to add support for folder operations within the 0g storage client.
// The 0g storage node has no built-in support for folder operations, so this package provides a way to
// represent folder hierarchies and perform operations such as adding, updating, and deleting files
// within a directory structure. The folder structure is serialized as JSON, merged into a protocol buffer
// based format (UnixFs), which can be stored and retrieved from the 0g storage node.
//
// The main features of this package include:
//
//   - Defining the FileNode structure, which represents files and directories in a nested hierarchical format.
//   - Serializing a FileNode structure into a JSON representation, assembled into a UnixFS raw file that can be
//     uploaded to a 0g storage node.
//   - Deserializing the JSON representation back into a FileNode structure for further operations.
//   - Supporting operations such as adding, updating, and deleting files within the directory structure.
//   - Providing methods to traverse and visualize the directory structure.
//
// This package enables the 0g storage client to handle complex directory structures by using a simple and
// extensible data format, allowing for more advanced file system operations beyond single file storage.
package unixfs
