package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/server"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/ethereum/go-ethereum/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/exp/rand"

	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	opmetrics "github.com/ethereum-optimism/optimism/op-service/metrics"

	"github.com/stretchr/testify/require"
)

const (
	privateKey = "SIGNER_PRIVATE_KEY"
	ethRPC     = "ETHEREUM_RPC"
	transport  = "http"
	svcName    = "eigenda_proxy"
	host       = "127.0.0.1"
	holeskyDA  = "disperser-holesky.eigenda.xyz:443"
)

type TestSuite struct {
	Ctx    context.Context
	Log    log.Logger
	Server *server.Server
}

func CreateTestSuite(t *testing.T, useMemory bool, useS3 bool) (TestSuite, func()) {

	ctx := context.Background()

	// load signer key from environment
	pk := os.Getenv(privateKey)
	if pk == "" && !useMemory {
		t.Fatal("SIGNER_PRIVATE_KEY environment variable not set")
	}

	// load node url from environment
	ethRPC := os.Getenv(ethRPC)
	if ethRPC == "" && !useMemory {
		t.Fatal("ETHEREUM_RPC environment variable is not set")
	}

	var pollInterval time.Duration
	if useMemory {
		pollInterval = time.Second * 1
	} else {
		pollInterval = time.Minute * 1
	}

	log := oplog.NewLogger(os.Stdout, oplog.CLIConfig{
		Level:  log.LevelDebug,
		Format: oplog.FormatLogFmt,
		Color:  true,
	}).New("role", svcName)

	eigendaCfg := server.Config{
		ClientConfig: clients.EigenDAClientConfig{
			RPC:                      holeskyDA,
			StatusQueryTimeout:       time.Minute * 45,
			StatusQueryRetryInterval: pollInterval,
			DisableTLS:               false,
			SignerPrivateKeyHex:      pk,
		},
		EthRPC:                 ethRPC,
		SvcManagerAddr:         "0xD4A7E1Bd8015057293f0D0A557088c286942e84b", // incompatible with non holeskly networks
		CacheDir:               "../resources/SRSTables",
		G1Path:                 "../resources/g1.point",
		MaxBlobLength:          "2mib",
		G2PowerOfTauPath:       "../resources/g2.point.powerOf2",
		PutBlobEncodingVersion: 0x00,
		MemstoreEnabled:        useMemory,
		MemstoreBlobExpiration: 14 * 24 * time.Hour,
		EthConfirmationDepth:   0,
	}

	if useMemory {
		eigendaCfg.ClientConfig.SignerPrivateKeyHex = "0000000000000000000100000000000000000000000000000000000000000000"
	}

	var cfg server.CLIConfig
	if useS3 {
		// generate random string
		bucketName := "eigenda-proxy-test-" + RandString(10)
		createS3Bucket(bucketName)

		cfg = server.CLIConfig{
			EigenDAConfig: eigendaCfg,
			S3Config: store.S3Config{
				Profiling:        true,
				Bucket:           bucketName,
				Path:             "",
				Endpoint:         "localhost:4566",
				AccessKeySecret:  "minioadmin",
				AccessKeyID:      "minioadmin",
				S3CredentialType: store.S3CredentialStatic,
				Backup:           false,
			},
		}
	} else {

		cfg = server.CLIConfig{
			EigenDAConfig: eigendaCfg,
			MetricsCfg:    opmetrics.CLIConfig{},
		}

	}

	store, err := server.LoadStoreRouter(
		cfg,
		ctx,
		log,
	)
	require.NoError(t, err)
	server := server.NewServer(host, 0, store, log, metrics.NoopMetrics)

	t.Log("Starting proxy server...")
	err = server.Start()
	require.NoError(t, err)

	kill := func() {
		if err := server.Stop(); err != nil {
			panic(err)
		}
	}

	return TestSuite{
		Ctx:    ctx,
		Log:    log,
		Server: server,
	}, kill
}

func (ts *TestSuite) Address() string {
	// read port from listener
	port := ts.Server.Port()

	return fmt.Sprintf("%s://%s:%d", transport, host, port)
}

func createS3Bucket(bucketName string) {
	// Initialize minio client object.
	endpoint := "localhost:4566"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(err)
	}

	location := "us-east-1"

	ctx := context.Background()
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			fmt.Printf("We already own %s\n", bucketName)
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("Successfully created %s\n", bucketName)
	}
}

func RandString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
