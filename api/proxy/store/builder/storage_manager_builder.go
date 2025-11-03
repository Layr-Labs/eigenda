//nolint:funlen // builder functions are expected to be long.
package builder

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"regexp"
	"slices"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	clients_v2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	metrics_v2 "github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	client_validator "github.com/Layr-Labs/eigenda/api/clients/v2/validator"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/store"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/eigenda"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/eigenda/verify"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore"
	memstore_v2 "github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore/v2"
	eigenda_v2 "github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/v2"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/s3"
	common_eigenda "github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierRouter"
	"github.com/prometheus/client_golang/prometheus"

	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
	"github.com/Layr-Labs/eigenda/core/payments/vault"
	core_v2 "github.com/Layr-Labs/eigenda/core/v2"
	kzgverifier "github.com/Layr-Labs/eigenda/encoding/v1/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/committer"
	kzgverifierv2 "github.com/Layr-Labs/eigenda/encoding/v2/kzg/verifier"
	rsv2 "github.com/Layr-Labs/eigenda/encoding/v2/rs"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	geth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

// BuildManagers builds separate cert and keccak managers
func BuildManagers(
	ctx context.Context,
	log logging.Logger,
	metrics metrics.Metricer,
	config Config,
	secrets common.SecretConfigV2,
	registry *prometheus.Registry,
	ethClient common_eigenda.EthClient,
) (*store.EigenDAManager, *store.KeccakManager, error) {
	var err error
	var s3Store *s3.Store
	var eigenDAV1Store common.EigenDAV1Store
	var eigenDAV2Store common.EigenDAV2Store

	if config.S3Config.Bucket != "" {
		log.Info("Using S3 storage backend")
		s3Store, err = s3.NewStore(config.S3Config)
		if err != nil {
			return nil, nil, fmt.Errorf("new S3 store: %w", err)
		}
	}

	v1Enabled := slices.Contains(config.StoreConfig.BackendsToEnable, common.V1EigenDABackend)
	v2Enabled := slices.Contains(config.StoreConfig.BackendsToEnable, common.V2EigenDABackend)

	if config.StoreConfig.DispersalBackend == common.V2EigenDABackend && !v2Enabled {
		return nil, nil, fmt.Errorf("dispersal backend is set to V2, but V2 backend is not enabled")
	} else if config.StoreConfig.DispersalBackend == common.V1EigenDABackend && !v1Enabled {
		return nil, nil, fmt.Errorf("dispersal backend is set to V1, but V1 backend is not enabled")
	}

	if v1Enabled {
		log.Info("Building EigenDA v1 storage backend")
		// The verifier doesn't support loading trailing g2 points from a separate file. If LoadG2Points is true, and
		// the user is using a slimmed down g2 SRS file, the verifier will encounter an error while trying to load g2
		// points. Since the verifier doesn't actually need g2 points, it's safe to force LoadG2Points to false, to
		// sidestep the issue entirely.
		kzgConfig := config.KzgConfig
		kzgConfig.LoadG2Points = false
		kzgVerifier, err := kzgverifier.NewVerifier(&kzgConfig, nil)
		if err != nil {
			return nil, nil, fmt.Errorf("new kzg verifier: %w", err)
		}
		eigenDAV1Store, err = buildEigenDAV1Backend(ctx, log, config, kzgVerifier)
		if err != nil {
			return nil, nil, fmt.Errorf("build v1 backend: %w", err)
		}
	}

	if v2Enabled {
		log.Info("Building EigenDA v2 storage backend")
		// kzgVerifier and encoder are only needed when validator retrieval is enabled
		var kzgVerifier *kzgverifierv2.Verifier
		if slices.Contains(config.ClientConfigV2.RetrieversToEnable, common.ValidatorRetrieverType) {
			kzgConfig := kzgverifierv2.ConfigFromV1KzgConfig(&config.KzgConfig)
			kzgVerifier, err = kzgverifierv2.NewVerifier(kzgConfig)
			if err != nil {
				return nil, nil, fmt.Errorf("new kzg verifier: %w", err)
			}
		}
		encoder, err := rsv2.NewEncoder(log, nil)
		if err != nil {
			return nil, nil, fmt.Errorf("new v2 encoder: %w", err)
		}
		eigenDAV2Store, err = buildEigenDAV2Backend(
			ctx, log, config, secrets, encoder, kzgVerifier, registry, ethClient)
		if err != nil {
			return nil, nil, fmt.Errorf("build v2 backend: %w", err)
		}
	}

	fallbacks := buildSecondaries(config.StoreConfig.FallbackTargets, s3Store)
	caches := buildSecondaries(config.StoreConfig.CacheTargets, s3Store)
	secondary := secondary.NewSecondaryManager(
		log,
		metrics,
		caches,
		fallbacks,
		config.StoreConfig.WriteOnCacheMiss,
		config.StoreConfig.ErrorOnSecondaryInsertFailure,
	)

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
		"read_fallback", len(fallbacks) > 0,
		"caching", len(caches) > 0,
		"async_secondary_writes", (secondary.Enabled() && config.StoreConfig.AsyncPutWorkers > 0),
		"error_on_secondary_insert_failure", config.StoreConfig.ErrorOnSecondaryInsertFailure,
		"verify_v1_certs", config.VerifierConfigV1.VerifyCerts,
	)

	certMgr, err := store.NewEigenDAManager(
		eigenDAV1Store,
		eigenDAV2Store,
		log,
		secondary,
		config.StoreConfig.DispersalBackend,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("new eigenda manager: %w", err)
	}

	keccakMgr, err := store.NewKeccakManager(s3Store, log)
	if err != nil {
		return nil, nil, fmt.Errorf("new keccak manager: %w", err)
	}

	return certMgr, keccakMgr, nil
}

