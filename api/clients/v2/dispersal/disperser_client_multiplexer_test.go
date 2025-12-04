package dispersal

import (
	"slices"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/disperser"
	authv2 "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/committer"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
)

func createTestMultiplexer(
	t *testing.T,
	config *DisperserClientMultiplexerConfig,
) (*DisperserClientMultiplexer, *disperser.MockDisperserRegistry) {
	mockRegistry := disperser.NewMockDisperserRegistry()
	logger := common.TestLogger(t)

	privateKey := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	signer, err := authv2.NewLocalBlobRequestSigner(privateKey)
	require.NoError(t, err)

	kzgCommitter, err := committer.NewFromConfig(committer.Config{
		SRSNumberToLoad:   8192,
		G1SRSPath:         "../../../../resources/srs/g1.point",
		G2SRSPath:         "../../../../resources/srs/g2.point",
		G2TrailingSRSPath: "../../../../resources/srs/g2.trailing.point",
	})
	require.NoError(t, err)

	mockRegistry.SetDefaultDispersers([]uint32{1, 2, 3})
	mockRegistry.SetOnDemandDispersers([]uint32{1, 3})
	mockRegistry.SetDisperserGrpcUri(1, "disperser1.example.com:50051")
	mockRegistry.SetDisperserGrpcUri(2, "disperser2.example.com:50051")
	mockRegistry.SetDisperserGrpcUri(3, "disperser3.example.com:50051")

	dcm, err := NewDisperserClientMultiplexer(
		logger,
		config,
		mockRegistry,
		signer,
		kzgCommitter,
		metrics.NoopDispersalMetrics,
		random.NewTestRandom().Rand,
	)
	require.NoError(t, err)

	// Create reputations for all dispersers
	ctx := t.Context()
	now := time.Now()
	_, err = dcm.GetDisperserClient(ctx, now, false)
	require.NoError(t, err)

	// Set up distinct reputations:
	// - Disperser 1: worst reputation (2 failures) - IS on-demand
	// - Disperser 2: best reputation (1 success) - NOT on-demand
	// - Disperser 3: second-worst reputation (1 failure) - IS on-demand
	// Only report outcomes for non-blacklisted dispersers
	if !slices.Contains(config.DisperserBlacklist, 1) {
		err = dcm.ReportDispersalOutcome(1, false, now)
		require.NoError(t, err)
		err = dcm.ReportDispersalOutcome(1, false, now)
		require.NoError(t, err)
	}
	if !slices.Contains(config.DisperserBlacklist, 2) {
		err = dcm.ReportDispersalOutcome(2, true, now)
		require.NoError(t, err)
	}
	if !slices.Contains(config.DisperserBlacklist, 3) {
		err = dcm.ReportDispersalOutcome(3, false, now)
		require.NoError(t, err)
	}

	return dcm, mockRegistry
}

func TestGetDisperserClient_WithOnDemandPaymentFilter(t *testing.T) {
	multiplexer, _ := createTestMultiplexer(t, DefaultDisperserClientMultiplexerConfig())

	now := time.Now()

	selections := make(map[uint32]int)
	for range 1000 {
		client, err := multiplexer.GetDisperserClient(t.Context(), now, true)
		require.NoError(t, err)
		selections[client.GetConfig().DisperserID]++
	}

	// Disperser 2 has best reputation but is NOT on-demand, should never be selected
	require.Equal(t, 0, selections[2], "disperser 2 should never be selected (not on-demand)")
}

func TestGetDisperserClient_CleansUpOutdatedClient(t *testing.T) {
	config := DefaultDisperserClientMultiplexerConfig()
	config.DisperserBlacklist = []uint32{1, 3} // Only disperser 2 is eligible

	multiplexer, registry := createTestMultiplexer(t, config)

	client1, err := multiplexer.GetDisperserClient(t.Context(), time.Now(), false)
	require.NoError(t, err)
	require.Equal(t, uint32(2), client1.GetConfig().DisperserID)
	require.Equal(t, "disperser2.example.com:50051", client1.GetConfig().GrpcUri)

	// Update disperser 2's URI
	registry.SetDisperserGrpcUri(2, "new-uri:50051")

	client2, err := multiplexer.GetDisperserClient(t.Context(), time.Now(), false)
	require.NoError(t, err)
	require.Equal(t, uint32(2), client2.GetConfig().DisperserID)
	require.Equal(t, "new-uri:50051", client2.GetConfig().GrpcUri, "should create new client with new URI")
	require.NotSame(t, client1, client2, "should be different client instance")
}

func TestGetDisperserClient_AdditionalDispersersAndBlacklist(t *testing.T) {
	config := DefaultDisperserClientMultiplexerConfig()
	config.AdditionalDispersers = []uint32{4}
	config.DisperserBlacklist = []uint32{2}

	multiplexer, registry := createTestMultiplexer(t, config)

	registry.SetDisperserGrpcUri(4, "disperser4.example.com:50051")

	now := time.Now()

	selections := make(map[uint32]int)
	for range 1000 {
		client, err := multiplexer.GetDisperserClient(t.Context(), now, false)
		require.NoError(t, err)
		selections[client.GetConfig().DisperserID]++
	}

	require.Equal(t, 0, selections[2], "disperser 2 should never be selected (blacklisted)")
	require.Equal(t, 0, selections[1], "disperser 1 should never be selected (filtered out due to reputation)")
	// Dispersers 3 and 4 should both be selected
	require.Greater(t, selections[3], 0, "disperser 3 should be selected")
	require.Greater(t, selections[4], selections[3], "disperser 4 should be selected more than disperser 3")
}

func TestGetDisperserClient_NoEligibleDispersers(t *testing.T) {
	config := DefaultDisperserClientMultiplexerConfig()
	multiplexer, registry := createTestMultiplexer(t, config)

	registry.SetDefaultDispersers([]uint32{})

	_, err := multiplexer.GetDisperserClient(t.Context(), time.Now(), false)
	require.Error(t, err)
}

func TestReportDispersalOutcome(t *testing.T) {
	config := DefaultDisperserClientMultiplexerConfig()
	multiplexer, _ := createTestMultiplexer(t, config)

	now := time.Now()

	err := multiplexer.ReportDispersalOutcome(1, true, now)
	require.NoError(t, err)

	err = multiplexer.ReportDispersalOutcome(1, false, now)
	require.NoError(t, err)

	err = multiplexer.ReportDispersalOutcome(99, true, now)
	require.Error(t, err, "should error for unknown disperser")
}

func TestClose(t *testing.T) {
	config := DefaultDisperserClientMultiplexerConfig()
	multiplexer, _ := createTestMultiplexer(t, config)

	err := multiplexer.Close()
	require.NoError(t, err)

	err = multiplexer.Close()
	require.NoError(t, err, "should be idempotent")

	_, err = multiplexer.GetDisperserClient(t.Context(), time.Now(), false)
	require.Error(t, err, "should block GetDisperserClient after close")

	err = multiplexer.ReportDispersalOutcome(1, true, time.Now())
	require.Error(t, err, "should block ReportDispersalOutcome after close")
}
