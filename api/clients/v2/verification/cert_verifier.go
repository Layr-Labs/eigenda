package verification

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/common"
	certVerifierBinding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifier"
	certVerifierV2Binding "github.com/Layr-Labs/eigenda/contracts/bindings/v2/EigenDACertVerifier"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// CertVerifier is responsible for making eth calls against version agnostic CertVerifier contracts to ensure
// cryptographic and structural integrity of EigenDA certificate types.
// The V3 cert verifier contract is located at:
// https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/periphery/cert/EigenDACertVerifier.sol
type CertVerifier struct {
	logger            logging.Logger
	ethClient         common.EthClient
	addressProvider   clients.CertVerifierAddressProvider
	v2VerifierBinding *certVerifierV2Binding.ContractEigenDACertVerifier

	// maps contract address to a ContractEigenDACertVerifierCaller object
	verifierCallers sync.Map
	// maps contract address to set of required quorums specified in the contract at that address
	requiredQuorums sync.Map
	// maps contract address to the confirmation threshold required by that address
	confirmationThresholds sync.Map
	// maps contract address to the cert version specified in the contract at that address
	versions sync.Map
}

// NewCertVerifier constructs a new CertVerifier instance
func NewCertVerifier(
	logger logging.Logger,
	ethClient common.EthClient,
	certVerifierAddressProvider clients.CertVerifierAddressProvider,
) (*CertVerifier, error) {
	return &CertVerifier{
		logger:            logger,
		ethClient:         ethClient,
		addressProvider:   certVerifierAddressProvider,
		v2VerifierBinding: certVerifierV2Binding.NewContractEigenDACertVerifier(),
	}, nil
}

// CheckDACert calls the CheckDACert view function on the EigenDACertVerifier contract.
// This method returns nil if the certificate is successfully verified; otherwise, it returns a
// [CertVerifierInvalidCertError] or [CertVerifierInternalError] error.
func (cv *CertVerifier) CheckDACert(
	ctx context.Context,
	cert coretypes.EigenDACert,
) error {
	// 1 - Normalize cert to V3
	certV3 := NormalizeCertV3(cert)

	// 2 - Call the contract method CheckDACert to verify the certificate
	// TODO: Determine adequate future proofing strategy for EigenDACertVerifierRouter to be compliant
	//       with future reference timestamp change which deprecates the reference block number
	//       used for quorum stake check-pointing.
	certVerifierAddr, err := cv.addressProvider.GetCertVerifierAddress(ctx, certV3.ReferenceBlockNumber())
	if err != nil {
		return &CertVerifierInternalError{Msg: "get verifier address", Err: err}
	}

	certBytes, err := certV3.Serialize(coretypes.CertSerializationABI)
	if err != nil {
		return &CertVerifierInternalError{Msg: "serialize cert", Err: err}
	}

	// TODO(ethenotethan): determine if there's any merit in passing call context
	// options (e.g, block number) to impose better determinism and safety on the simulation
	// call

	callMsgBytes, err := cv.v2VerifierBinding.TryPackCheckDACert(certBytes)
	if err != nil {
		return &CertVerifierInternalError{Msg: "pack checkDACert call", Err: err}
	}

	// TODO(ethenoethan): understand the best mechanisms for determining if the call ran into an
	// out-of-gas exception. Furthermore it's worth exploring whether an eth_simulateV1 rpc call
	// would provide better granularity and coverage while ensuring existing performance guarantees
	// see: https://www.quicknode.com/docs/ethereum/eth_simulateV1
	returnData, err := cv.ethClient.CallContract(ctx, ethereum.CallMsg{
		To:   &certVerifierAddr,
		Data: callMsgBytes,
	}, nil)
	if err != nil {
		return &CertVerifierInternalError{Msg: "checkDACert eth call", Err: err}
	}

	result, err := cv.v2VerifierBinding.UnpackCheckDACert(returnData)
	if err != nil {
		return &CertVerifierInternalError{Msg: "unpack checkDACert return data", Err: err}
	}

	// 3 - Cast result to structured enum type and check for not success status codes
	verifyResultCode := CheckDACertStatusCode(result)
	if verifyResultCode == StatusNullError {
		return &CertVerifierInternalError{Msg: fmt.Sprintf("checkDACert eth-call bug: %s", verifyResultCode.String())}
	} else if verifyResultCode != StatusSuccess {
		return &CertVerifierInvalidCertError{
			StatusCode: verifyResultCode,
			Msg:        verifyResultCode.String(),
		}
	}
	return nil
}

