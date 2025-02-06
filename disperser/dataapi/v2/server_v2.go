package v2

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/semver"
	disperserv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	docsv2 "github.com/Layr-Labs/eigenda/disperser/dataapi/docs/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
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

	cacheControlParam       = "Cache-Control"
	maxFeedBlobAge          = 300 // this is completely static
	maxOperatorsStakeAge    = 300 // not expect the stake change to happen frequently
	maxOperatorResponseAge  = 300 // this is completely static
	maxOperatorPortCheckAge = 60
	maxMetricAge            = 10
	maxThroughputAge        = 10
)

type (
	ErrorResponse struct {
		Error string `json:"error"`
	}

	SignedBatch struct {
		BatchHeader *corev2.BatchHeader `json:"batch_header"`
		Attestation *corev2.Attestation `json:"attestation"`
	}

	BlobResponse struct {
		BlobKey       string             `json:"blob_key"`
		BlobHeader    *corev2.BlobHeader `json:"blob_header"`
		Status        string             `json:"status"`
		DispersedAt   uint64             `json:"dispersed_at"`
		BlobSizeBytes uint64             `json:"blob_size_bytes"`
	}

	BlobCertificateResponse struct {
		Certificate *corev2.BlobCertificate `json:"blob_certificate"`
	}

	BlobAttestationInfoResponse struct {
		BlobKey         string                    `json:"blob_key"`
		BatchHeaderHash string                    `json:"batch_header_hash"`
		InclusionInfo   *corev2.BlobInclusionInfo `json:"blob_inclusion_info"`
		Attestation     *corev2.Attestation       `json:"attestation"`
	}

	BlobInfo struct {
		BlobKey      string                    `json:"blob_key"`
		BlobMetadata *disperserv2.BlobMetadata `json:"blob_metadata"`
	}
	BlobFeedResponse struct {
		Blobs           []BlobInfo `json:"blobs"`
		PaginationToken string     `json:"pagination_token"`
	}

	BatchResponse struct {
		BatchHeaderHash    string                      `json:"batch_header_hash"`
		SignedBatch        *SignedBatch                `json:"signed_batch"`
		BlobInclusionInfos []*corev2.BlobInclusionInfo `json:"blob_inclusion_infos"`
	}

	BatchInfo struct {
		BatchHeaderHash         string                  `json:"batch_header_hash"`
		BatchHeader             *corev2.BatchHeader     `json:"batch_header"`
		AttestedAt              uint64                  `json:"attested_at"`
		AggregatedSignature     *core.Signature         `json:"aggregated_signature"`
		QuorumNumbers           []core.QuorumID         `json:"quorum_numbers"`
		QuorumSignedPercentages map[core.QuorumID]uint8 `json:"quorum_signed_percentages"`
	}
	BatchFeedResponse struct {
		Batches []*BatchInfo `json:"batches"`
	}

	MetricSummary struct {
		AvgThroughput float64 `json:"avg_throughput"`
	}

	OperatorSigningInfo struct {
		OperatorId              string  `json:"operator_id"`
		OperatorAddress         string  `json:"operator_address"`
		QuorumId                uint8   `json:"quorum_id"`
		TotalUnsignedBatches    int     `json:"total_unsigned_batches"`
		TotalResponsibleBatches int     `json:"total_responsible_batches"`
		TotalBatches            int     `json:"total_batches"`
		SigningPercentage       float64 `json:"signing_percentage"`
		StakePercentage         float64 `json:"stake_percentage"`
	}
	OperatorsSigningInfoResponse struct {
		OperatorSigningInfo []*OperatorSigningInfo `json:"operator_signing_info"`
	}

	OperatorStake struct {
		QuorumId        string  `json:"quorum_id"`
		OperatorId      string  `json:"operator_id"`
		StakePercentage float64 `json:"stake_percentage"`
		Rank            int     `json:"rank"`
	}

	OperatorsStakeResponse struct {
		StakeRankedOperators map[string][]*OperatorStake `json:"stake_ranked_operators"`
	}

	// Operators' responses for a batch
	OperatorDispersalResponses struct {
		Responses []*corev2.DispersalResponse `json:"operator_dispersal_responses"`
	}

	OperatorPortCheckResponse struct {
		OperatorId      string `json:"operator_id"`
		DispersalSocket string `json:"dispersal_socket"`
		DispersalOnline bool   `json:"dispersal_online"`
		DispersalStatus string `json:"dispersal_status"`
		RetrievalSocket string `json:"retrieval_socket"`
		RetrievalOnline bool   `json:"retrieval_online"`
		RetrievalStatus string `json:"retrieval_status"`
	}

	SemverReportResponse struct {
		Semver map[string]*semver.SemverMetrics `json:"semver"`
	}

	Metric struct {
		Throughput float64 `json:"throughput"`
		CostInGas  float64 `json:"cost_in_gas"`
		// deprecated: use TotalStakePerQuorum instead. Remove when the frontend is updated.
		TotalStake          *big.Int                   `json:"total_stake"`
		TotalStakePerQuorum map[core.QuorumID]*big.Int `json:"total_stake_per_quorum"`
	}

	Throughput struct {
		Throughput float64 `json:"throughput"`
		Timestamp  uint64  `json:"timestamp"`
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

	operatorHandler *dataapi.OperatorHandler
	metricsHandler  *dataapi.MetricsHandler
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
) *ServerV2 {
	l := logger.With("component", "DataAPIServerV2")
	return &ServerV2{
		logger:            l,
		serverMode:        config.ServerMode,
		socketAddr:        config.SocketAddr,
		allowOrigins:      config.AllowOrigins,
		blobMetadataStore: blobMetadataStore,
		promClient:        promClient,
		subgraphClient:    subgraphClient,
		chainReader:       chainReader,
		chainState:        chainState,
		indexedChainState: indexedChainState,
		metrics:           metrics,
		operatorHandler:   dataapi.NewOperatorHandler(l, metrics, chainReader, chainState, indexedChainState, subgraphClient),
		metricsHandler:    dataapi.NewMetricsHandler(promClient),
	}
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
			blobs.GET("/feed", s.FetchBlobFeedHandler)
			blobs.GET("/:blob_key", s.FetchBlobHandler)
			blobs.GET("/:blob_key/certificate", s.FetchBlobCertificateHandler)
			blobs.GET("/:blob_key/attestation-info", s.FetchBlobAttestationInfo)
		}
		batches := v2.Group("/batches")
		{
			batches.GET("/feed", s.FetchBatchFeedHandler)
			batches.GET("/:batch_header_hash", s.FetchBatchHandler)
		}
		operators := v2.Group("/operators")
		{
			operators.GET("/signing-info", s.FetchOperatorSigningInfo)
			operators.GET("/stake", s.FetchOperatorsStake)
			operators.GET("/nodeinfo", s.FetchOperatorsNodeInfo)
			operators.GET("/reachability", s.CheckOperatorsReachability)
			operators.GET("/response/:batch_header_hash", s.FetchOperatorsResponses)
		}
		metrics := v2.Group("/metrics")
		{
			metrics.GET("/summary", s.FetchMetricsSummaryHandler)
			metrics.GET("/timeseries/throughput", s.FetchMetricsThroughputTimeseriesHandler)
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

// FetchBlobFeedHandler godoc
//
//	@Summary	Fetch blob feed
//	@Tags		Blobs
//	@Produce	json
//	@Param		end					query		string	false	"Fetch blobs up to the end time (ISO 8601 format: 2006-01-02T15:04:05Z) [default: now]"
//	@Param		interval			query		int		false	"Fetch blobs starting from an interval (in seconds) before the end time [default: 3600]"
//	@Param		pagination_token	query		string	false	"Fetch blobs starting from the pagination token (exclusively). Overrides the interval param if specified [default: empty]"
//	@Param		limit				query		int		false	"The maximum number of blobs to fetch. System max (1000) if limit <= 0 [default: 20; max: 1000]"
//	@Success	200					{object}	BlobFeedResponse
//	@Failure	400					{object}	ErrorResponse	"error: Bad request"
//	@Failure	404					{object}	ErrorResponse	"error: Not found"
//	@Failure	500					{object}	ErrorResponse	"error: Server error"
//	@Router		/blobs/feed [get]
func (s *ServerV2) FetchBlobFeedHandler(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("FetchBlobFeedHandler", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	var err error

	now := time.Now()
	oldestTime := now.Add(-maxBlobAge)

	endTime := now
	if c.Query("end") != "" {
		endTime, err = time.Parse("2006-01-02T15:04:05Z", c.Query("end"))
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeedHandler")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse end param: %w", err))
			return
		}
		if endTime.Before(oldestTime) {
			s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeedHandler")
			invalidParamsErrorResponse(c, fmt.Errorf("end time cannot be more than 14 days in the past, found: %s", c.Query("end")))
			return
		}
	}

	interval := 3600
	if c.Query("interval") != "" {
		interval, err = strconv.Atoi(c.Query("interval"))
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeedHandler")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse interval param: %w", err))
			return
		}
		if interval <= 0 {
			s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeedHandler")
			invalidParamsErrorResponse(c, fmt.Errorf("interval must be greater than 0, found: %d", interval))
			return
		}
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeedHandler")
		invalidParamsErrorResponse(c, fmt.Errorf("failed to parse limit param: %w", err))
		return
	}
	if limit <= 0 || limit > maxNumBlobsPerBlobFeedResponse {
		limit = maxNumBlobsPerBlobFeedResponse
	}

	paginationCursor := blobstore.BlobFeedCursor{
		RequestedAt: 0,
	}
	if c.Query("pagination_token") != "" {
		cursor, err := paginationCursor.FromCursorKey(c.Query("pagination_token"))
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeedHandler")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse the pagination token: %w", err))
			return
		}
		paginationCursor = *cursor
	}

	startTime := endTime.Add(-time.Duration(interval) * time.Second)
	if startTime.Before(oldestTime) {
		startTime = oldestTime
	}
	startCursor := blobstore.BlobFeedCursor{
		RequestedAt: uint64(startTime.UnixNano()),
	}
	if startCursor.LessThan(&paginationCursor) {
		startCursor = paginationCursor
	}
	endCursor := blobstore.BlobFeedCursor{
		RequestedAt: uint64(endTime.UnixNano()),
	}

	blobs, paginationToken, err := s.blobMetadataStore.GetBlobMetadataByRequestedAt(c.Request.Context(), startCursor, endCursor, limit)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBlobFeedHandler")
		errorResponse(c, fmt.Errorf("failed to fetch feed from blob metadata store: %w", err))
		return
	}

	token := ""
	if paginationToken != nil {
		token = paginationToken.ToCursorKey()
	}
	blobInfo := make([]BlobInfo, len(blobs))
	for i := 0; i < len(blobs); i++ {
		bk, err := blobs[i].BlobHeader.BlobKey()
		if err != nil {
			s.metrics.IncrementFailedRequestNum("FetchBlobFeedHandler")
			errorResponse(c, fmt.Errorf("failed to serialize blob key: %w", err))
			return
		}
		blobInfo[i].BlobKey = bk.Hex()
		blobInfo[i].BlobMetadata = blobs[i]
	}
	response := &BlobFeedResponse{
		Blobs:           blobInfo,
		PaginationToken: token,
	}
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxFeedBlobAge))
	s.metrics.IncrementSuccessfulRequestNum("FetchBlobFeedHandler")
	c.JSON(http.StatusOK, response)
}

