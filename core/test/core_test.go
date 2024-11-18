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
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/gammazero/workerpool"
	"github.com/hashicorp/go-multierror"
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
		LoadG2Points:    true,
	}

	p, err := prover.NewProver(config, nil)
	if err != nil {
		return nil, nil, err
	}

	v, err := verifier.NewVerifier(config, nil)
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

	data = codec.ConvertByPaddingEmptyByte(data)

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
func prepareBatch(t *testing.T, operatorCount uint, blobs []core.Blob, bn uint) ([]core.EncodedBlob, core.BatchHeader, *mock.ChainDataMock) {

	cst, err := mock.MakeChainDataMock(map[uint8]int{
		0: int(operatorCount),
		1: int(operatorCount),
		2: int(operatorCount),
	})
	assert.NoError(t, err)

	batchHeader := core.BatchHeader{
		ReferenceBlockNumber: bn,
		BatchRoot:            [32]byte{},
	}

	numBlob := len(blobs)
	encodedBlobs := make([]core.EncodedBlob, numBlob)
	blobHeaders := make([]*core.BlobHeader, numBlob)

	for z, blob := range blobs {

		blobHeader := &core.BlobHeader{
			QuorumInfos: make([]*core.BlobQuorumInfo, 0),
		}
		blobHeaders[z] = blobHeader

		encodedBlob := core.EncodedBlob{
			BlobHeader:               blobHeader,
			EncodedBundlesByOperator: make(map[core.OperatorID]core.EncodedBundles),
		}
		encodedBlobs[z] = encodedBlob

		for _, securityParam := range blob.RequestHeader.SecurityParams {

			quorumID := securityParam.QuorumID
			quorums := []core.QuorumID{quorumID}

			state, err := cst.GetOperatorState(context.Background(), bn, quorums)
			if err != nil {
				t.Fatal(err)
			}

			blobSize := uint(len(blob.Data))
			blobLength := encoding.GetBlobLength(blobSize)

			chunkLength, err := asn.CalculateChunkLength(state, blobLength, 0, securityParam)
			if err != nil {
				t.Fatal(err)
			}

			quorumHeader := &core.BlobQuorumInfo{
				SecurityParam: core.SecurityParam{
					QuorumID:              quorumID,
					AdversaryThreshold:    securityParam.AdversaryThreshold,
					ConfirmationThreshold: securityParam.ConfirmationThreshold,
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
			bytes := make([][]byte, 0, len(chunks))
			for _, c := range chunks {
				serialized, err := c.Serialize()
				if err != nil {
					t.Fatal(err)
				}
				bytes = append(bytes, serialized)
			}

			blobHeader.BlobCommitments = encoding.BlobCommitments{
				Commitment:       commitments.Commitment,
				LengthCommitment: commitments.LengthCommitment,
				LengthProof:      commitments.LengthProof,
				Length:           commitments.Length,
			}

			blobHeader.QuorumInfos = append(blobHeader.QuorumInfos, quorumHeader)

			for id, assignment := range assignments {
				chunksData := &core.ChunksData{
					Format:   core.GobChunkEncodingFormat,
					ChunkLen: int(chunkLength),
					Chunks:   bytes[assignment.StartIndex : assignment.StartIndex+assignment.NumChunks],
				}
				_, ok := encodedBlob.EncodedBundlesByOperator[id]
				if !ok {
					encodedBlob.EncodedBundlesByOperator[id] = map[core.QuorumID]*core.ChunksData{
						quorumID: chunksData,
					}
				} else {
					encodedBlob.EncodedBundlesByOperator[id][quorumID] = chunksData
				}
			}

		}

	}

	// Set the batch root

	_, err = batchHeader.SetBatchRoot(blobHeaders)
	if err != nil {
		t.Fatal(err)
	}

	return encodedBlobs, batchHeader, cst

}

// checkBatchByUniversalVerifier runs the verification logic for each DA node in the current OperatorState, and returns an error if any of
// the DA nodes' validation checks fails
func checkBatchByUniversalVerifier(cst core.IndexedChainState, encodedBlobs []core.EncodedBlob, header core.BatchHeader, pool common.WorkerPool) error {
	val := core.NewShardValidator(v, asn, cst, [32]byte{})

	quorums := []core.QuorumID{0, 1}
	state, _ := cst.GetIndexedOperatorState(context.Background(), header.ReferenceBlockNumber, quorums)
	numBlob := len(encodedBlobs)

	var errList *multierror.Error

	for id := range state.IndexedOperators {
		val.UpdateOperatorID(id)
		blobMessages := make([]*core.BlobMessage, numBlob)
		for z, encodedBlob := range encodedBlobs {
			bundles, err := new(core.Bundles).FromEncodedBundles(encodedBlob.EncodedBundlesByOperator[id])
			if err != nil {
				return err
			}
			blobMessages[z] = &core.BlobMessage{
				BlobHeader: encodedBlob.BlobHeader,
				Bundles:    bundles,
			}
		}
		err := val.ValidateBatch(&header, blobMessages, state.OperatorState, pool)
		if err != nil {
			errList = multierror.Append(errList, err)
		}
	}

	return errList.ErrorOrNil()

}

func TestValidationSucceeds(t *testing.T) {

	operatorCounts := []uint{1, 2, 4, 10, 30}

	numBlob := 3 // must be greater than 0
	blobLengths := []int{1, 64, 1000}

	securityParams := []*core.SecurityParam{
		{
			QuorumID:              0,
			AdversaryThreshold:    50,
			ConfirmationThreshold: 100,
		},
		{
			QuorumID:              1,
			AdversaryThreshold:    80,
			ConfirmationThreshold: 90,
		},
	}

	bn := uint(0)

	pool := workerpool.New(1)

	for _, operatorCount := range operatorCounts {

		// batch can only be tested per operatorCount, because the assignment would be wrong otherwise
		blobs := make([]core.Blob, 0)
		for _, blobLength := range blobLengths {
			for i := 0; i < numBlob; i++ {
				blobs = append(blobs, makeTestBlob(t, blobLength, securityParams))
			}
		}

		blobMessages, header, cst := prepareBatch(t, operatorCount, blobs, bn)

		t.Run(fmt.Sprintf("universal verifier operatorCount=%v over %v blobs", operatorCount, len(blobs)), func(t *testing.T) {
			err := checkBatchByUniversalVerifier(cst, blobMessages, header, pool)
			assert.NoError(t, err)
		})

	}

}

func TestImproperBatchHeader(t *testing.T) {

	operatorCount := uint(10)

	numBlob := 3 // must be greater than 0
	blobLengths := []int{1, 64, 1000}

	securityParams := []*core.SecurityParam{
		{
			QuorumID:              0,
			AdversaryThreshold:    50,
			ConfirmationThreshold: 100,
		},
		{
			QuorumID:              1,
			AdversaryThreshold:    80,
			ConfirmationThreshold: 90,
		},
	}

	bn := uint(0)

	pool := workerpool.New(1)

	// batch can only be tested per operatorCount, because the assignment would be wrong otherwise
	blobs := make([]core.Blob, 0)
	for _, blobLength := range blobLengths {
		for i := 0; i < numBlob; i++ {
			blobs = append(blobs, makeTestBlob(t, blobLength, securityParams))
		}
	}

	blobMessages, header, cst := prepareBatch(t, operatorCount, blobs, bn)

	// Leave out a blob
	err := checkBatchByUniversalVerifier(cst, blobMessages[:len(blobMessages)-2], header, pool)
	assert.Error(t, err)

	// Add an extra blob
	headers := make([]*core.BlobHeader, len(blobs)-1)
	for i := range headers {
		headers[i] = blobMessages[i].BlobHeader
	}

	_, err = header.SetBatchRoot(headers)
	assert.NoError(t, err)

	err = checkBatchByUniversalVerifier(cst, blobMessages, header, pool)
	assert.Error(t, err)

}
