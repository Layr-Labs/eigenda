package controller_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/testutils"
	testrandom "github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	"github.com/stretchr/testify/require"
)

func createOperatorID(i int) core.OperatorID {
	var operatorID core.OperatorID
	copy(operatorID[:], fmt.Sprintf("operator-%d", i))
	return operatorID
}

func createBatchHeaderHash(testRandom *testrandom.TestRandom) [32]byte {
	return [32]byte(testRandom.Bytes(32))
}

func createSigningMessage(
	operatorID core.OperatorID,
	keypair *core.KeyPair,
	headerHash [32]byte,
	withError bool,
) core.SigningMessage {
	var err error
	if withError {
		err = errors.New("simulated error")
	}

	return core.SigningMessage{
		Signature:            keypair.SignMessage(headerHash),
		Operator:             operatorID,
		BatchHeaderHash:      headerHash,
		AttestationLatencyMs: 10.0,
		TimeReceived:         time.Now(),
		Err:                  err,
	}
}

func createIndexedOperatorState(
	t *testing.T,
	testRandom *testrandom.TestRandom,
	operatorCount int,
	quorumCount int,
) (*core.IndexedOperatorState, map[core.OperatorID]*core.KeyPair) {
	quorumOperatorInfo := make(map[core.QuorumID]*core.OperatorInfo)
	quorumOperators := make(map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo)
	quorumAggregatePubkeys := make(map[core.QuorumID]*core.G1Point)

	operatorKeys := make(map[core.OperatorID]*core.KeyPair)

	// create operators
	operatorInfo := make(map[core.OperatorID]*core.IndexedOperatorInfo)
	for i := 0; i < operatorCount; i++ {
		operatorID := createOperatorID(i)
		keypair, err := core.GenRandomBlsKeys()
		require.NoError(t, err)

		operatorKeys[operatorID] = keypair

		operatorInfo[operatorID] = &core.IndexedOperatorInfo{
			PubkeyG1: keypair.GetPubKeyG1(),
			PubkeyG2: keypair.GetPubKeyG2(),
			Socket:   "127.0.0.1:9000",
		}
	}

	// create quorums
	for quorumIndex := 0; quorumIndex < quorumCount; quorumIndex++ {
		quorumID := core.QuorumID(quorumIndex)
		quorumOperators[quorumID] = make(map[core.OperatorID]*core.OperatorInfo)
		quorumOperatorInfo[quorumID] = &core.OperatorInfo{
			Stake: big.NewInt(0),
			Index: 0,
		}

		operatorQuorumIndex := 0
		for operatorID, indexedOperatorInfo := range operatorInfo {
			// each operator has a 50% chance of being in a given quorum, except for operator 0, which is always in the
			// quorum. this is to guarantee that there is never an empty quorum
			if operatorID != createOperatorID(0) && testRandom.Bool() {
				continue
			}

			operatorStake := big.NewInt(testRandom.Int64Range(1, 1000))
			quorumOperators[quorumID][operatorID] = &core.OperatorInfo{
				Stake: operatorStake,
				Index: uint(operatorQuorumIndex),
			}
			quorumOperatorInfo[quorumID].Stake.Add(quorumOperatorInfo[quorumID].Stake, operatorStake)

			_, exists := quorumAggregatePubkeys[quorumID]
			if exists {
				quorumAggregatePubkeys[quorumID].Add(indexedOperatorInfo.PubkeyG1)
			} else {
				quorumAggregatePubkeys[quorumID] = indexedOperatorInfo.PubkeyG1.Clone()
			}

			operatorQuorumIndex++
		}
	}

	return &core.IndexedOperatorState{
		OperatorState: &core.OperatorState{
			Operators:   quorumOperators,
			Totals:      quorumOperatorInfo,
			BlockNumber: uint(testRandom.Uint32n(1000)),
		},
		IndexedOperators: operatorInfo,
		AggKeys:          quorumAggregatePubkeys,
	}, operatorKeys
}

