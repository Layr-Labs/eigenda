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
	dat, err = mock.MakeChainDataMock(10)
	if err != nil {
		panic(err)
	}
	logger := logging.NewNoopLogger()
	transactor := &mock.MockTransactor{}
	transactor.On("OperatorIDToAddress").Return(gethcommon.Address{}, nil)
	agg, err = core.NewStdSignatureAggregator(logger, transactor)
	if err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}

func simulateOperators(state mock.PrivateOperatorState, message [32]byte, update chan core.SignerMessage, advCount uint) {

	count := 0

	// Simulate the operators signing the message.
	// In real life, the ordering will be random, but we simulate the signing in a fixed order
	// to simulate stakes deterministically
	for i := 0; i < len(state.PrivateOperators); i++ {
		id := makeOperatorId(i)
		op := state.PrivateOperators[id]
		sig := op.KeyPair.SignMessage(message)
		if count < len(state.IndexedOperators)-int(advCount) {
			update <- core.SignerMessage{
				Signature: sig,
				Operator:  id,
				Err:       nil,
			}
		} else {
			update <- core.SignerMessage{
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
		meetsQuorum    bool
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
			meetsQuorum:    true,
		},
		{
			name: "Succeeds when 9/10 operators sign at quorum threshold 80",
			quorums: []core.QuorumResult{
				{
					QuorumID:      0,
					PercentSigned: 80,
				},
			},
			adversaryCount: 1,
			expectedErr:    nil,
			meetsQuorum:    true,
		},
		{
			name: "Fails when 8/10 operators sign at quorum threshold 90",
			quorums: []core.QuorumResult{
				{
					QuorumID:      0,
					PercentSigned: 90,
				},
			},
			adversaryCount: 2,
			expectedErr:    nil,
			meetsQuorum:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			state := dat.GetTotalOperatorState(context.Background(), 0)

			update := make(chan core.SignerMessage)
			message := [32]byte{1, 2, 3, 4, 5, 6}

			go simulateOperators(*state, message, update, tt.adversaryCount)

			quorumIDs := make([]core.QuorumID, len(tt.quorums))
			for ind, quorum := range tt.quorums {
				quorumIDs[ind] = quorum.QuorumID
			}

			sigAgg, err := agg.AggregateSignatures(context.Background(), state.IndexedOperatorState, quorumIDs, message, update)
			assert.NoError(t, err)

			for _, quorum := range tt.quorums {
				if tt.meetsQuorum {
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

	update := make(chan core.SignerMessage)
	message := [32]byte{1, 2, 3, 4, 5, 6}

	go simulateOperators(*state, message, update, 4)

	quorums := []core.QuorumID{0}

	sigAgg, err := agg.AggregateSignatures(context.Background(), state.IndexedOperatorState, quorums, message, update)
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
