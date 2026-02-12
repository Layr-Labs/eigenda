package node

import (
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/core/payments/vault"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func testLogger(t *testing.T) logging.Logger {
	t.Helper()
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)
	return logger
}

// makeTestMetrics creates a minimal Metrics with only the fields needed for testing.
func makeTestMetrics(t *testing.T) *Metrics {
	t.Helper()
	reg := prometheus.NewRegistry()
	return &Metrics{
		ReachabilityGauge: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "reachability_status",
			},
			[]string{"service"},
		),
		AccuSocketUpdates: promauto.With(reg).NewCounter(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "socket_updates_total",
			},
		),
	}
}

// --- computeMemoryPoolSize ---

func TestComputeMemoryPoolSize_ConstantSize(t *testing.T) {
	logger := testLogger(t)
	size, err := computeMemoryPoolSize(logger, "test pool", 1024, 0.5, 4096)
	require.NoError(t, err)
	assert.Equal(t, uint64(1024), size)
}

func TestComputeMemoryPoolSize_Fraction(t *testing.T) {
	logger := testLogger(t)
	size, err := computeMemoryPoolSize(logger, "test pool", 0, 0.25, 4096)
	require.NoError(t, err)
	assert.Equal(t, uint64(1024), size)
}

func TestComputeMemoryPoolSize_ZeroFraction(t *testing.T) {
	logger := testLogger(t)
	size, err := computeMemoryPoolSize(logger, "test pool", 0, 0.0, 4096)
	require.NoError(t, err)
	assert.Equal(t, uint64(0), size)
}

func TestComputeMemoryPoolSize_FractionTooHigh(t *testing.T) {
	logger := testLogger(t)
	_, err := computeMemoryPoolSize(logger, "test pool", 0, 1.5, 4096)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be between 0.0 and 1.0")
}

func TestComputeMemoryPoolSize_NegativeFraction(t *testing.T) {
	logger := testLogger(t)
	_, err := computeMemoryPoolSize(logger, "test pool", 0, -0.1, 4096)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be between 0.0 and 1.0")
}

// --- buildSocket ---

func TestBuildSocket(t *testing.T) {
	n := &Node{
		Config: &Config{
			Hostname:        "myhost.com",
			V2DispersalPort: "32005",
			V2RetrievalPort: "32006",
		},
	}
	socket := n.buildSocket()
	// Format: host:v2Dispersal;v2Retrieval;v2Dispersal;v2Retrieval
	assert.Contains(t, socket, "myhost.com")
	assert.Contains(t, socket, "32005")
	assert.Contains(t, socket, "32006")
}

// --- processReachabilityResponse ---

func TestProcessReachabilityResponse_AllOnline(t *testing.T) {
	logger := testLogger(t)
	n := &Node{
		Logger:  logger,
		Metrics: makeTestMetrics(t),
	}
	resp := OperatorReachabilityResponse{
		DispersalOnline: true,
		DispersalSocket: "host:32005",
		DispersalStatus: "SERVING",
		RetrievalOnline: true,
		RetrievalSocket: "host:32006",
		RetrievalStatus: "SERVING",
	}
	n.processReachabilityResponse("v2", resp)
}

func TestProcessReachabilityResponse_AllOffline(t *testing.T) {
	logger := testLogger(t)
	n := &Node{
		Logger:  logger,
		Metrics: makeTestMetrics(t),
	}
	resp := OperatorReachabilityResponse{
		DispersalOnline: false,
		DispersalSocket: "host:32005",
		DispersalStatus: "UNREACHABLE",
		RetrievalOnline: false,
		RetrievalSocket: "host:32006",
		RetrievalStatus: "UNREACHABLE",
	}
	n.processReachabilityResponse("v1", resp)
}

// --- startNodeAPI ---

func TestStartNodeAPI_Disabled(t *testing.T) {
	logger := testLogger(t)
	n := &Node{
		Config: &Config{EnableNodeApi: false},
		Logger: logger,
	}
	// Should not panic when NodeApi is nil and EnableNodeApi is false.
	n.startNodeAPI()
}

// --- startMetrics ---

func TestStartMetrics_Disabled(t *testing.T) {
	logger := testLogger(t)
	n := &Node{
		Config: &Config{EnableMetrics: false},
		Logger: logger,
	}
	n.startMetrics()
}

// --- startPprof ---

func TestStartPprof_Disabled(t *testing.T) {
	logger := testLogger(t)
	n := &Node{
		Config: &Config{
			EnablePprof:   false,
			PprofHttpPort: "6060",
		},
		Logger: logger,
	}
	n.startPprof()
}

// --- startNodeIPUpdater ---

