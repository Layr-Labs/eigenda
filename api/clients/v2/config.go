package clients

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
)

// EigenDAClientConfig contains configuration values for EigenDAClient
type EigenDAClientConfig struct {
	// The blob encoding version to use when writing and reading blobs
	BlobEncodingVersion codecs.BlobEncodingVersion

	// The Ethereum RPC URL to use for querying the Ethereum blockchain.
	EthRpcUrl string

	// The address of the EigenDABlobVerifier contract
	EigenDABlobVerifierAddr string

	// BlobPolynomialForm is the form that the blob polynomial is commited to and dispersed in, as well as the form the
	// blob polynomial will be received in from the relay.
	//
	// The chosen form dictates how the KZG commitment made to the blob can be used. If the polynomial is in Coeff form
	// when committed to, then it will be possible to open points on the KZG commitment to prove that the field elements
	// correspond to the commitment. If the polynomial is in Eval form when committed to, then it will not be possible
	// to create a commitment opening: the blob will need to be supplied in its entirety to perform a verification that
	// any part of the data matches the KZG commitment.
	BlobPolynomialForm codecs.PolynomialForm

	// The timeout duration for relay calls
	RelayTimeout time.Duration

	// The timeout duration for contract calls
	ContractCallTimeout time.Duration
}
