package v2

import (
	"context"
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"math/rand"
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
	relayClient clients.RelayClient
}

// BuildEigenDAClient builds an EigenDAClient from config structs.
func BuildEigenDAClient(
	log logging.Logger,
	config *EigenDAClientConfig,
	relayClientConfig *clients.RelayClientConfig) (*EigenDAClient, error) {

	relayClient, err := clients.NewRelayClient(relayClientConfig, log)
	if err != nil {
		return nil, fmt.Errorf("new relay client: %w", err)
	}

	codec, err := createCodec(config)
	if err != nil {
		return nil, err
	}

	return NewEigenDAClient(log, rand.New(rand.NewSource(rand.Int63())), config, relayClient, codec)
}

// NewEigenDAClient assembles an EigenDAClient from subcomponents that have already been constructed and initialized.
func NewEigenDAClient(
	log logging.Logger,
	random *rand.Rand,
	config *EigenDAClientConfig,
	relayClient clients.RelayClient,
	codec codecs.BlobCodec) (*EigenDAClient, error) {

	return &EigenDAClient{
		log:         log,
		random:      random,
		config:      config,
		codec:       codec,
		relayClient: relayClient,
	}, nil
}

// GetBlob iteratively attempts to retrieve a given blob with key blobKey from the relays listed in the blobCertificate.
//
// The relays are attempted in random order.
//
// The returned blob is decoded.
func (c *EigenDAClient) GetBlob(
	ctx context.Context,
	blobKey core.BlobKey,
	blobCertificate core.BlobCertificate) ([]byte, error) {

	relayKeyCount := len(blobCertificate.RelayKeys)

	if relayKeyCount == 0 {
		return nil, errors.New("relay key count is zero")
	}

	// create a randomized array of indices, so that it isn't always the first relay in the list which gets hit
	indices := c.random.Perm(relayKeyCount)

	// TODO (litt3): consider creating a utility which can deprioritize relays that fail to respond (or respond maliciously)

	// iterate over relays in random order, until we are able to get the blob from someone
	for _, val := range indices {
		relayKey := blobCertificate.RelayKeys[val]

		data, err := c.getBlobWithTimeout(ctx, relayKey, blobKey)

		// if GetBlob returned an error, try calling a different relay
		if err != nil {
			c.log.Warn("blob couldn't be retrieved from relay", "blobKey", blobKey, "relayKey", relayKey, "error", err)
			continue
		}

		// An honest relay should never send an empty blob
		if len(data) == 0 {
			c.log.Warn("blob received from relay had length 0", "blobKey", blobKey, "relayKey", relayKey)
			continue
		}

		// An honest relay should never send a blob which cannot be decoded
		decodedData, err := c.codec.DecodeBlob(data)
		if err != nil {
			c.log.Warn("error decoding blob from relay", "blobKey", blobKey, "relayKey", relayKey, "error", err)
			continue
		}

		return decodedData, nil
	}

	return nil, fmt.Errorf("unable to retrieve blob from any relay. relay count: %d", relayKeyCount)
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

	switch config.PointVerificationMode {
	case NoIFFT:
		return codecs.NewNoIFFTCodec(lowLevelCodec), nil
	case IFFT:
		return codecs.NewIFFTCodec(lowLevelCodec), nil
	default:
		return nil, fmt.Errorf("unsupported point verification mode: %d", config.PointVerificationMode)
	}
}
