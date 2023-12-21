package integration

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/encoding"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/gammazero/workerpool"
	"github.com/stretchr/testify/assert"

	"github.com/Layr-Labs/eigenda/pkg/encoding/kzgEncoder"
)

var (
	enc core.Encoder
	asn core.AssignmentCoordinator = &core.StdAssignmentCoordinator{}
)

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	os.Exit(code)
}

func setup(m *testing.M) {

	var err error
	enc, err = makeTestEncoder()
	if err != nil {
		panic("failed to start localstack container")
	}
}

// makeTestEncoder makes an encoder currently using the only supported backend.
func makeTestEncoder() (core.Encoder, error) {
	config := kzgEncoder.KzgConfig{
		G1Path:    "../../inabox/resources/kzg/g1.point",
		G2Path:    "../../inabox/resources/kzg/g2.point",
		CacheDir:  "../../inabox/resources/kzg/SRSTables",
		SRSOrder:  3000,
		NumWorker: uint64(runtime.GOMAXPROCS(0)),
	}

	return encoding.NewEncoder(encoding.EncoderConfig{KzgConfig: config})

}

func makeTestBlob(t *testing.T, length int, securityParams []*core.SecurityParam) core.Blob {

	data := make([]byte, length)
	_, err := rand.Read(data)
	if err != nil {
		t.Fatal(err)
	}

	blob := core.Blob{
		RequestHeader: core.BlobRequestHeader{
			SecurityParams: securityParams,
		},
		Data: data,
	}
	return blob
}

// prepareBatch takes in multiple blob, encodes them, generates the associated assignments, and the batch header.
// These are the products that a disperser will need in order to disperse data to the DA nodes.
func prepareBatch(t *testing.T, cst core.IndexedChainState, blobs []core.Blob, quorumIndex uint, quantizationFactor uint, bn uint) ([]core.EncodedBlob, core.BatchHeader) {

	batchHeader := core.BatchHeader{
		ReferenceBlockNumber: bn,
		BatchRoot:            [32]byte{},
	}

	numBlob := len(blobs)
	var encodedBlobs []core.EncodedBlob = make([]core.EncodedBlob, numBlob)

	for z, blob := range blobs {
		quorumID := blob.RequestHeader.SecurityParams[quorumIndex].QuorumID
		quorums := []core.QuorumID{quorumID}

		state, err := cst.GetOperatorState(context.Background(), bn, quorums)
		if err != nil {
			t.Fatal(err)
		}

		assignments, info, err := asn.GetAssignments(state, quorumID, quantizationFactor)
		if err != nil {
			t.Fatal(err)
		}

		blobSize := uint(len(blob.Data))
		blobLength := core.GetBlobLength(blobSize)
		adversaryThreshold := blob.RequestHeader.SecurityParams[quorumIndex].AdversaryThreshold
		quorumThreshold := blob.RequestHeader.SecurityParams[quorumIndex].QuorumThreshold

		numOperators := uint(len(state.Operators[quorumID]))

		chunkLength, err := asn.GetMinimumChunkLength(numOperators, blobLength, quantizationFactor, quorumThreshold, adversaryThreshold)
		if err != nil {
			t.Fatal(err)
		}

		params, err := core.GetEncodingParams(chunkLength, info.TotalChunks)
		if err != nil {
			t.Fatal(err)
		}

		commitments, chunks, err := enc.Encode(blob.Data, params)
		if err != nil {
			t.Fatal(err)
		}

		quorumHeader := &core.BlobQuorumInfo{
			SecurityParam: core.SecurityParam{
				QuorumID:           quorumID,
				AdversaryThreshold: adversaryThreshold,
				QuorumThreshold:    quorumThreshold,
			},
			QuantizationFactor: quantizationFactor,
			EncodedBlobLength:  params.ChunkLength * quantizationFactor * numOperators,
		}

		blobHeader := &core.BlobHeader{
			BlobCommitments: core.BlobCommitments{
				Commitment:  commitments.Commitment,
				LengthProof: commitments.LengthProof,
				Length:      commitments.Length,
			},
			QuorumInfos: []*core.BlobQuorumInfo{quorumHeader},
		}

		var encodedBlob core.EncodedBlob = make(map[core.OperatorID]*core.BlobMessage, len(assignments))
		for id, assignment := range assignments {
			bundles := map[core.QuorumID]core.Bundle{
				quorumID: chunks[assignment.StartIndex : assignment.StartIndex+assignment.NumChunks],
			}
			encodedBlob[id] = &core.BlobMessage{
				BlobHeader: blobHeader,
				Bundles:    bundles,
			}
		}
		encodedBlobs[z] = encodedBlob

	}

	return encodedBlobs, batchHeader

}

