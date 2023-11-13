package encoder

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/big"
	"runtime"
	"testing"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/encoding"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	pb "github.com/Layr-Labs/eigenda/disperser/proto/protogen/encoder"
	"github.com/Layr-Labs/eigenda/pkg/encoding/kzgEncoder"
	"github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

var (
	gettysburgAddressBytes = []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")
)

var logger = &cmock.Logger{}

func makeTestEncoder() (*encoding.Encoder, ServerConfig) {
	kzgConfig := kzgEncoder.KzgConfig{
		G1Path:    "../../inabox/resources/kzg/g1.point",
		G2Path:    "../../inabox/resources/kzg/g2.point",
		CacheDir:  "../../inabox/resources/kzg/SRSTables",
		SRSOrder:  3000,
		NumWorker: uint64(runtime.GOMAXPROCS(0)),
	}

	encodingConfig := encoding.EncoderConfig{KzgConfig: kzgConfig}

	encoder, _ := encoding.NewEncoder(encodingConfig)
	encoderServerConfig := ServerConfig{
		GrpcPort:              "3000",
		MaxConcurrentRequests: 16,
		RequestPoolSize:       32,
	}

	return encoder, encoderServerConfig
}

var testEncoder, testServerConfig = makeTestEncoder()

func getTestData() (core.Blob, core.EncodingParams) {
	var quorumID core.QuorumID = 0
	var quantizationFactor uint = 2
	var adversaryThreshold uint8 = 80
	var quorumThreshold uint8 = 90
	securityParams := []*core.SecurityParam{
		{
			QuorumID:           quorumID,
			AdversaryThreshold: adversaryThreshold,
		},
	}

	testBlob := core.Blob{
		RequestHeader: core.BlobRequestHeader{
			SecurityParams: securityParams,
		},
		Data: gettysburgAddressBytes,
	}

	indexedChainState, _ := coremock.NewChainDataMock(core.OperatorIndex(10))
	operatorState, err := indexedChainState.GetOperatorState(context.Background(), uint(0), []core.QuorumID{quorumID})
	if err != nil {
		log.Fatalf("failed to get operator state: %s", err)
	}
	coordinator := &core.StdAssignmentCoordinator{}
	_, info, err := coordinator.GetAssignments(operatorState, quorumID, quantizationFactor)
	if err != nil {
		log.Fatal(err)
	}

	blobSize := uint(len(testBlob.Data))

	blobLength := core.GetBlobLength(uint(blobSize))
	numOperators := uint(len(operatorState.Operators[quorumID]))
	chunkLength, err := coordinator.GetMinimumChunkLength(numOperators, blobLength, quantizationFactor, quorumThreshold, adversaryThreshold)
	if err != nil {
		log.Fatal(err)
	}

	testEncodingParams, _ := core.GetEncodingParams(chunkLength, info.TotalChunks)

	return testBlob, testEncodingParams
}

func newEncoderTestServer(t *testing.T) *Server {
	metrics := NewMetrics("9000", logger)
	return NewServer(testServerConfig, logger, testEncoder, metrics)
}

func TestEncodeBlob(t *testing.T) {
	server := newEncoderTestServer(t)
	testBlobData, testEncodingParams := getTestData()

	testEncodingParamsProto := &pb.EncodingParams{
		ChunkLength: uint32(testEncodingParams.ChunkLength),
		NumChunks:   uint32(testEncodingParams.NumChunks),
	}

	encodeBlobRequestProto := &pb.EncodeBlobRequest{
		Data:           []byte(testBlobData.Data),
		EncodingParams: testEncodingParamsProto,
	}

	reply, err := server.EncodeBlob(context.Background(), encodeBlobRequestProto)
	assert.NoError(t, err)
	assert.NotNil(t, reply.Chunks)

	// Decode Server Data
	var chunksData []*core.Chunk

	for i := range reply.Chunks {
		chunkSerialized, _ := new(core.Chunk).Deserialize(reply.GetChunks()[i])
		// perform an operation
		chunksData = append(chunksData, chunkSerialized)
	}
	assert.NotNil(t, chunksData)

	// Indices obtained from Encoder_Test
	indices := []core.ChunkNumber{
		0, 1, 2, 3, 4, 5, 6, 7,
	}

	maxInputSize := uint64(len(gettysburgAddressBytes)) + 10
	decoded, err := testEncoder.Decode(chunksData, indices, testEncodingParams, maxInputSize)
	assert.Nil(t, err)
	recovered := bytes.TrimRight(decoded, "\x00")
	assert.Equal(t, recovered, gettysburgAddressBytes)
}

func TestThrottling(t *testing.T) {
	var X1, Y1 fp.Element
	X1 = *X1.SetBigInt(big.NewInt(1))
	Y1 = *Y1.SetBigInt(big.NewInt(2))

	metrics := NewMetrics("9000", logger)
	concurrentRequests := 2
	requestPoolSize := 4
	encoder := &encoding.MockEncoder{
		Delay: 500 * time.Millisecond,
	}
	blobCommitment := core.BlobCommitments{
		Commitment: &core.Commitment{
			G1Point: &bn254.G1Point{
				X: X1,
				Y: Y1,
			},
		},
		LengthProof: &core.Commitment{
			G1Point: &bn254.G1Point{
				X: X1,
				Y: Y1,
			},
		},
		Length: 10,
	}
	encoder.On("Encode", mock.Anything, mock.Anything).Return(blobCommitment, []*core.Chunk{}, nil)
	encoderServerConfig := ServerConfig{
		GrpcPort:              "3000",
		MaxConcurrentRequests: concurrentRequests,
		RequestPoolSize:       requestPoolSize,
	}
	s := NewServer(encoderServerConfig, logger, encoder, metrics)
	testBlobData, testEncodingParams := getTestData()

	testEncodingParamsProto := &pb.EncodingParams{
		ChunkLength: uint32(testEncodingParams.ChunkLength),
		NumChunks:   uint32(testEncodingParams.NumChunks),
	}

	encodeBlobRequestProto := &pb.EncodeBlobRequest{
		Data:           []byte(testBlobData.Data),
		EncodingParams: testEncodingParamsProto,
	}

	errs := make([]error, requestPoolSize+1)
	done := make(chan struct{}, requestPoolSize+1)
	for i := 0; i < requestPoolSize+1; i++ {
		go func(i int) {
			timeout := 200 * time.Millisecond
			fmt.Println("Making request", i, timeout)
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			_, err := s.EncodeBlob(ctx, encodeBlobRequestProto)
			errs[i] = err
			done <- struct{}{}
		}(i)

		time.Sleep(10 * time.Millisecond)
	}

	for i := 0; i < requestPoolSize+1; i++ {
		<-done
	}

	for i := 0; i < requestPoolSize+1; i++ {
		fmt.Println(errs[i])
	}

	for i := 0; i < requestPoolSize+1; i++ {
		err := errs[i]
		if i < concurrentRequests {
			assert.NoError(t, err)
		} else if i >= requestPoolSize {
			assert.ErrorContains(t, err, "too many requests")
		} else {
			assert.ErrorIs(t, err, context.DeadlineExceeded)
		}
	}
}
