package e2e

import (
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	_ "github.com/Layr-Labs/eigenda/api/clients/v2/verification" // imported for docstring link
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda/api/proxy/test/testutils"
	bindings "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDACertTypeBindings"
	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
)

// TODO: update this to test all 4 derivation error cases.
//
// RBN Recency Check is only available for V2
// Contract Test here refers to https://pactflow.io/blog/what-is-contract-testing/, not evm contracts.
func TestOPContractTestRBNRecentyCheck(t *testing.T) {
	t.Parallel()
	if testutils.GetBackend() == testutils.MemstoreBackend {
		t.Skip("Don't run for memstore backend, since rbn recency check is only implemented for eigenda v2 backend")
	}

	var testTable = []struct {
		name                 string
		RBNRecencyWindowSize uint64
		certRBN              uint32
		certL1IBN            uint64
		requireErrorFn       func(t *testing.T, err error)
	}{
		{
			name:                 "RBN recency check failed",
			RBNRecencyWindowSize: 100,
			certRBN:              100,
			certL1IBN:            201,
			requireErrorFn: func(t *testing.T, err error) {
				// expect proxy to return a 418 error which the client converts to this structured error
				var dropEigenDACommitmentErr altda.DropEigenDACommitmentError
				require.ErrorAs(t, err, &dropEigenDACommitmentErr)
				require.Equal(t,
					int(coretypes.ErrRecencyCheckFailedDerivationError.StatusCode),
					dropEigenDACommitmentErr.StatusCode)
			},
		},
		{
			name:                 "RBN recency check passed",
			RBNRecencyWindowSize: 100,
			certRBN:              100,
			certL1IBN:            199,
			requireErrorFn: func(t *testing.T, err error) {
				// After RBN check succeeds, CertVerifier.checkDACert contract call is made,
				// which returns a [verification.CertVerificationFailedError] with StatusCode 2 (inclusion proof
				// invalid). This gets converted to a [eigendav2store.ErrInvalidCertDerivationError] which gets marshalled
				// and returned as the body of a 418 response by the proxy.
				var dropEigenDACommitmentErr altda.DropEigenDACommitmentError
				require.ErrorAs(t, err, &dropEigenDACommitmentErr)
				require.Equal(t, int(coretypes.ErrInvalidCertDerivationError.StatusCode), dropEigenDACommitmentErr.StatusCode)
			},
		},
		{
			name:                 "RBN recency check skipped - Proxy set window size 0",
			RBNRecencyWindowSize: 0,
			certRBN:              100,
			certL1IBN:            201,
			requireErrorFn: func(t *testing.T, err error) {
				// After RBN check succeeds, CertVerifier.checkDACert contract call is made,
				// which returns a [verification.CertVerificationFailedError] with StatusCode 2 (inclusion proof
				// invalid). This gets converted to a [eigendav2store.ErrInvalidCertDerivationError] which gets marshalled
				// and returned as the body of a 418 response by the proxy.
				var dropEigenDACommitmentErr altda.DropEigenDACommitmentError
				require.ErrorAs(t, err, &dropEigenDACommitmentErr)
				require.Equal(t, int(coretypes.ErrInvalidCertDerivationError.StatusCode), dropEigenDACommitmentErr.StatusCode)
			},
		},
		{
			name:                 "RBN recency check skipped - client set IBN to 0",
			RBNRecencyWindowSize: 100,
			certRBN:              100,
			certL1IBN:            0,
			requireErrorFn: func(t *testing.T, err error) {
				// After RBN check succeeds, CertVerifier.checkDACert contract call is made,
				// which returns a [verification.CertVerificationFailedError] with StatusCode 2 (inclusion proof
				// invalid). This gets converted to a [eigendav2store.ErrInvalidCertDerivationError] which gets marshalled
				// and returned as the body of a 418 response by the proxy.
				var dropEigenDACommitmentErr altda.DropEigenDACommitmentError
				require.ErrorAs(t, err, &dropEigenDACommitmentErr)
				require.Equal(t, int(coretypes.ErrInvalidCertDerivationError.StatusCode), dropEigenDACommitmentErr.StatusCode)
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			t.Log("Running test: ", tt.name)
			testCfg := testutils.NewTestConfig(
				testutils.GetBackend(),
				common.V2EigenDABackend,
				[]common.EigenDABackend{common.V2EigenDABackend})
			tsConfig := testutils.BuildTestSuiteConfig(testCfg)
			tsConfig.StoreBuilderConfig.ClientConfigV2.RBNRecencyWindowSize = tt.RBNRecencyWindowSize
			ts, kill := testutils.CreateTestSuite(tsConfig)
			t.Cleanup(kill)

			// Build + Serialize (empty) cert with the given RBN
			certV3 := coretypes.EigenDACertV3{
				BatchHeader: bindings.EigenDATypesV2BatchHeaderV2{
					ReferenceBlockNumber: tt.certRBN,
				},
			}
			serializedCertV3, err := rlp.EncodeToBytes(certV3)
			require.NoError(t, err)
			// altdaCommitment is what is returned by the proxy
			altdaCommitment, err := commitments.EncodeCommitment(
				certs.NewVersionedCert(serializedCertV3, certs.V2VersionByte),
				commitments.OptimismGenericCommitmentMode)
			require.NoError(t, err)
			// the op client expects a typed commitment, so we have to decode the altdaCommitment
			commitmentData, err := altda.DecodeCommitmentData(altdaCommitment)
			require.NoError(t, err)

			daClient := altda.NewDAClient(ts.Address(), false, false)
			_, err = daClient.GetInput(ts.Ctx, commitmentData, tt.certL1IBN)
			tt.requireErrorFn(t, err)
		})
	}
}

