package eigenda

import (
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/stretchr/testify/require"
)

func TestVerifyCertRBNRecencyCheck(t *testing.T) {

	testTable := []struct {
		name                  string
		certRBN               uint64
		certL1IBN             uint64
		rbnRecencyWindowSize  uint64
		expectError           bool
		expectedErrorContains string
	}{
		{
			name:                  "input sanization: certRBN should always be > 0",
			certRBN:               0,
			certL1IBN:             100,
			rbnRecencyWindowSize:  100,
			expectError:           true,
			expectedErrorContains: "bug",
		},
		{
			name:                 "input sanization: certL1IBN=0 should skip the test (return nil)",
			certRBN:              100,
			certL1IBN:            0,
			rbnRecencyWindowSize: 100,
			expectError:          false,
		},
		{
			name:                 "input sanization: rbnRecencyWindowSize=0 should skip the test (return nil)",
			certRBN:              100,
			certL1IBN:            101,
			rbnRecencyWindowSize: 0,
			expectError:          false,
		},
		{
			name:                  "input sanization: certL1IBN should always be > certRBN (when != 0)",
			certRBN:               100,
			certL1IBN:             100,
			rbnRecencyWindowSize:  100,
			expectError:           true,
			expectedErrorContains: "bug",
		},
		{
			name:                 "ok: certL1IBN = certRBN + rbnRecencyWindowSize",
			certRBN:              100,
			certL1IBN:            200,
			rbnRecencyWindowSize: 100,
			expectError:          false,
		},
		{
			name:                  "error: certL1IBN > certRBN + rbnRecencyWindowSize",
			certRBN:               100,
			certL1IBN:             201,
			rbnRecencyWindowSize:  100,
			expectError:           true,
			expectedErrorContains: coretypes.NewRBNRecencyCheckFailedError(100, 201, 100).Error(),
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			err := verifyCertRBNRecencyCheck(test.certRBN, test.certL1IBN, test.rbnRecencyWindowSize)
			if test.expectError {
				require.ErrorContains(t, err, test.expectedErrorContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
