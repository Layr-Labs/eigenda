package e2e

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/server"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/redis"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/s3"
	"github.com/Layr-Labs/eigenda-proxy/utils"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
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
	eigendaCfg.RedisConfig = redis.Config{
		Endpoint: "127.0.0.1:9001",
		Password: "",
		DB:       0,
		Eviction: 10 * time.Minute,
		Profile:  true,
	}
	return server.CLIConfig{
		EigenDAConfig: eigendaCfg,
	}
}

func createS3Config(eigendaCfg server.Config) server.CLIConfig {
	// generate random string
	bucketName := "eigenda-proxy-test-" + RandString(10)
	createS3Bucket(bucketName)

	eigendaCfg.S3Config = s3.Config{
		Profiling:       true,
		Bucket:          bucketName,
		Path:            "",
		Endpoint:        "localhost:4566",
		AccessKeySecret: "minioadmin",
		AccessKeyID:     "minioadmin",
		CredentialType:  s3.CredentialTypeStatic,
		Backup:          false,
	}
	return server.CLIConfig{
		EigenDAConfig: eigendaCfg,
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

	maxBlobLengthBytes, err := utils.ParseBytesAmount("16mib")
	require.NoError(t, err)
	eigendaCfg := server.Config{
		EdaClientConfig: clients.EigenDAClientConfig{
			RPC:                      holeskyDA,
			StatusQueryTimeout:       time.Minute * 45,
			StatusQueryRetryInterval: pollInterval,
			DisableTLS:               false,
			SignerPrivateKeyHex:      pk,
		},
		VerifierConfig: verify.Config{
			VerifyCerts:          false,
			RPCURL:               ethRPC,
			SvcManagerAddr:       "0xD4A7E1Bd8015057293f0D0A557088c286942e84b", // incompatible with non holeskly networks
			EthConfirmationDepth: 0,
			KzgConfig: &kzg.KzgConfig{
				G1Path:          "../resources/g1.point",
				G2PowerOf2Path:  "../resources/g2.point.powerOf2",
				CacheDir:        "../resources/SRSTables",
				SRSOrder:        268435456,
				SRSNumberToLoad: maxBlobLengthBytes / 32,
				NumWorker:       uint64(runtime.GOMAXPROCS(0)),
			},
		},
		MemstoreEnabled: testCfg.UseMemory,
		MemstoreConfig: memstore.Config{
			BlobExpiration:   testCfg.Expiration,
			MaxBlobSizeBytes: maxBlobLengthBytes,
		},
	}

	if testCfg.UseMemory {
		eigendaCfg.EdaClientConfig.SignerPrivateKeyHex = "0000000000000000000100000000000000000000000000000000000000000000"
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
