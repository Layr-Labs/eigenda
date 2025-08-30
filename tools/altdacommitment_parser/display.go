package altdacommitment_parser

import (
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// DisplayCommitmentInfo displays the parsed commitment structure information
func DisplayCommitmentInfo(parsed *ParsedCommitment) {
	fmt.Printf("Decoded hex string to binary (%d bytes)\n", len(parsed.OriginalData))
	fmt.Printf("Commitment Structure Analysis:\n")
	fmt.Printf("  Mode: %s\n", parsed.Mode)
	fmt.Printf("  Total Size: %d bytes\n", len(parsed.OriginalData))

	if parsed.CommitTypeByte != nil {
		fmt.Printf("  Commit Type Byte: 0x%02x\n", *parsed.CommitTypeByte)
	}
	if parsed.DALayerByte != nil {
		fmt.Printf("  DA Layer Byte: 0x%02x", *parsed.DALayerByte)
		if *parsed.DALayerByte == 0x00 {
			fmt.Printf(" (EigenDA)")
		}
		fmt.Printf("\n")
	}
	if parsed.VersionByte != nil {
		fmt.Printf("  Version Byte: 0x%02x (%s)\n", byte(*parsed.VersionByte), parsed.VersionByte.String())
	}
	fmt.Printf("  Certificate Data Size: %d bytes\n\n", len(parsed.CertificateData))
}

// displayCertificateData displays the parsed certificate data
func DisplayCertificateData(result *ParseResult) error {
	// Determine expected certificate version based on version byte
	var expectedVersion string
	if result.Commitment.VersionByte != nil {
		switch *result.Commitment.VersionByte {
		case V0VersionByte:
			expectedVersion = "V1/V2"
		case V1VersionByte:
			expectedVersion = "V2"
		case V2VersionByte:
			expectedVersion = "V2/V3"
		}
	}

	// Display the certificate based on version
	if result.Version == "V3" && result.CertV3 != nil {
		fmt.Printf("Successfully parsed as EigenDA Certificate V3")
		if expectedVersion != "" {
			fmt.Printf(" (Expected: %s)", expectedVersion)
		}
		fmt.Printf("\n\n")
		displayCertV3(result.CertV3)
		return nil
	}

	if result.Version == "V2" && result.CertV2 != nil {
		fmt.Printf("Successfully parsed as EigenDA Certificate V2")
		if expectedVersion != "" {
			fmt.Printf(" (Expected: %s)", expectedVersion)
		}
		fmt.Printf("\n\n")
		displayCertV2(result.CertV2)
		return nil
	}

	return fmt.Errorf("no valid certificate found in parse result")
}

// displayCertV3 creates a nicely formatted table display for V3 certificates
func displayCertV3(cert *coretypes.EigenDACertV3) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleDefault)
	t.Style().Title.Align = text.AlignCenter

	// Set column widths to ensure consistent display with truncated long numbers
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: 30, WidthMin: 30}, // Field column - fixed 30 characters
		{Number: 2, WidthMax: 80},               // Value column - back to 80 chars with truncation handling
	})

	// Main certificate info
	t.SetTitle("EigenDA Certificate V3 Details")
	t.AppendHeader(table.Row{"Field", "Value"})

	// Blob Inclusion Info
	t.AppendSeparator()
	section := "BLOB INCLUSION INFO"
	t.AppendRow(table.Row{section, section}, table.RowConfig{
		AutoMerge:      true,
		AutoMergeAlign: text.AlignCenter,
	})
	t.AppendSeparator()

	blobCert := &cert.BlobInclusionInfo.BlobCertificate
	t.AppendRow(table.Row{"Blob Index", fmt.Sprintf("%d", cert.BlobInclusionInfo.BlobIndex)})
	t.AppendRow(table.Row{"Inclusion Proof", formatByteSlice(cert.BlobInclusionInfo.InclusionProof)})

	// Blob Header
	section = "BLOB HEADER"
	t.AppendSeparator()
	t.AppendRow(table.Row{section, section}, table.RowConfig{
		AutoMerge:      true,
		AutoMergeAlign: text.AlignCenter,
	})
	t.AppendSeparator()

	blobHeader := &blobCert.BlobHeader
	t.AppendRow(table.Row{"Version", fmt.Sprintf("%d", blobHeader.Version)})
	t.AppendRow(table.Row{"Quorum Numbers", formatByteSlice(blobHeader.QuorumNumbers)})
	t.AppendRow(table.Row{"Payment Header Hash", formatByteArray32(blobHeader.PaymentHeaderHash)})

	// Commitment details
	section = "BLOB COMMITMENT"
	commitment := &blobHeader.Commitment
	t.AppendSeparator()
	t.AppendRow(table.Row{section, section}, table.RowConfig{
		AutoMerge:      true,
		AutoMergeAlign: text.AlignCenter,
	})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Commitment X", formatBigInt(commitment.Commitment.X)})
	t.AppendRow(table.Row{"Commitment Y", formatBigInt(commitment.Commitment.Y)})
	t.AppendRow(table.Row{"Length Commitment X", formatBigIntArray(commitment.LengthCommitment.X)})
	t.AppendRow(table.Row{"Length Commitment Y", formatBigIntArray(commitment.LengthCommitment.Y)})
	t.AppendRow(table.Row{"Length Proof X", formatBigIntArray(commitment.LengthProof.X)})
	t.AppendRow(table.Row{"Length Proof Y", formatBigIntArray(commitment.LengthProof.Y)})
	t.AppendRow(table.Row{"Length", fmt.Sprintf("%d", commitment.Length)})

	// Blob certificate signature and relay keys
	section = "BLOB CERTIFICATE"
	t.AppendSeparator()
	t.AppendRow(table.Row{section, section}, table.RowConfig{
		AutoMerge:      true,
		AutoMergeAlign: text.AlignCenter,
	})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Signature", formatByteSlice(blobCert.Signature)})
	t.AppendRow(table.Row{"Relay Keys", formatRelayKeys(blobCert.RelayKeys)})

	// Batch Header
	section = "BATCH HEADER"
	t.AppendSeparator()
	t.AppendRow(table.Row{section, section}, table.RowConfig{
		AutoMerge:      true,
		AutoMergeAlign: text.AlignCenter,
	})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Batch Root", formatByteArray32(cert.BatchHeader.BatchRoot)})
	t.AppendRow(table.Row{"Reference Block Number", fmt.Sprintf("%d", cert.BatchHeader.ReferenceBlockNumber)})

	// Non-Signer Stakes and Signature
	section = "NON-SIGNER STAKES & SIGNATURE"
	nonSigner := &cert.NonSignerStakesAndSignature
	t.AppendSeparator()
	t.AppendRow(table.Row{section, section}, table.RowConfig{
		AutoMerge:      true,
		AutoMergeAlign: text.AlignCenter,
	})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Non-Signer Quorum Bitmap Indices", formatUint32Slice(nonSigner.NonSignerQuorumBitmapIndices)})
	t.AppendRow(table.Row{"Non-Signer Pubkeys Count", fmt.Sprintf("%d", len(nonSigner.NonSignerPubkeys))})
	t.AppendRow(table.Row{"Quorum APKs Count", fmt.Sprintf("%d", len(nonSigner.QuorumApks))})
	t.AppendRow(table.Row{"APK G2 X", formatBigIntArray(nonSigner.ApkG2.X)})
	t.AppendRow(table.Row{"APK G2 Y", formatBigIntArray(nonSigner.ApkG2.Y)})
	t.AppendRow(table.Row{"Sigma X", formatBigInt(nonSigner.Sigma.X)})
	t.AppendRow(table.Row{"Sigma Y", formatBigInt(nonSigner.Sigma.Y)})
	t.AppendRow(table.Row{"Quorum APK Indices", formatUint32Slice(nonSigner.QuorumApkIndices)})
	t.AppendRow(table.Row{"Total Stake Indices", formatUint32Slice(nonSigner.TotalStakeIndices)})
	t.AppendRow(table.Row{"Non-Signer Stake Indices", formatUint32SliceSlice(nonSigner.NonSignerStakeIndices)})

	// Signed Quorum Numbers
	section = "SIGNED QUORUM NUMBERS"
	t.AppendSeparator()
	t.AppendRow(table.Row{section, section}, table.RowConfig{
		AutoMerge:      true,
		AutoMergeAlign: text.AlignCenter,
	})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Signed Quorum Numbers", formatByteSlice(cert.SignedQuorumNumbers)})

	t.Render()
}

