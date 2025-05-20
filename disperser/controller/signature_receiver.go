package controller

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// signatureReceiver is a struct for receiving SigningMessages for a single batch. It should never be instantiated
// manually: it exists only as a helper struct for the ReceiveSignatures method.
type signatureReceiver struct {
	logger logging.Logger
	// metrics may be nil, in which case no metrics will be reported
	metrics *dispatcherMetrics

	// indexedOperatorState contains operator information including pubkeys, stakes, and quorum membership
	indexedOperatorState *core.IndexedOperatorState

	// validSignerMap tracks which operators have already submitted valid signatures
	validSignerMap map[core.OperatorID]bool
	// signatureMessageReceived tracks which operators have submitted signature messages, whether valid or invalid.
	// this is tracked separately from signerMap, since signerMap only includes valid signatures
	signatureMessageReceived map[core.OperatorID]bool
	// aggregateSignatures stores the accumulated BLS signatures for each quorum
	aggregateSignatures map[core.QuorumID]*core.Signature
	// aggregateSignersG2PubKeys stores the accumulated G2 public keys of signers for each quorum
	aggregateSignersG2PubKeys map[core.QuorumID]*core.G2Point

	// stakeSigned tracks the total stake that has signed for each quorum
	stakeSigned map[core.QuorumID]*big.Int

	// batchHeaderHash is the hash of the batch header that operators are signing
	batchHeaderHash [32]byte
	// signingMessageChan is the channel through which SigningMessages are received
	signingMessageChan chan core.SigningMessage
	// quorumIDs is a sorted list of quorum IDs for which signatures are being collected
	quorumIDs []core.QuorumID

	// tickInterval determines how frequently intermediate attestations are yielded
	tickInterval time.Duration

	// attestationUpdateStart is initialized when we first start receiving signatures, and is updated each time an
	// attestation is yielded. This is used to track how long it takes to yield each attestation.
	attestationUpdateStart time.Time

	// significantSigningThresholdPercentage is a configurable "important" signing threshold. Right now, it's being
	// used to track signing metrics, to understand system performance. If the value is 0, then special handling for
	// the threshold is disabled.
	// TODO (litt3): this might eventually be used to cause special case handling at an "important" threshold, e.g.
	//  "update the attestation as soon as the threshold is reached."
	significantSigningThresholdPercentage uint8

	// significantSigningThresholdReachedTime tracks when each quorum's signing percentage first reached or exceeded the
	// significantSigningThresholdPercentage
	significantSigningThresholdReachedTime map[core.QuorumID]time.Time

	// Tracks whether there are new signatures that have been gathered but not aggregated.
	newSignaturesGathered bool

	// The number of attestations received and processed so far.
	attestationUpdateCount int

	// A ticker used to periodically yield QuorumAttestations.
	ticker *time.Ticker

	// The number of errors encountered while processing SigningMessages.
	errorCount int
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
	metrics *dispatcherMetrics,
	indexedOperatorState *core.IndexedOperatorState,
	batchHeaderHash [32]byte,
	signingMessageChan chan core.SigningMessage,
	tickInterval time.Duration,
	significantSigningThresholdPercentage uint8,
) (chan *core.QuorumAttestation, error) {
	sortedQuorumIDs, err := getSortedQuorumIDs(indexedOperatorState)
	if err != nil {
		return nil, fmt.Errorf("get sorted quorum ids: %w", err)
	}

	validSignerMap := make(map[core.OperatorID]bool)
	signatureMessageReceived := make(map[core.OperatorID]bool)
	aggregateSignatures := make(map[core.QuorumID]*core.Signature, len(sortedQuorumIDs))
	aggregateSignersG2PubKeys := make(map[core.QuorumID]*core.G2Point, len(sortedQuorumIDs))

	// initialized stakeSigned map with 0 stake signed for each quorum
	stakeSigned := make(map[core.QuorumID]*big.Int, len(sortedQuorumIDs))
	for _, quorumID := range sortedQuorumIDs {
		stakeSigned[quorumID] = big.NewInt(0)
	}

	significantSigningThresholdReachedTime := make(map[core.QuorumID]time.Time, len(sortedQuorumIDs))

	receiver := &signatureReceiver{
		logger:                                 logger,
		metrics:                                metrics,
		indexedOperatorState:                   indexedOperatorState,
		aggregateSignatures:                    aggregateSignatures,
		validSignerMap:                         validSignerMap,
		signatureMessageReceived:               signatureMessageReceived,
		aggregateSignersG2PubKeys:              aggregateSignersG2PubKeys,
		stakeSigned:                            stakeSigned,
		batchHeaderHash:                        batchHeaderHash,
		signingMessageChan:                     signingMessageChan,
		quorumIDs:                              sortedQuorumIDs,
		tickInterval:                           tickInterval,
		significantSigningThresholdPercentage:  significantSigningThresholdPercentage,
		significantSigningThresholdReachedTime: significantSigningThresholdReachedTime,
		ticker:                                 time.NewTicker(tickInterval),
	}

	attestationChan := make(chan *core.QuorumAttestation, len(indexedOperatorState.IndexedOperators))
	go receiver.receiveSigningMessages(ctx, attestationChan)

	return attestationChan, nil
}

