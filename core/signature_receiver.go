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
	logger               logging.Logger
	indexedOperatorState *IndexedOperatorState

	signerMap                 map[OperatorID]bool
	aggregateSignatures       map[QuorumID]*Signature
	aggregateSignersG2PubKeys map[QuorumID]*G2Point

	stakeSigned map[QuorumID]*big.Int

	batchHeaderHash    [32]byte
	signingMessageChan chan SigningMessage
	quorumIDs          []QuorumID
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
) (chan *QuorumAttestation, error) {
	sortedQuorumIDs, err := getSortedQuorumIDs(indexedOperatorState)
	if err != nil {
		return nil, fmt.Errorf("get sorted quorum ids: %w", err)
	}

	signerMap := make(map[OperatorID]bool)
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
		signerMap:                 signerMap,
		aggregateSignersG2PubKeys: aggregateSignersG2PubKeys,
		stakeSigned:               stakeSigned,
		batchHeaderHash:           batchHeaderHash,
		signingMessageChan:        signingMessageChan,
		quorumIDs:                 sortedQuorumIDs,
	}

	attestationChan := make(chan *QuorumAttestation, len(signerMap))
	go receiver.receiveSigningMessages(ctx, attestationChan)

	return attestationChan, nil
}

// receiveSigningMessages receives SigningMessages, and sends QuorumAttestations to the input attestationChan
func (sr *signatureReceiver) receiveSigningMessages(ctx context.Context, attestationChan chan *QuorumAttestation) {
	// this ticker causes QuorumAttestations to be sent periodically
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer close(attestationChan)

	operatorCount := len(sr.indexedOperatorState.IndexedOperators)
	signingMessageCount := 0
	errorCount := 0
	newSignaturesGathered := false

	// we expect a single SigningMessage from each operator
	for signingMessageCount < operatorCount {
		contextExpired := false
		select {
		case <-ctx.Done():
			sr.logger.Infof(
				`global batch attestation timeout exceeded for batch %s. Recieved and processed %d/%d signing
						messages. %d of the signing messages caused an error during processing`,
				hex.EncodeToString(sr.batchHeaderHash[:]), signingMessageCount, operatorCount, errorCount)
			contextExpired = true
		case signingMessage := <-sr.signingMessageChan:
			signingMessageCount++
			err := sr.processSigningMessage(signingMessage)
			if err != nil {
				errorCount++
				sr.logger.Warn("process signing message",
					"batchHeaderHash", hex.EncodeToString(sr.batchHeaderHash[:]),
					"operatorID", signingMessage.Operator.Hex(),
					"attestationLatencyMs", signingMessage.AttestationLatencyMs,
					"error", err)
				continue
			}
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

		if contextExpired {
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
func (sr *signatureReceiver) processSigningMessage(signingMessage SigningMessage) error {
	indexedOperatorInfo, err := sr.checkSigningMessage(signingMessage)
	if err != nil {
		return fmt.Errorf("check signing message: %w", err)
	}

	// record that we've received a message from this operator
	sr.signerMap[signingMessage.Operator] = true

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

// checkSigningMessage checks the input SigningMessage, and returns an error if any check fails
func (sr *signatureReceiver) checkSigningMessage(signingMessage SigningMessage) (*IndexedOperatorInfo, error) {
	if seen := sr.signerMap[signingMessage.Operator]; seen {
		return nil, fmt.Errorf("duplicate message from operator")
	}

	if signingMessage.Err != nil {
		return nil, signingMessage.Err
	}

	indexedOperatorInfo, found := sr.indexedOperatorState.IndexedOperators[signingMessage.Operator]
	if !found {
		return nil, fmt.Errorf("operator not found in state")
	}

	operatorPubkey := indexedOperatorInfo.PubkeyG2
	if !signingMessage.Signature.Verify(operatorPubkey, sr.batchHeaderHash) {
		return nil, fmt.Errorf("signature verification for pubkey %s",
			hex.EncodeToString(operatorPubkey.Serialize()))
	}

	return indexedOperatorInfo, nil
}

// submitAttestation aggregates and submits a QuorumAttestation representing the most up-to-date aggregates
func (sr *signatureReceiver) submitAttestation(attestationChan chan *QuorumAttestation) {
	nonSignerMap := make(map[OperatorID]*G1Point)
	for operatorID, operatorInfo := range sr.indexedOperatorState.IndexedOperators {
		_, found := sr.signerMap[operatorID]
		if !found {
			nonSignerMap[operatorID] = operatorInfo.PubkeyG1
		}
	}

	quorumResults := make(map[QuorumID]*QuorumResult)
	for _, quorumID := range sr.quorumIDs {
		quorumResult, err := sr.computeQuorumResult(quorumID, nonSignerMap)
		if err != nil {
			sr.logger.Error("compute quorum result",
				"quorumID", quorumID, "batchHeaderHash", sr.batchHeaderHash)
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
		quorumAggPubKeyCopy[quorumID] = g1Point
	}
	aggregateSignersG2PubKeysCopy := make(map[QuorumID]*G2Point, len(sr.aggregateSignersG2PubKeys))
	for quorumID, aggregatePubkey := range sr.aggregateSignersG2PubKeys {
		aggregateSignersG2PubKeysCopy[quorumID] = aggregatePubkey
	}
	aggregateSignaturesCopy := make(map[QuorumID]*Signature, len(sr.aggregateSignatures))
	for quorumID, aggregateSignature := range sr.aggregateSignatures {
		aggregateSignaturesCopy[quorumID] = aggregateSignature
	}
	signerMapCopy := make(map[OperatorID]bool, len(sr.signerMap))
	for operatorID, signed := range sr.signerMap {
		signerMapCopy[operatorID] = signed
	}

	attestationChan <- &QuorumAttestation{
		QuorumAggPubKey:  quorumAggPubKeyCopy,
		SignersAggPubKey: aggregateSignersG2PubKeysCopy,
		AggSignature:     aggregateSignaturesCopy,
		QuorumResults:    quorumResults,
		SignerMap:        signerMapCopy,
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
		return nil, fmt.Errorf("quorum %v has 0%% signed percentage", quorumID)
	}

	// clone the quorum aggregate G1 pubkey, so that we can safely subtract non-signer pubkeys to yield the aggregate
	// G1 pubkey of all the signers
	aggregateSignersG1PubKey := sr.indexedOperatorState.AggKeys[quorumID].Clone()
	for nonSignerOperatorID, nonSignerPubKey := range nonSignerMap {
		quorumOperatorInfo := sr.indexedOperatorState.Operators[quorumID]
		if _, ok := quorumOperatorInfo[nonSignerOperatorID]; ok {
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
