package clients

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/hashicorp/go-multierror"
	"math/rand"
)

// EigenDAClientV2 provides the ability to get blobs from the relay subsystem, and to send new blobs to the disperser.
type EigenDAClientV2 interface {
	// GetBlob iteratively attempts to retrieve a given blob with key blobKey from the relays listed in the blobCertificate.
	GetBlob(ctx context.Context, blobKey corev2.BlobKey, blobCertificate corev2.BlobCertificate) ([]byte, error)
	GetCodec() codecs.BlobCodec
	Close() error
}

type eigenDAClientV2 struct {
	clientConfig EigenDAClientConfig
	log          logging.Logger
	relayClient  RelayClient
	codec        codecs.BlobCodec
}

var _ EigenDAClientV2 = &eigenDAClientV2{}

// BuildEigenDAClientV2 builds an EigenDAClientV2 from config structs.
func BuildEigenDAClientV2(
	log logging.Logger,
	clientConfig EigenDAClientConfig,
	relayClientConfig RelayClientConfig) (EigenDAClientV2, error) {

	err := clientConfig.CheckAndSetDefaults()
	if err != nil {
		return nil, err
	}

	relayClient, err := NewRelayClient(&relayClientConfig, log)

	lowLevelCodec, err := codecs.BlobEncodingVersionToCodec(clientConfig.PutBlobEncodingVersion)
	if err != nil {
		return nil, fmt.Errorf("create low level codec: %w", err)
	}
	var codec codecs.BlobCodec
	if clientConfig.DisablePointVerificationMode {
		codec = codecs.NewNoIFFTCodec(lowLevelCodec)
	} else {
		codec = codecs.NewIFFTCodec(lowLevelCodec)
	}

	return NewEigenDAClientV2(log, clientConfig, relayClient, codec)
}

// NewEigenDAClientV2 assembles an EigenDAClientV2 from subcomponents that have already been constructed and initialized.
func NewEigenDAClientV2(
	log logging.Logger,
	clientConfig EigenDAClientConfig,
	relayClient RelayClient,
	codec codecs.BlobCodec) (EigenDAClientV2, error) {

	return &eigenDAClientV2{
		clientConfig: clientConfig,
		log:          log,
		relayClient:  relayClient,
		codec:        codec,
	}, nil
}

// GetBlob iteratively attempts to retrieve a given blob with key blobKey from the relays listed in the blobCertificate.
//
// The relays are attempted in random order.
//
// The returned blob is decoded.
func (c *eigenDAClientV2) GetBlob(ctx context.Context, blobKey corev2.BlobKey, blobCertificate corev2.BlobCertificate) ([]byte, error) {
	// create a randomized array of indices, so that it isn't always the first relay in the list which gets hit
	random := rand.New(rand.NewSource(rand.Int63()))
	relayKeyCount := len(blobCertificate.RelayKeys)
	var indices []int
	for i := 0; i < relayKeyCount; i++ {
		indices = append(indices, i)
	}
	random.Shuffle(len(indices), func(i int, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	// TODO (litt3): consider creating a utility which can deprioritize relays that fail to respond (or respond maliciously)

	// iterate over relays in random order, until we are able to get the blob from someone
	for _, val := range indices {
		relayKey := blobCertificate.RelayKeys[val]

		data, err := c.relayClient.GetBlob(ctx, relayKey, blobKey)

		// if GetBlob returned an error, try calling a different relay
		if err != nil {
			// TODO: should this log type be downgraded to debug to avoid log spam? I'm not sure how frequent retrieval
			//  from a relay will fail in practice (?)
			c.log.Info("blob couldn't be retrieved from relay", "blobKey", blobKey, "relayKey", relayKey)
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
			c.log.Warn("error decoding blob", "blobKey", blobKey, "relayKey", relayKey)
			continue
		}

		return decodedData, nil
	}

	return nil, fmt.Errorf("unable to retrieve blob from any relay")
}

func (c *eigenDAClientV2) GetCodec() codecs.BlobCodec {
	return c.codec
}

func (c *eigenDAClientV2) Close() error {
	var errList *multierror.Error

	// TODO: this is using a multierror, since there will be more subcomponents requiring closing after adding PUT functionality
	relayClientErr := c.relayClient.Close()
	if relayClientErr != nil {
		errList = multierror.Append(errList, relayClientErr)
	}

	if errList != nil {
		return errList.ErrorOrNil()
	}

	return nil
}