// FetchBlobHandler godoc
//
//	@Summary	Fetch blob metadata by blob key
//	@Tags		Blobs
//	@Produce	json
//	@Param		blob_key	path		string	true	"Blob key in hex string"
//	@Success	200			{object}	BlobResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/blobs/{blob_key} [get]
func (s *ServerV2) FetchBlobHandler(c *gin.Context) {
	start := time.Now()
	blobKey, err := corev2.HexToBlobKey(c.Param("blob_key"))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchBlob")
		errorResponse(c, err)
		return
	}
	metadata, err := s.blobMetadataStore.GetBlobMetadata(c.Request.Context(), blobKey)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBlob")
		errorResponse(c, err)
		return
	}
	bk, err := metadata.BlobHeader.BlobKey()
	if err != nil || bk != blobKey {
		s.metrics.IncrementFailedRequestNum("FetchBlob")
		errorResponse(c, err)
		return
	}
	response := &BlobResponse{
		BlobKey:       bk.Hex(),
		BlobHeader:    metadata.BlobHeader,
		Status:        metadata.BlobStatus.String(),
		DispersedAt:   metadata.RequestedAt,
		BlobSizeBytes: metadata.BlobSize,
	}
	s.metrics.IncrementSuccessfulRequestNum("FetchBlob")
	s.metrics.ObserveLatency("FetchBlob", float64(time.Since(start).Milliseconds()))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxFeedBlobAge))
	c.JSON(http.StatusOK, response)
}