// EstimateGasCheckDACert uses eth_estimateGas to estimate the gas requirements for a CheckDACert call.
func (cv *CertVerifier) EstimateGasCheckDACert(
	ctx context.Context,
	cert coretypes.EigenDACert,
) (uint64, error) {
	// Normalize cert to V3
	certV3 := NormalizeCertV3(cert)

	certVerifierAddress, err := cv.addressProvider.GetCertVerifierAddress(
		ctx,
		certV3.ReferenceBlockNumber(),
	)
	if err != nil {
		return 0, fmt.Errorf("get cert verifier address: %w", err)
	}

	certBytes, err := certV3.Serialize(coretypes.CertSerializationABI)
	if err != nil {
		return 0, fmt.Errorf("serialize cert: %w", err)
	}

	// Pack the checkDACert method call data
	abi, err := certVerifierBinding.ContractEigenDACertVerifierMetaData.GetAbi()
	if err != nil {
		return 0, fmt.Errorf("get contract ABI: %w", err)
	}

	callData, err := abi.Pack("checkDACert", certBytes)
	if err != nil {
		return 0, fmt.Errorf("pack checkDACert call data: %w", err)
	}

	callMsg := ethereum.CallMsg{
		To:   &certVerifierAddress,
		Data: callData,
	}

	// Estimate gas using eth_estimateGas
	gasEstimate, err := cv.ethClient.EstimateGas(ctx, callMsg)
	if err != nil {
		cv.logger.Error(
			"eth_estimateGas",
			"to", callMsg.To.Hex(),
			"data", fmt.Sprintf("0x%x", callMsg.Data),
		)
		return 0, fmt.Errorf("estimate gas for checkDACert: %w", err)
	}

	return gasEstimate, nil
}

// GetQuorumNumbersRequired returns the set of quorum numbers that must be set in the BlobHeader, and verified in
// VerifyCert and CheckDACert.
//
// This method will return required quorum numbers from an internal cache if they are already known for the currently
// active cert verifier. Otherwise, this method will query the required quorum numbers from the currently active
// cert verifier, and cache the result for future use.
func (cv *CertVerifier) GetQuorumNumbersRequired(ctx context.Context) ([]uint8, error) {
	// get the latest cert verifier address from the address provider

	blockNum, err := cv.ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("get latest block number: %w", err)
	}

	certVerifierAddress, err := cv.addressProvider.GetCertVerifierAddress(ctx, blockNum.NumberU64())
	if err != nil {
		return nil, fmt.Errorf("get cert verifier address: %w", err)
	}

	// if the quorum numbers for the active cert verifier address have already been cached, return them immediately
	cachedQuorumNumbers, ok := cv.requiredQuorums.Load(certVerifierAddress)
	if ok {
		castQuorums, ok := cachedQuorumNumbers.([]uint8)
		if !ok {
			return nil, fmt.Errorf("expected quorum numbers to be []uint8")
		}
		return castQuorums, nil
	}

	// quorum numbers weren't cached, so proceed to fetch them
	certVerifierCaller, err := cv.getVerifierCallerFromAddress(certVerifierAddress)
	if err != nil {
		return nil, fmt.Errorf("get verifier caller from address: %w", err)
	}

	quorumNumbersRequired, err := certVerifierCaller.QuorumNumbersRequired(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("get quorum numbers required: %w", err)
	}

	cv.requiredQuorums.Store(certVerifierAddress, quorumNumbersRequired)

	return quorumNumbersRequired, nil
}

// getVerifierCallerFromAddress returns a ContractEigenDACertVerifier that corresponds to the input contract
// address
//
// This method caches ContractEigenDACertVerifier instances, since their construction requires acquiring a lock
// and parsing json, and is therefore non-trivially expensive.
func (cv *CertVerifier) getVerifierCallerFromAddress(
	certVerifierAddress gethcommon.Address,
) (*certVerifierBinding.ContractEigenDACertVerifier, error) {
	existingCallerAny, valueExists := cv.verifierCallers.Load(certVerifierAddress)
	if valueExists {
		existingCaller, ok := existingCallerAny.(*certVerifierBinding.ContractEigenDACertVerifier)
		if !ok {
			return nil, fmt.Errorf(
				"value in verifierCallers wasn't of type ContractEigenDACertVerifier. this should be impossible")
		}
		return existingCaller, nil
	}

	certVerifierCaller, err := certVerifierBinding.NewContractEigenDACertVerifier(certVerifierAddress, cv.ethClient)
	if err != nil {
		return nil, fmt.Errorf("bind to verifier contract at %s: %w", certVerifierAddress, err)
	}

	cv.verifierCallers.Store(certVerifierAddress, certVerifierCaller)
	return certVerifierCaller, nil
}