// receiveSigningMessages receives SigningMessages, and sends QuorumAttestations to the input attestationChan
func (sr *signatureReceiver) receiveSigningMessages(ctx context.Context, attestationChan chan *core.QuorumAttestation) {
	defer sr.ticker.Stop()
	defer close(attestationChan)

	// the number of attestations submitted by this method
	defer func() {
		if sr.metrics != nil {
			sr.reportThresholdSignedToDoneLatency()
			sr.metrics.reportAttestationUpdateCount(float64(sr.attestationUpdateCount))
		}
	}()

	sr.attestationUpdateStart = time.Now()

	operatorCount := len(sr.indexedOperatorState.IndexedOperators)

	// we expect a single SigningMessage from each operator
	for len(sr.signatureMessageReceived) < operatorCount {
		breakLoop := false
		select {
		case <-ctx.Done():
			sr.logger.Infof(
				"global batch attestation timeout exceeded for batch %s. Received and processed %d/%d signing "+
					"messages. %d of the signing messages caused an error during processing", hex.EncodeToString(
					sr.batchHeaderHash[:]), len(sr.signatureMessageReceived), operatorCount, sr.errorCount)
			breakLoop = true
			break
		case signingMessage, ok := <-sr.signingMessageChan:
			if !ok {
				sr.logger.Errorf(
					"signing message channel closed for batch %s. Received and processed %d/%d signing "+
						"messages. %d of the signing messages caused an error during processing",
					hex.EncodeToString(sr.batchHeaderHash[:]),
					len(sr.signatureMessageReceived),
					operatorCount,
					sr.errorCount)
				breakLoop = true
				break
			}

			sr.handleNextSignature(signingMessage, attestationChan)

		// The ticker case is intentionally ordered after the message receiving case. If there are SigningMessages
		// waiting to be handled, we shouldn't delay their processing for the sake of yielding a QuorumAttestation.
		// The most likely time for there to be a backlog of SigningMessages is early-on in the signature gathering
		// process, when we are unlikely to have reached a threshold of signatures anyway.
		case <-sr.ticker.C:
			sr.buildAndSubmitAttestation(attestationChan)
		}

		if breakLoop {
			break
		}
	}

	// Aggregate any remaining signatures and submit an attestation.
	sr.buildAndSubmitAttestation(attestationChan)
}

func (sr *signatureReceiver) handleNextSignature(
	signingMessage core.SigningMessage,
	attestationChan chan *core.QuorumAttestation,
) {

	if signingMessage.TimeReceived.IsZero() {
		sr.logger.Errorf("signing message from %s time received is zero in batch %s. "+
			"This shouldn't be possible.",
			signingMessage.Operator.Hex(),
			hex.EncodeToString(sr.batchHeaderHash[:]))
	} else if sr.metrics != nil {
		sr.metrics.reportSigningMessageChannelLatency(time.Since(signingMessage.TimeReceived))
	}

	indexedOperatorInfo, found := sr.indexedOperatorState.IndexedOperators[signingMessage.Operator]
	if !found {
		sr.logger.Error("operator not found in state",
			"batchHeaderHash", hex.EncodeToString(sr.batchHeaderHash[:]),
			"operatorID", signingMessage.Operator.Hex(),
			"attestationLatencyMs", signingMessage.AttestationLatencyMs)
		return
	}

	if seen := sr.signatureMessageReceived[signingMessage.Operator]; seen {
		sr.logger.Error("duplicate message from operator",
			"batchHeaderHash", hex.EncodeToString(sr.batchHeaderHash[:]),
			"operatorID", signingMessage.Operator.Hex(),
			"attestationLatencyMs", signingMessage.AttestationLatencyMs)
		return
	}

	// this map records messages received, whether the messages are valid or not
	sr.signatureMessageReceived[signingMessage.Operator] = true

	thresholdCrossed, err := sr.processSigningMessage(signingMessage, indexedOperatorInfo)
	if err != nil {
		sr.errorCount++
		sr.logger.Warn("error processing signing message",
			"batchHeaderHash", hex.EncodeToString(sr.batchHeaderHash[:]),
			"operatorID", signingMessage.Operator.Hex(),
			"attestationLatencyMs", signingMessage.AttestationLatencyMs,
			"error", err)
		return
	}

	sr.validSignerMap[signingMessage.Operator] = true
	sr.newSignaturesGathered = true

	if thresholdCrossed {
		// Immediately build and submit an attestation.
		sr.buildAndSubmitAttestation(attestationChan)
		// Delay the next tick since we just submitted an attestation.
		sr.ticker.Reset(sr.tickInterval)
	}
}

