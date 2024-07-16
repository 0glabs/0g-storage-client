package node

import (
	"github.com/0glabs/0g-storage-client/core/merkle"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Status sync status of 0g storage node
type Status struct {
	ConnectedPeers uint        `json:"connectedPeers"`
	LogSyncHeight  uint64      `json:"logSyncHeight"`
	LogSyncBlock   common.Hash `json:"logSyncBlock"`
}

// Transaction on-chain transaction about a file
type Transaction struct {
	StreamIds       []*hexutil.Big `json:"streamIds"`       // transaction related stream id, used for KV operations
	Data            []byte         `json:"data"`            // in-place data
	DataMerkleRoot  common.Hash    `json:"dataMerkleRoot"`  // data merkle root
	StartEntryIndex uint64         `json:"startEntryIndex"` // start entry index in on-chain flow contract
	Size            uint64         `json:"size"`            // file size in bytes
	Seq             uint64         `json:"seq"`             // sequence id in on-chain flow contract
}

// FileInfo information about a file responsed from 0g storage node
type FileInfo struct {
	Tx             Transaction `json:"tx"`             // on-chain transaction
	Finalized      bool        `json:"finalized"`      // whether the file has been finalized in the storage node
	IsCached       bool        `json:"isCached"`       // whether the file is cached in the storage node
	UploadedSegNum uint64      `json:"uploadedSegNum"` // the number of uploaded segments
}

// SegmentWithProof data segment with merkle proof
type SegmentWithProof struct {
	Root     common.Hash  `json:"root"`     // file merkle root
	Data     []byte       `json:"data"`     // segment data
	Index    uint64       `json:"index"`    // segment index
	Proof    merkle.Proof `json:"proof"`    // segment merkle proof
	FileSize uint64       `json:"fileSize"` // file size
}

// Value KV value
type Value struct {
	Version uint64 `json:"version"` // key version
	Data    []byte `json:"data"`    // value data
	Size    uint64 `json:"size"`    // value total size
}

// KeyValue KV key and value
type KeyValue struct {
	Version uint64 `json:"version"` // key version
	Key     []byte `json:"key"`     // value key
	Data    []byte `json:"data"`    // value data
	Size    uint64 `json:"size"`    // value total size
}
