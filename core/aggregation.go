package core

import (
	"bytes"
	"errors"
	"math/big"
	"sort"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	ErrPubKeysNotEqual     = errors.New("public keys are not equal")
	ErrInsufficientEthSigs = errors.New("insufficient eth signatures")
	ErrAggSigNotValid      = errors.New("aggregated signature is not valid")
)

type SignerMessage struct {
	Signature *Signature
	Operator  OperatorID
	Err       error
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
}

// SignatureAggregator is an interface for aggregating the signatures returned by DA nodes so that they can be verified by the DA contract
type SignatureAggregator interface {

	// AggregateSignatures blocks until it receives a response for each operator in the operator state via messageChan, and then returns the aggregated signature.
	// If the aggregated signature is invalid, an error is returned.
	AggregateSignatures(state *IndexedOperatorState, quorumIDs []QuorumID, message [32]byte, messageChan chan SignerMessage) (*SignatureAggregation, error)
}

type StdSignatureAggregator struct {
	Logger common.Logger
}

func NewStdSignatureAggregator(logger common.Logger) *StdSignatureAggregator {
	return &StdSignatureAggregator{
		Logger: logger,
	}
}

var _ SignatureAggregator = (*StdSignatureAggregator)(nil)

func (a *StdSignatureAggregator) AggregateSignatures(state *IndexedOperatorState, quorumIDs []QuorumID, message [32]byte, messageChan chan SignerMessage) (*SignatureAggregation, error) {

	// TODO: Add logging

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
		r := <-messageChan
		operatorIDHex := hexutil.Encode(r.Operator[:])
		socket := ""
		if op, ok := state.IndexedOperators[r.Operator]; ok {
			socket = op.Socket
		}
		if r.Err != nil {
			a.Logger.Warn("[AggregateSignatures] error returned from messageChan", "operator", operatorIDHex, "socket", socket, "err", r.Err)
			continue
		}

		op, found := state.IndexedOperators[r.Operator]
		if !found {
			a.Logger.Error("Operator not found in state", "operator", operatorIDHex, "socket", socket)
			continue
		}

		// Verify Signature
		sig := r.Signature
		ok := sig.Verify(op.PubkeyG2, message)
		if !ok {
			a.Logger.Error("Signature is not valid", "operator", operatorIDHex, "socket", socket, "pubkey", hexutil.Encode(op.PubkeyG2.Serialize()))
			continue
		}

		a.Logger.Info("[AggregateSignatures] received signature from operator", "operator", operatorIDHex, "socket", socket)

		for ind, id := range quorumIDs {

			// Get stake amounts for operator
			ops := state.Operators[id]
			opInfo, ok := ops[r.Operator]

			// If operator is not in quorum, skip
			if !ok {
				a.Logger.Error("Operator not found in quorum", "operator", operatorIDHex, "socket", socket)
				continue
			}

			signerMap[r.Operator] = true

			// Add to stake signed
			stakeSigned[ind].Add(stakeSigned[ind], opInfo.Stake)

			// Add to agg signature
			if aggSigs[ind] == nil {
				aggSigs[ind] = &Signature{sig.Deserialize(sig.Serialize())}
				aggPubKeys[ind] = op.PubkeyG2.Deserialize(op.PubkeyG2.Serialize())
			} else {
				aggSigs[ind].Add(sig.G1Point)
				aggPubKeys[ind].Add(op.PubkeyG2)
			}
		}
	}

	// Aggregate Non signer Pubkey Id
	nonSignerKeys := make([]*G1Point, 0)
	nonSignerOperatorIds := make([]OperatorID, 0)

	for id, op := range state.IndexedOperators {
		_, found := signerMap[id]
		a.Logger.Trace("[state.IndexedOperators]", "operator", hexutil.Encode(id[:]), "G1X", op.PubkeyG1.X.Text(16), "G1Y", op.PubkeyG1.Y.Text(16))
		if !found {
			nonSignerKeys = append(nonSignerKeys, op.PubkeyG1)
			a.Logger.Trace("[state.IndexedOperators] Non signer found", "operator", hexutil.Encode(id[:]), "G1X", op.PubkeyG1.X.Text(16), "G1Y", op.PubkeyG1.Y.Text(16))
			nonSignerOperatorIds = append(nonSignerOperatorIds, id)
		}
	}

	quorumAggPubKeys := make([]*G1Point, len(quorumIDs))

	// Validate the amount signed and aggregate signatures for each quorum
	quorumResults := make(map[QuorumID]*QuorumResult)

	for ind, id := range quorumIDs {
		// Check that quorum has sufficient stake
		percent := GetSignedPercentage(state.OperatorState, id, stakeSigned[ind])
		quorumResults[id] = &QuorumResult{
			QuorumID:      id,
			PercentSigned: percent,
		}

		// Verify that the aggregated public key for the quorum matches the on-chain quorum aggregate public key sans non-signers of the quorum
		quorumAggKey := state.AggKeys[id]
		quorumAggPubKeys[ind] = quorumAggKey

		signersAggKey := quorumAggKey.Deserialize(quorumAggKey.Serialize())
		for opInd, nsk := range nonSignerKeys {
			ops := state.Operators[id]
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
	// ref: https://github.com/Layr-Labs/eigenlayer-contracts/blob/master/src/contracts/middleware/BLSSignatureChecker.sol#L99
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
	}, nil

}

func GetStakeThreshold(state *OperatorState, quorum QuorumID, quorumThreshold uint8) *big.Int {

	// Get stake threshold
	quorumThresholdBig := new(big.Int).SetUint64(uint64(quorumThreshold))
	stakeThreshold := new(big.Int)
	stakeThreshold.Mul(quorumThresholdBig, state.Totals[quorum].Stake)
	stakeThreshold = roundUpDivideBig(stakeThreshold, new(big.Int).SetUint64(PercentMultiplier))

	return stakeThreshold
}

func GetSignedPercentage(state *OperatorState, quorum QuorumID, stakeAmount *big.Int) uint8 {

	stakeAmount = stakeAmount.Mul(stakeAmount, new(big.Int).SetUint64(PercentMultiplier))
	quorumThresholdBig := stakeAmount.Div(stakeAmount, state.Totals[quorum].Stake)

	quorumThreshold := uint8(quorumThresholdBig.Uint64())

	return quorumThreshold
}
