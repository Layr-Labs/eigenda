package deployment

import (
	"context"
	"fmt"

	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common"
	routerbindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierRouter"
	verifierv1bindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV1"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// CertVerificationComponents holds all certification verification components
type CertVerificationComponents struct {
	CertBuilder                      *clientsv2.CertBuilder
	RouterCertVerifier               *verification.CertVerifier
	StaticCertVerifier               *verification.CertVerifier
	EigenDACertVerifierV1            *verifierv1bindings.ContractEigenDACertVerifierV1
	EigenDACertVerifierRouter        *routerbindings.ContractEigenDACertVerifierRouterTransactor
	EigenDACertVerifierRouterCaller  *routerbindings.ContractEigenDACertVerifierRouterCaller
}

// InitializeCertVerification sets up the certification verification components
func InitializeCertVerification(
	ctx context.Context,
	ethClient common.EthClient,
	logger logging.Logger,
	contracts *EigenDAContracts,
) (*CertVerificationComponents, error) {
	if contracts == nil {
		return nil, fmt.Errorf("EigenDA contracts are required for cert verification")
	}

	components := &CertVerificationComponents{}

	// Initialize EigenDACertVerifierV1
	eigenDACertVerifierV1, err := verifierv1bindings.NewContractEigenDACertVerifierV1(
		gethcommon.HexToAddress(contracts.EigenDAV1CertVerifier),
		ethClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create EigenDACertVerifierV1: %w", err)
	}
	components.EigenDACertVerifierV1 = eigenDACertVerifierV1

	// Initialize CertBuilder
	certBuilder, err := clientsv2.NewCertBuilder(
		logger,
		gethcommon.HexToAddress(contracts.OperatorStateRetriever),
		gethcommon.HexToAddress(contracts.RegistryCoordinator),
		ethClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create CertBuilder: %w", err)
	}
	components.CertBuilder = certBuilder

	// Initialize router address provider and cert verifier if router is available
	if contracts.EigenDACertVerifierRouter != "" {
		routerAddress := gethcommon.HexToAddress(contracts.EigenDACertVerifierRouter)
		
		// Create router transactor binding
		eigenDACertVerifierRouter, err := routerbindings.NewContractEigenDACertVerifierRouterTransactor(
			routerAddress,
			ethClient,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create EigenDACertVerifierRouter transactor: %w", err)
		}
		components.EigenDACertVerifierRouter = eigenDACertVerifierRouter

		// Create router caller binding
		eigenDACertVerifierRouterCaller, err := routerbindings.NewContractEigenDACertVerifierRouterCaller(
			routerAddress,
			ethClient,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create EigenDACertVerifierRouter caller: %w", err)
		}
		components.EigenDACertVerifierRouterCaller = eigenDACertVerifierRouterCaller

		// Create router address provider
		routerAddressProvider, err := verification.BuildRouterAddressProvider(
			routerAddress,
			ethClient,
			logger,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to build router address provider: %w", err)
		}

		// Create router cert verifier
		routerCertVerifier, err := verification.NewCertVerifier(
			logger,
			ethClient,
			routerAddressProvider,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create router cert verifier: %w", err)
		}
		components.RouterCertVerifier = routerCertVerifier
	}

	// Initialize static cert verifier using V2 cert verifier address
	certVerifierAddress := contracts.EigenDAV2CertVerifier
	if certVerifierAddress == "" {
		// Fallback to V1 if V2 is not available
		certVerifierAddress = contracts.EigenDAV1CertVerifier
	}

	if certVerifierAddress != "" {
		staticAddressProvider := verification.NewStaticCertVerifierAddressProvider(
			gethcommon.HexToAddress(certVerifierAddress),
		)

		staticCertVerifier, err := verification.NewCertVerifier(
			logger,
			ethClient,
			staticAddressProvider,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create static cert verifier: %w", err)
		}
		components.StaticCertVerifier = staticCertVerifier
	}

	return components, nil
}