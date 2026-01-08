package altdacommitment_parser

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda/tools/integration_utils/flags"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/urfave/cli"
)

// PrefixMetadata holds the parsed prefix information
type PrefixMetadata struct {
	Mode           commitments.CommitmentMode
	CommitTypeByte *byte
	DALayerByte    *byte
	CertVersion    certs.VersionByte
	OriginalSize   int
}

// DisplayAltDACommitmentFromHex parses an EigenDA AltDA commitment from a hex-encoded RLP string
// and prints a nicely formatted display of its contents to stdout
func DisplayAltDACommitmentFromHex(ctx *cli.Context) error {
	hexString := ctx.String(flags.CertHexFlag.Name)

	// Use the parser library to parse the certificate
	prefix, versionedCert, err := ParseAltDACommitmentFromHex(hexString)
	if err != nil {
		return fmt.Errorf("failed to parse cert prefix: %w", err)
	}

	// Display the parsed prefix information
	DisplayPrefixInfo(prefix)

	// Display the certificate data (handles V2, V3, and V4)
	if err := DisplayCertData(versionedCert.SerializedCert); err != nil {
		return fmt.Errorf("failed to display certificate data: %w", err)
	}
	return nil
}

// ParseAltDACommitmentFromHex parses an prefix and certificate from a hex-encoded RLP string
func ParseAltDACommitmentFromHex(hexString string) (*PrefixMetadata, *certs.VersionedCert, error) {
	// Process the hex string to get binary data
	data, err := ProcessHexString(hexString)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to process hex string: %w", err)
	}
	if len(data) == 0 {
		return nil, nil, fmt.Errorf("empty data")
	}

	// determine commitment mode
	mode, err := determineCommitmentMode(data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to determine commitment mode: %w", err)
	}

	// parse cert
	var versionedCert *certs.VersionedCert
	var prefix PrefixMetadata
	prefix.Mode = mode
	// length of binary data on L1
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

	return &prefix, versionedCert, nil
}

// ProcessHexString processes a hex-encoded string and returns binary data for RLP decoding
func ProcessHexString(hexString string) ([]byte, error) {
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

// determineCommitmentMode uses RLP validation to distinguish between [commitments.StandardCommitmentMode]
// and [commitments.OptimismGenericCommitmentMode]. The standard commitment with cert version 1 and Optimism
// Generic Commitment produce a leading byte 1.
// Without asking user to indicate the type, we use the following test for which commitment a serialized altda
// commitment belongs. In RLP spec, https://ethereum.org/en/developers/docs/data-structures-and-encoding/rlp/.
// By RLP decode, a standard commitment cannot possibly have a leading 0 in its rlp encoded data, unless the data
// to be serialized contains a single byte.
func determineCommitmentMode(data []byte) (commitments.CommitmentMode, error) {
	// for the smaller standard commitment, we assume it must have at least 3 bytes. Which is pretty reasonable
	// given the size of a cert is far greater than 3.
	// standard commitment = [version_byte][rlp_certificate]. Size of 3 eliminates the case which rlp_certificate
	// is a single byte and therefore rlp_certificate cannot start with 0 byte. Given this case is elimniated,
	// the data must either be a [commitments.OptimismGenericCommitmentMode] or a incorrect altda commitment
	if len(data) <= 3 {
		return "", fmt.Errorf("insufficient data")
	}

	if commitments.OPCommitmentByte(data[0]) == commitments.OPKeccak256CommitmentByte {
		return "", fmt.Errorf("OP Keccak commitment not supported for not containing altda commitment")
	}

	// First, try to parse as Standard mode: [version_byte][rlp_certificate]
	if isValidRLP(data[1:]) {
		return commitments.StandardCommitmentMode, nil
	}

	// If Standard mode RLP validation failed, check for Optimism Generic mode
	// Optimism Generic: [0x01][da_layer_byte][version_byte][rlp_certificate]
	if isValidRLP(data[3:]) {
		return commitments.OptimismGenericCommitmentMode, nil
	} else {
		// If we can't determine the mode conclusively
		return "", fmt.Errorf("cannot determine commitment mode for a data of size %v", len(data))
	}
}

// isValidRLP attempts to validate if the data is valid RLP encoding
func isValidRLP(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	// Try to decode as both V2, V3 and V4 certificates to validate RLP structure
	var certV4 coretypes.EigenDACertV4
	if err := rlp.DecodeBytes(data, &certV4); err == nil {
		return true
	}

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

