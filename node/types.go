package node

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/zero-gravity-labs/zerog-storage-client/core/merkle"
)

type Status struct {
	ConnectedPeers uint        `json:"connectedPeers"`
	LogSyncHeight  uint64      `json:"logSyncHeight"`
	LogSyncBlock   common.Hash `json:"logSyncBlock"`
}

type Transaction struct {
	StreamIds       []*hexutil.Big `json:"streamIds"`
	Data            []byte         `json:"data"` // in-place data
	DataMerkleRoot  common.Hash    `json:"dataMerkleRoot"`
	StartEntryIndex uint64         `json:"startEntryIndex"`
	Size            uint64         `json:"size"` // file size in bytes
	Seq             uint64         `json:"seq"`
}

type FileInfo struct {
	Tx             Transaction `json:"tx"`
	Finalized      bool        `json:"finalized"`
	IsCached       bool        `json:"isCached"`
	UploadedSegNum uint64      `json:"uploadedSegNum"`
}

type SegmentWithProof struct {
	Root     common.Hash  `json:"root"`     // file merkle root
	Data     []byte       `json:"data"`     // segment data
	Index    uint64       `json:"index"`    // segment index
	Proof    merkle.Proof `json:"proof"`    // segment merkle proof
	FileSize uint64       `json:"fileSize"` // file size
}

type Value struct {
	Version uint64 `json:"version"` // key version
	Data    []byte `json:"data"`    // value data
	Size    uint64 `json:"size"`    // value total size
}

type KeyValue struct {
	Version uint64 `json:"version"` // key version
	Key     []byte `json:"key"`     // value key
	Data    []byte `json:"data"`    // value data
	Size    uint64 `json:"size"`    // value total size
}