// checkBatch runs the verification logic for each DA node in the current OperatorState, and returns an error if any of
// the DA nodes' validation checks fails
func checkBatch(t *testing.T, cst core.IndexedChainState, encodedBlob core.EncodedBlob, header core.BatchHeader) {
	val := core.NewChunkValidator(enc, asn, cst, [32]byte{})

	quorums := []core.QuorumID{0}
	state, _ := cst.GetIndexedOperatorState(context.Background(), header.ReferenceBlockNumber, quorums)

	for id := range state.IndexedOperators {
		val.UpdateOperatorID(id)
		blobMessage := encodedBlob[id]
		err := val.ValidateBlob(blobMessage, state.OperatorState)
		assert.NoError(t, err)
	}

}

// checkBatchByUniversalVerifier runs the verification logic for each DA node in the current OperatorState, and returns an error if any of
// the DA nodes' validation checks fails
func checkBatchByUniversalVerifier(t *testing.T, cst core.IndexedChainState, encodedBlobs []core.EncodedBlob, header core.BatchHeader, pool common.WorkerPool) {
	val := core.NewChunkValidator(enc, asn, cst, [32]byte{})

	quorums := []core.QuorumID{0}
	state, _ := cst.GetIndexedOperatorState(context.Background(), header.ReferenceBlockNumber, quorums)
	numBlob := len(encodedBlobs)

	for id := range state.IndexedOperators {
		val.UpdateOperatorID(id)
		var blobMessages []*core.BlobMessage = make([]*core.BlobMessage, numBlob)
		for z, encodedBlob := range encodedBlobs {
			blobMessages[z] = encodedBlob[id]
		}

		err := val.ValidateBatch(blobMessages, state.OperatorState, pool)
		assert.NoError(t, err)
	}

}

func TestCoreLibrary(t *testing.T) {

	numBlob := 1 // must be greater than 0
	blobLengths := []int{1, 64, 1000}
	quantizationFactors := []uint{1, 10}
	operatorCounts := []uint{1, 2, 4, 10, 30}

	securityParams := []*core.SecurityParam{
		{
			QuorumID:           0,
			AdversaryThreshold: 50,
			QuorumThreshold:    100,
		},
		{
			QuorumID:           0,
			AdversaryThreshold: 80,
			QuorumThreshold:    90,
		},
	}

	quorumIndex := uint(0)
	bn := uint(0)

	pool := workerpool.New(1)

	for _, operatorCount := range operatorCounts {
		cst, err := mock.NewChainDataMock(core.OperatorIndex(operatorCount))
		assert.NoError(t, err)
		batches := make([]core.EncodedBlob, 0)
		batchHeader := core.BatchHeader{
			ReferenceBlockNumber: bn,
			BatchRoot:            [32]byte{},
		}
		// batch can only be tested per operatorCount, because the assignment would be wrong otherwise
		for _, blobLength := range blobLengths {

			for _, quantizationFactor := range quantizationFactors {
				for _, securityParam := range securityParams {

					t.Run(fmt.Sprintf("blobLength=%v, quantizationFactor=%v, operatorCount=%v, securityParams=%v", blobLength, quantizationFactor, operatorCount, securityParam), func(t *testing.T) {

						blobs := make([]core.Blob, numBlob)
						for i := 0; i < numBlob; i++ {
							blobs[i] = makeTestBlob(t, blobLength, []*core.SecurityParam{securityParam})
						}

						batch, header := prepareBatch(t, cst, blobs, quorumIndex, quantizationFactor, bn)
						batches = append(batches, batch...)

						checkBatch(t, cst, batch[0], header)
					})
				}

			}

		}
		t.Run(fmt.Sprintf("universal verifier operatorCount=%v over %v blobs", operatorCount, len(batches)), func(t *testing.T) {
			checkBatchByUniversalVerifier(t, cst, batches, batchHeader, pool)
		})

	}

}

func TestParseOperatorSocket(t *testing.T) {
	operatorSocket := "localhost:1234;5678"
	host, dispersalPort, retrievalPort, err := core.ParseOperatorSocket(operatorSocket)
	assert.NoError(t, err)
	assert.Equal(t, "localhost", host)
	assert.Equal(t, "1234", dispersalPort)
	assert.Equal(t, "5678", retrievalPort)

	_, _, _, err = core.ParseOperatorSocket("localhost:12345678")
	assert.NotNil(t, err)
	assert.Equal(t, "invalid socket address format, missing retrieval port: localhost:12345678", err.Error())

	_, _, _, err = core.ParseOperatorSocket("localhost1234;5678")
	assert.NotNil(t, err)
	assert.Equal(t, "invalid socket address format: localhost1234;5678", err.Error())
}
