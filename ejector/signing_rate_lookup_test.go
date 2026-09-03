package ejector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	dataapiv2 "github.com/Layr-Labs/eigenda/disperser/dataapi/v2"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/stretchr/testify/require"
)

func TestDataApiV2LookupUsesBlobQuorumResponsibility(t *testing.T) {
	operatorID := strings.Repeat("01", 32)
	perfectOperatorID := strings.Repeat("02", 32)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v2/operators/signing-info", r.URL.Path)
		require.Equal(t, blobQuorumSigningInfoAccountingMode, r.URL.Query().Get("accounting"))
		require.Equal(t, "43200", r.URL.Query().Get("interval"))
		require.Empty(t, r.URL.Query().Get("nonsigner_only"))

		err := json.NewEncoder(w).Encode(dataapiv2.OperatorsSigningInfoResponse{
			OperatorSigningInfo: []*dataapiv2.OperatorSigningInfo{
				{
					OperatorId:              operatorID,
					QuorumId:                0,
					TotalUnsignedBatches:    0,
					TotalResponsibleBatches: 10,
					TotalBatches:            100,
				},
				{
					OperatorId:              operatorID,
					QuorumId:                1,
					TotalUnsignedBatches:    1,
					TotalResponsibleBatches: 1,
					TotalBatches:            100,
				},
				{
					OperatorId:              perfectOperatorID,
					QuorumId:                0,
					TotalUnsignedBatches:    0,
					TotalResponsibleBatches: 5,
					TotalBatches:            100,
				},
			},
		})
		require.NoError(t, err)
	}))
	defer server.Close()

	lookup := NewDataApiSigningRateLookup(common.TestLogger(t), server.URL, time.Second)
	signingRates, err := lookup.GetSigningRates(
		12*time.Hour,
		nil,
		ProtocolVersionV2,
		true,
	)
	require.NoError(t, err)
	require.Len(t, signingRates, 1)

	validatorID, err := core.OperatorIDFromHex(operatorID)
	require.NoError(t, err)
	require.Equal(t, validatorID[:], signingRates[0].GetValidatorId())
	require.Equal(t, uint64(10), signingRates[0].GetSignedBatches())
	require.Equal(t, uint64(1), signingRates[0].GetUnsignedBatches())
	require.Equal(t, uint64(10), signingRates[0].GetSignedBytes())
	require.Equal(t, uint64(1), signingRates[0].GetUnsignedBytes())
}

func TestEvaluateValidatorsRetainsPerfectV2Signers(t *testing.T) {
	v1Lookup := &recordingSigningRateLookup{}
	v2Lookup := &recordingSigningRateLookup{}
	ejector := &Ejector{
		logger:                     common.TestLogger(t),
		signingRateLookupV1:        v1Lookup,
		signingRateLookupV2:        v2Lookup,
		ejectionCriteriaTimeWindow: time.Hour,
	}

	require.NoError(t, ejector.evaluateValidators())
	require.Equal(t, []bool{true}, v1Lookup.omitPerfectSigners)
	require.Equal(t, []bool{false}, v2Lookup.omitPerfectSigners)
}

func TestTranslateV2ToProtoRejectsInvalidResponsibilityCounts(t *testing.T) {
	_, err := translateV2ToProto(&dataapiv2.OperatorSigningInfo{
		OperatorId:              strings.Repeat("01", 32),
		TotalUnsignedBatches:    2,
		TotalResponsibleBatches: 1,
	})
	require.ErrorContains(t, err, "invalid batch counts")
}

type recordingSigningRateLookup struct {
	omitPerfectSigners []bool
}

func (r *recordingSigningRateLookup) GetSigningRates(
	_ time.Duration,
	_ []core.QuorumID,
	_ ProtocolVersion,
	omitPerfectSigners bool,
) ([]*validator.ValidatorSigningRate, error) {
	r.omitPerfectSigners = append(r.omitPerfectSigners, omitPerfectSigners)
	return nil, nil
}

func TestDataApiLookup(t *testing.T) {
	test.SkipInCI(t)

	logger := common.TestLogger(t)
	url := "https://dataapi.eigenda.xyz"

	lookup := NewDataApiSigningRateLookup(logger, url, 100*time.Second)

	signingRates, err := lookup.GetSigningRates(1*time.Hour, []core.QuorumID{0, 1}, ProtocolVersionV2, false)
	require.NoError(t, err)

	sortByUnsignedBytesDescending(signingRates)

	for i, rate := range signingRates {
		validatorID := core.OperatorID(rate.GetValidatorId())

		fmt.Printf("%d: %s\n", i, validatorID.Hex())
		fmt.Printf("        SignedBatches: %d\n", rate.GetSignedBatches())
		fmt.Printf("        UnsignedBatches: %d\n", rate.GetUnsignedBatches())
		fmt.Printf("        SignedBytes: %d\n", rate.GetSignedBytes())
		fmt.Printf("        UnsignedBytes: %d\n", rate.GetUnsignedBytes())
		fmt.Printf("        SigningLatency: %d\n", rate.GetSigningLatency())
	}
}
