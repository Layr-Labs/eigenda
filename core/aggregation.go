package core

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sort"

	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	lru "github.com/hashicorp/golang-lru/v2"
)

const maxNumOperatorAddresses = 300

var (
	ErrPubKeysNotEqual     = errors.New("public keys are not equal")
	ErrInsufficientEthSigs = errors.New("insufficient eth signatures")
	ErrAggSigNotValid      = errors.New("aggregated signature is not valid")
)

type SigningMessage struct {
	Signature       *Signature
	Operator        OperatorID
	BatchHeaderHash [32]byte
	Err             error
}

// SignatureAggregation contains the results of aggregating signatures from a set of operators
type SignatureAggregation struct {
	// NonSigners contains the public keys of the operators that did not sign the message
	NonSigners []*G1Point
	// QuorumAggPubKeys contains the aggregated public keys for all of the operators each quorum,
	// Including those that did not sign
	QuorumAggPubKeys []*G1Point
	// AggPubKey is the aggregated public key for all of the operators that signed the message,
	// further aggregated across the quorums; operators signing for multiple quorums will be included in
	// the aggregation multiple times
	AggPubKey *G2Point
	// AggSignature is the aggregated signature for all of the operators that signed the message, mirroring the
	// AggPubKey.
	AggSignature *Signature
	// QuorumResults contains the quorum ID and the amount signed for each quorum
	QuorumResults map[QuorumID]*QuorumResult
	// SignerMap contains the operator IDs that signed the message
	SignerMap map[OperatorID]bool
}

// SignatureAggregator is an interface for aggregating the signatures returned by DA nodes so that they can be verified by the DA contract
type SignatureAggregator interface {

	// AggregateSignatures blocks until it receives a response for each operator in the operator state via messageChan, and then returns the aggregated signature.
	// If the aggregated signature is invalid, an error is returned.
	AggregateSignatures(ctx context.Context, state *IndexedOperatorState, quorumIDs []QuorumID, message [32]byte, messageChan chan SigningMessage) (*SignatureAggregation, error)
}

type StdSignatureAggregator struct {
	Logger     logging.Logger
	Transactor Transactor
	// OperatorAddresses contains the ethereum addresses of the operators corresponding to their operator IDs
	OperatorAddresses *lru.Cache[OperatorID, gethcommon.Address]
}

func NewStdSignatureAggregator(logger logging.Logger, transactor Transactor) (*StdSignatureAggregator, error) {
	operatorAddrs, err := lru.New[OperatorID, gethcommon.Address](maxNumOperatorAddresses)
	if err != nil {
		return nil, err
	}

	return &StdSignatureAggregator{
		Logger:            logger.With("component", "SignatureAggregator"),
		Transactor:        transactor,
		OperatorAddresses: operatorAddrs,
	}, nil
}

var _ SignatureAggregator = (*StdSignatureAggregator)(nil)

