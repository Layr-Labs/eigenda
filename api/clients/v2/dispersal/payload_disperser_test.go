package dispersal

import (
	"math/big"
	"testing"

	dispgrpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestVerifyReceivedBlobKey(t *testing.T) {
	blobCommitments := encoding.BlobCommitments{
		Commitment:       &encoding.G1Commitment{},
		LengthCommitment: &encoding.G2Commitment{},
		LengthProof:      &encoding.LengthProof{},
		Length:           4,
	}

	quorumNumbers := make([]core.QuorumID, 1)
	quorumNumbers[0] = 8

	paymentMetadata := core.PaymentMetadata{
		AccountID:         gethcommon.Address{1},
		Timestamp:         5,
		CumulativePayment: big.NewInt(6),
	}

	blobHeader := &corev2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: blobCommitments,
		QuorumNumbers:   quorumNumbers,
		PaymentMetadata: paymentMetadata,
	}

	realKey, err := blobHeader.BlobKey()
	require.NoError(t, err)

	reply := dispgrpc.DisperseBlobReply{
		BlobKey: realKey[:],
	}

	verifiedKey, err := verifyReceivedBlobKey(blobHeader, &reply)
	require.NoError(t, err)
	require.Equal(t, realKey, verifiedKey)

	blobHeader.BlobVersion = 1
	_, err = verifyReceivedBlobKey(blobHeader, &reply)
	require.Error(t, err, "Any modification to the header should cause verification to fail")
}
