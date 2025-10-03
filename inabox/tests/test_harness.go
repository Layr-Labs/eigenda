package integration

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients"
	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common"
	routerbindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierRouter"
	verifierv1bindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV1"
	paymentvaultbindings "github.com/Layr-Labs/eigenda/contracts/bindings/PaymentVault"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
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
	TestConfig        *deploy.Config
	TemplateName      string
	TestName          string
	InMemoryBlobStore bool
	LocalStackPort    string

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
	PayloadDisperser           *payloaddispersal.PayloadDisperser

	// Core components
	ChainReader core.Reader

	// PaymentVault interaction
	PaymentVaultTransactor *paymentvaultbindings.ContractPaymentVaultTransactor

	// Transaction options - specific to test
	DeployerTransactorOpts *bind.TransactOpts

	// Test-specific configuration
	NumConfirmations int
	NumRetries       int

	// Chain ID for this test context
	ChainID *big.Int

	// Test account ID used for dispersals and reservations
	TestAccountID gethcommon.Address
}

// Cleanup releases resources held by the TestHarness
func (tc *TestHarness) Cleanup() {
	// Clean up any test-specific resources if needed
	// Most will be garbage collected, but connections will be closed when EthClient is garbage collected
}

// Updates the reservation for the test account on the PaymentVault contract
func (tc *TestHarness) UpdateReservationOnChain(ctx context.Context, t *testing.T, r *reservation.Reservation) error {
	quorumNumbers := r.GetQuorumNumbers()
	quorumSplits := calculateQuorumSplits(len(quorumNumbers))

	newReservation := paymentvaultbindings.IPaymentVaultReservation{
		SymbolsPerSecond: r.GetSymbolsPerSecond(),
		StartTimestamp:   uint64(r.GetStartTime().Unix()),
		EndTimestamp:     uint64(r.GetEndTime().Unix()),
		QuorumNumbers:    quorumNumbers,
		QuorumSplits:     quorumSplits,
	}

	tx, err := tc.PaymentVaultTransactor.SetReservation(
		tc.DeployerTransactorOpts,
		tc.TestAccountID,
		newReservation,
	)
	if err != nil {
		return fmt.Errorf("set reservation: %w", err)
	}

	MineAnvilBlocks(t, tc.RPCClient, 1)
	receipt, err := bind.WaitMined(ctx, tc.EthClient, tx)
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