// GetConfirmationThreshold returns the ConfirmationThreshold that corresponds to the input reference block number.
//
// This method will return the confirmation threshold from an internal cache if it is already known for the cert
// verifier which corresponds to the input reference block number. Otherwise, this method will query the confirmation
// threshold and cache the result for future use.
func (cv *CertVerifier) GetConfirmationThreshold(ctx context.Context, referenceBlockNumber uint64) (uint8, error) {
	certVerifierAddress, err := cv.addressProvider.GetCertVerifierAddress(ctx, referenceBlockNumber)
	if err != nil {
		return 0, fmt.Errorf("get cert verifier address: %w", err)
	}

	// if the confirmation threshold for the active cert verifier address has already been cached, return it immediately
	cachedThreshold, ok := cv.confirmationThresholds.Load(certVerifierAddress)
	if ok {
		castThreshold, ok := cachedThreshold.(uint8)
		if !ok {
			return 0, fmt.Errorf("expected confirmation threshold to be uint8")
		}
		return castThreshold, nil
	}

	// confirmation threshold wasn't cached, so proceed to fetch it
	certVerifierCaller, err := cv.getVerifierCallerFromAddress(certVerifierAddress)
	if err != nil {
		return 0, fmt.Errorf("get verifier caller from address: %w", err)
	}

	securityThresholds, err := certVerifierCaller.SecurityThresholds(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, fmt.Errorf("get security thresholds via contract call: %w", err)
	}

	cv.confirmationThresholds.Store(certVerifierAddress, securityThresholds.ConfirmationThreshold)

	return securityThresholds.ConfirmationThreshold, nil
}

// GetCertVersion returns the CertVersion that corresponds to the input reference block number.
//
// This method will return the version from an internal cache if it is already known for the cert
// verifier which corresponds to the input reference block number. Otherwise, this method will query the version
// and cache the result for future use.
func (cv *CertVerifier) GetCertVersion(ctx context.Context, referenceBlockNumber uint64) (uint8, error) {
	certVerifierAddress, err := cv.addressProvider.GetCertVerifierAddress(ctx, referenceBlockNumber)
	if err != nil {
		return 0, fmt.Errorf("get cert verifier address: %w", err)
	}

	// if the version for the active cert verifier address has already been cached, return it immediately
	cachedVersion, ok := cv.versions.Load(certVerifierAddress)
	if ok {
		castVersion, ok := cachedVersion.(uint8)
		if !ok {
			return 0, fmt.Errorf("expected version to be uint64")
		}
		return castVersion, nil
	}

	// version wasn't cached, so proceed to fetch it
	certVerifierCaller, err := cv.getVerifierCallerFromAddress(certVerifierAddress)
	if err != nil {
		return 0, fmt.Errorf("get verifier caller from address: %w", err)
	}

	version, err := certVerifierCaller.CertVersion(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, fmt.Errorf("get version via contract call: %w", err)
	}

	cv.versions.Store(certVerifierAddress, version)

	return version, nil
}

// NormalizeCertV3 returns a EigenDACertV3 for a given EigenDACert
//
// This method normalizes a given EigenDACert (V2 or V3) to V3. If a V2 cert is given
// it is converted to V3 then returned, otherwise the given V3 cert is returned. All
// other versions will result in a panic.
func NormalizeCertV3(cert coretypes.EigenDACert) *coretypes.EigenDACertV3 {
	// switch on the certificate type to determine which contract to call
	var certV3 *coretypes.EigenDACertV3
	switch cert := cert.(type) {
	case *coretypes.EigenDACertV3:
		certV3 = cert
	case *coretypes.EigenDACertV2:
		// EigenDACertV3 is the only version that is supported by the CheckDACert function
		// but the V2 cert is a simple permutation of the V3 cert fields, so we convert it.
		certV3 = cert.ToV3()
	default:
		// If golang had enums the world would be a better place.
		panic(fmt.Sprintf("unsupported cert version: %T. All cert versions that we can "+
			"construct offchain should have a CertVerifier contract which we can call to "+
			"verify the certificate", cert))
	}

	return certV3
}
