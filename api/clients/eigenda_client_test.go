package clients_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	clientsmock "github.com/Layr-Labs/eigenda/api/clients/mock"
	"github.com/Layr-Labs/eigenda/api/grpc/common"
	grpcdisperser "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPutRetrieveBlobIFFTSuccess(t *testing.T) {
	disperserClient := clientsmock.NewMockDisperserClient()
	expectedBlobStatus := disperser.Processing
	(disperserClient.On("DisperseBlobAuthenticated", mock.Anything, mock.Anything, mock.Anything).
		Return(&expectedBlobStatus, []byte("mock-request-id"), nil))
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_PROCESSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_DISPERSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_CONFIRMED}, nil).Once())
	finalizedBlobInfo := &grpcdisperser.BlobInfo{
		BlobHeader: &grpcdisperser.BlobHeader{
			Commitment: &common.G1Commitment{X: []byte{0x00, 0x00, 0x00, 0x00}, Y: []byte{0x01, 0x00, 0x00, 0x00}},
			BlobQuorumParams: []*grpcdisperser.BlobQuorumParam{
				{
					QuorumNumber: 0,
				},
				{
					QuorumNumber: 1,
				},
			},
		},
		BlobVerificationProof: &grpcdisperser.BlobVerificationProof{
			BlobIndex: 100,
			BatchMetadata: &grpcdisperser.BatchMetadata{
				BatchHeaderHash: []byte("mock-batch-header-hash"),
				BatchHeader: &grpcdisperser.BatchHeader{
					ReferenceBlockNumber: 200,
				},
			},
		},
	}
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_FINALIZED, Info: finalizedBlobInfo}, nil).Once())
	(disperserClient.On("RetrieveBlob", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil).Once()) // pass nil in as the return blob to tell the mock to return the corresponding blob
	logger := log.NewLogger(log.DiscardHandler())
	eigendaClient := clients.EigenDAClient{
		Log: logger,
		Config: clients.EigenDAClientConfig{
			RPC:                          "localhost:51001",
			StatusQueryTimeout:           10 * time.Minute,
			StatusQueryRetryInterval:     50 * time.Millisecond,
			ResponseTimeout:              10 * time.Second,
			CustomQuorumIDs:              []uint{},
			SignerPrivateKeyHex:          "75f9e29cac7f5774d106adb355ef294987ce39b7863b75bb3f2ea42ca160926d",
			DisableTLS:                   false,
			PutBlobEncodingVersion:       codecs.DefaultBlobEncoding,
			DisablePointVerificationMode: false,
			WaitForFinalization:          true,
		},
		Client: disperserClient,
		Codec:  codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec()),
	}
	expectedBlob := []byte("dc49e7df326cfb2e7da5cf68f263e1898443ec2e862350606e7dfbda55ad10b5d61ed1d54baf6ae7a86279c1b4fa9c49a7de721dacb211264c1f5df31bade51c")
	blobInfo, err := eigendaClient.PutBlob(context.Background(), expectedBlob)
	require.NoError(t, err)
	require.NotNil(t, blobInfo)
	assert.Equal(t, finalizedBlobInfo, blobInfo)

	resultBlob, err := eigendaClient.GetBlob(context.Background(), []byte("mock-batch-header-hash"), 100)
	require.NoError(t, err)
	require.Equal(t, expectedBlob, resultBlob)
}

