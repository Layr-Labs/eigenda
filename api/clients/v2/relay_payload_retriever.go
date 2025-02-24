package clients

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/Layr-Labs/eigenda/api/clients/v2/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	core "github.com/Layr-Labs/eigenda/core/v2"
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

	err := relayPayloadRetrieverConfig.checkAndSetDefaults()
	if err != nil {
		return nil, fmt.Errorf("check and set RelayPayloadRetrieverConfig config: %w", err)
	}

	relayClient, err := NewRelayClient(relayClientConfig, log)
	if err != nil {
		return nil, fmt.Errorf("new relay client: %w", err)
	}

	return NewRelayPayloadRetriever(
		log,
		rand.New(rand.NewSource(rand.Int63())),
		relayPayloadRetrieverConfig,
		relayClient,
		g1Srs)
}

// NewRelayPayloadRetriever assembles a RelayPayloadRetriever from subcomponents that have already been constructed and
// initialized.
func NewRelayPayloadRetriever(
	log logging.Logger,
	random *rand.Rand,
	relayPayloadRetrieverConfig RelayPayloadRetrieverConfig,
	relayClient RelayClient,
	g1Srs []bn254.G1Affine) (*RelayPayloadRetriever, error) {

	err := relayPayloadRetrieverConfig.checkAndSetDefaults()
	if err != nil {
		return nil, fmt.Errorf("check and set RelayPayloadRetrieverConfig config: %w", err)
	}

	return &RelayPayloadRetriever{
		log:         log,
		random:      random,
		config:      relayPayloadRetrieverConfig,
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
	eigenDACert *verification.EigenDACert) (*codecs.Payload, error) {

	blobKey, err := eigenDACert.ComputeBlobKey()
	if err != nil {
		return nil, fmt.Errorf("compute blob key: %w", err)
	}

	relayKeys := eigenDACert.BlobInclusionInfo.BlobCertificate.RelayKeys
	relayKeyCount := len(relayKeys)
	if relayKeyCount == 0 {
		return nil, errors.New("relay key count is zero")
	}

	blobCommitments, err := verification.BlobCommitmentsBindingToInternal(
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

		blobLengthSymbols := eigenDACert.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Length

		blob, err := pr.retrieveBlobWithTimeout(ctx, relayKey, blobKey, blobLengthSymbols)
		// if GetBlob returned an error, try calling a different relay
		if err != nil {
			pr.log.Warn(
				"blob couldn't be retrieved from relay",
				"blobKey", blobKey.Hex(),
				"relayKey", relayKey,
				"error", err)
			continue
		}

		// TODO (litt3): eventually, we should make GenerateAndCompareBlobCommitment accept a blob, instead of the
		//  serialization of a blob. Commitment generation operates on field elements, which is how a blob is stored
		//  under the hood, so it's actually duplicating work to serialize the blob here. I'm declining to make this
		//  change now, to limit the size of the refactor PR.
		valid, err := verification.GenerateAndCompareBlobCommitment(
			pr.g1Srs,
			blob.Serialize(),
			blobCommitments.Commitment)
		if err != nil {
			pr.log.Warn(
				"generate and compare blob commitment",
				"blobKey", blobKey.Hex(), "relayKey", relayKey, "error", err)
			continue
		}
		if !valid {
			pr.log.Warn(
				"generated commitment doesn't match cert commitment",
				"blobKey", blobKey.Hex(), "relayKey", relayKey)
			continue
		}

		payload, err := blob.ToPayload(pr.config.PayloadPolynomialForm)
		if err != nil {
			pr.log.Error(
				`Commitment verification was successful, but conversion from blob to payload failed!
					This is likely a problem with the local configuration, but could potentially indicate
					malicious dispersed data. It should not be possible for a commitment to verify for an
					invalid blob!`,
				"blobKey", blobKey.Hex(), "relayKey", relayKey, "eigenDACert", eigenDACert, "error", err)
			return nil, fmt.Errorf("decode blob: %w", err)
		}

		return payload, nil
	}

	return nil, fmt.Errorf("unable to retrieve blob %v from any relay. relay count: %d", blobKey.Hex(), relayKeyCount)
}

// retrieveBlobWithTimeout attempts to retrieve a blob from a given relay, and times out based on config.FetchTimeout
func (pr *RelayPayloadRetriever) retrieveBlobWithTimeout(
	ctx context.Context,
	relayKey core.RelayKey,
	blobKey *core.BlobKey,
	blobLengthSymbols uint32,
) (*codecs.Blob, error) {

	timeoutCtx, cancel := context.WithTimeout(ctx, pr.config.RelayTimeout)
	defer cancel()

	blobBytes, err := pr.relayClient.GetBlob(timeoutCtx, relayKey, *blobKey)
	if err != nil {
		return nil, fmt.Errorf("get blob from relay: %w", err)
	}

	blob, err := codecs.DeserializeBlob(blobBytes, blobLengthSymbols)
	if err != nil {
		return nil, fmt.Errorf("deserialize blob: %w", err)
	}

	return blob, nil
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