// getSortedQuorumIDs returns a sorted slice of QuorumIDs from the state
func getSortedQuorumIDs(state *core.IndexedOperatorState) ([]core.QuorumID, error) {
	quorumIDs := make([]core.QuorumID, 0, len(state.Operators))
	for quorumID := range state.Operators {
		quorumIDs = append(quorumIDs, quorumID)
	}
	slices.Sort(quorumIDs)

	if len(quorumIDs) == 0 {
		return nil, errors.New("number of quorums must be greater than zero")
	}

	return quorumIDs, nil
}

// processSigningMessage accepts a SigningMessage, verifies it, and updates the signatureReceiver aggregates
// accordingly. Returns true if any quorums cross their signing threshold as a result of processing this message.
func (sr *signatureReceiver) processSigningMessage(
	signingMessage core.SigningMessage,
	indexedOperatorInfo *core.IndexedOperatorInfo,
) (bool, error) {
	processSigningMessageStart := time.Now()
	defer func() {
		if sr.metrics != nil {
			sr.metrics.reportProcessSigningMessageLatency(time.Since(processSigningMessageStart))
		}
	}()

	if signingMessage.Err != nil {
		return false, fmt.Errorf("signingMessage contained error: %w", signingMessage.Err)
	}

	operatorPubkey := indexedOperatorInfo.PubkeyG2
	if !signingMessage.Signature.Verify(operatorPubkey, sr.batchHeaderHash) {
		return false, fmt.Errorf("signature verification with pubkey %s", hex.EncodeToString(operatorPubkey.Serialize()))
	}

	thresholdCrossed := false
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
			sr.aggregateSignatures[quorumID] = &core.Signature{G1Point: signingMessage.Signature.Clone()}
			sr.aggregateSignersG2PubKeys[quorumID] = indexedOperatorInfo.PubkeyG2.Clone()
		} else {
			sr.aggregateSignatures[quorumID].Add(signingMessage.Signature.G1Point)
			sr.aggregateSignersG2PubKeys[quorumID].Add(indexedOperatorInfo.PubkeyG2)
		}

		thresholdCrossed = thresholdCrossed || sr.checkSigningPercentage(quorumID)
	}

	return thresholdCrossed, nil
}

// buildAndSubmitAttestation aggregates and submits a QuorumAttestation representing the most up-to-date aggregates
func (sr *signatureReceiver) buildAndSubmitAttestation(attestationChan chan *core.QuorumAttestation) {

	if !sr.newSignaturesGathered {
		// no work to be done
		return
	}
	sr.newSignaturesGathered = false
	sr.attestationUpdateCount++

	submitAttestationStart := time.Now()
	defer func() {
		if sr.metrics != nil {
			sr.metrics.reportAttestationBuildingLatency(time.Since(submitAttestationStart))
		}
	}()

	nonSignerMap := make(map[core.OperatorID]*core.G1Point)
	// operators that aren't in the validSignerMap are "non-signers"
	for operatorID, operatorInfo := range sr.indexedOperatorState.IndexedOperators {
		_, found := sr.validSignerMap[operatorID]
		if !found {
			nonSignerMap[operatorID] = operatorInfo.PubkeyG1
		}
	}

	quorumResults := make(map[core.QuorumID]*core.QuorumResult)
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
	quorumAggPubKeyCopy := make(map[core.QuorumID]*core.G1Point, len(sr.indexedOperatorState.AggKeys))
	for quorumID, g1Point := range sr.indexedOperatorState.AggKeys {
		quorumAggPubKeyCopy[quorumID] = g1Point.Clone()
	}
	aggregateSignersG2PubKeysCopy := make(map[core.QuorumID]*core.G2Point, len(sr.aggregateSignersG2PubKeys))
	for quorumID, aggregatePubkey := range sr.aggregateSignersG2PubKeys {
		aggregateSignersG2PubKeysCopy[quorumID] = aggregatePubkey.Clone()
	}
	aggregateSignaturesCopy := make(map[core.QuorumID]*core.Signature, len(sr.aggregateSignatures))
	for quorumID, aggregateSignature := range sr.aggregateSignatures {
		aggregateSignaturesCopy[quorumID] = &core.Signature{G1Point: aggregateSignature.Clone()}
	}
	validSignerMapCopy := make(map[core.OperatorID]bool, len(sr.validSignerMap))
	for operatorID, signed := range sr.validSignerMap {
		validSignerMapCopy[operatorID] = signed
	}

	attestationChan <- &core.QuorumAttestation{
		QuorumAggPubKey:  quorumAggPubKeyCopy,
		SignersAggPubKey: aggregateSignersG2PubKeysCopy,
		AggSignature:     aggregateSignaturesCopy,
		QuorumResults:    quorumResults,
		SignerMap:        validSignerMapCopy,
	}

	if sr.metrics != nil {
		sr.metrics.reportAttestationUpdateLatency(time.Since(sr.attestationUpdateStart))
	}
	sr.attestationUpdateStart = time.Now()
}

