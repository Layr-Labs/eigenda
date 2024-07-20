package e2e_test

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/client"
	"github.com/Layr-Labs/eigenda-proxy/e2e"
	"github.com/Layr-Labs/eigenda-proxy/utils"
	op_plasma "github.com/ethereum-optimism/optimism/op-plasma"
	"github.com/stretchr/testify/require"
)

func useMemory() bool {
	return !runTestnetIntegrationTests
}

func TestOptimismClientWithS3Backend(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	ts, kill := e2e.CreateTestSuite(t, useMemory(), true)
	defer kill()

	daClient := op_plasma.NewDAClient(ts.Address(), false, true)

	testPreimage := []byte(e2e.RandString(100))

	commit, err := daClient.SetInput(ts.Ctx, testPreimage)
	require.NoError(t, err)

	preimage, err := daClient.GetInput(ts.Ctx, commit)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)
}

func TestOptimismClientWithEigenDABackend(t *testing.T) {
	// this test asserts that the data can be posted/read to EigenDA with a concurrent S3 backend configured

	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	ts, kill := e2e.CreateTestSuite(t, useMemory(), true)
	defer kill()

	daClient := op_plasma.NewDAClient(ts.Address(), false, false)

	testPreimage := []byte(e2e.RandString(100))

	t.Log("Setting input data on proxy server...")
	commit, err := daClient.SetInput(ts.Ctx, testPreimage)
	require.NoError(t, err)

	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetInput(ts.Ctx, commit)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)
}

func TestProxyClient(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	ts, kill := e2e.CreateTestSuite(t, useMemory(), false)
	defer kill()

	cfg := &client.Config{
		URL: ts.Address(),
	}
	daClient := client.New(cfg)

	testPreimage := []byte(e2e.RandString(100))

	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, testPreimage)
	require.NoError(t, err)

	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)
}

func TestProxyClientWithLargeBlob(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	ts, kill := e2e.CreateTestSuite(t, useMemory(), false)
	defer kill()

	cfg := &client.Config{
		URL: ts.Address(),
	}
	daClient := client.New(cfg)
	//  2MB blob
	testPreimage := []byte(e2e.RandString(2000000))

	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, testPreimage)
	require.NoError(t, err)

	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)
}

func TestProxyClientWithOversizedBlob(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	ts, kill := e2e.CreateTestSuite(t, useMemory(), false)
	defer kill()

	cfg := &client.Config{
		URL: ts.Address(),
	}
	daClient := client.New(cfg)
	//  2MB blob
	testPreimage := []byte(e2e.RandString(200000000))

	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, testPreimage)
	require.Empty(t, blobInfo)
	require.Error(t, err)

	oversizedError := false
	if strings.Contains(err.Error(), "blob is larger than max blob size") {
		oversizedError = true
	}

	if strings.Contains(err.Error(), "blob size cannot exceed 2 MiB") {
		oversizedError = true
	}

	require.True(t, oversizedError)

}

func TestProxyClient_MultiSameContentBlobs_SameBatch(t *testing.T) {
	t.Skip("Skipping test until fix is applied to holesky")


	t.Parallel()

	ts, kill := e2e.CreateTestSuite(t, useMemory(), false)
	defer kill()

	cfg := &client.Config{
		URL: ts.Address(),
	}
	
	errChan := make(chan error, 10)
	var wg sync.WaitGroup

	// disperse 10 blobs with the same content in the same batch
	for i := 0; i < 4; i ++ {
		wg.Add(1)
		go func(){
			defer wg.Done()
			daClient := client.New(cfg)
			testPreimage := []byte("hellooooooooooo world!")
		
			t.Log("Setting input data on proxy server...")
			blobInfo, err := daClient.SetData(ts.Ctx, testPreimage)
			if err != nil {
				errChan <- err
				return
			}
		
			t.Log("Getting input data from proxy server...")
			preimage, err := daClient.GetData(ts.Ctx, blobInfo)
			if err != nil {
				errChan <- err
				return
			}
			
			if !utils.EqualSlices(preimage, testPreimage) {
				errChan <- fmt.Errorf("expected preimage %s, got %s", testPreimage, preimage)
				return
			}
		}()
	}

	timedOut := waitTimeout(&wg, 10*time.Minute)
	if timedOut {
		t.Fatal("timed out waiting for parallel tests to complete")
	}

	if len(errChan) > 0 {
		// iterate over channel and log errors 
		for i := 0; i < len(errChan); i++ {
			err := <-errChan
			t.Log(err.Error())
			t.Fail()
		}
	}
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
    c := make(chan struct{})
    go func() {
        defer close(c)
        wg.Wait()
    }()
    select {
    case <-c:
        return false
    case <-time.After(timeout):
        return true
    }
}