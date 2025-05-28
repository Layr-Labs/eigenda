package blobstore

import "errors"

var (
	ErrBlobNotFound           = errors.New("blob not found")
	ErrMetadataNotFound       = errors.New("metadata not found")
	ErrAlreadyExists          = errors.New("record already exists")
	ErrInvalidStateTransition = errors.New("invalid state transition")
)
