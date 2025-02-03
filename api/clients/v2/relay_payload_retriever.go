package clients

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	verifiercontract "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifier"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
)

// RelayPayloadRetriever provides the ability to get payloads from the relay subsystem.
//
// This struct is goroutine safe.
type RelayPayloadRetriever struct {
	log logging.Logger
	// random doesn't need to be cryptographically secure, as it's only used to distribute load across relays.
	// Not all methods on Rand are guaranteed goroutine safe: if additional usages of random are added, they
	// must be evaluated for thread safety.
	random      *rand.Rand
	config      RelayPayloadRetrieverConfig
	codec       codecs.BlobCodec
	relayClient RelayClient
	g1Srs       []bn254.G1Affine
}

var _ PayloadRetriever = &RelayPayloadRetriever{}

// BuildRelayPayloadRetriever builds a RelayPayloadRetriever from config structs.
func BuildRelayPayloadRetriever(
	log logging.Logger,
	relayPayloadRetrieverConfig RelayPayloadRetrieverConfig,
	relayClientConfig *RelayClientConfig,
	g1Srs []bn254.G1Affine) (*RelayPayloadRetriever, error) {

	relayClient, err := NewRelayClient(relayClientConfig, log)
	if err != nil {
		return nil, fmt.Errorf("new relay client: %w", err)
	}

	codec, err := codecs.CreateCodec(
		relayPayloadRetrieverConfig.PayloadPolynomialForm,
		relayPayloadRetrieverConfig.BlobEncodingVersion)
	if err != nil {
		return nil, err
	}

	return NewRelayPayloadRetriever(
		log,
		rand.New(rand.NewSource(rand.Int63())),
		relayPayloadRetrieverConfig,
		relayClient,
		codec,
		g1Srs)
}

// NewRelayPayloadRetriever assembles a RelayPayloadRetriever from subcomponents that have already been constructed and
// initialized.
func NewRelayPayloadRetriever(
	log logging.Logger,
	random *rand.Rand,
	relayPayloadRetrieverConfig RelayPayloadRetrieverConfig,
	relayClient RelayClient,
	codec codecs.BlobCodec,
	g1Srs []bn254.G1Affine) (*RelayPayloadRetriever, error) {

	err := relayPayloadRetrieverConfig.checkAndSetDefaults()
	if err != nil {
		return nil, fmt.Errorf("check and set RelayPayloadRetrieverConfig config: %w", err)
	}

	return &RelayPayloadRetriever{
		log:         log,
		random:      random,
		config:      relayPayloadRetrieverConfig,
		codec:       codec,
		relayClient: relayClient,
		g1Srs:       g1Srs,
	}, nil
}

// GetPayload iteratively attempts to fetch a given blob with key blobKey from relays that have it, as claimed by the
// blob certificate. The relays are attempted in random order.
//
// If the blob is successfully retrieved, then the blob is verified against the certificate. If the verification
// succeeds, the blob is decoded to yield the payload (the original user data, with no padding or any modification),
// and the payload is returned.
//
// This method does NOT verify the eigenDACert on chain: it is assumed that the input eigenDACert has already been
// verified prior to calling this method.
func (pr *RelayPayloadRetriever) GetPayload(
	ctx context.Context,
	eigenDACert *verification.EigenDACert) ([]byte, error) {

	blobKey, err := eigenDACert.ComputeBlobKey()
	if err != nil {
		return nil, fmt.Errorf("compute blob key: %w", err)
	}

	relayKeys := eigenDACert.BlobInclusionInfo.BlobCertificate.RelayKeys
	relayKeyCount := len(relayKeys)
	if relayKeyCount == 0 {
		return nil, errors.New("relay key count is zero")
	}

	blobCommitments, err := blobCommitmentsBindingToInternal(
		&eigenDACert.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment)

	if err != nil {
		return nil, fmt.Errorf("blob commitments binding to internal: %w", err)
	}

	// create a randomized array of indices, so that it isn't always the first relay in the list which gets hit
	indices := pr.random.Perm(relayKeyCount)

	// TODO (litt3): consider creating a utility which deprioritizes relays that fail to respond (or respond maliciously),
	//  and prioritizes relays with lower latencies.

	// iterate over relays in random order, until we are able to get the blob from someone
	for _, val := range indices {
		relayKey := relayKeys[val]

		blob, err := pr.getBlobWithTimeout(ctx, relayKey, blobKey)
		// if GetBlob returned an error, try calling a different relay
		if err != nil {
			pr.log.Warn(
				"blob couldn't be retrieved from relay",
				"blobKey", blobKey.Hex(),
				"relayKey", relayKey,
				"error", err)
			continue
		}

		err = pr.verifyBlobAgainstCert(blobKey, relayKey, blob, blobCommitments.Commitment, blobCommitments.Length)

		// An honest relay should never send a blob which doesn't verify against the cert
		if err != nil {
			pr.log.Warn("verify blob from relay against cert: %w", err)
			continue
		}

		payload, err := pr.codec.DecodeBlob(blob)
		if err != nil {
			pr.log.Error(
				`Cert verification was successful, but decode blob failed!
					This is likely a problem with the local blob codec configuration,
					but could potentially indicate a maliciously generated blob certificate.
					It should not be possible for an honestly generated certificate to verify
					for an invalid blob!`,
				"blobKey", blobKey.Hex(), "relayKey", relayKey, "eigenDACert", eigenDACert, "error", err)
			return nil, fmt.Errorf("decode blob: %w", err)
		}

		return payload, nil
	}

	return nil, fmt.Errorf("unable to retrieve blob %v from any relay. relay count: %d", blobKey.Hex(), relayKeyCount)
}

