package dataapi

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginswagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

const (
	maxWorkerPoolLimit   = 10
	maxQueryBatchesLimit = 2
)

var errNotFound = errors.New("not found")

type EigenDAServiceChecker interface {
	CheckHealth(ctx context.Context, serviceName string) (*grpc_health_v1.HealthCheckResponse, error)
	CloseConnections() error
}

type (
	BlobMetadataResponse struct {
		BlobKey                 string                    `json:"blob_key"`
		BatchHeaderHash         string                    `json:"batch_header_hash"`
		BlobIndex               uint32                    `json:"blob_index"`
		SignatoryRecordHash     string                    `json:"signatory_record_hash"`
		ReferenceBlockNumber    uint32                    `json:"reference_block_number"`
		BatchRoot               string                    `json:"batch_root"`
		BlobInclusionProof      string                    `json:"blob_inclusion_proof"`
		BlobCommitment          *encoding.BlobCommitments `json:"blob_commitment"`
		BatchId                 uint32                    `json:"batch_id"`
		ConfirmationBlockNumber uint32                    `json:"confirmation_block_number"`
		ConfirmationTxnHash     string                    `json:"confirmation_txn_hash"`
		Fee                     string                    `json:"fee"`
		SecurityParams          []*core.SecurityParam     `json:"security_params"`
		RequestAt               uint64                    `json:"requested_at"`
		BlobStatus              disperser.BlobStatus      `json:"blob_status"`
	}

	Metric struct {
		Throughput float64 `json:"throughput"`
		CostInGas  float64 `json:"cost_in_gas"`
		TotalStake uint64  `json:"total_stake"`
	}

	Throughput struct {
		Throughput float64 `json:"throughput"`
		Timestamp  uint64  `json:"timestamp"`
	}

	Meta struct {
		Size int `json:"size"`
	}

	BlobsResponse struct {
		Meta Meta                    `json:"meta"`
		Data []*BlobMetadataResponse `json:"data"`
	}

	OperatorNonsigningPercentageMetrics struct {
		OperatorId           string  `json:"operator_id"`
		QuorumId             uint8   `json:"quorum_id"`
		TotalUnsignedBatches int     `json:"total_unsigned_batches"`
		TotalBatches         int     `json:"total_batches"`
		Percentage           float64 `json:"percentage"`
	}

	OperatorsNonsigningPercentage struct {
		Meta Meta                                   `json:"meta"`
		Data []*OperatorNonsigningPercentageMetrics `json:"data"`
	}

	DeregisteredOperatorMetadata struct {
		OperatorId           string `json:"operator_id"`
		BlockNumber          uint   `json:"block_number"`
		Socket               string `json:"socket"`
		IsOnline             bool   `json:"is_online"`
		OperatorProcessError string `json:"operator_process_error"`
	}

	DeregisteredOperatorsResponse struct {
		Meta Meta                            `json:"meta"`
		Data []*DeregisteredOperatorMetadata `json:"data"`
	}

	ServiceAvailability struct {
		ServiceName   string `json:"service_name"`
		ServiceStatus string `json:"service_status"`
	}

	ServiceAvailabilityResponse struct {
		Meta Meta                   `json:"meta"`
		Data []*ServiceAvailability `json:"data"`
	}

	ErrorResponse struct {
		Error string `json:"error"`
	}

	server struct {
		serverMode     string
		socketAddr     string
		allowOrigins   []string
		logger         logging.Logger
		blobstore      disperser.BlobStore
		promClient     PrometheusClient
		subgraphClient SubgraphClient
		transactor     core.Transactor
		chainState     core.ChainState

		metrics               *Metrics
		disperserHostName     string
		churnerHostName       string
		eigenDAServiceChecker EigenDAServiceChecker
	}
)