// Test that proxy DerivationErrors are correctly parsed as DropCommitmentErrors on op side,
// for parsing and cert validation errors.
func TestOPContractTestValidAndInvalidCertErrors(t *testing.T) {
	t.Parallel()
	if testutils.GetBackend() == testutils.MemstoreBackend {
		t.Skip("Don't run for memstore backend, since verifying certs is only done for eigenda v2 backend")
	}

	var testTable = []struct {
		name           string
		certCreationFn func() ([]byte, error)
		requireErrorFn func(t *testing.T, err error)
	}{
		{
			// TODO: need to figure out why this is happening, since ErrNotFound is supposed to be a keccak only error.
			// Seems like op-client allows submitting an empty cert, and because its not a valid cert request, it gets
			// matched by proxy's keccak commitment handler, which returns ErrNotFound (there is no such key in the store).
			// I think this is ok behavior... since it would be a bug to submit an empty cert....?
			// But need to think about this more.
			name: "empty cert returns ErrNotFound",
			certCreationFn: func() ([]byte, error) {
				return []byte{}, nil
			},
			requireErrorFn: func(t *testing.T, err error) {
				require.ErrorIs(t, err, altda.ErrNotFound)
			},
		},
		{
			name: "cert parsing error",
			certCreationFn: func() ([]byte, error) {
				cert := make([]byte, 10)
				return cert, nil
			},
			requireErrorFn: func(t *testing.T, err error) {
				var dropEigenDACommitmentErr altda.DropEigenDACommitmentError
				require.ErrorAs(t, err, &dropEigenDACommitmentErr)
				require.Equal(t, int(coretypes.ErrCertParsingFailedDerivationError.StatusCode), dropEigenDACommitmentErr.StatusCode)
			},
		},
		{
			name: "invalid (default) cert",
			certCreationFn: func() ([]byte, error) {
				// Build + Serialize invalid default cert
				certV3 := coretypes.EigenDACertV3{}
				serializedCertV3, err := rlp.EncodeToBytes(certV3)
				if err != nil {
					return nil, err
				}
				return serializedCertV3, nil
			},
			requireErrorFn: func(t *testing.T, err error) {
				var dropEigenDACommitmentErr altda.DropEigenDACommitmentError
				require.ErrorAs(t, err, &dropEigenDACommitmentErr)
				require.Equal(t, int(coretypes.ErrInvalidCertDerivationError.StatusCode), dropEigenDACommitmentErr.StatusCode)
			},
		},
	}

	testCfg := testutils.NewTestConfig(
		testutils.GetBackend(),
		common.V2EigenDABackend,
		[]common.EigenDABackend{common.V2EigenDABackend})
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	ts, kill := testutils.CreateTestSuite(tsConfig)
	t.Cleanup(kill)

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			t.Log("Running test: ", tt.name)
			serializedCert, err := tt.certCreationFn()
			require.NoError(t, err)

			altdaCommitment, err := commitments.EncodeCommitment(
				certs.NewVersionedCert(serializedCert, certs.V2VersionByte),
				commitments.OptimismGenericCommitmentMode)
			require.NoError(t, err)
			// the op client expects a typed commitment, so we have to decode the altdaCommitment
			commitmentData, err := altda.DecodeCommitmentData(altdaCommitment)
			require.NoError(t, err)

			daClient := altda.NewDAClient(ts.Address(), false, false)
			_, err = daClient.GetInput(ts.Ctx, commitmentData, 0)

			tt.requireErrorFn(t, err)
		})
	}

}

func TestOPContractTestBlobDecodingErrors(t *testing.T) {
	// Writing this test is a lot more involved... because we need to populate mock relay backends
	// that would return a blob that doesn't decode properly.
	// Probably will require adding this after we've created a better test suite framework for the eigenda clients.
	t.Skip("TODO: implement blob decoding errors test")
}
