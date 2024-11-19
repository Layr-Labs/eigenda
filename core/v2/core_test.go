package v2_test

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
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gammazero/workerpool"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

var (
	dat *mock.ChainDataMock
	agg core.SignatureAggregator

	p encoding.Prover
	v encoding.Verifier

	GETTYSBURG_ADDRESS_BYTES = []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")
)

func TestMain(m *testing.M) {
	var err error
	dat, err = mock.MakeChainDataMock(map[uint8]int{
		0: 6,
		1: 3,
	})
	if err != nil {
		panic(err)
	}
	logger := logging.NewNoopLogger()
	reader := &mock.MockWriter{}
	reader.On("OperatorIDToAddress").Return(gethcommon.Address{}, nil)
	agg, err = core.NewStdSignatureAggregator(logger, reader)
	if err != nil {
		panic(err)
	}

	p, v, err = makeTestComponents()
	if err != nil {
		panic("failed to start localstack container")
	}

	code := m.Run()
	os.Exit(code)
}

// makeTestComponents makes a prover and verifier currently using the only supported backend.
func makeTestComponents() (encoding.Prover, encoding.Verifier, error) {
	config := &kzg.KzgConfig{
		G1Path:          "../../inabox/resources/kzg/g1.point.300000",
		G2Path:          "../../inabox/resources/kzg/g2.point.300000",
		CacheDir:        "../../inabox/resources/kzg/SRSTables",
		SRSOrder:        8192,
		SRSNumberToLoad: 8192,
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

func makeTestBlob(t *testing.T, p encoding.Prover, version corev2.BlobVersion, length int, quorums []core.QuorumID) (corev2.BlobCertificate, []byte) {

	data := make([]byte, length*31)
	_, err := rand.Read(data)
	if err != nil {
		t.Fatal(err)
	}

	data = codec.ConvertByPaddingEmptyByte(data)

	commitments, err := p.GetCommitments(data)
	if err != nil {
		t.Fatal(err)
	}

	header := corev2.BlobCertificate{
		BlobHeader: &corev2.BlobHeader{
			BlobVersion:     version,
			QuorumNumbers:   quorums,
			BlobCommitments: commitments,
		},
	}

	return header, data

}

// prepareBlobs takes in multiple blob, encodes them, generates the associated assignments, and the batch header.
// These are the products that a disperser will need in order to disperse data to the DA nodes.
func prepareBlobs(
	t *testing.T,
	operatorCount uint,
	certs []corev2.BlobCertificate,
	blobs [][]byte,
	referenceBlockNumber uint64,
) (map[core.OperatorID][]*corev2.BlobShard, core.IndexedChainState) {

	cst, err := mock.MakeChainDataMock(map[uint8]int{
		0: int(operatorCount),
		1: int(operatorCount),
		2: int(operatorCount),
	})
	assert.NoError(t, err)

	blobsMap := make([]map[core.QuorumID]map[core.OperatorID][]*encoding.Frame, 0, len(certs))

	for z, cert := range certs {
		blob := blobs[z]
		header := cert.BlobHeader

		params, err := header.GetEncodingParams()
		if err != nil {
			t.Fatal(err)
		}

		chunks, err := p.GetFrames(blob, params)
		if err != nil {
			t.Fatal(err)
		}

		state, err := cst.GetOperatorState(context.Background(), uint(referenceBlockNumber), header.QuorumNumbers)
		if err != nil {
			t.Fatal(err)
		}

		blobMap := make(map[core.QuorumID]map[core.OperatorID][]*encoding.Frame)

		for _, quorum := range header.QuorumNumbers {

			assignments, err := corev2.GetAssignments(state, header.BlobVersion, quorum)
			if err != nil {
				t.Fatal(err)
			}

			blobMap[quorum] = make(map[core.OperatorID][]*encoding.Frame)

			for opID, assignment := range assignments {

				blobMap[quorum][opID] = chunks[assignment.StartIndex : assignment.StartIndex+assignment.NumChunks]

			}

		}

		blobsMap = append(blobsMap, blobMap)
	}

	// Invert the blobsMap
	inverseMap := make(map[core.OperatorID][]*corev2.BlobShard)
	for blobIndex, blobMap := range blobsMap {
		for quorum, operatorMap := range blobMap {
			for operatorID, frames := range operatorMap {

				if _, ok := inverseMap[operatorID]; !ok {
					inverseMap[operatorID] = make([]*corev2.BlobShard, 0)
				}
				if len(inverseMap[operatorID]) < blobIndex+1 {
					inverseMap[operatorID] = append(inverseMap[operatorID], &corev2.BlobShard{
						BlobCertificate: &certs[blobIndex],
						Bundles:         make(map[core.QuorumID]core.Bundle),
					})
				}
				if len(frames) == 0 {
					continue
				}
				inverseMap[operatorID][blobIndex].Bundles[quorum] = append(inverseMap[operatorID][blobIndex].Bundles[quorum], frames...)

			}
		}
	}

	return inverseMap, cst

}

// checkBatchByUniversalVerifier runs the verification logic for each DA node in the current OperatorState, and returns an error if any of
// the DA nodes' validation checks fails
func checkBatchByUniversalVerifier(
	cst core.IndexedChainState,
	packagedBlobs map[core.OperatorID][]*corev2.BlobShard,
	pool common.WorkerPool,
) error {

	ctx := context.Background()

	quorums := []core.QuorumID{0, 1}
	state, _ := cst.GetIndexedOperatorState(context.Background(), 0, quorums)

	var errList *multierror.Error

	for id := range state.IndexedOperators {

		val := corev2.NewShardValidator(v, id)

		blobs := packagedBlobs[id]

		err := val.ValidateBlobs(ctx, blobs, pool, state.OperatorState)
		if err != nil {
			errList = multierror.Append(errList, err)
		}
	}

	return errList.ErrorOrNil()

}

func TestValidationSucceeds(t *testing.T) {

	// operatorCounts := []uint{1, 2, 4, 10, 30}

	// numBlob := 3 // must be greater than 0
	// blobLengths := []int{1, 32, 128}

	operatorCounts := []uint{4}

	numBlob := 1 // must be greater than 0
	blobLengths := []int{1, 2}

	quorumNumbers := []core.QuorumID{0, 1}

	bn := uint64(1000)

	version := corev2.BlobVersion(0)

	pool := workerpool.New(1)

	for _, operatorCount := range operatorCounts {

		// batch can only be tested per operatorCount, because the assignment would be wrong otherwise
		headers := make([]corev2.BlobCertificate, 0)
		blobs := make([][]byte, 0)
		for _, blobLength := range blobLengths {
			for i := 0; i < numBlob; i++ {
				header, data := makeTestBlob(t, p, version, blobLength, quorumNumbers)
				headers = append(headers, header)
				blobs = append(blobs, data)
			}
		}

		packagedBlobs, cst := prepareBlobs(t, operatorCount, headers, blobs, bn)

		t.Run(fmt.Sprintf("universal verifier operatorCount=%v over %v blobs", operatorCount, len(blobs)), func(t *testing.T) {
			err := checkBatchByUniversalVerifier(cst, packagedBlobs, pool)
			assert.NoError(t, err)
		})

	}

}
