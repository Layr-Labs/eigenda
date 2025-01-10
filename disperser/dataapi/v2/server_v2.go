package v2

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/semver"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	docsv2 "github.com/Layr-Labs/eigenda/disperser/dataapi/docs/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

var errNotFound = errors.New("not found")

const (
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

	BlobVerificationInfoResponse struct {
		VerificationInfo *corev2.BlobVerificationInfo `json:"blob_verification_info"`
	}

	BatchResponse struct {
		BatchHeaderHash       string                         `json:"batch_header_hash"`
		SignedBatch           *SignedBatch                   `json:"signed_batch"`
		BlobVerificationInfos []*corev2.BlobVerificationInfo `json:"blob_verification_infos"`
	}

	MetricSummary struct {
		AvgThroughput float64 `json:"avg_throughput"`
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

	OperatorDispersalResponse struct {
		Response *corev2.DispersalResponse `json:"operator_dispersal_response"`
	}

	OperatorPortCheckResponse struct {
		OperatorId      string `json:"operator_id"`
		DispersalSocket string `json:"dispersal_socket"`
		RetrievalSocket string `json:"retrieval_socket"`
		DispersalOnline bool   `json:"dispersal_online"`
		RetrievalOnline bool   `json:"retrieval_online"`
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
	basePath := "/api/v2"
	docsv2.SwaggerInfoV2.BasePath = basePath
	docsv2.SwaggerInfoV2.Host = os.Getenv("SWAGGER_HOST")

	v2 := router.Group(basePath)
	{
		blob := v2.Group("/blob")
		{
			blob.GET("/blobs/feed", s.FetchBlobFeedHandler)
			blob.GET("/blobs/:blob_key", s.FetchBlobHandler)
			blob.GET("/blobs/:blob_key/certificate", s.FetchBlobCertificateHandler)
			blob.GET("/blobs/:blob_key/verification-info", s.FetchBlobVerificationInfoHandler)
		}
		batch := v2.Group("/batch")
		{
			batch.GET("/batches/feed", s.FetchBatchFeedHandler)
			batch.GET("/batches/:batch_header_hash", s.FetchBatchHandler)
		}
		operators := v2.Group("/operators")
		{
			operators.GET("/nonsigners", s.FetchNonSingers)
			operators.GET("/stake", s.FetchOperatorsStake)
			operators.GET("/nodeinfo", s.FetchOperatorsNodeInfo)
			operators.GET("/reachability", s.CheckOperatorsReachability)
			operators.GET("/response/:batch_header_hash", s.FetchOperatorResponse)
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

	config := cors.DefaultConfig()
	config.AllowOrigins = s.allowOrigins
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "HEAD", "OPTIONS"}

	if s.serverMode != gin.ReleaseMode {
		config.AllowOrigins = []string{"*"}
	}
	router.Use(cors.New(config))

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

func (s *ServerV2) FetchBlobFeedHandler(c *gin.Context) {
	errorResponse(c, errors.New("FetchBlobFeedHandler unimplemented"))
}

// FetchBlobHandler godoc
//
//	@Summary	Fetch blob metadata by blob key
//	@Tags		Blob
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
//	@Tags		Blob
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

// FetchBlobVerificationInfoHandler godoc
//
//	@Summary	Fetch blob verification info by blob key and batch header hash
//	@Tags		Blob
//	@Produce	json
//	@Param		blob_key			path		string	true	"Blob key in hex string"
//	@Param		batch_header_hash	path		string	true	"Batch header hash in hex string"
//
//	@Success	200					{object}	BlobVerificationInfoResponse
//	@Failure	400					{object}	ErrorResponse	"error: Bad request"
//	@Failure	404					{object}	ErrorResponse	"error: Not found"
//	@Failure	500					{object}	ErrorResponse	"error: Server error"
//	@Router		/blobs/{blob_key}/verification-info [get]
func (s *ServerV2) FetchBlobVerificationInfoHandler(c *gin.Context) {
	start := time.Now()
	blobKey, err := corev2.HexToBlobKey(c.Param("blob_key"))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchBlobVerificationInfo")
		errorResponse(c, err)
		return
	}
	batchHeaderHashHex := c.Query("batch_header_hash")
	batchHeaderHash, err := dataapi.ConvertHexadecimalToBytes([]byte(batchHeaderHashHex))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchBlobVerificationInfo")
		errorResponse(c, err)
		return
	}
	bvi, err := s.blobMetadataStore.GetBlobVerificationInfo(c.Request.Context(), blobKey, batchHeaderHash)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBlobVerificationInfo")
		errorResponse(c, err)
		return
	}
	response := &BlobVerificationInfoResponse{
		VerificationInfo: bvi,
	}
	s.metrics.IncrementSuccessfulRequestNum("FetchBlobVerificationInfo")
	s.metrics.ObserveLatency("FetchBlobVerificationInfo", float64(time.Since(start).Milliseconds()))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxFeedBlobAge))
	c.JSON(http.StatusOK, response)
}

func (s *ServerV2) FetchBatchFeedHandler(c *gin.Context) {
	errorResponse(c, errors.New("FetchBatchFeedHandler unimplemented"))
}

// FetchBatchHandler godoc
//
//	@Summary	Fetch batch by the batch header hash
//	@Tags		Batch
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
	// TODO: support fetch of blob verification info
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

// FetchOperatorResponse godoc
//
//	@Summary	Fetch operator attestation response for a batch
//	@Tags		Operators
//	@Produce	json
//	@Param		batch_header_hash	path		string	true	"Batch header hash in hex string"
//	@Param		operator_id			query		string	false	"Operator ID in hex string"
//	@Success	200					{object}	OperatorDispersalResponse
//	@Failure	400					{object}	ErrorResponse	"error: Bad request"
//	@Failure	404					{object}	ErrorResponse	"error: Not found"
//	@Failure	500					{object}	ErrorResponse	"error: Server error"
//	@Router		/operators/{batch_header_hash} [get]
func (s *ServerV2) FetchOperatorResponse(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("FetchOperatorResponse", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	batchHeaderHashHex := c.Param("batch_header_hash")
	batchHeaderHash, err := dataapi.ConvertHexadecimalToBytes([]byte(batchHeaderHashHex))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchOperatorResponse")
		errorResponse(c, errors.New("invalid batch header hash"))
		return
	}
	operatorIdStr := c.DefaultQuery("operator_id", "")
	operatorId, err := core.OperatorIDFromHex(operatorIdStr)
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchOperatorResponse")
		errorResponse(c, errors.New("invalid operatorId"))
		return
	}

	operatorResponse, err := s.blobMetadataStore.GetDispersalResponse(c.Request.Context(), batchHeaderHash, operatorId)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchOperatorResponse")
		errorResponse(c, err)
		return
	}
	response := &OperatorDispersalResponse{
		Response: operatorResponse,
	}
	s.metrics.IncrementSuccessfulRequestNum("FetchOperatorResponse")
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxOperatorResponseAge))
	c.JSON(http.StatusOK, response)
}

// CheckOperatorsReachability godoc
//
//	@Summary	Operator node reachability check
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
	portCheckResponse, err := s.operatorHandler.ProbeOperatorHosts(c.Request.Context(), operatorId)
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

func (s *ServerV2) FetchNonSingers(c *gin.Context) {
	errorResponse(c, errors.New("FetchNonSingers unimplemented"))
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
