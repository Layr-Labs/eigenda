package clients

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
)

// PayloadClientConfig contains configuration values that are needed by both PayloadRetriever and PayloadDisperser
type PayloadClientConfig struct {
	// The blob encoding version to use when writing and reading blobs
	BlobEncodingVersion codecs.BlobEncodingVersion

	// The Ethereum RPC URL to use for querying an Ethereum network
	EthRpcUrl string

	// The address of the EigenDACertVerifier contract
	EigenDACertVerifierAddr string

	// PayloadPolynomialForm is the initial form of a Payload after being encoded. The configured form does not imply
	// any restrictions on the contents of a payload: it merely dictates how payload data is treated after being
	// encoded.
	//
	// Since blobs sent to the disperser must be in coefficient form, the initial form of the encoded payload dictates
	// what data processing must be performed during blob construction.
	//
	// The chosen form also dictates how the KZG commitment made to the blob can be used. If the encoded payload starts
	// in PolynomialFormEval (meaning the data WILL be IFFTed before computing the commitment) then it will be possible
	// to open points on the KZG commitment to prove that the field elements correspond to the commitment. If the
	// encoded payload starts in PolynomialFormCoeff (meaning the data will NOT be IFFTed before computing the
	// commitment) then it will not be possible to create a commitment opening: the blob will need to be supplied in its
	// entirety to perform a verification that any part of the data matches the KZG commitment.
	PayloadPolynomialForm codecs.PolynomialForm

	// The timeout duration for contract calls
	ContractCallTimeout time.Duration

	// The BlobVersion to use when creating new blobs, or interpreting blob bytes.
	//
	// BlobVersion needs to point to a version defined in the threshold registry contract.
	// https://github.com/Layr-Labs/eigenda/blob/3ed9ef6ed3eb72c46ce3050eb84af28f0afdfae2/contracts/src/interfaces/IEigenDAThresholdRegistry.sol#L6
	BlobVersion v2.BlobVersion
}

// PayloadRetrieverConfig contains an embedded PayloadClientConfig, plus all additional configuration values needed
// by a PayloadRetriever
type PayloadRetrieverConfig struct {
	PayloadClientConfig

	// The timeout duration for relay calls to retrieve blobs.
	RelayTimeout time.Duration
}

// PayloadDisperserConfig contains an embedded PayloadClientConfig, plus all additional configuration values needed
// by a PayloadDisperser
type PayloadDisperserConfig struct {
	PayloadClientConfig

	// DisperseBlobTimeout is the duration after which the PayloadDisperser will time out, when trying to disperse a
	// blob
	DisperseBlobTimeout time.Duration

	// BlobCertifiedTimeout is the duration after which the PayloadDisperser will time out, while polling
	// the disperser for blob status, waiting for BlobStatus_CERTIFIED
	BlobCertifiedTimeout time.Duration

	// BlobStatusPollInterval is the tick rate for the PayloadDisperser to use, while polling the disperser with
	// GetBlobStatus.
	BlobStatusPollInterval time.Duration

	// Quorums is the set of quorums that need to have a threshold of signatures for an EigenDA cert to successfully
	// verify.
	//
	// TODO: Clients are currently charged for QuorumIDs 0 and 1 regardless of whether or not they are included in this
	//  array. Therefore, if 0 and 1 aren't included in this array, you are missing out on security that your are paying
	//  for. The strategy for how to handle this field in the context of rollups is still in flux: this comment should
	//  be revisited and revised as necessary.
	Quorums []core.QuorumID
}

// GetDefaultPayloadClientConfig creates a PayloadClientConfig with default values
//
// NOTE: EthRpcUrl and EigenDACertVerifierAddr do not have defined defaults. These must always be specifically configured.
func getDefaultPayloadClientConfig() *PayloadClientConfig {
	return &PayloadClientConfig{
		BlobEncodingVersion:   codecs.DefaultBlobEncoding,
		PayloadPolynomialForm: codecs.PolynomialFormEval,
		ContractCallTimeout:   5 * time.Second,
		BlobVersion:           0,
	}
}