// FetchBlobCertificateHandler godoc
//
//	@Summary	Fetch blob certificate by blob key v2
//	@Tags		Blobs
//	@Produce	json
//	@Param		blob_key	path		string	true	"Blob key in hex string"
//	@Success	200			{object}	BlobCertificateResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/blobs/{blob_key}/certificate [get]
func (s *ServerV2) FetchBlobCertificateHandler(c *gin.Context) {
	start := time.Now()
	blobKey, err := corev2.HexToBlobKey(c.Param("blob_key"))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchBlobCertificate")
		errorResponse(c, err)
		return
	}
	cert, _, err := s.blobMetadataStore.GetBlobCertificate(c.Request.Context(), blobKey)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBlobCertificate")
		errorResponse(c, err)
		return
	}
	response := &BlobCertificateResponse{
		Certificate: cert,
	}
	s.metrics.IncrementSuccessfulRequestNum("FetchBlobCertificate")
	s.metrics.ObserveLatency("FetchBlobCertificate", float64(time.Since(start).Milliseconds()))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxFeedBlobAge))
	c.JSON(http.StatusOK, response)
}

// FetchBlobAttestationInfo godoc
//
//	@Summary	Fetch attestation info for a blob
//	@Tags		Blobs
//	@Produce	json
//	@Param		blob_key	path		string	true	"Blob key in hex string"
//	@Success	200			{object}	BlobAttestationInfoResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/blobs/{blob_key}/attestation-info [get]
func (s *ServerV2) FetchBlobAttestationInfo(c *gin.Context) {
	start := time.Now()
	blobKey, err := corev2.HexToBlobKey(c.Param("blob_key"))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchBlobAttestationInfo")
		invalidParamsErrorResponse(c, fmt.Errorf("failed to parse blob_key param: %w", err))
		return
	}

	attestationInfo, err := s.blobMetadataStore.GetBlobAttestationInfo(c.Request.Context(), blobKey)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBlobAttestationInfo")
		errorResponse(c, fmt.Errorf("failed to fetch blob attestation info: %w", err))
		return
	}

	batchHeaderHash, err := attestationInfo.InclusionInfo.BatchHeader.Hash()
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBlobAttestationInfo")
		errorResponse(c, fmt.Errorf("failed to get batch header hash from blob inclusion info: %w", err))
		return
	}

	response := &BlobAttestationInfoResponse{
		BlobKey:         blobKey.Hex(),
		BatchHeaderHash: hex.EncodeToString(batchHeaderHash[:]),
		InclusionInfo:   attestationInfo.InclusionInfo,
		Attestation:     attestationInfo.Attestation,
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchBlobAttestationInfo")
	s.metrics.ObserveLatency("FetchBlobAttestationInfo", float64(time.Since(start).Milliseconds()))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxFeedBlobAge))
	c.JSON(http.StatusOK, response)
}

