package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	commonv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	docsv2 "github.com/Layr-Labs/eigenda/disperser/dataapi/docs/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru/v2"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

var errNotFound = errors.New("not found")

const (
	maxBlobAge = 14 * 24 * time.Hour

	// The max number of blobs to return from blob feed API, regardless of the time
	// range or "limit" param.
	maxNumBlobsPerBlobFeedResponse = 1000

	// The max number of batches to return from batch feed API, regardless of the time
	// range or "limit" param.
	maxNumBatchesPerBatchFeedResponse = 1000

	// The quorum IDs that are allowed to query for signing info are [0, maxQuorumIDAllowed]
	maxQuorumIDAllowed = 2

	// Suppose 1 batch/s, we cache 2 days worth of batch attestations.
	// Suppose 1KB for each attestation, this will be 173MB memory.
	maxNumBatchesToCache = 3600 * 24 * 2

	// Cache ~10mins worth of blobs for KV lookups
	maxNumKVBlobsToCache = 100 * 600
	// Cache ~1h worth of batches for KV lookups
	maxNumKVBatchesToCache = 3600

	cacheControlParam = "Cache-Control"

	// Static content
	maxBlobDataAge                  = 300
	maxBatchDataAge                 = 300
	maxOperatorDispersalResponseAge = 300

	// Rarely changing content
	maxOperatorsStakeAge    = 300 // not expect the stake changes frequently
	maxOperatorPortCheckAge = 60  // not expect validator port changes frequently, but it's consequential to have right port

	// Live content
	maxMetricAge        = 5
	maxThroughputAge    = 5
	maxBlobFeedAge      = 5
	maxBatchFeedAge     = 5
	maxDispersalFeedAge = 5
	maxSigningInfoAge   = 5
)

type (
	ErrorResponse struct {
		Error string `json:"error"`
	}
)

type ServerV2 struct {
	serverMode   string
	socketAddr   string
	allowOrigins []string
	logger       logging.Logger

	blobMetadataStore *blobstore.BlobMetadataStore
	subgraphClient    dataapi.SubgraphClient
	chainReader       core.Reader
	chainState        core.ChainState
	indexedChainState core.IndexedChainState
	promClient        dataapi.PrometheusClient
	metrics           *dataapi.Metrics
	meterer           *meterer.Meterer

	operatorHandler *dataapi.OperatorHandler
	metricsHandler  *dataapi.MetricsHandler

	// Feed cache
	batchFeedCache *FeedCache[corev2.Attestation]

	// KV caches for blobs, keyed by blobkey
	blobMetadataCache                *lru.Cache[string, *commonv2.BlobMetadata]
	blobAttestationInfoCache         *lru.Cache[string, *commonv2.BlobAttestationInfo]
	blobCertificateCache             *lru.Cache[string, *corev2.BlobCertificate]
	blobAttestationInfoResponseCache *lru.Cache[string, *BlobAttestationInfoResponse]

	// KV caches for batches, keyed by batch header hash
	batchResponseCache *lru.Cache[string, *BatchResponse]
}

func NewServerV2(
	config dataapi.Config,
	blobMetadataStore *blobstore.BlobMetadataStore,
	promClient dataapi.PrometheusClient,
	subgraphClient dataapi.SubgraphClient,
	chainReader core.Reader,
	chainState core.ChainState,
	indexedChainState core.IndexedChainState,
	logger logging.Logger,
	metrics *dataapi.Metrics,
	meterer *meterer.Meterer,
) (*ServerV2, error) {
	l := logger.With("component", "DataAPIServerV2")

	getBatchTimestampFn := func(item *corev2.Attestation) time.Time {
		return time.Unix(0, int64(item.AttestedAt))
	}
	fetchBatchFn := func(ctx context.Context, start, end time.Time, order FetchOrder, limit int) ([]*corev2.Attestation, error) {
		if order == Ascending {
			return blobMetadataStore.GetAttestationByAttestedAtForward(
				ctx, uint64(start.UnixNano())-1, uint64(end.UnixNano()), limit,
			)
		}
		return blobMetadataStore.GetAttestationByAttestedAtBackward(
			ctx, uint64(end.UnixNano()), uint64(start.UnixNano())-1, limit,
		)
	}
	batchFeedCache := NewFeedCache[corev2.Attestation](
		maxNumBatchesToCache,
		fetchBatchFn,
		getBatchTimestampFn,
		metrics.BatchFeedCacheMetrics,
	)

	blobMetadataCache, err := lru.New[string, *commonv2.BlobMetadata](maxNumKVBlobsToCache)
	if err != nil {
		return nil, fmt.Errorf("failed to create blobMetadataCache: %w", err)
	}
	blobAttestationInfoCache, err := lru.New[string, *commonv2.BlobAttestationInfo](maxNumKVBlobsToCache)
	if err != nil {
		return nil, fmt.Errorf("failed to create blobAttestationInfoCache: %w", err)
	}
	blobCertificateCache, err := lru.New[string, *corev2.BlobCertificate](maxNumKVBlobsToCache)
	if err != nil {
		return nil, fmt.Errorf("failed to create blobCertificateCache: %w", err)
	}
	blobAttestationInfoResponseCache, err := lru.New[string, *BlobAttestationInfoResponse](maxNumKVBlobsToCache)
	if err != nil {
		return nil, fmt.Errorf("failed to create blobAttestationInfoResponseCache: %w", err)
	}

	batchResponseCache, err := lru.New[string, *BatchResponse](maxNumKVBatchesToCache)
	if err != nil {
		return nil, fmt.Errorf("failed to create batchResponseCache: %w", err)
	}

	operatorHandler, err := dataapi.NewOperatorHandler(l, metrics, chainReader, chainState, indexedChainState, subgraphClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create operatorHandler: %w", err)
	}

	return &ServerV2{
		logger:                           l,
		serverMode:                       config.ServerMode,
		socketAddr:                       config.SocketAddr,
		allowOrigins:                     config.AllowOrigins,
		blobMetadataStore:                blobMetadataStore,
		promClient:                       promClient,
		subgraphClient:                   subgraphClient,
		chainReader:                      chainReader,
		chainState:                       chainState,
		indexedChainState:                indexedChainState,
		metrics:                          metrics,
		operatorHandler:                  operatorHandler,
		metricsHandler:                   dataapi.NewMetricsHandler(promClient, dataapi.V2),
		batchFeedCache:                   batchFeedCache,
		blobMetadataCache:                blobMetadataCache,
		blobAttestationInfoCache:         blobAttestationInfoCache,
		blobCertificateCache:             blobCertificateCache,
		blobAttestationInfoResponseCache: blobAttestationInfoResponseCache,
		batchResponseCache:               batchResponseCache,
		meterer:                          meterer,
	}, nil
}

