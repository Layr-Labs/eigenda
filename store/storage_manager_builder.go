package store

import (
	"context"
	"fmt"
	"math/rand"
	"slices"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/eigenda"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore/memconfig"
	memstore_v2 "github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore/v2"
	eigenda_v2 "github.com/Layr-Labs/eigenda-proxy/store/generated_key/v2"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/redis"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/s3"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api/clients"
	clients_v2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	client_validator "github.com/Layr-Labs/eigenda/api/clients/v2/validator"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	common_eigenda "github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/core/eth"
	core_v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	kzgverifier "github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	geth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// StorageManagerBuilder centralizes dependency initialization.
// It ensures proper typing and avoids redundant logic scattered across functions.
type StorageManagerBuilder struct {
	ctx     context.Context
	log     logging.Logger
	metrics metrics.Metricer

	// configs that are used for both v1 and v2
	managerCfg      Config
	memConfig       *memconfig.SafeConfig
	memstoreEnabled bool
	kzgConfig       kzg.KzgConfig

	// v1 specific configs
	v1ClientCfg   common.ClientConfigV1
	v1VerifierCfg verify.Config

	// v2 specific configs
	v2ClientCfg common.ClientConfigV2
	v2SecretCfg common.SecretConfigV2
}

// NewStorageManagerBuilder creates a builder which knows how to build an IManager
func NewStorageManagerBuilder(
	ctx context.Context,
	log logging.Logger,
	metrics metrics.Metricer,
	managerConfig Config,
	memConfig *memconfig.SafeConfig,
	memstoreEnabled bool,
	kzgConfig kzg.KzgConfig,
	v1ClientCfg common.ClientConfigV1,
	v1VerifierCfg verify.Config,
	v2ClientCfg common.ClientConfigV2,
	v2SecretCfg common.SecretConfigV2,
) *StorageManagerBuilder {
	return &StorageManagerBuilder{
		ctx,
		log,
		metrics,
		managerConfig,
		memConfig,
		memstoreEnabled,
		kzgConfig,
		v1ClientCfg,
		v1VerifierCfg,
		v2ClientCfg,
		v2SecretCfg,
	}
}

// Build builds the storage manager object
func (smb *StorageManagerBuilder) Build(ctx context.Context) (*Manager, error) {
	var err error
	var s3Store *s3.Store
	var redisStore *redis.Store
	var eigenDAV1Store, eigenDAV2Store common.EigenDAStore

	if smb.managerCfg.S3Config.Bucket != "" {
		smb.log.Info("Using S3 storage backend")
		s3Store, err = s3.NewStore(smb.managerCfg.S3Config)
		if err != nil {
			return nil, err
		}
	}

	if smb.managerCfg.RedisConfig.Endpoint != "" {
		smb.log.Info("Using Redis storage backend")
		redisStore, err = redis.NewStore(&smb.managerCfg.RedisConfig)
		if err != nil {
			return nil, err
		}
	}

	v1Enabled := slices.Contains(smb.managerCfg.BackendsToEnable, common.V1EigenDABackend)
	v2Enabled := slices.Contains(smb.managerCfg.BackendsToEnable, common.V2EigenDABackend)

	if smb.managerCfg.DispersalBackend == common.V2EigenDABackend && !v2Enabled {
		return nil, fmt.Errorf("dispersal backend is set to V2, but V2 backend is not enabled")
	} else if smb.managerCfg.DispersalBackend == common.V1EigenDABackend && !v1Enabled {
		return nil, fmt.Errorf("dispersal backend is set to V1, but V1 backend is not enabled")
	}

	var kzgVerifier *kzgverifier.Verifier
	// there are two cases in which we need to construct the kzgVerifier:
	// 1. V1
	// 2. V2, when validator retrieval is enabled
	if v1Enabled || v2Enabled && slices.Contains(smb.v2ClientCfg.RetrieversToEnable, common.ValidatorRetrieverType) {
		// The verifier doesn't support loading trailing g2 points from a separate file. If LoadG2Points is true, and
		// the user is using a slimmed down g2 SRS file, the verifier will encounter an error while trying to load g2
		// points. Since the verifier doesn't actually need g2 points, it's safe to force LoadG2Points to false, to
		// sidestep the issue entirely.
		kzgConfig := smb.kzgConfig
		kzgConfig.LoadG2Points = false

		kzgVerifier, err = kzgverifier.NewVerifier(&kzgConfig, nil)
		if err != nil {
			return nil, fmt.Errorf("new kzg verifier: %w", err)
		}
	}

	if v1Enabled {
		smb.log.Info("Building EigenDA v1 storage backend")
		eigenDAV1Store, err = smb.buildEigenDAV1Backend(ctx, kzgVerifier)
		if err != nil {
			return nil, fmt.Errorf("build v1 backend: %w", err)
		}
	}

	if v2Enabled {
		smb.log.Info("Building EigenDA v2 storage backend")
		eigenDAV2Store, err = smb.buildEigenDAV2Backend(ctx, kzgVerifier)
		if err != nil {
			return nil, fmt.Errorf("build v2 backend: %w", err)
		}
	}

	fallbacks := smb.buildSecondaries(smb.managerCfg.FallbackTargets, s3Store, redisStore)
	caches := smb.buildSecondaries(smb.managerCfg.CacheTargets, s3Store, redisStore)
	secondary := NewSecondaryManager(smb.log, smb.metrics, caches, fallbacks)

	if secondary.Enabled() { // only spin-up go routines if secondary storage is enabled
		smb.log.Info("Starting secondary write loop(s)", "count", smb.managerCfg.AsyncPutWorkers)

		for i := 0; i < smb.managerCfg.AsyncPutWorkers; i++ {
			go secondary.WriteSubscriptionLoop(ctx)
		}
	}

	smb.log.Info(
		"Created storage backends",
		"eigenda_v1", eigenDAV1Store != nil,
		"eigenda_v2", eigenDAV2Store != nil,
		"s3", s3Store != nil,
		"redis", redisStore != nil,
		"read_fallback", len(fallbacks) > 0,
		"caching", len(caches) > 0,
		"async_secondary_writes", (secondary.Enabled() && smb.managerCfg.AsyncPutWorkers > 0),
		"verify_v1_certs", smb.v1VerifierCfg.VerifyCerts,
	)

	return NewManager(eigenDAV1Store, eigenDAV2Store, s3Store, smb.log, secondary, smb.managerCfg.DispersalBackend)
}

// buildSecondaries ... Creates a slice of secondary targets used for either read
// failover or caching
func (smb *StorageManagerBuilder) buildSecondaries(
	targets []string,
	s3Store common.PrecomputedKeyStore,
	redisStore *redis.Store,
) []common.PrecomputedKeyStore {
	stores := make([]common.PrecomputedKeyStore, len(targets))

	for i, target := range targets {
		//nolint:exhaustive // TODO: implement additional secondaries
		switch common.StringToBackendType(target) {
		case common.RedisBackendType:
			if redisStore == nil {
				panic(fmt.Sprintf("Redis backend not configured: %s", target))
			}
			stores[i] = redisStore
		case common.S3BackendType:
			if s3Store == nil {
				panic(fmt.Sprintf("S3 backend not configured: %s", target))
			}
			stores[i] = s3Store

		default:
			panic(fmt.Sprintf("Invalid backend target: %s", target))
		}
	}
	return stores
}

// buildEigenDAV2Backend ... Builds EigenDA V2 storage backend
func (smb *StorageManagerBuilder) buildEigenDAV2Backend(
	ctx context.Context,
	kzgVerifier *kzgverifier.Verifier,
) (common.EigenDAStore, error) {
	// This is a bit of a hack. The kzg config is used by both v1 AND v2, but the `LoadG2Points` field has special
	// requirements. For v1, it must always be false. For v2, it must always be true. Ideally, we would modify
	// the underlying core library to be more flexible, but that is a larger change for another time. As a stopgap, we
	// simply set this value to whatever it needs to be prior to using it.
	config := smb.kzgConfig
	config.LoadG2Points = true

	kzgProver, err := prover.NewProver(&config, nil)
	if err != nil {
		return nil, fmt.Errorf("new KZG prover: %w", err)
	}

	if smb.memstoreEnabled {
		return memstore_v2.New(smb.ctx, smb.log, smb.memConfig, kzgProver.Srs.G1)
	}

	ethClient, err := smb.buildEthClient()
	if err != nil {
		return nil, fmt.Errorf("build eth client: %w", err)
	}

	certVerifierAddressProvider := verification.NewStaticCertVerifierAddressProvider(
		geth_common.HexToAddress(smb.v2ClientCfg.EigenDACertVerifierAddress))

	certVerifier, err := verification.NewCertVerifier(
		smb.log, ethClient, certVerifierAddressProvider)
	if err != nil {
		return nil, fmt.Errorf("new cert verifier: %w", err)
	}

	ethReader, err := smb.buildEthReader(ethClient)
	if err != nil {
		return nil, fmt.Errorf("build eth reader: %w", err)
	}

	var retrievers []clients_v2.PayloadRetriever
	for _, retrieverType := range smb.v2ClientCfg.RetrieversToEnable {
		switch retrieverType {
		case common.RelayRetrieverType:
			smb.log.Info("Initializing relay payload retriever")
			relayPayloadRetriever, err := smb.buildRelayPayloadRetriever(
				ethClient, kzgProver.Srs.G1, ethReader.GetRelayRegistryAddress())
			if err != nil {
				return nil, fmt.Errorf("build relay payload retriever: %w", err)
			}
			retrievers = append(retrievers, relayPayloadRetriever)
		case common.ValidatorRetrieverType:
			smb.log.Info("Initializing validator payload retriever")
			validatorPayloadRetriever, err := smb.buildValidatorPayloadRetriever(
				ethClient, ethReader, kzgVerifier, kzgProver.Srs.G1)
			if err != nil {
				return nil, fmt.Errorf("build validator payload retriever: %w", err)
			}
			retrievers = append(retrievers, validatorPayloadRetriever)
		default:
			return nil, fmt.Errorf("unknown retriever type: %s", retrieverType)
		}
	}

	// Ensure at least one retriever is configured
	if len(retrievers) == 0 {
		return nil, fmt.Errorf("no payload retrievers enabled, please enable at least one retriever type")
	}

	payloadDisperser, err := smb.buildPayloadDisperser(ctx, ethClient, kzgProver, certVerifier)
	if err != nil {
		return nil, fmt.Errorf("build payload disperser: %w", err)
	}

	eigenDAV2Store, err := eigenda_v2.NewStore(
		smb.log,
		smb.v2ClientCfg.PutTries,
		payloadDisperser,
		retrievers,
		certVerifier)
	if err != nil {
		return nil, fmt.Errorf("create v2 store: %w", err)
	}

	return eigenDAV2Store, nil
}

// buildEigenDAV1Backend ... Builds EigenDA V1 storage backend
func (smb *StorageManagerBuilder) buildEigenDAV1Backend(
	ctx context.Context,
	kzgVerifier *kzgverifier.Verifier,
) (common.EigenDAStore, error) {
	verifier, err := verify.NewVerifier(&smb.v1VerifierCfg, kzgVerifier, smb.log)
	if err != nil {
		return nil, fmt.Errorf("new verifier: %w", err)
	}

	if smb.v1VerifierCfg.VerifyCerts {
		smb.log.Info("Certificate verification with Ethereum enabled")
	} else {
		smb.log.Warn("Certificate verification disabled. This can result in invalid EigenDA certificates being accredited.")
	}

	if smb.memstoreEnabled {
		smb.log.Info("Using memstore backend for EigenDA V1")
		return memstore.New(ctx, verifier, smb.log, smb.memConfig)
	}
	// EigenDAV1 backend dependency injection
	var client *clients.EigenDAClient

	client, err = clients.NewEigenDAClient(smb.log, smb.v1ClientCfg.EdaClientCfg)
	if err != nil {
		return nil, err
	}

	storeConfig, err := eigenda.NewStoreConfig(
		smb.v1ClientCfg.MaxBlobSizeBytes,
		smb.v1VerifierCfg.EthConfirmationDepth,
		smb.v1ClientCfg.EdaClientCfg.StatusQueryTimeout,
		smb.v1ClientCfg.PutTries,
	)
	if err != nil {
		return nil, fmt.Errorf("create v1 store config: %w", err)
	}

	return eigenda.NewStore(
		client,
		verifier,
		smb.log,
		storeConfig,
	)
}

func (smb *StorageManagerBuilder) buildEthClient() (common_eigenda.EthClient, error) {
	gethCfg := geth.EthClientConfig{
		RPCURLs: []string{smb.v2SecretCfg.EthRPCURL},
	}

	ethClient, err := geth.NewClient(gethCfg, geth_common.Address{}, 0, smb.log)
	if err != nil {
		return nil, fmt.Errorf("create geth client: %w", err)
	}

	return ethClient, nil
}

func (smb *StorageManagerBuilder) buildRelayPayloadRetriever(
	ethClient common_eigenda.EthClient,
	g1Srs []bn254.G1Affine,
	relayRegistryAddress geth_common.Address,
) (*payloadretrieval.RelayPayloadRetriever, error) {
	relayClient, err := smb.buildRelayClient(ethClient, relayRegistryAddress)
	if err != nil {
		return nil, fmt.Errorf("build relay client: %w", err)
	}

	relayPayloadRetriever, err := payloadretrieval.NewRelayPayloadRetriever(
		smb.log,
		//nolint:gosec // disable G404: this doesn't need to be cryptographically secure
		rand.New(rand.NewSource(time.Now().UnixNano())),
		smb.v2ClientCfg.RelayPayloadRetrieverCfg,
		relayClient,
		g1Srs)
	if err != nil {
		return nil, fmt.Errorf("new relay payload retriever: %w", err)
	}

	return relayPayloadRetriever, nil
}

func (smb *StorageManagerBuilder) buildRelayClient(
	ethClient common_eigenda.EthClient,
	relayRegistryAddress geth_common.Address,
) (relay.RelayClient, error) {
	relayURLProvider, err := relay.NewRelayUrlProvider(ethClient, relayRegistryAddress)
	if err != nil {
		return nil, fmt.Errorf("new relay url provider: %w", err)
	}

	relayCfg := &relay.RelayClientConfig{
		UseSecureGrpcFlag: smb.v2ClientCfg.DisperserClientCfg.UseSecureGrpcFlag,
		// we should never expect a message greater than our allowed max blob size.
		// 10% of max blob size is added for additional safety
		MaxGRPCMessageSize: uint(smb.v2ClientCfg.MaxBlobSizeBytes + (smb.v2ClientCfg.MaxBlobSizeBytes / 10)),
	}

	relayClient, err := relay.NewRelayClient(relayCfg, smb.log, relayURLProvider)
	if err != nil {
		return nil, fmt.Errorf("new relay client: %w", err)
	}

	return relayClient, nil
}

// buildValidatorPayloadRetriever constructs a ValidatorPayloadRetriever for retrieving
// payloads directly from EigenDA validators
func (smb *StorageManagerBuilder) buildValidatorPayloadRetriever(
	ethClient common_eigenda.EthClient,
	ethReader *eth.Reader,
	kzgVerifier *kzgverifier.Verifier,
	g1Srs []bn254.G1Affine,
) (*payloadretrieval.ValidatorPayloadRetriever, error) {
	chainState := eth.NewChainState(ethReader, ethClient)

	retrievalClient := client_validator.NewValidatorClient(
		smb.log,
		ethReader,
		chainState,
		kzgVerifier,
		client_validator.DefaultClientConfig(),
		nil,
	)

	// Create validator payload retriever
	validatorRetriever, err := payloadretrieval.NewValidatorPayloadRetriever(
		smb.log,
		smb.v2ClientCfg.ValidatorPayloadRetrieverCfg,
		retrievalClient,
		g1Srs,
	)
	if err != nil {
		return nil, fmt.Errorf("new validator payload retriever: %w", err)
	}

	return validatorRetriever, nil
}

func (smb *StorageManagerBuilder) buildEthReader(ethClient common_eigenda.EthClient) (*eth.Reader, error) {
	ethReader, err := eth.NewReader(
		smb.log,
		ethClient,
		smb.v2ClientCfg.BLSOperatorStateRetrieverAddr,
		smb.v2ClientCfg.EigenDAServiceManagerAddr,
	)
	if err != nil {
		return nil, fmt.Errorf("new reader: %w", err)
	}

	return ethReader, nil
}

func (smb *StorageManagerBuilder) buildPayloadDisperser(
	ctx context.Context,
	ethClient common_eigenda.EthClient,
	kzgProver *prover.Prover,
	certVerifier *verification.CertVerifier,
) (*payloaddispersal.PayloadDisperser, error) {
	signer, err := smb.buildLocalSigner(ctx, ethClient)
	if err != nil {
		return nil, fmt.Errorf("build local signer: %w", err)
	}

	disperserClient, err := clients_v2.NewDisperserClient(&smb.v2ClientCfg.DisperserClientCfg, signer, kzgProver, nil)
	if err != nil {
		return nil, fmt.Errorf("new disperser client: %w", err)
	}

	payloadDisperser, err := payloaddispersal.NewPayloadDisperser(
		smb.log,
		smb.v2ClientCfg.PayloadDisperserCfg,
		disperserClient,
		certVerifier,
		nil)
	if err != nil {
		return nil, fmt.Errorf("new payload disperser: %w", err)
	}

	return payloadDisperser, nil
}

// buildLocalSigner attempts to check the pending balance of the created signer account. If the check fails, or if the
// balance is determined to be 0, the user is warned with a log. This method doesn't return an error based on this
// check:
// it's possible that a user could want to set up a signer before it's actually ready to be used
//
// TODO: the checks performed in this method could be improved in the future, e.g. by checking payment vault state,
// or by accessing the disperser accountant
func (smb *StorageManagerBuilder) buildLocalSigner(
	ctx context.Context,
	ethClient common_eigenda.EthClient,
) (core_v2.BlobRequestSigner, error) {
	signer, err := auth.NewLocalBlobRequestSigner(smb.v2SecretCfg.SignerPaymentKey)
	if err != nil {
		return nil, fmt.Errorf("new local blob request signer: %w", err)
	}

	accountID := crypto.PubkeyToAddress(signer.PrivateKey.PublicKey)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	pendingBalance, err := ethClient.PendingBalanceAt(ctxWithTimeout, accountID)

	switch {
	case err != nil:
		smb.log.Errorf("get pending balance for accountID %v: %v", accountID, err)
	case pendingBalance == nil:
		smb.log.Errorf(
			"get pending balance for accountID %v didn't return an error, but pending balance is nil", accountID)
	case pendingBalance.Sign() <= 0:
		smb.log.Warnf("pending balance for accountID %v is zero", accountID)
	}

	return signer, nil
}
