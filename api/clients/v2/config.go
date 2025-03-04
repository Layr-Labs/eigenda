package clients

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
)

// PayloadClientConfig contains configuration values that are needed by both PayloadRetriever and PayloadDisperser
type PayloadClientConfig struct {
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

	// The BlobVersion to use when creating new blobs, or interpreting blob bytes.
	//
	// BlobVersion needs to point to a version defined in the threshold registry contract.
	// https://github.com/Layr-Labs/eigenda/blob/3ed9ef6ed3eb72c46ce3050eb84af28f0afdfae2/contracts/src/interfaces/IEigenDAThresholdRegistry.sol#L6
	BlobVersion v2.BlobVersion
}

// RelayPayloadRetrieverConfig contains an embedded PayloadClientConfig, plus all additional configuration values needed
// by a RelayPayloadRetriever
type RelayPayloadRetrieverConfig struct {
	PayloadClientConfig

	// The timeout duration for relay calls to retrieve blobs.
	RelayTimeout time.Duration
}

// ValidatorPayloadRetrieverConfig contains an embedded PayloadClientConfig, plus all additional configuration values
// needed by a ValidatorPayloadRetriever
type ValidatorPayloadRetrieverConfig struct {
	PayloadClientConfig

	// The timeout duration for retrieving chunks from a given quorum, and reassembling the chunks into a blob.
	// Once this timeout triggers, the retriever will give up on the quorum, and retry with the next quorum (if one exists)
	RetrievalTimeout time.Duration
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

	// The timeout duration for contract calls
	ContractCallTimeout time.Duration
}

// GetDefaultPayloadClientConfig creates a PayloadClientConfig with default values
func GetDefaultPayloadClientConfig() *PayloadClientConfig {
	return &PayloadClientConfig{
		PayloadPolynomialForm: codecs.PolynomialFormEval,
		BlobVersion:           0,
	}
}

// GetDefaultRelayPayloadRetrieverConfig creates a RelayPayloadRetrieverConfig with default values
//
// NOTE: EigenDACertVerifierAddr does not have a defined default. It must always be specifically configured.
func GetDefaultRelayPayloadRetrieverConfig() *RelayPayloadRetrieverConfig {
	return &RelayPayloadRetrieverConfig{
		PayloadClientConfig: *GetDefaultPayloadClientConfig(),
		RelayTimeout:        5 * time.Second,
	}
}

// checkAndSetDefaults checks an existing config struct. It performs one of the following actions for any contained 0 values:
//
// 1. If 0 is an acceptable value for the field, do nothing.
// 2. If 0 is NOT an acceptable value for the field, and a default value is defined, then set it to the default.
// 3. If 0 is NOT an acceptable value for the field, and a default value is NOT defined, return an error.
func (rc *RelayPayloadRetrieverConfig) checkAndSetDefaults() error {
	defaultConfig := GetDefaultRelayPayloadRetrieverConfig()
	if rc.RelayTimeout == 0 {
		rc.RelayTimeout = defaultConfig.RelayTimeout
	}

	return nil
}

// GetDefaultValidatorPayloadRetrieverConfig creates a ValidatorPayloadRetrieverConfig with default values
func GetDefaultValidatorPayloadRetrieverConfig() *ValidatorPayloadRetrieverConfig {
	return &ValidatorPayloadRetrieverConfig{
		PayloadClientConfig: *GetDefaultPayloadClientConfig(),
		RetrievalTimeout:    30 * time.Second,
	}
}

// checkAndSetDefaults checks an existing config struct. It performs one of the following actions for any contained 0 values:
//
// 1. If 0 is an acceptable value for the field, do nothing.
// 2. If 0 is NOT an acceptable value for the field, and a default value is defined, then set it to the default.
// 3. If 0 is NOT an acceptable value for the field, and a default value is NOT defined, return an error.
func (rc *ValidatorPayloadRetrieverConfig) checkAndSetDefaults() error {
	defaultConfig := GetDefaultValidatorPayloadRetrieverConfig()
	if rc.RetrievalTimeout == 0 {
		rc.RetrievalTimeout = defaultConfig.RetrievalTimeout
	}

	return nil
}

// GetDefaultPayloadDisperserConfig creates a PayloadDisperserConfig with default values
//
// NOTE: EigenDACertVerifierAddr does not have a defined default. It must always be specifically configured.
func GetDefaultPayloadDisperserConfig() *PayloadDisperserConfig {
	return &PayloadDisperserConfig{
		PayloadClientConfig:    *GetDefaultPayloadClientConfig(),
		DisperseBlobTimeout:    2 * time.Minute,
		BlobCertifiedTimeout:   2 * time.Minute,
		BlobStatusPollInterval: 1 * time.Second,
		ContractCallTimeout:    5 * time.Second,
	}
}

// checkAndSetDefaults checks an existing config struct. It performs one of the following actions for any contained 0 values:
//
// 1. If 0 is an acceptable value for the field, do nothing.
// 2. If 0 is NOT an acceptable value for the field, and a default value is defined, then set it to the default.
// 3. If 0 is NOT an acceptable value for the field, and a default value is NOT defined, return an error.
func (dc *PayloadDisperserConfig) checkAndSetDefaults() error {
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

	if dc.ContractCallTimeout == 0 {
		dc.ContractCallTimeout = defaultConfig.ContractCallTimeout
	}

	return nil
}
