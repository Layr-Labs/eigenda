package integration_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	grpcdisperser "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/api/grpc/retriever"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// This should match with the configuration used for the disperser
	perUserThroughput   = 1000 // Bytes per second
	dispersalMultiplier = 1
	retrievalMultiplier = 2 // https://github.com/Layr-Labs/eigenda/blob/master/node/cmd/main.go#L25
)

type result struct {
	*grpcdisperser.BlobInfo
	data []byte
	err  error
}

func disperse(t *testing.T, ctx context.Context, client clients.DisperserClient, resultChan chan result, data []byte, param core.SecurityParam) {

	blobStatus, key, err := client.DisperseBlob(ctx, data, []uint8{param.QuorumID})
	if err != nil {
		resultChan <- result{
			err: err,
		}
		return
	}

	assert.NoError(t, err)
	assert.NotNil(t, key)
	assert.NotNil(t, blobStatus)
	assert.Equal(t, *blobStatus, disperser.Processing)

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			t.Error("timed out waiting for dispersed blob to confirm")
			return
		case <-ticker.C:
			reply, err := client.GetBlobStatus(ctx, key)
			assert.NoError(t, err)
			assert.NotNil(t, reply)
			blobStatus, err = disperser.FromBlobStatusProto(reply.GetStatus())
			assert.NoError(t, err)
			if *blobStatus == disperser.Confirmed {
				blobInfo := reply.GetInfo()

				resultChan <- result{
					BlobInfo: blobInfo,
					data:     data,
				}

				return
			}
		}
	}
}

func retrieve(t *testing.T, ctx context.Context, client retriever.RetrieverClient, result result) error {

	reply, err := client.RetrieveBlob(ctx, &retriever.BlobRequest{
		BatchHeaderHash:      result.BlobVerificationProof.BatchMetadata.BatchHeaderHash,
		BlobIndex:            result.BlobVerificationProof.BlobIndex,
		ReferenceBlockNumber: result.BlobVerificationProof.BatchMetadata.ConfirmationBlockNumber,
		QuorumId:             result.BlobHeader.BlobQuorumParams[0].QuorumNumber,
	})
	if err != nil {
		return err
	}
	assert.NotNil(t, reply)
	if reply != nil {
		assert.Equal(t, result.data, reply.Data[:len(result.data)])
	}
	return nil
}

type ratelimitTestCase struct {
	numDispersal      int
	numRetrieval      int
	dispersalInterval time.Duration
	retrievalInterval time.Duration
	pause             time.Duration
	blobSize          int
	param             core.SecurityParam
}

func testRatelimit(t *testing.T, testConfig *deploy.Config, c ratelimitTestCase) (int, int) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	disp, err := clients.NewDisperserClient(&clients.Config{
		Hostname: "localhost",
		Port:     testConfig.Dispersers[0].DISPERSER_SERVER_GRPC_PORT,
		Timeout:  10 * time.Second,
	}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, disp)

	data := make([]byte, c.blobSize)
	_, err = rand.Read(data)
	assert.NoError(t, err)

	dispersalTicker := time.NewTicker(c.dispersalInterval)
	defer dispersalTicker.Stop()
	resultChan := make(chan result, c.numDispersal)
	go func() {
		for i := 0; i < c.numDispersal; i++ {
			<-dispersalTicker.C
			go disperse(t, ctx, disp, resultChan, data, c.param)
		}
	}()

	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%v", testConfig.Retriever.RETRIEVER_GRPC_PORT),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)
	defer func() { _ = conn.Close() }()
	ret := retriever.NewRetrieverClient(conn)
	assert.NotNil(t, ret)

	dispersalErrors := 0
	retrievalErrors := 0

	time.Sleep(c.pause)
	retrievalTicker := time.NewTicker(c.retrievalInterval)
	defer retrievalTicker.Stop()

	for i := 0; i < c.numDispersal; i++ {

		if i < c.numRetrieval {
			<-retrievalTicker.C

			result := <-resultChan
			if result.err != nil {
				dispersalErrors++
			} else {
				err := retrieve(t, ctx, ret, result)
				if err != nil {
					retrievalErrors++
				}
			}
		} else {
			result := <-resultChan
			if result.err != nil {
				dispersalErrors++
			}
		}
	}

	return dispersalErrors, retrievalErrors

}

