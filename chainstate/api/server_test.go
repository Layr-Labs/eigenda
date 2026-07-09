package api

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Layr-Labs/eigenda/chainstate"
	"github.com/Layr-Labs/eigenda/chainstate/store"
	"github.com/Layr-Labs/eigenda/chainstate/types"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, st store.Store) *Server {
	t.Helper()
	return NewServer(chainstate.DefaultIndexerConfig(), st, test.GetLogger())
}

func get(t *testing.T, s *Server, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	return rec
}

func TestQuorumAPKRejectsBadQuorumID(t *testing.T) {
	ctx := context.Background()
	st := store.NewMemoryStore()
	require.NoError(t, st.SaveQuorumAPK(ctx, &types.QuorumAPK{QuorumID: 0, BlockNumber: 100}))
	s := newTestServer(t, st)

	for _, endpoint := range []string{"/api/v1/quorum-apk", "/api/v1/quorum-apk/history"} {
		// Missing, non-numeric, and out-of-uint8-range values must all be 400,
		// never silently served as quorum 0.
		badQueries := []string{
			"block_number=100",
			"quorum_id=abc&block_number=100",
			"quorum_id=256&block_number=100",
		}
		for _, query := range badQueries {
			rec := get(t, s, endpoint+"?"+query)
			require.Equal(t, http.StatusBadRequest, rec.Code, "%s?%s", endpoint, query)
		}
	}

	// A valid quorum_id works.
	rec := get(t, s, "/api/v1/quorum-apk?quorum_id=0&block_number=100")
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestGetOperatorStatusCodes(t *testing.T) {
	ctx := context.Background()
	st := store.NewMemoryStore()
	id := core.OperatorID{1}
	require.NoError(t, st.SaveOperator(ctx, &types.Operator{ID: id, Socket: "host:1"}))
	s := newTestServer(t, st)

	rec := get(t, s, "/api/v1/operators/"+id.Hex())
	require.Equal(t, http.StatusOK, rec.Code)

	unknown := core.OperatorID{2}
	rec = get(t, s, "/api/v1/operators/"+unknown.Hex())
	require.Equal(t, http.StatusNotFound, rec.Code)

	rec = get(t, s, "/api/v1/operators/not-hex")
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// failingStore wraps a Store, forcing GetOperator to return a non-not-found
// error, to verify the handler distinguishes internal failures from absence.
type failingStore struct {
	store.Store
}

func (f *failingStore) GetOperator(ctx context.Context, id core.OperatorID) (*types.Operator, error) {
	return nil, errors.New("backend unavailable")
}

func TestGetOperatorInternalErrorIsNot404(t *testing.T) {
	s := newTestServer(t, &failingStore{Store: store.NewMemoryStore()})

	rec := get(t, s, "/api/v1/operators/"+core.OperatorID{1}.Hex())
	require.Equal(t, http.StatusInternalServerError, rec.Code,
		"a store failure must not masquerade as operator-not-found")
}

// TestOperatorJSONFieldNames pins the documented wire format (see
// chainstate/README.md). If this test fails, either the marshaler or the
// README must change — together.
func TestOperatorJSONFieldNames(t *testing.T) {
	ctx := context.Background()
	st := store.NewMemoryStore()

	keyPair, err := core.GenRandomBlsKeys()
	require.NoError(t, err)
	id := keyPair.GetPubKeyG1().GetOperatorID()

	require.NoError(t, st.SaveOperator(ctx, &types.Operator{
		ID:                      id,
		Address:                 common.HexToAddress("0x1234"),
		BLSPubKeyG1:             keyPair.GetPubKeyG1(),
		BLSPubKeyG2:             keyPair.GetPubKeyG2(),
		Socket:                  "host:1",
		RegisteredAtBlockNumber: 50,
		QuorumIDs:               []core.QuorumID{0, 1},
	}))
	s := newTestServer(t, st)

	rec := get(t, s, "/api/v1/operators/"+id.Hex())
	require.Equal(t, http.StatusOK, rec.Code)

	var body map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	for _, field := range []string{
		"id", "address", "bls_pubkey_g1", "bls_pubkey_g2", "socket",
		"registered_at_block_number", "deregistered_at_block_number",
		"quorum_ids", "registered_tx_hash", "deregistered_tx_hash",
	} {
		require.Contains(t, body, field)
	}

	// quorum_ids must be a JSON array of numbers, not a base64 string.
	var quorums []uint16
	require.NoError(t, json.Unmarshal(body["quorum_ids"], &quorums))
	require.Equal(t, []uint16{0, 1}, quorums)

	// id must be a 0x-prefixed hex string.
	var idStr string
	require.NoError(t, json.Unmarshal(body["id"], &idStr))
	require.Equal(t, "0x"+id.Hex(), idStr)
}

func TestEjectionAndSocketUpdateJSONFieldNames(t *testing.T) {
	ctx := context.Background()
	st := store.NewMemoryStore()
	id := core.OperatorID{1}

	require.NoError(t, st.SaveEjection(ctx, &types.OperatorEjection{
		OperatorID:  id,
		QuorumIDs:   []core.QuorumID{0},
		BlockNumber: 100,
		LogIndex:    3,
	}))
	require.NoError(t, st.SaveSocketUpdate(ctx, &types.OperatorSocketUpdate{
		OperatorID:  id,
		Socket:      "host:1",
		BlockNumber: 100,
		LogIndex:    2,
	}))
	s := newTestServer(t, st)

	rec := get(t, s, "/api/v1/ejections/"+id.Hex())
	require.Equal(t, http.StatusOK, rec.Code)
	var ejBody struct {
		Ejections []map[string]json.RawMessage `json:"ejections"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &ejBody))
	require.Len(t, ejBody.Ejections, 1)
	for _, field := range []string{"operator_id", "quorum_ids", "block_number", "log_index", "tx_hash", "ejected_at"} {
		require.Contains(t, ejBody.Ejections[0], field)
	}

	rec = get(t, s, "/api/v1/socket-updates/"+id.Hex())
	require.Equal(t, http.StatusOK, rec.Code)
	var suBody struct {
		SocketUpdates []map[string]json.RawMessage `json:"socket_updates"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &suBody))
	require.Len(t, suBody.SocketUpdates, 1)
	for _, field := range []string{"operator_id", "socket", "block_number", "log_index", "tx_hash", "updated_at"} {
		require.Contains(t, suBody.SocketUpdates[0], field)
	}
}

func TestQuorumAPKJSONFieldNames(t *testing.T) {
	ctx := context.Background()
	st := store.NewMemoryStore()
	require.NoError(t, st.SaveQuorumAPK(ctx, &types.QuorumAPK{
		QuorumID:    1,
		BlockNumber: 100,
		TotalStake:  bigFromString(t, "1000000000000000000000000"),
	}))
	s := newTestServer(t, st)

	rec := get(t, s, "/api/v1/quorum-apk?quorum_id=1&block_number=100")
	require.Equal(t, http.StatusOK, rec.Code)

	var body map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	for _, field := range []string{"quorum_id", "block_number", "apk", "total_stake", "updated_at"} {
		require.Contains(t, body, field)
	}

	// total_stake must be a decimal string: wei-scale values overflow the
	// 2^53 safe-integer range of many JSON consumers.
	var stake string
	require.NoError(t, json.Unmarshal(body["total_stake"], &stake))
	require.Equal(t, "1000000000000000000000000", stake)
}

func TestListOperatorsInvalidQuorumID(t *testing.T) {
	s := newTestServer(t, store.NewMemoryStore())

	rec := get(t, s, "/api/v1/operators?quorum_id=abc")
	require.Equal(t, http.StatusBadRequest, rec.Code)

	// quorum_id is optional here: omitting it lists all operators.
	rec = get(t, s, "/api/v1/operators")
	require.Equal(t, http.StatusOK, rec.Code)
}

func bigFromString(t *testing.T, s string) *big.Int {
	t.Helper()
	v, ok := new(big.Int).SetString(s, 10)
	require.True(t, ok)
	return v
}