func assertAttestationCorrectness(
	t *testing.T,
	attestationToVerify *core.QuorumAttestation,
	indexedOperatorState *core.IndexedOperatorState,
	operatorKeys map[core.OperatorID]*core.KeyPair,
	operatorSignatures map[core.OperatorID]*core.Signature,
) {
	for quorumID, quorumOperators := range indexedOperatorState.Operators {
		var expectedQuorumPubkeyAggregate *core.G1Point
		var expectedQuorumSignerPubkeyAggregate *core.G2Point
		var expectedQuorumSignatureAggregate *core.Signature
		expectedStakeSigned := uint64(0)
		for operatorID, operatorInfo := range quorumOperators {
			// pubkey of every operator is included, regardless of whether they signed or not
			if expectedQuorumPubkeyAggregate == nil {
				expectedQuorumPubkeyAggregate = operatorKeys[operatorID].GetPubKeyG1().Clone()
			} else {
				expectedQuorumPubkeyAggregate.Add(operatorKeys[operatorID].GetPubKeyG1())
			}

			if !attestationToVerify.SignerMap[operatorID] {
				// the rest of the aggregates are only for signers
				continue
			}

			if expectedQuorumSignerPubkeyAggregate == nil {
				expectedQuorumSignerPubkeyAggregate = operatorKeys[operatorID].GetPubKeyG2().Clone()
			} else {
				expectedQuorumSignerPubkeyAggregate.Add(operatorKeys[operatorID].GetPubKeyG2())
			}

			if expectedQuorumSignatureAggregate == nil {
				expectedQuorumSignatureAggregate = &core.Signature{G1Point: operatorSignatures[operatorID].Clone()}
			} else {
				expectedQuorumSignatureAggregate.Add(operatorSignatures[operatorID].G1Point)
			}

			expectedStakeSigned += operatorInfo.Stake.Uint64()

			_, actuallySigned := operatorSignatures[operatorID]
			require.True(t, actuallySigned)
		}

		expectedPercentSigned := uint8(expectedStakeSigned * 100 / indexedOperatorState.Totals[quorumID].Stake.Uint64())

		require.Equal(t, expectedQuorumPubkeyAggregate, attestationToVerify.QuorumAggPubKey[quorumID])
		require.Equal(t, expectedQuorumSignerPubkeyAggregate, attestationToVerify.SignersAggPubKey[quorumID])
		require.Equal(t, expectedQuorumSignatureAggregate, attestationToVerify.AggSignature[quorumID])
		require.Equal(t, expectedPercentSigned, attestationToVerify.QuorumResults[quorumID].PercentSigned)
		require.Equal(t, quorumID, attestationToVerify.QuorumResults[quorumID].QuorumID)
	}
}

// Test basic signature receiving functionality without concurrency
func TestReceiveSignatures_Basic(t *testing.T) {
	testRandom := testrandom.NewTestRandom()

	operatorCount := 3
	quorumCount := 2

	indexedOperatorState, operatorKeys := createIndexedOperatorState(t, testRandom, operatorCount, quorumCount)

	batchHeaderHash := createBatchHeaderHash(testRandom)
	signingMessageChan := make(chan core.SigningMessage, 3)

	attestationChan, err := controller.ReceiveSignatures(
		context.Background(),
		testutils.GetLogger(),
		nil,
		indexedOperatorState,
		batchHeaderHash,
		signingMessageChan,
		50*time.Millisecond,
		55)
	require.NoError(t, err)

	// send signing messages from each operator
	operatorSignatures := make(map[core.OperatorID]*core.Signature)
	for operatorID := range indexedOperatorState.IndexedOperators {
		signingMessage := createSigningMessage(operatorID, operatorKeys[operatorID], batchHeaderHash, false)
		signingMessageChan <- signingMessage
		operatorSignatures[operatorID] = signingMessage.Signature
	}

	for attestation := range attestationChan {
		assertAttestationCorrectness(t, attestation, indexedOperatorState, operatorKeys, operatorSignatures)
	}
}

// Test receiving signatures with an error in one of the signing messages
func TestReceiveSignatures_WithError(t *testing.T) {
	testRandom := testrandom.NewTestRandom()

	operatorCount := 3
	quorumCount := 2

	indexedOperatorState, operatorKeys := createIndexedOperatorState(t, testRandom, operatorCount, quorumCount)

	batchHeaderHash := createBatchHeaderHash(testRandom)
	signingMessageChan := make(chan core.SigningMessage, operatorCount)

	attestationChan, err := controller.ReceiveSignatures(
		context.Background(),
		testutils.GetLogger(),
		nil,
		indexedOperatorState,
		batchHeaderHash,
		signingMessageChan,
		50*time.Millisecond,
		55)
	require.NoError(t, err)

	// Send signing messages with one error
	operatorSignatures := make(map[core.OperatorID]*core.Signature)
	for operatorID := range indexedOperatorState.IndexedOperators {
		withError := operatorID == createOperatorID(0)
		signingMessage := createSigningMessage(operatorID, operatorKeys[operatorID], batchHeaderHash, withError)
		signingMessageChan <- signingMessage
		if !withError {
			operatorSignatures[operatorID] = signingMessage.Signature
		}
	}

	for attestation := range attestationChan {
		assertAttestationCorrectness(t, attestation, indexedOperatorState, operatorKeys, operatorSignatures)
	}
}

