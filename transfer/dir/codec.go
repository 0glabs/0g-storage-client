package dir

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"encoding/json"
	"math"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

var (
	// Assert that customEncoding implements the correct interfaces.
	_ encoding.BinaryMarshaler   = (*FsNode)(nil)
	_ encoding.BinaryUnmarshaler = (*FsNode)(nil)

	CodecVersion    = uint16(1)
	CodecMagicBytes = crypto.Keccak256([]byte("0g-storage-client-dir-codec"))
)

// MarshalBinary implements the encoding.BinaryMarshaler interface.
// It encodes the FsNode into a binary format.
func (node *FsNode) MarshalBinary() ([]byte, error) {
	// Serialize the FsNode to JSON
	mdata, err := json.Marshal(node)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to marshal `FsNode` to JSON")
	}

	// Check if the json metadata is too large
	if len(mdata) > math.MaxUint32 {
		return nil, errors.New("the json marshalled data is too large")
	}

	// Calculate the total length needed for the byte slice:
	// MagicBytes + CodecVersion (2 bytes) + Metadata Length (4 bytes) + JSON Metadata
	totalLength := len(CodecMagicBytes) + 2 + 4 + len(mdata)

	// Create the byte slice with the calculated total length
	data := make([]byte, totalLength)

	// Offset to track the position in the byte slice
	offset := int64(0)

	// Write magic bytes
	copy(data[offset:], CodecMagicBytes)
	offset += int64(len(CodecMagicBytes))

	// Write codec version
	binary.BigEndian.PutUint16(data[offset:], CodecVersion)
	offset += 2

	// Write metadata length
	binary.BigEndian.PutUint32(data[offset:], uint32(len(mdata)))
	offset += 4

	// Write JSON data
	copy(data[offset:], mdata)
	return data, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// It decodes the FsNode from a binary format.
func (node *FsNode) UnmarshalBinary(data []byte) error {
	// Verify magic bytes
	if len(data) < len(CodecMagicBytes) {
		return errors.New("not enough data to read magic bytes")
	}

	magicBytes := data[:len(CodecMagicBytes)]
	if !bytes.Equal(magicBytes, CodecMagicBytes) {
		return errors.New("invalid magic bytes")
	}
	data = data[len(CodecMagicBytes):]

	// Verify codec version
	if len(data) < 2 {
		return errors.New("not enough data to read codec version")
	}
	version := binary.BigEndian.Uint16(data[:2])
	if version != CodecVersion {
		return errors.Errorf("unsupported codec version: got %d, expected %d", version, CodecVersion)
	}
	data = data[2:]

	// Read metadata length
	if len(data) < 4 {
		return errors.New("not enough data to read metadata length")
	}
	metalen := binary.BigEndian.Uint32(data[:4])
	data = data[4:]

	// Read JSON metadata
	if uint32(len(data)) < metalen {
		return errors.New("not enough data to read JSON metadata")
	}
	data = data[:metalen]

	// Deserialize the FsNode from JSON metadata
	if err := json.Unmarshal(data, node); err != nil {
		return errors.WithMessage(err, "failed to unmarshal `FsNode` from JSON")
	}
	return nil
}
