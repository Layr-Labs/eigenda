package e2e

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/config"
	proxy_metrics "github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/server"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/redis"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/s3"
	"github.com/Layr-Labs/eigenda-proxy/verify/v1"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/exp/rand"

	miniotc "github.com/testcontainers/testcontainers-go/modules/minio"
	redistc "github.com/testcontainers/testcontainers-go/modules/redis"
)

const (
	privateKey         = "SIGNER_PRIVATE_KEY"
	ethRPC             = "ETHEREUM_RPC"
	transport          = "http"
	svcName            = "eigenda_proxy"
	host               = "127.0.0.1"
	v1DisperserHolesky = "disperser-holesky.eigenda.xyz:443"
)

var (
	// set by startMinioContainer
	minioEndpoint = ""
	// set by startRedisContainer
	redisEndpoint = ""
)

// TODO: we shouldn't start the containers in the init function like this.
// Need to find a better way to start the containers and set the endpoints.
// Even better would be for the endpoints not to be global variables injected into the test configs.
// Starting the containers on init like this also makes it harder to import this file into other tests.
func init() {
	err := startMinIOContainer()
	if err != nil {
		panic(err)
	}
	err = startRedisContainer()
	if err != nil {
		panic(err)
	}
}

// startMinIOContainer starts a MinIO container and sets the minioEndpoint global variable
func startMinIOContainer() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	minioContainer, err := miniotc.Run(
		ctx,
		"minio/minio:RELEASE.2024-10-02T17-50-41Z",
		miniotc.WithUsername("minioadmin"),
		miniotc.WithPassword("minioadmin"),
	)
	if err != nil {
		return fmt.Errorf("failed to start MinIO container: %w", err)
	}

	endpoint, err := minioContainer.Endpoint(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to get MinIO endpoint: %w", err)
	}

	minioEndpoint = strings.TrimPrefix(endpoint, "http://")
	return nil
}

// startRedisContainer starts a Redis container and sets the redisEndpoint global variable
func startRedisContainer() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	redisContainer, err := redistc.Run(
		ctx,
		"docker.io/redis:7",
	)
	if err != nil {
		return fmt.Errorf("failed to start Redis container: %w", err)
	}

	endpoint, err := redisContainer.Endpoint(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to get Redis endpoint: %w", err)
	}
	redisEndpoint = endpoint
	return nil
}

type Cfg struct {
	UseV2            bool
	UseMemory        bool
	Expiration       time.Duration
	WriteThreadCount int
	// at most one of the below options should be true
	UseKeccak256ModeS3 bool
	UseS3Caching       bool
	UseRedisCaching    bool
	UseS3Fallback      bool
}

func TestConfig(useMemory bool, useV2 bool) Cfg {
	return Cfg{
		UseV2:              useV2,
		UseMemory:          useMemory,
		Expiration:         14 * 24 * time.Hour,
		UseKeccak256ModeS3: false,
		UseS3Caching:       false,
		UseRedisCaching:    false,
		UseS3Fallback:      false,
		WriteThreadCount:   0,
	}
}

func createRedisConfig(eigendaCfg config.ProxyConfig) config.AppConfig {
	eigendaCfg.StorageConfig.RedisConfig = redis.Config{
		Endpoint: redisEndpoint,
		Password: "",
		DB:       0,
		Eviction: 10 * time.Minute,
	}
	return config.AppConfig{
		EigenDAConfig: eigendaCfg,
	}
}

func createS3Config(eigendaCfg config.ProxyConfig) config.AppConfig {
	// generate random string
	bucketName := "eigenda-proxy-test-" + RandStr(10)
	createS3Bucket(bucketName)

	eigendaCfg.StorageConfig.S3Config = s3.Config{
		Bucket:          bucketName,
		Path:            "",
		Endpoint:        minioEndpoint,
		EnableTLS:       false,
		AccessKeySecret: "minioadmin",
		AccessKeyID:     "minioadmin",
		CredentialType:  s3.CredentialTypeStatic,
	}
	return config.AppConfig{
		EigenDAConfig: eigendaCfg,
	}
}

