package core

import "errors"

// DataInMemory implement of IterableData, the underlying is memory data
type DataInMemory struct {
	underlying []byte
	offset     int64
	size       int64
	paddedSize uint64
}

var _ IterableData = (*DataInMemory)(nil)

// NewDataInMemory creates DataInMemory from given data
func NewDataInMemory(data []byte) (*DataInMemory, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	return &DataInMemory{
		underlying: data,
		offset:     0,
		size:       int64(len(data)),
		paddedSize: IteratorPaddedSize(int64(len(data)), true),
	}, nil
}

func (data *DataInMemory) Read(buf []byte, offset int64) (int, error) {
	n := copy(buf, data.underlying[data.offset+offset:])
	return n, nil
}

func (data *DataInMemory) NumChunks() uint64 {
	return NumSplits(int64(len(data.underlying)), DefaultChunkSize)
}

func (data *DataInMemory) NumSegments() uint64 {
	return NumSplits(int64(len(data.underlying)), DefaultSegmentSize)
}

func (data *DataInMemory) Size() int64 {
	return data.size
}

func (data *DataInMemory) Offset() int64 {
	return data.offset
}

func (data *DataInMemory) PaddedSize() uint64 {
	return data.paddedSize
}

func (data *DataInMemory) Split(fragmentSize int64) []IterableData {
	fragments := make([]IterableData, 0)
	for offset := data.offset; offset < data.offset+data.size; offset += fragmentSize {
		size := min(data.size-offset, fragmentSize)
		fragment := &DataInMemory{
			underlying: data.underlying,
			offset:     offset,
			size:       size,
			paddedSize: IteratorPaddedSize(size, true),
		}
		fragments = append(fragments, fragment)
	}
	return fragments
}
