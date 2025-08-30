package altdacommitment_parser

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/ethereum/go-ethereum/rlp"
)

// CommitmentMode represents different types of EigenDA commitments
type CommitmentMode int

const (
	StandardMode CommitmentMode = iota
	OptimismGenericMode
	UnknownMode
)

func (c CommitmentMode) String() string {
	switch c {
	case StandardMode:
		return "Standard"
	case OptimismGenericMode:
		return "Optimism Generic"
	default:
		return "Unknown"
	}
}

// VersionByte represents EigenDA certificate versions
type VersionByte byte

const (
	V0VersionByte VersionByte = 0x00 // EigenDA V1
	V1VersionByte VersionByte = 0x01 // EigenDA V2 legacy
	V2VersionByte VersionByte = 0x02 // EigenDA V2 with V3 cert support
)

func (v VersionByte) String() string {
	switch v {
	case V0VersionByte:
		return "V0 (EigenDA V1)"
	case V1VersionByte:
		return "V1 (EigenDA V2 Legacy)"
	case V2VersionByte:
		return "V2 (EigenDA V2 with V3 Cert)"
	default:
		return fmt.Sprintf("Unknown (0x%02x)", byte(v))
	}
}

// ParsedCommitment holds the parsed commitment information
type ParsedCommitment struct {
	Mode            CommitmentMode
	CommitTypeByte  *byte
	DALayerByte     *byte
	VersionByte     *VersionByte
	CertificateData []byte
	OriginalData    []byte
}

// ParseResult holds the final parsing result
type ParseResult struct {
	Commitment *ParsedCommitment
	CertV2     *coretypes.EigenDACertV2
	CertV3     *coretypes.EigenDACertV3
	Version    string // "V2" or "V3"
}

// ParseCertFromHex parses an EigenDA certificate from a hex-encoded RLP string
func ParseCertFromHex(hexString string) (*ParseResult, error) {
	// Process the hex string to get binary data
	data, err := processHexString(hexString)
	if err != nil {
		return nil, fmt.Errorf("failed to process hex string: %w", err)
	}

	// Parse the commitment structure to extract prefix bytes and certificate data
	parsed, err := parseCommitmentStructure(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse commitment structure: %w", err)
	}

	// Try to parse the certificate data
	result := &ParseResult{
		Commitment: parsed,
	}

	err = parseCertificateData(parsed.CertificateData, parsed.VersionByte, result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate data: %w", err)
	}

	return result, nil
}

// ParseCertFromBytes parses an EigenDA certificate from binary data
func ParseCertFromBytes(data []byte) (*ParseResult, error) {
	// Determine if the input is hex-encoded or binary data
	processedData, err := processFileData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to process data: %w", err)
	}

	// Parse the commitment structure to extract prefix bytes and certificate data
	parsed, err := parseCommitmentStructure(processedData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse commitment structure: %w", err)
	}

	// Try to parse the certificate data
	result := &ParseResult{
		Commitment: parsed,
	}

	err = parseCertificateData(parsed.CertificateData, parsed.VersionByte, result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate data: %w", err)
	}

	return result, nil
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

// parseCommitmentStructure parses the commitment structure to identify mode and extract certificate data
func parseCommitmentStructure(data []byte) (*ParsedCommitment, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	// Step 1: Determine commitment mode
	mode, err := determineCommitmentMode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to determine commitment mode: %w", err)
	}

	// Step 2: Extract certificate data based on determined mode
	parsed, err := convertToCommitmentData(data, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to convert commitment data: %w", err)
	}

	return parsed, nil
}

// determineCommitmentMode uses RLP validation to distinguish between Standard and Optimism Generic modes
func determineCommitmentMode(data []byte) (CommitmentMode, error) {
	if len(data) == 0 {
		return UnknownMode, fmt.Errorf("empty data")
	}

	// First, try to parse as Standard mode: [version_byte][rlp_certificate]
	if len(data) > 1 {
		if isValidRLP(data[1:]) {
			return StandardMode, nil
		}
	}

	// If Standard mode RLP validation failed, check for Optimism Generic mode
	// Optimism Generic: [0x01][da_layer_byte][version_byte][rlp_certificate]
	if len(data) >= 3 && isValidRLP(data[3:]) {
		return OptimismGenericMode, nil
	} else {
		// If we can't determine the mode conclusively
		return UnknownMode, fmt.Errorf("cannot determine commitment mode")
	}
}

// convertToCommitmentData extracts the commitment data based on the determined mode
func convertToCommitmentData(data []byte, mode CommitmentMode) (*ParsedCommitment, error) {
	parsed := &ParsedCommitment{
		OriginalData: data,
		Mode:         mode,
	}

	switch mode {
	case StandardMode:
		// Standard mode: [version_byte][rlp_certificate]
		versionByte := VersionByte(data[0])
		parsed.VersionByte = &versionByte
		parsed.CertificateData = data[1:]

	case OptimismGenericMode:
		// Optimism Generic mode: [0x01][da_layer_byte][version_byte][rlp_certificate]
		if len(data) < 3 {
			return nil, fmt.Errorf("insufficient data for Optimism Generic mode: need at least 3 bytes, got %d", len(data))
		}
		commitTypeByte := data[0]
		daLayerByte := data[1]
		versionByte := VersionByte(data[2])
		parsed.CommitTypeByte = &commitTypeByte
		parsed.DALayerByte = &daLayerByte
		parsed.VersionByte = &versionByte
		parsed.CertificateData = data[3:]

	default:
		return nil, fmt.Errorf("unsupported commitment mode: %v", mode)
	}

	return parsed, nil
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

// parseCertificateData attempts to parse the certificate data as V2 or V3
func parseCertificateData(data []byte, versionByte *VersionByte, result *ParseResult) error {
	if len(data) == 0 {
		return fmt.Errorf("no certificate data to parse")
	}

	// Try to parse as V3 cert first
	var certV3 coretypes.EigenDACertV3
	err := rlp.DecodeBytes(data, &certV3)
	if err == nil {
		result.CertV3 = &certV3
		result.Version = "V3"
		return nil
	}

	// Try to parse as V2 cert
	var certV2 coretypes.EigenDACertV2
	err = rlp.DecodeBytes(data, &certV2)
	if err == nil {
		result.CertV2 = &certV2
		result.Version = "V2"
		return nil
	}

	return fmt.Errorf("failed to parse certificate as V2 or V3: %w", err)
}

// processFileData determines if the input is hex-encoded or binary data
func processFileData(rawData []byte) ([]byte, error) {
	// Check if the data looks like hex by examining the first few bytes
	if isHexData(rawData) {
		// Remove common hex prefixes and whitespace
		hexStr := strings.TrimSpace(string(rawData))
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

	// Assume it's already binary RLP data
	return rawData, nil
}

// isHexData attempts to determine if the input data is hex-encoded
func isHexData(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	// Convert to string for easier analysis
	str := string(data)

	// Check for hex prefix
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		return true
	}

	// Check if the majority of characters are valid hex characters
	hexChars := 0
	validChars := 0
	for _, b := range data {
		if (b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F') {
			hexChars++
		}
		if (b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F') ||
			b == ' ' || b == '\n' || b == '\r' || b == '\t' {
			validChars++
		}
	}

	// If most characters are hex and the valid character ratio is high, likely hex
	if len(data) > 10 && float64(hexChars)/float64(len(data)) > 0.7 &&
		float64(validChars)/float64(len(data)) > 0.9 {
		return true
	}

	return false
}
