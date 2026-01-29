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
	"github.com/Layr-Labs/eigenda/api/clients/v2/dispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/config"
	enablement "github.com/Layr-Labs/eigenda/api/proxy/config/enablement"
	proxy_metrics "github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/arbitrum_altda"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/rest"
	"github.com/Layr-Labs/eigenda/api/proxy/store"
	"github.com/Layr-Labs/eigenda/api/proxy/store/builder"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/eigenda/verify"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/s3"
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v1/kzg"
	integration "github.com/Layr-Labs/eigenda/inabox/tests"
	"github.com/ethereum/go-ethereum/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	miniotc "github.com/testcontainers/testcontainers-go/modules/minio"
)

const (
	minioAdmin       = "minioadmin"
	backendEnvVar    = "BACKEND"
	privateKeyEnvVar = "SIGNER_PRIVATE_KEY"
	EthRPCEnvVar     = "ETHEREUM_RPC"
	transport        = "http"
	host             = "127.0.0.1"

	// CertVerifier and SvcManager addresses are still specified by hand for V1.
	// Probably not worth the effort to force use of EigenDADirectory for V1.
	disperserSepoliaHostname   = "disperser-testnet-sepolia.eigenda.xyz"
	sepoliaEigenDADirectory    = "0x9620dC4B3564198554e4D2b06dEFB7A369D90257"
	sepoliaCertVerifierAddress = "0x58D2B844a894f00b7E6F9F492b9F43aD54Cd4429"
	sepoliaSvcManagerAddress   = "0x3a5acf46ba6890B8536420F4900AC9BC45Df4764"

	disperserHoodiTestnetHostname   = "disperser-hoodi.eigenda.xyz"
	hoodiTestnetEigenDADirectory    = "0x5a44e56e88abcf610c68340c6814ae7f5c4369fd"
	hoodiTestnetCertVerifierAddress = "0xD82d14F1c6d1403E95Cd9EC40CBb6463E27C1c5F"
	hoodiTestnetSvcManagerAddress   = "0x3FF2204A567C15dC3731140B95362ABb4b17d8ED"

	disperserHoodiPreprodHostname   = "disperser-v2-preprod-hoodi.eigenda.xyz"
	hoodiPreprodEigenDADirectory    = "0xbFa1b820bb302925a3eb98C8836a95361FB75b87"
	hoodiPreprodCertVerifierAddress = "0xb64101890d15499790d665f9863ede1278ce553d"
	hoodiPreprodSvcManagerAddress   = "0x9F3A67f1b56d0B21115A54356c02B2d77f39EA8a"
)