func (a *StdSignatureAggregator) AggregateSignatures(ctx context.Context, state *IndexedOperatorState, quorumIDs []QuorumID, message [32]byte, messageChan chan SigningMessage) (*SignatureAggregation, error) {
	// TODO: Add logging

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

	stakeSigned := make([]*big.Int, len(quorumIDs))
	for ind := range quorumIDs {
		stakeSigned[ind] = big.NewInt(0)
	}
	aggSigs := make([]*Signature, len(quorumIDs))
	aggPubKeys := make([]*G2Point, len(quorumIDs))

	signerMap := make(map[OperatorID]bool)

	// Aggregate Signatures
	numOperators := len(state.IndexedOperators)

	for numReply := 0; numReply < numOperators; numReply++ {
		var err error
		r := <-messageChan
		operatorIDHex := r.Operator.Hex()
		operatorAddr, ok := a.OperatorAddresses.Get(r.Operator)
		if !ok && a.Transactor != nil {
			operatorAddr, err = a.Transactor.OperatorIDToAddress(ctx, r.Operator)
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
			a.Logger.Warn("error returned from messageChan", "operatorID", operatorIDHex, "operatorAddress", operatorAddr, "socket", socket, "batchHeaderHash", batchHeaderHashHex, "err", r.Err)
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
		for ind, quorumID := range quorumIDs {
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
			stakeSigned[ind].Add(stakeSigned[ind], opInfo.Stake)

			// Add to agg signature
			if aggSigs[ind] == nil {
				aggSigs[ind] = &Signature{sig.Clone()}
				aggPubKeys[ind] = op.PubkeyG2.Clone()
			} else {
				aggSigs[ind].Add(sig.G1Point)
				aggPubKeys[ind].Add(op.PubkeyG2)
			}
		}
		a.Logger.Info("received signature from operator", "operatorID", operatorIDHex, "operatorAddress", operatorAddr, "socket", socket, "quorumIDs", fmt.Sprint(operatorQuorums), "batchHeaderHash", batchHeaderHashHex)
	}

	// Aggregate Non signer Pubkey Id
	nonSignerKeys := make([]*G1Point, 0)
	nonSignerOperatorIds := make([]OperatorID, 0)

	for id, op := range state.IndexedOperators {
		_, found := signerMap[id]
		if !found {
			nonSignerKeys = append(nonSignerKeys, op.PubkeyG1)
			nonSignerOperatorIds = append(nonSignerOperatorIds, id)
		}
	}

	quorumAggPubKeys := make([]*G1Point, len(quorumIDs))

	// Validate the amount signed and aggregate signatures for each quorum
	quorumResults := make(map[QuorumID]*QuorumResult)

	for ind, quorumID := range quorumIDs {
		// Check that quorum has sufficient stake
		percent := GetSignedPercentage(state.OperatorState, quorumID, stakeSigned[ind])
		quorumResults[quorumID] = &QuorumResult{
			QuorumID:      quorumID,
			PercentSigned: percent,
		}

		// Verify that the aggregated public key for the quorum matches the on-chain quorum aggregate public key sans non-signers of the quorum
		quorumAggKey := state.AggKeys[quorumID]
		quorumAggPubKeys[ind] = quorumAggKey

		signersAggKey := quorumAggKey.Clone()
		for opInd, nsk := range nonSignerKeys {
			ops := state.Operators[quorumID]
			if _, ok := ops[nonSignerOperatorIds[opInd]]; ok {
				signersAggKey.Sub(nsk)
			}
		}

		if aggPubKeys[ind] == nil {
			return nil, ErrAggSigNotValid
		}

		ok, err := signersAggKey.VerifyEquivalence(aggPubKeys[ind])
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrPubKeysNotEqual
		}

		// Verify the aggregated signature for the quorum
		ok = aggSigs[ind].Verify(aggPubKeys[ind], message)
		if !ok {
			return nil, ErrAggSigNotValid
		}
	}

	// Aggregate the aggregated signatures. We reuse the first aggregated signature as the accumulator
	for i := 1; i < len(aggSigs); i++ {
		aggSigs[0].Add(aggSigs[i].G1Point)
	}

	// Aggregate the aggregated public keys. We reuse the first aggregated public key as the accumulator
	for i := 1; i < len(aggPubKeys); i++ {
		aggPubKeys[0].Add(aggPubKeys[i])
	}

	// sort non signer keys according to how it's checked onchain
	// ref: https://github.com/Layr-Labs/eigenlayer-middleware/blob/m2-mainnet/src/BLSSignatureChecker.sol#L99
	sort.Slice(nonSignerKeys, func(i, j int) bool {
		hash1 := nonSignerKeys[i].Hash()
		hash2 := nonSignerKeys[j].Hash()
		// sort in accending order
		return bytes.Compare(hash1[:], hash2[:]) == -1
	})

	return &SignatureAggregation{
		NonSigners:       nonSignerKeys,
		QuorumAggPubKeys: quorumAggPubKeys,
		AggPubKey:        aggPubKeys[0],
		AggSignature:     aggSigs[0],
		QuorumResults:    quorumResults,
		SignerMap:        signerMap,
	}, nil

}

func GetStakeThreshold(state *OperatorState, quorum QuorumID, quorumThreshold uint8) *big.Int {

	// Get stake threshold
	quorumThresholdBig := new(big.Int).SetUint64(uint64(quorumThreshold))
	stakeThreshold := new(big.Int)
	stakeThreshold.Mul(quorumThresholdBig, state.Totals[quorum].Stake)
	stakeThreshold = roundUpDivideBig(stakeThreshold, new(big.Int).SetUint64(percentMultiplier))

	return stakeThreshold
}

func GetSignedPercentage(state *OperatorState, quorum QuorumID, stakeAmount *big.Int) uint8 {

	stakeAmount = stakeAmount.Mul(stakeAmount, new(big.Int).SetUint64(percentMultiplier))
	quorumThresholdBig := stakeAmount.Div(stakeAmount, state.Totals[quorum].Stake)

	quorumThreshold := uint8(quorumThresholdBig.Uint64())

	return quorumThreshold
}