func TestPutRetrieveBlobIFFTNoDecodeSuccess(t *testing.T) {
	disperserClient := clientsmock.NewMockDisperserClient()
	expectedBlobStatus := disperser.Processing
	(disperserClient.On("DisperseBlobAuthenticated", mock.Anything, mock.Anything, mock.Anything).
		Return(&expectedBlobStatus, []byte("mock-request-id"), nil))
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_PROCESSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_DISPERSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_CONFIRMED}, nil).Once())
	finalizedBlobInfo := &grpcdisperser.BlobInfo{
		BlobHeader: &grpcdisperser.BlobHeader{
			Commitment: &common.G1Commitment{X: []byte{0x00, 0x00, 0x00, 0x00}, Y: []byte{0x01, 0x00, 0x00, 0x00}},
			BlobQuorumParams: []*grpcdisperser.BlobQuorumParam{
				{
					QuorumNumber: 0,
				},
				{
					QuorumNumber: 1,
				},
			},
		},
		BlobVerificationProof: &grpcdisperser.BlobVerificationProof{
			BlobIndex: 100,
			BatchMetadata: &grpcdisperser.BatchMetadata{
				BatchHeaderHash: []byte("mock-batch-header-hash"),
				BatchHeader: &grpcdisperser.BatchHeader{
					ReferenceBlockNumber: 200,
				},
			},
		},
	}
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_FINALIZED, Info: finalizedBlobInfo}, nil).Once())
	(disperserClient.On("RetrieveBlob", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil).Once()) // pass nil in as the return blob to tell the mock to return the corresponding blob
	logger := log.NewLogger(log.DiscardHandler())
	ifftCodec := codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec())
	eigendaClient := clients.EigenDAClient{
		Log: logger,
		Config: clients.EigenDAClientConfig{
			RPC:                          "localhost:51001",
			StatusQueryTimeout:           10 * time.Minute,
			StatusQueryRetryInterval:     50 * time.Millisecond,
			ResponseTimeout:              10 * time.Second,
			CustomQuorumIDs:              []uint{},
			SignerPrivateKeyHex:          "75f9e29cac7f5774d106adb355ef294987ce39b7863b75bb3f2ea42ca160926d",
			DisableTLS:                   false,
			PutBlobEncodingVersion:       codecs.DefaultBlobEncoding,
			DisablePointVerificationMode: false,
			WaitForFinalization:          true,
		},
		Client: disperserClient,
		Codec:  ifftCodec,
	}
	expectedBlob := []byte("dc49e7df326cfb2e7da5cf68f263e1898443ec2e862350606e7dfbda55ad10b5d61ed1d54baf6ae7a86279c1b4fa9c49a7de721dacb211264c1f5df31bade51c")
	blobInfo, err := eigendaClient.PutBlob(context.Background(), expectedBlob)
	require.NoError(t, err)
	require.NotNil(t, blobInfo)
	assert.Equal(t, finalizedBlobInfo, blobInfo)

	resultBlob, err := eigendaClient.GetBlob(context.Background(), []byte("mock-batch-header-hash"), 100)
	require.NoError(t, err)
	encodedBlob, err := ifftCodec.EncodeBlob(resultBlob)
	require.NoError(t, err)

	resultBlob, err = codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec()).DecodeBlob(encodedBlob)
	require.NoError(t, err)
	require.Equal(t, expectedBlob, resultBlob)
}

func TestPutRetrieveBlobNoIFFTSuccess(t *testing.T) {
	disperserClient := clientsmock.NewMockDisperserClient()
	expectedBlobStatus := disperser.Processing
	(disperserClient.On("DisperseBlobAuthenticated", mock.Anything, mock.Anything, mock.Anything).
		Return(&expectedBlobStatus, []byte("mock-request-id"), nil))
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_PROCESSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_DISPERSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_CONFIRMED}, nil).Once())
	finalizedBlobInfo := &grpcdisperser.BlobInfo{
		BlobHeader: &grpcdisperser.BlobHeader{
			Commitment: &common.G1Commitment{X: []byte{0x00, 0x00, 0x00, 0x00}, Y: []byte{0x01, 0x00, 0x00, 0x00}},
			BlobQuorumParams: []*grpcdisperser.BlobQuorumParam{
				{
					QuorumNumber: 0,
				},
				{
					QuorumNumber: 1,
				},
			},
		},
		BlobVerificationProof: &grpcdisperser.BlobVerificationProof{
			BlobIndex: 100,
			BatchMetadata: &grpcdisperser.BatchMetadata{
				BatchHeaderHash: []byte("mock-batch-header-hash"),
				BatchHeader: &grpcdisperser.BatchHeader{
					ReferenceBlockNumber: 200,
				},
			},
		},
	}
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_FINALIZED, Info: finalizedBlobInfo}, nil).Once())
	(disperserClient.On("RetrieveBlob", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil).Once()) // pass nil in as the return blob to tell the mock to return the corresponding blob
	logger := log.NewLogger(log.DiscardHandler())
	eigendaClient := clients.EigenDAClient{
		Log: logger,
		Config: clients.EigenDAClientConfig{
			RPC:                          "localhost:51001",
			StatusQueryTimeout:           10 * time.Minute,
			StatusQueryRetryInterval:     50 * time.Millisecond,
			ResponseTimeout:              10 * time.Second,
			CustomQuorumIDs:              []uint{},
			SignerPrivateKeyHex:          "75f9e29cac7f5774d106adb355ef294987ce39b7863b75bb3f2ea42ca160926d",
			DisableTLS:                   false,
			PutBlobEncodingVersion:       codecs.DefaultBlobEncoding,
			DisablePointVerificationMode: true,
			WaitForFinalization:          true,
		},
		Client: disperserClient,
		Codec:  codecs.NewNoIFFTCodec(codecs.NewDefaultBlobCodec()),
	}
	expectedBlob := []byte("dc49e7df326cfb2e7da5cf68f263e1898443ec2e862350606e7dfbda55ad10b5d61ed1d54baf6ae7a86279c1b4fa9c49a7de721dacb211264c1f5df31bade51c")
	blobInfo, err := eigendaClient.PutBlob(context.Background(), expectedBlob)
	require.NoError(t, err)
	require.NotNil(t, blobInfo)
	assert.Equal(t, finalizedBlobInfo, blobInfo)

	resultBlob, err := eigendaClient.GetBlob(context.Background(), []byte("mock-batch-header-hash"), 100)
	require.NoError(t, err)
	require.Equal(t, expectedBlob, resultBlob)
}

