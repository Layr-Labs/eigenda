package clients

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common/geth"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// RelayPayloadRetriever provides the ability to get payloads from the relay subsystem.
//
// This struct is goroutine safe.
type RelayPayloadRetriever struct {
	log logging.Logger
	// random doesn't need to be cryptographically secure, as it's only used to distribute load across relays.
	// Not all methods on Rand are guaranteed goroutine safe: if additional usages of random are added, they
	// must be evaluated for thread safety.
	random       *rand.Rand
	config       RelayPayloadRetrieverConfig
	codec        codecs.BlobCodec
	relayClient  RelayClient
	g1Srs        []bn254.G1Affine
	certVerifier verification.ICertVerifier
}

var _ PayloadRetriever = &RelayPayloadRetriever{}

// BuildRelayPayloadRetriever builds a RelayPayloadRetriever from config structs.
func BuildRelayPayloadRetriever(
	log logging.Logger,
	relayPayloadRetrieverConfig RelayPayloadRetrieverConfig,
	ethConfig geth.EthClientConfig,
	relayClientConfig *RelayClientConfig,
	g1Srs []bn254.G1Affine) (*RelayPayloadRetriever, error) {

	relayClient, err := NewRelayClient(relayClientConfig, log)
	if err != nil {
		return nil, fmt.Errorf("new relay client: %w", err)
	}

	ethClient, err := geth.NewClient(ethConfig, gethcommon.Address{}, 0, log)
	if err != nil {
		return nil, fmt.Errorf("new eth client: %w", err)
	}

	certVerifier, err := verification.NewCertVerifier(*ethClient, relayPayloadRetrieverConfig.EigenDACertVerifierAddr)
	if err != nil {
		return nil, fmt.Errorf("new cert verifier: %w", err)
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
		certVerifier,
		codec,
		g1Srs)
}

// NewRelayPayloadRetriever assembles a RelayPayloadRetriever from subcomponents that have already been constructed and initialized.
func NewRelayPayloadRetriever(
	log logging.Logger,
	random *rand.Rand,
	relayPayloadRetrieverConfig RelayPayloadRetrieverConfig,
	relayClient RelayClient,
	certVerifier verification.ICertVerifier,
	codec codecs.BlobCodec,
	g1Srs []bn254.G1Affine) (*RelayPayloadRetriever, error) {

	err := relayPayloadRetrieverConfig.checkAndSetDefaults()
	if err != nil {
		return nil, fmt.Errorf("check and set RelayRelayPayloadRetrieverConfig config: %w", err)
	}

	return &RelayPayloadRetriever{
		log:          log,
		random:       random,
		config:       relayPayloadRetrieverConfig,
		codec:        codec,
		relayClient:  relayClient,
		certVerifier: certVerifier,
		g1Srs:        g1Srs,
	}, nil
}

// GetPayload iteratively attempts to fetch a given blob with key blobKey from relays that have it, as claimed by the
// blob certificate. The relays are attempted in random order.
//
// If the blob is successfully retrieved, then the blob is verified. If the verification succeeds, the blob is decoded
// to yield the payload (the original user data, with no padding or any modification), and the payload is returned.
func (pr *RelayPayloadRetriever) GetPayload(
	ctx context.Context,
	eigenDACert *verification.EigenDACert) ([]byte, error) {

	blobKey, err := eigenDACert.ComputeBlobKey()
	if err != nil {
		return nil, fmt.Errorf("compute blob key: %w", err)
	}

	err = pr.verifyCertWithTimeout(ctx, eigenDACert)
	if err != nil {
		return nil, fmt.Errorf("verify cert with timeout for blobKey %v: %w", blobKey.Hex(), err)
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

		err = verification.CheckBlobLength(blob, blobCommitments.Length)
		if err != nil {
			pr.log.Warn("check blob length", "blobKey", blobKey.Hex(), "relayKey", relayKey, "error", err)
			continue
		}
		err = verification.VerifyBlobAgainstCert(blobKey, blob, blobCommitments.Commitment, pr.g1Srs)
		if err != nil {
			pr.log.Warn("verify blob against cert", "blobKey", blobKey.Hex(), "relayKey", relayKey, "error", err)
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

// getBlobWithTimeout attempts to get a blob from a given relay, and times out based on config.FetchTimeout
func (pr *RelayPayloadRetriever) getBlobWithTimeout(
	ctx context.Context,
	relayKey core.RelayKey,
	blobKey *core.BlobKey) ([]byte, error) {

	timeoutCtx, cancel := context.WithTimeout(ctx, pr.config.RelayTimeout)
	defer cancel()

	return pr.relayClient.GetBlob(timeoutCtx, relayKey, *blobKey)
}

// verifyCertWithTimeout verifies an EigenDACert by making a call to VerifyCertV2.
//
// This method times out after the duration configured in RelayPayloadRetrieverConfig.ContractCallTimeout
func (pr *RelayPayloadRetriever) verifyCertWithTimeout(
	ctx context.Context,
	eigenDACert *verification.EigenDACert,
) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, pr.config.ContractCallTimeout)
	defer cancel()

	return pr.certVerifier.VerifyCertV2(timeoutCtx, eigenDACert)
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
