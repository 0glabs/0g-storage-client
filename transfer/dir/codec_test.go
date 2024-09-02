package dir_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/0glabs/0g-storage-client/transfer/dir"
	"github.com/ethereum/go-ethereum/common"
)

func TestEncodeDecodeFsNode(t *testing.T) {
	// Create a sample FsNode structure
	rootNode := dir.FsNode{
		Name: "root",
		Type: dir.FileTypeDirectory,
		Hash: common.HexToHash("0x12"),
		Entries: []*dir.FsNode{
			{
				Name: "file1.txt",
				Type: dir.FileTypeFile,
				Hash: common.HexToHash("0x23"),
				Size: 1024,
			},
			{
				Name: "rawfile2.bin",
				Type: dir.FileTypeRaw,
				Hash: common.HexToHash("0x34"),
				Size: 7,
				Data: []byte("rawdata"),
			},
			{
				Name: "symlink",
				Type: dir.FileTypeSymbolic,
				Hash: common.HexToHash("0x45"),
				Link: "/path/to/target",
			},
			{
				Name: "subdir",
				Type: dir.FileTypeDirectory,
				Hash: common.HexToHash("0x56"),
				Entries: []*dir.FsNode{
					{
						Name: "file2.txt",
						Type: dir.FileTypeFile,
						Hash: common.HexToHash("0x67"),
						Size: 2048,
					},
					{
						Name: "rawfile.bin",
						Type: dir.FileTypeRaw,
						Hash: common.HexToHash("0x78"),
						Size: 7,
						Data: []byte("rawdata"),
					},
				},
			},
		},
	}

	// Encode the FsNode to bytes
	encodedData, err := rootNode.MarshalBinary()
	if err != nil {
		t.Fatalf("Binary marshal failed: %v", err)
	}

	// Decode the bytes back into an FsNode
	var decodedNode dir.FsNode
	err = decodedNode.UnmarshalBinary(encodedData)
	if err != nil {
		t.Fatalf("Binary unmarshal failed: %v", err)
	}

	// Compare the original and decoded FsNode to ensure they are the same
	if !reflect.DeepEqual(rootNode, decodedNode) {
		t.Errorf("Expected FsNode to be equal, but got %v and %v", rootNode, decodedNode)
	}
}

func TestInvalidMagicBytes(t *testing.T) {
	// Create a sample FsNode structure and encode it
	originalNode := dir.FsNode{
		Name: "testfile.txt",
		Type: dir.FileTypeFile,
		Hash: common.HexToHash("0x1234"),
		Size: 1024,
	}
	encodedData, err := originalNode.MarshalBinary()
	if err != nil {
		t.Fatalf("EncodeFsNode failed: %v", err)
	}

	// Modify the magic bytes to simulate corruption
	encodedData[0] ^= 0xFF

	// Attempt to decode the corrupted data
	var decodedNode dir.FsNode
	err = decodedNode.UnmarshalBinary(encodedData)
	if err == nil || err.Error() != "invalid magic bytes" {
		t.Fatalf("expected error 'invalid magic bytes', got %v", err)
	}
}

func TestInvalidVersion(t *testing.T) {
	// Create a sample FsNode structure and encode it
	originalNode := dir.FsNode{
		Name: "testfile.txt",
		Type: dir.FileTypeFile,
		Hash: common.HexToHash("0x1234"),
		Size: 1024,
	}
	encodedData, err := originalNode.MarshalBinary()
	if err != nil {
		t.Fatalf("EncodeFsNode failed: %v", err)
	}

	// Modify the version bytes to simulate an unsupported version
	encodedData[len(dir.CodecMagicBytes)] ^= 0xFF

	// Attempt to decode the data with the modified version
	var decodedNode dir.FsNode
	err = decodedNode.UnmarshalBinary(encodedData)
	if err == nil || !bytes.Contains([]byte(err.Error()), []byte("unsupported codec version")) {
		t.Fatalf("expected version error, got %v", err)
	}
}

func TestIncompleteRawData(t *testing.T) {
	// Create an FsNode structure with raw data
	originalNode := dir.FsNode{
		Name: "rawfile",
		Type: dir.FileTypeRaw,
		Hash: common.HexToHash("0x1234"),
		Size: 1024,
		Data: make([]byte, 1024),
	}

	// Encode the FsNode to bytes
	encodedData, err := originalNode.MarshalBinary()
	if err != nil {
		t.Fatalf("Binary marshal failed: %v", err)
	}

	// Truncate the encoded data to simulate incomplete raw data
	truncatedData := encodedData[:len(encodedData)-10]

	// Attempt to decode the truncated data
	var decodedNode dir.FsNode
	err = decodedNode.UnmarshalBinary(truncatedData)
	if err == nil || !bytes.Contains([]byte(err.Error()), []byte("expected to read")) {
		t.Fatalf("expected read error, got %v", err)
	}
}

func TestEmptyNodeEncodeDecode(t *testing.T) {
	// Create an empty FsNode
	originalNode := dir.FsNode{}

	// Encode the FsNode to bytes
	encodedData, err := originalNode.MarshalBinary()
	if err != nil {
		t.Fatalf("Binary marshal failed: %v", err)
	}

	// Decode the bytes back into an FsNode
	var decodedNode dir.FsNode
	err = decodedNode.UnmarshalBinary(encodedData)
	if err != nil {
		t.Fatalf("Binary unmarshal failed: %v", err)
	}

	// Compare the original and decoded FsNode to ensure they are the same
	if !reflect.DeepEqual(originalNode, decodedNode) {
		t.Errorf("Expected FsNode to be equal, but got %v and %v", originalNode, decodedNode)
	}
}