func TestPutBlobFailDispersal(t *testing.T) {
	disperserClient := clientsmock.NewMockDisperserClient()
	(disperserClient.On("DisperseBlobAuthenticated", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil, fmt.Errorf("error dispersing")))
	logger := log.NewLogger(log.DiscardHandler())
	eigendaClient := clients.EigenDAClient{
		Log: logger,
		Config: clients.EigenDAClientConfig{
			RPC:                      "localhost:51001",
			StatusQueryTimeout:       10 * time.Minute,
			StatusQueryRetryInterval: 50 * time.Millisecond,
			ResponseTimeout:          10 * time.Second,
			CustomQuorumIDs:          []uint{},
			SignerPrivateKeyHex:      "75f9e29cac7f5774d106adb355ef294987ce39b7863b75bb3f2ea42ca160926d",
			DisableTLS:               false,
			PutBlobEncodingVersion:   codecs.DefaultBlobEncoding,
			WaitForFinalization:      true,
		},
		Client: disperserClient,
		Codec:  codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec()),
	}
	blobInfo, err := eigendaClient.PutBlob(context.Background(), []byte("hello"))
	require.Error(t, err)
	require.Nil(t, blobInfo)
}

func TestPutBlobFailureInsufficentSignatures(t *testing.T) {
	disperserClient := clientsmock.NewMockDisperserClient()
	expectedBlobStatus := disperser.Processing
	(disperserClient.On("DisperseBlobAuthenticated", mock.Anything, mock.Anything, mock.Anything).
		Return(&expectedBlobStatus, []byte("mock-request-id"), nil))
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_PROCESSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_DISPERSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_INSUFFICIENT_SIGNATURES}, nil).Once())
	logger := log.NewLogger(log.DiscardHandler())
	eigendaClient := clients.EigenDAClient{
		Log: logger,
		Config: clients.EigenDAClientConfig{
			RPC:                      "localhost:51001",
			StatusQueryTimeout:       10 * time.Minute,
			StatusQueryRetryInterval: 50 * time.Millisecond,
			ResponseTimeout:          10 * time.Second,
			CustomQuorumIDs:          []uint{},
			SignerPrivateKeyHex:      "75f9e29cac7f5774d106adb355ef294987ce39b7863b75bb3f2ea42ca160926d",
			DisableTLS:               false,
			PutBlobEncodingVersion:   codecs.DefaultBlobEncoding,
			WaitForFinalization:      true,
		},
		Client: disperserClient,
		Codec:  codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec()),
	}
	blobInfo, err := eigendaClient.PutBlob(context.Background(), []byte("hello"))
	require.Error(t, err)
	require.Nil(t, blobInfo)
}

func TestPutBlobFailureGeneral(t *testing.T) {
	disperserClient := clientsmock.NewMockDisperserClient()
	expectedBlobStatus := disperser.Processing
	(disperserClient.On("DisperseBlobAuthenticated", mock.Anything, mock.Anything, mock.Anything).
		Return(&expectedBlobStatus, []byte("mock-request-id"), nil))
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_PROCESSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_DISPERSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_FAILED}, nil).Once())
	logger := log.NewLogger(log.DiscardHandler())
	eigendaClient := clients.EigenDAClient{
		Log: logger,
		Config: clients.EigenDAClientConfig{
			RPC:                      "localhost:51001",
			StatusQueryTimeout:       10 * time.Minute,
			StatusQueryRetryInterval: 50 * time.Millisecond,
			ResponseTimeout:          10 * time.Second,
			CustomQuorumIDs:          []uint{},
			SignerPrivateKeyHex:      "75f9e29cac7f5774d106adb355ef294987ce39b7863b75bb3f2ea42ca160926d",
			DisableTLS:               false,
			PutBlobEncodingVersion:   codecs.DefaultBlobEncoding,
			WaitForFinalization:      true,
		},
		Client: disperserClient,
		Codec:  codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec()),
	}
	blobInfo, err := eigendaClient.PutBlob(context.Background(), []byte("hello"))
	require.Error(t, err)
	require.Nil(t, blobInfo)
}