// displayCertV2 creates a nicely formatted table display for V2 certificates
func displayCertV2(cert *coretypes.EigenDACertV2) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleDefault)
	t.Style().Title.Align = text.AlignCenter

	// Set column widths to ensure consistent display with truncated long numbers
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: 30, WidthMin: 30}, // Field column - fixed 30 characters
		{Number: 2, WidthMax: 80},               // Value column - back to 80 chars with truncation handling
	})

	// Main certificate info
	t.SetTitle("EigenDA Certificate V2 Details")
	t.AppendHeader(table.Row{"Field", "Value"})

	// Similar structure but for V2 types
	// Note: V2 and V3 have very similar structure, main difference is in the type bindings

	// Blob Inclusion Info
	t.AppendSeparator()
	t.AppendRow(table.Row{"BLOB INCLUSION INFO", ""})
	t.AppendSeparator()

	blobCert := &cert.BlobInclusionInfo.BlobCertificate
	t.AppendRow(table.Row{"Blob Index", fmt.Sprintf("%d", cert.BlobInclusionInfo.BlobIndex)})
	t.AppendRow(table.Row{"Inclusion Proof", formatByteSlice(cert.BlobInclusionInfo.InclusionProof)})

	// Blob Header
	t.AppendSeparator()
	t.AppendRow(table.Row{"BLOB HEADER", ""})
	t.AppendSeparator()

	blobHeader := &blobCert.BlobHeader
	t.AppendRow(table.Row{"Version", fmt.Sprintf("%d", blobHeader.Version)})
	t.AppendRow(table.Row{"Quorum Numbers", formatByteSlice(blobHeader.QuorumNumbers)})
	t.AppendRow(table.Row{"Payment Header Hash", formatByteArray32(blobHeader.PaymentHeaderHash)})

	// Commitment details
	commitment := &blobHeader.Commitment
	t.AppendSeparator()
	t.AppendRow(table.Row{"BLOB COMMITMENT", ""})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Commitment X", formatBigInt(commitment.Commitment.X)})
	t.AppendRow(table.Row{"Commitment Y", formatBigInt(commitment.Commitment.Y)})
	t.AppendRow(table.Row{"Length Commitment X", formatBigIntArray(commitment.LengthCommitment.X)})
	t.AppendRow(table.Row{"Length Commitment Y", formatBigIntArray(commitment.LengthCommitment.Y)})
	t.AppendRow(table.Row{"Length Proof X", formatBigIntArray(commitment.LengthProof.X)})
	t.AppendRow(table.Row{"Length Proof Y", formatBigIntArray(commitment.LengthProof.Y)})
	t.AppendRow(table.Row{"Length", fmt.Sprintf("%d", commitment.Length)})

	// Blob certificate signature and relay keys
	t.AppendSeparator()
	t.AppendRow(table.Row{"BLOB CERTIFICATE", ""})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Signature", formatByteSlice(blobCert.Signature)})
	t.AppendRow(table.Row{"Relay Keys", formatUint32Slice(blobCert.RelayKeys)})

	// Batch Header
	t.AppendSeparator()
	t.AppendRow(table.Row{"BATCH HEADER", ""})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Batch Root", formatByteArray32(cert.BatchHeader.BatchRoot)})
	t.AppendRow(table.Row{"Reference Block Number", fmt.Sprintf("%d", cert.BatchHeader.ReferenceBlockNumber)})

	// Non-Signer Stakes and Signature
	nonSigner := &cert.NonSignerStakesAndSignature
	t.AppendSeparator()
	t.AppendRow(table.Row{"NON-SIGNER STAKES & SIGNATURE", ""})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Non-Signer Quorum Bitmap Indices", formatUint32Slice(nonSigner.NonSignerQuorumBitmapIndices)})
	t.AppendRow(table.Row{"Non-Signer Pubkeys Count", fmt.Sprintf("%d", len(nonSigner.NonSignerPubkeys))})
	t.AppendRow(table.Row{"Quorum APKs Count", fmt.Sprintf("%d", len(nonSigner.QuorumApks))})
	t.AppendRow(table.Row{"APK G2 X", formatBigIntArray(nonSigner.ApkG2.X)})
	t.AppendRow(table.Row{"APK G2 Y", formatBigIntArray(nonSigner.ApkG2.Y)})
	t.AppendRow(table.Row{"Sigma X", formatBigInt(nonSigner.Sigma.X)})
	t.AppendRow(table.Row{"Sigma Y", formatBigInt(nonSigner.Sigma.Y)})
	t.AppendRow(table.Row{"Quorum APK Indices", formatUint32Slice(nonSigner.QuorumApkIndices)})
	t.AppendRow(table.Row{"Total Stake Indices", formatUint32Slice(nonSigner.TotalStakeIndices)})
	t.AppendRow(table.Row{"Non-Signer Stake Indices", formatUint32SliceSlice(nonSigner.NonSignerStakeIndices)})

	// Signed Quorum Numbers
	t.AppendSeparator()
	t.AppendRow(table.Row{"SIGNED QUORUM NUMBERS", ""})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Signed Quorum Numbers", formatByteSlice(cert.SignedQuorumNumbers)})

	t.Render()
}

