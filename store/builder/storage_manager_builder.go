package builder

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"slices"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/eigenda"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/eigenda/verify"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore"
	memstore_v2 "github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore/v2"
	eigenda_v2 "github.com/Layr-Labs/eigenda-proxy/store/generated_key/v2"
	"github.com/Layr-Labs/eigenda-proxy/store/secondary"
	"github.com/Layr-Labs/eigenda-proxy/store/secondary/redis"
	"github.com/Layr-Labs/eigenda-proxy/store/secondary/s3"
	"github.com/Layr-Labs/eigenda/api/clients"
	clients_v2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	client_validator "github.com/Layr-Labs/eigenda/api/clients/v2/validator"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	common_eigenda "github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierRouter"

	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/core/eth"
	core_v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	kzgverifier "github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	geth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

// BuildStoreManager is the main builder for proxy's store.
// It builds all the different store clients, and injects them into
// a new store manager, which it returns when successful.
func BuildStoreManager(
	ctx context.Context,
	log logging.Logger,
	metrics metrics.Metricer,
	config Config,
	secrets common.SecretConfigV2,
) (*store.Manager, error) {
	var err error
	var s3Store *s3.Store
	var redisStore *redis.Store
	var eigenDAV1Store common.EigenDAV1Store
	var eigenDAV2Store common.EigenDAV2Store

	if config.S3Config.Bucket != "" {
		log.Info("Using S3 storage backend")
		s3Store, err = s3.NewStore(config.S3Config)
		if err != nil {
			return nil, err
		}
	}

	if config.RedisConfig.Endpoint != "" {
		log.Info("Using Redis storage backend")
		redisStore, err = redis.NewStore(&config.RedisConfig)
		if err != nil {
			return nil, err
		}
	}

	v1Enabled := slices.Contains(config.StoreConfig.BackendsToEnable, common.V1EigenDABackend)
	v2Enabled := slices.Contains(config.StoreConfig.BackendsToEnable, common.V2EigenDABackend)

	if config.StoreConfig.DispersalBackend == common.V2EigenDABackend && !v2Enabled {
		return nil, fmt.Errorf("dispersal backend is set to V2, but V2 backend is not enabled")
	} else if config.StoreConfig.DispersalBackend == common.V1EigenDABackend && !v1Enabled {
		return nil, fmt.Errorf("dispersal backend is set to V1, but V1 backend is not enabled")
	}

	var kzgVerifier *kzgverifier.Verifier
	// there are two cases in which we need to construct the kzgVerifier:
	// 1. V1
	// 2. V2, when validator retrieval is enabled
	if v1Enabled ||
		v2Enabled && slices.Contains(config.ClientConfigV2.RetrieversToEnable, common.ValidatorRetrieverType) {
		// The verifier doesn't support loading trailing g2 points from a separate file. If LoadG2Points is true, and
		// the user is using a slimmed down g2 SRS file, the verifier will encounter an error while trying to load g2
		// points. Since the verifier doesn't actually need g2 points, it's safe to force LoadG2Points to false, to
		// sidestep the issue entirely.
		kzgConfig := config.KzgConfig
		kzgConfig.LoadG2Points = false

		kzgVerifier, err = kzgverifier.NewVerifier(&kzgConfig, nil)
		if err != nil {
			return nil, fmt.Errorf("new kzg verifier: %w", err)
		}
	}

	if v1Enabled {
		log.Info("Building EigenDA v1 storage backend")
		eigenDAV1Store, err = buildEigenDAV1Backend(ctx, log, config, kzgVerifier)
		if err != nil {
			return nil, fmt.Errorf("build v1 backend: %w", err)
		}
	}

	if v2Enabled {
		log.Info("Building EigenDA v2 storage backend")
		eigenDAV2Store, err = buildEigenDAV2Backend(ctx, log, config, secrets, kzgVerifier)
		if err != nil {
			return nil, fmt.Errorf("build v2 backend: %w", err)
		}
	}

	fallbacks := buildSecondaries(config.StoreConfig.FallbackTargets, s3Store, redisStore)
	caches := buildSecondaries(config.StoreConfig.CacheTargets, s3Store, redisStore)
	secondary := secondary.NewSecondaryManager(log, metrics, caches, fallbacks)

	if secondary.Enabled() { // only spin-up go routines if secondary storage is enabled
		log.Info("Starting secondary write loop(s)", "count", config.StoreConfig.AsyncPutWorkers)

		for i := 0; i < config.StoreConfig.AsyncPutWorkers; i++ {
			go secondary.WriteSubscriptionLoop(ctx)
		}
	}

	log.Info(
		"Created storage backends",
		"eigenda_v1", eigenDAV1Store != nil,
		"eigenda_v2", eigenDAV2Store != nil,
		"s3", s3Store != nil,
		"redis", redisStore != nil,
		"read_fallback", len(fallbacks) > 0,
		"caching", len(caches) > 0,
		"async_secondary_writes", (secondary.Enabled() && config.StoreConfig.AsyncPutWorkers > 0),
		"verify_v1_certs", config.VerifierConfigV1.VerifyCerts,
	)

	return store.NewManager(
		eigenDAV1Store,
		eigenDAV2Store,
		s3Store,
		log,
		secondary,
		config.StoreConfig.DispersalBackend,
	)
}