// Test behavior when receiving duplicate signing messages
func TestReceiveSignatures_DuplicateMessage(t *testing.T) {
	testRandom := testrandom.NewTestRandom()

	operatorCount := 3
	quorumCount := 2

	indexedOperatorState, operatorKeys := createIndexedOperatorState(t, testRandom, operatorCount, quorumCount)

	batchHeaderHash := createBatchHeaderHash(testRandom)
	signingMessageChan := make(chan core.SigningMessage, operatorCount+1) // One extra for duplicate

	attestationChan, err := controller.ReceiveSignatures(
		context.Background(),
		testutils.GetLogger(),
		nil,
		indexedOperatorState,
		batchHeaderHash,
		signingMessageChan,
		50*time.Millisecond,
		55)
	require.NoError(t, err)

	// Send signing messages from each operator
	operatorSignatures := make(map[core.OperatorID]*core.Signature)
	for operatorID := range indexedOperatorState.IndexedOperators {
		signingMessage := createSigningMessage(operatorID, operatorKeys[operatorID], batchHeaderHash, false)
		signingMessageChan <- signingMessage
		operatorSignatures[operatorID] = signingMessage.Signature

		// send one duplicate
		if operatorID == createOperatorID(0) {
			signingMessage := createSigningMessage(operatorID, operatorKeys[operatorID], batchHeaderHash, false)
			signingMessageChan <- signingMessage
		}
	}

	for attestation := range attestationChan {
		assertAttestationCorrectness(t, attestation, indexedOperatorState, operatorKeys, operatorSignatures)
	}
}

// Test context cancellation behavior
func TestReceiveSignatures_ContextCancellation(t *testing.T) {
	testRandom := testrandom.NewTestRandom()

	operatorCount := 3
	quorumCount := 2

	indexedOperatorState, operatorKeys := createIndexedOperatorState(t, testRandom, operatorCount, quorumCount)

	batchHeaderHash := createBatchHeaderHash(testRandom)
	signingMessageChan := make(chan core.SigningMessage, operatorCount)

	ctx, cancel := context.WithCancel(context.Background())
	attestationChan, err := controller.ReceiveSignatures(
		ctx,
		testutils.GetLogger(),
		nil,
		indexedOperatorState,
		batchHeaderHash,
		signingMessageChan,
		50*time.Millisecond,
		55)
	require.NoError(t, err)

	// Send only 1 signing message
	operatorSignatures := make(map[core.OperatorID]*core.Signature)
	operatorID := createOperatorID(0)
	signingMessage := createSigningMessage(operatorID, operatorKeys[operatorID], batchHeaderHash, false)
	signingMessageChan <- signingMessage
	operatorSignatures[operatorID] = signingMessage.Signature

	attestation := <-attestationChan

	cancel()

	assertAttestationCorrectness(t, attestation, indexedOperatorState, operatorKeys, operatorSignatures)
}

// Test concurrent signature receiving with a large number of operators
func TestReceiveSignatures_Concurrency(t *testing.T) {
	testRandom := testrandom.NewTestRandom()

	const operatorCount = 100
	const quorumCount = 10
	const errorProbability = 0.05
	const invalidSignatureProbability = 0.05

	indexedOperatorState, operatorKeys := createIndexedOperatorState(t, testRandom, operatorCount, quorumCount)

	batchHeaderHash := createBatchHeaderHash(testRandom)
	signingMessageChan := make(chan core.SigningMessage, operatorCount)

	attestationChan, err := controller.ReceiveSignatures(
		context.Background(),
		testutils.GetLogger(),
		nil,
		indexedOperatorState,
		batchHeaderHash,
		signingMessageChan,
		1*time.Millisecond,
		55)
	require.NoError(t, err)

	attestationCount := atomic.Int32{}

	operatorSignatures := make(map[core.OperatorID]*core.Signature)
	signatureMapMutex := sync.Mutex{}

	// Start a goroutine to collect attestations
	attestationsDone := make(chan struct{})
	go func() {
		for attestation := range attestationChan {
			attestationCount.Add(1)

			signatureMapMutex.Lock()
			assertAttestationCorrectness(
				t,
				attestation,
				indexedOperatorState,
				operatorKeys,
				operatorSignatures)
			signatureMapMutex.Unlock()
		}

		attestationsDone <- struct{}{}
	}()

	for operatorID := range indexedOperatorState.IndexedOperators {
		boundID := operatorID
		go func() {
			time.Sleep(time.Duration(testRandom.Uint32n(10)) * time.Millisecond)

			// some signing messages will contain an error
			withError := testRandom.Float64() < errorProbability

			hashToSign := batchHeaderHash
			// some signing messages will be invalid
			if testRandom.Float64() < invalidSignatureProbability {
				hashToSign = createBatchHeaderHash(testRandom)
			}

			signingMessage := createSigningMessage(boundID, operatorKeys[boundID], hashToSign, withError)
			signingMessageChan <- signingMessage

			if !withError && hashToSign == batchHeaderHash {
				signatureMapMutex.Lock()
				defer signatureMapMutex.Unlock()
				operatorSignatures[boundID] = signingMessage.Signature
			}
		}()
	}

	// Wait for all attestations to be processed
	<-attestationsDone

	require.Greater(t, attestationCount.Load(), int32(1), "Should have received multiple attestations")
}
