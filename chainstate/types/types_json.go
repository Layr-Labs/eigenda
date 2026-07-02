package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/common"
)

// This file defines the JSON representation of the indexed types, used both by
// the HTTP API and by the JSON persister's snapshots. Custom marshalers are
// required because the default encoding of several fields does not match the
// documented API (see the README's endpoint examples):
//
//   - core.OperatorID is a [32]byte and would marshal as a JSON array of 32
//     numbers; the API documents a 0x-prefixed hex string.
//   - []core.QuorumID is a []uint8, which encoding/json marshals as a base64
//     string; the API documents an array of numbers.
//   - *big.Int would marshal as a JSON number, which overflows the 2^53 safe
//     integer range of many JSON consumers for wei-scale stakes; the API
//     documents a decimal string.

// quorumIDsToInts widens []core.QuorumID (an alias of []uint8, which would be
// base64-encoded) to a slice that marshals as a JSON array of numbers.
func quorumIDsToInts(quorumIDs []core.QuorumID) []uint16 {
	out := make([]uint16, len(quorumIDs))
	for i, q := range quorumIDs {
		out[i] = uint16(q)
	}
	return out
}

// intsToQuorumIDs is the inverse of quorumIDsToInts. It rejects values that do
// not fit in a QuorumID rather than silently truncating.
func intsToQuorumIDs(ints []uint16) ([]core.QuorumID, error) {
	out := make([]core.QuorumID, len(ints))
	for i, v := range ints {
		if v > 255 {
			return nil, fmt.Errorf("quorum ID %d out of range", v)
		}
		out[i] = core.QuorumID(v)
	}
	return out, nil
}

// bigIntToString renders a possibly-nil big.Int as a decimal string.
func bigIntToString(v *big.Int) string {
	if v == nil {
		return ""
	}
	return v.String()
}

// stringToBigInt parses a decimal string produced by bigIntToString.
func stringToBigInt(s string) (*big.Int, error) {
	if s == "" {
		return nil, nil
	}
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid big integer %q", s)
	}
	return v, nil
}

type operatorJSON struct {
	ID                        string         `json:"id"`
	Address                   common.Address `json:"address"`
	BLSPubKeyG1               *core.G1Point  `json:"bls_pubkey_g1"`
	BLSPubKeyG2               *core.G2Point  `json:"bls_pubkey_g2"`
	Socket                    string         `json:"socket"`
	RegisteredAtBlockNumber   uint64         `json:"registered_at_block_number"`
	DeregisteredAtBlockNumber *uint64        `json:"deregistered_at_block_number"`
	QuorumIDs                 []uint16       `json:"quorum_ids"`
	RegisteredTxHash          common.Hash    `json:"registered_tx_hash"`
	DeregisteredTxHash        *common.Hash   `json:"deregistered_tx_hash"`
}

// MarshalJSON implements json.Marshaler.
func (o Operator) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(operatorJSON{
		ID:                        "0x" + o.ID.Hex(),
		Address:                   o.Address,
		BLSPubKeyG1:               o.BLSPubKeyG1,
		BLSPubKeyG2:               o.BLSPubKeyG2,
		Socket:                    o.Socket,
		RegisteredAtBlockNumber:   o.RegisteredAtBlockNumber,
		DeregisteredAtBlockNumber: o.DeregisteredAtBlockNumber,
		QuorumIDs:                 quorumIDsToInts(o.QuorumIDs),
		RegisteredTxHash:          o.RegisteredTxHash,
		DeregisteredTxHash:        o.DeregisteredTxHash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operator: %w", err)
	}
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (o *Operator) UnmarshalJSON(data []byte) error {
	var aux operatorJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("failed to unmarshal operator: %w", err)
	}
	id, err := core.OperatorIDFromHex(aux.ID)
	if err != nil {
		return fmt.Errorf("invalid operator ID: %w", err)
	}
	quorumIDs, err := intsToQuorumIDs(aux.QuorumIDs)
	if err != nil {
		return fmt.Errorf("invalid quorum IDs: %w", err)
	}
	*o = Operator{
		ID:                        id,
		Address:                   aux.Address,
		BLSPubKeyG1:               aux.BLSPubKeyG1,
		BLSPubKeyG2:               aux.BLSPubKeyG2,
		Socket:                    aux.Socket,
		RegisteredAtBlockNumber:   aux.RegisteredAtBlockNumber,
		DeregisteredAtBlockNumber: aux.DeregisteredAtBlockNumber,
		QuorumIDs:                 quorumIDs,
		RegisteredTxHash:          aux.RegisteredTxHash,
		DeregisteredTxHash:        aux.DeregisteredTxHash,
	}
	return nil
}