func (s *ServerV2) Start() error {
	if s.serverMode == gin.ReleaseMode {
		// optimize performance and disable debug features.
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add recovery middleware (best practice according to Cursor)
	router.Use(gin.Recovery())

	basePath := "/api/v2"
	docsv2.SwaggerInfoV2.BasePath = basePath
	docsv2.SwaggerInfoV2.Host = os.Getenv("SWAGGER_HOST")

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = s.allowOrigins
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}

	if s.serverMode != gin.ReleaseMode {
		config.AllowOrigins = []string{"*"}
	}

	// Apply CORS middleware before routes
	router.Use(cors.New(config))

	// Add OPTIONS handlers for all routes
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	v2 := router.Group(basePath)
	{
		blobs := v2.Group("/blobs")
		{
			blobs.GET("/feed", s.FetchBlobFeed)
			blobs.GET("/:blob_key", s.FetchBlob)
			blobs.GET("/:blob_key/certificate", s.FetchBlobCertificate)
			blobs.GET("/:blob_key/attestation-info", s.FetchBlobAttestationInfo)
		}
		batches := v2.Group("/batches")
		{
			batches.GET("/feed", s.FetchBatchFeed)
			batches.GET("/:batch_header_hash", s.FetchBatch)
		}
		accounts := v2.Group("/accounts")
		{
			accounts.GET("/:account_id/blobs", s.FetchAccountBlobFeed)
			accounts.GET("/:account_id/payment-state", s.FetchAccountPaymentState)
			accounts.GET("/:account_id/reservation/usage", s.FetchAccountReservationUsage)
		}
		operators := v2.Group("/operators")
		{
			operators.GET("/:operator_id/dispersals", s.FetchOperatorDispersalFeed)
			operators.GET("/:operator_id/dispersals/:batch_header_hash/response", s.FetchOperatorDispersalResponse)
			operators.GET("/signing-info", s.FetchOperatorSigningInfo)
			operators.GET("/stake", s.FetchOperatorsStake)
			operators.GET("/node-info", s.FetchOperatorsNodeInfo)
			operators.GET("/liveness", s.CheckOperatorsLiveness)
		}
		metrics := v2.Group("/metrics")
		{
			metrics.GET("/summary", s.FetchMetricsSummary)
			metrics.GET("/timeseries/throughput", s.FetchMetricsThroughputTimeseries)
			metrics.GET("/timeseries/network-signing-rate", s.FetchNetworkSigningRate)
		}
		swagger := v2.Group("/swagger")
		{
			swagger.GET("/*any", ginswagger.WrapHandler(swaggerfiles.Handler, ginswagger.InstanceName("V2"), ginswagger.URL("/api/v2/swagger/doc.json")))

		}
	}

	router.GET("/", func(g *gin.Context) {
		g.JSON(http.StatusAccepted, gin.H{"status": "OK"})
	})

	router.Use(logger.SetLogger(
		logger.WithSkipPath([]string{"/"}),
	))

	srv := &http.Server{
		Addr:              s.socketAddr,
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	errChan := run(s.logger, srv)
	return <-errChan
}

func errorResponse(c *gin.Context, err error) {
	_ = c.Error(err)
	var code int
	switch {
	case errors.Is(err, errNotFound):
		code = http.StatusNotFound
	default:
		code = http.StatusInternalServerError
	}
	c.JSON(code, ErrorResponse{
		Error: err.Error(),
	})
}

func invalidParamsErrorResponse(c *gin.Context, err error) {
	_ = c.Error(err)
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: err.Error(),
	})
}

func run(logger logging.Logger, httpServer *http.Server) <-chan error {
	errChan := make(chan error, 1)
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-ctx.Done()

		logger.Info("shutdown signal received")

		defer func() {
			stop()
			close(errChan)
		}()

		if err := httpServer.Shutdown(context.Background()); err != nil {
			errChan <- err
		}
		logger.Info("shutdown completed")
	}()

	go func() {
		logger.Info("server v2 running", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	return errChan
}

func (s *ServerV2) Shutdown() error {
	return nil
}

func safeAccess(data map[string]map[uint8]int, i string, j uint8) (int, bool) {
	innerMap, ok := data[i]
	if !ok {
		return 0, false // Key i does not exist
	}
	val, ok := innerMap[j]
	if !ok {
		return 0, false // Key j does not exist in the inner map
	}
	return val, true
}