// Formatting helper functions
func formatByteSlice(data []byte) string {
	if len(data) == 0 {
		return "[]"
	}
	return fmt.Sprintf("0x%s", hex.EncodeToString(data))
}

func formatByteArray32(data [32]byte) string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(data[:]))
}

func formatBigInt(val interface{}) string {
	if val == nil {
		return "nil"
	}

	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return "nil"
	}

	str := fmt.Sprintf("%v", val)
	return str
}

func formatBigIntArray(val interface{}) string {
	if val == nil {
		return "nil"
	}

	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Slice && v.Len() > 0 {
		elements := make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			str := fmt.Sprintf("%v", v.Index(i).Interface())
			elements[i] = str
		}
		// Use newlines to separate array elements so each big integer is on its own line
		return fmt.Sprintf("[\n  %s\n]", strings.Join(elements, ",\n  "))
	}

	return fmt.Sprintf("%v", val)
}

func formatUint32Slice(data []uint32) string {
	if len(data) == 0 {
		return "[]"
	}

	strs := make([]string, len(data))
	for i, v := range data {
		strs[i] = fmt.Sprintf("%d", v)
	}
	return fmt.Sprintf("[%s]", strings.Join(strs, ", "))
}

func formatUint32SliceSlice(data [][]uint32) string {
	if len(data) == 0 {
		return "[]"
	}

	strs := make([]string, len(data))
	for i, slice := range data {
		strs[i] = formatUint32Slice(slice)
	}
	return fmt.Sprintf("[%s]", strings.Join(strs, ", "))
}

func formatRelayKeys(keys interface{}) string {
	v := reflect.ValueOf(keys)
	if v.Kind() != reflect.Slice {
		return fmt.Sprintf("%v", keys)
	}

	if v.Len() == 0 {
		return "[]"
	}

	strs := make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		strs[i] = fmt.Sprintf("%v", v.Index(i).Interface())
	}
	return fmt.Sprintf("[%s]", strings.Join(strs, ", "))
}
