package verification

import (
	"context"
	"fmt"
	"sync"

	clients "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/common"
	verifierBindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// LegacyCertVerifier is responsible for making eth calls against the LegacyCertVerifier contract to ensure cryptographic and
// structural integrity of V2 certificates. Currently this only used for eigenda v2 rollup tesnets using initial v2 release of eigenda-proxy.
// This will get deprecated in a future core release.
//
// The legacy cert verifier contract is located at https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/periphery/cert/legacy/v2/EigenDACertVerifierV2.sol
type LegacyCertVerifier struct {
	logger                      logging.Logger
	ethClient                   common.EthClient
	certVerifierAddressProvider clients.CertVerifierAddressProvider
	// maps contract address to a ContractEigenDACertVerifierV2Caller object
	verifierCallers sync.Map
}

// NewLegacyCertVerifier constructs a CertVerifier
func NewLegacyCertVerifier(
	logger logging.Logger,
	// the eth client, which should already be set up
	ethClient common.EthClient,
	certVerifierAddressProvider clients.CertVerifierAddressProvider,
) (*LegacyCertVerifier, error) {
	return &LegacyCertVerifier{
		logger:                      logger,
		ethClient:                   ethClient,
		certVerifierAddressProvider: certVerifierAddressProvider,
	}, nil
}

// VerifyCertV2 calls the VerifyCertV2 view function on the EigenDACertVerifier contract.
//
// This method returns nil if the cert is successfully verified. Otherwise, it returns an error.
func (cv *LegacyCertVerifier) VerifyCertV2(ctx context.Context, eigenDACert *coretypes.EigenDACertV2) error {
	referenceBlockNumber := uint64(eigenDACert.BatchHeader.ReferenceBlockNumber)

	certVerifierCaller, err := cv.getVerifierCallerFromBlockNumber(ctx, referenceBlockNumber)
	if err != nil {
		return fmt.Errorf("get verifier caller: %w", err)
	}

	err = certVerifierCaller.VerifyDACertV2(
		&bind.CallOpts{Context: ctx},
		eigenDACert.BatchHeader,
		eigenDACert.BlobInclusionInfo,
		eigenDACert.NonSignerStakesAndSignature,
		eigenDACert.SignedQuorumNumbers)

	if err != nil {
		return fmt.Errorf("verify cert v2: %w", err)
	}

	return nil
}


// getVerifierCallerFromBlockNumber returns a ContractEigenDACertVerifierV2Caller that corresponds to the input reference
// block number.
//
// This method caches ContractEigenDACertVerifierV2Caller instances, since their construction requires acquiring a lock
// and parsing json, and is therefore non-trivially expensive.
func (cv *LegacyCertVerifier) getVerifierCallerFromBlockNumber(
	ctx context.Context,
	referenceBlockNumber uint64,
) (*verifierBindings.ContractEigenDACertVerifierV2Caller, error) {
	certVerifierAddress, err := cv.certVerifierAddressProvider.GetCertVerifierAddress(ctx, referenceBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("get cert verifier address: %w", err)
	}

	return cv.getVerifierCallerFromAddress(certVerifierAddress)
}

// getVerifierCallerFromAddress returns a ContractEigenDACertVerifierV2Caller that corresponds to the input contract
// address
//
// This method caches ContractEigenDACertVerifierV2Caller instances, since their construction requires acquiring a lock
// and parsing json, and is therefore non-trivially expensive.
func (cv *LegacyCertVerifier) getVerifierCallerFromAddress(
	certVerifierAddress gethcommon.Address,
) (*verifierBindings.ContractEigenDACertVerifierV2Caller, error) {
	existingCallerAny, valueExists := cv.verifierCallers.Load(certVerifierAddress)
	if valueExists {
		existingCaller, ok := existingCallerAny.(*verifierBindings.ContractEigenDACertVerifierV2Caller)
		if !ok {
			return nil, fmt.Errorf(
				"value in verifierCallers wasn't of type ContractEigenDACertVerifierV2Caller. this should be impossible")
		}
		return existingCaller, nil
	}

	certVerifierCaller, err := verifierBindings.NewContractEigenDACertVerifierV2Caller(certVerifierAddress, cv.ethClient)
	if err != nil {
		return nil, fmt.Errorf("bind to verifier contract at %s: %w", certVerifierAddress, err)
	}

	cv.verifierCallers.Store(certVerifierAddress, certVerifierCaller)
	return certVerifierCaller, nil
}