func TestPutBlobFailureUnknown(t *testing.T) {
	disperserClient := clientsmock.NewMockDisperserClient()
	expectedBlobStatus := disperser.Processing
	(disperserClient.On("DisperseBlobAuthenticated", mock.Anything, mock.Anything, mock.Anything).
		Return(&expectedBlobStatus, []byte("mock-request-id"), nil))
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_PROCESSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_DISPERSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_UNKNOWN}, nil).Once())
	logger := log.NewLogger(log.DiscardHandler())
	eigendaClient := clients.EigenDAClient{
		Log: logger,
		Config: clients.EigenDAClientConfig{
			RPC:                      "localhost:51001",
			StatusQueryTimeout:       10 * time.Minute,
			StatusQueryRetryInterval: 50 * time.Millisecond,
			ResponseTimeout:          10 * time.Second,
			CustomQuorumIDs:          []uint{},
			SignerPrivateKeyHex:      "75f9e29cac7f5774d106adb355ef294987ce39b7863b75bb3f2ea42ca160926d",
			DisableTLS:               false,
			PutBlobEncodingVersion:   codecs.DefaultBlobEncoding,
			WaitForFinalization:      true,
		},
		Client: disperserClient,
		Codec:  codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec()),
	}
	blobInfo, err := eigendaClient.PutBlob(context.Background(), []byte("hello"))
	require.Error(t, err)
	require.Nil(t, blobInfo)
}

func TestPutBlobFinalizationTimeout(t *testing.T) {
	disperserClient := clientsmock.NewMockDisperserClient()
	expectedBlobStatus := disperser.Processing
	(disperserClient.On("DisperseBlobAuthenticated", mock.Anything, mock.Anything, mock.Anything).
		Return(&expectedBlobStatus, []byte("mock-request-id"), nil))
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_PROCESSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_PROCESSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_PROCESSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_PROCESSING}, nil).Once())
	logger := log.NewLogger(log.DiscardHandler())
	eigendaClient := clients.EigenDAClient{
		Log: logger,
		Config: clients.EigenDAClientConfig{
			RPC:                      "localhost:51001",
			StatusQueryTimeout:       200 * time.Millisecond,
			StatusQueryRetryInterval: 51 * time.Millisecond,
			ResponseTimeout:          10 * time.Second,
			CustomQuorumIDs:          []uint{},
			SignerPrivateKeyHex:      "75f9e29cac7f5774d106adb355ef294987ce39b7863b75bb3f2ea42ca160926d",
			DisableTLS:               false,
			PutBlobEncodingVersion:   codecs.DefaultBlobEncoding,
			WaitForFinalization:      true,
		},
		Client: disperserClient,
		Codec:  codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec()),
	}
	blobInfo, err := eigendaClient.PutBlob(context.Background(), []byte("hello"))
	require.Error(t, err)
	require.Nil(t, blobInfo)
}