type quorumAPKJSON struct {
	QuorumID    uint16        `json:"quorum_id"`
	BlockNumber uint64        `json:"block_number"`
	APK         *core.G1Point `json:"apk"`
	TotalStake  string        `json:"total_stake"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// MarshalJSON implements json.Marshaler.
func (q QuorumAPK) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(quorumAPKJSON{
		QuorumID:    uint16(q.QuorumID),
		BlockNumber: q.BlockNumber,
		APK:         q.APK,
		TotalStake:  bigIntToString(q.TotalStake),
		UpdatedAt:   q.UpdatedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal quorum APK: %w", err)
	}
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (q *QuorumAPK) UnmarshalJSON(data []byte) error {
	var aux quorumAPKJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("failed to unmarshal quorum APK: %w", err)
	}
	if aux.QuorumID > 255 {
		return fmt.Errorf("quorum ID %d out of range", aux.QuorumID)
	}
	totalStake, err := stringToBigInt(aux.TotalStake)
	if err != nil {
		return fmt.Errorf("invalid total stake: %w", err)
	}
	*q = QuorumAPK{
		QuorumID:    core.QuorumID(aux.QuorumID),
		BlockNumber: aux.BlockNumber,
		APK:         aux.APK,
		TotalStake:  totalStake,
		UpdatedAt:   aux.UpdatedAt,
	}
	return nil
}

type operatorSocketUpdateJSON struct {
	OperatorID  string      `json:"operator_id"`
	Socket      string      `json:"socket"`
	BlockNumber uint64      `json:"block_number"`
	TxHash      common.Hash `json:"tx_hash"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// MarshalJSON implements json.Marshaler.
func (u OperatorSocketUpdate) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(operatorSocketUpdateJSON{
		OperatorID:  "0x" + u.OperatorID.Hex(),
		Socket:      u.Socket,
		BlockNumber: u.BlockNumber,
		TxHash:      u.TxHash,
		UpdatedAt:   u.UpdatedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal socket update: %w", err)
	}
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *OperatorSocketUpdate) UnmarshalJSON(data []byte) error {
	var aux operatorSocketUpdateJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("failed to unmarshal socket update: %w", err)
	}
	id, err := core.OperatorIDFromHex(aux.OperatorID)
	if err != nil {
		return fmt.Errorf("invalid operator ID: %w", err)
	}
	*u = OperatorSocketUpdate{
		OperatorID:  id,
		Socket:      aux.Socket,
		BlockNumber: aux.BlockNumber,
		TxHash:      aux.TxHash,
		UpdatedAt:   aux.UpdatedAt,
	}
	return nil
}

type operatorEjectionJSON struct {
	OperatorID  string      `json:"operator_id"`
	QuorumIDs   []uint16    `json:"quorum_ids"`
	BlockNumber uint64      `json:"block_number"`
	TxHash      common.Hash `json:"tx_hash"`
	EjectedAt   time.Time   `json:"ejected_at"`
}

// MarshalJSON implements json.Marshaler.
func (e OperatorEjection) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(operatorEjectionJSON{
		OperatorID:  "0x" + e.OperatorID.Hex(),
		QuorumIDs:   quorumIDsToInts(e.QuorumIDs),
		BlockNumber: e.BlockNumber,
		TxHash:      e.TxHash,
		EjectedAt:   e.EjectedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ejection: %w", err)
	}
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *OperatorEjection) UnmarshalJSON(data []byte) error {
	var aux operatorEjectionJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("failed to unmarshal ejection: %w", err)
	}
	id, err := core.OperatorIDFromHex(aux.OperatorID)
	if err != nil {
		return fmt.Errorf("invalid operator ID: %w", err)
	}
	quorumIDs, err := intsToQuorumIDs(aux.QuorumIDs)
	if err != nil {
		return fmt.Errorf("invalid quorum IDs: %w", err)
	}
	*e = OperatorEjection{
		OperatorID:  id,
		QuorumIDs:   quorumIDs,
		BlockNumber: aux.BlockNumber,
		TxHash:      aux.TxHash,
		EjectedAt:   aux.EjectedAt,
	}
	return nil
}
