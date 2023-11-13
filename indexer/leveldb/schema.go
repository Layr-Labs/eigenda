package leveldb

import (
	"github.com/Layr-Labs/eigenda/indexer"
	"math"
	"reflect"
	"strconv"
)

var (
	headerKeyPrefix    = []byte("h-")
	finalizedHeaderKey = []byte("latest-finalized-header")
)

func newHeaderKey(v uint64) []byte {
	return append(headerKeyPrefix, newHeaderKeySuffix(v)...)
}

func newHeaderKeySuffix(v uint64) []byte {
	return []byte(strconv.FormatUint(math.MaxUint64-v, 16))
}

func newAccumulatorKey(acc indexer.Accumulator, header *indexer.Header) []byte {
	return append(newAccumulatorKeyPrefix(acc), newHeaderKeySuffix(header.Number)...)
}

func newAccumulatorKeyPrefix(acc indexer.Accumulator) []byte {
	accTyp := reflect.TypeOf(acc)
	if accTyp.Kind() == reflect.Pointer {
		accTyp = accTyp.Elem()
	}
	return []byte("a-" + accTyp.Name() + "-")
}

type headerEntry struct {
	Header          *indexer.Header
	AccumulatorKeys [][]byte
}

func newHeaderEntry(header *indexer.Header) *headerEntry {
	return &headerEntry{Header: header}
}

func (e *headerEntry) UpdateAccumulatorKeys(key []byte) *headerEntry {
	e.AccumulatorKeys = append(e.AccumulatorKeys, key)
	return e
}

type accumulatorEntry struct {
	HeaderNumber    uint64
	AccumulatorData []byte
}

func newAccumulatorEntry(headerNo uint64, accData []byte) *accumulatorEntry {
	return &accumulatorEntry{
		HeaderNumber:    headerNo,
		AccumulatorData: accData,
	}
}
