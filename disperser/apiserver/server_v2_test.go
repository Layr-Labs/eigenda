package apiserver_test

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/disperser/apiserver"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigensdk-go/logging"

	pbcommon "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	pbv2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/stretchr/testify/assert"
)

func TestV2DisperseBlob(t *testing.T) {
	data := make([]byte, 3*1024)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	data = codec.ConvertByPaddingEmptyByte(data)
	_, err = dispersalServerV2.DisperseBlob(context.Background(), &pbv2.DisperseBlobRequest{
		Data:       data,
		BlobHeader: &pbcommon.BlobHeader{},
	})
	assert.ErrorContains(t, err, "not implemented")
}

func TestV2GetBlobStatus(t *testing.T) {
	_, err := dispersalServerV2.GetBlobStatus(context.Background(), &pbv2.BlobStatusRequest{
		BlobKey: []byte{1},
	})
	assert.ErrorContains(t, err, "not implemented")
}

func TestV2GetBlobCommitment(t *testing.T) {
	_, err := dispersalServerV2.GetBlobCommitment(context.Background(), &pbv2.BlobCommitmentRequest{
		Data: []byte{1},
	})
	assert.ErrorContains(t, err, "not implemented")
}

func newTestServerV2() *apiserver.DispersalServerV2 {
	logger := logging.NewNoopLogger()
	return apiserver.NewDispersalServerV2(disperser.ServerConfig{
		GrpcPort:    "51002",
		GrpcTimeout: 1 * time.Second,
	}, logger)
}
