package common

import "errors"

var (
	ErrBlobNotFound     = errors.New("blob not found")
	ErrMetadataNotFound = errors.New("metadata not found")
	ErrAlreadyExists    = errors.New("record already exists")
)
