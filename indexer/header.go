package indexer

import (
	"bytes"
	"errors"
)

var (
	ErrHeadersUnordered = errors.New("headers unordered")
	ErrHeaderNotFound   = errors.New("header not found")
)

type Header struct {
	BlockHash     [32]byte
	PrevBlockHash [32]byte
	Number        uint64
	Finalized     bool
	CurrentFork   string
	IsUpgrade     bool
}

func (h *Header) After(prev *Header) bool {
	return h.PrevBlockHash == prev.BlockHash
}

func (h *Header) BlockHashIs(hash []byte) bool {
	return bytes.Equal(h.BlockHash[:], hash)
}

func (h *Header) Equals(other *Header) bool {
	return h.BlockHash == other.BlockHash
}

type Headers []*Header

func (hh Headers) Empty() bool {
	return hh.Len() == 0
}

// Len returns the number of headers in the header list.
func (hh Headers) Len() int {
	return len(hh)
}

// First returns the first header in the header list.
func (hh Headers) First() *Header {
	return hh[0]
}

// Last returns the last header in the header list.
func (hh Headers) Last() *Header {
	return hh[len(hh)-1]
}

func (hh Headers) OK() error {
	if !hh.IsOrdered() {
		return ErrHeadersUnordered
	}
	return nil
}

// IsOrdered tells whether a list of headers is a proper chain
func (hh Headers) IsOrdered() bool {
	for ind := 1; ind < len(hh); ind++ {
		if hh[ind].PrevBlockHash != hh[ind-1].BlockHash {
			return false
		}
	}
	return true
}

// GetHeaderByNumber gives the header with a given number. Assumes headers are ordered
func (hh Headers) GetHeaderByNumber(number uint64) (*Header, error) {
	offset := int(hh[0].Number)
	ind := int(number) - offset
	if ind < 0 || ind >= len(hh) {
		return nil, ErrHeaderNotFound
	}

	return hh[ind], nil
}
