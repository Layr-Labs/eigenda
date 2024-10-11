package corev2_test

import (
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/chainio/mock"
	corev2 "github.com/Layr-Labs/eigenda/corev2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

var (
	dat *mock.ChainDataMock
	agg corev2.SignatureAggregator

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
	reader := &mock.MockReader{}
	reader.On("OperatorIDToAddress").Return(gethcommon.Address{}, nil)
	agg, err = corev2.NewStdSignatureAggregator(logger, reader)
	if err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}