func NewServer(
	config Config,
	blobstore disperser.BlobStore,
	promClient PrometheusClient,
	subgraphClient SubgraphClient,
	transactor core.Transactor,
	chainState core.ChainState,
	logger logging.Logger,
	metrics *Metrics,
	grpcConn GRPCConn,
	eigenDAServiceChecker EigenDAServiceChecker,

) *server {
	// Initialize the health checker service for EigenDA services
	if grpcConn == nil {
		grpcConn = &GRPCDialerSkipTLS{}
	}

	if eigenDAServiceChecker == nil {

		eigenDAServiceChecker = NewEigenDAServiceHealthCheck(grpcConn, config.DisperserHostname, config.ChurnerHostname)
	}

	return &server{
		logger:                logger,
		serverMode:            config.ServerMode,
		socketAddr:            config.SocketAddr,
		allowOrigins:          config.AllowOrigins,
		blobstore:             blobstore,
		promClient:            promClient,
		subgraphClient:        subgraphClient,
		transactor:            transactor,
		chainState:            chainState,
		metrics:               metrics,
		disperserHostName:     config.DisperserHostname,
		churnerHostName:       config.ChurnerHostname,
		eigenDAServiceChecker: eigenDAServiceChecker,
	}
}

func (s *server) Start() error {
	if s.serverMode == gin.ReleaseMode {
		// optimize performance and disable debug features.
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	basePath := "/api/v1"
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Host = os.Getenv("SWAGGER_HOST")

	v1 := router.Group(basePath)
	{
		feed := v1.Group("/feed")
		{
			feed.GET("/blobs", s.FetchBlobsHandler)
			feed.GET("/blobs/:blob_key", s.FetchBlobHandler)
		}
		operatorsInfo := v1.Group("/operators-info")
		{
			operatorsInfo.GET("/deregistered-operators", s.FetchDeregisteredOperators)
		}
		serviceAvailability := v1.Group("/eigenda")
		{
			serviceAvailability.GET("/service-availability", s.GetEigenDAServiceAvailability)
		}
		metrics := v1.Group("/metrics")
		{
			metrics.GET("/", s.FetchMetricsHandler)
			metrics.GET("/throughput", s.FetchMetricsTroughputHandler)
			metrics.GET("/non-signers", s.FetchNonSigners)
			metrics.GET("/operator-nonsigning-percentage", s.FetchOperatorsNonsigningPercentageHandler)
		}
		swagger := v1.Group("/swagger")
		{
			swagger.GET("/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
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
	config.AllowMethods = []string{"GET", "HEAD", "OPTIONS"}

	if s.serverMode != gin.ReleaseMode {
		config.AllowOrigins = []string{"*"}
	}
	router.Use(cors.New(config))

	srv := &http.Server{
		Addr:              s.socketAddr,
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	errChan := run(s.logger, srv)
	return <-errChan
}

func (s *server) Shutdown() error {

	if s.eigenDAServiceChecker != nil {
		err := s.eigenDAServiceChecker.CloseConnections()

		if err != nil {
			s.logger.Error("Failed to close connections", "error", err)
			return err
		}
	}

	return nil
}

// FetchBlobHandler godoc
//
//	@Summary	Fetch blob metadata by blob key
//	@Tags		Feed
//	@Produce	json
//	@Param		blob_key	path		string	true	"Blob Key"
//	@Success	200			{object}	BlobMetadataResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/feed/blobs/{blob_key} [get]
func (s *server) FetchBlobHandler(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("FetchBlob", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	blobKey := c.Param("blob_key")

	metadata, err := s.getBlob(c.Request.Context(), blobKey)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBlob")
		errorResponse(c, err)
		return
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchBlob")
	c.JSON(http.StatusOK, metadata)
}

// FetchBlobsHandler godoc
//
//	@Summary	Fetch blobs metadata list
//	@Tags		Feed
//	@Produce	json
//	@Param		limit	query		int	false	"Limit [default: 10]"
//	@Success	200		{object}	BlobsResponse
//	@Failure	400		{object}	ErrorResponse	"error: Bad request"
//	@Failure	404		{object}	ErrorResponse	"error: Not found"
//	@Failure	500		{object}	ErrorResponse	"error: Server error"
//	@Router		/feed/blobs [get]
func (s *server) FetchBlobsHandler(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("FetchBlobs", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	metadatas, err := s.getBlobs(c.Request.Context(), limit)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBlobs")
		errorResponse(c, err)
		return
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchBlobs")
	c.JSON(http.StatusOK, BlobsResponse{
		Meta: Meta{
			Size: len(metadatas),
		},
		Data: metadatas,
	})
}

// FetchMetricsHandler godoc
//
//	@Summary	Fetch metrics
//	@Tags		Metrics
//	@Produce	json
//	@Param		start	query		int	false	"Start unix timestamp [default: 1 hour ago]"
//	@Param		end		query		int	false	"End unix timestamp [default: unix time now]"
//	@Param		limit	query		int	false	"Limit [default: 10]"
//	@Success	200		{object}	Metric
//	@Failure	400		{object}	ErrorResponse	"error: Bad request"
//	@Failure	404		{object}	ErrorResponse	"error: Not found"
//	@Failure	500		{object}	ErrorResponse	"error: Server error"
//	@Router		/metrics  [get]
func (s *server) FetchMetricsHandler(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("FetchMetrics", f*1000) // make milliseconds
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
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit == 0 {
		limit = 10
	}

	metric, err := s.getMetric(c.Request.Context(), start, end, limit)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchMetrics")
		errorResponse(c, err)
		return
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchMetrics")
	c.JSON(http.StatusOK, metric)
}

// FetchMetricsTroughputHandler godoc
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
//	@Router		/metrics/throughput  [get]
func (s *server) FetchMetricsTroughputHandler(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("FetchMetricsTroughput", f*1000) // make milliseconds
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

	ths, err := s.getThroughput(c.Request.Context(), start, end)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchMetricsTroughput")
		errorResponse(c, err)
		return
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchMetricsTroughput")
	c.JSON(http.StatusOK, ths)
}

// FetchNonSigners godoc
//
//	@Summary	Fetch non signers
//	@Tags		Metrics
//	@Produce	json
//	@Param		interval	query		int	false	"Interval to query for non signers in seconds [default: 3600]"
//	@Success	200			{object}	[]NonSigner
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/metrics/non-signers  [get]
func (s *server) FetchNonSigners(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("FetchNonSigners", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	interval, err := strconv.ParseInt(c.DefaultQuery("interval", "3600"), 10, 64)
	if err != nil || interval == 0 {
		interval = 3600
	}
	metric, err := s.getNonSigners(c.Request.Context(), interval)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchNonSigners")
		errorResponse(c, err)
		return
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchNonSigners")
	c.JSON(http.StatusOK, metric)
}

// FetchOperatorsNonsigningPercentageHandler godoc
//
//	@Summary	Fetch operators non signing percentage
//	@Tags		Metrics
//	@Produce	json
//	@Param		interval	query		int	false	"Interval to query for operators nonsigning percentage [default: 3600]"
//	@Success	200			{object}	OperatorsNonsigningPercentage
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/metrics/operator-nonsigning-percentage  [get]
func (s *server) FetchOperatorsNonsigningPercentageHandler(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("FetchOperatorsNonsigningPercentageHandler", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	interval, err := strconv.ParseInt(c.DefaultQuery("interval", "3600"), 10, 64)
	if err != nil || interval == 0 {
		interval = 3600
	}
	metric, err := s.getOperatorNonsigningRate(c.Request.Context(), interval)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchOperatorsNonsigningPercentageHandler")
		errorResponse(c, err)
		return
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchOperatorsNonsigningPercentageHandler")
	c.JSON(http.StatusOK, metric)
}

// FetchDeregisteredOperators godoc
//
//	@Summary	Fetch list of operators that have been deregistered for days. Days is a query parameter with a default value of 14 and max value of 30.
//	@Tags		OperatorsInfo
//	@Produce	json
//	@Success	200	{object}	DeregisteredOperatorsResponse
//	@Failure	400	{object}	ErrorResponse	"error: Bad request"
//	@Failure	404	{object}	ErrorResponse	"error: Not found"
//	@Failure	500	{object}	ErrorResponse	"error: Server error"
//	@Router		/operators-info/deregistered-operators [get]
func (s *server) FetchDeregisteredOperators(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("FetchDeregisteredOperators", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	// Get query parameters
	// Default Value 14 days
	days := c.DefaultQuery("days", "14") // If not specified, defaults to 14

	// Convert days to integer
	daysInt, err := strconv.Atoi(days)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'days' parameter"})
		return
	}

	if daysInt > 30 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'days' parameter. Max value is 30"})
		return
	}

	operatorMetadatas, err := s.getDeregisteredOperatorForDays(c.Request.Context(), int32(daysInt))
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchDeregisteredOperators")
		errorResponse(c, err)
		return
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchDeregisteredOperators")
	c.JSON(http.StatusOK, DeregisteredOperatorsResponse{
		Meta: Meta{
			Size: len(operatorMetadatas),
		},
		Data: operatorMetadatas,
	})
}

// GetEigenDAServiceAvailability godoc
//
//	@Summary	Get status of public EigenDA services.
//	@Tags		ServiceAvailability
//	@Produce	json
//	@Success	200	{object}	ServiceAvailabilityResponse
//	@Failure	400	{object}	ErrorResponse	"error: Bad request"
//	@Failure	404	{object}	ErrorResponse	"error: Not found"
//	@Failure	500	{object}	ErrorResponse	"error: Server error"
//	@Router		/eigenda/service-availability [get]
func (s *server) GetEigenDAServiceAvailability(c *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("GetEigenDAServiceAvailability", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	// Get query parameters to filter services
	serviceName := c.DefaultQuery("service-name", "") // If not specified, default to return all services

	// If service name is not specified, return all services
	services := []string{}

	if serviceName == "disperser" {
		services = append(services, "Disperser")
	} else if serviceName == "churner" {
		services = append(services, "Churner")
	} else if serviceName == "" {
		services = append(services, "Disperser", "Churner")
	}
	s.logger.Info("Getting service availability for", "services", strings.Join(services, ", "))

	availabilityStatuses, err := s.getServiceAvailability(c.Request.Context(), services)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("GetEigenDAServiceAvailability")
		errorResponse(c, err)
		return
	}

	s.metrics.IncrementSuccessfulRequestNum("GetEigenDAServiceAvailability")

	// Set the status code to 503 if any of the services are not serving
	availabilityStatus := http.StatusOK
	for _, status := range availabilityStatuses {
		if status.ServiceStatus == "NOT_SERVING" {
			availabilityStatus = http.StatusServiceUnavailable
			break
		}

		if status.ServiceStatus == "UNKNOWN" {
			availabilityStatus = http.StatusInternalServerError
			break
		}

	}

	c.JSON(availabilityStatus, ServiceAvailabilityResponse{
		Meta: Meta{
			Size: len(availabilityStatuses),
		},
		Data: availabilityStatuses,
	})
}

func (s *server) getBlobMetadataByBatchesWithLimit(ctx context.Context, limit int) ([]*Batch, []*disperser.BlobMetadata, error) {
	var (
		blobMetadatas   = make([]*disperser.BlobMetadata, 0)
		batches         = make([]*Batch, 0)
		blobKeyPresence = make(map[string]struct{})
		batchPresence   = make(map[string]struct{})
	)

	for skip := 0; len(blobMetadatas) < limit && skip < limit; skip += maxQueryBatchesLimit {
		batchesWithLimit, err := s.subgraphClient.QueryBatchesWithLimit(ctx, maxQueryBatchesLimit, skip)
		if err != nil {
			s.logger.Error("Failed to query batches", "error", err)
			return nil, nil, err
		}

		if len(batchesWithLimit) == 0 {
			break
		}

		for i := range batchesWithLimit {
			s.logger.Debug("Getting blob metadata", "batchHeaderHash", batchesWithLimit[i].BatchHeaderHash)
			var (
				batch = batchesWithLimit[i]
			)
			if batch == nil {
				continue
			}
			batchHeaderHash, err := ConvertHexadecimalToBytes(batch.BatchHeaderHash)
			if err != nil {
				s.logger.Error("Failed to convert batch header hash to hex string", "error", err)
				continue
			}
			batchKey := string(batchHeaderHash[:])
			if _, found := batchPresence[batchKey]; !found {
				batchPresence[batchKey] = struct{}{}
			} else {
				// The batch has processed, skip it.
				s.logger.Error("Getting duplicate batch from the graph", "batch header hash", batchKey)
				continue
			}

			metadatas, err := s.blobstore.GetAllBlobMetadataByBatch(ctx, batchHeaderHash)
			if err != nil {
				s.logger.Error("Failed to get blob metadata", "error", err)
				continue
			}
			for _, bm := range metadatas {
				blobKey := bm.GetBlobKey().String()
				if _, found := blobKeyPresence[blobKey]; !found {
					blobKeyPresence[blobKey] = struct{}{}
					blobMetadatas = append(blobMetadatas, bm)
				} else {
					s.logger.Error("Getting duplicate blob key from the blobstore", "blobkey", blobKey)
				}
			}
			batches = append(batches, batch)
			if len(blobMetadatas) >= limit {
				break
			}
		}
	}

	if len(blobMetadatas) >= limit {
		blobMetadatas = blobMetadatas[:limit]
	}

	return batches, blobMetadatas, nil
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
		logger.Info("server running", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	return errChan
}
