package dir

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

var (
	CodecVersion    = uint16(1)
	CodecMagicBytes = crypto.Keccak256([]byte("0g-storage-client-dir-codec"))
)

// EncodeFsNode serializes an FsNode into a custom binary format with magic bytes and versioning.
func EncodeFsNode(node *FsNode) ([]byte, error) {
	// Serialize the FsNode to JSON
	mdata, err := json.Marshal(node)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal FsNode to JSON")
	}

	// Prepare the buffer for the final binary data
	buf := new(bytes.Buffer)

	// Write magic bytes
	if _, err := buf.Write(CodecMagicBytes); err != nil {
		return nil, errors.Wrap(err, "failed to write magic bytes")
	}

	// Write codec version
	if err := binary.Write(buf, binary.BigEndian, CodecVersion); err != nil {
		return nil, errors.Wrap(err, "failed to write codec version")
	}

	// Write the JSON data
	if _, err := buf.Write(mdata); err != nil {
		return nil, errors.Wrap(err, "failed to write JSON data")
	}

	return buf.Bytes(), nil
}

// DecodeFsNode deserializes an FsNode from a custom binary format with magic bytes and versioning.
func DecodeFsNode(data []byte) (*FsNode, error) {
	// Create a reader from the input data
	buf := bytes.NewReader(data)

	// Verify magic bytes
	magicBytes := make([]byte, len(CodecMagicBytes))
	if _, err := io.ReadFull(buf, magicBytes); err != nil {
		return nil, errors.Wrap(err, "failed to read magic bytes")
	}
	if !bytes.Equal(magicBytes, CodecMagicBytes) {
		return nil, errors.New("invalid magic bytes")
	}

	// Verify codec version
	var version uint16
	if err := binary.Read(buf, binary.BigEndian, &version); err != nil {
		return nil, errors.Wrap(err, "failed to read codec version")
	}
	if version != CodecVersion {
		return nil, errors.Errorf("unsupported codec version: got %d, expected %d", version, CodecVersion)
	}

	// Read the remaining data (the JSON-encoded FsNode)
	mdata, err := io.ReadAll(buf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read JSON data")
	}

	// Deserialize the FsNode from JSON data
	var node FsNode
	if err := json.Unmarshal(mdata, &node); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal FsNode from JSON")
	}

	return &node, nil
}
