package gateway

import "github.com/0glabs/0g-storage-client/common/api"

var (
	ErrFileNotFound        = api.NewBusinessError(101, "File not found", nil)
	ErrFilePruned          = api.NewBusinessError(102, "File already pruned", nil)
	ErrFileNotFinalized    = api.NewBusinessError(103, "File not finalized yet", nil)
	ErrFilePathNotFound    = api.NewBusinessError(104, "File path not found", nil)
	ErrFileSizeTooLarge    = api.NewBusinessError(105, "File size too large", nil)
	ErrFileTypeUnsupported = api.NewBusinessError(106, "File type unsupported", nil)
)