func TestSuiteConfig(testCfg Cfg) config.AppConfig {
	// load signer key from environment
	pk := os.Getenv(privateKey)
	if pk == "" && !testCfg.UseMemory {
		panic("SIGNER_PRIVATE_KEY environment variable not set")
	}

	// load node url from environment
	ethRPC := os.Getenv(ethRPC)
	if ethRPC == "" && !testCfg.UseMemory {
		panic("ETHEREUM_RPC environment variable is not set")
	}

	var pollInterval time.Duration
	if testCfg.UseMemory {
		pollInterval = time.Second * 1
	} else {
		pollInterval = time.Minute * 1
	}

	maxBlobLengthBytes, err := common.ParseBytesAmount("16mib")
	if err != nil {
		panic(err)
	}

	svcManagerAddr := "0xD4A7E1Bd8015057293f0D0A557088c286942e84b" // holesky testnet
	eigendaCfg := config.ProxyConfig{
		ServerConfig: server.Config{
			DisperseV2: testCfg.UseV2,
			Host:       host,
			Port:       0,
		},
		EdaClientConfigV1: clients.EigenDAClientConfig{
			RPC:                      v1DisperserHolesky,
			StatusQueryTimeout:       time.Minute * 45,
			StatusQueryRetryInterval: pollInterval,
			DisableTLS:               false,
			SignerPrivateKeyHex:      pk,
			EthRpcUrl:                ethRPC,
			SvcManagerAddr:           svcManagerAddr,
		},
		EdaVerifierConfigV1: verify.Config{
			VerifyCerts:          false,
			RPCURL:               ethRPC,
			SvcManagerAddr:       svcManagerAddr,
			EthConfirmationDepth: 1,
			KzgConfig: &kzg.KzgConfig{
				G1Path:          "../resources/g1.point",
				G2PowerOf2Path:  "../resources/g2.point.powerOf2",
				CacheDir:        "../resources/SRSTables",
				SRSOrder:        268435456,
				SRSNumberToLoad: maxBlobLengthBytes / 32,
				NumWorker:       uint64(runtime.GOMAXPROCS(0)), // #nosec G115
			},
			WaitForFinalization: false,
		},
		MemstoreEnabled: testCfg.UseMemory,
		MemstoreConfig: memconfig.NewSafeConfig(
			memconfig.Config{
				BlobExpiration:   testCfg.Expiration,
				MaxBlobSizeBytes: maxBlobLengthBytes,
			}),

		EdaClientConfigV2: common.ClientConfigV2{
			Enabled: testCfg.UseV2,
		},
		StorageConfig: store.Config{
			AsyncPutWorkers: testCfg.WriteThreadCount,
		},

		EigenDAV2Enabled: testCfg.UseV2,
	}

	if testCfg.UseMemory {
		eigendaCfg.EdaClientConfigV1.SignerPrivateKeyHex = "0000000000000000000100000000000000000000000000000000000000000000"
	}

	var cfg config.AppConfig
	switch {
	case testCfg.UseKeccak256ModeS3:
		cfg = createS3Config(eigendaCfg)

	case testCfg.UseS3Caching:
		eigendaCfg.StorageConfig.CacheTargets = []string{"S3"}
		cfg = createS3Config(eigendaCfg)

	case testCfg.UseS3Fallback:
		eigendaCfg.StorageConfig.FallbackTargets = []string{"S3"}
		cfg = createS3Config(eigendaCfg)

	case testCfg.UseRedisCaching:
		eigendaCfg.StorageConfig.CacheTargets = []string{"redis"}
		cfg = createRedisConfig(eigendaCfg)

	default:
		cfg = config.AppConfig{
			EigenDAConfig: eigendaCfg,
			MetricsCfg:    proxy_metrics.Config{},
		}
	}

	return cfg
}

func TestSuiteSecretConfig(testCfg Cfg) common.SecretConfigV2 {
	// load signer key from environment
	signerPrivateKey := os.Getenv(privateKey)
	if signerPrivateKey == "" && !testCfg.UseMemory {
		panic("SIGNER_PRIVATE_KEY environment variable not set")
	}

	return common.SecretConfigV2{
		SignerPaymentKey: signerPrivateKey,
	}
}

type TestSuite struct {
	Ctx     context.Context
	Log     logging.Logger
	Server  *server.Server
	Metrics *proxy_metrics.EmulatedMetricer
}

func CreateTestSuite(testSuiteCfg config.AppConfig, secretConfigV2 common.SecretConfigV2) (TestSuite, func()) {
	log := logging.NewTextSLogger(os.Stdout, &logging.SLoggerOptions{})

	metrics := proxy_metrics.NewEmulatedMetricer()
	ctx := context.Background()

	storageManager, err := store.NewStorageManagerBuilder(
		ctx,
		log,
		metrics,
		testSuiteCfg.EigenDAConfig.StorageConfig,
		testSuiteCfg.EigenDAConfig.EdaVerifierConfigV1,
		testSuiteCfg.EigenDAConfig.EdaClientConfigV1,
		testSuiteCfg.EigenDAConfig.EdaClientConfigV2,
		secretConfigV2,
		testSuiteCfg.EigenDAConfig.MemstoreConfig,
		testSuiteCfg.EigenDAConfig.PutRetries,
		testSuiteCfg.EigenDAConfig.MaxBlobSizeBytes,
	).Build(ctx)
	if err != nil {
		panic(err)
	}

	proxySvr := server.NewServer(testSuiteCfg.EigenDAConfig.ServerConfig, storageManager, log, metrics)

	log.Info("Starting proxy server...")
	r := mux.NewRouter()
	proxySvr.RegisterRoutes(r)
	err = proxySvr.Start(r)
	if err != nil {
		panic(err)
	}

	if testSuiteCfg.EigenDAConfig.MemstoreEnabled {
		memconfig.NewHandlerHTTP(log, testSuiteCfg.EigenDAConfig.MemstoreConfig).RegisterMemstoreConfigHandlers(r)
	}

	kill := func() {
		if err := proxySvr.Stop(); err != nil {
			log.Error("failed to stop proxy server", "err", err)
		}
	}

	return TestSuite{
		Ctx:     ctx,
		Log:     log,
		Server:  proxySvr,
		Metrics: metrics,
	}, kill
}

func (ts *TestSuite) Address() string {
	// read port from listener
	port := ts.Server.Port()

	return fmt.Sprintf("%s://%s:%d", transport, host, port)
}

func createS3Bucket(bucketName string) {
	// Initialize minio client object.
	endpoint := minioEndpoint
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false

	minioClient, err := minio.New(
		endpoint, &minio.Options{
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

func RandStr(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandBytes(n int) []byte {
	return []byte(RandStr(n))
}
