package certs

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
)

// Version byte that prefixes serialized EigenDACert to identify their type.
type VersionByte byte

const (
	// EigenDA V1
	V0VersionByte VersionByte = iota
	// All future CertVersions will be against EigenDA V2 Blazar (https://docs.eigenda.xyz/releases/blazar)
	V1VersionByte
	V2VersionByte
)

// versionByteString returns a string representation of the version byte for display
func (v VersionByte) VersionByteString() string {
	switch v {
	case V0VersionByte:
		return "EigenDA V1"
	case V1VersionByte:
		return "EigenDA V2 Legacy"
	case V2VersionByte:
		return "EigenDA V2 with V3 Cert"
	default:
		return fmt.Sprintf("Unknown (0x%02x)", byte(v))
	}
}

// IntoCertVersion converts from a version byte into a
// DA Cert type version enum
func (v VersionByte) IntoCertVersion() (coretypes.CertificateVersion, error) {
	switch v {
	case V0VersionByte:
		return 0, fmt.Errorf("V0 DA Commit version corresponds to EigenDAV1 which is unsupported for CertVersion")
	case V1VersionByte:
		return coretypes.VersionTwoCert, nil
	case V2VersionByte:
		return coretypes.VersionThreeCert, nil
	default:
		return 0, fmt.Errorf("unknown version byte (0x%02x)", byte(v))
	}
}

func ByteToVersion(b byte) (VersionByte, error) {
	switch b {
	case byte(V0VersionByte):
		return V0VersionByte, nil
	case byte(V1VersionByte):
		return V1VersionByte, nil
	case byte(V2VersionByte):
		return V2VersionByte, nil
	default:
		return 0, fmt.Errorf("unknown EigenDA cert version: %d", b)
	}
}

type VersionedCert struct {
	Version        VersionByte
	SerializedCert []byte
}

// NewVersionedCert creates a new EigenDA VersionedCert that holds the certVersion
// and a serialized certificate of that version.
func NewVersionedCert(serializedCert []byte, certVersion VersionByte) VersionedCert {
	return VersionedCert{
		Version:        certVersion,
		SerializedCert: serializedCert,
	}
}

// Encode adds a commitment type prefix self describing the commitment.
func (c VersionedCert) Encode() []byte {
	return append([]byte{byte(c.Version)}, c.SerializedCert...)
}
