package clients

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common/geth"
	contractEigenDABlobVerifier "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDABlobVerifier"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// EigenDAClient provides the ability to get payloads from the relay subsystem, and to send new payloads to the disperser.
//
// This struct is goroutine safe.
type EigenDAClient struct {
	log logging.Logger
	// doesn't need to be cryptographically secure, as it's only used to distribute load across relays
	random       *rand.Rand
	clientConfig *EigenDAClientConfig
	codec        codecs.BlobCodec
	relayClient  RelayClient
	g1Srs        []bn254.G1Affine
	blobVerifier verification.IBlobVerifier
}

// BuildEigenDAClient builds an EigenDAClient from config structs.
func BuildEigenDAClient(
	log logging.Logger,
	clientConfig *EigenDAClientConfig,
	ethConfig geth.EthClientConfig,
	relayClientConfig *RelayClientConfig,
	g1Srs []bn254.G1Affine) (*EigenDAClient, error) {

	relayClient, err := NewRelayClient(relayClientConfig, log)
	if err != nil {
		return nil, fmt.Errorf("new relay client: %w", err)
	}

	ethClient, err := geth.NewClient(ethConfig, gethcommon.Address{}, 0, log)
	if err != nil {
		return nil, fmt.Errorf("new eth client: %w", err)
	}

	blobVerifier, err := verification.NewBlobVerifier(*ethClient, clientConfig.EigenDABlobVerifierAddr)
	if err != nil {
		return nil, fmt.Errorf("new blob verifier: %w", err)
	}

	codec, err := createCodec(clientConfig)
	if err != nil {
		return nil, err
	}

	return NewEigenDAClient(
		log,
		rand.New(rand.NewSource(rand.Int63())),
		clientConfig,
		relayClient,
		blobVerifier,
		codec,
		g1Srs)
}

// NewEigenDAClient assembles an EigenDAClient from subcomponents that have already been constructed and initialized.
func NewEigenDAClient(
	log logging.Logger,
	random *rand.Rand,
	clientConfig *EigenDAClientConfig,
	relayClient RelayClient,
	blobVerifier verification.IBlobVerifier,
	codec codecs.BlobCodec,
	g1Srs []bn254.G1Affine) (*EigenDAClient, error) {

	return &EigenDAClient{
		log:          log,
		random:       random,
		clientConfig: clientConfig,
		codec:        codec,
		relayClient:  relayClient,
		blobVerifier: blobVerifier,
		g1Srs:        g1Srs,
	}, nil
}

// GetPayload iteratively attempts to fetch a given blob with key blobKey from relays that have it, as claimed by the
// blob certificate. The relays are attempted in random order.
//
// If the blob is successfully retrieved, then the blob is verified. If the verification succeeds, the blob is decoded
// to yield the payload (the original user data), and the payload is returned.
func (c *EigenDAClient) GetPayload(
	ctx context.Context,
	blobKey core.BlobKey,
	eigenDACert *verification.EigenDACert) ([]byte, error) {

	err := c.verifyCertWithTimeout(ctx, eigenDACert)
	if err != nil {
		return nil, fmt.Errorf("verify cert with timeout for blobKey %v: %w", blobKey, err)
	}

	relayKeys := eigenDACert.BlobVerificationProof.BlobCertificate.RelayKeys
	relayKeyCount := len(relayKeys)
	if relayKeyCount == 0 {
		return nil, errors.New("relay key count is zero")
	}

	blobCommitments, err := blobCommitmentsBindingToInternal(
		&eigenDACert.BlobVerificationProof.BlobCertificate.BlobHeader.Commitment)

	if err != nil {
		return nil, fmt.Errorf("blob commitments binding to internal: %w", err)
	}

	// create a randomized array of indices, so that it isn't always the first relay in the list which gets hit
	indices := c.random.Perm(relayKeyCount)

	// TODO (litt3): consider creating a utility which deprioritizes relays that fail to respond (or respond maliciously),
	//  and prioritizes relays with lower latencies.

	// iterate over relays in random order, until we are able to get the blob from someone
	for _, val := range indices {
		relayKey := relayKeys[val]

		blob, err := c.getBlobWithTimeout(ctx, relayKey, blobKey)
		// if GetBlob returned an error, try calling a different relay
		if err != nil {
			c.log.Warn("blob couldn't be retrieved from relay", "blobKey", blobKey, "relayKey", relayKey, "error", err)
			continue
		}

		err = c.verifyBlobAgainstCert(blobKey, relayKey, blob, blobCommitments.Commitment, blobCommitments.Length)

		// An honest relay should never send a blob which doesn't verify
		if err != nil {
			c.log.Warn("verify blob from relay: %w", err)
			continue
		}

		payload, err := c.codec.DecodeBlob(blob)
		if err != nil {
			c.log.Error(
				`Blob verification was successful, but decode blob failed!
					This is likely a problem with the local blob codec configuration,
					but could potentially indicate a maliciously generated blob certificate.
					It should not be possible for an honestly generated certificate to verify
					for an invalid blob!`,
				"blobKey", blobKey, "relayKey", relayKey, "eigenDACert", eigenDACert, "error", err)
			return nil, fmt.Errorf("decode blob: %w", err)
		}

		return payload, nil
	}

	return nil, fmt.Errorf("unable to retrieve blob %v from any relay. relay count: %d", blobKey, relayKeyCount)
}