// FetchOperatorSigningInfo godoc
//
//	@Summary	Fetch operators signing info
//	@Tags		Operators
//	@Produce	json
//	@Param		end				query		string	false	"Fetch operators signing info up to the end time (ISO 8601 format: 2006-01-02T15:04:05Z) [default: now]"
//	@Param		interval		query		int		false	"Fetch operators signing info starting from an interval (in seconds) before the end time [default: 3600]"
//	@Param		quorums			query		string	false	"Comma separated list of quorum IDs to fetch signing info for [default: 0,1]"
//	@Param		nonsigner_only	query		boolean	false	"Whether to only return operators with signing rate less than 100% [default: false]"
//	@Success	200				{object}	OperatorsSigningInfoResponse
//	@Failure	400				{object}	ErrorResponse	"error: Bad request"
//	@Failure	404				{object}	ErrorResponse	"error: Not found"
//	@Failure	500				{object}	ErrorResponse	"error: Server error"
//	@Router		/operators/signing-info [get]
func (s *ServerV2) FetchOperatorSigningInfo(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("FetchOperatorSigningInfo", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	var err error

	now := time.Now()
	oldestTime := now.Add(-maxBlobAge)

	endTime := now
	if c.Query("end") != "" {
		endTime, err = time.Parse("2006-01-02T15:04:05Z", c.Query("end"))
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchOperatorSigningInfo")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse end param: %w", err))
			return
		}
		if endTime.Before(oldestTime) {
			s.metrics.IncrementInvalidArgRequestNum("FetchOperatorSigningInfo")
			invalidParamsErrorResponse(
				c, fmt.Errorf("end time cannot be more than 14 days in the past, found: %s", c.Query("end")),
			)
			return
		}
	}

	interval := 3600
	if c.Query("interval") != "" {
		interval, err = strconv.Atoi(c.Query("interval"))
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchOperatorSigningInfo")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse interval param: %w", err))
			return
		}
		if interval <= 0 {
			s.metrics.IncrementInvalidArgRequestNum("FetchOperatorSigningInfo")
			invalidParamsErrorResponse(c, fmt.Errorf("interval must be greater than 0, found: %d", interval))
			return
		}
	}

	quorumStr := "0,1"
	if c.Query("quorums") != "" {
		quorumStr = c.Query("quorums")
	}
	quorums := strings.Split(quorumStr, ",")
	quorumsSeen := make(map[uint8]struct{}, 0)
	for _, idStr := range quorums {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchOperatorSigningInfo")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse the provided quorum: %s", quorumStr))
			return
		}
		if id < 0 || id > maxQuorumIDAllowed {
			s.metrics.IncrementInvalidArgRequestNum("FetchOperatorSigningInfo")
			invalidParamsErrorResponse(
				c, fmt.Errorf("the quorumID must be in range [0, %d], found: %d", maxQuorumIDAllowed, id),
			)
			return
		}
		quorumsSeen[uint8(id)] = struct{}{}
	}
	quorumIds := make([]uint8, 0, len(quorumsSeen))
	for q := range quorumsSeen {
		quorumIds = append(quorumIds, q)
	}

	nonsignerOnly := false
	if c.Query("nonsigner_only") != "" {
		nonsignerOnlyStr := c.Query("nonsigner_only")
		nonsignerOnly, err = strconv.ParseBool(nonsignerOnlyStr)
		if err != nil {
			invalidParamsErrorResponse(c, errors.New("the nonsigner_only param must be \"true\" or \"false\""))
			return
		}
	}

	startTime := endTime.Add(-time.Duration(interval) * time.Second)
	if startTime.Before(oldestTime) {
		startTime = oldestTime
	}

	attestations, err := s.blobMetadataStore.GetAttestationByAttestedAt(
		c.Request.Context(), uint64(startTime.UnixNano())+1, uint64(endTime.UnixNano()), -1,
	)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchOperatorSigningInfo")
		errorResponse(c, fmt.Errorf("failed to fetch attestation feed from blob metadata store: %w", err))
		return
	}

	signingInfo, err := s.computeOperatorsSigningInfo(c.Request.Context(), attestations, quorumIds, nonsignerOnly)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchOperatorSigningInfo")
		errorResponse(c, fmt.Errorf("failed to compute the operators signing info: %w", err))
		return
	}
	response := OperatorsSigningInfoResponse{
		OperatorSigningInfo: signingInfo,
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchOperatorSigningInfo")
	c.JSON(http.StatusOK, response)
}

