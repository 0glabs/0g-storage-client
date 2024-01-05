package core

type Iterator interface {
	Next() (bool, error)
	Current() []byte
}

func IteratorPaddedSize(dataSize int64, flowPadding bool) uint64 {
	chunks := NumSplits(dataSize, DefaultChunkSize)
	var paddedSize uint64
	if flowPadding {
		paddedChunks, _ := ComputePaddedSize(chunks)
		paddedSize = paddedChunks * DefaultChunkSize
	} else {
		paddedSize = chunks * DefaultChunkSize
	}
	return paddedSize
}
