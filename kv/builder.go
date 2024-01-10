package kv

import (
	"errors"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const maxSetSize = 1 << 16 // 64K

const maxKeySize = 1 << 24 // 16.7M

const maxQuerySize = 1024 * 256

var errSizeTooLarge = errors.New("size too large")

var errKeyTooLarge = errors.New("key too large")

var errKeyIsEmpty = errors.New("key is empty")

type builder struct {
	streamIds map[common.Hash]bool // to build tags
}

func (builder *builder) AddStreamId(streamId common.Hash) {
	builder.streamIds[streamId] = true
}

func (builder *builder) BuildTags(sorted ...bool) []byte {
	var ids []common.Hash

	for k := range builder.streamIds {
		ids = append(ids, k)
	}

	if len(sorted) > 0 {
		if sorted[0] {
			sort.SliceStable(ids, func(i, j int) bool {
				return ids[i].Hex() < ids[j].Hex()
			})
		}
	}

	return CreateTags(ids...)
}

type StreamDataBuilder struct {
	AccessControlBuilder
	version uint64
	reads   map[common.Hash]map[string]bool
	writes  map[common.Hash]map[string][]byte
}

func NewStreamDataBuilder(version uint64) *StreamDataBuilder {
	return &StreamDataBuilder{
		AccessControlBuilder: AccessControlBuilder{
			builder: builder{
				streamIds: make(map[common.Hash]bool),
			},
			controls: make([]AccessControl, 0),
		},
		version: version,
		reads:   make(map[common.Hash]map[string]bool),
		writes:  make(map[common.Hash]map[string][]byte),
	}
}

func (builder *StreamDataBuilder) Build(sorted ...bool) (*StreamData, error) {
	var err error
	data := StreamData{
		Version: builder.version,
	}

	// controls
	if data.Controls, err = builder.AccessControlBuilder.Build(); err != nil {
		return nil, err
	}

	// reads
	for streamId, keys := range builder.reads {
		for k := range keys {
			key := hexutil.MustDecode(k)
			if len(key) > maxKeySize {
				return nil, errKeyTooLarge
			}
			if len(key) == 0 {
				return nil, errKeyIsEmpty
			}
			data.Reads = append(data.Reads, StreamRead{
				StreamId: streamId,
				Key:      key,
			})

			if len(data.Reads) > maxSetSize {
				return nil, errSizeTooLarge
			}
		}
	}

	// writes
	for streamId, keys := range builder.writes {
		for k, d := range keys {
			key := hexutil.MustDecode(k)
			if len(key) > maxKeySize {
				return nil, errKeyTooLarge
			}
			if len(key) == 0 {
				return nil, errKeyIsEmpty
			}
			data.Writes = append(data.Writes, StreamWrite{
				StreamId: streamId,
				Key:      key,
				Data:     d,
			})

			if len(data.Writes) > maxSetSize {
				return nil, errSizeTooLarge
			}
		}
	}

	if len(sorted) > 0 {
		if sorted[0] {
			sort.SliceStable(data.Reads, func(i, j int) bool {
				streamIdI := data.Reads[i].StreamId.Hex()
				streamIdJ := data.Reads[j].StreamId.Hex()
				if streamIdI == streamIdJ {
					return hexutil.Encode(data.Reads[i].Key) < hexutil.Encode(data.Reads[j].Key)
				} else {
					return streamIdI < streamIdJ
				}
			})
			sort.SliceStable(data.Writes, func(i, j int) bool {
				streamIdI := data.Writes[i].StreamId.Hex()
				streamIdJ := data.Writes[j].StreamId.Hex()
				if streamIdI == streamIdJ {
					return hexutil.Encode(data.Writes[i].Key) < hexutil.Encode(data.Writes[j].Key)
				} else {
					return streamIdI < streamIdJ
				}
			})
		}
	}

	return &data, nil
}

func (builder *StreamDataBuilder) SetVersion(version uint64) *StreamDataBuilder {
	builder.version = version
	return builder
}

func (builder *StreamDataBuilder) Watch(streamId common.Hash, key []byte) *StreamDataBuilder {
	if keys, ok := builder.reads[streamId]; ok {
		keys[hexutil.Encode(key)] = true
	} else {
		builder.reads[streamId] = make(map[string]bool)
		builder.reads[streamId][hexutil.Encode(key)] = true
	}

	return builder
}

func (builder *StreamDataBuilder) Set(streamId common.Hash, key []byte, data []byte) *StreamDataBuilder {
	builder.AddStreamId(streamId)

	if keys, ok := builder.writes[streamId]; ok {
		keys[hexutil.Encode(key)] = data
	} else {
		builder.writes[streamId] = make(map[string][]byte)
		builder.writes[streamId][hexutil.Encode(key)] = data
	}

	return builder
}

type AccessControlBuilder struct {
	builder
	controls []AccessControl
}

func (builder *AccessControlBuilder) Build() ([]AccessControl, error) {
	if len(builder.controls) > maxSetSize {
		return nil, errSizeTooLarge
	}

	return builder.controls, nil
}

func (builder *AccessControlBuilder) withControl(t AccessControlType, streamId common.Hash, account *common.Address, key []byte) *AccessControlBuilder {
	builder.AddStreamId(streamId)

	builder.controls = append(builder.controls, AccessControl{
		Type:     t,
		StreamId: streamId,
		Account:  account,
		Key:      key,
	})

	return builder
}

func (builder *AccessControlBuilder) GrantAdminRole(streamId common.Hash, account common.Address) *AccessControlBuilder {
	return builder.withControl(AclTypeGrantAdminRole, streamId, &account, nil)
}

func (builder *AccessControlBuilder) RenounceAdminRole(streamId common.Hash) *AccessControlBuilder {
	return builder.withControl(AclTypeRenounceAdminRole, streamId, nil, nil)
}

func (builder *AccessControlBuilder) SetKeyToSpecial(streamId common.Hash, key []byte) *AccessControlBuilder {
	return builder.withControl(AclTypeSetKeyToSpecial, streamId, nil, key)
}

func (builder *AccessControlBuilder) SetKeyToNormal(streamId common.Hash, key []byte) *AccessControlBuilder {
	return builder.withControl(AclTypeSetKeyToNormal, streamId, nil, key)
}

func (builder *AccessControlBuilder) GrantWriteRole(streamId common.Hash, account common.Address) *AccessControlBuilder {
	return builder.withControl(AclTypeGrantWriteRole, streamId, &account, nil)
}

func (builder *AccessControlBuilder) RevokeWriteRole(streamId common.Hash, account common.Address) *AccessControlBuilder {
	return builder.withControl(AclTypeRevokeWriteRole, streamId, &account, nil)
}

func (builder *AccessControlBuilder) RenounceWriteRole(streamId common.Hash) *AccessControlBuilder {
	return builder.withControl(AclTypeRenounceWriteRole, streamId, nil, nil)
}

func (builder *AccessControlBuilder) GrantSpecialWriteRole(streamId common.Hash, key []byte, account common.Address) *AccessControlBuilder {
	return builder.withControl(AclTypeGrantSpecialWriteRole, streamId, &account, key)
}

func (builder *AccessControlBuilder) RevokeSpecialWriteRole(streamId common.Hash, key []byte, account common.Address) *AccessControlBuilder {
	return builder.withControl(AclTypeRevokeSpecialWriteRole, streamId, &account, key)
}

func (builder *AccessControlBuilder) RenounceSpecialWriteRole(streamId common.Hash, key []byte) *AccessControlBuilder {
	return builder.withControl(AclTypeRenounceSpecialWriteRole, streamId, nil, key)
}
