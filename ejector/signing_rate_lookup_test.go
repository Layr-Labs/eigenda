package ejector

import (
	"fmt"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/stretchr/testify/require"
)

func TestDataApiLookup(t *testing.T) {
	test.SkipInCI(t)

	logger := common.TestLogger(t)
	url := "https://dataapi.eigenda.xyz"

	lookup := NewDynamoSigningRateLookup(logger, url, 100*time.Second)

	signingRates, err := lookup.GetSigningRates(1*time.Hour, []core.QuorumID{0, 1}, ProtocolVersionV2, false)
	require.NoError(t, err)

	for i, rate := range signingRates {
		validatorID := core.OperatorID(rate.ValidatorId)

		fmt.Printf("%d: %s\n", i, validatorID.Hex())
		fmt.Printf("        SignedBatches: %d\n", rate.SignedBatches)
		fmt.Printf("        UnsignedBatches: %d\n", rate.UnsignedBatches)
		fmt.Printf("        SignedBytes: %d\n", rate.SignedBytes)
		fmt.Printf("        UnsignedBytes: %d\n", rate.UnsignedBytes)
		fmt.Printf("        SigningLatency: %d\n", rate.SigningLatency)
	}
}
