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
	originalNode := &dir.FsNode{
		Name: "testfile.txt",
		Type: dir.File,
		Hash: common.HexToHash("0x1234"),
		Size: 1024,
	}

	// Encode the FsNode to bytes
	encodedData, err := dir.EncodeFsNode(originalNode)
	if err != nil {
		t.Fatalf("EncodeFsNode failed: %v", err)
	}

	// Decode the bytes back into an FsNode
	decodedNode, err := dir.DecodeFsNode(encodedData)
	if err != nil {
		t.Fatalf("DecodeFsNode failed: %v", err)
	}

	// Compare the original and decoded FsNode to ensure they are the same
	if !reflect.DeepEqual(originalNode, decodedNode) {
		t.Errorf("expected FsNode to be equal, but got %v and %v", originalNode, decodedNode)
	}
}

func TestInvalidMagicBytes(t *testing.T) {
	// Create a sample FsNode structure and encode it
	originalNode := &dir.FsNode{
		Name: "testfile.txt",
		Type: dir.File,
		Hash: common.HexToHash("0x1234"),
		Size: 1024,
	}
	encodedData, err := dir.EncodeFsNode(originalNode)
	if err != nil {
		t.Fatalf("EncodeFsNode failed: %v", err)
	}

	// Modify the magic bytes to simulate corruption
	encodedData[0] ^= 0xFF

	// Attempt to decode the corrupted data
	_, err = dir.DecodeFsNode(encodedData)
	if err == nil || err.Error() != "invalid magic bytes" {
		t.Fatalf("expected error 'invalid magic bytes', got %v", err)
	}
}

func TestInvalidVersion(t *testing.T) {
	// Create a sample FsNode structure and encode it
	originalNode := &dir.FsNode{
		Name: "testfile.txt",
		Type: dir.File,
		Hash: common.HexToHash("0x1234"),
		Size: 1024,
	}
	encodedData, err := dir.EncodeFsNode(originalNode)
	if err != nil {
		t.Fatalf("EncodeFsNode failed: %v", err)
	}

	// Modify the version bytes to simulate an unsupported version
	encodedData[len(dir.CodecMagicBytes)] ^= 0xFF

	// Attempt to decode the data with the modified version
	_, err = dir.DecodeFsNode(encodedData)
	if err == nil || !bytes.Contains([]byte(err.Error()), []byte("unsupported codec version")) {
		t.Fatalf("expected version error, got %v", err)
	}
}
