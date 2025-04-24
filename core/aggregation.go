package core

import (
	"bytes"
	"context"
	"fmt"
	"sort"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

type SigningMessage struct {
	Signature       *Signature
	Operator        OperatorID
	BatchHeaderHash [32]byte
	// Undefined if this value <= 0.
	AttestationLatencyMs float64
	Err                  error
}

// QuorumAttestation contains the results of aggregating signatures from a set of operators by quorums
// It also returns map of all signers across all quorums
type QuorumAttestation struct {
	// QuorumAggPubKeys contains the aggregated public keys for all of the operators each quorum,
	// including those that did not sign
	QuorumAggPubKey map[QuorumID]*G1Point
	// SignersAggPubKey is the aggregated public key for all of the operators that signed the message by each quorum
	SignersAggPubKey map[QuorumID]*G2Point
	// AggSignature is the aggregated signature for all of the operators that signed the message for each quorum,
	// mirroring the SignersAggPubKey.
	AggSignature map[QuorumID]*Signature
	// QuorumResults contains the quorum ID and the amount signed for each quorum
	QuorumResults map[QuorumID]*QuorumResult
	// SignerMap contains the operator IDs that signed the message
	SignerMap map[OperatorID]bool
}

// SignatureAggregation contains the results of aggregating signatures from a set of operators across multiple quorums
type SignatureAggregation struct {
	// NonSigners contains the public keys of the operators that did not sign the message
	NonSigners []*G1Point
	// QuorumAggPubKeys contains the aggregated public keys for all of the operators each quorum,
	// Including those that did not sign
	QuorumAggPubKeys map[QuorumID]*G1Point
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

// SignatureAggregator is an interface for aggregating the signatures returned by DA nodes
// so that they can be verified by the DA contract
type SignatureAggregator interface {
	// ReceiveSignatures blocks until it receives a response for each operator in the operator state via messageChan,
	// and then returns the attestation result by quorum.
	//
	// This function accepts two contexts. ctx is the background context. attestationCtx is a context that is cancelled
	// once the attestation period is over. If the attestationCtx is cancelled, the function will stop waiting for
	// responses and return the result of the signatures received so far.
	//
	// TODO (litt3): this method is only used by V1. When V1 support is removed, this method should be removed.
	ReceiveSignatures(
		ctx context.Context,
		state *IndexedOperatorState,
		message [32]byte,
		messageChan chan SigningMessage,
	) (*QuorumAttestation, error)
	// AggregateSignatures takes attestation result by quorum and aggregates the signatures across them.
	// If the aggregated signature is invalid, an error is returned.
	AggregateSignatures(
		ctx context.Context,
		ics IndexedChainState,
		referenceBlockNumber uint,
		quorumAttestation *QuorumAttestation,
		quorumIDs []QuorumID,
	) (*SignatureAggregation, error)
}

type StdSignatureAggregator struct {
	Logger     logging.Logger
	Transactor Reader
}

func NewStdSignatureAggregator(logger logging.Logger, transactor Reader) (*StdSignatureAggregator, error) {
	return &StdSignatureAggregator{
		Logger: logger.With(
			"component", "SignatureAggregator"),
		Transactor:        transactor,
	}, nil
}

var _ SignatureAggregator = (*StdSignatureAggregator)(nil)

func (a *StdSignatureAggregator) ReceiveSignatures(
	ctx context.Context,
	state *IndexedOperatorState,
	message [32]byte,
	messageChan chan SigningMessage,
) (*QuorumAttestation, error) {
	attestationChan, err := ReceiveSignatures(ctx, a.Logger, state, message, messageChan)
	if err != nil {
		return nil, fmt.Errorf("receive signatures: %w", err)
	}

	var finalAttestation *QuorumAttestation
	for receivedAttestation := range attestationChan {
		finalAttestation = receivedAttestation
	}

	return finalAttestation, nil
}

func (a *StdSignatureAggregator) AggregateSignatures(
	ctx context.Context,
	ics IndexedChainState,
	referenceBlockNumber uint,
	quorumAttestation *QuorumAttestation,
	quorumIDs []QuorumID,
) (*SignatureAggregation, error) {
	// Aggregate the aggregated signatures. We reuse the first aggregated signature as the accumulator
	var aggSig *Signature
	for _, quorumID := range quorumIDs {
		if quorumAttestation.AggSignature[quorumID] == nil {
			a.Logger.Error("cannot aggregate signature for quorum because aggregated signature is nil",
				"quorumID", quorumID)
			continue
		}
		sig := quorumAttestation.AggSignature[quorumID]
		if aggSig == nil {
			aggSig = &Signature{sig.G1Point.Clone()}
		} else {
			aggSig.Add(sig.G1Point)
		}
	}

	// Aggregate the aggregated public keys. We reuse the first aggregated public key as the accumulator
	var aggPubKey *G2Point
	for _, quorumID := range quorumIDs {
		if quorumAttestation.SignersAggPubKey[quorumID] == nil {
			a.Logger.Error("cannot aggregate public key for quorum because signers aggregated public key is nil",
				"quorumID", quorumID)
			continue
		}
		apk := quorumAttestation.SignersAggPubKey[quorumID]
		if aggPubKey == nil {
			aggPubKey = apk.Clone()
		} else {
			aggPubKey.Add(apk)
		}
	}

	nonSignerKeys := make([]*G1Point, 0)
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
		// sort in ascending order
		return bytes.Compare(hash1[:], hash2[:]) == -1
	})

	quorumAggKeys := make(map[QuorumID]*G1Point, len(quorumIDs))
	for _, quorumID := range quorumIDs {
		if quorumAttestation.QuorumAggPubKey[quorumID] == nil {
			a.Logger.Error("cannot aggregate public key for quorum because aggregated public key is nil",
				"quorumID", quorumID)
			continue
		}
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
