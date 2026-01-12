package altdacommitment_parser

import (
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	certTypesBinding "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDACertTypeBindings"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// DisplayPrefixInfo displays the parsed commitment structure information
func DisplayPrefixInfo(parsed *PrefixMetadata) {
	fmt.Printf("Decoded hex string to binary (%d bytes)\n", parsed.OriginalSize)
	fmt.Printf("Commitment Structure Analysis:\n")
	fmt.Printf("  Mode: %s\n", parsed.Mode)

	if parsed.CommitTypeByte != nil {
		fmt.Printf("  Commitment Type Byte: 0x%02x\n", *parsed.CommitTypeByte)
	}
	if parsed.DALayerByte != nil {
		fmt.Printf("  DA Layer Byte: 0x%02x", *parsed.DALayerByte)
		if *parsed.DALayerByte == 0x00 {
			fmt.Printf(" (EigenDA)")
		}
		fmt.Printf("\n")
	}
	versionByte := parsed.CertVersion
	fmt.Printf("  Version Byte: 0x%02x (%s)\n", byte(versionByte), versionByte.VersionByteString())
}

// DisplayCertData creates a nicely formatted table display for V2, V3, or V4 certificates.
// It takes raw certificate bytes and attempts to parse as V4, then V3, then V2.
func DisplayCertData(certBytes []byte) error {
	if len(certBytes) == 0 {
		return fmt.Errorf("no certificate data to parse")
	}

	// Try to parse as V4 first
	var certV4 coretypes.EigenDACertV4
	err := rlp.DecodeBytes(certBytes, &certV4)
	if err == nil {
		displayCert(&certV4)
		return nil
	}

	// Try to parse as V3
	var certV3 coretypes.EigenDACertV3
	err = rlp.DecodeBytes(certBytes, &certV3)
	if err == nil {
		displayCert(&certV3)
		return nil
	}

	// Try to parse as V2 and convert to V3 for display
	var certV2 coretypes.EigenDACertV2
	err = rlp.DecodeBytes(certBytes, &certV2)
	if err == nil {
		certV3 := certV2.ToV3()
		displayCert(certV3)
		return nil
	}

	return fmt.Errorf("failed to parse certificate as V2, V3, or V4: %w", err)
}

// displayCert creates a nicely formatted table display for V3 or V4 certificates
func displayCert(cert interface{}) {
	// Extract common fields using type switch
	var blobInclusionInfo *certTypesBinding.EigenDATypesV2BlobInclusionInfo
	var batchHeader *certTypesBinding.EigenDATypesV2BatchHeaderV2
	var nonSignerStakesAndSignature *certTypesBinding.EigenDATypesV1NonSignerStakesAndSignature
	var signedQuorumNumbers []byte
	var offchainDerivationVersion *uint16
	var title string

	switch c := cert.(type) {
	case *coretypes.EigenDACertV3:
		blobInclusionInfo = &c.BlobInclusionInfo
		batchHeader = &c.BatchHeader
		nonSignerStakesAndSignature = &c.NonSignerStakesAndSignature
		signedQuorumNumbers = c.SignedQuorumNumbers
		title = "EigenDA Certificate V3 Details"
	case *coretypes.EigenDACertV4:
		blobInclusionInfo = &c.BlobInclusionInfo
		batchHeader = &c.BatchHeader
		nonSignerStakesAndSignature = &c.NonSignerStakesAndSignature
		signedQuorumNumbers = c.SignedQuorumNumbers
		offchainDerivationVersion = &c.OffchainDerivationVersion
		title = "EigenDA Certificate V4 Details"
	default:
		fmt.Printf("Unsupported certificate type: %T\n", cert)
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleDefault)
	t.Style().Title.Align = text.AlignCenter

	// Set column widths to ensure consistent display with truncated long numbers
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: 35, WidthMin: 35}, // Field column - fixed 35 characters
		{Number: 2, WidthMax: 80},               // Value column - back to 80 chars with truncation handling
	})

	// Main certificate info
	t.SetTitle(title)
	t.AppendHeader(table.Row{"Field", "Value"})

	// Blob Inclusion Info
	t.AppendSeparator()
	section := "BLOB INCLUSION INFO"
	t.AppendRow(table.Row{section, section}, table.RowConfig{
		AutoMerge:      true,
		AutoMergeAlign: text.AlignCenter,
	})
	t.AppendSeparator()

	blobCert := &blobInclusionInfo.BlobCertificate
	t.AppendRow(table.Row{"Blob Index", fmt.Sprintf("%d", blobInclusionInfo.BlobIndex)})
	t.AppendRow(table.Row{"Inclusion Proof", formatByteSlice(blobInclusionInfo.InclusionProof)})

	// Blob Header
	section = "BLOB HEADER"
	t.AppendSeparator()
	t.AppendRow(table.Row{section, section}, table.RowConfig{
		AutoMerge:      true,
		AutoMergeAlign: text.AlignCenter,
	})
	t.AppendSeparator()

	blobHeader := &blobCert.BlobHeader
	t.AppendRow(table.Row{"Blob Params Version", fmt.Sprintf("%d", blobHeader.Version)})
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
	t.AppendRow(table.Row{"Account ECDSA Signature", formatByteSlice(blobCert.Signature)})
	t.AppendRow(table.Row{"Relay Keys", formatRelayKeys(blobCert.RelayKeys)})

	// Batch Header
	section = "BATCH HEADER"
	t.AppendSeparator()
	t.AppendRow(table.Row{section, section}, table.RowConfig{
		AutoMerge:      true,
		AutoMergeAlign: text.AlignCenter,
	})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Batch Root", formatByteArray32(batchHeader.BatchRoot)})
	t.AppendRow(table.Row{"Reference Block Number", fmt.Sprintf("%d", batchHeader.ReferenceBlockNumber)})

	// Non-Signer Stakes and BLS Signature
	section = "NON-SIGNER STAKES & BLS SIGNATURE"
	t.AppendSeparator()
	t.AppendRow(table.Row{section, section}, table.RowConfig{
		AutoMerge:      true,
		AutoMergeAlign: text.AlignCenter,
	})
	t.AppendSeparator()
	t.AppendRow(table.Row{
		"Non-Signer Quorum Bitmap Indices",
		formatUint32Slice(nonSignerStakesAndSignature.NonSignerQuorumBitmapIndices),
	})
	t.AppendRow(table.Row{
		"Non-Signer Pubkeys Count",
		fmt.Sprintf("%d", len(nonSignerStakesAndSignature.NonSignerPubkeys)),
	})
	t.AppendRow(table.Row{"Quorum APKs Count", fmt.Sprintf("%d", len(nonSignerStakesAndSignature.QuorumApks))})
	t.AppendRow(table.Row{"APK G2 X", formatBigIntArray(nonSignerStakesAndSignature.ApkG2.X)})
	t.AppendRow(table.Row{"APK G2 Y", formatBigIntArray(nonSignerStakesAndSignature.ApkG2.Y)})
	t.AppendRow(table.Row{"Sigma X", formatBigInt(nonSignerStakesAndSignature.Sigma.X)})
	t.AppendRow(table.Row{"Sigma Y", formatBigInt(nonSignerStakesAndSignature.Sigma.Y)})
	t.AppendRow(table.Row{"Quorum APK Indices", formatUint32Slice(nonSignerStakesAndSignature.QuorumApkIndices)})
	t.AppendRow(table.Row{"Total Stake Indices", formatUint32Slice(nonSignerStakesAndSignature.TotalStakeIndices)})
	t.AppendRow(table.Row{
		"Non-Signer Stake Indices",
		formatUint32SliceSlice(nonSignerStakesAndSignature.NonSignerStakeIndices),
	})

	// Signed Quorum Numbers
	section = "SIGNED QUORUM NUMBERS"
	t.AppendSeparator()
	t.AppendRow(table.Row{section, section}, table.RowConfig{
		AutoMerge:      true,
		AutoMergeAlign: text.AlignCenter,
	})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Signed Quorum Numbers", formatByteSlice(signedQuorumNumbers)})

	// V4-specific fields
	if offchainDerivationVersion != nil {
		section = "OFFCHAIN DERIVATION VERSION"
		t.AppendSeparator()
		t.AppendRow(table.Row{section, section}, table.RowConfig{
			AutoMerge:      true,
			AutoMergeAlign: text.AlignCenter,
		})
		t.AppendSeparator()
		t.AppendRow(table.Row{"Offchain Derivation Version", fmt.Sprintf("%d", *offchainDerivationVersion)})
	}

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