func TestPutBlobIndividualRequestTimeout(t *testing.T) {
	disperserClient := clientsmock.NewMockDisperserClient()
	expectedBlobStatus := disperser.Processing
	(disperserClient.On("DisperseBlobAuthenticated", mock.Anything, mock.Anything, mock.Anything).
		Return(&expectedBlobStatus, []byte("mock-request-id"), nil))
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			time.Sleep(100 * time.Millisecond) // Simulate a 100ms delay, which should fail the request
		}).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_PROCESSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_DISPERSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_CONFIRMED}, nil).Once())
	finalizedBlobInfo := &grpcdisperser.BlobInfo{
		BlobHeader: &grpcdisperser.BlobHeader{
			Commitment: &common.G1Commitment{X: []byte{0x00, 0x00, 0x00, 0x00}, Y: []byte{0x01, 0x00, 0x00, 0x00}},
			BlobQuorumParams: []*grpcdisperser.BlobQuorumParam{
				{
					QuorumNumber: 0,
				},
				{
					QuorumNumber: 1,
				},
			},
		},
		BlobVerificationProof: &grpcdisperser.BlobVerificationProof{
			BlobIndex: 100,
			BatchMetadata: &grpcdisperser.BatchMetadata{
				BatchHeaderHash: []byte("mock-batch-header-hash"),
				BatchHeader: &grpcdisperser.BatchHeader{
					ReferenceBlockNumber: 200,
				},
			},
		},
	}
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_FINALIZED, Info: finalizedBlobInfo}, nil).Once())
	logger := log.NewLogger(log.DiscardHandler())
	eigendaClient := clients.EigenDAClient{
		Log: logger,
		Config: clients.EigenDAClientConfig{
			RPC:                      "localhost:51001",
			StatusQueryTimeout:       10 * time.Minute,
			StatusQueryRetryInterval: 50 * time.Millisecond,
			ResponseTimeout:          50 * time.Millisecond, // very low timeout
			CustomQuorumIDs:          []uint{},
			SignerPrivateKeyHex:      "75f9e29cac7f5774d106adb355ef294987ce39b7863b75bb3f2ea42ca160926d",
			DisableTLS:               false,
			PutBlobEncodingVersion:   codecs.DefaultBlobEncoding,
			WaitForFinalization:      true,
		},
		Client: disperserClient,
		Codec:  codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec()),
	}
	blobInfo, err := eigendaClient.PutBlob(context.Background(), []byte("hello"))

	// despite initial timeout it should succeed
	require.NoError(t, err)
	require.NotNil(t, blobInfo)
	assert.Equal(t, finalizedBlobInfo, blobInfo)
}

func TestPutBlobTotalTimeout(t *testing.T) {
	disperserClient := clientsmock.NewMockDisperserClient()
	expectedBlobStatus := disperser.Processing
	(disperserClient.On("DisperseBlobAuthenticated", mock.Anything, mock.Anything, mock.Anything).
		Return(&expectedBlobStatus, []byte("mock-request-id"), nil))
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			time.Sleep(100 * time.Millisecond) // Simulate a 100ms delay, which should fail the request
		}).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_PROCESSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_DISPERSING}, nil).Once())
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_CONFIRMED}, nil).Once())
	finalizedBlobInfo := &grpcdisperser.BlobInfo{
		BlobHeader: &grpcdisperser.BlobHeader{
			Commitment: &common.G1Commitment{X: []byte{0x00, 0x00, 0x00, 0x00}, Y: []byte{0x01, 0x00, 0x00, 0x00}},
			BlobQuorumParams: []*grpcdisperser.BlobQuorumParam{
				{
					QuorumNumber: 0,
				},
				{
					QuorumNumber: 1,
				},
			},
		},
		BlobVerificationProof: &grpcdisperser.BlobVerificationProof{
			BlobIndex: 100,
			BatchMetadata: &grpcdisperser.BatchMetadata{
				BatchHeaderHash: []byte("mock-batch-header-hash"),
				BatchHeader: &grpcdisperser.BatchHeader{
					ReferenceBlockNumber: 200,
				},
			},
		},
	}
	(disperserClient.On("GetBlobStatus", mock.Anything, mock.Anything).
		Return(&grpcdisperser.BlobStatusReply{Status: grpcdisperser.BlobStatus_FINALIZED, Info: finalizedBlobInfo}, nil).Once())
	logger := log.NewLogger(log.DiscardHandler())
	eigendaClient := clients.EigenDAClient{
		Log: logger,
		Config: clients.EigenDAClientConfig{
			RPC:                      "localhost:51001",
			StatusQueryTimeout:       100 * time.Millisecond, // low total timeout
			StatusQueryRetryInterval: 50 * time.Millisecond,
			ResponseTimeout:          10 * time.Minute,
			CustomQuorumIDs:          []uint{},
			SignerPrivateKeyHex:      "75f9e29cac7f5774d106adb355ef294987ce39b7863b75bb3f2ea42ca160926d",
			DisableTLS:               false,
			PutBlobEncodingVersion:   codecs.DefaultBlobEncoding,
			WaitForFinalization:      true,
		},
		Client: disperserClient,
		Codec:  codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec()),
	}
	blobInfo, err := eigendaClient.PutBlob(context.Background(), []byte("hello"))

	// should timeout even though it would have finalized eventually
	require.Error(t, err)
	require.Nil(t, blobInfo)
}