// FetchBatchFeedHandler godoc
//
//	@Summary	Fetch batch feed
//	@Tags		Batches
//	@Produce	json
//	@Param		end			query		string	false	"Fetch batches up to the end time (ISO 8601 format: 2006-01-02T15:04:05Z) [default: now]"
//	@Param		interval	query		int		false	"Fetch batches starting from an interval (in seconds) before the end time [default: 3600]"
//	@Param		limit		query		int		false	"The maximum number of batches to fetch. System max (1000) if limit <= 0 [default: 20; max: 1000]"
//	@Success	200			{object}	BatchFeedResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/batches/feed [get]
func (s *ServerV2) FetchBatchFeedHandler(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("FetchBatchFeedHandler", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	var err error

	now := time.Now()
	oldestTime := now.Add(-maxBlobAge)

	endTime := now
	if c.Query("end") != "" {
		endTime, err = time.Parse("2006-01-02T15:04:05Z", c.Query("end"))
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchBatchFeedHandler")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse end param: %w", err))
			return
		}
		if endTime.Before(oldestTime) {
			s.metrics.IncrementInvalidArgRequestNum("FetchBatchFeedHandler")
			invalidParamsErrorResponse(c, fmt.Errorf("end time cannot be more than 14 days in the past, found: %s", c.Query("end")))
			return
		}
	}

	interval := 3600
	if c.Query("interval") != "" {
		interval, err = strconv.Atoi(c.Query("interval"))
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchBatchFeedHandler")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse interval param: %w", err))
			return
		}
		if interval <= 0 {
			s.metrics.IncrementInvalidArgRequestNum("FetchBatchFeedHandler")
			invalidParamsErrorResponse(c, fmt.Errorf("interval must be greater than 0, found: %d", interval))
			return
		}
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchBatchFeedHandler")
		invalidParamsErrorResponse(c, fmt.Errorf("failed to parse limit param: %w", err))
		return
	}
	if limit <= 0 || limit > maxNumBatchesPerBatchFeedResponse {
		limit = maxNumBatchesPerBatchFeedResponse
	}

	startTime := endTime.Add(-time.Duration(interval) * time.Second)
	attestations, err := s.blobMetadataStore.GetAttestationByAttestedAt(c.Request.Context(), uint64(startTime.UnixNano())+1, uint64(endTime.UnixNano()), limit)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBatchFeedHandler")
		errorResponse(c, fmt.Errorf("failed to fetch feed from blob metadata store: %w", err))
		return
	}

	batches := make([]*BatchInfo, len(attestations))
	for i, at := range attestations {
		batchHeaderHash, err := at.BatchHeader.Hash()
		if err != nil {
			s.metrics.IncrementFailedRequestNum("FetchBatchFeedHandler")
			errorResponse(c, fmt.Errorf("failed to compute batch header hash from batch header: %w", err))
			return
		}

		batches[i] = &BatchInfo{
			BatchHeaderHash:         hex.EncodeToString(batchHeaderHash[:]),
			BatchHeader:             at.BatchHeader,
			AttestedAt:              at.AttestedAt,
			AggregatedSignature:     at.Sigma,
			QuorumNumbers:           at.QuorumNumbers,
			QuorumSignedPercentages: at.QuorumResults,
		}
	}
	response := &BatchFeedResponse{
		Batches: batches,
	}
	s.metrics.IncrementSuccessfulRequestNum("FetchBatchFeedHandler")
	c.JSON(http.StatusOK, response)
}

// FetchBatchHandler godoc
//
//	@Summary	Fetch batch by the batch header hash
//	@Tags		Batches
//	@Produce	json
//	@Param		batch_header_hash	path		string	true	"Batch header hash in hex string"
//	@Success	200					{object}	BatchResponse
//	@Failure	400					{object}	ErrorResponse	"error: Bad request"
//	@Failure	404					{object}	ErrorResponse	"error: Not found"
//	@Failure	500					{object}	ErrorResponse	"error: Server error"
//	@Router		/batches/{batch_header_hash} [get]
func (s *ServerV2) FetchBatchHandler(c *gin.Context) {
	start := time.Now()
	batchHeaderHashHex := c.Param("batch_header_hash")
	batchHeaderHash, err := dataapi.ConvertHexadecimalToBytes([]byte(batchHeaderHashHex))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchBatch")
		errorResponse(c, errors.New("invalid batch header hash"))
		return
	}
	batchHeader, attestation, err := s.blobMetadataStore.GetSignedBatch(c.Request.Context(), batchHeaderHash)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBatch")
		errorResponse(c, err)
		return
	}
	// TODO: support fetch of blob inclusion info
	batchResponse := &BatchResponse{
		BatchHeaderHash: batchHeaderHashHex,
		SignedBatch: &SignedBatch{
			BatchHeader: batchHeader,
			Attestation: attestation,
		},
	}
	s.metrics.IncrementSuccessfulRequestNum("FetchBatch")
	s.metrics.ObserveLatency("FetchBatch", float64(time.Since(start).Milliseconds()))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxFeedBlobAge))
	c.JSON(http.StatusOK, batchResponse)
}

// FetchOperatorsStake godoc
//
//	@Summary	Operator stake distribution query
//	@Tags		Operators
//	@Produce	json
//	@Param		operator_id	query		string	false	"Operator ID in hex string [default: all operators if unspecified]"
//	@Success	200			{object}	OperatorsStakeResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/operators/stake [get]
func (s *ServerV2) FetchOperatorsStake(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("FetchOperatorsStake", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	operatorId := c.DefaultQuery("operator_id", "")
	s.logger.Info("getting operators stake distribution", "operatorId", operatorId)

	operatorsStakeResponse, err := s.operatorHandler.GetOperatorsStake(c.Request.Context(), operatorId)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchOperatorsStake")
		errorResponse(c, fmt.Errorf("failed to get operator stake - %s", err))
		return
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchOperatorsStake")
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxOperatorsStakeAge))
	c.JSON(http.StatusOK, operatorsStakeResponse)
}

// FetchOperatorsNodeInfo godoc
//
//	@Summary	Active operator semver
//	@Tags		Operators
//	@Produce	json
//	@Success	200	{object}	SemverReportResponse
//	@Failure	500	{object}	ErrorResponse	"error: Server error"
//	@Router		/operators/nodeinfo [get]
func (s *ServerV2) FetchOperatorsNodeInfo(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("FetchOperatorsNodeInfo", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	report, err := s.operatorHandler.ScanOperatorsHostInfo(c.Request.Context())
	if err != nil {
		s.logger.Error("failed to scan operators host info", "error", err)
		s.metrics.IncrementFailedRequestNum("FetchOperatorsNodeInfo")
		errorResponse(c, err)
	}
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxOperatorPortCheckAge))
	c.JSON(http.StatusOK, report)
}

