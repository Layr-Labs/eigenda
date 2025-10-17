package encoder_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/big"
	"runtime"
	"testing"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/encoder"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	encmock "github.com/Layr-Labs/eigenda/disperser/mock"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/codec"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	gettysburgAddressBytes = []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")
)

func makeTestProverV1(numPoint uint64) (*prover.Prover, encoder.ServerConfig) {
	kzgConfig := &kzg.KzgConfig{
		G1Path:          "../../resources/srs/g1.point",
		G2Path:          "../../resources/srs/g2.point",
		CacheDir:        "../../resources/srs/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: numPoint,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    true,
	}

	p, _ := prover.NewProver(kzgConfig, nil)
	encoderServerConfig := encoder.ServerConfig{
		GrpcPort:              "3000",
		MaxConcurrentRequests: 16,
		RequestPoolSize:       32,
	}

	return p, encoderServerConfig
}

var testProver, testServerConfig = makeTestProverV1(3000)

func getTestData(t *testing.T) (core.Blob, encoding.EncodingParams) {
	t.Helper()
	ctx := t.Context()

	var quorumID core.QuorumID = 0
	var adversaryThreshold uint8 = 80
	var quorumThreshold uint8 = 90
	securityParams := []*core.SecurityParam{
		{
			QuorumID:              quorumID,
			ConfirmationThreshold: quorumThreshold,
			AdversaryThreshold:    adversaryThreshold,
		},
	}

	testBlob := core.Blob{
		RequestHeader: core.BlobRequestHeader{
			SecurityParams: securityParams,
		},
		Data: codec.ConvertByPaddingEmptyByte(gettysburgAddressBytes),
	}

	indexedChainState, _ := coremock.MakeChainDataMock(map[uint8]int{
		0: 10,
		1: 10,
		2: 10,
	})
	operatorState, err := indexedChainState.GetOperatorState(ctx, uint(0), []core.QuorumID{quorumID})
	if err != nil {
		log.Fatalf("failed to get operator state: %s", err)
	}
	coordinator := &core.StdAssignmentCoordinator{}

	blobSize := uint32(len(testBlob.Data))
	blobLength := encoding.GetBlobLength(blobSize)

	chunkLength, err := coordinator.CalculateChunkLength(operatorState, uint(blobLength), 0, securityParams[0])
	if err != nil {
		log.Fatal(err)
	}

	blobQuorumInfo := &core.BlobQuorumInfo{
		SecurityParam: *securityParams[0],
		ChunkLength:   chunkLength,
	}

	_, info, err := coordinator.GetAssignments(operatorState, uint(blobLength), blobQuorumInfo)
	if err != nil {
		log.Fatal(err)
	}

	testEncodingParams := encoding.ParamsFromMins(uint64(chunkLength), info.TotalChunks)

	return testBlob, testEncodingParams
}

func newEncoderTestServer(t *testing.T) *encoder.EncoderServer {
	metrics := encoder.NewMetrics(prometheus.NewRegistry(), "9000", logger)
	return encoder.NewEncoderServer(testServerConfig, logger, testProver, metrics, nil)
}

func TestEncodeBlobV1(t *testing.T) {
	server := newEncoderTestServer(t)
	testBlobData, testEncodingParams := getTestData(t)

	testEncodingParamsProto := &pb.EncodingParams{
		ChunkLength: uint32(testEncodingParams.ChunkLength),
		NumChunks:   uint32(testEncodingParams.NumChunks),
	}

	encodeBlobRequestProto := &pb.EncodeBlobRequest{
		Data:           []byte(testBlobData.Data),
		EncodingParams: testEncodingParamsProto,
	}

	reply, err := server.EncodeBlob(t.Context(), encodeBlobRequestProto)
	assert.NoError(t, err)
	assert.NotNil(t, reply.GetChunks())

	// Decode Server Data
	var chunksData []*encoding.Frame

	for i := range reply.GetChunks() {
		chunkSerialized, _ := new(encoding.Frame).DeserializeGob(reply.GetChunks()[i])
		// perform an operation
		chunksData = append(chunksData, chunkSerialized)
	}
	assert.NotNil(t, chunksData)

	// Indices obtained from Encoder_Test
	indices := make([]encoding.ChunkNumber, len(reply.GetChunks()))
	for i := range indices {
		indices[i] = encoding.ChunkNumber(i)
	}

	maxInputSize := uint64(len(testBlobData.Data)) + 10
	decoded, err := testProver.Decode(chunksData, indices, testEncodingParams, maxInputSize)
	assert.Nil(t, err)

	recovered := codec.RemoveEmptyByteFromPaddedBytes(decoded)

	restored := bytes.TrimRight(recovered, "\x00")
	assert.Equal(t, restored, gettysburgAddressBytes)
}