// verifyBlobAgainstCert verifies the blob received from a relay against the certificate.
//
// The following verifications are performed in this method:
// 1. Verify that the blob isn't empty
// 2. Verify the blob against the cert's kzg commitment
// 3. Verify that the blob length is less than or equal to the cert's blob length
//
// If all verifications succeed, the method returns nil. Otherwise, it returns an error.
func (pr *RelayPayloadRetriever) verifyBlobAgainstCert(
	blobKey *core.BlobKey,
	relayKey core.RelayKey,
	blob []byte,
	kzgCommitment *encoding.G1Commitment,
	blobLength uint) error {

	// An honest relay should never send an empty blob
	if len(blob) == 0 {
		return fmt.Errorf("blob %v received from relay %v had length 0", blobKey.Hex(), relayKey)
	}

	// TODO: in the future, this will be optimized to use fiat shamir transformation for verification, rather than
	//  regenerating the commitment: https://github.com/Layr-Labs/eigenda/issues/1037
	valid, err := verification.GenerateAndCompareBlobCommitment(pr.g1Srs, blob, kzgCommitment)
	if err != nil {
		return fmt.Errorf(
			"generate and compare commitment for blob %v received from relay %v: %w",
			blobKey.Hex(),
			relayKey,
			err)
	}

	if !valid {
		return fmt.Errorf("commitment for blob %v is invalid for bytes received from relay %v", blobKey.Hex(), relayKey)
	}

	// Checking that the length returned by the relay is <= the length claimed in the BlobCommitments is sufficient
	// here: it isn't necessary to verify the length proof itself, since this will have been done by DA nodes prior to
	// signing for availability.
	//
	// Note that the length in the commitment is the length of the blob in symbols
	if uint(len(blob)) > blobLength*encoding.BYTES_PER_SYMBOL {
		return fmt.Errorf(
			"length for blob %v (%d bytes) received from relay %v is greater than claimed blob length (%d bytes)",
			blobKey.Hex(),
			len(blob),
			relayKey,
			blobLength*encoding.BYTES_PER_SYMBOL)
	}

	return nil
}

// getBlobWithTimeout attempts to get a blob from a given relay, and times out based on config.RelayTimeout
func (pr *RelayPayloadRetriever) getBlobWithTimeout(
	ctx context.Context,
	relayKey core.RelayKey,
	blobKey *core.BlobKey) ([]byte, error) {

	timeoutCtx, cancel := context.WithTimeout(ctx, pr.config.RelayTimeout)
	defer cancel()

	return pr.relayClient.GetBlob(timeoutCtx, relayKey, *blobKey)
}

// Close is responsible for calling close on all internal clients. This method will do its best to close all internal
// clients, even if some closes fail.
//
// Any and all errors returned from closing internal clients will be joined and returned.
//
// This method should only be called once.
func (pr *RelayPayloadRetriever) Close() error {
	err := pr.relayClient.Close()
	if err != nil {
		return fmt.Errorf("close relay client: %w", err)
	}

	return nil
}

// blobCommitmentsBindingToInternal converts a blob commitment from an eigenDA cert into the internal
// encoding.BlobCommitments type
func blobCommitmentsBindingToInternal(
	blobCommitmentBinding *verifiercontract.BlobCommitment,
) (*encoding.BlobCommitments, error) {

	blobCommitment, err := encoding.BlobCommitmentsFromProtobuf(
		verification.BlobCommitmentBindingToProto(blobCommitmentBinding))

	if err != nil {
		return nil, fmt.Errorf("blob commitments from protobuf: %w", err)
	}

	return blobCommitment, nil
}
