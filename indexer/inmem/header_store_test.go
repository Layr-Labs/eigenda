package inmem

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/stretchr/testify/assert"
)

type mockAccumulator struct{}

type object struct {
	ID   int
	Name string
}

func (acc mockAccumulator) InitializeObject(header indexer.Header) (indexer.AccumulatorObject, error) {
	return nil, nil
}

func (acc mockAccumulator) UpdateObject(object indexer.AccumulatorObject, header *indexer.Header, event indexer.Event) (indexer.AccumulatorObject, error) {
	return nil, nil
}

func (acc mockAccumulator) SerializeObject(obj indexer.AccumulatorObject, fork indexer.UpgradeFork) ([]byte, error) {
	return encode(obj)
}

func (acc mockAccumulator) DeserializeObject(data []byte, fork indexer.UpgradeFork) (indexer.AccumulatorObject, error) {
	obj := object{}
	err := decode(data, &obj)
	return obj, err
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

	return NewHeaderStore()
}

func newTestHeaders(t *testing.T, fork int) indexer.Headers {
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
	headers := newTestHeaders(t, 1)
	fork := newTestHeaders(t, 2)

	tests := []struct {
		name        string
		store       func(t *testing.T) *HeaderStore
		headers     indexer.Headers
		expected    indexer.Headers
		expectedErr error
	}{
		{
			name:        "add no headers",
			store:       newTestStore,
			headers:     indexer.Headers{},
			expected:    indexer.Headers{},
			expectedErr: nil,
		},
		{
			name:        "add headers to empty store",
			store:       newTestStore,
			headers:     headers,
			expected:    headers,
			expectedErr: nil,
		},
		{
			name: "add headers to the end of non-empty store",
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
			name: "add no new headers",
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
			name: "add intersecting headers to non-empty store",
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
			name: "add reorged headers to non-empty store",
			store: func(t *testing.T) *HeaderStore {
				t.Helper()
				store := newTestStore(t)
				_, _ = store.AddHeaders(headers[:7])
				return store
			},
			headers:     fork,
			expected:    fork[5:],
			expectedErr: nil,
		},
		{
			name: "add non-intersecting headers to non-empty store",
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
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestHeaderStore_GetLatestHeader(t *testing.T) {
	headers := newTestHeaders(t, 1)
	for i := 0; i <= 5; i++ {
		headers[i].Finalized = true
	}

	tests := []struct {
		name        string
		store       func(t *testing.T) *HeaderStore
		headers     indexer.Headers
		finalized   bool
		expected    *indexer.Header
		expectedErr error
	}{
		{
			name: "latest header",
			store: func(t *testing.T) *HeaderStore {
				t.Helper()
				store := newTestStore(t)
				_, err := store.AddHeaders(headers)
				assert.NoError(t, err)
				return store
			},
			headers:     headers,
			finalized:   false,
			expected:    headers[len(headers)-1],
			expectedErr: nil,
		},
		{
			name: "latest finalized header",
			store: func(t *testing.T) *HeaderStore {
				t.Helper()
				store := newTestStore(t)
				_, err := store.AddHeaders(headers)
				assert.NoError(t, err)
				return store
			},
			headers:     headers,
			finalized:   true,
			expected:    headers[5],
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store(t)
			got, err := store.GetLatestHeader(tt.finalized)
			assert.Equal(t, tt.expected, got)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestHeaderStore_AttachObject(t *testing.T) {
	accum := mockAccumulator{}
	headers := newTestHeaders(t, 1)

	type object struct {
		ID   int
		Name string
	}

	tests := []struct {
		name     string
		store    func(t *testing.T) *HeaderStore
		header   *indexer.Header
		object   indexer.AccumulatorObject
		expected error
	}{
		{
			name: "attach to existing header",
			store: func(t *testing.T) *HeaderStore {
				store := newTestStore(t)
				_, err := store.AddHeaders(headers)
				assert.NoError(t, err)
				return store
			},
			header:   headers[0],
			object:   object{ID: 1000, Name: "object-1"},
			expected: nil,
		},
		{
			name:     "attach to non-existing header",
			store:    newTestStore,
			header:   headers[0],
			object:   object{ID: 1001, Name: "object-2"},
			expected: ErrHeaderNotFound,
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

	accum := mockAccumulator{}
	headers := newTestHeaders(t, 1)
	object1 := object{ID: 1000, Name: "object-1"}

	tests := []struct {
		name        string
		store       func(t *testing.T) *HeaderStore
		header      *indexer.Header
		expected    indexer.AccumulatorObject
		expectedErr error
	}{
		{
			name: "get existing object",
			store: func(t *testing.T) *HeaderStore {
				store := newTestStore(t)

				_, err := store.AddHeaders(headers)
				assert.NoError(t, err)
				err = store.AttachObject(object1, headers[0], accum)
				assert.NoError(t, err)

				return store
			},
			header:      headers[0],
			expected:    object1,
			expectedErr: nil,
		},
		{
			name:        "get non-existing object",
			store:       newTestStore,
			header:      headers[1],
			expected:    nil,
			expectedErr: ErrObjectNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got indexer.AccumulatorObject
			store := tt.store(t)

			got, _, err := store.GetObject(tt.header, accum)
			assert.Equal(t, tt.expected, got)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
