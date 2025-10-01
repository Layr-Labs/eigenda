package v2_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/mock"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/committer"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/test"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gammazero/workerpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dat *mock.ChainDataMock
	agg core.SignatureAggregator

	p *prover.Prover
	c *committer.Committer
	v *verifier.Verifier

	GETTYSBURG_ADDRESS_BYTES = []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")

	blobParams = &core.BlobVersionParameters{
		NumChunks:       8192,
		CodingRate:      8,
		MaxNumOperators: 2048,
	}
	blobParamsMap = corev2.NewBlobVersionParameterMap(map[corev2.BlobVersion]*core.BlobVersionParameters{
		0: blobParams,
	})
	quorumNumbers = []core.QuorumID{0, 1, 2}
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
	logger := test.GetLogger()
	reader := &mock.MockWriter{}
	reader.On("OperatorIDToAddress").Return(gethcommon.Address{}, nil)
	agg, err = core.NewStdSignatureAggregator(logger, reader)
	if err != nil {
		panic(err)
	}

	p, c, v, err = makeTestComponents()
	if err != nil {
		panic("failed to start localstack container: " + err.Error())
	}

	code := m.Run()
	os.Exit(code)
}

// makeTestComponents makes a prover and verifier currently using the only supported backend.
func makeTestComponents() (*prover.Prover, *committer.Committer, *verifier.Verifier, error) {
	proverConfig := &prover.KzgConfig{
		G1Path:          "../../resources/srs/g1.point",
		G2Path:          "../../resources/srs/g2.point",
		G2TrailingPath:  "../../resources/srs/g2.trailing.point",
		CacheDir:        "../../resources/srs/SRSTables",
		SRSNumberToLoad: 8192,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    true,
	}
	verifierConfig := &verifier.KzgConfig{
		SRSNumberToLoad: 8192,
		G1Path:          "../../resources/srs/g1.point",
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	p, err := prover.NewProver(proverConfig, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("new prover: %w", err)
	}

	c, err := committer.New(p.Srs.G1, p.Srs.G2, p.G2Trailing)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("new committer: %w", err)
	}

	v, err := verifier.NewVerifier(verifierConfig, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("new verifier: %w", err)
	}

	return p, c, v, nil
}

func makeTestBlob(
	t *testing.T, p *prover.Prover, version corev2.BlobVersion, length int, quorums []core.QuorumID,
) (corev2.BlobCertificate, []byte) {

	data := make([]byte, length*31)
	_, err := rand.Read(data)
	if err != nil {
		t.Fatal(err)
	}

	data = codec.ConvertByPaddingEmptyByte(data)

	commitments, err := c.GetCommitmentsForPaddedLength(data)
	if err != nil {
		t.Fatal(err)
	}

	header := corev2.BlobCertificate{
		BlobHeader: &corev2.BlobHeader{
			BlobVersion:     version,
			QuorumNumbers:   quorums,
			BlobCommitments: commitments,
			PaymentMetadata: core.PaymentMetadata{
				AccountID:         gethcommon.HexToAddress("0x123"),
				Timestamp:         5,
				CumulativePayment: big.NewInt(100),
			},
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
	t.Helper()
	ctx := t.Context()

	cst, err := mock.MakeChainDataMock(map[uint8]int{
		0: int(operatorCount),
		1: int(operatorCount),
		2: int(operatorCount) / 2,
	})
	assert.NoError(t, err)

	blobsMap := make(map[core.OperatorID][]*corev2.BlobShard)

	for z := range certs {
		cert := certs[z]
		blob := blobs[z]
		header := cert.BlobHeader

		params, err := corev2.GetEncodingParams(header.BlobCommitments.Length, blobParams)
		require.NoError(t, err)
		chunks, err := p.GetFrames(blob, params)
		require.NoError(t, err)
		state, err := cst.GetOperatorState(ctx, uint(referenceBlockNumber), header.QuorumNumbers)

		require.NoError(t, err)

		assignments, err := corev2.GetAssignmentsForBlob(state, blobParams, header.QuorumNumbers)
		require.NoError(t, err)

		for opID, assignment := range assignments {

			if _, ok := blobsMap[opID]; !ok {
				blobsMap[opID] = make([]*corev2.BlobShard, 0)
			}

			shard := &corev2.BlobShard{
				BlobCertificate: &cert,
				Bundle:          make([]*encoding.Frame, assignment.NumChunks()),
			}
			for i := uint32(0); i < assignment.NumChunks(); i++ {
				shard.Bundle[i] = chunks[assignment.Indices[i]]
			}

			blobsMap[opID] = append(blobsMap[opID], shard)
		}
	}

	return blobsMap, cst

}

// checkBatchByUniversalVerifier runs the verification logic for each DA node in the current OperatorState, and returns an error if any of
// the DA nodes' validation checks fails
func checkBatchByUniversalVerifier(
	t *testing.T,
	cst core.IndexedChainState,
	packagedBlobs map[core.OperatorID][]*corev2.BlobShard,
	pool common.WorkerPool,
) {
	t.Helper()
	ctx := t.Context()
	state, _ := cst.GetIndexedOperatorState(ctx, 0, quorumNumbers)

	for id := range state.IndexedOperators {
		val := corev2.NewShardValidator(v, id, test.GetLogger())
		blobs := packagedBlobs[id]
		st, err := cst.GetOperatorState(ctx, 0, quorumNumbers)
		require.NoError(t, err)
		err = val.ValidateBlobs(ctx, blobs, blobParamsMap, pool, st)
		require.NoError(t, err)
	}
}

func TestValidationSucceeds(t *testing.T) {

	operatorCounts := []uint{2, 10}
	numBlob := 1 // must be greater than 0
	blobLengths := []int{1, 2}
	bn := uint64(1000)

	version := corev2.BlobVersion(0)

	pool := workerpool.New(1)

	for _, operatorCount := range operatorCounts {

		// batch can only be tested per operatorCount, because the assignment would be wrong otherwise
		certs := make([]corev2.BlobCertificate, 0)
		blobs := make([][]byte, 0)
		for _, blobLength := range blobLengths {
			for i := 0; i < numBlob; i++ {
				cert, data := makeTestBlob(t, p, version, blobLength, quorumNumbers)
				certs = append(certs, cert)
				blobs = append(blobs, data)
			}
		}

		packagedBlobs, cst := prepareBlobs(t, operatorCount, certs, blobs, bn)

		t.Run(fmt.Sprintf("universal verifier operatorCount=%v over %v blobs", operatorCount, len(blobs)), func(t *testing.T) {
			checkBatchByUniversalVerifier(t, cst, packagedBlobs, pool)
		})

	}

}
