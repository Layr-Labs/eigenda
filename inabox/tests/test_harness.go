package integration

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	routerbindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierRouter"
	verifierv1bindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV1"
	paymentvaultbindings "github.com/Layr-Labs/eigenda/contracts/bindings/PaymentVault"
	"github.com/Layr-Labs/eigenda/core"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
	"github.com/Layr-Labs/eigenda/core/payments/vault"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/testcontainers/testcontainers-go"
)

// InfrastructureHarness contains the shared infrastructure components
// that are global across all tests (external dependencies)
type InfrastructureHarness struct {
	// Shared docker network. Currently the only users of this network are the anvil chain and the graph node.
	SharedNetwork *testcontainers.DockerNetwork

	// Chain related components
	ChainHarness ChainHarness

	// Operator related components
	OperatorHarness OperatorHarness

	// EigenDA components (includes relays)
	DisperserHarness DisperserHarness

	// Proxy
	// TODO: Add harness when we need it

	// Legacy deployment configuration
	TestConfig     *deploy.Config
	TemplateName   string
	TestName       string
	LocalStackPort string

	// Logger for the infrastructure components
	Logger logging.Logger

	// Context for managing infrastructure lifecycle
	Ctx    context.Context
	Cancel context.CancelFunc
}

// TestHarness contains all the components that should be created fresh for each test
type TestHarness struct {
	// Ethereum clients
	EthClient common.EthClient
	RPCClient common.RPCEthClient

	// Verifiers and builders
	CertBuilder                     *clientsv2.CertBuilder
	RouterCertVerifier              *verification.CertVerifier
	StaticCertVerifier              *verification.CertVerifier
	EigenDACertVerifierRouter       *routerbindings.ContractEigenDACertVerifierRouterTransactor
	EigenDACertVerifierRouterCaller *routerbindings.ContractEigenDACertVerifierRouterCaller
	EigenDACertVerifierV1           *verifierv1bindings.ContractEigenDACertVerifierV1

	// Retrieval clients
	RetrievalClient            clients.RetrievalClient
	RelayRetrievalClientV2     *payloadretrieval.RelayPayloadRetriever
	ValidatorRetrievalClientV2 *payloadretrieval.ValidatorPayloadRetriever
	// Tests can use this default payload disperser directly, or create custom payload dispersers via
	// CreatePayloadDisperser().
	PayloadDisperser *payloaddispersal.PayloadDisperser

	// Core components
	ChainReader       core.Reader
	ContractDirectory *directory.ContractDirectory

	// PaymentVault interaction
	PaymentVaultTransactor *paymentvaultbindings.ContractPaymentVaultTransactor

	// Transaction options - specific to test
	DeployerTransactorOpts *bind.TransactOpts
	// Access to the TransactOpts must be synchronized if transactions from the same account are submitted
	// in parallel. The internal logic for determining nonce isn't threadsafe.
	deployerTransactOptsLock sync.Mutex

	// Test-specific configuration
	NumConfirmations int
	NumRetries       int

	// Chain ID for this test context
	ChainID *big.Int
}

// Cleanup releases resources held by the TestHarness
func (tc *TestHarness) Cleanup() {
	// Clean up any test-specific resources if needed
	// Most will be garbage collected, but connections will be closed when EthClient is garbage collected
}

// Provides thread-safe access to the deployer TransactOpts.
//
// Returns the TransactOpts and an unlock function that MUST be called when done.
//
// TODO(litt3): This is a bit of a hack. The returned struct doesn't have a populated nonce field: the nonce is
// populated by the ethereum client iff the nonce within TransactOpts is nil. An alternate strategy to the one used here
// would be to keep track of nonce internally instead of relying on the eth client, thus hiding any synchronization
// logic from the user of the utility. But I struggled to get that working, and decided to go with what worked for now.
// A future task could be to improve the user experience by hiding the sync logic.
func (tc *TestHarness) GetDeployerTransactOpts() (*bind.TransactOpts, func()) {
	tc.deployerTransactOptsLock.Lock()
	return tc.DeployerTransactorOpts, func() {
		tc.deployerTransactOptsLock.Unlock()
	}
}