// verifyBlobAgainstCert verifies the blob received from a relay against the certificate.
//
// The following verifications are performed in this method:
// 1. Verify that the blob isn't empty
// 2. Verify the blob against the cert's kzg commitment
// 3. Verify that the blob length is less than or equal to the cert's blob length
//
// If all verifications succeed, the method returns nil. Otherwise, it returns an error.
func (c *EigenDAClient) verifyBlobAgainstCert(
	blobKey core.BlobKey,
	relayKey core.RelayKey,
	blob []byte,
	kzgCommitment *encoding.G1Commitment,
	blobLength uint) error {

	// An honest relay should never send an empty blob
	if len(blob) == 0 {
		return fmt.Errorf("blob %v received from relay %v had length 0", blobKey, relayKey)
	}

	// TODO: in the future, this will be optimized to use fiat shamir transformation for verification, rather than
	//  regenerating the commitment: https://github.com/Layr-Labs/eigenda/issues/1037
	valid, err := verification.GenerateAndCompareBlobCommitment(c.g1Srs, blob, kzgCommitment)
	if err != nil {
		return fmt.Errorf(
			"generate and compare commitment for blob %v received from relay %v: %w",
			blobKey,
			relayKey,
			err)
	}

	if !valid {
		return fmt.Errorf("commitment for blob %v is invalid for bytes received from relay %v", blobKey, relayKey)
	}

	// Checking that the length returned by the relay is <= the length claimed in the BlobCommitments is sufficient
	// here: it isn't necessary to verify the length proof itself, since this will have been done by DA nodes prior to
	// signing for availability.
	//
	// Note that the length in the commitment is the length of the blob in symbols
	if uint(len(blob)) > blobLength*encoding.BYTES_PER_SYMBOL {
		return fmt.Errorf(
			"length for blob %v (%d bytes) received from relay %v is greater than claimed blob length (%d bytes)",
			blobKey,
			len(blob),
			relayKey,
			blobLength*encoding.BYTES_PER_SYMBOL)
	}

	return nil
}

// getBlobWithTimeout attempts to get a blob from a given relay, and times out based on config.RelayTimeout
func (c *EigenDAClient) getBlobWithTimeout(
	ctx context.Context,
	relayKey core.RelayKey,
	blobKey core.BlobKey) ([]byte, error) {

	timeoutCtx, cancel := context.WithTimeout(ctx, c.clientConfig.RelayTimeout)
	defer cancel()

	return c.relayClient.GetBlob(timeoutCtx, relayKey, blobKey)
}

// verifyCertWithTimeout verifies an EigenDACert by making a call to VerifyBlobV2.
//
// This method times out after the duration configured in clientConfig.ContractCallTimeout
func (c *EigenDAClient) verifyCertWithTimeout(
	ctx context.Context,
	eigenDACert *verification.EigenDACert,
) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, c.clientConfig.ContractCallTimeout)
	defer cancel()

	return c.blobVerifier.VerifyBlobV2(timeoutCtx, eigenDACert)
}

// Close is responsible for calling close on all internal clients. This method will do its best to close all internal
// clients, even if some closes fail.
//
// Any and all errors returned from closing internal clients will be joined and returned.
//
// This method should only be called once.
func (c *EigenDAClient) Close() error {
	relayClientErr := c.relayClient.Close()

	// TODO: this is using join, since there will be more subcomponents requiring closing after adding PUT functionality
	return errors.Join(relayClientErr)
}

// createCodec creates the codec based on client config values
func createCodec(config *EigenDAClientConfig) (codecs.BlobCodec, error) {
	lowLevelCodec, err := codecs.BlobEncodingVersionToCodec(config.BlobEncodingVersion)
	if err != nil {
		return nil, fmt.Errorf("create low level codec: %w", err)
	}

	switch config.PayloadPolynomialForm {
	case codecs.PolynomialFormCoeff:
		// Data must NOT be IFFTed during blob construction, since the payload is already in PolynomialFormCoeff after
		// being encoded.
		return codecs.NewNoIFFTCodec(lowLevelCodec), nil
	case codecs.PolynomialFormEval:
		// Data MUST be IFFTed during blob construction, since the payload is in PolynomialFormEval after being encoded,
		// but must be in PolynomialFormCoeff to produce a valid blob.
		return codecs.NewIFFTCodec(lowLevelCodec), nil
	default:
		return nil, fmt.Errorf("unsupported polynomial form: %d", config.PayloadPolynomialForm)
	}
}

// blobCommitmentsBindingToInternal converts a blob commitment from an eigenDA cert into the internal
// encoding.BlobCommitments type
func blobCommitmentsBindingToInternal(
	blobCommitmentBinding *contractEigenDABlobVerifier.BlobCommitment,
) (*encoding.BlobCommitments, error) {

	blobCommitment, err := encoding.BlobCommitmentsFromProtobuf(
		verification.BlobCommitmentBindingToProto(blobCommitmentBinding))

	if err != nil {
		return nil, fmt.Errorf("blob commitments from protobuf: %w", err)
	}

	return blobCommitment, nil
}
