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
// Right now this test relies on our op fork's InvalidCommitmentError definition, which we are
// changing in https://github.com/Layr-Labs/optimism/pull/50...
// Prob best to merge that PR first, and then update this test suite.
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
			name:                 "RBN recency check failed - invalid cert",
			RBNRecencyWindowSize: 100,
			certRBN:              100,
			certL1IBN:            201,
			requireErrorFn: func(t *testing.T, err error) {
				// expect proxy to return a 418 error which the client converts to this structured error
				var invalidCommitmentErr altda.InvalidCommitmentError
				require.ErrorAs(t, err, &invalidCommitmentErr)
				require.Equal(t,
					int(coretypes.ErrRecencyCheckFailedDerivationError.StatusCode),
					invalidCommitmentErr.StatusCode)
			},
		},
		{
			name:                 "RBN recency check passed - valid cert",
			RBNRecencyWindowSize: 100,
			certRBN:              100,
			certL1IBN:            199,
			requireErrorFn: func(t *testing.T, err error) {
				// After RBN check succeeds, CertVerifier.checkDACert contract call is made,
				// which returns a [verification.CertVerificationFailedError] with StatusCode 2 (inclusion proof
				// invalid). This gets converted to a [eigendav2store.ErrInvalidCertDerivationError] which gets marshalled
				// and returned as the body of a 418 response by the proxy.
				var invalidCommitmentErr altda.InvalidCommitmentError
				require.ErrorAs(t, err, &invalidCommitmentErr)
				require.Equal(t, int(coretypes.ErrInvalidCertDerivationError.StatusCode), invalidCommitmentErr.StatusCode)
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
				var invalidCommitmentErr altda.InvalidCommitmentError
				require.ErrorAs(t, err, &invalidCommitmentErr)
				require.Equal(t, int(coretypes.ErrInvalidCertDerivationError.StatusCode), invalidCommitmentErr.StatusCode)
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
				var invalidCommitmentErr altda.InvalidCommitmentError
				require.ErrorAs(t, err, &invalidCommitmentErr)
				require.Equal(t, int(coretypes.ErrInvalidCertDerivationError.StatusCode), invalidCommitmentErr.StatusCode)
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
