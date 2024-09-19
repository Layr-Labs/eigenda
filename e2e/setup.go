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

type Cfg struct {
	UseMemory  bool
	Expiration time.Duration
	// at most one of the below options should be true
	UseKeccak256ModeS3 bool
	UseS3Caching       bool
	UseRedisCaching    bool
	UseS3Fallback      bool
}

func TestConfig(useMemory bool) *Cfg {
	return &Cfg{
		UseMemory:          useMemory,
		Expiration:         14 * 24 * time.Hour,
		UseKeccak256ModeS3: false,
		UseS3Caching:       false,
		UseRedisCaching:    false,
		UseS3Fallback:      false,
	}
}

func createRedisConfig(eigendaCfg server.Config) server.CLIConfig {
	return server.CLIConfig{
		EigenDAConfig: eigendaCfg,
		RedisCfg: store.RedisConfig{
			Endpoint: "127.0.0.1:9001",
			Password: "",
			DB:       0,
			Eviction: 10 * time.Minute,
			Profile:  true,
		},
	}
}

func createS3Config(eigendaCfg server.Config) server.CLIConfig {
	// generate random string
	bucketName := "eigenda-proxy-test-" + RandString(10)
	createS3Bucket(bucketName)

	return server.CLIConfig{
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
}

func TestSuiteConfig(t *testing.T, testCfg *Cfg) server.CLIConfig {
	// load signer key from environment
	pk := os.Getenv(privateKey)
	if pk == "" && !testCfg.UseMemory {
		t.Fatal("SIGNER_PRIVATE_KEY environment variable not set")
	}

	// load node url from environment
	ethRPC := os.Getenv(ethRPC)
	if ethRPC == "" && !testCfg.UseMemory {
		t.Fatal("ETHEREUM_RPC environment variable is not set")
	}

	var pollInterval time.Duration
	if testCfg.UseMemory {
		pollInterval = time.Second * 1
	} else {
		pollInterval = time.Minute * 1
	}

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
		MaxBlobLength:          "16mib",
		G2PowerOfTauPath:       "../resources/g2.point.powerOf2",
		PutBlobEncodingVersion: 0x00,
		MemstoreEnabled:        testCfg.UseMemory,
		MemstoreBlobExpiration: testCfg.Expiration,
		EthConfirmationDepth:   0,
	}

	if testCfg.UseMemory {
		eigendaCfg.ClientConfig.SignerPrivateKeyHex = "0000000000000000000100000000000000000000000000000000000000000000"
	}

	var cfg server.CLIConfig
	switch {
	case testCfg.UseKeccak256ModeS3:
		cfg = createS3Config(eigendaCfg)

	case testCfg.UseS3Caching:
		eigendaCfg.CacheTargets = []string{"S3"}
		cfg = createS3Config(eigendaCfg)

	case testCfg.UseS3Fallback:
		eigendaCfg.FallbackTargets = []string{"S3"}
		cfg = createS3Config(eigendaCfg)

	case testCfg.UseRedisCaching:
		eigendaCfg.CacheTargets = []string{"redis"}
		cfg = createRedisConfig(eigendaCfg)

	default:
		cfg = server.CLIConfig{
			EigenDAConfig: eigendaCfg,
			MetricsCfg:    opmetrics.CLIConfig{},
		}
	}

	return cfg
}

type TestSuite struct {
	Ctx    context.Context
	Log    log.Logger
	Server *server.Server
}

func CreateTestSuite(t *testing.T, testSuiteCfg server.CLIConfig) (TestSuite, func()) {
	log := oplog.NewLogger(os.Stdout, oplog.CLIConfig{
		Level:  log.LevelDebug,
		Format: oplog.FormatLogFmt,
		Color:  true,
	}).New("role", svcName)

	ctx := context.Background()
	store, err := server.LoadStoreRouter(
		ctx,
		testSuiteCfg,
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
			log.Info(fmt.Sprintf("We already own %s\n", bucketName))
		} else {
			panic(err)
		}
	} else {
		log.Info(fmt.Sprintf("Successfully created %s\n", bucketName))
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
