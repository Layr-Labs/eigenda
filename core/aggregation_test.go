package core_test

import (
	"context"
	"errors"
	"math/big"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

var (
	dat *mock.ChainDataMock
	agg core.SignatureAggregator

	GETTYSBURG_ADDRESS_BYTES = []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")
)

func TestMain(m *testing.M) {
	var err error
	dat, err = mock.MakeChainDataMock(map[uint8]int{
		0: 6,
		1: 3,
	})
	if err != nil {
		panic(err)
	}
	logger := logging.NewNoopLogger()
	transactor := &mock.MockWriter{}
	transactor.On("OperatorIDToAddress").Return(gethcommon.Address{}, nil)
	agg, err = core.NewStdSignatureAggregator(logger, transactor)
	if err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}

func simulateOperators(state mock.PrivateOperatorState, message [32]byte, update chan core.SigningMessage, advCount uint) {

	count := 0

	// Simulate the operators signing the message.
	// In real life, the ordering will be random, but we simulate the signing in a fixed order
	// to simulate stakes deterministically
	for i := 0; i < len(state.PrivateOperators); i++ {
		id := mock.MakeOperatorId(i)
		op := state.PrivateOperators[id]
		sig := op.KeyPair.SignMessage(message)
		if count < len(state.IndexedOperators)-int(advCount) {
			update <- core.SigningMessage{
				Signature: sig,
				Operator:  id,
				Err:       nil,
			}
		} else {
			update <- core.SigningMessage{
				Signature: nil,
				Operator:  id,
				Err:       errors.New("adversary"),
			}
		}

		count += 1
	}
}

func TestAggregateSignaturesStatus(t *testing.T) {

	tests := []struct {
		name           string
		quorums        []core.QuorumResult
		adversaryCount uint
		expectedErr    error
		meetsQuorum    []bool
	}{
		{
			name: "Succeeds when all operators sign at quorum threshold 100",
			quorums: []core.QuorumResult{
				{
					QuorumID:      0,
					PercentSigned: 100,
				},
			},
			adversaryCount: 0,
			expectedErr:    nil,
			meetsQuorum:    []bool{true},
		},
		{
			name: "Succeeds when 5/6 operators sign at quorum threshold 70",
			quorums: []core.QuorumResult{
				{
					QuorumID:      0,
					PercentSigned: 70,
				},
			},
			adversaryCount: 1,
			expectedErr:    nil,
			meetsQuorum:    []bool{true},
		},
		{
			name: "Fails when 4/6 operators sign at quorum threshold 90",
			quorums: []core.QuorumResult{
				{
					QuorumID:      0,
					PercentSigned: 90,
				},
			},
			adversaryCount: 2,
			expectedErr:    nil,
			meetsQuorum:    []bool{false},
		},
		{
			name: "Fails when 5/6 operators sign at quorum threshold 80 for 2 quorums",
			quorums: []core.QuorumResult{
				{
					QuorumID:      0,
					PercentSigned: 80,
				},
				{
					QuorumID:      1,
					PercentSigned: 80,
				},
			},
			adversaryCount: 1,
			expectedErr:    nil,
			meetsQuorum:    []bool{false, true},
		},
		{
			name: "Succeeds when 5/6 operators sign at quorum threshold 70 and 100",
			quorums: []core.QuorumResult{
				{
					QuorumID:      0,
					PercentSigned: 70,
				},
				{
					QuorumID:      1,
					PercentSigned: 100,
				},
			},
			adversaryCount: 1,
			expectedErr:    nil,
			meetsQuorum:    []bool{true, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := dat.GetTotalOperatorStateWithQuorums(context.Background(), 0, []core.QuorumID{0, 1})

			update := make(chan core.SigningMessage)
			message := [32]byte{1, 2, 3, 4, 5, 6}

			go simulateOperators(*state, message, update, tt.adversaryCount)

			quorumIDs := make([]core.QuorumID, len(tt.quorums))
			for ind, quorum := range tt.quorums {
				quorumIDs[ind] = quorum.QuorumID
			}

			numOpr := 0
			for _, quorum := range tt.quorums {
				if len(dat.Stakes[quorum.QuorumID]) > numOpr {
					numOpr = len(dat.Stakes[quorum.QuorumID])
				}
			}

			aq, err := agg.ReceiveSignatures(context.Background(), state.IndexedOperatorState, message, update)
			assert.NoError(t, err)
			assert.Len(t, aq.SignerMap, numOpr-int(tt.adversaryCount))
			assert.Len(t, aq.AggSignature, 2)
			assert.Len(t, aq.QuorumAggPubKey, 2)
			assert.Len(t, aq.SignersAggPubKey, 2)
			assert.Len(t, aq.QuorumResults, 2)
			for i, q := range tt.quorums {
				assert.NotNil(t, aq.AggSignature[q.QuorumID])
				assert.NotNil(t, aq.QuorumAggPubKey[q.QuorumID])
				assert.NotNil(t, aq.SignersAggPubKey[q.QuorumID])
				if tt.meetsQuorum[i] {
					assert.GreaterOrEqual(t, aq.QuorumResults[q.QuorumID].PercentSigned, q.PercentSigned)
				} else {
					assert.Less(t, aq.QuorumResults[q.QuorumID].PercentSigned, q.PercentSigned)
				}
			}

			sigAgg, err := agg.AggregateSignatures(context.Background(), dat, 0, aq, quorumIDs)
			assert.NoError(t, err)

			for i, quorum := range tt.quorums {
				if tt.meetsQuorum[i] {
					assert.GreaterOrEqual(t, sigAgg.QuorumResults[quorum.QuorumID].PercentSigned, quorum.PercentSigned)
				} else {
					assert.Less(t, sigAgg.QuorumResults[quorum.QuorumID].PercentSigned, quorum.PercentSigned)
				}
			}
		})
	}

}

func TestSortNonsigners(t *testing.T) {
	state := dat.GetTotalOperatorState(context.Background(), 0)

	update := make(chan core.SigningMessage)
	message := [32]byte{1, 2, 3, 4, 5, 6}

	go simulateOperators(*state, message, update, 4)

	quorums := []core.QuorumID{0}

	aq, err := agg.ReceiveSignatures(context.Background(), state.IndexedOperatorState, message, update)
	assert.NoError(t, err)
	sigAgg, err := agg.AggregateSignatures(context.Background(), dat, 0, aq, quorums)
	assert.NoError(t, err)

	for i := range sigAgg.NonSigners {
		if i == 0 {
			continue
		}
		prevHash := sigAgg.NonSigners[i-1].Hash()
		currHash := sigAgg.NonSigners[i].Hash()
		prevHashInt := new(big.Int).SetBytes(prevHash[:])
		currHashInt := new(big.Int).SetBytes(currHash[:])
		assert.Equal(t, currHashInt.Cmp(prevHashInt), 1)
	}
}

func TestFilterQuorums(t *testing.T) {
	allQuorums := []core.QuorumID{0, 1}
	state := dat.GetTotalOperatorStateWithQuorums(context.Background(), 0, allQuorums)

	update := make(chan core.SigningMessage)
	message := [32]byte{1, 2, 3, 4, 5, 6}
	advCount := 4
	go simulateOperators(*state, message, update, uint(advCount))

	numOpr := 0
	for _, quorum := range allQuorums {
		if len(dat.Stakes[quorum]) > numOpr {
			numOpr = len(dat.Stakes[quorum])
		}
	}

	aq, err := agg.ReceiveSignatures(context.Background(), state.IndexedOperatorState, message, update)
	assert.NoError(t, err)
	assert.Len(t, aq.SignerMap, numOpr-advCount)
	assert.Equal(t, aq.SignerMap, map[core.OperatorID]bool{
		mock.MakeOperatorId(0): true,
		mock.MakeOperatorId(1): true,
	})
	assert.Contains(t, aq.AggSignature, core.QuorumID(0))
	assert.Contains(t, aq.AggSignature, core.QuorumID(1))
	assert.Equal(t, aq.QuorumAggPubKey, map[core.QuorumID]*core.G1Point{
		core.QuorumID(0): state.IndexedOperatorState.AggKeys[0],
		core.QuorumID(1): state.IndexedOperatorState.AggKeys[1],
	})
	aggSignerPubKey0 := state.IndexedOperatorState.IndexedOperators[mock.MakeOperatorId(0)].PubkeyG2.Clone()
	aggSignerPubKey0.Add(state.IndexedOperatorState.IndexedOperators[mock.MakeOperatorId(1)].PubkeyG2)
	aggSignerPubKey1 := state.IndexedOperatorState.IndexedOperators[mock.MakeOperatorId(0)].PubkeyG2.Clone()
	aggSignerPubKey1.Add(state.IndexedOperatorState.IndexedOperators[mock.MakeOperatorId(1)].PubkeyG2)
	assert.Contains(t, aq.SignersAggPubKey, core.QuorumID(0))
	assert.Equal(t, aq.SignersAggPubKey[core.QuorumID(0)], aggSignerPubKey0)
	assert.Contains(t, aq.SignersAggPubKey, core.QuorumID(1))
	assert.Equal(t, aq.SignersAggPubKey[core.QuorumID(1)], aggSignerPubKey1)
	assert.Equal(t, aq.QuorumResults[core.QuorumID(0)].PercentSigned, uint8(14))
	assert.Equal(t, aq.QuorumResults[core.QuorumID(1)].PercentSigned, uint8(50))

	// Only consider quorum 0
	quorums := []core.QuorumID{0}
	sigAgg, err := agg.AggregateSignatures(context.Background(), dat, 0, aq, quorums)
	assert.NoError(t, err)
	assert.Len(t, sigAgg.NonSigners, 4)
	assert.ElementsMatch(t, sigAgg.NonSigners, []*core.G1Point{
		state.IndexedOperatorState.IndexedOperators[mock.MakeOperatorId(2)].PubkeyG1,
		state.IndexedOperatorState.IndexedOperators[mock.MakeOperatorId(3)].PubkeyG1,
		state.IndexedOperatorState.IndexedOperators[mock.MakeOperatorId(4)].PubkeyG1,
		state.IndexedOperatorState.IndexedOperators[mock.MakeOperatorId(5)].PubkeyG1,
	})
	assert.Len(t, sigAgg.QuorumAggPubKeys, 1)
	assert.Contains(t, sigAgg.QuorumAggPubKeys, core.QuorumID(0))
	assert.Equal(t, sigAgg.QuorumAggPubKeys[0], state.IndexedOperatorState.AggKeys[0])

	assert.Equal(t, sigAgg.AggPubKey, aggSignerPubKey0)
	expectedAggSignerKey := sigAgg.QuorumAggPubKeys[0].Clone()
	for _, nsk := range sigAgg.NonSigners {
		expectedAggSignerKey.Sub(nsk)
	}
	ok, err := expectedAggSignerKey.VerifyEquivalence(sigAgg.AggPubKey)
	assert.NoError(t, err)
	assert.True(t, ok)
	ok = sigAgg.AggSignature.Verify(sigAgg.AggPubKey, message)
	assert.True(t, ok)
	assert.Len(t, sigAgg.QuorumResults, 1)
	assert.Contains(t, sigAgg.QuorumResults, core.QuorumID(0))
	assert.Equal(t, sigAgg.QuorumResults[0].QuorumID, core.QuorumID(0))
	assert.Equal(t, sigAgg.QuorumResults[0].PercentSigned, core.QuorumID(14))

	// Only consider quorum 1
	quorums = []core.QuorumID{1}
	sigAgg, err = agg.AggregateSignatures(context.Background(), dat, 0, aq, quorums)
	assert.NoError(t, err)
	assert.Len(t, sigAgg.NonSigners, 1)
	assert.ElementsMatch(t, sigAgg.NonSigners, []*core.G1Point{
		state.IndexedOperatorState.IndexedOperators[mock.MakeOperatorId(2)].PubkeyG1,
	})
	assert.Len(t, sigAgg.QuorumAggPubKeys, 1)
	assert.Contains(t, sigAgg.QuorumAggPubKeys, core.QuorumID(1))
	assert.Equal(t, sigAgg.QuorumAggPubKeys[1], state.IndexedOperatorState.AggKeys[1])

	assert.Equal(t, sigAgg.AggPubKey, aggSignerPubKey1)
	expectedAggSignerKey = sigAgg.QuorumAggPubKeys[1].Clone()
	for _, nsk := range sigAgg.NonSigners {
		expectedAggSignerKey.Sub(nsk)
	}
	ok, err = expectedAggSignerKey.VerifyEquivalence(sigAgg.AggPubKey)
	assert.NoError(t, err)
	assert.True(t, ok)
	ok = sigAgg.AggSignature.Verify(sigAgg.AggPubKey, message)
	assert.True(t, ok)
	assert.Len(t, sigAgg.QuorumResults, 1)
	assert.Contains(t, sigAgg.QuorumResults, core.QuorumID(1))
	assert.Equal(t, sigAgg.QuorumResults[1].QuorumID, core.QuorumID(1))
	assert.Equal(t, sigAgg.QuorumResults[1].PercentSigned, core.QuorumID(50))
}
