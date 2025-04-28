package core

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

// signatureReceiver is a struct for receiving SigningMessages for a single batch. It should never be instantiated
// manually: it exists only as a helper struct for the ReceiveSignatures method.
type signatureReceiver struct {
	logger logging.Logger
	// indexedOperatorState contains operator information including pubkeys, stakes, and quorum membership
	indexedOperatorState *IndexedOperatorState

	// validSignerMap tracks which operators have already submitted valid signatures
	validSignerMap map[OperatorID]bool
	// signatureMessageReceived tracks which operators have submitted signature messages, whether valid or invalid.
	// this is tracked separately from signerMap, since signerMap only includes valid signatures
	signatureMessageReceived map[OperatorID]bool
	// aggregateSignatures stores the accumulated BLS signatures for each quorum
	aggregateSignatures map[QuorumID]*Signature
	// aggregateSignersG2PubKeys stores the accumulated G2 public keys of signers for each quorum
	aggregateSignersG2PubKeys map[QuorumID]*G2Point

	// stakeSigned tracks the total stake that has signed for each quorum
	stakeSigned map[QuorumID]*big.Int

	// batchHeaderHash is the hash of the batch header that operators are signing
	batchHeaderHash [32]byte
	// signingMessageChan is the channel through which SigningMessages are received
	signingMessageChan chan SigningMessage
	// quorumIDs is a sorted list of quorum IDs for which signatures are being collected
	quorumIDs []QuorumID

	// tickInterval determines how frequently intermediate attestations are yielded
	tickInterval time.Duration
}

// ReceiveSignatures receives SigningMessages over the signingMessageChan, and yields QuorumAttestations produced
// from these SigningMessages.
//
// The yielded QuorumAttestations contain aggregate signing data from all SigningMessages received thus far. Each
// QuorumAttestation will have incorporated more SigningMessages than the previously yielded QuorumAttestation.
//
// This channel will be closed when one of the following conditions is met:
// 1. The global attestation timeout is exceeded
// 2. A SigningMessage from every Operator has been received and processed
//
// Before being closed, the QuorumAttestation chan will have returned a QuorumAttestation containing data from every
// gathered SigningMessage.
func ReceiveSignatures(
	ctx context.Context,
	logger logging.Logger,
	indexedOperatorState *IndexedOperatorState,
	batchHeaderHash [32]byte,
	signingMessageChan chan SigningMessage,
	tickInterval time.Duration,
) (chan *QuorumAttestation, error) {
	sortedQuorumIDs, err := getSortedQuorumIDs(indexedOperatorState)
	if err != nil {
		return nil, fmt.Errorf("get sorted quorum ids: %w", err)
	}

	validSignerMap := make(map[OperatorID]bool)
	signatureMessageReceived := make(map[OperatorID]bool)
	aggregateSignatures := make(map[QuorumID]*Signature, len(sortedQuorumIDs))
	aggregateSignersG2PubKeys := make(map[QuorumID]*G2Point, len(sortedQuorumIDs))

	// initialized stakeSigned map with 0 stake signed for each quorum
	stakeSigned := make(map[QuorumID]*big.Int, len(sortedQuorumIDs))
	for _, quorumID := range sortedQuorumIDs {
		stakeSigned[quorumID] = big.NewInt(0)
	}

	receiver := &signatureReceiver{
		logger:                    logger,
		indexedOperatorState:      indexedOperatorState,
		aggregateSignatures:       aggregateSignatures,
		validSignerMap:            validSignerMap,
		signatureMessageReceived:  signatureMessageReceived,
		aggregateSignersG2PubKeys: aggregateSignersG2PubKeys,
		stakeSigned:               stakeSigned,
		batchHeaderHash:           batchHeaderHash,
		signingMessageChan:        signingMessageChan,
		quorumIDs:                 sortedQuorumIDs,
		tickInterval:              tickInterval,
	}

	attestationChan := make(chan *QuorumAttestation, len(indexedOperatorState.IndexedOperators))
	go receiver.receiveSigningMessages(ctx, attestationChan)

	return attestationChan, nil
}

