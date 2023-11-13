package indexer

import "errors"

var (
	ErrNotImplemented         = errors.New("not implemented")
	ErrIncorrectObject        = errors.New("incorrect object")
	ErrUnrecognizedFork       = errors.New("unrecognized fork")
	ErrHeadersNotOrdered      = errors.New("headers not ordered")
	ErrIncorrectEvent         = errors.New("incorrect event payload")
	ErrOperatorNotFound       = errors.New("operator not found")
	ErrWrongObjectFromIndexer = errors.New("indexer returned error of wrong type")
)
