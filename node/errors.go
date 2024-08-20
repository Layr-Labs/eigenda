package node

import (
	"errors"
)

var (
	ErrKeyAlreadyExists     = errors.New("commit already exists as key in kvstore")
	ErrKeyNotFound          = errors.New("commit not found in kvstore")
	ErrKeyExpired           = errors.New("commit is expired")
	ErrKeyNotFoundOrExpired = errors.New("data is either expired or not found")
)
