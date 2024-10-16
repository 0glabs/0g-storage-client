package core

import (
	"io"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	// ErrFileRequired is returned when manipulate on a folder.
	ErrFileRequired = errors.New("file required")

	// ErrFileEmpty is returned when empty file opened.
	ErrFileEmpty = errors.New("file is empty")
)

// File implement of IterableData, the underlying is a file on disk
type File struct {
	os.FileInfo
	underlying *os.File
	paddedSize uint64
}

var _ IterableData = (*File)(nil)

func (file *File) Read(buf []byte, offset int64) (int, error) {
	n, err := file.underlying.ReadAt(buf, offset)
	// unexpected IO error
	if !errors.Is(err, io.EOF) {
		return 0, err
	}
	return n, nil
}

func Exists(name string) (bool, error) {
	file, err := os.Open(name)
	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	defer file.Close()

	return true, nil
}

// Open create a File from a file on disk
func Open(name string) (*File, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, ErrFileRequired
	}

	if info.Size() == 0 {
		return nil, ErrFileEmpty
	}

	return &File{
		FileInfo:   info,
		underlying: file,
		paddedSize: IteratorPaddedSize(info.Size(), true),
	}, nil
}

// MerkleRoot returns the merkle root hash of a file on disk
func MerkleRoot(filename string) (common.Hash, error) {
	file, err := Open(filename)
	if err != nil {
		return common.Hash{}, errors.WithMessage(err, "failed to open file")
	}
	defer file.Close()

	// Generate the Merkle tree from the file content
	tree, err := MerkleTree(file)
	if err != nil {
		return common.Hash{}, errors.WithMessage(err, "failed to create merkle tree")
	}

	// Return the root hash of the Merkle tree
	return tree.Root(), nil
}

func (file *File) Close() error {
	return file.underlying.Close()
}

func (file *File) NumChunks() uint64 {
	return NumSplits(file.Size(), DefaultChunkSize)
}

func (file *File) NumSegments() uint64 {
	return NumSplits(file.Size(), DefaultSegmentSize)
}

func (file *File) PaddedSize() uint64 {
	return file.paddedSize
}

func (file *File) Iterate(offset int64, batch int64, flowPadding bool) Iterator {
	if batch%DefaultChunkSize > 0 {
		panic("batch size should align with chunk size")
	}
	dataSize := file.Size()
	return &FileIterator{
		file:       file.underlying,
		buf:        make([]byte, batch),
		offset:     offset,
		fileSize:   dataSize,
		paddedSize: IteratorPaddedSize(dataSize, flowPadding),
	}
}

type FileIterator struct {
	file       *os.File
	buf        []byte // buffer to read data from file
	bufSize    int    // actual data size in buffer
	fileSize   int64
	paddedSize uint64
	offset     int64 // offset to read data
}

var _ Iterator = (*FileIterator)(nil)

func (it *FileIterator) Next() (bool, error) {
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

func (it *FileIterator) clearBuffer() {
	it.bufSize = 0
}

func (it *FileIterator) paddingZeros(length int) {
	paddingZeros(it.buf, it.bufSize, length)
	it.bufSize += length
	it.offset += int64(length)
}

func (it *FileIterator) Current() []byte {
	return it.buf[:it.bufSize]
}
