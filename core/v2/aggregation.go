package corev2

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"sort"

	"github.com/Layr-Labs/eigenda/chainio"
	"github.com/Layr-Labs/eigenda/crypto/ecc/bn254"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	lru "github.com/hashicorp/golang-lru/v2"
)

const percentMultiplier = 100

const maxNumOperatorAddresses = 300

var (
	ErrPubKeysNotEqual     = errors.New("public keys are not equal")
	ErrInsufficientEthSigs = errors.New("insufficient eth signatures")
	ErrAggPubKeyNotValid   = errors.New("aggregated public key is not valid")
	ErrAggSigNotValid      = errors.New("aggregated signature is not valid")
)

type SigningMessage struct {
	Signature       *bn254.Signature
	Operator        OperatorID
	BatchHeaderHash [32]byte
	// Undefined if this value <= 0.
	AttestationLatencyMs float64
	Err                  error
}

// QuorumResult contains the quorum ID and the amount signed for the quorum
type QuorumResult struct {
	QuorumID QuorumID
	// PercentSigned is percentage of the total stake for the quorum that signed for a particular batch.
	PercentSigned uint8
}

// QuorumAttestation contains the results of aggregating signatures from a set of operators by quorums
// It also returns map of all signers across all quorums
type QuorumAttestation struct {
	// QuorumAggPubKeys contains the aggregated public keys for all of the operators each quorum,
	// including those that did not sign
	QuorumAggPubKey map[QuorumID]*bn254.G1Point
	// SignersAggPubKey is the aggregated public key for all of the operators that signed the message by each quorum
	SignersAggPubKey map[QuorumID]*bn254.G2Point
	// AggSignature is the aggregated signature for all of the operators that signed the message for each quorum, mirroring the
	// SignersAggPubKey.
	AggSignature map[QuorumID]*bn254.Signature
	// QuorumResults contains the quorum ID and the amount signed for each quorum
	QuorumResults map[QuorumID]*QuorumResult
	// SignerMap contains the operator IDs that signed the message
	SignerMap map[OperatorID]bool
}

// SignatureAggregation contains the results of aggregating signatures from a set of operators across multiple quorums
type SignatureAggregation struct {
	// NonSigners contains the public keys of the operators that did not sign the message
	NonSigners []*bn254.G1Point
	// QuorumAggPubKeys contains the aggregated public keys for all of the operators each quorum,
	// Including those that did not sign
	QuorumAggPubKeys map[QuorumID]*bn254.G1Point
	// AggPubKey is the aggregated public key for all of the operators that signed the message,
	// further aggregated across the quorums; operators signing for multiple quorums will be included in
	// the aggregation multiple times
	AggPubKey *bn254.G2Point
	// AggSignature is the aggregated signature for all of the operators that signed the message, mirroring the
	// AggPubKey.
	AggSignature *bn254.Signature
	// QuorumResults contains the quorum ID and the amount signed for each quorum
	QuorumResults map[QuorumID]*QuorumResult
}

// SignatureAggregator is an interface for aggregating the signatures returned by DA nodes so that they can be verified by the DA contract
type SignatureAggregator interface {
	// ReceiveSignatures blocks until it receives a response for each operator in the operator state via messageChan, and then returns the attestation result by quorum.
	ReceiveSignatures(ctx context.Context, state *chainio.IndexedOperatorState, message [32]byte, messageChan chan SigningMessage) (*QuorumAttestation, error)
	// AggregateSignatures takes attestation result by quorum and aggregates the signatures across them.
	// If the aggregated signature is invalid, an error is returned.
	AggregateSignatures(ctx context.Context, ics chainio.IndexedChainState, referenceBlockNumber uint, quorumAttestation *QuorumAttestation, quorumIDs []QuorumID) (*SignatureAggregation, error)
}

type StdSignatureAggregator struct {
	Logger      logging.Logger
	ChainReader chainio.Reader
	// OperatorAddresses contains the ethereum addresses of the operators corresponding to their operator IDs
	OperatorAddresses *lru.Cache[OperatorID, gethcommon.Address]
}

func NewStdSignatureAggregator(logger logging.Logger, reader chainio.Reader) (*StdSignatureAggregator, error) {
	operatorAddrs, err := lru.New[OperatorID, gethcommon.Address](maxNumOperatorAddresses)
	if err != nil {
		return nil, err
	}

	return &StdSignatureAggregator{
		Logger:            logger.With("component", "SignatureAggregator"),
		ChainReader:       reader,
		OperatorAddresses: operatorAddrs,
	}, nil
}

var _ SignatureAggregator = (*StdSignatureAggregator)(nil)