var (
	disperserPort = "443"
	// set by startMinioContainer
	minioEndpoint = ""
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

type Backend int

const (
	SepoliaBackend Backend = iota + 1
	MemstoreBackend
	HoodiTestnetBackend
	HoodiPreprodBackend
	InaboxBackend
)

func (b Backend) SupportsEigenDAV1() bool {
	switch b {
	// technically HoodiTestnet supports V1 but there's 0 rollup usage
	case HoodiTestnetBackend, HoodiPreprodBackend, InaboxBackend:
		return false

	case SepoliaBackend, MemstoreBackend:
		return true

	default:
		panic("unknown backend type can't be inferred")
	}
}

// ParseBackend converts a string to a Backend enum (case insensitive)
func ParseBackend(inputString string) (Backend, error) {
	switch strings.ToLower(inputString) {
	case "sepolia":
		return SepoliaBackend, nil
	case "memstore":
		return MemstoreBackend, nil
	case "hoodi-testnet":
		return HoodiTestnetBackend, nil
	case "hoodi-preprod":
		return HoodiPreprodBackend, nil
	case "inabox":
		return InaboxBackend, nil

	default:
		return 0, fmt.Errorf("invalid backend: %s", inputString)
	}
}

func GetBackend() Backend {
	backend, err := ParseBackend(os.Getenv(backendEnvVar))
	if err != nil {
		panic(fmt.Sprintf("BACKEND must be = memstore|hoodi-testnet|hoodi-preprod|sepolia|inabox. parse backend error: %v", err))
	}
	return backend
}

type TestConfig struct {
	EnabledRestAPIs  *enablement.RestApisEnabled
	BackendsToEnable []common.EigenDABackend
	DispersalBackend common.EigenDABackend
	Backend          Backend
	Retrievers       []common.RetrieverType
	Expiration       time.Duration
	MaxBlobLength    string
	WriteThreadCount int
	WriteOnCacheMiss bool
	// at most one of the below options should be true
	UseKeccak256ModeS3            bool
	UseS3Caching                  bool
	UseS3Fallback                 bool
	ErrorOnSecondaryInsertFailure bool

	ClientLedgerMode     clientledger.ClientLedgerMode
	VaultMonitorInterval time.Duration
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
		EnabledRestAPIs: &enablement.RestApisEnabled{
			Admin:               false,
			OpGenericCommitment: true,
			OpKeccakCommitment:  true,
			StandardCommitment:  true,
		},
		BackendsToEnable:              backendsToEnable,
		DispersalBackend:              dispersalBackend,
		Backend:                       backend,
		Retrievers:                    []common.RetrieverType{common.RelayRetrieverType, common.ValidatorRetrieverType},
		Expiration:                    14 * 24 * time.Hour,
		UseKeccak256ModeS3:            false,
		UseS3Caching:                  false,
		UseS3Fallback:                 false,
		WriteThreadCount:              0,
		WriteOnCacheMiss:              false,
		ErrorOnSecondaryInsertFailure: false,
		ClientLedgerMode:              clientledger.ClientLedgerModeReservationOnly,
		VaultMonitorInterval:          30 * time.Second,
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
	useInabox := testCfg.Backend == InaboxBackend
	pk := os.Getenv(privateKeyEnvVar)
	ethRPC := ""

	if useInabox {
		// inabox tests always use v2 backend
		testCfg.DispersalBackend = common.V2EigenDABackend
		testCfg.BackendsToEnable = []common.EigenDABackend{common.V2EigenDABackend}

		// use inabox ethRPC
		ethRPC = "http://localhost:8545"

		// use inabox default private key
		pk = integration.GetDefaultTestPayloadDisperserConfig().PrivateKey

		// TODO(iquidus): initialize on-demand balance for on-demand test
	} else {
		ethRPC = os.Getenv(EthRPCEnvVar)
		if ethRPC == "" && !useMemory {
			panic("ETHEREUM_RPC environment variable is not set")
		}
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
	var eigenDADirectory string
	switch testCfg.Backend {
	case MemstoreBackend:
		break // no need to set these fields for local tests
	case SepoliaBackend:
		disperserHostname = disperserSepoliaHostname
		certVerifierAddress = sepoliaCertVerifierAddress
		svcManagerAddress = sepoliaSvcManagerAddress
		eigenDADirectory = sepoliaEigenDADirectory

	case HoodiTestnetBackend:
		disperserHostname = disperserHoodiTestnetHostname
		certVerifierAddress = hoodiTestnetCertVerifierAddress
		svcManagerAddress = hoodiTestnetSvcManagerAddress
		eigenDADirectory = hoodiTestnetEigenDADirectory

	case HoodiPreprodBackend:
		disperserHostname = disperserHoodiPreprodHostname
		certVerifierAddress = hoodiPreprodCertVerifierAddress
		svcManagerAddress = hoodiPreprodSvcManagerAddress
		eigenDADirectory = hoodiPreprodEigenDADirectory
	case InaboxBackend:
		// TODO(iquidus): set these values dynamically when inabox backend is ready
		disperserHostname = "localhost"
		disperserPort = "32005"
		certVerifierAddress = "0x99bbA657f2BbC93c02D617f8bA121cB8Fc104Acf"
		svcManagerAddress = ""
		eigenDADirectory = "0x1613beB3B2C4f22Ee086B2b38C1476A3cE7f78E8"

	default:
		panic("Unsupported backend")
	}
	payloadClientConfig := clientsv2.PayloadClientConfig{
		PayloadPolynomialForm: codecs.PolynomialFormEval,
		BlobVersion:           0,
	}

	builderConfig := builder.Config{
		StoreConfig: store.Config{
			AsyncPutWorkers:               testCfg.WriteThreadCount,
			BackendsToEnable:              testCfg.BackendsToEnable,
			DispersalBackend:              testCfg.DispersalBackend,
			WriteOnCacheMiss:              testCfg.WriteOnCacheMiss,
			ErrorOnSecondaryInsertFailure: testCfg.ErrorOnSecondaryInsertFailure,
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
			SRSOrder:        encoding.SRSOrder,
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
			DisperserClientCfg: dispersal.DisperserClientConfig{
				GrpcUri:           fmt.Sprintf("%s:%s", disperserHostname, disperserPort),
				UseSecureGrpcFlag: !useInabox,
				DisperserID:       0,
				ChainID:           nil, // Will be populated after eth client is created
			},
			PayloadDisperserCfg: dispersal.PayloadDisperserConfig{
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
			EigenDADirectory:                   eigenDADirectory,
			RetrieversToEnable:                 testCfg.Retrievers,
			ClientLedgerMode:                   testCfg.ClientLedgerMode,
			VaultMonitorInterval:               testCfg.VaultMonitorInterval,
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
	}

	secretConfig := common.SecretConfigV2{
		SignerPaymentKey: pk,
		EthRPCURL:        ethRPC,
	}

	return config.AppConfig{
		StoreBuilderConfig: builderConfig,
		SecretConfig:       secretConfig,
		EnabledServersConfig: &enablement.EnabledServersConfig{
			Metric:        false,
			ArbCustomDA:   false,
			RestAPIConfig: *testCfg.EnabledRestAPIs,
		},
		MetricsSvrConfig: proxy_metrics.Config{},
		RestSvrCfg: rest.Config{
			Host:        host,
			Port:        0,
			APIsEnabled: testCfg.EnabledRestAPIs,
		},
		ArbCustomDASvrCfg: arbitrum_altda.Config{
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
