package clients

import (
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
	//  array. A decision still needs to be made for how we want to handle this. Should this field be called
	//  `CustomQuorums`, and we simply append any values contained onto [0, 1]? Or should we require users to include 0
	//  and 1 here, and throw an error if they don't?
	Quorums []core.QuorumID
}
