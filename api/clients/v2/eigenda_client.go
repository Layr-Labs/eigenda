package clients

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
)

// EigenDAClient provides the ability to get blobs from the relay subsystem, and to send new blobs to the disperser.
//
// This struct is not threadsafe.
type EigenDAClient struct {
	log logging.Logger
	// doesn't need to be cryptographically secure, as it's only used to distribute load across relays
	random      *rand.Rand
	config      *EigenDAClientConfig
	codec       codecs.BlobCodec
	relayClient RelayClient
	g1Srs       []bn254.G1Affine
}

// BuildEigenDAClient builds an EigenDAClient from config structs.
func BuildEigenDAClient(
	log logging.Logger,
	config *EigenDAClientConfig,
	relayClientConfig *RelayClientConfig,
	g1Srs []bn254.G1Affine) (*EigenDAClient, error) {

	relayClient, err := NewRelayClient(relayClientConfig, log)
	if err != nil {
		return nil, fmt.Errorf("new relay client: %w", err)
	}

	codec, err := createCodec(config)
	if err != nil {
		return nil, err
	}

	return NewEigenDAClient(log, rand.New(rand.NewSource(rand.Int63())), config, relayClient, codec, g1Srs)
}

// NewEigenDAClient assembles an EigenDAClient from subcomponents that have already been constructed and initialized.
func NewEigenDAClient(
	log logging.Logger,
	random *rand.Rand,
	config *EigenDAClientConfig,
	relayClient RelayClient,
	codec codecs.BlobCodec,
	g1Srs []bn254.G1Affine) (*EigenDAClient, error) {

	return &EigenDAClient{
		log:         log,
		random:      random,
		config:      config,
		codec:       codec,
		relayClient: relayClient,
		g1Srs:       g1Srs,
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
	blobCertificate *core.BlobCertificate) ([]byte, error) {

	relayKeys := blobCertificate.RelayKeys
	relayKeyCount := len(relayKeys)

	if relayKeyCount == 0 {
		return nil, errors.New("relay key count is zero")
	}

	// create a randomized array of indices, so that it isn't always the first relay in the list which gets hit
	indices := c.random.Perm(relayKeyCount)

	// TODO (litt3): consider creating a utility which can deprioritize relays that fail to respond (or respond maliciously)

	// iterate over relays in random order, until we are able to get the blob from someone
	for _, val := range indices {
		relayKey := relayKeys[val]

		blob, err := c.getBlobWithTimeout(ctx, relayKey, blobKey)

		// if GetBlob returned an error, try calling a different relay
		if err != nil {
			c.log.Warn("blob couldn't be retrieved from relay", "blobKey", blobKey, "relayKey", relayKey, "error", err)
			continue
		}

		// An honest relay should never send a blob which doesn't verify
		if !c.verifyBlobFromRelay(blobKey, relayKey, blob, blobCertificate) {
			// specifics are logged in verifyBlobFromRelay
			continue
		}

		// An honest relay should never send a blob which cannot be decoded into a payload
		payload, err := c.codec.DecodeBlob(blob)
		if err != nil {
			c.log.Warn("error decoding blob from relay", "blobKey", blobKey, "relayKey", relayKey, "error", err)
			continue
		}

		return payload, nil
	}

	return nil, fmt.Errorf("unable to retrieve blob %v from any relay. relay count: %d", blobKey, relayKeyCount)
}

// verifyBlobFromRelay performs the necessary verification after having retrieved a blob from a relay.
//
// If all verifications succeed, the method returns true. Otherwise, it logs a warning and returns false.
func (c *EigenDAClient) verifyBlobFromRelay(
	blobKey core.BlobKey,
	relayKey core.RelayKey,
	blob []byte,
	blobCertificate *core.BlobCertificate) bool {

	// An honest relay should never send an empty blob
	if len(blob) == 0 {
		c.log.Warn("blob received from relay had length 0", "blobKey", blobKey, "relayKey", relayKey)
		return false
	}

	blobCommitments := blobCertificate.BlobHeader.BlobCommitments

	// TODO: in the future, this will be optimized to use fiat shamir transformation for verification, rather than
	//  regenerating the commitment: https://github.com/Layr-Labs/eigenda/issues/1037
	valid, err := verification.GenerateAndCompareBlobCommitment(
		c.g1Srs,
		blob,
		blobCommitments.Commitment)
	if err != nil {
		c.log.Warn(
			"error generating commitment from received blob",
			"blobKey",
			blobKey,
			"relayKey",
			relayKey,
			"error",
			err)
		return false
	}

	if !valid {
		c.log.Warn(
			"blob commitment is invalid for received bytes",
			"blobKey",
			blobKey,
			"relayKey",
			relayKey)
		return false
	}

	// checking that the length returned by the relay matches the claimed length is sufficient here: it isn't necessary
	// to verify the length proof itself, since this will have been done by DA nodes prior to signing for availability.
	if uint(len(blob)) != blobCommitments.Length {
		c.log.Warn(
			"blob length doesn't match length claimed in blob commitments",
			"blobKey",
			blobKey,
			"relayKey",
			relayKey,
			"blobLength",
			len(blob),
			"blobCommitments.Length",
			blobCommitments.Length)
		return false
	}

	// TODO: call verifyBlobV2

	return true
}

// getBlobWithTimeout attempts to get a blob from a given relay, and times out based on config.RelayTimeout
func (c *EigenDAClient) getBlobWithTimeout(
	ctx context.Context,
	relayKey core.RelayKey,
	blobKey core.BlobKey) ([]byte, error) {

	timeoutCtx, cancel := context.WithTimeout(ctx, c.config.RelayTimeout)
	defer cancel()

	return c.relayClient.GetBlob(timeoutCtx, relayKey, blobKey)
}

// GetCodec returns the codec the client uses for encoding and decoding blobs
func (c *EigenDAClient) GetCodec() codecs.BlobCodec {
	return c.codec
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

	switch config.BlobPolynomialForm {
	case codecs.Eval:
		// a blob polynomial is already in Eval form after being encoded. Therefore, we use the NoIFFTCodec, which
		// doesn't do any further conversion.
		return codecs.NewNoIFFTCodec(lowLevelCodec), nil
	case codecs.Coeff:
		// a blob polynomial starts in Eval form after being encoded. Therefore, we use the IFFT codec to transform
		// the blob into Coeff form after initial encoding. This codec also transforms the Coeff form received from the
		// relay back into Eval form when decoding.
		return codecs.NewIFFTCodec(lowLevelCodec), nil
	default:
		return nil, fmt.Errorf("unsupported polynomial form: %d", config.BlobPolynomialForm)
	}
}
