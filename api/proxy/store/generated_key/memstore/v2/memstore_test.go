package memstore

import (
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/require"
)

var (
	testLogger = logging.NewTextSLogger(os.Stdout, &logging.SLoggerOptions{})
)

const (
	testPreimage = "Four score and seven years ago"
)

func getDefaultMemStoreTestConfig() *memconfig.SafeConfig {
	return memconfig.NewSafeConfig(memconfig.Config{
		MaxBlobSizeBytes: 1024 * 1024,
		BlobExpiration:   0,
		PutLatency:       0,
		GetLatency:       0,
	})
}

func TestGetSet(t *testing.T) {
	g1Srs, err := kzg.ReadG1Points("../../../../resources/g1.point", 3000, 2)
	require.NoError(t, err)

	require.NoError(t, err)

	msV2 := New(
		t.Context(),
		testLogger,
		getDefaultMemStoreTestConfig(),
		g1Srs,
	)

	expected := []byte(testPreimage)
	versionedCert, err := msV2.Put(t.Context(), expected, coretypes.CertSerializationRLP)
	require.NoError(t, err)

	actual, err := msV2.Get(t.Context(), versionedCert, coretypes.CertSerializationRLP, false)
	require.NoError(t, err)
	require.Equal(t, expected, actual)

	// Test getting the encoded payload
	encodedPayload, err := msV2.Get(t.Context(), versionedCert, coretypes.CertSerializationRLP, true)
	require.NoError(t, err)
	require.NotEqual(t, expected, encodedPayload)
}

func TestGetSetV3Cert(t *testing.T) {
	g1Srs, err := kzg.ReadG1Points("../../../../resources/g1.point", 3000, 2)
	require.NoError(t, err)

	config := getDefaultMemStoreTestConfig()
	// Configure to use V3 certs
	err = config.SetCertVersion(coretypes.VersionThreeCert)
	require.NoError(t, err)

	msV2 := New(
		t.Context(),
		testLogger,
		config,
		g1Srs,
	)

	expected := []byte(testPreimage)
	versionedCert, err := msV2.Put(t.Context(), expected, coretypes.CertSerializationRLP)
	require.NoError(t, err)

	// Verify the version byte is correct for V3
	require.Equal(t, byte(0x2), byte(versionedCert.Version), "V3 cert should use V2VersionByte (0x2)")

	actual, err := msV2.Get(t.Context(), versionedCert, coretypes.CertSerializationRLP, false)
	require.NoError(t, err)
	require.Equal(t, expected, actual)

	// Test getting the encoded payload
	encodedPayload, err := msV2.Get(t.Context(), versionedCert, coretypes.CertSerializationRLP, true)
	require.NoError(t, err)
	require.NotEqual(t, expected, encodedPayload)
}

func TestGetSetV4Cert(t *testing.T) {
	g1Srs, err := kzg.ReadG1Points("../../../../resources/g1.point", 3000, 2)
	require.NoError(t, err)

	config := getDefaultMemStoreTestConfig()
	// Explicitly configure to use V4 certs
	err = config.SetCertVersion(coretypes.VersionFourCert)
	require.NoError(t, err)

	msV2 := New(
		t.Context(),
		testLogger,
		config,
		g1Srs,
	)

	expected := []byte(testPreimage)
	versionedCert, err := msV2.Put(t.Context(), expected, coretypes.CertSerializationRLP)
	require.NoError(t, err)

	// Verify the version byte is correct for V4
	require.Equal(t, byte(0x3), byte(versionedCert.Version), "V4 cert should use V3VersionByte (0x3)")

	actual, err := msV2.Get(t.Context(), versionedCert, coretypes.CertSerializationRLP, false)
	require.NoError(t, err)
	require.Equal(t, expected, actual)

	// Test getting the encoded payload
	encodedPayload, err := msV2.Get(t.Context(), versionedCert, coretypes.CertSerializationRLP, true)
	require.NoError(t, err)
	require.NotEqual(t, expected, encodedPayload)
}

func TestSwitchCertVersion(t *testing.T) {
	g1Srs, err := kzg.ReadG1Points("../../../../resources/g1.point", 3000, 2)
	require.NoError(t, err)

	config := getDefaultMemStoreTestConfig()
	msV2 := New(
		t.Context(),
		testLogger,
		config,
		g1Srs,
	)

	expected := []byte(testPreimage)

	// Store with V4 (default)
	versionedCertV4, err := msV2.Put(t.Context(), expected, coretypes.CertSerializationRLP)
	require.NoError(t, err)
	require.Equal(t, byte(0x3), byte(versionedCertV4.Version), "Should use V3VersionByte for V4 cert")

	// Switch to V3
	err = config.SetCertVersion(coretypes.VersionThreeCert)
	require.NoError(t, err)

	// Store with V3
	versionedCertV3, err := msV2.Put(t.Context(), expected, coretypes.CertSerializationRLP)
	require.NoError(t, err)
	require.Equal(t, byte(0x2), byte(versionedCertV3.Version), "Should use V2VersionByte for V3 cert")

	// Verify both can be retrieved correctly regardless of current config
	actualV4, err := msV2.Get(t.Context(), versionedCertV4, coretypes.CertSerializationRLP, false)
	require.NoError(t, err)
	require.Equal(t, expected, actualV4)

	actualV3, err := msV2.Get(t.Context(), versionedCertV3, coretypes.CertSerializationRLP, false)
	require.NoError(t, err)
	require.Equal(t, expected, actualV3)
}