// FetchOperatorsResponses godoc
//
//	@Summary	Fetch operator attestation response for a batch
//	@Tags		Operators
//	@Produce	json
//	@Param		batch_header_hash	path		string	true	"Batch header hash in hex string"
//	@Param		operator_id			query		string	false	"Operator ID in hex string [default: all operators if unspecified]"
//	@Success	200					{object}	OperatorDispersalResponses
//	@Failure	400					{object}	ErrorResponse	"error: Bad request"
//	@Failure	404					{object}	ErrorResponse	"error: Not found"
//	@Failure	500					{object}	ErrorResponse	"error: Server error"
//	@Router		/operators/{batch_header_hash} [get]
func (s *ServerV2) FetchOperatorsResponses(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("FetchOperatorsResponses", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	batchHeaderHashHex := c.Param("batch_header_hash")
	batchHeaderHash, err := dataapi.ConvertHexadecimalToBytes([]byte(batchHeaderHashHex))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchOperatorsResponses")
		errorResponse(c, errors.New("invalid batch header hash"))
		return
	}
	operatorIdStr := c.DefaultQuery("operator_id", "")

	operatorResponses := make([]*corev2.DispersalResponse, 0)
	if operatorIdStr == "" {
		res, err := s.blobMetadataStore.GetDispersalResponses(c.Request.Context(), batchHeaderHash)
		if err != nil {
			s.metrics.IncrementFailedRequestNum("FetchOperatorsResponses")
			errorResponse(c, err)
			return
		}
		operatorResponses = append(operatorResponses, res...)
	} else {
		operatorId, err := core.OperatorIDFromHex(operatorIdStr)
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchOperatorsResponses")
			errorResponse(c, errors.New("invalid operatorId"))
			return
		}

		res, err := s.blobMetadataStore.GetDispersalResponse(c.Request.Context(), batchHeaderHash, operatorId)
		if err != nil {
			s.metrics.IncrementFailedRequestNum("FetchOperatorsResponses")
			errorResponse(c, err)
			return
		}
		operatorResponses = append(operatorResponses, res)
	}
	response := &OperatorDispersalResponses{
		Responses: operatorResponses,
	}
	s.metrics.IncrementSuccessfulRequestNum("FetchOperatorsResponses")
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxOperatorResponseAge))
	c.JSON(http.StatusOK, response)
}

// CheckOperatorsReachability godoc
//
//	@Summary	Operator v2 node reachability check
//	@Tags		Operators
//	@Produce	json
//	@Param		operator_id	query		string	false	"Operator ID in hex string [default: all operators if unspecified]"
//	@Success	200			{object}	OperatorPortCheckResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/operators/reachability [get]
func (s *ServerV2) CheckOperatorsReachability(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("OperatorPortCheck", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	operatorId := c.DefaultQuery("operator_id", "")
	s.logger.Info("checking operator ports", "operatorId", operatorId)
	portCheckResponse, err := s.operatorHandler.ProbeV2OperatorPorts(c.Request.Context(), operatorId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			err = errNotFound
			s.logger.Warn("operator not found", "operatorId", operatorId)
			s.metrics.IncrementNotFoundRequestNum("OperatorPortCheck")
		} else {
			s.logger.Error("operator port check failed", "error", err)
			s.metrics.IncrementFailedRequestNum("OperatorPortCheck")
		}
		errorResponse(c, err)
		return
	}
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxOperatorPortCheckAge))
	c.JSON(http.StatusOK, portCheckResponse)
}

// FetchMetricsSummaryHandler godoc
//
//	@Summary	Fetch metrics summary
//	@Tags		Metrics
//	@Produce	json
//	@Param		start	query		int	false	"Start unix timestamp [default: 1 hour ago]"
//	@Param		end		query		int	false	"End unix timestamp [default: unix time now]"
//	@Success	200		{object}	Metric
//	@Failure	400		{object}	ErrorResponse	"error: Bad request"
//	@Failure	404		{object}	ErrorResponse	"error: Not found"
//	@Failure	500		{object}	ErrorResponse	"error: Server error"
//	@Router		/metrics/summary  [get]
func (s *ServerV2) FetchMetricsSummaryHandler(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("FetchMetricsSummary", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	now := time.Now()
	start, err := strconv.ParseInt(c.DefaultQuery("start", "0"), 10, 64)
	if err != nil || start == 0 {
		start = now.Add(-time.Hour * 1).Unix()
	}

	end, err := strconv.ParseInt(c.DefaultQuery("end", "0"), 10, 64)
	if err != nil || end == 0 {
		end = now.Unix()
	}

	avgThroughput, err := s.metricsHandler.GetAvgThroughput(c.Request.Context(), start, end)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchMetricsSummary")
		errorResponse(c, err)
		return
	}

	metricSummary := &MetricSummary{
		AvgThroughput: avgThroughput,
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchMetricsSummary")
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxMetricAge))
	c.JSON(http.StatusOK, metricSummary)
}

// FetchMetricsThroughputTimeseriesHandler godoc
//
//	@Summary	Fetch throughput time series
//	@Tags		Metrics
//	@Produce	json
//	@Param		start	query		int	false	"Start unix timestamp [default: 1 hour ago]"
//	@Param		end		query		int	false	"End unix timestamp [default: unix time now]"
//	@Success	200		{object}	[]Throughput
//	@Failure	400		{object}	ErrorResponse	"error: Bad request"
//	@Failure	404		{object}	ErrorResponse	"error: Not found"
//	@Failure	500		{object}	ErrorResponse	"error: Server error"
//	@Router		/metrics/timeseries/throughput  [get]
func (s *ServerV2) FetchMetricsThroughputTimeseriesHandler(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("FetchMetricsThroughputTimeseriesHandler", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	now := time.Now()
	start, err := strconv.ParseInt(c.DefaultQuery("start", "0"), 10, 64)
	if err != nil || start == 0 {
		start = now.Add(-time.Hour * 1).Unix()
	}

	end, err := strconv.ParseInt(c.DefaultQuery("end", "0"), 10, 64)
	if err != nil || end == 0 {
		end = now.Unix()
	}

	ths, err := s.metricsHandler.GetThroughputTimeseries(c.Request.Context(), start, end)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchMetricsThroughputTimeseriesHandler")
		errorResponse(c, err)
		return
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchMetricsThroughputTimeseriesHandler")
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxThroughputAge))
	c.JSON(http.StatusOK, ths)
}