// buildSecondaries ... Creates a slice of secondary targets used for either read
// failover or caching
func buildSecondaries(
	targets []string,
	s3Store common.SecondaryStore,
	redisStore *redis.Store,
) []common.SecondaryStore {
	stores := make([]common.SecondaryStore, len(targets))

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

// A regexp matching "execution reverted" errors returned from the parent chain RPC.
var executionRevertedRegexp = regexp.MustCompile(`(?i)execution reverted|VM execution error\.?`)

// IsExecutionReverted returns true if the error is an "execution reverted" error
// or if the error is a rpc.Error with ErrorCode 3.
// Taken from
func isExecutionReverted(err error) bool {
	if executionRevertedRegexp.MatchString(err.Error()) {
		return true
	}
	var rpcError rpc.Error
	ok := errors.As(err, &rpcError)
	if ok && rpcError.ErrorCode() == 3 {
		return true
	}
	return false
}

// buildEigenDAV2Backend ... Builds EigenDA V2 storage backend
func buildEigenDAV2Backend(
	ctx context.Context,
	log logging.Logger,
	config Config,
	secrets common.SecretConfigV2,
	kzgVerifier *kzgverifier.Verifier,
) (common.EigenDAV2Store, error) {
	// This is a bit of a hack. The kzg config is used by both v1 AND v2, but the `LoadG2Points` field has special
	// requirements. For v1, it must always be false. For v2, it must always be true. Ideally, we would modify
	// the underlying core library to be more flexible, but that is a larger change for another time. As a stopgap, we
	// simply set this value to whatever it needs to be prior to using it.
	kzgConfig := config.KzgConfig
	kzgConfig.LoadG2Points = true

	kzgProver, err := prover.NewProver(&kzgConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("new KZG prover: %w", err)
	}

	if config.MemstoreEnabled {
		return memstore_v2.New(ctx, log, config.MemstoreConfig, kzgProver.Srs.G1)
	}

	ethClient, err := buildEthClient(ctx, log, secrets, config.ClientConfigV2.EigenDANetwork)
	if err != nil {
		return nil, fmt.Errorf("build eth client: %w", err)
	}

	routerOrImmutableVerifierAddr := geth_common.HexToAddress(config.ClientConfigV2.EigenDACertVerifierOrRouterAddress)
	caller, err := binding.NewContractEigenDACertVerifierRouterCaller(routerOrImmutableVerifierAddr, ethClient)
	if err != nil {
		return nil, fmt.Errorf("new cert verifier router caller: %w", err)
	}

	isRouter := true
	// Check if the router address is actually a router. if method `getCertVerifierAt` fails, it means that the
	// address is not a router, and we should treat it as an immutable cert verifier instead
	_, err = caller.GetCertVerifierAt(&bind.CallOpts{Context: ctx}, 0)
	switch {
	case err != nil && isExecutionReverted(err):
		log.Warnf("EigenDA cert verifier router address was detected to not be a router at address (%s), "+
			"using it as an immutable cert verifier instead", routerOrImmutableVerifierAddr.Hex())
		isRouter = false
	case err != nil:
		return nil, fmt.Errorf("failed to determine whether cert verifier is immutable or "+
			"deployed behind a router at address (%s) : %w", routerOrImmutableVerifierAddr.Hex(), err)
	default:
		log.Infof("EigenDA cert verifier address was detected as an EigenDACertVerifierRouter "+
			"at address (%s), using it as such", routerOrImmutableVerifierAddr.Hex())
	}

	var provider clients_v2.CertVerifierAddressProvider
	if !isRouter {
		provider = verification.NewStaticCertVerifierAddressProvider(
			routerOrImmutableVerifierAddr)
	} else {
		provider, err = verification.BuildRouterAddressProvider(
			routerOrImmutableVerifierAddr,
			ethClient,
			log,
		)

		if err != nil {
			return nil, fmt.Errorf("build router address provider: %w", err)
		}
	}

	ethReader, err := buildEthReader(log, config.ClientConfigV2, ethClient)
	if err != nil {
		return nil, fmt.Errorf("build eth reader: %w", err)
	}
	certVerifier, err := verification.NewCertVerifier(
		log,
		ethClient,
		provider,
	)
	if err != nil {
		return nil, fmt.Errorf("new cert verifier: %w", err)
	}

	var retrievers []clients_v2.PayloadRetriever
	for _, retrieverType := range config.ClientConfigV2.RetrieversToEnable {
		switch retrieverType {
		case common.RelayRetrieverType:
			log.Info("Initializing relay payload retriever")
			relayPayloadRetriever, err := buildRelayPayloadRetriever(
				log, config.ClientConfigV2, ethClient, kzgProver.Srs.G1, ethReader.GetRelayRegistryAddress())
			if err != nil {
				return nil, fmt.Errorf("build relay payload retriever: %w", err)
			}
			retrievers = append(retrievers, relayPayloadRetriever)
		case common.ValidatorRetrieverType:
			log.Info("Initializing validator payload retriever")
			validatorPayloadRetriever, err := buildValidatorPayloadRetriever(
				log, config.ClientConfigV2, ethClient, ethReader, kzgVerifier, kzgProver.Srs.G1)
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

	payloadDisperser, err := buildPayloadDisperser(
		ctx,
		log,
		config.ClientConfigV2,
		secrets,
		ethClient,
		kzgProver,
		certVerifier,
		ethReader,
	)
	if err != nil {
		return nil, fmt.Errorf("build payload disperser: %w", err)
	}

	eigenDAV2Store, err := eigenda_v2.NewStore(
		log,
		config.ClientConfigV2.PutTries,
		config.ClientConfigV2.RBNRecencyWindowSize,
		payloadDisperser,
		retrievers,
		certVerifier,
	)
	if err != nil {
		return nil, fmt.Errorf("create v2 store: %w", err)
	}

	return eigenDAV2Store, nil
}

// buildEigenDAV1Backend ... Builds EigenDA V1 storage backend
func buildEigenDAV1Backend(
	ctx context.Context,
	log logging.Logger,
	config Config,
	kzgVerifier *kzgverifier.Verifier,
) (common.EigenDAV1Store, error) {
	verifier, err := verify.NewVerifier(&config.VerifierConfigV1, kzgVerifier, log)
	if err != nil {
		return nil, fmt.Errorf("new verifier: %w", err)
	}

	if config.VerifierConfigV1.VerifyCerts {
		log.Info("Certificate verification with Ethereum enabled")
	} else {
		log.Warn("Certificate verification disabled. This can result in invalid EigenDA certificates being accredited.")
	}

	if config.MemstoreEnabled {
		log.Info("Using memstore backend for EigenDA V1")
		return memstore.New(ctx, verifier, log, config.MemstoreConfig)
	}
	// EigenDAV1 backend dependency injection
	var client *clients.EigenDAClient

	client, err = clients.NewEigenDAClient(log, config.ClientConfigV1.EdaClientCfg)
	if err != nil {
		return nil, err
	}

	storeConfig, err := eigenda.NewStoreConfig(
		config.ClientConfigV1.MaxBlobSizeBytes,
		config.VerifierConfigV1.EthConfirmationDepth,
		config.ClientConfigV1.EdaClientCfg.StatusQueryTimeout,
		config.ClientConfigV1.PutTries,
	)
	if err != nil {
		return nil, fmt.Errorf("create v1 store config: %w", err)
	}

	return eigenda.NewStore(
		client,
		verifier,
		log,
		storeConfig,
	)
}

func buildEthClient(ctx context.Context, log logging.Logger, secretConfigV2 common.SecretConfigV2,
	expectedNetwork common.EigenDANetwork) (common_eigenda.EthClient, error) {
	gethCfg := geth.EthClientConfig{
		RPCURLs: []string{secretConfigV2.EthRPCURL},
	}

	ethClient, err := geth.NewClient(gethCfg, geth_common.Address{}, 0, log)
	if err != nil {
		return nil, fmt.Errorf("create geth client: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID from ETH RPC: %w", err)
	}

	log.Info("Using chain id: %d", chainID.Uint64())

	// Validate that the chain ID matches the expected network
	if expectedNetwork != "" {
		actualNetwork, err := common.EigenDANetworkFromChainID(chainID.String())
		if err != nil {
			return nil, fmt.Errorf("unknown chain ID %s: %w", chainID.String(), err)
		}
		if actualNetwork != expectedNetwork {
			return nil, fmt.Errorf("network mismatch: expected %s (based on configuration), but ETH RPC"+
				"returned chain ID %s which corresponds to %s",
				expectedNetwork, chainID.String(), actualNetwork)
		}

		log.Info("Detected EigenDA network: %s. Will use for reading network default values if overrides"+
			"aren't provided.", expectedNetwork.String())
	}

	return ethClient, nil
}

func buildRelayPayloadRetriever(
	log logging.Logger,
	clientConfigV2 common.ClientConfigV2,
	ethClient common_eigenda.EthClient,
	g1Srs []bn254.G1Affine,
	relayRegistryAddress geth_common.Address,
) (*payloadretrieval.RelayPayloadRetriever, error) {
	relayClient, err := buildRelayClient(log, clientConfigV2, ethClient, relayRegistryAddress)
	if err != nil {
		return nil, fmt.Errorf("build relay client: %w", err)
	}

	relayPayloadRetriever, err := payloadretrieval.NewRelayPayloadRetriever(
		log,
		//nolint:gosec // disable G404: this doesn't need to be cryptographically secure
		rand.New(rand.NewSource(time.Now().UnixNano())),
		clientConfigV2.RelayPayloadRetrieverCfg,
		relayClient,
		g1Srs)
	if err != nil {
		return nil, fmt.Errorf("new relay payload retriever: %w", err)
	}

	return relayPayloadRetriever, nil
}

func buildRelayClient(
	log logging.Logger,
	clientConfigV2 common.ClientConfigV2,
	ethClient common_eigenda.EthClient,
	relayRegistryAddress geth_common.Address,
) (relay.RelayClient, error) {
	relayURLProvider, err := relay.NewRelayUrlProvider(ethClient, relayRegistryAddress)
	if err != nil {
		return nil, fmt.Errorf("new relay url provider: %w", err)
	}

	relayCfg := &relay.RelayClientConfig{
		UseSecureGrpcFlag: clientConfigV2.DisperserClientCfg.UseSecureGrpcFlag,
		// we should never expect a message greater than our allowed max blob size.
		// 10% of max blob size is added for additional safety
		MaxGRPCMessageSize: uint(clientConfigV2.MaxBlobSizeBytes + (clientConfigV2.MaxBlobSizeBytes / 10)),
	}

	relayClient, err := relay.NewRelayClient(relayCfg, log, relayURLProvider)
	if err != nil {
		return nil, fmt.Errorf("new relay client: %w", err)
	}

	return relayClient, nil
}

// buildValidatorPayloadRetriever constructs a ValidatorPayloadRetriever for retrieving
// payloads directly from EigenDA validators
func buildValidatorPayloadRetriever(
	log logging.Logger,
	clientConfigV2 common.ClientConfigV2,
	ethClient common_eigenda.EthClient,
	ethReader *eth.Reader,
	kzgVerifier *kzgverifier.Verifier,
	g1Srs []bn254.G1Affine,
) (*payloadretrieval.ValidatorPayloadRetriever, error) {
	chainState := eth.NewChainState(ethReader, ethClient)

	retrievalClient := client_validator.NewValidatorClient(
		log,
		ethReader,
		chainState,
		kzgVerifier,
		client_validator.DefaultClientConfig(),
		nil,
	)

	// Create validator payload retriever
	validatorRetriever, err := payloadretrieval.NewValidatorPayloadRetriever(
		log,
		clientConfigV2.ValidatorPayloadRetrieverCfg,
		retrievalClient,
		g1Srs,
	)
	if err != nil {
		return nil, fmt.Errorf("new validator payload retriever: %w", err)
	}

	return validatorRetriever, nil
}

func buildEthReader(log logging.Logger,
	clientConfigV2 common.ClientConfigV2,
	ethClient common_eigenda.EthClient,
) (*eth.Reader, error) {
	ethReader, err := eth.NewReader(
		log,
		ethClient,
		clientConfigV2.BLSOperatorStateRetrieverAddr,
		clientConfigV2.EigenDAServiceManagerAddr,
	)
	if err != nil {
		return nil, fmt.Errorf("new reader: %w", err)
	}

	return ethReader, nil
}

func buildPayloadDisperser(
	ctx context.Context,
	log logging.Logger,
	clientConfigV2 common.ClientConfigV2,
	secrets common.SecretConfigV2,
	ethClient common_eigenda.EthClient,
	kzgProver *prover.Prover,
	certVerifier *verification.CertVerifier,
	ethReader *eth.Reader,
) (*payloaddispersal.PayloadDisperser, error) {
	signer, err := buildLocalSigner(ctx, log, secrets, ethClient)
	if err != nil {
		return nil, fmt.Errorf("build local signer: %w", err)
	}

	disperserClient, err := clients_v2.NewDisperserClient(
		log,
		&clientConfigV2.DisperserClientCfg,
		signer,
		kzgProver,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("new disperser client: %w", err)
	}

	blockNumMonitor, err := verification.NewBlockNumberMonitor(
		log,
		ethClient,
		time.Second*1, // NOTE: this polling interval works for e.g Ethereum but is too slow for L2 chains
		//       which have block times of 2 seconds or less.
	)
	if err != nil {
		return nil, fmt.Errorf("new block number monitor: %w", err)
	}

	certBuilder, err := clients_v2.NewCertBuilder(
		log,
		geth_common.HexToAddress(clientConfigV2.BLSOperatorStateRetrieverAddr),
		ethReader.GetRegistryCoordinatorAddress(),
		ethClient,
	)

	if err != nil {
		return nil, fmt.Errorf("new cert builder: %w", err)
	}

	payloadDisperser, err := payloaddispersal.NewPayloadDisperser(
		log,
		clientConfigV2.PayloadDisperserCfg,
		disperserClient,
		blockNumMonitor,
		certBuilder,
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
func buildLocalSigner(
	ctx context.Context,
	log logging.Logger,
	secrets common.SecretConfigV2,
	ethClient common_eigenda.EthClient,
) (core_v2.BlobRequestSigner, error) {
	signer, err := auth.NewLocalBlobRequestSigner(secrets.SignerPaymentKey)
	if err != nil {
		return nil, fmt.Errorf("new local blob request signer: %w", err)
	}

	accountID := crypto.PubkeyToAddress(signer.PrivateKey.PublicKey)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	pendingBalance, err := ethClient.PendingBalanceAt(ctxWithTimeout, accountID)

	switch {
	case err != nil:
		log.Errorf("get pending balance for accountID %v: %v", accountID, err)
	case pendingBalance == nil:
		log.Errorf(
			"get pending balance for accountID %v didn't return an error, but pending balance is nil", accountID)
	case pendingBalance.Sign() <= 0:
		log.Warnf("pending balance for accountID %v is zero", accountID)
	}

	return signer, nil
}