// Updates the reservation for the specified account on the PaymentVault contract
func (tc *TestHarness) UpdateReservationOnChain(
	t *testing.T,
	accountID gethcommon.Address,
	reservation *reservation.Reservation,
) error {
	quorumNumbers := reservation.GetQuorumNumbers()
	quorumSplits := calculateQuorumSplits(len(quorumNumbers))

	newReservation := paymentvaultbindings.IPaymentVaultReservation{
		SymbolsPerSecond: reservation.GetSymbolsPerSecond(),
		StartTimestamp:   uint64(reservation.GetStartTime().Unix()),
		EndTimestamp:     uint64(reservation.GetEndTime().Unix()),
		QuorumNumbers:    quorumNumbers,
		QuorumSplits:     quorumSplits,
	}

	opts, unlock := tc.GetDeployerTransactOpts()
	defer unlock()

	tx, err := tc.PaymentVaultTransactor.SetReservation(
		opts,
		accountID,
		newReservation,
	)
	if err != nil {
		return fmt.Errorf("set reservation: %w", err)
	}

	receipt, err := bind.WaitMined(t.Context(), tc.EthClient, tx)
	if err != nil {
		return fmt.Errorf("wait mined: %w", err)
	}

	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed")
	}

	return nil
}

// calculateQuorumSplits creates equal percentage splits for all quorums
// The splits will sum to 100, with any remainder going to the first quorum
func calculateQuorumSplits(numQuorums int) []byte {
	quorumSplits := make([]byte, numQuorums)
	if numQuorums > 0 {
		splitValue := byte(100 / numQuorums)
		remainder := byte(100 % numQuorums)
		for i := range quorumSplits {
			quorumSplits[i] = splitValue
			if i == 0 {
				quorumSplits[i] += remainder // Add remainder to first quorum
			}
		}
	}
	return quorumSplits
}

// Creates a new PayloadDisperser and configures the client according to the provided configuration.
func (tc *TestHarness) CreatePayloadDisperser(
	ctx context.Context,
	logger logging.Logger,
	config TestPayloadDisperserConfig,
) (*payloaddispersal.PayloadDisperser, error) {
	blockMonitor, err := verification.NewBlockNumberMonitor(logger, tc.EthClient, time.Second*1)
	if err != nil {
		return nil, fmt.Errorf("create block number monitor: %w", err)
	}

	if config.PrivateKey == "" {
		return nil, fmt.Errorf("private key must be provided")
	}

	signer, err := auth.NewLocalBlobRequestSigner(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("create blob request signer: %w", err)
	}

	disperserClientConfig := &clientsv2.DisperserClientConfig{
		Hostname: "localhost",
		Port:     strconv.Itoa(int(config.APIServerPort)),
	}

	accountId, err := signer.GetAccountID()
	if err != nil {
		return nil, fmt.Errorf("error getting account ID: %w", err)
	}

	accountant := clientsv2.NewAccountant(
		accountId,
		nil,
		nil,
		0,
		0,
		0,
		0,
		metrics.NoopAccountantMetrics,
	)

	disperserClient, err := clientsv2.NewDisperserClient(
		logger,
		disperserClientConfig,
		signer,
		nil, // no prover so will query disperser for generating commitments
		accountant,
		metrics.NoopDispersalMetrics,
	)
	if err != nil {
		return nil, fmt.Errorf("create disperser client: %w", err)
	}

	payloadDisperserConfig := payloaddispersal.PayloadDisperserConfig{
		PayloadClientConfig:    *clientsv2.GetDefaultPayloadClientConfig(),
		DisperseBlobTimeout:    2 * time.Minute,
		BlobCompleteTimeout:    2 * time.Minute,
		BlobStatusPollInterval: 1 * time.Second,
		ContractCallTimeout:    5 * time.Second,
	}

	// Create ClientLedger based on configured mode
	var clientLedger *clientledger.ClientLedger
	if config.ClientLedgerMode != clientledger.ClientLedgerModeLegacy {
		paymentVaultAddr, err := tc.ContractDirectory.GetContractAddress(ctx, directory.PaymentVault)
		if err != nil {
			return nil, fmt.Errorf("get PaymentVault address: %w", err)
		}

		clientLedger, err = buildClientLedger(
			ctx,
			logger,
			tc.EthClient,
			paymentVaultAddr,
			accountId,
			config.ClientLedgerMode,
			disperserClient,
		)
		if err != nil {
			return nil, fmt.Errorf("build client ledger: %w", err)
		}
	}

	payloadDisperser, err := payloaddispersal.NewPayloadDisperser(
		logger,
		payloadDisperserConfig,
		disperserClient,
		blockMonitor,
		tc.CertBuilder,
		tc.RouterCertVerifier,
		clientLedger,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create payload disperser: %w", err)
	}

	return payloadDisperser, nil
}