func (s *ServerV2) computeOperatorsSigningInfo(
	ctx context.Context,
	attestations []*corev2.Attestation,
	quorumIDs []uint8,
	nonsignerOnly bool,
) ([]*OperatorSigningInfo, error) {
	if len(attestations) == 0 {
		return nil, errors.New("no attestations to compute signing info")
	}

	// Compute the block number range [startBlock, endBlock] (both inclusive) when the
	// attestations have happened.
	startBlock := attestations[0].ReferenceBlockNumber
	endBlock := attestations[0].ReferenceBlockNumber
	for i := range attestations {
		if startBlock > attestations[i].ReferenceBlockNumber {
			startBlock = attestations[i].ReferenceBlockNumber
		}
		if endBlock < attestations[i].ReferenceBlockNumber {
			endBlock = attestations[i].ReferenceBlockNumber
		}
	}

	// Get quorum change events in range [startBlock+1, endBlock].
	// We don't need the events at startBlock because we'll fetch all active operators and
	// quorums at startBlock.
	operatorQuorumEvents, err := s.subgraphClient.QueryOperatorQuorumEvent(ctx, uint32(startBlock+1), uint32(endBlock))
	if err != nil {
		return nil, err
	}

	// Get operators of interest to compute signing info, which includes:
	// - operators that were active at startBlock
	// - operators that joined after startBlock
	operatorList, err := s.getOperatorsOfInterest(
		ctx, startBlock, endBlock, quorumIDs, operatorQuorumEvents,
	)
	if err != nil {
		return nil, err
	}

	// Create operators' quorum intervals: OperatorQuorumIntervals[op][q] is a sequence of
	// increasing and non-overlapping block intervals during which the operator "op" is
	// registered in quorum "q".
	operatorQuorumIntervals, _, err := s.operatorHandler.CreateOperatorQuorumIntervals(
		ctx, operatorList, operatorQuorumEvents, uint32(startBlock), uint32(endBlock),
	)
	if err != nil {
		return nil, err
	}

	// Compute num batches failed, where numFailed[op][q] is the number of batches
	// failed to sign for quorum "q" by operator "op".
	numFailed := computeNumFailed(attestations, operatorQuorumIntervals)

	// Compute num batches responsible, where numResponsible[op][q] is the number of batches
	// that operator "op" are responsible for in quorum "q".
	numResponsible := computeNumResponsible(attestations, operatorQuorumIntervals)

	totalNumBatchesPerQuorum := computeTotalNumBatchesPerQuorum(attestations)

	state, err := s.chainState.GetOperatorState(ctx, uint(endBlock), quorumIDs)
	if err != nil {
		return nil, err
	}
	signingInfo := make([]*OperatorSigningInfo, 0)
	for _, op := range operatorList.GetOperatorIds() {
		for _, q := range quorumIDs {
			operatorId := op.Hex()

			numShouldHaveSigned := 0
			if num, exist := safeAccess(numResponsible, operatorId, q); exist {
				numShouldHaveSigned = num
			}
			// The operator op received no batch that it should sign.
			if numShouldHaveSigned == 0 {
				continue
			}

			numFailedToSign := 0
			if num, exist := safeAccess(numFailed, operatorId, q); exist {
				numFailedToSign = num
			}

			if nonsignerOnly && numFailedToSign == 0 {
				continue
			}

			operatorAddress, ok := operatorList.GetAddress(operatorId)
			if !ok {
				// This should never happen (becuase OperatorList ensures the 1:1 mapping
				// between ID and address), but we don't fail the entire request, just
				// mark internal error for the address field to signal the issue.
				operatorAddress = "Unexpected internal error"
				s.logger.Error("Internal error: failed to find address for operatorId", "operatorId", operatorId)
			}

			// Signing percentage with 2 decimal (e.g. 95.75, which means 95.75%)
			signingPercentage := math.Round(
				(float64(numShouldHaveSigned-numFailedToSign)/float64(numShouldHaveSigned))*100*100,
			) / 100

			stakePercentage := float64(0)
			if stake, ok := state.Operators[q][op]; ok {
				totalStake := new(big.Float).SetInt(state.Totals[q].Stake)
				stakePercentage, _ = new(big.Float).Quo(
					new(big.Float).SetInt(stake.Stake),
					totalStake).Float64()
			}

			si := &OperatorSigningInfo{
				OperatorId:              operatorId,
				OperatorAddress:         operatorAddress,
				QuorumId:                q,
				TotalUnsignedBatches:    numFailedToSign,
				TotalResponsibleBatches: numShouldHaveSigned,
				TotalBatches:            totalNumBatchesPerQuorum[q],
				SigningPercentage:       signingPercentage,
				StakePercentage:         stakePercentage,
			}
			signingInfo = append(signingInfo, si)
		}
	}

	// Sort by descending order of signing rate and then ascending order of <quorumId, operatorId>.
	sort.Slice(signingInfo, func(i, j int) bool {
		if signingInfo[i].SigningPercentage == signingInfo[j].SigningPercentage {
			if signingInfo[i].OperatorId == signingInfo[j].OperatorId {
				return signingInfo[i].QuorumId < signingInfo[j].QuorumId
			}
			return signingInfo[i].OperatorId < signingInfo[j].OperatorId
		}
		return signingInfo[i].SigningPercentage > signingInfo[j].SigningPercentage
	})

	return signingInfo, nil
}

