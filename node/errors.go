package node

import (
	"errors"
)

var (
	ErrKeyAlreadyExists     = errors.New("commit already exists as key in db")
	ErrKeyNotFound          = errors.New("commit not found in db")
	ErrKeyExpired           = errors.New("commit is expired")
	ErrKeyNotFoundOrExpired = errors.New("data is either expired or not found")
)