func (a *StdSignatureAggregator) ReceiveSignatures(ctx context.Context, state *chainio.IndexedOperatorState, message [32]byte, messageChan chan SigningMessage) (*QuorumAttestation, error) {
	quorumIDs := make([]QuorumID, 0, len(state.AggKeys))
	for quorumID := range state.Operators {
		quorumIDs = append(quorumIDs, quorumID)
	}
	slices.Sort(quorumIDs)

	if len(quorumIDs) == 0 {
		return nil, errors.New("the number of quorums must be greater than zero")
	}

	// Ensure all quorums are found in state
	for _, id := range quorumIDs {
		_, found := state.Operators[id]
		if !found {
			return nil, errors.New("quorum not found")
		}
	}

	stakeSigned := make(map[QuorumID]*big.Int, len(quorumIDs))
	for _, quorumID := range quorumIDs {
		stakeSigned[quorumID] = big.NewInt(0)
	}
	aggSigs := make(map[QuorumID]*bn254.Signature, len(quorumIDs))
	aggPubKeys := make(map[QuorumID]*bn254.G2Point, len(quorumIDs))
	signerMap := make(map[OperatorID]bool)

	// Aggregate Signatures
	numOperators := len(state.IndexedOperators)

	for numReply := 0; numReply < numOperators; numReply++ {
		var err error
		r := <-messageChan
		operatorIDHex := chainio.GetOperatorHex(r.Operator)
		operatorAddr, ok := a.OperatorAddresses.Get(r.Operator)
		if !ok && a.ChainReader != nil {
			operatorAddr, err = a.ChainReader.OperatorIDToAddress(ctx, r.Operator)
			if err != nil {
				a.Logger.Error("failed to get operator address from registry", "operatorID", operatorIDHex)
				operatorAddr = gethcommon.Address{}
			} else {
				a.OperatorAddresses.Add(r.Operator, operatorAddr)
			}
		} else if !ok {
			operatorAddr = gethcommon.Address{}
		}

		socket := ""
		if op, ok := state.IndexedOperators[r.Operator]; ok {
			socket = op.Socket
		}
		batchHeaderHashHex := hex.EncodeToString(r.BatchHeaderHash[:])
		if r.Err != nil {
			a.Logger.Warn("error returned from messageChan", "operatorID", operatorIDHex, "operatorAddress", operatorAddr, "socket", socket, "batchHeaderHash", batchHeaderHashHex, "attestationLatencyMs", r.AttestationLatencyMs, "err", r.Err)
			continue
		}

		op, found := state.IndexedOperators[r.Operator]
		if !found {
			a.Logger.Error("Operator not found in state", "operatorID", operatorIDHex, "operatorAddress", operatorAddr, "socket", socket)
			continue
		}

		// Verify Signature
		sig := r.Signature
		ok = sig.Verify(op.PubkeyG2, message)
		if !ok {
			a.Logger.Error("signature is not valid", "operatorID", operatorIDHex, "operatorAddress", operatorAddr, "socket", socket, "pubkey", hexutil.Encode(op.PubkeyG2.Serialize()))
			continue
		}

		operatorQuorums := make([]uint8, 0, len(quorumIDs))
		for _, quorumID := range quorumIDs {
			// Get stake amounts for operator
			ops := state.Operators[quorumID]
			opInfo, ok := ops[r.Operator]
			// If operator is not in quorum, skip
			if !ok {
				continue
			}
			operatorQuorums = append(operatorQuorums, quorumID)

			signerMap[r.Operator] = true

			// Add to stake signed
			stakeSigned[quorumID].Add(stakeSigned[quorumID], opInfo.Stake)

			// Add to agg signature
			if aggSigs[quorumID] == nil {
				aggSigs[quorumID] = &bn254.Signature{sig.Clone()}
				aggPubKeys[quorumID] = op.PubkeyG2.Clone()
			} else {
				aggSigs[quorumID].Add(sig.G1Point)
				aggPubKeys[quorumID].Add(op.PubkeyG2)
			}
		}
		a.Logger.Info("received signature from operator", "operatorID", operatorIDHex, "operatorAddress", operatorAddr, "socket", socket, "quorumIDs", fmt.Sprint(operatorQuorums), "batchHeaderHash", batchHeaderHashHex, "attestationLatencyMs", r.AttestationLatencyMs)
	}

	// Aggregate Non signer Pubkey Id
	nonSignerKeys := make([]*bn254.G1Point, 0)
	nonSignerOperatorIds := make([]OperatorID, 0)

	for id, op := range state.IndexedOperators {
		_, found := signerMap[id]
		if !found {
			nonSignerKeys = append(nonSignerKeys, op.PubkeyG1)
			nonSignerOperatorIds = append(nonSignerOperatorIds, id)
		}
	}

	quorumAggPubKeys := make(map[QuorumID]*bn254.G1Point, len(quorumIDs))

	// Validate the amount signed and aggregate signatures for each quorum
	quorumResults := make(map[QuorumID]*QuorumResult)

	for _, quorumID := range quorumIDs {
		// Check that quorum has sufficient stake
		percent := GetSignedPercentage(state.OperatorState, quorumID, stakeSigned[quorumID])
		quorumResults[quorumID] = &QuorumResult{
			QuorumID:      quorumID,
			PercentSigned: percent,
		}

		if percent == 0 {
			a.Logger.Warn("no stake signed for quorum", "quorumID", quorumID)
			continue
		}

		// Verify that the aggregated public key for the quorum matches the on-chain quorum aggregate public key sans non-signers of the quorum
		quorumAggKey := state.AggKeys[quorumID]
		quorumAggPubKeys[quorumID] = quorumAggKey

		signersAggKey := quorumAggKey.Clone()
		for opInd, nsk := range nonSignerKeys {
			ops := state.Operators[quorumID]
			if _, ok := ops[nonSignerOperatorIds[opInd]]; ok {
				signersAggKey.Sub(nsk)
			}
		}

		if aggPubKeys[quorumID] == nil {
			return nil, ErrAggPubKeyNotValid
		}

		ok, err := signersAggKey.VerifyEquivalence(aggPubKeys[quorumID])
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrPubKeysNotEqual
		}

		// Verify the aggregated signature for the quorum
		ok = aggSigs[quorumID].Verify(aggPubKeys[quorumID], message)
		if !ok {
			return nil, ErrAggSigNotValid
		}
	}

	return &QuorumAttestation{
		QuorumAggPubKey:  quorumAggPubKeys,
		SignersAggPubKey: aggPubKeys,
		AggSignature:     aggSigs,
		QuorumResults:    quorumResults,
		SignerMap:        signerMap,
	}, nil
}

