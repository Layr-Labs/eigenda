package clients

import (
	"math/big"
	"testing"

	v2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/stretchr/testify/require"
)

func TestVerifyReceivedBlobKey(t *testing.T) {
	blobCommitments := encoding.BlobCommitments{
		Commitment: &encoding.G1Commitment{},
		LengthCommitment: &encoding.G2Commitment{},
		LengthProof: &encoding.LengthProof{},
		Length: 4,
	}

	quorumNumbers := make([]core.QuorumID, 1)
	quorumNumbers[0] = 8

	paymentMetadata := core.PaymentMetadata{
		AccountID: "asdf",
		ReservationPeriod: 5,
		CumulativePayment: big.NewInt(6),
	}

	blobHeader := &corev2.BlobHeader{
		BlobVersion: 0,
		BlobCommitments: blobCommitments,
		QuorumNumbers: quorumNumbers,
		PaymentMetadata: paymentMetadata,
		Salt: 9,
	}

	realKey, err := blobHeader.BlobKey()
	require.NoError(t, err)

	reply := v2.DisperseBlobReply{
		BlobKey: realKey[:],
	}

	require.NoError(t, verifyReceivedBlobKey(blobHeader, &reply))

	blobHeader.BlobVersion = 1
	require.Error(t, verifyReceivedBlobKey(blobHeader, &reply),
		"Any modification to the header should cause verification to fail")
}
