package integration_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	disperserpb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/wealdtech/go-merkletree/v2"
	"github.com/wealdtech/go-merkletree/v2/keccak256"
)

// FuzzInaboxIntegration tests the Inabox v2 integration with randomized inputs.
func FuzzInaboxIntegration(f *testing.F) {
	// Define seed inputs with edge cases for data1
	seedInputs := [][]byte{
		{},                              // 0 bytes
		make([]byte, 1),                 // Minimal non-empty input
		bytes.Repeat([]byte{0x00}, 100), // All zeros
		bytes.Repeat([]byte{0xFF}, 100), // All 0xFF bytes
		[]byte("Hello, 世界"),           // Unicode characters
		make([]byte, 1024),              // 1KB of zeros
		func() []byte { // Random 4MB
			data := make([]byte, 4*1024*1024) // 4MB
			_, err := rand.Read(data)
			if err != nil {
				panic(fmt.Sprintf("Failed to generate random seed data: %v", err))
			}
			return data
		}(),
	}

	// Add seed inputs to the fuzzer
	for _, seed := range seedInputs {
		f.Add(seed)
	}

	// Define the fuzzing function
	f.Fuzz(func(t *testing.T, data1 []byte) {
		// Limit data1 to a maximum of 4MB
		const maxSize = 4 * 1024 * 1024 // 4MB
		if len(data1) > maxSize {
			data1 = data1[:maxSize]
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Initialize the signer with the provided private key
		privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded"
		signer := auth.NewLocalBlobRequestSigner(privateKeyHex)

		// Create the Disperser client
		disp, err := clients.NewDisperserClient(&clients.DisperserClientConfig{
			Hostname: "localhost",
			Port:     "32005",
		}, signer, nil, nil)
		if err != nil {
			t.Errorf("Failed to create Disperser client: %v", err)
			return
		}
		if disp == nil {
			t.Errorf("Disperser client is nil")
			return
		}

		// Pad data1 as required by the codec
		paddedData1 := codec.ConvertByPaddingEmptyByte(data1)

		// Disperse the blob
		blobStatus, key, err := disp.DisperseBlob(ctx, paddedData1, 0, []uint8{0, 1}, 0)
		if err != nil {
			t.Errorf("DisperseBlob for data1 failed: %v", err)
			return
		}
		if key.Hex() == "" {
			t.Errorf("DisperseBlob for data1 returned nil key")
			return
		}
		if blobStatus == nil {
			t.Errorf("DisperseBlob for data1 returned nil status")
			return
		}
		if *blobStatus != dispv2.Queued {
			t.Errorf("Expected blobStatus to be Queued, got %v", *blobStatus)
			return
		}

		// Poll for certification
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		var reply *disperserpb.BlobStatusReply
		var blobCert *corev2.BlobCertificate
		var batchHeader *commonpb.BatchHeader

		for loop := true; loop; {
			select {
			case <-ctx.Done():
				t.Error("Test timed out while waiting for blob certification")
				return
			case <-ticker.C:
				reply, err = disp.GetBlobStatus(context.Background(), key)
				if err != nil {
					t.Errorf("GetBlobStatus failed: %v", err)
					return
				}
				if reply == nil {
					t.Error("GetBlobStatus returned nil reply")
					return
				}
				status, err := dispv2.BlobStatusFromProtobuf(reply.GetStatus())
				if err != nil {
					t.Errorf("BlobStatusFromProtobuf failed: %v", err)
					return
				}

				if status != dispv2.Certified {
					// Blob not yet certified; continue polling
					continue
				}

				// Process the certified blob
				batchHeader = reply.GetSignedBatch().GetHeader()
				if batchHeader == nil {
					t.Error("batchHeader is nil")
					return
				}
				if batchHeader.GetBatchRoot() == nil {
					t.Error("batchHeader.BatchRoot is nil")
					return
				}
				if batchHeader.GetReferenceBlockNumber() <= 0 {
					t.Errorf("batchHeader.ReferenceBlockNumber expected > 0, got %d", batchHeader.GetReferenceBlockNumber())
					return
				}

				attestation := reply.GetSignedBatch().GetAttestation()
				if attestation == nil {
					t.Error("attestation is nil")
					return
				}
				if !bytes.Equal(attestation.QuorumSignedPercentages, []byte{100, 100}) {
					t.Errorf("Expected quorumSignedPercentages to be [100, 100], got %v", attestation.QuorumSignedPercentages)
					return
				}

				blobVerification := reply.GetBlobVerificationInfo()
				if blobVerification == nil {
					t.Error("blobVerification is nil")
					return
				}
				if blobVerification.GetBlobCertificate() == nil {
					t.Error("blobVerification.BlobCertificate is nil")
					return
				}

				blobCert, err = corev2.BlobCertificateFromProtobuf(blobVerification.GetBlobCertificate())
				if err != nil {
					t.Errorf("BlobCertificateFromProtobuf failed: %v", err)
					return
				}

				inclusionProofBytes := blobVerification.GetInclusionProof()
				blobIndex := blobVerification.GetBlobIndex()
				proof, err := core.DeserializeMerkleProof(inclusionProofBytes, uint64(blobIndex))
				if err != nil {
					t.Errorf("DeserializeMerkleProof failed: %v", err)
					return
				}

				certHash, err := blobCert.Hash()
				if err != nil {
					t.Errorf("BlobCertificate Hash failed: %v", err)
					return
				}

				verified, err := merkletree.VerifyProofUsing(certHash[:], false, proof, [][]byte{batchHeader.BatchRoot}, keccak256.New())
				if err != nil {
					t.Errorf("VerifyProofUsing failed: %v", err)
					return
				}
				if !verified {
					t.Error("Merkle proof verification failed")
					return
				}

				// Blob is certified and verified; exit the loop
				loop = false
			}
		}

		// Test retrieval from relay
		relayClient, err := clients.NewRelayClient(&clients.RelayClientConfig{
			Sockets: relays,
		}, logger)
		if err != nil {
			t.Errorf("Failed to create Relay client: %v", err)
			return
		}

		blobRelays := make(map[corev2.RelayKey]struct{}, len(blobCert.RelayKeys))
		for _, k := range blobCert.RelayKeys {
			blobRelays[corev2.RelayKey(k)] = struct{}{}
		}

		for relayKey := range relays { // Ensure 'relays' is defined
			blob, err := relayClient.GetBlob(ctx, relayKey, key)
			if _, ok := blobRelays[corev2.RelayKey(relayKey)]; ok {
				if err != nil {
					t.Errorf("GetBlob from relay %v failed: %v", relayKey, err)
					continue
				}
				if !bytes.Equal(blob, paddedData1) {
					t.Errorf("Retrieved blob data mismatch from relay %v", relayKey)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error when retrieving blob from non-associated relay %v", relayKey)
				}
			}
		}

		// Test retrieval from retrievalClientV2
		// Ensure 'retrievalClientV2' is defined or initialized appropriately
		b, err := retrievalClientV2.GetBlob(ctx, blobCert.BlobHeader, batchHeader.ReferenceBlockNumber, 0)
		if err != nil {
			t.Errorf("retrievalClientV2.GetBlob failed: %v", err)
		} else {
			restored := bytes.TrimRight(b, "\x00")
			if !bytes.Equal(restored, paddedData1) {
				t.Errorf("Restored blob data does not match original data")
			}
		}

		// additional retrieval checks if necessary
	})
}