func TestStartNodeIPUpdater_Disabled(t *testing.T) {
	logger := testLogger(t)
	n := &Node{
		Config: &Config{PubIPCheckInterval: 0},
		Logger: logger,
	}
	n.startNodeIPUpdater()
}

// --- startOnDemandMeterer ---

func TestStartOnDemandMeterer_NilMeterer(t *testing.T) {
	n := &Node{
		onDemandMeterer: nil,
		Config:          &Config{},
	}
	n.startOnDemandMeterer(t.Context())
}

func TestStartOnDemandMeterer_ZeroInterval(t *testing.T) {
	ctx := t.Context()
	pv := vault.NewTestPaymentVault()
	pv.SetGlobalSymbolsPerSecond(10)
	pv.SetGlobalRatePeriodInterval(1)
	pv.SetMinNumSymbols(1)

	m, err := meterer.NewOnDemandMeterer(ctx, pv, time.Now, nil, 1.0)
	require.NoError(t, err)

	n := &Node{
		onDemandMeterer: m,
		Config: &Config{
			OnDemandMeterRefreshInterval: 0,
		},
	}
	n.startOnDemandMeterer(ctx)
}

// --- checkNodeReachability ---

func TestCheckNodeReachability_Disabled(t *testing.T) {
	logger := testLogger(t)
	n := &Node{
		Config: &Config{ReachabilityPollIntervalSec: 0},
		Logger: logger,
	}
	// ReachabilityPollIntervalSec == 0 causes immediate return.
	n.checkNodeReachability("api/v2/operators/liveness")
}

func TestCheckNodeReachability_NoDataApiUrl(t *testing.T) {
	logger := testLogger(t)
	n := &Node{
		Config: &Config{
			ReachabilityPollIntervalSec: 30,
			DataApiUrl:                  "",
		},
		Logger: logger,
	}
	// Empty DataApiUrl causes immediate return.
	n.checkNodeReachability("api/v2/operators/liveness")
}

// --- checkValidatorRegistration ---

func TestCheckValidatorRegistration_SocketMatch_KnownChain(t *testing.T) {
	logger := testLogger(t)
	tx := &coremock.MockWriter{}
	socket := "myhost:32005;32006;32005;32006"
	tx.On("GetOperatorSocket", mock.Anything, mock.Anything).Return(socket, nil)

	n := &Node{
		CTX:        t.Context(),
		Config:     &Config{ID: core.OperatorID{1}},
		Logger:     logger,
		Transactor: tx,
		ChainID:    big.NewInt(1), // known chain ID → logs EigenDA URL
	}
	n.checkValidatorRegistration(socket)
	tx.AssertExpectations(t)
}

func TestCheckValidatorRegistration_SocketMismatch_UnknownChain(t *testing.T) {
	logger := testLogger(t)
	tx := &coremock.MockWriter{}
	tx.On("GetOperatorSocket", mock.Anything, mock.Anything).Return("other:1111;2222;1111;2222", nil)

	n := &Node{
		CTX:        t.Context(),
		Config:     &Config{ID: core.OperatorID{1}},
		Logger:     logger,
		Transactor: tx,
		ChainID:    big.NewInt(99999), // unknown chain ID
	}
	n.checkValidatorRegistration("expected:32005;32006;32005;32006")
	tx.AssertExpectations(t)
}

func TestCheckValidatorRegistration_TransactorError(t *testing.T) {
	logger := testLogger(t)
	tx := &coremock.MockWriter{}
	tx.On("GetOperatorSocket", mock.Anything, mock.Anything).Return("", assert.AnError)

	n := &Node{
		CTX:        t.Context(),
		Config:     &Config{ID: core.OperatorID{1}},
		Logger:     logger,
		Transactor: tx,
		ChainID:    big.NewInt(17000), // known chain (holesky)
	}
	n.checkValidatorRegistration("myhost:32005;32006;32005;32006")
	tx.AssertExpectations(t)
}

// --- MeterOnDemandDispersal nil meterer ---

func TestMeterOnDemandDispersal_NilMeterer(t *testing.T) {
	n := &Node{onDemandMeterer: nil}
	_, err := n.MeterOnDemandDispersal(100)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not configured")
}

// --- updateSocketAddress ---

func TestUpdateSocketAddress_NoChange(t *testing.T) {
	logger := testLogger(t)
	n := &Node{
		Config:        &Config{},
		Logger:        logger,
		Metrics:       makeTestMetrics(t),
		CurrentSocket: "host:32005;32006;32005;32006",
	}
	// Same socket → should be a no-op.
	n.updateSocketAddress(t.Context(), "host:32005;32006;32005;32006")
	assert.Equal(t, "host:32005;32006;32005;32006", n.CurrentSocket)
}

