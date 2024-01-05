package core

type DataInMemory struct {
	underlying []byte
}

var _ IterableData = (*DataInMemory)(nil)

func NewDataInMemory(data []byte) *DataInMemory {
	return &DataInMemory{
		underlying: data,
	}
}

func (data *DataInMemory) NumChunks() uint64 {
	return NumSplits(int64(len(data.underlying)), DefaultChunkSize)
}

func (data *DataInMemory) NumSegments() uint64 {
	return NumSplits(int64(len(data.underlying)), DefaultSegmentSize)
}

func (data *DataInMemory) Size() int64 {
	return int64(len(data.underlying))
}

func (data *DataInMemory) Iterate(offset int64, batch int64, flowPadding bool) Iterator {
	if batch%DefaultChunkSize > 0 {
		panic("batch size should align with chunk size")
	}
	dataSize := int64(len(data.underlying))
	return &MemoryDataIterator{
		data:       data,
		buf:        make([]byte, batch),
		offset:     offset,
		dataSize:   dataSize,
		paddedSize: IteratorPaddedSize(dataSize, flowPadding),
	}
}

type MemoryDataIterator struct {
	data       *DataInMemory
	buf        []byte // buffer to read data from file
	bufSize    int    // actual data size in buffer
	dataSize   int64
	paddedSize uint64
	offset     int64 // offset to read data
}

var _ Iterator = (*MemoryDataIterator)(nil)

func (it *MemoryDataIterator) Next() (bool, error) {
	// Reject invalid offset
	if it.offset < 0 || uint64(it.offset) >= it.paddedSize {
		return false, nil
	}

	var expectedBufSize int
	maxAvailableLength := it.paddedSize - uint64(it.offset)
	if maxAvailableLength >= uint64(len(it.buf)) {
		expectedBufSize = len(it.buf)
	} else {
		expectedBufSize = int(maxAvailableLength)
	}

	if it.offset >= it.dataSize {
		it.paddingZeros(expectedBufSize)
		return true, nil
	}

	n := copy(it.buf, it.data.underlying[it.offset:])
	it.offset += int64(n)
	it.bufSize = n

	if n == expectedBufSize {
		return true, nil
	}

	it.paddingZeros(expectedBufSize - n)

	return true, nil
}

func (it *MemoryDataIterator) paddingZeros(length int) {
	startOffset := it.bufSize
	for i := 0; i < length; i++ {
		it.buf[startOffset+i] = 0
	}
	it.bufSize += length
	it.offset += int64(length)
}

func (it *MemoryDataIterator) Current() []byte {
	return it.buf[:it.bufSize]
}
