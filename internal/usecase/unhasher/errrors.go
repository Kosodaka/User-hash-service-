package unhasher

import "errors"

var (
	ErrQueryIsEmpty           = errors.New("uc - query is empty")
	ErrQueryWithInjection     = errors.New("uc - query contains statement that have data modification")
	ErrOnlySelectAllowed      = errors.New("uc - only select statements are allowed")
	ErrEmptyDomain            = errors.New("uc - empty domain")
	ErrFailedUnhashing        = errors.New("uc - failed to unhash data")
	ErrCreeateTempDir         = errors.New("uc - failed to create temp dir")
	ErrCreateInitialFile      = errors.New("uc - failed to create initial file")
	ErrCreateNewFile          = errors.New("uc - failed to create new file")
	ErrWriteDataToFile        = errors.New("uc -failed to write data to file")
	ErrCreateFileForRemaining = errors.New("uc - failed to create file for remaining data")
	ErrWriteRemainingData     = errors.New("uc - failed to write remaining data")
)
