package leveldb

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

type mockAccumulatorObjectV1 struct {
	Balance uint64
}

type mockAccumulatorObjectV2 struct {
	Balance uint64
}

type mockAccumulator struct{}

func (acc mockAccumulator) InitializeObject(_ indexer.Header) (indexer.AccumulatorObject, error) {
	return nil, nil
}

func (acc mockAccumulator) UpdateObject(_ indexer.AccumulatorObject, _ *indexer.Header, _ indexer.Event) (indexer.AccumulatorObject, error) {
	return nil, nil
}

func (acc mockAccumulator) SerializeObject(object indexer.AccumulatorObject, _ indexer.UpgradeFork) ([]byte, error) {
	return encode(object)
}

func (acc mockAccumulator) DeserializeObject(data []byte, fork indexer.UpgradeFork) (indexer.AccumulatorObject, error) {
	var objV1 mockAccumulatorObjectV1
	var objV2 mockAccumulatorObjectV2

	switch fork {
	case "genesis":
		err := decode(data, &objV1)
		return objV1, err

	case "exodus":
		err := decode(data, &objV2)
		return objV2, err

	default:
		return nil, errors.New("unknown fork")
	}
}

func blockHash(t *testing.T, hash string) [32]byte {
	t.Helper()
	var hashBytes [32]byte

	v, err := hex.DecodeString(hash)
	assert.NoError(t, err)

	copy(hashBytes[:], v)
	return hashBytes
}

func newTestStore(t *testing.T) *HeaderStore {
	t.Helper()

	s, err := NewHeaderStore("", func(path string) (*leveldb.DB, error) {
		return leveldb.Open(storage.NewMemStorage(), &opt.Options{Filter: filter.NewBloomFilter(10)})
	})
	assert.NoError(t, err)

	return s
}

func newTestHeaders(t *testing.T) indexer.Headers {
	t.Helper()

	var headerList []map[string]any

	data, err := os.ReadFile("testdata/headers.json")
	assert.NoError(t, err)
	assert.NoError(t, json.Unmarshal(data, &headerList))

	var res indexer.Headers

	for i := len(headerList) - 1; i >= 0; i-- {
		header := headerList[i]
		res = append(res, &indexer.Header{
			BlockHash:     blockHash(t, header["BlockHash"].(string)),
			PrevBlockHash: blockHash(t, header["PrevBlockHash"].(string)),
			Number:        uint64(header["Number"].(float64)),
			CurrentFork:   header["CurrentFork"].(string),
			IsUpgrade:     header["IsUpgrade"].(bool),
		})
	}

	return res
}

func newTestHeadersWithFork(t *testing.T, fork int) indexer.Headers {
	t.Helper()

	var headerList []map[string]any

	var data []byte
	var err error
	if fork == 1 {
		data, err = os.ReadFile("testdata/fork1.json")
	} else if fork == 2 {
		data, err = os.ReadFile("testdata/fork2.json")
	}

	assert.NoError(t, err)
	assert.NoError(t, json.Unmarshal(data, &headerList))

	var res indexer.Headers

	for i := len(headerList) - 1; i >= 0; i-- {
		header := headerList[i]
		res = append(res, &indexer.Header{
			BlockHash:     blockHash(t, header["BlockHash"].(string)),
			PrevBlockHash: blockHash(t, header["PrevBlockHash"].(string)),
			Number:        uint64(header["Number"].(float64)),
		})
	}

	return res
}

