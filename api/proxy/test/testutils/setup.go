package testutils

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/config"
	"github.com/Layr-Labs/eigenda/api/proxy/config/eigendaflags"
	proxy_metrics "github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/server"
	"github.com/Layr-Labs/eigenda/api/proxy/store"
	"github.com/Layr-Labs/eigenda/api/proxy/store/builder"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/eigenda/verify"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/redis"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/s3"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/ethereum/go-ethereum/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/exp/rand"

	miniotc "github.com/testcontainers/testcontainers-go/modules/minio"
	redistc "github.com/testcontainers/testcontainers-go/modules/redis"
)

const (
	minioAdmin       = "minioadmin"
	backendEnvVar    = "BACKEND"
	privateKeyEnvVar = "SIGNER_PRIVATE_KEY"
	ethRPCEnvVar     = "ETHEREUM_RPC"
	transport        = "http"
	host             = "127.0.0.1"
	disperserPort    = "443"

	disperserPreprodHostname                = "disperser-preprod-holesky.eigenda.xyz"
	preprodCertVerifierAddress              = "0xCCFE3d87fB7D369f1eeE65221a29A83f1323043C"
	preprodSvcManagerAddress                = "0x54A03db2784E3D0aCC08344D05385d0b62d4F432"
	preprodBLSOperatorStateRetrieverAddress = "0x003497Dd77E5B73C40e8aCbB562C8bb0410320E7"
	preprodEigenDADirectory                 = "0xfB676e909f376efFDbDee7F17342aCF55f6Ec502"

	disperserTestnetHostname                = "disperser-testnet-holesky.eigenda.xyz"
	testnetCertVerifierAddress              = "0xd305aeBcdEc21D00fDF8796CE37d0e74836a6B6e"
	testnetSvcManagerAddress                = "0xD4A7E1Bd8015057293f0D0A557088c286942e84b"
	testnetBLSOperatorStateRetrieverAddress = "0x003497Dd77E5B73C40e8aCbB562C8bb0410320E7"
	testnetEigenDADirectory                 = "0x90776Ea0E99E4c38aA1Efe575a61B3E40160A2FE"

	disperserSepoliaHostname                = "disperser-testnet-sepolia.eigenda.xyz"
	sepoliaCertVerifierAddress              = "0x58D2B844a894f00b7E6F9F492b9F43aD54Cd4429"
	sepoliaSvcManagerAddress                = "0x3a5acf46ba6890B8536420F4900AC9BC45Df4764"
	sepoliaBLSOperatorStateRetrieverAddress = "0x22478d082E9edaDc2baE8443E4aC9473F6E047Ff"
	sepoliaEigenDADirectory                 = "0x9620dC4B3564198554e4D2b06dEFB7A369D90257"
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
	// TODO: we should pass in the test.Test here and using t.Context() instead of creating a new context.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	minioContainer, err := miniotc.Run(
		ctx,
		"minio/minio:RELEASE.2024-10-02T17-50-41Z",
		miniotc.WithUsername(minioAdmin),
		miniotc.WithPassword(minioAdmin),
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
	// TODO: we should pass in the test.Test here and using t.Context() instead of creating a new context.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

type Backend int

const (
	TestnetBackend Backend = iota + 1
	PreprodBackend
	SepoliaBackend
	MemstoreBackend
)

// ParseBackend converts a string to a Backend enum (case insensitive)
func ParseBackend(inputString string) (Backend, error) {
	switch strings.ToLower(inputString) {
	case "testnet":
		return TestnetBackend, nil
	case "preprod":
		return PreprodBackend, nil
	case "sepolia":
		return SepoliaBackend, nil
	case "memstore":
		return MemstoreBackend, nil
	default:
		return 0, fmt.Errorf("invalid backend: %s", inputString)
	}
}

func GetBackend() Backend {
	backend, err := ParseBackend(os.Getenv(backendEnvVar))
	if err != nil {
		panic(fmt.Sprintf("BACKEND must be = memstore|testnet|sepolia|preprod. parse backend error: %v", err))
	}
	return backend
}

type TestConfig struct {
	BackendsToEnable []common.EigenDABackend
	DispersalBackend common.EigenDABackend
	Backend          Backend
	Retrievers       []common.RetrieverType
	Expiration       time.Duration
	MaxBlobLength    string
	WriteThreadCount int
	// at most one of the below options should be true
	UseKeccak256ModeS3 bool
	UseS3Caching       bool
	UseRedisCaching    bool
	UseS3Fallback      bool
}

// NewTestConfig returns a new TestConfig
func NewTestConfig(
	backend Backend,
	dispersalBackend common.EigenDABackend,
	// if backendsToEnable is nil, then this method will simply enable whichever backend is being dispersed to
	backendsToEnable []common.EigenDABackend,
) TestConfig {
	if backendsToEnable == nil {
		if dispersalBackend == common.V2EigenDABackend {
			backendsToEnable = []common.EigenDABackend{common.V2EigenDABackend}
		} else {
			backendsToEnable = []common.EigenDABackend{common.V1EigenDABackend}
		}
	}

	return TestConfig{
		BackendsToEnable:   backendsToEnable,
		DispersalBackend:   dispersalBackend,
		Backend:            backend,
		Retrievers:         []common.RetrieverType{common.RelayRetrieverType, common.ValidatorRetrieverType},
		Expiration:         14 * 24 * time.Hour,
		UseKeccak256ModeS3: false,
		UseS3Caching:       false,
		UseRedisCaching:    false,
		UseS3Fallback:      false,
		WriteThreadCount:   0,
	}
}

func createRedisConfig() redis.Config {
	return redis.Config{
		Endpoint: redisEndpoint,
		Password: "",
		DB:       0,
		Eviction: 10 * time.Minute,
	}
}

func createS3Config() s3.Config {
	// generate random string
	bucketName := "eigenda-proxy-test-" + RandStr(10)
	createS3Bucket(bucketName)

	return s3.Config{
		Bucket:          bucketName,
		Path:            "",
		Endpoint:        minioEndpoint,
		EnableTLS:       false,
		AccessKeySecret: "minioadmin",
		AccessKeyID:     "minioadmin",
		CredentialType:  s3.CredentialTypeStatic,
	}
}

// nolint: funlen
func BuildTestSuiteConfig(testCfg TestConfig) config.AppConfig {
	useMemory := testCfg.Backend == MemstoreBackend
	pk := os.Getenv(privateKeyEnvVar)
	if pk == "" && !useMemory {
		panic("SIGNER_PRIVATE_KEY environment variable not set")
	}

	ethRPC := os.Getenv(ethRPCEnvVar)
	if ethRPC == "" && !useMemory {
		panic("ETHEREUM_RPC environment variable is not set")
	}

	var pollInterval time.Duration
	if useMemory {
		pollInterval = time.Second * 1
	} else {
		pollInterval = time.Minute * 1
	}

	maxBlobLength := testCfg.MaxBlobLength
	if maxBlobLength == "" {
		maxBlobLength = "1mib"
	}
	maxBlobLengthBytes, err := common.ParseBytesAmount(maxBlobLength)
	if err != nil {
		panic(err)
	}

	var disperserHostname string
	var certVerifierAddress string
	var svcManagerAddress string
	var blsOperatorStateRetrieverAddress string
	var eigenDADirectory string
	switch testCfg.Backend {
	case MemstoreBackend:
		break // no need to set these fields for local tests
	case PreprodBackend:
		disperserHostname = disperserPreprodHostname
		certVerifierAddress = preprodCertVerifierAddress
		svcManagerAddress = preprodSvcManagerAddress
		blsOperatorStateRetrieverAddress = preprodBLSOperatorStateRetrieverAddress
		eigenDADirectory = preprodEigenDADirectory
	case TestnetBackend:
		disperserHostname = disperserTestnetHostname
		certVerifierAddress = testnetCertVerifierAddress
		svcManagerAddress = testnetSvcManagerAddress
		blsOperatorStateRetrieverAddress = testnetBLSOperatorStateRetrieverAddress
		eigenDADirectory = testnetEigenDADirectory
	case SepoliaBackend:
		disperserHostname = disperserSepoliaHostname
		certVerifierAddress = sepoliaCertVerifierAddress
		svcManagerAddress = sepoliaSvcManagerAddress
		blsOperatorStateRetrieverAddress = sepoliaBLSOperatorStateRetrieverAddress
		eigenDADirectory = sepoliaEigenDADirectory
	default:
		panic("Unsupported backend")
	}
	payloadClientConfig := clientsv2.PayloadClientConfig{
		PayloadPolynomialForm: codecs.PolynomialFormEval,
		BlobVersion:           0,
	}
	builderConfig := builder.Config{
		StoreConfig: store.Config{
			AsyncPutWorkers:  testCfg.WriteThreadCount,
			BackendsToEnable: testCfg.BackendsToEnable,
			DispersalBackend: testCfg.DispersalBackend,
		},
		ClientConfigV1: common.ClientConfigV1{
			EdaClientCfg: clients.EigenDAClientConfig{
				RPC:                      disperserHostname + ":" + disperserPort,
				StatusQueryTimeout:       time.Minute * 45,
				StatusQueryRetryInterval: pollInterval,
				DisableTLS:               false,
				SignerPrivateKeyHex:      pk,
				EthRpcUrl:                ethRPC,
				SvcManagerAddr:           svcManagerAddress,
			},
			MaxBlobSizeBytes: maxBlobLengthBytes,
			PutTries:         3,
		},
		VerifierConfigV1: verify.Config{
			VerifyCerts:          true,
			RPCURL:               ethRPC,
			SvcManagerAddr:       svcManagerAddress,
			EthConfirmationDepth: 1,
			WaitForFinalization:  false,
			MaxBlobSizeBytes:     maxBlobLengthBytes,
		},
		KzgConfig: kzg.KzgConfig{
			G1Path:          "../../resources/g1.point",
			G2Path:          "../../resources/g2.point",
			G2TrailingPath:  "../../resources/g2.trailing.point",
			CacheDir:        "../../resources/SRSTables",
			SRSOrder:        eigendaflags.SrsOrder,
			SRSNumberToLoad: maxBlobLengthBytes / 32,
			NumWorker:       uint64(runtime.GOMAXPROCS(0)), // #nosec G115
			LoadG2Points:    true,
		},
		MemstoreConfig: memconfig.NewSafeConfig(
			memconfig.Config{
				BlobExpiration:   testCfg.Expiration,
				MaxBlobSizeBytes: maxBlobLengthBytes,
			}),
		MemstoreEnabled: useMemory,
		ClientConfigV2: common.ClientConfigV2{
			DisperserClientCfg: clientsv2.DisperserClientConfig{
				Hostname:          disperserHostname,
				Port:              disperserPort,
				UseSecureGrpcFlag: true,
			},
			PayloadDisperserCfg: payloaddispersal.PayloadDisperserConfig{
				PayloadClientConfig:    payloadClientConfig,
				DisperseBlobTimeout:    5 * time.Minute,
				BlobCompleteTimeout:    5 * time.Minute,
				BlobStatusPollInterval: 1 * time.Second,
				ContractCallTimeout:    5 * time.Second,
			},
			RelayPayloadRetrieverCfg: payloadretrieval.RelayPayloadRetrieverConfig{
				PayloadClientConfig: payloadClientConfig,
				RelayTimeout:        5 * time.Second,
			},
			PutTries:                           3,
			MaxBlobSizeBytes:                   maxBlobLengthBytes,
			EigenDACertVerifierOrRouterAddress: certVerifierAddress,
			BLSOperatorStateRetrieverAddr:      blsOperatorStateRetrieverAddress,
			EigenDAServiceManagerAddr:          svcManagerAddress,
			EigenDADirectory:                   eigenDADirectory,
			RetrieversToEnable:                 testCfg.Retrievers,
		},
	}
	if useMemory {
		builderConfig.ClientConfigV1.EdaClientCfg.SignerPrivateKeyHex =
			"0000000000000000000100000000000000000000000000000000000000000000"
		builderConfig.ClientConfigV1.EdaClientCfg.SvcManagerAddr = "0x00000000069"
		builderConfig.KzgConfig.LoadG2Points = false
		builderConfig.VerifierConfigV1.VerifyCerts = false
	}
	switch {
	case testCfg.UseKeccak256ModeS3:
		builderConfig.S3Config = createS3Config()
	case testCfg.UseS3Caching:
		builderConfig.StoreConfig.CacheTargets = []string{"S3"}
		builderConfig.S3Config = createS3Config()
	case testCfg.UseS3Fallback:
		builderConfig.StoreConfig.FallbackTargets = []string{"S3"}
		builderConfig.S3Config = createS3Config()
	case testCfg.UseRedisCaching:
		builderConfig.StoreConfig.CacheTargets = []string{"redis"}
		builderConfig.RedisConfig = createRedisConfig()
	}
	secretConfig := common.SecretConfigV2{
		SignerPaymentKey: pk,
		EthRPCURL:        ethRPC,
	}
	return config.AppConfig{
		StoreBuilderConfig:  builderConfig,
		SecretConfig:        secretConfig,
		MetricsServerConfig: proxy_metrics.Config{},
		ServerConfig: server.Config{
			Host: host,
			Port: 0,
		},
	}
}
func createS3Bucket(bucketName string) {
	// Initialize minio client object.
	endpoint := minioEndpoint
	accessKeyID := minioAdmin
	secretAccessKey := minioAdmin
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
