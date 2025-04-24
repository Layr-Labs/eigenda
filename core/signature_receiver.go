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

// TODO
// This struct is not threadsafe. it's only meant to be used once
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

func ReceiveSignatures(
	ctx context.Context,
	logger logging.Logger,
	indexedOperatorState *IndexedOperatorState,
	batchHeaderHash [32]byte,
	signingMessageChan chan SigningMessage,
) (chan *QuorumAttestation, error) {
	quorumIDs, err := getSortedQuorumIDs(indexedOperatorState)
	if err != nil {
		return nil, fmt.Errorf("get sorted quorum ids: %w", err)
	}

	aggregateSignatures := make(map[QuorumID]*Signature, len(quorumIDs))
	aggregateSignersG2PubKeys := make(map[QuorumID]*G2Point, len(quorumIDs))

	// initialized stakeSigned map with 0 stake signed for each quorum
	stakeSigned := make(map[QuorumID]*big.Int, len(quorumIDs))
	for _, quorumID := range quorumIDs {
		stakeSigned[quorumID] = big.NewInt(0)
	}

	receiver := &signatureReceiver{
		logger:                    logger,
		indexedOperatorState:      indexedOperatorState,
		aggregateSignatures:       aggregateSignatures,
		aggregateSignersG2PubKeys: aggregateSignersG2PubKeys,
		stakeSigned:               stakeSigned,
		batchHeaderHash:           batchHeaderHash,
		signingMessageChan:        signingMessageChan,
		quorumIDs:                 quorumIDs,
	}

	attestationChan := make(chan *QuorumAttestation)
	receiver.receiveSigningMessages(ctx, attestationChan, hex.EncodeToString(batchHeaderHash[:]))

	return attestationChan, nil
}

func (sr *signatureReceiver) receiveSigningMessages(
	ctx context.Context,
	attestationChan chan *QuorumAttestation,
	batchHeaderHashHex string,
) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer close(attestationChan)

	operatorCount := len(sr.indexedOperatorState.IndexedOperators)
	signingMessageCount := 0
	errorCount := 0
	for signingMessageCount < operatorCount {
		var signingMessage SigningMessage
		select {
		case <-ctx.Done():
			sr.logger.Infof(
				`global batch attestation timeout exceeded for batch %s. Recieved and processed %d/%d signing
						messages. %d of the signing messages caused an error during processing`,
				batchHeaderHashHex, signingMessageCount, operatorCount, errorCount)
		case signingMessage = <-sr.signingMessageChan:
			signingMessageCount++
			err := sr.processSigningMessage(signingMessage)
			if err != nil {
				errorCount++
				sr.logger.Warn("process signing message",
					"batchHeaderHash", batchHeaderHashHex,
					"operatorID", signingMessage.Operator.Hex(),
					"attestationLatencyMs", signingMessage.AttestationLatencyMs,
					"error", err)
				continue
			}
		case <-ticker.C:
			quorumResults := sr.computeQuorumResults()

			attestationChan <- &QuorumAttestation{
				// TODO: is this ok? semantics are changed from before: we used to exclude aggregate keys of quorums that had no
				//  signatures, but I don't see why that case should be special
				QuorumAggPubKey:  sr.indexedOperatorState.AggKeys,
				SignersAggPubKey: sr.aggregateSignersG2PubKeys,
				AggSignature:     sr.aggregateSignatures,
				QuorumResults:    quorumResults,
				SignerMap:        sr.signerMap,
			}
		}
	}
}

func (sr *signatureReceiver) computeQuorumResults() map[QuorumID]*QuorumResult {
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
			sr.logger.Warn("compute quorum result failed",
				"quorumID", quorumID, "batchHeaderHash", sr.batchHeaderHash)
			continue
		}
		quorumResults[quorumID] = quorumResult
	}

	return quorumResults
}

func (sr *signatureReceiver) computeQuorumResult(
	quorumID QuorumID,
	nonSignerMap map[OperatorID]*G1Point,
) (*QuorumResult, error) {
	signedPercentage := getSignedPercentage(
		sr.indexedOperatorState.Totals[quorumID].Stake, sr.stakeSigned[quorumID])

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
		return nil, errors.New("invalid aggregate signer G2 public key")
	}

	ok, err := aggregateSignersG1PubKey.VerifyEquivalence(sr.aggregateSignersG2PubKeys[quorumID])
	if err != nil {
		return nil, fmt.Errorf("verify pubkey equivalence: %w", err)
	}
	if !ok {
		return nil, errors.New("aggregate signers G1 pubkey is not equivalent to aggregate signers G2 pubkey")
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

func (sr *signatureReceiver) processSigningMessage(signingMessage SigningMessage) error {
	indexedOperatorInfo, err := sr.checkSigningMessage(signingMessage)
	if err != nil {
		return fmt.Errorf("check signing message: %w", err)
	}

	sr.signerMap[signingMessage.Operator] = true

	for _, quorumID := range sr.quorumIDs {
		quorumOperators := sr.indexedOperatorState.Operators[quorumID]
		quorumOperatorInfo, isOperatorInQuorum := quorumOperators[signingMessage.Operator]
		// if the operator which sent the signing message isn't in a given quorum, then we shouldn't make any
		// changes to the aggregates that are tracked on a per-quorum basis
		if !isOperatorInQuorum {
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

func getSignedPercentage(totalStake *big.Int, signedStake *big.Int) uint8 {
	if totalStake.Cmp(big.NewInt(0)) == 0 {
		return 0
	}

	signedStake = signedStake.Mul(signedStake, new(big.Int).SetUint64(percentMultiplier))
	quorumThresholdBig := signedStake.Div(signedStake, totalStake)

	quorumThreshold := uint8(quorumThresholdBig.Uint64())

	return quorumThreshold
}