// checkAndSetDefaults checks an existing config struct and performs the following actions:
//
// 1. If a config value is 0, and a 0 value makes sense, do nothing.
// 2. If a config value is 0, but a 0 value doesn't make sense and a default value is defined, then set it to the default.
// 3. If a config value is 0, but a 0 value doesn't make sense and a default value isn't defined, return an error.
func (cc *PayloadClientConfig) checkAndSetDefaults() error {
	// BlobEncodingVersion may be 0, so don't do anything

	if cc.EthRpcUrl == "" {
		return fmt.Errorf("EthRpcUrl is required")
	}

	if cc.EigenDACertVerifierAddr == "" {
		return fmt.Errorf("EigenDACertVerifierAddr is required")
	}

	// Nothing to do for PayloadPolynomialForm

	defaultConfig := getDefaultPayloadClientConfig()

	if cc.ContractCallTimeout == 0 {
		cc.ContractCallTimeout = defaultConfig.ContractCallTimeout
	}

	// BlobVersion may be 0, so don't do anything

	return nil
}

// GetDefaultPayloadRetrieverConfig creates a PayloadRetrieverConfig with default values
//
// NOTE: EthRpcUrl and EigenDACertVerifierAddr do not have defined defaults. These must always be specifically configured.
func GetDefaultPayloadRetrieverConfig() *PayloadRetrieverConfig {
	return &PayloadRetrieverConfig{
		PayloadClientConfig: *getDefaultPayloadClientConfig(),
		RelayTimeout:        5 * time.Second,
	}
}

// checkAndSetDefaults checks an existing config struct and performs the following actions:
//
// 1. If a config value is 0, and a 0 value makes sense, do nothing.
// 2. If a config value is 0, but a 0 value doesn't make sense and a default value is defined, then set it to the default.
// 3. If a config value is 0, but a 0 value doesn't make sense and a default value isn't defined, return an error.
func (rc *PayloadRetrieverConfig) checkAndSetDefaults() error {
	err := rc.PayloadClientConfig.checkAndSetDefaults()
	if err != nil {
		return err
	}

	defaultConfig := GetDefaultPayloadRetrieverConfig()
	if rc.RelayTimeout == 0 {
		rc.RelayTimeout = defaultConfig.RelayTimeout
	}

	return nil
}

// GetDefaultPayloadDisperserConfig creates a PayloadDisperserConfig with default values
//
// NOTE: EthRpcUrl and EigenDACertVerifierAddr do not have defined defaults. These must always be specifically configured.
func GetDefaultPayloadDisperserConfig() *PayloadDisperserConfig {
	return &PayloadDisperserConfig{
		PayloadClientConfig:    *getDefaultPayloadClientConfig(),
		DisperseBlobTimeout:    5 * time.Second,
		BlobCertifiedTimeout:   10 * time.Second,
		BlobStatusPollInterval: 1 * time.Second,
		Quorums:                []core.QuorumID{0, 1},
	}
}

// checkAndSetDefaults checks an existing config struct and performs the following actions:
//
// 1. If a config value is 0, and a 0 value makes sense, do nothing.
// 2. If a config value is 0, but a 0 value doesn't make sense and a default value is defined, then set it to the default.
// 3. If a config value is 0, but a 0 value doesn't make sense and a default value isn't defined, return an error.
func (dc *PayloadDisperserConfig) checkAndSetDefaults() error {
	err := dc.PayloadClientConfig.checkAndSetDefaults()
	if err != nil {
		return err
	}

	defaultConfig := GetDefaultPayloadDisperserConfig()

	if dc.DisperseBlobTimeout == 0 {
		dc.DisperseBlobTimeout = defaultConfig.DisperseBlobTimeout
	}

	if dc.BlobCertifiedTimeout == 0 {
		dc.BlobCertifiedTimeout = defaultConfig.BlobCertifiedTimeout
	}

	if dc.BlobStatusPollInterval == 0 {
		dc.BlobStatusPollInterval = defaultConfig.BlobStatusPollInterval
	}

	// Quorums may be empty, so don't do anything

	return nil
}
