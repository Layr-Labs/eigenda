package codecs

import (
	"fmt"
)

type PayloadEncodingVersion uint8

const (
	// PayloadEncodingVersion0 entails a 32 byte header = [0x00, version byte, big-endian uint32 len of payload, 0x00, 0x00,...]
	// followed by the encoded data [0x00, 31 bytes of data, 0x00, 31 bytes of data,...]
	//
	// Each group of 32 bytes starts with a 0x00 byte so that they can be parsed as valid bn254 field elements.
	PayloadEncodingVersion0 PayloadEncodingVersion = 0x0
)

type BlobCodec interface {
	DecodeBlob(encodedData []byte) ([]byte, error)
	EncodeBlob(rawData []byte) ([]byte, error)
}

func BlobEncodingVersionToCodec(version PayloadEncodingVersion) (BlobCodec, error) {
	switch version {
	case PayloadEncodingVersion0:
		return DefaultBlobCodec{}, nil
	default:
		return nil, fmt.Errorf("unsupported blob encoding version: %x", version)
	}
}

func GenericDecodeBlob(data []byte) ([]byte, error) {
	if len(data) <= 32 {
		return nil, fmt.Errorf("data is not of length greater than 32 bytes: %d", len(data))
	}
	// version byte is stored in [1], because [0] is always 0 to ensure the codecBlobHeader is a valid bn254 element
	// see https://github.com/Layr-Labs/eigenda/blob/master/api/clients/codecs/default_blob_codec.go#L21
	// TODO: we should prob be working over a struct with methods such as GetBlobEncodingVersion() to prevent index errors
	version := PayloadEncodingVersion(data[1])
	codec, err := BlobEncodingVersionToCodec(version)
	if err != nil {
		return nil, err
	}

	data, err = codec.DecodeBlob(data)
	if err != nil {
		return nil, fmt.Errorf("unable to decode blob: %w", err)
	}

	return data, nil
}
