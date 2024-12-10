package dataapi

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/docs"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

type ServerInterface interface {
	Start() error
	Shutdown() error
}

type serverv2 struct {
	serverMode   string
	socketAddr   string
	allowOrigins []string
	logger       logging.Logger

	blobMetadataStore *blobstore.BlobMetadataStore
	subgraphClient    SubgraphClient
	chainReader       core.Reader
	chainState        core.ChainState
	indexedChainState core.IndexedChainState
	promClient        PrometheusClient
	metrics           *Metrics
}

func NewServerV2(
	config Config,
	blobMetadataStore *blobstore.BlobMetadataStore,
	promClient PrometheusClient,
	subgraphClient SubgraphClient,
	chainReader core.Reader,
	chainState core.ChainState,
	indexedChainState core.IndexedChainState,
	logger logging.Logger,
	metrics *Metrics,
) *serverv2 {
	return &serverv2{
		logger:            logger.With("component", "DataAPIServerV2"),
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
	}
}

func (s *serverv2) Start() error {
	if s.serverMode == gin.ReleaseMode {
		// optimize performance and disable debug features.
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	basePath := "/api/v2"
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Host = os.Getenv("SWAGGER_HOST")
	v2 := router.Group(basePath)
	{
		feed := v2.Group("/feed")
		{
			// Blob feed
			feed.GET("/blobs", s.FetchBlobsHandler)
			feed.GET("/blobs/:blob_key", s.FetchBlobHandler)
			// Batch feed
			feed.GET("/batches", s.FetchBatchesHandler)
			feed.GET("/batches/:batch_header_hash", s.FetchBatchHandler)
		}
		operators := v2.Group("/operators")
		{
			operators.GET("/non-signers", s.FetchNonSingers)
			operators.GET("/stake", s.FetchOperatorsStake)
			operators.GET("/nodeinfo", s.FetchOperatorsNodeInfo)
			operators.GET("/reachability", s.CheckOperatorsReachability)
		}
		metrics := v2.Group("/metrics")
		{
			metrics.GET("/overview", s.FetchMetricsOverviewHandler)
			metrics.GET("/throughput", s.FetchMetricsThroughputHandler)
		}
		swagger := v2.Group("/swagger")
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

func (s *serverv2) Shutdown() error {
	return nil
}

func (s *serverv2) FetchBlobsHandler(c *gin.Context) {
	errorResponse(c, errors.New("FetchBlobsHandler unimplemented"))
}

func (s *serverv2) FetchBlobHandler(c *gin.Context) {
	errorResponse(c, errors.New("FetchBlobHandler unimplemented"))
}

func (s *serverv2) FetchBatchesHandler(c *gin.Context) {
	errorResponse(c, errors.New("FetchBatchesHandler unimplemented"))
}

func (s *serverv2) FetchBatchHandler(c *gin.Context) {
	errorResponse(c, errors.New("FetchBatchHandler unimplemented"))
}

func (s *serverv2) FetchOperatorsStake(c *gin.Context) {
	errorResponse(c, errors.New("FetchOperatorsStake unimplemented"))
}

func (s *serverv2) FetchOperatorsNodeInfo(c *gin.Context) {
	errorResponse(c, errors.New("FetchOperatorsNodeInfo unimplemented"))
}

func (s *serverv2) CheckOperatorsReachability(c *gin.Context) {
	errorResponse(c, errors.New("CheckOperatorsReachability unimplemented"))
}

func (s *serverv2) FetchNonSingers(c *gin.Context) {
	errorResponse(c, errors.New("FetchNonSingers unimplemented"))
}

func (s *serverv2) FetchMetricsOverviewHandler(c *gin.Context) {
	errorResponse(c, errors.New("FetchMetricsOverviewHandler unimplemented"))
}

func (s *serverv2) FetchMetricsThroughputHandler(c *gin.Context) {
	errorResponse(c, errors.New("FetchMetricsThroughputHandler unimplemented"))
}
