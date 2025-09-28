package integration_test

import (
	"context"
	"math/big"
	"net"

	"github.com/Layr-Labs/eigenda/api/clients"
	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common"
	routerbindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierRouter"
	verifierv1bindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV1"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/testcontainers/testcontainers-go"
	"google.golang.org/grpc"
)

// InfrastructureHarness contains the shared infrastructure components
// that are global across all tests (external dependencies)
type InfrastructureHarness struct {
	// Infrastructure containers - truly global
	AnvilContainer      *testbed.AnvilContainer
	GraphNodeContainer  *testbed.GraphNodeContainer
	LocalstackContainer *testbed.LocalStackContainer
	ChainDockerNetwork  *testcontainers.DockerNetwork

	// EigenDA components
	// TODO: Should EigenDA components be their own test harness?
	ChurnerServer     *grpc.Server
	ChurnerListener   net.Listener
	OperatorInstances []*OperatorInstance

	// Global configuration
	TemplateName      string
	TestName          string
	InMemoryBlobStore bool
	LocalStackPort    string

	// DynamoDB table names (global for the test suite)
	MetadataTableName   string
	BucketTableName     string
	MetadataTableNameV2 string

	// Deployment configuration (shared)
	TestConfig *deploy.Config

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

	// Transaction options - specific to test
	DeployerTransactorOpts *bind.TransactOpts

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
