package integration

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/encoding"
	"github.com/Layr-Labs/eigenda/core/mock"
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

// prepareBatch takes in a single blob, encodes it, generates the associated assignments, and the batch header.
// These are the products that a disperser will need in order to disperse data to the DA nodes.
func prepareBatch(t *testing.T, cst core.IndexedChainState, blob core.Blob, quorumIndex uint, quantizationFactor uint, bn uint) (core.EncodedBlob, core.BatchHeader) {

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

	batchHeader := core.BatchHeader{
		ReferenceBlockNumber: bn,
		BatchRoot:            [32]byte{},
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

	return encodedBlob, batchHeader

}

// checkBatch runs the verification logic for each DA node in the current OperatorState, and returns an error if any of
// the DA nodes' validation checks fails
func checkBatch(t *testing.T, cst core.IndexedChainState, encodedBlob core.EncodedBlob, header core.BatchHeader) {
	val := core.NewChunkValidator(enc, asn, cst, [32]byte{})

	quorums := []core.QuorumID{0}
	state, _ := cst.GetIndexedOperatorState(context.Background(), header.ReferenceBlockNumber, quorums)

	for id := range state.IndexedOperators {

		blobMessage := encodedBlob[id]

		val.UpdateOperatorID(id)
		err := val.ValidateBlob(blobMessage, state.OperatorState)
		assert.NoError(t, err)
	}

}

func TestCoreLibrary(t *testing.T) {

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

	for _, blobLength := range blobLengths {
		for _, quantizationFactor := range quantizationFactors {
			for _, operatorCount := range operatorCounts {
				for _, securityParam := range securityParams {

					t.Run(fmt.Sprintf("blobLength=%v, quantizationFactor=%v, operatorCount=%v, securityParams=%v", blobLength, quantizationFactor, operatorCount, securityParam), func(t *testing.T) {

						blob := makeTestBlob(t, blobLength, []*core.SecurityParam{securityParam})

						cst, err := mock.NewChainDataMock(core.OperatorIndex(operatorCount))
						assert.NoError(t, err)

						quorumIndex := uint(0)
						bn := uint(0)

						batch, header := prepareBatch(t, cst, blob, quorumIndex, quantizationFactor, bn)

						checkBatch(t, cst, batch, header)
					})
				}
			}
		}
	}

}
