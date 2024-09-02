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
		return nil, errors.WithMessage(err, "failed to marshal FsNode to JSON")
	}

	// Check if the metadata is too large
	if len(mdata) > math.MaxUint32 {
		return nil, errors.New("metadata too large")
	}

	// Calculate the total length needed for the byte slice:
	// MagicBytes + CodecVersion (2 bytes) + Metadata Length (4 bytes) + JSON Metadata + Optional Raw Data
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

	// Write raw data to buffer
	buf := bytes.NewBuffer(data)
	if err := node.writeRawData(buf); err != nil {
		return nil, errors.WithMessage(err, "failed to write raw file data")
	}

	return buf.Bytes(), nil
}

// writeRawData recursively writes raw data to the provided buffer for FsNodes of type FileTypeRaw.
func (node *FsNode) writeRawData(buf *bytes.Buffer) error {
	switch node.Type {
	case FileTypeRaw:
		if _, err := buf.Write(node.Data); err != nil {
			return err
		}
	case FileTypeDirectory:
		for _, entry := range node.Entries {
			if err := entry.writeRawData(buf); err != nil {
				return err
			}
		}
	}

	return nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// It decodes the FsNode from a binary format.
func (node *FsNode) UnmarshalBinary(data []byte) error {
	offset := int64(0)
	datalen := int64(len(data))

	// Verify magic bytes
	if datalen < offset+int64(len(CodecMagicBytes)) {
		return errors.New("not enough data to read magic bytes")
	}

	magicBytes := data[offset : offset+int64(len(CodecMagicBytes))]
	if !bytes.Equal(magicBytes, CodecMagicBytes) {
		return errors.New("invalid magic bytes")
	}
	offset += int64(len(CodecMagicBytes))

	// Verify codec version
	if datalen < offset+2 {
		return errors.New("not enough data to read codec version")
	}
	version := binary.BigEndian.Uint16(data[offset : offset+2])
	if version != CodecVersion {
		return errors.Errorf("unsupported codec version: got %d, expected %d", version, CodecVersion)
	}
	offset += 2

	// Read metadata length
	if datalen < offset+4 {
		return errors.New("not enough data to read metadata length")
	}
	metadataLength := binary.BigEndian.Uint32(data[offset : offset+4])
	offset += 4

	// Read JSON metadata
	if datalen < offset+int64(metadataLength) {
		return errors.New("not enough data to read JSON metadata")
	}
	mdata := data[offset : offset+int64(metadataLength)]
	offset += int64(metadataLength)

	// Deserialize the FsNode from JSON metadata
	if err := json.Unmarshal(mdata, node); err != nil {
		return errors.WithMessage(err, "failed to unmarshal FsNode from JSON")
	}

	// Read and decode raw data
	buf := bytes.NewReader(nil)
	if offset < datalen {
		buf = bytes.NewReader(data[offset:])
	}
	if err := node.readRawData(buf); err != nil {
		return errors.WithMessage(err, "failed to read raw file data")
	}

	return nil
}

// readRawData recursively reads raw data for FsNodes of type FileTypeRaw.
func (node *FsNode) readRawData(buf *bytes.Reader) error {
	switch node.Type {
	case FileTypeRaw:
		data := make([]byte, node.Size)
		n, err := buf.Read(data)
		if err != nil {
			return err
		}
		if n != int(node.Size) {
			return errors.Errorf("expected to read %d bytes, but read %d bytes", node.Size, n)
		}

		node.Data = data
	case FileTypeDirectory:
		for _, entry := range node.Entries {
			if err := entry.readRawData(buf); err != nil {
				return err
			}
		}
	}

	return nil
}