func buildClientLedger(
	ctx context.Context,
	logger logging.Logger,
	ethClient common.EthClient,
	paymentVaultAddr gethcommon.Address,
	accountID gethcommon.Address,
	mode clientledger.ClientLedgerMode,
	disperserClient *clientsv2.DisperserClient,
) (*clientledger.ClientLedger, error) {
	paymentVault, err := vault.NewPaymentVault(logger, ethClient, paymentVaultAddr)
	if err != nil {
		return nil, fmt.Errorf("new payment vault: %w", err)
	}

	minNumSymbols, err := paymentVault.GetMinNumSymbols(ctx)
	if err != nil {
		return nil, fmt.Errorf("get min num symbols: %w", err)
	}

	var reservationLedger *reservation.ReservationLedger
	var onDemandLedger *ondemand.OnDemandLedger

	// Build reservation ledger if needed
	needsReservation := mode == clientledger.ClientLedgerModeReservationOnly ||
		mode == clientledger.ClientLedgerModeReservationAndOnDemand
	if needsReservation {
		reservationLedger, err = buildReservationLedger(ctx, paymentVault, accountID, minNumSymbols)
		if err != nil {
			return nil, fmt.Errorf("build reservation ledger: %w", err)
		}
	}

	// Build on-demand ledger if needed
	needsOnDemand := mode == clientledger.ClientLedgerModeOnDemandOnly ||
		mode == clientledger.ClientLedgerModeReservationAndOnDemand
	if needsOnDemand {
		onDemandLedger, err = buildOnDemandLedger(ctx, paymentVault, accountID, minNumSymbols, disperserClient)
		if err != nil {
			return nil, fmt.Errorf("build on-demand ledger: %w", err)
		}
	}

	ledger := clientledger.NewClientLedger(
		ctx,
		logger,
		metrics.NoopAccountantMetrics,
		accountID,
		mode,
		reservationLedger,
		onDemandLedger,
		time.Now,
		paymentVault,
		1*time.Second, // update interval for vault monitoring
	)

	return ledger, nil
}

func buildReservationLedger(
	ctx context.Context,
	paymentVault payments.PaymentVault,
	accountID gethcommon.Address,
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
		true,
		ratelimit.OverfillOncePermitted,
		10*time.Second,
	)
	if err != nil {
		return nil, fmt.Errorf("new reservation ledger config: %w", err)
	}

	reservationLedger, err := reservation.NewReservationLedger(*reservationConfig, time.Now())
	if err != nil {
		return nil, fmt.Errorf("new reservation ledger: %w", err)
	}

	return reservationLedger, nil
}

func buildOnDemandLedger(
	ctx context.Context,
	paymentVault payments.PaymentVault,
	accountID gethcommon.Address,
	minNumSymbols uint32,
	disperserClient *clientsv2.DisperserClient,
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
		return nil, fmt.Errorf("on-demand ledger from value: %w", err)
	}

	return onDemandLedger, nil
}
