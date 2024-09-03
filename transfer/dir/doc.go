// Package dir provides support for folder operations within the 0g storage client. Since the 0g storage
// node does not natively support folder operations, this package introduces a way to represent and manage
// hierarchical directory structures, including encoding, decoding, and comparing directories and files.
// The folder structure is serialized into a binary format that includes a JSON representation of the data,
// which can be stored on and retrieved from the 0g storage node.
//
// The main features of this package include:
//
//   - Defining the FsNode structure, which models files and directories in a nested, hierarchical format.
//   - Serializing an FsNode structure into a binary format that includes a JSON representation, suitable
//     for storage on a 0g storage node.
//   - Deserializing the binary format back into an FsNode structure for further manipulation and operations.
//   - Supporting the comparison of two directory structures to identify differences such as added, removed,
//     or modified files.
//
// This package enables the 0g storage client to manage complex directory structures using an efficient
// and extensible data format, facilitating advanced file system operations beyond single file storage.
package dir