// buildSecondaries ... Creates a slice of secondary targets used for either read
// failover or caching
func buildSecondaries(
	targets []string,
	s3Store common.SecondaryStore,
) []common.SecondaryStore {
	stores := make([]common.SecondaryStore, len(targets))

	for i, target := range targets {
		//nolint:exhaustive // TODO: implement additional secondaries
		switch common.StringToBackendType(target) {
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
	encoder *rsv2.Encoder,
	kzgVerifier *kzgverifierv2.Verifier,
	registry *prometheus.Registry,
	ethClient common_eigenda.EthClient,
) (common.EigenDAV2Store, error) {
	kzgCommitter, err := committer.NewFromConfig(committer.Config{
		G1SRSPath:         config.KzgConfig.G1Path,
		G2SRSPath:         config.KzgConfig.G2Path,
		G2TrailingSRSPath: config.KzgConfig.G2TrailingPath,
		SRSNumberToLoad:   config.KzgConfig.SRSNumberToLoad,
	})
	if err != nil {
		return nil, fmt.Errorf("new kzg committer: %w", err)
	}

	if config.MemstoreEnabled {
		return memstore_v2.New(ctx, log, config.MemstoreConfig, kzgVerifier.G1SRS), nil
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
	certVerifier, err := verification.NewCertVerifier(
		log,
		ethClient,
		provider,
	)
	if err != nil {
		return nil, fmt.Errorf("new cert verifier: %w", err)
	}

	if !isRouter {
		// We call GetCertVersion to ensure that the cert verifier is of a supported version. See
		// https://github.com/Layr-Labs/eigenda/blob/d0a14fa44/contracts/src/integrations/cert/interfaces/IVersionedEigenDACertVerifier.sol#L12
		// https://github.com/Layr-Labs/eigenda/blob/d0a14fa44/contracts/src/integrations/cert/EigenDACertVerifier.sol#L79
		// We pass in block 0 because a static certVerifierAddress provider is used when not using a router,
		// so the block number is not relevant.
		certVersion, err := certVerifier.GetCertVersion(ctx, 0)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to eth-call certVersion(), meaning that you either have network problems with your eth node, or "+
					"%s is not a CertVerifier version >= V3, which is required by this version of proxy: %w",
				routerOrImmutableVerifierAddr.Hex(), err)
		}
		// Note that we also support certV2s, just not V2 CertVerifiers.
		// This is because we transform certV2s into certV3s and verified using the CertVerifierV3 contract.
		// However, the serialization logic, as well as some functions needed during the dispersal path (eg. requiredQuorums),
		// are only compatible/available with CertVerifier V3, hence the requirement here.
		if certVersion != 3 {
			return nil, fmt.Errorf("this version of proxy is only compatible with CertVerifier V3 : cert verifier at address %s is version %d",
				routerOrImmutableVerifierAddr.Hex(), certVersion)
		}
	}

	var eigenDAServiceManagerAddr, operatorStateRetrieverAddr geth_common.Address
	contractDirectory, err := directory.NewContractDirectory(ctx, log, ethClient,
		geth_common.HexToAddress(config.ClientConfigV2.EigenDADirectory))
	if err != nil {
		return nil, fmt.Errorf("new contract directory: %w", err)
	}
	eigenDAServiceManagerAddr, err = contractDirectory.GetContractAddress(ctx, directory.ServiceManager)
	if err != nil {
		return nil, fmt.Errorf("get eigenDAServiceManagerAddr: %w", err)
	}
	operatorStateRetrieverAddr, err = contractDirectory.GetContractAddress(ctx, directory.OperatorStateRetriever)
	if err != nil {
		return nil, fmt.Errorf("get OperatorStateRetriever addr: %w", err)
	}
	registryCoordinator, err := contractDirectory.GetContractAddress(ctx, directory.RegistryCoordinator)
	if err != nil {
		return nil, fmt.Errorf("get registryCoordinator: %w", err)
	}

	retrievalMetrics := metrics_v2.NewRetrievalMetrics(registry)

	var retrievers []clients_v2.PayloadRetriever
	for _, retrieverType := range config.ClientConfigV2.RetrieversToEnable {
		switch retrieverType {
		case common.RelayRetrieverType:
			log.Info("Initializing relay payload retriever")
			relayRegistryAddr, err := contractDirectory.GetContractAddress(ctx, directory.RelayRegistry)
			if err != nil {
				return nil, fmt.Errorf("get relay registry address: %w", err)
			}
			relayPayloadRetriever, err := buildRelayPayloadRetriever(
				log, config.ClientConfigV2, ethClient, kzgVerifier.G1SRS, relayRegistryAddr, retrievalMetrics)
			if err != nil {
				return nil, fmt.Errorf("build relay payload retriever: %w", err)
			}
			retrievers = append(retrievers, relayPayloadRetriever)
		case common.ValidatorRetrieverType:
			log.Info("Initializing validator payload retriever")
			validatorPayloadRetriever, err := buildValidatorPayloadRetriever(
				log, config.ClientConfigV2, ethClient,
				operatorStateRetrieverAddr, eigenDAServiceManagerAddr,
				encoder, kzgVerifier, kzgVerifier.G1SRS, retrievalMetrics)
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

	var payloadDisperser *payloaddispersal.PayloadDisperser

	if secrets.SignerPaymentKey == "" {
		log.Warn("No SignerPaymentKey provided: EigenDA V2 backend configured in read-only mode")
	} else {
		log.Info("SignerPaymentKey available: EigenDA V2 backend configured with write support")
		payloadDisperser, err = buildPayloadDisperser(
			ctx,
			log,
			config.ClientConfigV2,
			secrets,
			ethClient,
			kzgCommitter,
			contractDirectory,
			certVerifier,
			operatorStateRetrieverAddr,
			registryCoordinator,
			registry,
		)
		if err != nil {
			return nil, fmt.Errorf("build payload disperser: %w", err)
		}
	}

	eigenDAV2Store, err := eigenda_v2.NewStore(
		log,
		payloadDisperser,
		config.ClientConfigV2.PutTries,
		certVerifier,
		config.ClientConfigV2.RBNRecencyWindowSize,
		retrievers,
		// PayloadDisperserCfg.ContractCallTimeout is set by the --eigenda.v2.contract-call-timeout flag, the value
		// is not read into any other configs. For simplicity the PayloadDisperserCfg value is reused here.
		config.ClientConfigV2.PayloadDisperserCfg.ContractCallTimeout,
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

func buildRelayPayloadRetriever(
	log logging.Logger,
	clientConfigV2 common.ClientConfigV2,
	ethClient common_eigenda.EthClient,
	g1Srs []bn254.G1Affine,
	relayRegistryAddr geth_common.Address,
	metrics metrics_v2.RetrievalMetricer,
) (*payloadretrieval.RelayPayloadRetriever, error) {
	relayClient, err := buildRelayClient(log, clientConfigV2, ethClient, relayRegistryAddr)
	if err != nil {
		return nil, fmt.Errorf("build relay client: %w", err)
	}

	relayPayloadRetriever, err := payloadretrieval.NewRelayPayloadRetriever(
		log,
		//nolint:gosec // disable G404: this doesn't need to be cryptographically secure
		rand.New(rand.NewSource(time.Now().UnixNano())),
		clientConfigV2.RelayPayloadRetrieverCfg,
		relayClient,
		g1Srs,
		metrics)
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
		ConnectionPoolSize: clientConfigV2.RelayConnectionPoolSize,
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
	operatorStateRetrieverAddr geth_common.Address,
	eigenDAServiceManagerAddr geth_common.Address,
	encoder *rsv2.Encoder,
	kzgVerifier *kzgverifierv2.Verifier,
	g1Srs []bn254.G1Affine,
	metrics metrics_v2.RetrievalMetricer,
) (*payloadretrieval.ValidatorPayloadRetriever, error) {
	ethReader, err := eth.NewReader(
		log,
		ethClient,
		operatorStateRetrieverAddr.String(),
		eigenDAServiceManagerAddr.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("new reader: %w", err)
	}
	chainState := eth.NewChainState(ethReader, ethClient)

	retrievalClient := client_validator.NewValidatorClient(
		log,
		ethReader,
		chainState,
		encoder,
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
		metrics,
	)
	if err != nil {
		return nil, fmt.Errorf("new validator payload retriever: %w", err)
	}

	return validatorRetriever, nil
}

func buildPayloadDisperser(
	ctx context.Context,
	log logging.Logger,
	clientConfigV2 common.ClientConfigV2,
	secrets common.SecretConfigV2,
	ethClient common_eigenda.EthClient,
	kzgCommitter *committer.Committer,
	contractDirectory *directory.ContractDirectory,
	certVerifier *verification.CertVerifier,
	operatorStateRetrieverAddr geth_common.Address,
	registryCoordinatorAddr geth_common.Address,
	registry *prometheus.Registry,
) (*payloaddispersal.PayloadDisperser, error) {
	signer, err := buildLocalSigner(ctx, log, secrets, ethClient)
	if err != nil {
		return nil, fmt.Errorf("build local signer: %w", err)
	}

	accountId, err := signer.GetAccountID()
	if err != nil {
		return nil, fmt.Errorf("error getting account ID: %w", err)
	}

	log.Infof("Using account ID %s", accountId.Hex())

	accountantMetrics := metrics_v2.NewAccountantMetrics(registry)
	dispersalMetrics := metrics_v2.NewDispersalMetrics(registry)

	var accountant *clients_v2.Accountant
	// The legacy `Accountant` is only initialized if using legacy payments.
	//
	// There isn't an `else` statement here, because `ClientLedger` (responsible for the new payment system)
	// construction is handled below by the `buildClientLedger` helper function. The `ClientLedger` cannot be built
	// here in the same place as the `Accountant` because it requires the `disperserClient` be already built, and the
	// `Accountant`, if being used, is a part of the `disperserClient`
	if clientConfigV2.ClientLedgerMode == clientledger.ClientLedgerModeLegacy {
		// The accountant is populated lazily by disperserClient.PopulateAccountant
		accountant = clients_v2.NewUnpopulatedAccountant(accountId, accountantMetrics)
	}

	disperserClient, err := clients_v2.NewDisperserClient(
		log,
		&clientConfigV2.DisperserClientCfg,
		signer,
		kzgCommitter,
		accountant,
		dispersalMetrics,
	)
	if err != nil {
		return nil, fmt.Errorf("new disperser client: %w", err)
	}

	clientLedger, err := buildClientLedger(
		ctx,
		log,
		clientConfigV2,
		ethClient,
		accountId,
		contractDirectory,
		accountantMetrics,
		time.Now,
		disperserClient,
	)
	if err != nil {
		return nil, fmt.Errorf("build client ledger: %w", err)
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
		log, operatorStateRetrieverAddr, registryCoordinatorAddr, ethClient)
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
		clientLedger,
		registry)
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

// buildReservationLedger creates a reservation ledger for a given account
func buildReservationLedger(
	ctx context.Context,
	paymentVault payments.PaymentVault,
	accountID geth_common.Address,
	now time.Time,
	minNumSymbols uint32,
) (*reservation.ReservationLedger, error) {
	reservationData, err := paymentVault.GetReservation(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get reservation: %w", err)
	}
	if reservationData == nil {
		return nil, fmt.Errorf("no reservation found for account %s", accountID.Hex())
	}

	clientReservation, err := reservation.NewReservation(
		reservationData.SymbolsPerSecond,
		time.Unix(int64(reservationData.StartTimestamp), 0),
		time.Unix(int64(reservationData.EndTimestamp), 0),
		reservationData.QuorumNumbers,
	)
	if err != nil {
		return nil, fmt.Errorf("new reservation: %w", err)
	}

	reservationConfig, err := reservation.NewReservationLedgerConfig(
		*clientReservation,
		minNumSymbols,
		// start full since reservation usage isn't persisted: assume the worst case (heavy usage before startup)
		true,
		// this is a parameter for flexibility, but there aren't plans to operate with anything other than this value
		ratelimit.OverfillOncePermitted,
		// TODO(litt3): once the checkpointed onchain config registry is ready, that should be used
		// instead of hardcoding. At that point, this field will be removed from the config struct
		// entirely, and the value will be fetched dynamically at runtime.
		30*time.Second,
	)
	if err != nil {
		return nil, fmt.Errorf("new reservation ledger config: %w", err)
	}

	reservationLedger, err := reservation.NewReservationLedger(*reservationConfig, now)
	if err != nil {
		return nil, fmt.Errorf("new reservation ledger: %w", err)
	}

	return reservationLedger, nil
}

// buildOnDemandLedger creates an on-demand ledger for a given account
func buildOnDemandLedger(
	ctx context.Context,
	paymentVault payments.PaymentVault,
	accountID geth_common.Address,
	minNumSymbols uint32,
	disperserClient *clients_v2.DisperserClient,
) (*ondemand.OnDemandLedger, error) {
	pricePerSymbol, err := paymentVault.GetPricePerSymbol(ctx)
	if err != nil {
		return nil, fmt.Errorf("get price per symbol: %w", err)
	}

	totalDeposits, err := paymentVault.GetTotalDeposit(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get total deposit from vault: %w", err)
	}

	paymentState, err := disperserClient.GetPaymentState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get payment state from disperser: %w", err)
	}

	var cumulativePayment *big.Int
	if paymentState.GetCumulativePayment() == nil {
		cumulativePayment = big.NewInt(0)
	} else {
		cumulativePayment = new(big.Int).SetBytes(paymentState.GetCumulativePayment())
	}

	onDemandLedger, err := ondemand.OnDemandLedgerFromValue(
		totalDeposits,
		new(big.Int).SetUint64(pricePerSymbol),
		minNumSymbols,
		cumulativePayment,
	)
	if err != nil {
		return nil, fmt.Errorf("new on-demand ledger: %w", err)
	}

	return onDemandLedger, nil
}

// buildClientLedger creates a ClientLedger for managing payment state
// Returns nil for legacy mode
func buildClientLedger(
	ctx context.Context,
	log logging.Logger,
	config common.ClientConfigV2,
	ethClient common_eigenda.EthClient,
	accountID geth_common.Address,
	contractDirectory *directory.ContractDirectory,
	accountantMetrics metrics_v2.AccountantMetricer,
	getNow func() time.Time,
	disperserClient *clients_v2.DisperserClient,
) (*clientledger.ClientLedger, error) {
	if config.ClientLedgerMode == clientledger.ClientLedgerModeLegacy {
		return nil, nil
	}
	paymentVaultAddr, err := contractDirectory.GetContractAddress(ctx, directory.PaymentVault)
	if err != nil {
		return nil, fmt.Errorf("get PaymentVault address: %w", err)
	}

	paymentVault, err := vault.NewPaymentVault(log, ethClient, paymentVaultAddr)
	if err != nil {
		return nil, fmt.Errorf("new payment vault: %w", err)
	}

	now := getNow()

	minNumSymbols, err := paymentVault.GetMinNumSymbols(ctx)
	if err != nil {
		return nil, fmt.Errorf("get min num symbols: %w", err)
	}

	var reservationLedger *reservation.ReservationLedger
	var onDemandLedger *ondemand.OnDemandLedger
	switch config.ClientLedgerMode {
	case clientledger.ClientLedgerModeLegacy:
		panic("impossible case- this is checked at the start of the method")
	case clientledger.ClientLedgerModeReservationOnly:
		reservationLedger, err = buildReservationLedger(ctx, paymentVault, accountID, now, minNumSymbols)
		if err != nil {
			return nil, fmt.Errorf("build reservation ledger: %w", err)
		}
	case clientledger.ClientLedgerModeOnDemandOnly:
		onDemandLedger, err = buildOnDemandLedger(ctx, paymentVault, accountID, minNumSymbols, disperserClient)
		if err != nil {
			return nil, fmt.Errorf("build on-demand ledger: %w", err)
		}

	case clientledger.ClientLedgerModeReservationAndOnDemand:
		reservationLedger, err = buildReservationLedger(ctx, paymentVault, accountID, now, minNumSymbols)
		if err != nil {
			return nil, fmt.Errorf("build reservation ledger: %w", err)
		}
		onDemandLedger, err = buildOnDemandLedger(ctx, paymentVault, accountID, minNumSymbols, disperserClient)
		if err != nil {
			return nil, fmt.Errorf("build on-demand ledger: %w", err)
		}

	default:
		return nil, fmt.Errorf("unexpected client ledger mode: %s", config.ClientLedgerMode)
	}

	ledger := clientledger.NewClientLedger(
		ctx,
		log,
		accountantMetrics,
		accountID,
		config.ClientLedgerMode,
		reservationLedger,
		onDemandLedger,
		getNow,
		paymentVault,
		config.VaultMonitorInterval,
	)

	return ledger, nil
}
