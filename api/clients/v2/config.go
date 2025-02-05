package clients

import (
	"errors"
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

	// BlockNumberPollInterval is how frequently to check latest block number when waiting for the internal eth client
	// to advance to a certain block.
	//
	// If this is configured to be <= 0, then contract calls which require the internal eth client to have reached a
	// certain block height will fail if the internal client is behind.
	BlockNumberPollInterval time.Duration

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

	// The address of the BlsOperatorStateRetriever contract
	BlsOperatorStateRetrieverAddr string

	// The address of the EigenDAServiceManager contract
	EigenDAServiceManagerAddr string

	// The number of simultaneous connections to use when fetching chunks during validator retrieval
	ConnectionCount uint
}

// PayloadDisperserConfig contains an embedded PayloadClientConfig, plus all additional configuration values needed
// by a PayloadDisperser
type PayloadDisperserConfig struct {
	PayloadClientConfig

	// SignerPaymentKey is the private key used for signing payment authorization headers
	SignerPaymentKey string

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
		BlobEncodingVersion:     codecs.DefaultBlobEncoding,
		PayloadPolynomialForm:   codecs.PolynomialFormEval,
		ContractCallTimeout:     5 * time.Second,
		BlockNumberPollInterval: 1 * time.Second,
		BlobVersion:             0,
	}
}

// checkAndSetDefaults checks an existing config struct. It performs one of the following actions for any contained 0 values:
//
// 1. If 0 is an acceptable value for the field, do nothing.
// 2. If 0 is NOT an acceptable value for the field, and a default value is defined, then set it to the default.
// 3. If 0 is NOT an acceptable value for the field, and a default value is NOT defined, return an error.
func (cc *PayloadClientConfig) checkAndSetDefaults() error {
	// BlobEncodingVersion may be 0, so don't do anything

	if cc.EthRpcUrl == "" {
		return errors.New("EthRpcUrl is required")
	}

	if cc.EigenDACertVerifierAddr == "" {
		return errors.New("EigenDACertVerifierAddr is required")
	}

	// Nothing to do for PayloadPolynomialForm

	defaultConfig := getDefaultPayloadClientConfig()

	if cc.ContractCallTimeout == 0 {
		cc.ContractCallTimeout = defaultConfig.ContractCallTimeout
	}

	// BlockNumberPollInterval may be 0, so don't do anything

	// BlobVersion may be 0, so don't do anything

	return nil
}

// GetDefaultRelayPayloadRetrieverConfig creates a RelayPayloadRetrieverConfig with default values
//
// NOTE: EthRpcUrl and EigenDACertVerifierAddr do not have defined defaults. These must always be specifically configured.
func GetDefaultRelayPayloadRetrieverConfig() *RelayPayloadRetrieverConfig {
	return &RelayPayloadRetrieverConfig{
		PayloadClientConfig: *getDefaultPayloadClientConfig(),
		RelayTimeout:        5 * time.Second,
	}
}

// checkAndSetDefaults checks an existing config struct. It performs one of the following actions for any contained 0 values:
//
// 1. If 0 is an acceptable value for the field, do nothing.
// 2. If 0 is NOT an acceptable value for the field, and a default value is defined, then set it to the default.
// 3. If 0 is NOT an acceptable value for the field, and a default value is NOT defined, return an error.
func (rc *RelayPayloadRetrieverConfig) checkAndSetDefaults() error {
	err := rc.PayloadClientConfig.checkAndSetDefaults()
	if err != nil {
		return err
	}

	defaultConfig := GetDefaultRelayPayloadRetrieverConfig()
	if rc.RelayTimeout == 0 {
		rc.RelayTimeout = defaultConfig.RelayTimeout
	}

	return nil
}

// GetDefaultValidatorPayloadRetrieverConfig creates a ValidatorPayloadRetrieverConfig with default values
//
// NOTE: The following fields do not have defined defaults and must always be specifically configured:
// - EthRpcUrl
// - EigenDACertVerifierAddr
// - BlsOperatorStateRetrieverAddr
// - EigenDAServiceManagerAddr
func GetDefaultValidatorPayloadRetrieverConfig() *ValidatorPayloadRetrieverConfig {
	return &ValidatorPayloadRetrieverConfig{
		PayloadClientConfig: *getDefaultPayloadClientConfig(),
		RetrievalTimeout:    20 * time.Second,
		ConnectionCount:     10,
	}
}

// checkAndSetDefaults checks an existing config struct. It performs one of the following actions for any contained 0 values:
//
// 1. If 0 is an acceptable value for the field, do nothing.
// 2. If 0 is NOT an acceptable value for the field, and a default value is defined, then set it to the default.
// 3. If 0 is NOT an acceptable value for the field, and a default value is NOT defined, return an error.
func (rc *ValidatorPayloadRetrieverConfig) checkAndSetDefaults() error {
	err := rc.PayloadClientConfig.checkAndSetDefaults()
	if err != nil {
		return err
	}

	if rc.BlsOperatorStateRetrieverAddr == "" {
		return errors.New("BlsOperatorStateRetrieverAddr is required")
	}

	if rc.EigenDAServiceManagerAddr == "" {
		return errors.New("EigenDAServiceManagerAddr is required")
	}

	defaultConfig := GetDefaultValidatorPayloadRetrieverConfig()
	if rc.RetrievalTimeout == 0 {
		rc.RetrievalTimeout = defaultConfig.RetrievalTimeout
	}
	if rc.ConnectionCount == 0 {
		rc.ConnectionCount = defaultConfig.ConnectionCount
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

// checkAndSetDefaults checks an existing config struct. It performs one of the following actions for any contained 0 values:
//
// 1. If 0 is an acceptable value for the field, do nothing.
// 2. If 0 is NOT an acceptable value for the field, and a default value is defined, then set it to the default.
// 3. If 0 is NOT an acceptable value for the field, and a default value is NOT defined, return an error.
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