func (a *StdSignatureAggregator) AggregateSignatures(ctx context.Context, ics chainio.IndexedChainState, referenceBlockNumber uint, quorumAttestation *QuorumAttestation, quorumIDs []QuorumID) (*SignatureAggregation, error) {
	// Aggregate the aggregated signatures. We reuse the first aggregated signature as the accumulator
	var aggSig *bn254.Signature
	for _, quorumID := range quorumIDs {
		sig := quorumAttestation.AggSignature[quorumID]
		if aggSig == nil {
			aggSig = &bn254.Signature{sig.G1Point.Clone()}
		} else {
			aggSig.Add(sig.G1Point)
		}
	}

	// Aggregate the aggregated public keys. We reuse the first aggregated public key as the accumulator
	var aggPubKey *bn254.G2Point
	for _, quorumID := range quorumIDs {
		apk := quorumAttestation.SignersAggPubKey[quorumID]
		if aggPubKey == nil {
			aggPubKey = apk.Clone()
		} else {
			aggPubKey.Add(apk)
		}
	}

	nonSignerKeys := make([]*bn254.G1Point, 0)
	indexedOperatorState, err := ics.GetIndexedOperatorState(ctx, referenceBlockNumber, quorumIDs)
	if err != nil {
		return nil, err
	}
	for id, op := range indexedOperatorState.IndexedOperators {
		_, found := quorumAttestation.SignerMap[id]
		if !found {
			nonSignerKeys = append(nonSignerKeys, op.PubkeyG1)
		}
	}

	// sort non signer keys according to how it's checked onchain
	// ref: https://github.com/Layr-Labs/eigenlayer-middleware/blob/m2-mainnet/src/BLSSignatureChecker.sol#L99
	sort.Slice(nonSignerKeys, func(i, j int) bool {
		hash1 := nonSignerKeys[i].Hash()
		hash2 := nonSignerKeys[j].Hash()
		// sort in accending order
		return bytes.Compare(hash1[:], hash2[:]) == -1
	})

	quorumAggKeys := make(map[QuorumID]*bn254.G1Point, len(quorumIDs))
	for _, quorumID := range quorumIDs {
		quorumAggKeys[quorumID] = quorumAttestation.QuorumAggPubKey[quorumID]
	}

	quorumResults := make(map[QuorumID]*QuorumResult, len(quorumIDs))
	for _, quorumID := range quorumIDs {
		quorumResults[quorumID] = quorumAttestation.QuorumResults[quorumID]
	}

	return &SignatureAggregation{
		NonSigners:       nonSignerKeys,
		QuorumAggPubKeys: quorumAggKeys,
		AggPubKey:        aggPubKey,
		AggSignature:     aggSig,
		QuorumResults:    quorumResults,
	}, nil

}

func GetStakeThreshold(state *chainio.OperatorState, quorum QuorumID, quorumThreshold uint8) *big.Int {

	// Get stake threshold
	quorumThresholdBig := new(big.Int).SetUint64(uint64(quorumThreshold))
	stakeThreshold := new(big.Int)
	stakeThreshold.Mul(quorumThresholdBig, state.Totals[quorum].Stake)
	stakeThreshold = RoundUpDivideBig(stakeThreshold, new(big.Int).SetUint64(percentMultiplier))

	return stakeThreshold
}

func GetSignedPercentage(state *chainio.OperatorState, quorum QuorumID, stakeAmount *big.Int) uint8 {

	stakeAmount = stakeAmount.Mul(stakeAmount, new(big.Int).SetUint64(percentMultiplier))
	quorumThresholdBig := stakeAmount.Div(stakeAmount, state.Totals[quorum].Stake)

	quorumThreshold := uint8(quorumThresholdBig.Uint64())

	return quorumThreshold
}
