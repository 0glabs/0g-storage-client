package file

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

type Iterator struct {
	file       *os.File
	buf        []byte // buffer to read data from file
	bufSize    int    // actual data size in buffer
	fileSize   int64
	paddedSize uint64
	offset     int64 // offset to read data
}

func NewSegmentIterator(file *os.File, fileSize int64, offset int64, flowPadding bool) *Iterator {
	return NewIterator(file, fileSize, offset, DefaultSegmentSize, flowPadding)
}

func NewIterator(file *os.File, fileSize int64, offset int64, batch int64, flowPadding bool) *Iterator {
	if batch%DefaultChunkSize > 0 {
		panic("batch size should align with chunk size")
	}

	buf := make([]byte, batch)

	chunks := numSplits(fileSize, DefaultChunkSize)
	var paddedSize uint64
	if flowPadding {
		paddedChunks, _ := computePaddedSize(chunks)
		paddedSize = paddedChunks * DefaultChunkSize
	} else {
		paddedSize = chunks * DefaultChunkSize
	}

	return &Iterator{
		file:       file,
		buf:        buf,
		offset:     offset,
		fileSize:   fileSize,
		paddedSize: paddedSize,
	}
}

func (it *Iterator) Next() (bool, error) {
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

	it.clearBuffer()

	if it.offset >= it.fileSize {
		it.paddingZeros(expectedBufSize)
		return true, nil
	}

	n, err := it.file.ReadAt(it.buf, it.offset)
	it.bufSize = n
	it.offset += int64(n)

	// not reach EOF
	if n == expectedBufSize {
		return true, nil
	}

	// unexpected IO error
	if !errors.Is(err, io.EOF) {
		return false, err
	}

	if n > expectedBufSize {
		// should never happen
		panic("load more data from file than expected")
	}

	it.paddingZeros(expectedBufSize - n)

	return true, nil
}

func (it *Iterator) clearBuffer() {
	it.bufSize = 0
}

func (it *Iterator) paddingZeros(length int) {
	startOffset := it.bufSize
	for i := 0; i < length; i++ {
		it.buf[startOffset+i] = 0
	}
	it.bufSize += length
	it.offset += int64(length)
}

func (it *Iterator) Current() []byte {
	return it.buf[:it.bufSize]
}