func TestRatelimit(t *testing.T) {

	t.Skip("Manual test for now")

	rootPath := "../../"
	testname, err := deploy.GetLatestTestDirectory(rootPath)
	if err != nil {
		t.Fatal(err)
	}
	testConfig := deploy.NewTestConfig(testname, rootPath)

	if testConfig.Dispersers[0].DISPERSER_SERVER_PER_USER_UNAUTH_BYTE_RATE != fmt.Sprint(perUserThroughput) {
		t.Fatalf("per user throughput should be %v", perUserThroughput)
	}
	if testConfig.Dispersers[0].DISPERSER_SERVER_BUCKET_MULTIPLIERS != fmt.Sprint(dispersalMultiplier) {
		t.Fatalf("dispersal multiplier should be %v", dispersalMultiplier)
	}

	t.Run("no ratelimiting when dispersing and retrieving within rate", func(t *testing.T) {

		t.Skip("Manual test for now")

		testCase := ratelimitTestCase{
			numDispersal:      10,
			numRetrieval:      10,
			dispersalInterval: time.Second,
			retrievalInterval: 500 * time.Millisecond,
			pause:             0,
			blobSize:          500,
			param: core.SecurityParam{
				QuorumID:              0,
				AdversaryThreshold:    50,
				ConfirmationThreshold: 100,
			},
		}

		const encodedBlobSize = 1000

		assert.EqualValues(t, perUserThroughput*dispersalMultiplier/1000, encodedBlobSize/testCase.dispersalInterval.Milliseconds())
		assert.EqualValues(t, perUserThroughput*retrievalMultiplier/1000, encodedBlobSize/testCase.retrievalInterval.Milliseconds())

		dispersalErrors, retrievalErrors := testRatelimit(t, testConfig, testCase)

		assert.Equal(t, 0, dispersalErrors)
		assert.Equal(t, 0, retrievalErrors)

	})

	t.Run("dispersal ratelimiting when dispersing above rate", func(t *testing.T) {

		t.Skip("Manual test for now")

		testCase := ratelimitTestCase{
			numDispersal:      10,
			numRetrieval:      0,
			dispersalInterval: time.Second,
			retrievalInterval: time.Second,
			pause:             0,
			blobSize:          1000,
			param: core.SecurityParam{
				QuorumID:              0,
				AdversaryThreshold:    50,
				ConfirmationThreshold: 100,
			},
		}

		const encodedBlobSize = 2000
		const overageFactor = 2

		assert.EqualValues(t, overageFactor*perUserThroughput*dispersalMultiplier/1000, encodedBlobSize/testCase.dispersalInterval.Milliseconds())

		dispersalErrors, retrievalErrors := testRatelimit(t, testConfig, testCase)

		fmt.Println("Dispersal Ratelimited: ", dispersalErrors)
		assert.Greater(t, dispersalErrors, 0)
		assert.Equal(t, 0, retrievalErrors)

	})

	t.Run("retrieval ratelimiting when retrieving above rate", func(t *testing.T) {

		t.Skip("Manual test for now")

		testCase := ratelimitTestCase{
			numDispersal:      10,
			numRetrieval:      10,
			dispersalInterval: 2 * time.Second,
			retrievalInterval: 500 * time.Millisecond,
			pause:             20 * time.Second,
			blobSize:          1000,
			param: core.SecurityParam{
				QuorumID:              0,
				AdversaryThreshold:    50,
				ConfirmationThreshold: 100,
			},
		}

		const encodedBlobSize = 2000
		const overageFactor = 2

		assert.EqualValues(t, perUserThroughput*dispersalMultiplier/1000, encodedBlobSize/testCase.dispersalInterval.Milliseconds())
		assert.EqualValues(t, overageFactor*perUserThroughput*retrievalMultiplier/1000, encodedBlobSize/testCase.retrievalInterval.Milliseconds())

		dispersalErrors, retrievalErrors := testRatelimit(t, testConfig, testCase)

		fmt.Println("Retrieval Ratelimited: ", retrievalErrors)
		assert.Equal(t, 0, dispersalErrors)
		assert.Greater(t, retrievalErrors, 0)

	})

	t.Run("ratelimiting when dispersing greater than blob rate", func(t *testing.T) {

		t.Skip("Manual test for now")

		testCase := ratelimitTestCase{
			numDispersal:      200,
			numRetrieval:      0,
			dispersalInterval: 450 * time.Millisecond,
			retrievalInterval: 500 * time.Millisecond,
			pause:             0,
			blobSize:          5,
			param: core.SecurityParam{
				QuorumID:              0,
				AdversaryThreshold:    50,
				ConfirmationThreshold: 100,
			},
		}

		dispersalErrors, retrievalErrors := testRatelimit(t, testConfig, testCase)

		fmt.Println("Dispersal Ratelimited: ", dispersalErrors)

		assert.Greater(t, dispersalErrors, 0)
		assert.Equal(t, 0, retrievalErrors)

	})

}