// receiveSigningMessages receives SigningMessages, and sends QuorumAttestations to the input attestationChan
func (sr *signatureReceiver) receiveSigningMessages(ctx context.Context, attestationChan chan *QuorumAttestation) {
	// this ticker causes QuorumAttestations to be sent periodically
	ticker := time.NewTicker(sr.tickInterval)
	defer ticker.Stop()
	defer close(attestationChan)

	operatorCount := len(sr.indexedOperatorState.IndexedOperators)
	errorCount := 0
	newSignaturesGathered := false

	// we expect a single SigningMessage from each operator
	for len(sr.signatureMessageReceived) < operatorCount {
		breakLoop := false
		select {
		case <-ctx.Done():
			sr.logger.Infof(
				"global batch attestation timeout exceeded for batch %s. Received and processed %d/%d signing "+
					"messages. %d of the signing messages caused an error during processing",
				hex.EncodeToString(sr.batchHeaderHash[:]), len(sr.signatureMessageReceived), operatorCount, errorCount)
			breakLoop = true
		case signingMessage, ok := <-sr.signingMessageChan:
			if !ok {
				sr.logger.Errorf(
					"signing message channel closed for batch %s. Received and processed %d/%d signing "+
						"messages. %d of the signing messages caused an error during processing",
					hex.EncodeToString(sr.batchHeaderHash[:]),
					len(sr.signatureMessageReceived),
					operatorCount,
					errorCount)
				breakLoop = true
				break
			}
			indexedOperatorInfo, found := sr.indexedOperatorState.IndexedOperators[signingMessage.Operator]
			if !found {
				sr.logger.Warn("operator not found in state",
					"batchHeaderHash", hex.EncodeToString(sr.batchHeaderHash[:]),
					"operatorID", signingMessage.Operator.Hex(),
					"attestationLatencyMs", signingMessage.AttestationLatencyMs)
				continue
			}

			if seen := sr.signatureMessageReceived[signingMessage.Operator]; seen {
				sr.logger.Warn("duplicate message from operator",
					"batchHeaderHash", hex.EncodeToString(sr.batchHeaderHash[:]),
					"operatorID", signingMessage.Operator.Hex(),
					"attestationLatencyMs", signingMessage.AttestationLatencyMs)
				continue
			}

			// this map records messages received, whether the messages are valid or not
			sr.signatureMessageReceived[signingMessage.Operator] = true

			err := sr.processSigningMessage(signingMessage, indexedOperatorInfo)
			if err != nil {
				errorCount++
				sr.logger.Warn("error processing signing message",
					"batchHeaderHash", hex.EncodeToString(sr.batchHeaderHash[:]),
					"operatorID", signingMessage.Operator.Hex(),
					"attestationLatencyMs", signingMessage.AttestationLatencyMs,
					"error", err)
				continue
			}

			// record that we've received a valid message from this operator
			sr.validSignerMap[signingMessage.Operator] = true
			newSignaturesGathered = true
		// The ticker case is intentionally ordered after the message receiving case. If there are SigningMessages
		// waiting to be handled, we shouldn't delay their processing for the sake of yielding a QuorumAttestation.
		// The most likely time for there to be a backlog of SigningMessages is early-on in the signature gathering
		// process, when we are unlikely to have reached a threshold of signatures anyway.
		case <-ticker.C:
			if !newSignaturesGathered {
				continue
			}

			sr.submitAttestation(attestationChan)
			newSignaturesGathered = false
		}

		if breakLoop {
			break
		}
	}

	if newSignaturesGathered {
		sr.submitAttestation(attestationChan)
	}
}