func TestUpdateSocketAddress_Changed(t *testing.T) {
	logger := testLogger(t)
	tx := &coremock.MockWriter{}
	tx.On("UpdateOperatorSocket", mock.Anything, mock.Anything).Return(nil)

	n := &Node{
		Config:        &Config{},
		Logger:        logger,
		Transactor:    tx,
		Metrics:       makeTestMetrics(t),
		CurrentSocket: "old:32005;32006;32005;32006",
	}
	n.updateSocketAddress(t.Context(), "new:32005;32006;32005;32006")
	assert.Equal(t, "new:32005;32006;32005;32006", n.CurrentSocket)
	tx.AssertExpectations(t)
}

func TestUpdateSocketAddress_TransactorError(t *testing.T) {
	logger := testLogger(t)
	tx := &coremock.MockWriter{}
	tx.On("UpdateOperatorSocket", mock.Anything, mock.Anything).Return(assert.AnError)

	n := &Node{
		Config:        &Config{},
		Logger:        logger,
		Transactor:    tx,
		Metrics:       makeTestMetrics(t),
		CurrentSocket: "old:32005;32006;32005;32006",
	}
	n.updateSocketAddress(t.Context(), "new:32005;32006;32005;32006")
	// Socket should NOT change when the transactor errors.
	assert.Equal(t, "old:32005;32006;32005;32006", n.CurrentSocket)
	tx.AssertExpectations(t)
}

// --- ValidateReservationPayment ---

func TestValidateReservationPayment_EmptyBatch(t *testing.T) {
	n := &Node{}
	batch := &corev2.Batch{BlobCertificates: []*corev2.BlobCertificate{}}
	// nil SequenceProbe is safe — SetStage handles nil receiver.
	err := n.ValidateReservationPayment(t.Context(), batch, nil)
	assert.NoError(t, err)
}

// --- startV2 ---

func TestStartV2_DisabledRefreshAndReachability(t *testing.T) {
	logger := testLogger(t)
	n := &Node{
		Config: &Config{
			OnchainStateRefreshInterval: 0,
			ReachabilityPollIntervalSec: 0,
		},
		Logger: logger,
	}
	// Both goroutines exit immediately because refresh interval <= 0 and poll interval == 0.
	n.startV2()
}

// --- GetReachabilityURL error path ---

func TestGetReachabilityURL_InvalidBase(t *testing.T) {
	// url.JoinPath returns an error for certain malformed URLs.
	_, err := GetReachabilityURL("://bad", "path", "op123")
	assert.Error(t, err)
}

// --- configureMemoryLimits ---

func TestConfigureMemoryLimits_ConstantSizes(t *testing.T) {
	logger := testLogger(t)
	config := &Config{
		GCSafetyBufferSizeBytes:      1024,
		LittDBReadCacheSizeBytes:     2048,
		LittDBWriteCacheSizeBytes:    2048,
		StoreChunksBufferSizeBytes:   4096,
	}
	err := configureMemoryLimits(logger, config)
	require.NoError(t, err)
	assert.Equal(t, uint64(1024), config.GCSafetyBufferSizeBytes)
	assert.Equal(t, uint64(2048), config.LittDBReadCacheSizeBytes)
	assert.Equal(t, uint64(2048), config.LittDBWriteCacheSizeBytes)
	assert.Equal(t, uint64(4096), config.StoreChunksBufferSizeBytes)
}

func TestConfigureMemoryLimits_FractionSizes(t *testing.T) {
	logger := testLogger(t)
	// Use small fractions that should not exceed system memory.
	config := &Config{
		GCSafetyBufferSizeFraction:      0.01,
		LittDBReadCacheSizeFraction:     0.01,
		LittDBWriteCacheSizeFraction:    0.01,
		StoreChunksBufferSizeFraction:   0.01,
	}
	err := configureMemoryLimits(logger, config)
	require.NoError(t, err)
	// Each should be 1% of system memory. Just verify they're non-zero and consistent.
	assert.Greater(t, config.GCSafetyBufferSizeBytes, uint64(0))
	assert.Greater(t, config.LittDBReadCacheSizeBytes, uint64(0))
	assert.Greater(t, config.LittDBWriteCacheSizeBytes, uint64(0))
	assert.Greater(t, config.StoreChunksBufferSizeBytes, uint64(0))
}

func TestConfigureMemoryLimits_InvalidFraction(t *testing.T) {
	logger := testLogger(t)
	config := &Config{
		GCSafetyBufferSizeFraction: 2.0, // invalid
	}
	err := configureMemoryLimits(logger, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to compute size")
}
