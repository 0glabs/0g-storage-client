package unixfs

import (
	"encoding/json"

	"github.com/ipfs/go-unixfsnode/data"
	"github.com/ipfs/go-unixfsnode/data/builder"
	"github.com/pkg/errors"
)

// EncodeFsNode serializes an FsNode to bytes, embedding it in a UnixFS structure.
func EncodeFsNode(node *FsNode) ([]byte, error) {
	// Serialize the FsNode to JSON
	mdata, err := json.Marshal(node)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to marshal FsNode to JSON")
	}

	// Build a UnixFS node containing the serialized FsNode
	ufs, err := builder.BuildUnixFS(func(b *builder.Builder) {
		// Set the data type to raw data
		builder.DataType(b, data.Data_Raw)
		// Set the serialized FsNode data
		builder.Data(b, mdata)
	})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to build UnixFS node")
	}

	// Encode the UnixFS node to bytes
	encodedData := data.EncodeUnixFSData(ufs)
	return encodedData, nil
}

// DecodeFsNode deserializes a FsNode from bytes, extracting it from a UnixFS structure.
func DecodeFsNode(b []byte) (*FsNode, error) {
	// Decode the UnixFS data from the byte slice
	ufs, err := data.DecodeUnixFSData(b)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to decode UnixFS data")
	}

	// Verify the data type is raw as expected
	if ufs.FieldDataType().Int() != data.Data_Raw {
		return nil, errors.New("unexpected data type, expected raw data")
	}

	// Deserialize the FsNode from the UnixFS data
	var node FsNode
	err = json.Unmarshal(ufs.FieldData().Must().Bytes(), &node)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to unmarshal FsNode from JSON")
	}

	return &node, nil
}