// getSortedQuorumIDs returns a sorted slice of QuorumIDs from the state
func getSortedQuorumIDs(state *IndexedOperatorState) ([]QuorumID, error) {
	quorumIDs := make([]QuorumID, 0, len(state.AggKeys))
	for quorumID := range state.Operators {
		quorumIDs = append(quorumIDs, quorumID)
	}
	slices.Sort(quorumIDs)

	if len(quorumIDs) == 0 {
		return nil, errors.New("number of quorums must be greater than zero")
	}

	return quorumIDs, nil
}

// processSigningMessage accepts a SigningMessage, verifies it, and updates the signatureReceiver aggregates accordingly
func (sr *signatureReceiver) processSigningMessage(
	signingMessage SigningMessage,
	indexedOperatorInfo *IndexedOperatorInfo,
) error {
	if signingMessage.Err != nil {
		return fmt.Errorf("signingMessage contained error: %w", signingMessage.Err)
	}

	operatorPubkey := indexedOperatorInfo.PubkeyG2
	if !signingMessage.Signature.Verify(operatorPubkey, sr.batchHeaderHash) {
		return fmt.Errorf("signature verification with pubkey %s", hex.EncodeToString(operatorPubkey.Serialize()))
	}

	for _, quorumID := range sr.quorumIDs {
		quorumOperators := sr.indexedOperatorState.Operators[quorumID]
		quorumOperatorInfo, isOperatorInQuorum := quorumOperators[signingMessage.Operator]
		if !isOperatorInQuorum {
			// if the operator which sent the signing message isn't in a given quorum, then we shouldn't make any
			// changes to the aggregates that are tracked on a per-quorum basis
			continue
		}

		sr.stakeSigned[quorumID].Add(sr.stakeSigned[quorumID], quorumOperatorInfo.Stake)

		if sr.aggregateSignatures[quorumID] == nil {
			sr.aggregateSignatures[quorumID] = &Signature{signingMessage.Signature.Clone()}
			sr.aggregateSignersG2PubKeys[quorumID] = indexedOperatorInfo.PubkeyG2.Clone()
		} else {
			sr.aggregateSignatures[quorumID].Add(signingMessage.Signature.G1Point)
			sr.aggregateSignersG2PubKeys[quorumID].Add(indexedOperatorInfo.PubkeyG2)
		}
	}

	return nil
}

// submitAttestation aggregates and submits a QuorumAttestation representing the most up-to-date aggregates
func (sr *signatureReceiver) submitAttestation(attestationChan chan *QuorumAttestation) {
	nonSignerMap := make(map[OperatorID]*G1Point)
	for operatorID, operatorInfo := range sr.indexedOperatorState.IndexedOperators {
		_, found := sr.validSignerMap[operatorID]
		if !found {
			nonSignerMap[operatorID] = operatorInfo.PubkeyG1
		}
	}

	quorumResults := make(map[QuorumID]*QuorumResult)
	for _, quorumID := range sr.quorumIDs {
		quorumResult, err := sr.computeQuorumResult(quorumID, nonSignerMap)
		if err != nil {
			sr.logger.Error("compute quorum result",
				"quorumID", quorumID,
				"batchHeaderHash", sr.batchHeaderHash,
				"error", err)
			continue
		}
		quorumResults[quorumID] = quorumResult
	}

	// Make copies of the maps that are populated while receiving signatures. The yielded QuorumAttestation will be
	// handled by a separate routine, so it's important that we don't mutate these maps after they are yielded.
	quorumAggPubKeyCopy := make(map[QuorumID]*G1Point, len(sr.indexedOperatorState.AggKeys))
	for quorumID, g1Point := range sr.indexedOperatorState.AggKeys {
		// TODO: is this ok? semantics are changed from before: we used to exclude aggregate keys of quorums that had no
		//  signatures, but I don't see why that case should be special.
		quorumAggPubKeyCopy[quorumID] = g1Point.Clone()
	}
	aggregateSignersG2PubKeysCopy := make(map[QuorumID]*G2Point, len(sr.aggregateSignersG2PubKeys))
	for quorumID, aggregatePubkey := range sr.aggregateSignersG2PubKeys {
		aggregateSignersG2PubKeysCopy[quorumID] = aggregatePubkey.Clone()
	}
	aggregateSignaturesCopy := make(map[QuorumID]*Signature, len(sr.aggregateSignatures))
	for quorumID, aggregateSignature := range sr.aggregateSignatures {
		aggregateSignaturesCopy[quorumID] = &Signature{aggregateSignature.Clone()}
	}
	validSignerMapCopy := make(map[OperatorID]bool, len(sr.validSignerMap))
	for operatorID, signed := range sr.validSignerMap {
		validSignerMapCopy[operatorID] = signed
	}

	attestationChan <- &QuorumAttestation{
		QuorumAggPubKey:  quorumAggPubKeyCopy,
		SignersAggPubKey: aggregateSignersG2PubKeysCopy,
		AggSignature:     aggregateSignaturesCopy,
		QuorumResults:    quorumResults,
		SignerMap:        validSignerMapCopy,
	}
}