func TestHeaderStore_AddHeaders(t *testing.T) {
	headers := newTestHeadersWithFork(t, 1)
	fork := newTestHeadersWithFork(t, 2)

	tests := []struct {
		name        string
		store       func(t *testing.T) *HeaderStore
		headers     indexer.Headers
		expected    indexer.Headers
		expectedErr error
	}{
		{
			name:        "add headers no headers",
			store:       newTestStore,
			headers:     indexer.Headers{},
			expected:    indexer.Headers{},
			expectedErr: nil,
		},
		{
			name:        "add headers to empty database",
			store:       newTestStore,
			headers:     headers,
			expected:    headers,
			expectedErr: nil,
		},
		{
			name: "add headers to the end of non-empty database",
			store: func(t *testing.T) *HeaderStore {
				t.Helper()
				store := newTestStore(t)
				_, _ = store.AddHeaders(headers[:5])
				return store
			},
			headers:     headers[5:],
			expected:    headers[5:],
			expectedErr: nil,
		},
		{
			name: "add headers no new headers",
			store: func(t *testing.T) *HeaderStore {
				t.Helper()
				store := newTestStore(t)
				_, _ = store.AddHeaders(headers[:5])
				return store
			},
			headers:     headers[:5],
			expected:    nil,
			expectedErr: nil,
		},
		{
			name: "add headers intersecting headers to non-empty database",
			store: func(t *testing.T) *HeaderStore {
				t.Helper()
				store := newTestStore(t)
				_, _ = store.AddHeaders(headers[:7])
				return store
			},
			headers:     headers,
			expected:    headers[7:],
			expectedErr: nil,
		},
		{
			name: "add headers reorged headers to non-empty database",
			store: func(t *testing.T) *HeaderStore {
				t.Helper()
				store := newTestStore(t)
				_, err := store.AddHeaders(headers[:7])
				assert.NoError(t, err)
				return store
			},
			headers:     fork,
			expected:    fork[5:],
			expectedErr: nil,
		},
		{
			name: "add headers non-intersecting headers to non-empty database",
			store: func(t *testing.T) *HeaderStore {
				t.Helper()
				store := newTestStore(t)
				_, _ = store.AddHeaders(headers[:5])
				return store
			},
			headers:     headers[6:],
			expected:    nil,
			expectedErr: ErrPrevBlockHashNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store(t)
			got, err := store.AddHeaders(tt.headers)
			assert.Equal(t, tt.expected, got)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestHeaderStore_GetLatestHeader(t *testing.T) {
	tests := []struct {
		name        string
		store       func(t *testing.T) *HeaderStore
		finalized   bool
		expected    *indexer.Header
		expectedErr error
	}{
		{
			name:        "latest header",
			store:       newTestStore,
			finalized:   false,
			expectedErr: nil,
		},
		{
			name:        "latest finalized header",
			store:       newTestStore,
			finalized:   true,
			expectedErr: nil,
		},
		{
			name: "latest finalized header after updating existing headers",
			store: func(t *testing.T) *HeaderStore {
				headers := newTestHeaders(t)
				store := newTestStore(t)
				_, err := store.AddHeaders(headers)
				assert.NoError(t, err)
				return store
			},
			finalized:   true,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store(t)
			defer store.Close()

			headers := newTestHeaders(t)

			for i := 0; i < headers.Len(); i++ {
				if tt.finalized {
					headers[i].Finalized = true
				}
				_, err := store.AddHeaders(headers[0 : i+1])
				assert.NoError(t, err)

				header, err := store.GetLatestHeader(tt.finalized)
				assert.NoError(t, err)
				assert.Equal(t, headers[i], header)
			}
		})
	}
}

func TestHeaderStore_AttachObject(t *testing.T) {
	accum := mockAccumulator{}
	headers := newTestHeaders(t)

	tests := []struct {
		name     string
		store    func(t *testing.T) *HeaderStore
		header   *indexer.Header
		object   mockAccumulatorObjectV1
		expected error
	}{
		{
			name: "attach to existing header",
			store: func(t *testing.T) *HeaderStore {
				store := newTestStore(t)
				headers[0].CurrentFork = "genesis"
				_, err := store.AddHeaders(headers)
				assert.NoError(t, err)
				return store
			},
			header:   headers[0],
			object:   mockAccumulatorObjectV1{Balance: 1000},
			expected: nil,
		},
		{
			name:     "attach to non-existing header",
			store:    newTestStore,
			header:   headers[0],
			object:   mockAccumulatorObjectV1{Balance: 1005},
			expected: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store(t)
			got := store.AttachObject(tt.object, tt.header, accum)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestHeaderStore_GetObject(t *testing.T) {
	headers := newTestHeadersWithFork(t, 1)
	fork := newTestHeadersWithFork(t, 2)

	accum := mockAccumulator{}
	object1 := mockAccumulatorObjectV1{Balance: 1000}
	object2 := mockAccumulatorObjectV2{Balance: 100}

	tests := []struct {
		name           string
		store          func(t *testing.T) *HeaderStore
		header         *indexer.Header
		expectedObject indexer.AccumulatorObject
		expectedHeader *indexer.Header
		expectedErr    error
	}{
		{
			name: "get existing object",
			store: func(t *testing.T) *HeaderStore {
				store := newTestStore(t)

				headers[0].CurrentFork = "genesis"

				_, err := store.AddHeaders(headers)
				assert.NoError(t, err)

				err = store.AttachObject(object1, headers[0], accum)
				assert.NoError(t, err)

				return store
			},
			header:         headers[0],
			expectedObject: object1,
			expectedHeader: headers[0],
			expectedErr:    nil,
		},
		{
			name:           "get non-existing object",
			store:          newTestStore,
			header:         headers[1],
			expectedObject: nil,
			expectedHeader: nil,
			expectedErr:    ErrNotFound,
		},
		{
			name: "get object from prior header",
			store: func(t *testing.T) *HeaderStore {
				store := newTestStore(t)

				headers[4].CurrentFork = "genesis"

				_, err := store.AddHeaders(headers)
				assert.NoError(t, err)

				err = store.AttachObject(object2, headers[7], accum)
				assert.NoError(t, err)
				err = store.AttachObject(object1, headers[4], accum)
				assert.NoError(t, err)

				return store
			},
			header:         headers[5],
			expectedObject: object1,
			expectedHeader: headers[4],
			expectedErr:    nil,
		},
		{
			name: "get object from latest header",
			store: func(t *testing.T) *HeaderStore {
				store := newTestStore(t)

				headers[7].CurrentFork = "exodus"

				_, err := store.AddHeaders(headers)
				assert.NoError(t, err)

				err = store.AttachObject(object2, headers[7], accum)
				assert.NoError(t, err)

				err = store.AttachObject(object1, headers[0], accum)
				assert.NoError(t, err)

				return store
			},
			header:         headers.Last(),
			expectedObject: object2,
			expectedHeader: headers[7],
			expectedErr:    nil,
		},
		{
			name: "get object after reorg",
			store: func(t *testing.T) *HeaderStore {
				store := newTestStore(t)

				_, err := store.AddHeaders(headers[:7])
				assert.NoError(t, err)

				err = store.AttachObject(object1, headers[6], accum)
				assert.NoError(t, err)

				_, err = store.AddHeaders(fork[5:])
				assert.NoError(t, err)

				return store
			},
			header:         headers.Last(),
			expectedObject: nil,
			expectedHeader: nil,
			expectedErr:    ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store(t)
			o, h, err := store.GetObject(tt.header, accum)
			assert.Equal(t, tt.expectedObject, o)
			assert.Equal(t, tt.expectedHeader, h)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
