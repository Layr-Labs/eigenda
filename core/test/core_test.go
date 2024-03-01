package core_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/gammazero/workerpool"
	"github.com/stretchr/testify/assert"
)

var (
	p   encoding.Prover
	v   encoding.Verifier
	asn core.AssignmentCoordinator = &core.StdAssignmentCoordinator{}
)

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	os.Exit(code)
}

func setup(m *testing.M) {

	var err error
	p, v, err = makeTestComponents()
	if err != nil {
		panic("failed to start localstack container")
	}
}

// makeTestComponents makes a prover and verifier currently using the only supported backend.
func makeTestComponents() (encoding.Prover, encoding.Verifier, error) {
	config := &kzg.KzgConfig{
		G1Path:          "../../inabox/resources/kzg/g1.point",
		G2Path:          "../../inabox/resources/kzg/g2.point",
		CacheDir:        "../../inabox/resources/kzg/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	p, err := prover.NewProver(config, true)
	if err != nil {
		return nil, nil, err
	}

	v, err := verifier.NewVerifier(config, true)
	if err != nil {
		return nil, nil, err
	}

	return p, v, nil
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
func prepareBatch(t *testing.T, cst core.IndexedChainState, blobs []core.Blob, quorumIndex uint, bn uint) ([]core.EncodedBlob, core.BatchHeader) {

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

		blobSize := uint(len(blob.Data))
		blobLength := encoding.GetBlobLength(blobSize)

		chunkLength, err := asn.CalculateChunkLength(state, blobLength, 0, blob.RequestHeader.SecurityParams[quorumIndex])
		if err != nil {
			t.Fatal(err)
		}

		quorumHeader := &core.BlobQuorumInfo{
			SecurityParam: core.SecurityParam{
				QuorumID:              quorumID,
				AdversaryThreshold:    blob.RequestHeader.SecurityParams[quorumIndex].AdversaryThreshold,
				ConfirmationThreshold: blob.RequestHeader.SecurityParams[quorumIndex].ConfirmationThreshold,
			},
			ChunkLength: chunkLength,
		}

		assignments, info, err := asn.GetAssignments(state, blobLength, quorumHeader)
		if err != nil {
			t.Fatal(err)
		}

		params := encoding.ParamsFromMins(chunkLength, info.TotalChunks)

		commitments, chunks, err := p.EncodeAndProve(blob.Data, params)
		if err != nil {
			t.Fatal(err)
		}

		blobHeader := &core.BlobHeader{
			BlobCommitments: encoding.BlobCommitments{
				Commitment:       commitments.Commitment,
				LengthCommitment: commitments.LengthCommitment,
				LengthProof:      commitments.LengthProof,
				Length:           commitments.Length,
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
	val := core.NewChunkValidator(v, asn, cst, [32]byte{})

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
	val := core.NewChunkValidator(v, asn, cst, [32]byte{})

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
	operatorCounts := []uint{1, 2, 4, 10, 30}

	securityParams := []*core.SecurityParam{
		{
			QuorumID:              0,
			AdversaryThreshold:    50,
			ConfirmationThreshold: 100,
		},
		{
			QuorumID:              0,
			AdversaryThreshold:    80,
			ConfirmationThreshold: 90,
		},
	}

	quorumIndex := uint(0)
	bn := uint(0)

	pool := workerpool.New(1)

	for _, operatorCount := range operatorCounts {
		cst, err := mock.MakeChainDataMock(core.OperatorIndex(operatorCount))
		assert.NoError(t, err)
		batches := make([]core.EncodedBlob, 0)
		batchHeader := core.BatchHeader{
			ReferenceBlockNumber: bn,
			BatchRoot:            [32]byte{},
		}
		// batch can only be tested per operatorCount, because the assignment would be wrong otherwise
		for _, blobLength := range blobLengths {

			for _, securityParam := range securityParams {

				t.Run(fmt.Sprintf("blobLength=%v, operatorCount=%v, securityParams=%v", blobLength, operatorCount, securityParam), func(t *testing.T) {

					blobs := make([]core.Blob, numBlob)
					for i := 0; i < numBlob; i++ {
						blobs[i] = makeTestBlob(t, blobLength, []*core.SecurityParam{securityParam})
					}

					batch, header := prepareBatch(t, cst, blobs, quorumIndex, bn)
					batches = append(batches, batch...)

					checkBatch(t, cst, batch[0], header)
				})
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