// computeQuorumResult creates a QuorumResult for a given quorum
func (sr *signatureReceiver) computeQuorumResult(
	quorumID QuorumID,
	nonSignerMap map[OperatorID]*G1Point,
) (*QuorumResult, error) {
	signedPercentage := getSignedPercentage(
		sr.stakeSigned[quorumID],
		sr.indexedOperatorState.Totals[quorumID].Stake)

	if signedPercentage == 0 {
		return &QuorumResult{
			QuorumID:      quorumID,
			PercentSigned: 0,
		}, nil
	}

	// clone the quorum aggregate G1 pubkey, so that we can safely subtract non-signer pubkeys to yield the aggregate
	// G1 pubkey of all the signers
	aggregateSignersG1PubKey := sr.indexedOperatorState.AggKeys[quorumID].Clone()
	for operatorID := range sr.indexedOperatorState.Operators[quorumID] {
		if nonSignerPubKey, ok := nonSignerMap[operatorID]; ok {
			aggregateSignersG1PubKey.Sub(nonSignerPubKey)
		}
	}

	if sr.aggregateSignersG2PubKeys[quorumID] == nil {
		return nil, errors.New("nil aggregate signer G2 public key")
	}

	ok, err := aggregateSignersG1PubKey.VerifyEquivalence(sr.aggregateSignersG2PubKeys[quorumID])
	if err != nil {
		return nil, fmt.Errorf("verify aggregate G1 and G2 pubkey equivalence: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf(
			"aggregate signers G1 pubkey is not equivalent to aggregate signers G2 pubkey: %s != %s",
			hex.EncodeToString(aggregateSignersG1PubKey.Serialize()),
			hex.EncodeToString(sr.aggregateSignersG2PubKeys[quorumID].Serialize()))
	}

	// Verify the aggregate signature for the quorum
	ok = sr.aggregateSignatures[quorumID].Verify(sr.aggregateSignersG2PubKeys[quorumID], sr.batchHeaderHash)
	if !ok {
		return nil, errors.New("aggregated signature is not valid")
	}

	return &QuorumResult{
		QuorumID:      quorumID,
		PercentSigned: signedPercentage,
	}, nil
}

// getSignedPercentage the amount is signedStake, and the totalStake. It returns a uint8 representing the percentage
// of the total stake that has signed.
func getSignedPercentage(signedStake *big.Int, totalStake *big.Int) uint8 {
	if totalStake.Cmp(big.NewInt(0)) == 0 {
		// avoid dividing by 0
		return 0
	}

	// the calculation being performed here is: signedStake * 100 / totalStake

	signedStakeNumerator := new(big.Int).Mul(signedStake, new(big.Int).SetUint64(percentMultiplier))
	quorumThreshold := uint8(new(big.Int).Div(signedStakeNumerator, totalStake).Uint64())

	return quorumThreshold
}
