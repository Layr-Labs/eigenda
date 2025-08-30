package altdacommitment_parser

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/ethereum/go-ethereum/rlp"
)

// PrefixMetadata holds the parsed commitment information
type PrefixMetadata struct {
	Mode           commitments.CommitmentMode
	CommitTypeByte *byte
	DALayerByte    *byte
	CertVersion    certs.VersionByte
	OriginalSize   int
}

// ParseCertFromHex parses an EigenDA certificate from a hex-encoded RLP string
func ParseCertFromHex(hexString string) (*PrefixMetadata, *certs.VersionedCert, error) {
	// Process the hex string to get binary data
	data, err := processHexString(hexString)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to process hex string: %w", err)
	}

	if len(data) == 0 {
		return nil, nil, fmt.Errorf("empty data")
	}

	// Step 1: Determine commitment mode
	mode, err := determineCommitmentMode(data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to determine commitment mode: %w", err)
	}

	// parse cert
	var versionedCert certs.VersionedCert
	var prefix PrefixMetadata
	prefix.Mode = mode
	prefix.OriginalSize = len(data)
	switch mode {
	case commitments.StandardCommitmentMode:
		// Standard mode: [version_byte][rlp_certificate]
		versionByte := certs.VersionByte(data[0])
		prefix.CertVersion = versionByte

		versionedCert = certs.NewVersionedCert(data[1:], versionByte)

	case commitments.OptimismGenericCommitmentMode:
		// Optimism Generic mode: [0x01][da_layer_byte][version_byte][rlp_certificate]
		if len(data) < 3 {
			return nil, nil, fmt.Errorf("insufficient data for Optimism Generic mode: need at least 3 bytes, got %d", len(data))
		}
		prefix.CommitTypeByte = &data[0]
		prefix.DALayerByte = &data[1]
		versionByte := certs.VersionByte(data[2])
		prefix.CertVersion = versionByte

		versionedCert = certs.NewVersionedCert(data[3:], versionByte)

	case commitments.OptimismKeccakCommitmentMode:
		// Optimism Keccak mode is not expected in this parser context but included for exhaustiveness
		return nil, nil, fmt.Errorf("OptimismKeccakCommitmentMode is not supported by this parser")

	default:
		return nil, nil, fmt.Errorf("unsupported commitment mode: %v", mode)
	}

	return &prefix, &versionedCert, nil
}

// processHexString processes a hex-encoded string and returns binary data for RLP decoding
func processHexString(hexString string) ([]byte, error) {
	// Remove common hex prefixes and whitespace
	hexStr := strings.TrimSpace(hexString)
	hexStr = strings.TrimPrefix(hexStr, "0x")
	hexStr = strings.TrimPrefix(hexStr, "0X")

	// Remove any whitespace, newlines, and other non-hex characters
	hexStr = strings.ReplaceAll(hexStr, " ", "")
	hexStr = strings.ReplaceAll(hexStr, "\n", "")
	hexStr = strings.ReplaceAll(hexStr, "\r", "")
	hexStr = strings.ReplaceAll(hexStr, "\t", "")

	// Decode hex string to binary data
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}

	return data, nil
}

// determineCommitmentMode uses RLP validation to distinguish between Standard and Optimism Generic modes
func determineCommitmentMode(data []byte) (commitments.CommitmentMode, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("empty data")
	}

	// First, try to parse as Standard mode: [version_byte][rlp_certificate]
	if len(data) > 1 {
		if isValidRLP(data[1:]) {
			return commitments.StandardCommitmentMode, nil
		}
	}

	// If Standard mode RLP validation failed, check for Optimism Generic mode
	// Optimism Generic: [0x01][da_layer_byte][version_byte][rlp_certificate]
	if len(data) >= 3 && isValidRLP(data[3:]) {
		return commitments.OptimismGenericCommitmentMode, nil
	} else {
		// If we can't determine the mode conclusively
		return "", fmt.Errorf("cannot determine commitment mode")
	}
}

// isValidRLP attempts to validate if the data is valid RLP encoding
func isValidRLP(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	// Try to decode as both V2 and V3 certificates to validate RLP structure
	var certV3 coretypes.EigenDACertV3
	if err := rlp.DecodeBytes(data, &certV3); err == nil {
		return true
	}

	var certV2 coretypes.EigenDACertV2
	if err := rlp.DecodeBytes(data, &certV2); err == nil {
		return true
	}

	return false
}

// ParseCertificateData attempts to parse the certificate data as V2 or V3 and returns V3
func ParseCertificateData(cert *certs.VersionedCert) (*coretypes.EigenDACertV3, error) {
	if len(cert.SerializedCert) == 0 {
		return nil, fmt.Errorf("no certificate data to parse")
	}

	// Try to parse as V3 cert first
	var certV3 coretypes.EigenDACertV3
	err := rlp.DecodeBytes(cert.SerializedCert, &certV3)
	if err == nil {
		return &certV3, nil
	}

	// Try to parse as V2 cert and convert to V3
	var certV2 coretypes.EigenDACertV2
	err = rlp.DecodeBytes(cert.SerializedCert, &certV2)
	if err == nil {
		// Convert V2 to V3
		return certV2.ToV3(), nil
	}

	return nil, fmt.Errorf("failed to parse certificate as V2 or V3: %w", err)
}
