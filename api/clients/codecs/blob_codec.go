package codecs

import (
	"fmt"
)

type BlobEncodingVersion byte

const (
	// This minimal blob encoding includes a version byte, a length uint32, and 31 byte field element mapping. It does not include IFFT padding + IFFT.
	DefaultBlobEncoding BlobEncodingVersion = 0x0
)

type BlobCodec interface {
	DecodeBlob(encodedData []byte) ([]byte, error)
	EncodeBlob(rawData []byte) ([]byte, error)
}

func BlobEncodingVersionToCodec(version BlobEncodingVersion) (BlobCodec, error) {
	switch version {
	case DefaultBlobEncoding:
		return DefaultBlobEncodingCodec{}, nil
	default:
		return nil, fmt.Errorf("unsupported blob encoding version: %x", version)
	}
}
