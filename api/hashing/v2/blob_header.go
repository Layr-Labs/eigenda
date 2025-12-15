package hashing

import (
	"fmt"
	"time"

	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/api/hashing/v2/serialize"
	"golang.org/x/crypto/sha3"
)

// BlobHeaderHashWithTimestamp is a tuple of a blob header hash and the timestamp of the blob header.
type BlobHeaderHashWithTimestamp struct {
	Hash      []byte
	Timestamp time.Time
}

// BlobHeadersHashesAndTimestamps returns a list of per-BlobHeader hashes (one per BlobCertificate)
// with the timestamp.
func BlobHeadersHashesAndTimestamps(request *grpc.StoreChunksRequest) ([]BlobHeaderHashWithTimestamp, error) {
	certs := request.GetBatch().GetBlobCertificates()
	out := make([]BlobHeaderHashWithTimestamp, len(certs))
	for i, cert := range certs {
		if cert == nil {
			return nil, fmt.Errorf("nil BlobCertificate at index %d", i)
		}
		header := cert.GetBlobHeader()
		if header == nil {
			return nil, fmt.Errorf("nil BlobHeader at index %d", i)
		}
		paymentHeader := header.GetPaymentHeader()
		if paymentHeader == nil {
			return nil, fmt.Errorf("nil PaymentHeader at index %d", i)
		}

		headerBytes, err := serialize.SerializeBlobHeader(header)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize blob header at index %d: %w", i, err)
		}
		hasher := sha3.New256()
		_, _ = hasher.Write(headerBytes)
		out[i] = BlobHeaderHashWithTimestamp{
			Hash: hasher.Sum(nil),
			Timestamp: time.Unix(
				paymentHeader.GetTimestamp()/int64(time.Second),
				paymentHeader.GetTimestamp()%int64(time.Second),
			),
		}
	}

	return out, nil
}