// getOperatorsOfInterest returns operators that we want to compute signing info for.
//
// This contains two parts:
// - the operators that were active at the startBlock
// - the operators that joined after startBlock
func (s *ServerV2) getOperatorsOfInterest(
	ctx context.Context,
	startBlock, endBlock uint64,
	quorumIDs []uint8,
	operatorQuorumEvents *dataapi.OperatorQuorumEvents,
) (*dataapi.OperatorList, error) {
	operatorList := dataapi.NewOperatorList()

	// The first part: active operators at startBlock
	operatorsByQuorum, err := s.chainReader.GetOperatorStakesForQuorums(ctx, quorumIDs, uint32(startBlock))
	if err != nil {
		return nil, err
	}
	operatorsSeen := make(map[core.OperatorID]struct{}, 0)
	for _, ops := range operatorsByQuorum {
		for _, op := range ops {
			operatorsSeen[op.OperatorID] = struct{}{}
		}
	}
	operatorIDs := make([]core.OperatorID, 0)
	for id := range operatorsSeen {
		operatorIDs = append(operatorIDs, id)
	}
	// Get the address for the operators.
	// operatorAddresses[i] is the address for operatorIDs[i].
	operatorAddresses, err := s.chainReader.BatchOperatorIDToAddress(ctx, operatorIDs)
	if err != nil {
		return nil, err
	}
	for i := range operatorIDs {
		operatorList.Add(operatorIDs[i], operatorAddresses[i].Hex())
	}

	// The second part: new operators after startBlock.
	newAddresses := make(map[string]struct{}, 0)
	for op := range operatorQuorumEvents.AddedToQuorum {
		if _, exist := operatorList.GetID(op); !exist {
			newAddresses[op] = struct{}{}
		}
	}
	for op := range operatorQuorumEvents.RemovedFromQuorum {
		if _, exist := operatorList.GetID(op); !exist {
			newAddresses[op] = struct{}{}
		}
	}
	addresses := make([]gethcommon.Address, 0, len(newAddresses))
	for addr := range newAddresses {
		addresses = append(addresses, gethcommon.HexToAddress(addr))
	}
	operatorIds, err := s.chainReader.BatchOperatorAddressToID(ctx, addresses)
	if err != nil {
		return nil, err
	}
	// We merge the new operators observed in AddedToQuorum and RemovedFromQuorum
	// into the operator set.
	for i := 0; i < len(operatorIds); i++ {
		operatorList.Add(operatorIds[i], addresses[i].Hex())
	}

	return operatorList, nil
}

func computeNumFailed(
	attestations []*corev2.Attestation,
	operatorQuorumIntervals dataapi.OperatorQuorumIntervals,
) map[string]map[uint8]int {
	numFailed := make(map[string]map[uint8]int)
	for _, at := range attestations {
		for _, pubkey := range at.NonSignerPubKeys {
			op := pubkey.GetOperatorID().Hex()
			// Note: avg number of quorums per operator is a small number, so use brute
			// force here (otherwise, we can create a map to make it more efficient)
			for _, operatorQuorum := range operatorQuorumIntervals.GetQuorums(
				op,
				uint32(at.ReferenceBlockNumber),
			) {
				for _, batchQuorum := range at.QuorumNumbers {
					if operatorQuorum == batchQuorum {
						if _, ok := numFailed[op]; !ok {
							numFailed[op] = make(map[uint8]int)
						}
						numFailed[op][operatorQuorum]++
						break
					}
				}
			}
		}
	}
	return numFailed
}

func computeNumResponsible(
	attestations []*corev2.Attestation,
	operatorQuorumIntervals dataapi.OperatorQuorumIntervals,
) map[string]map[uint8]int {
	// Create quorumBatches, where quorumBatches[q].AccuBatches is the total number of
	// batches in block interval [startBlock, b] for quorum "q".
	quorumBatches := dataapi.CreatQuorumBatches(dataapi.CreateQuorumBatchMapV2(attestations))

	numResponsible := make(map[string]map[uint8]int)
	for op, val := range operatorQuorumIntervals {
		if _, ok := numResponsible[op]; !ok {
			numResponsible[op] = make(map[uint8]int)
		}
		for q, intervals := range val {
			numBatches := 0
			if _, ok := quorumBatches[q]; ok {
				for _, interval := range intervals {
					numBatches += dataapi.ComputeNumBatches(
						quorumBatches[q], interval.StartBlock, interval.EndBlock,
					)
				}
			}
			numResponsible[op][q] = numBatches
		}
	}

	return numResponsible
}

func computeTotalNumBatchesPerQuorum(attestations []*corev2.Attestation) map[uint8]int {
	numBatchesPerQuorum := make(map[uint8]int)
	for _, at := range attestations {
		for _, q := range at.QuorumNumbers {
			numBatchesPerQuorum[q]++
		}
	}
	return numBatchesPerQuorum
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
