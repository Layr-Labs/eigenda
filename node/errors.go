package node

import (
	"errors"
)

var (
	ErrKeyNotFound = errors.New("commit not found in levelDB")
)
