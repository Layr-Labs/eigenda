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

	// The Ethereum RPC URL to use for querying the Ethereum blockchain.
	EthRpcUrl string

	// The address of the EigenDABlobVerifier contract
	EigenDABlobVerifierAddr string

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

	BlobVersion v2.BlobVersion

	Quorums []core.QuorumID
}