func TestThrottling(t *testing.T) {
	ctx := t.Context()
	var X1, Y1 fp.Element
	X1 = *X1.SetBigInt(big.NewInt(1))
	Y1 = *Y1.SetBigInt(big.NewInt(2))

	var lengthXA0, lengthXA1, lengthYA0, lengthYA1 fp.Element
	_, err := lengthXA0.SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781")
	assert.NoError(t, err)
	_, err = lengthXA1.SetString("11559732032986387107991004021392285783925812861821192530917403151452391805634")
	assert.NoError(t, err)
	_, err = lengthYA0.SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930")
	assert.NoError(t, err)
	_, err = lengthYA1.SetString("4082367875863433681332203403145435568316851327593401208105741076214120093531")
	assert.NoError(t, err)

	var lengthProof, lengthCommitment bn254.G2Affine
	lengthProof.X.A0 = lengthXA0
	lengthProof.X.A1 = lengthXA1
	lengthProof.Y.A0 = lengthYA0
	lengthProof.Y.A1 = lengthYA1

	lengthCommitment = lengthProof

	metrics := encoder.NewMetrics(prometheus.NewRegistry(), "9000", logger)
	concurrentRequests := 2
	requestPoolSize := 4
	mockEncoder := &encmock.MockEncoder{
		Delay: 500 * time.Millisecond,
	}

	blobCommitment := encoding.BlobCommitments{
		Commitment: &encoding.G1Commitment{
			X: X1,
			Y: Y1,
		},
		LengthCommitment: (*encoding.G2Commitment)(&lengthCommitment),
		LengthProof:      (*encoding.G2Commitment)(&lengthProof),
		Length:           10,
	}

	mockEncoder.On("EncodeAndProve", mock.Anything, mock.Anything).Return(blobCommitment, []*encoding.Frame{}, nil)
	encoderServerConfig := encoder.ServerConfig{
		GrpcPort:              "3000",
		MaxConcurrentRequests: concurrentRequests,
		RequestPoolSize:       requestPoolSize,
	}
	s := encoder.NewEncoderServer(encoderServerConfig, logger, mockEncoder, metrics, nil)
	testBlobData, testEncodingParams := getTestData(t)

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
			ctx, cancel := context.WithTimeout(ctx, timeout)
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

func TestEncoderPointsLoading(t *testing.T) {
	ctx := t.Context()
	// encoder 1 only loads 1500 points
	prover1, config1 := makeTestProverV1(1500)
	metrics := encoder.NewMetrics(prometheus.NewRegistry(), "9000", logger)
	server1 := encoder.NewEncoderServer(config1, logger, prover1, metrics, nil)

	testBlobData, testEncodingParams := getTestData(t)

	testEncodingParamsProto := &pb.EncodingParams{
		ChunkLength: uint32(testEncodingParams.ChunkLength),
		NumChunks:   uint32(testEncodingParams.NumChunks),
	}

	encodeBlobRequestProto := &pb.EncodeBlobRequest{
		Data:           []byte(testBlobData.Data),
		EncodingParams: testEncodingParamsProto,
	}

	reply1, err := server1.EncodeBlob(ctx, encodeBlobRequestProto)
	assert.NoError(t, err)
	assert.NotNil(t, reply1.GetChunks())

	// Decode Server Data
	var chunksData []*encoding.Frame

	for i := range reply1.GetChunks() {
		chunkSerialized, _ := new(encoding.Frame).DeserializeGob(reply1.GetChunks()[i])
		// perform an operation
		chunksData = append(chunksData, chunkSerialized)
	}
	assert.NotNil(t, chunksData)

	indices := make([]encoding.ChunkNumber, len(reply1.GetChunks()))
	for i := range indices {
		indices[i] = encoding.ChunkNumber(i)
	}

	maxInputSize := uint64(len(testBlobData.Data)) + 10
	decoded, err := testProver.Decode(chunksData, indices, testEncodingParams, maxInputSize)
	assert.Nil(t, err)

	recovered := codec.RemoveEmptyByteFromPaddedBytes(decoded)

	restored := bytes.TrimRight(recovered, "\x00")
	assert.Equal(t, restored, gettysburgAddressBytes)

	// encoder 2 only loads 2900 points
	encoder2, config2 := makeTestProverV1(2900)
	server2 := encoder.NewEncoderServer(config2, logger, encoder2, metrics, nil)

	reply2, err := server2.EncodeBlob(ctx, encodeBlobRequestProto)
	assert.NoError(t, err)
	assert.NotNil(t, reply2.GetChunks())

	for i := range reply2.GetChunks() {
		chunkSerialized, _ := new(encoding.Frame).DeserializeGob(reply2.GetChunks()[i])
		// perform an operation
		assert.Equal(t, len(chunkSerialized.Coeffs), len(chunksData[i].Coeffs))
		assert.Equal(t, chunkSerialized.Coeffs, chunksData[i].Coeffs)
		assert.Equal(t, chunkSerialized.Proof, chunksData[i].Proof)
	}

}