// computeQuorumResult creates a QuorumResult for a given quorum
func (sr *signatureReceiver) computeQuorumResult(
	quorumID core.QuorumID,
	nonSignerMap map[core.OperatorID]*core.G1Point,
) (*core.QuorumResult, error) {
	signedPercentage := getSignedPercentage(
		sr.stakeSigned[quorumID],
		sr.indexedOperatorState.Totals[quorumID].Stake)

	if signedPercentage == 0 {
		return &core.QuorumResult{
			QuorumID:      quorumID,
			PercentSigned: 0,
		}, nil
	}

	signerCount := 0

	// clone the quorum aggregate G1 pubkey, so that we can safely subtract non-signer pubkeys to yield the aggregate
	// G1 pubkey of all the signers
	aggregateSignersG1PubKey := sr.indexedOperatorState.AggKeys[quorumID].Clone()
	for operatorID := range sr.indexedOperatorState.Operators[quorumID] {
		operatorPubkey := sr.indexedOperatorState.IndexedOperators[operatorID].PubkeyG1

		if nonSignerPubKey, ok := nonSignerMap[operatorID]; ok {
			aggregateSignersG1PubKey.Sub(nonSignerPubKey)

			if !nonSignerPubKey.G1Affine.Equal(operatorPubkey.G1Affine) {
				sr.logger.Error("non-signer pubkey stored in non-signer map does not match indexed operator state pubkey",
					"pubkeyFromNonSignerMap", nonSignerPubKey.Serialize(),
					"pubkeyFromState", operatorPubkey.Serialize(),
				)
			}
		} else {
			signerCount++
		}
	}

	quorumOperatorCount := len(sr.indexedOperatorState.Operators[quorumID])
	nonSignerCount := len(nonSignerMap)

	stateOperatorCount := len(sr.indexedOperatorState.IndexedOperators)
	sr.logger.Debug("State details for quorum",
		"quorumID", quorumID,
		"totalStateOperatorCount", stateOperatorCount,
		"quorumOperatorCount", quorumOperatorCount,
		"quorumAggregateG1PubKey", sr.indexedOperatorState.AggKeys[quorumID].Serialize(),
		"signerCount", signerCount,
		"nonSignerCount", nonSignerCount,
		"batchHeaderHash", hex.EncodeToString(sr.batchHeaderHash[:]))

	if sr.aggregateSignersG2PubKeys[quorumID] == nil {
		return nil, errors.New("nil aggregate signer G2 public key")
	}

	ok, err := aggregateSignersG1PubKey.VerifyEquivalence(sr.aggregateSignersG2PubKeys[quorumID])
	if err != nil {
		return nil, fmt.Errorf("verify aggregate G1 and G2 pubkey equivalence: %w", err)
	}
	if !ok {
		sr.debugEquivalenceError(quorumID, nonSignerMap, aggregateSignersG1PubKey)

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

	return &core.QuorumResult{
		QuorumID:      quorumID,
		PercentSigned: signedPercentage,
	}, nil
}

// getSignedPercentage accepts the signedStake and the totalStake. It returns a uint8 representing the percentage
// of the total stake that has signed.
func getSignedPercentage(signedStake *big.Int, totalStake *big.Int) uint8 {
	if totalStake.Cmp(big.NewInt(0)) == 0 {
		// avoid dividing by 0
		return 0
	}

	// the calculation being performed here is: signedStake * 100 / totalStake

	signedStakeNumerator := new(big.Int).Mul(signedStake, new(big.Int).SetUint64(core.PercentMultiplier))
	quorumThreshold := uint8(new(big.Int).Div(signedStakeNumerator, totalStake).Uint64())

	return quorumThreshold
}

// checkSigningPercentage checks if the signing percentage for a quorum meets or exceeds the configured
// significantSigningThresholdPercentage, and records the time when the threshold was first crossed
// Returns true if the threshold was crossed, false otherwise. If called after the threshold was crossed, this
// method always returns false.
func (sr *signatureReceiver) checkSigningPercentage(quorumID core.QuorumID) bool {
	if sr.significantSigningThresholdPercentage == 0 {
		// if significantSigningThresholdPercentage is 0, skip
		return false
	}

	if !sr.significantSigningThresholdReachedTime[quorumID].IsZero() {
		// if significantSigningThresholdReachedTime[quorumID] has already been set, there is no need to check signing
		// percentage again, since the time has already been recorded
		return false
	}

	signedPercentage := getSignedPercentage(sr.stakeSigned[quorumID], sr.indexedOperatorState.Totals[quorumID].Stake)
	// check if the significantSigningThresholdPercentage has been crossed, and record the time if it has
	if signedPercentage >= sr.significantSigningThresholdPercentage {
		// Record the time when the threshold was first crossed
		sr.significantSigningThresholdReachedTime[quorumID] = time.Now()
		return true
	}

	return false
}

// reportThresholdSignedToDoneLatency calculates and reports the latency between the time when the
// significantSigningThresholdPercentage was first crossed, and now
func (sr *signatureReceiver) reportThresholdSignedToDoneLatency() {
	if sr.metrics == nil {
		return
	}

	for _, quorumID := range sr.quorumIDs {
		thresholdReachedTime := sr.significantSigningThresholdReachedTime[quorumID]
		if thresholdReachedTime.IsZero() {
			continue
		}

		sr.metrics.reportThresholdSignedToDoneLatency(quorumID, time.Since(thresholdReachedTime))
	}
}

// debugEquivalenceError is used to debug pubkey equivalence check failures by recomputing and comparing aggregate keys.
// Results are logged in this method.
func (sr *signatureReceiver) debugEquivalenceError(
	quorumID core.QuorumID,
	nonSignerMap map[core.OperatorID]*core.G1Point,
	aggregateSignersG1PubKey *core.G1Point,
) {
	var recomputedG1PubKeyAggregate *core.G1Point
	var recomputedSignerG1PubKeyAggregate *core.G1Point

	for operatorID := range sr.indexedOperatorState.Operators[quorumID] {
		operatorPubkey := sr.indexedOperatorState.IndexedOperators[operatorID].PubkeyG1

		if recomputedG1PubKeyAggregate == nil {
			recomputedG1PubKeyAggregate = operatorPubkey.Clone()
		} else {
			recomputedG1PubKeyAggregate.Add(operatorPubkey)
		}

		if _, ok := nonSignerMap[operatorID]; !ok {
			if recomputedSignerG1PubKeyAggregate == nil {
				recomputedSignerG1PubKeyAggregate = operatorPubkey.Clone()
			} else {
				recomputedSignerG1PubKeyAggregate.Add(operatorPubkey)
			}
		}
	}

	if recomputedG1PubKeyAggregate == nil {
		sr.logger.Error("recomputed aggregate G1 pubkey is nil. this shouldn't be possible")
	} else if !recomputedG1PubKeyAggregate.G1Affine.Equal(sr.indexedOperatorState.AggKeys[quorumID].G1Affine) {
		sr.logger.Error("recomputed aggregate G1 pubkey does not match indexed operator state aggregate G1 pubkey",
			"recomputedG1PubKeyAggregate", recomputedG1PubKeyAggregate.Serialize(),
			"indexedOperatorStateAggregateG1PubKey", sr.indexedOperatorState.AggKeys[quorumID].Serialize(),
			"quorumID", quorumID,
			"batchHeaderHash", hex.EncodeToString(sr.batchHeaderHash[:]))
	}

	if recomputedSignerG1PubKeyAggregate == nil {
		sr.logger.Error("recomputed aggregate signer G1 pubkey is nil. this shouldn't be possible")
	} else if !recomputedSignerG1PubKeyAggregate.G1Affine.Equal(aggregateSignersG1PubKey.G1Affine) {
		sr.logger.Error("recomputed aggregate signer G1 pubkey does not match key computed via subtraction",
			"recomputedSignerG1PubKeyAggregate", recomputedSignerG1PubKeyAggregate.Serialize(),
			"pubkeyComputedViaSubtraction", aggregateSignersG1PubKey.Serialize(),
		)
	}
}
