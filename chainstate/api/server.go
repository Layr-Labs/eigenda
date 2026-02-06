package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda/chainstate"
	"github.com/Layr-Labs/eigenda/chainstate/store"
	"github.com/Layr-Labs/eigenda/chainstate/types"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server is the HTTP API server for querying indexed chain state.
type Server struct {
	config *chainstate.IndexerConfig
	store  store.Store
	router *gin.Engine
	logger logging.Logger
	server *http.Server
}

// NewServer creates a new API server.
func NewServer(
	config *chainstate.IndexerConfig,
	store store.Store,
	logger logging.Logger,
) *Server {
	// Set Gin mode to release by default
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	server := &Server{
		config: config,
		store:  store,
		router: router,
		logger: logger.With("component", "ChainStateAPI"),
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures all API routes.
func (s *Server) setupRoutes() {
	api := s.router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", s.handleHealth)

		// Operator endpoints
		api.GET("/operators", s.handleListOperators)
		api.GET("/operators/:id", s.handleGetOperator)

		// Quorum APK endpoints
		api.GET("/quorum-apk", s.handleGetQuorumAPK)
		api.GET("/quorum-apk/history", s.handleListQuorumAPKs)

		// Ejection endpoints
		api.GET("/ejections", s.handleListEjections)
		api.GET("/ejections/:operator_id", s.handleListOperatorEjections)

		// Socket update endpoints
		api.GET("/socket-updates/:operator_id", s.handleListSocketUpdates)

		// Status endpoint
		api.GET("/status", s.handleStatus)
	}
}

// handleHealth returns the health status of the API server.
func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().UTC(),
	})
}

// handleStatus returns the current indexer status.
func (s *Server) handleStatus(c *gin.Context) {
	lastBlock, err := s.store.GetLastIndexedBlock(c.Request.Context())
	if err != nil {
		s.logger.Error("Failed to get last indexed block", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"last_indexed_block": lastBlock,
		"time":               time.Now().UTC(),
	})
}

// handleListOperators returns a paginated list of operators.
func (s *Server) handleListOperators(c *gin.Context) {
	// Parse query parameters
	registeredOnly := c.Query("registered") == "true"
	deregisteredOnly := c.Query("deregistered") == "true"
	limit := parseIntOr(c.Query("limit"), 100)
	offset := parseIntOr(c.Query("offset"), 0)

	var quorumID *core.QuorumID
	if qidStr := c.Query("quorum_id"); qidStr != "" {
		qid, err := strconv.ParseUint(qidStr, 10, 8)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quorum_id"})
			return
		}
		qidByte := core.QuorumID(qid)
		quorumID = &qidByte
	}

	filter := types.OperatorFilter{
		RegisteredOnly:   registeredOnly,
		DeregisteredOnly: deregisteredOnly,
		QuorumID:         quorumID,
		MinBlock:         parseUint64Or(c.Query("min_block"), 0),
		MaxBlock:         parseUint64Or(c.Query("max_block"), 0),
	}

	operators, err := s.store.ListOperators(c.Request.Context(), filter, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list operators", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list operators"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"operators": operators,
		"count":     len(operators),
		"limit":     limit,
		"offset":    offset,
	})
}

// handleGetOperator returns a single operator by ID.
func (s *Server) handleGetOperator(c *gin.Context) {
	idStr := c.Param("id")
	if len(idStr) != 64 && len(idStr) != 66 { // 32 bytes hex = 64 chars, or 66 with "0x" prefix
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid operator ID format (expected 32-byte hex)"})
		return
	}

	// Remove "0x" prefix if present
	if len(idStr) == 66 && idStr[:2] == "0x" {
		idStr = idStr[2:]
	}

	var operatorID core.OperatorID
	for i := 0; i < 32; i++ {
		_, err := fmt.Sscanf(idStr[i*2:i*2+2], "%02x", &operatorID[i])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid operator ID format"})
			return
		}
	}

	operator, err := s.store.GetOperator(c.Request.Context(), operatorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "operator not found"})
		return
	}

	c.JSON(http.StatusOK, operator)
}

// handleGetQuorumAPK returns the aggregate public key for a quorum at a specific block.
func (s *Server) handleGetQuorumAPK(c *gin.Context) {
	quorumID := parseUint8Or(c.Query("quorum_id"), 0)
	blockNumber := parseUint64Or(c.Query("block_number"), 0)

	if blockNumber == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "block_number is required"})
		return
	}

	apk, err := s.store.GetQuorumAPK(c.Request.Context(), quorumID, blockNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "quorum APK not found"})
		return
	}

	c.JSON(http.StatusOK, apk)
}

// handleListQuorumAPKs returns a list of quorum APK snapshots.
func (s *Server) handleListQuorumAPKs(c *gin.Context) {
	quorumID := parseUint8Or(c.Query("quorum_id"), 0)
	blockNumber := parseUint64Or(c.Query("block_number"), 0)
	minBlock := parseUint64Or(c.Query("min_block"), 0)
	maxBlock := parseUint64Or(c.Query("max_block"), 0)

	filter := types.QuorumAPKFilter{
		QuorumID:    core.QuorumID(quorumID),
		BlockNumber: blockNumber,
		MinBlock:    minBlock,
		MaxBlock:    maxBlock,
	}

	apks, err := s.store.ListQuorumAPKs(c.Request.Context(), filter)
	if err != nil {
		s.logger.Error("Failed to list quorum APKs", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list quorum APKs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"quorum_apks": apks,
		"count":       len(apks),
	})
}

// handleListEjections returns a paginated list of all ejections.
func (s *Server) handleListEjections(c *gin.Context) {
	limit := parseIntOr(c.Query("limit"), 100)
	offset := parseIntOr(c.Query("offset"), 0)

	ejections, err := s.store.ListEjections(c.Request.Context(), nil, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list ejections", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list ejections"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ejections": ejections,
		"count":     len(ejections),
		"limit":     limit,
		"offset":    offset,
	})
}

// handleListOperatorEjections returns ejections for a specific operator.
func (s *Server) handleListOperatorEjections(c *gin.Context) {
	idStr := c.Param("operator_id")
	if len(idStr) != 64 && len(idStr) != 66 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid operator ID format"})
		return
	}

	// Remove "0x" prefix if present
	if len(idStr) == 66 && idStr[:2] == "0x" {
		idStr = idStr[2:]
	}

	var operatorID core.OperatorID
	for i := 0; i < 32; i++ {
		_, err := fmt.Sscanf(idStr[i*2:i*2+2], "%02x", &operatorID[i])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid operator ID format"})
			return
		}
	}

	limit := parseIntOr(c.Query("limit"), 100)
	offset := parseIntOr(c.Query("offset"), 0)

	ejections, err := s.store.ListEjections(c.Request.Context(), &operatorID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list ejections", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list ejections"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ejections": ejections,
		"count":     len(ejections),
		"limit":     limit,
		"offset":    offset,
	})
}

// handleListSocketUpdates returns socket updates for a specific operator.
func (s *Server) handleListSocketUpdates(c *gin.Context) {
	idStr := c.Param("operator_id")
	if len(idStr) != 64 && len(idStr) != 66 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid operator ID format"})
		return
	}

	// Remove "0x" prefix if present
	if len(idStr) == 66 && idStr[:2] == "0x" {
		idStr = idStr[2:]
	}

	var operatorID core.OperatorID
	for i := 0; i < 32; i++ {
		_, err := fmt.Sscanf(idStr[i*2:i*2+2], "%02x", &operatorID[i])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid operator ID format"})
			return
		}
	}

	limit := parseIntOr(c.Query("limit"), 100)
	offset := parseIntOr(c.Query("offset"), 0)

	updates, err := s.store.ListSocketUpdates(c.Request.Context(), operatorID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list socket updates", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list socket updates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"socket_updates": updates,
		"count":          len(updates),
		"limit":          limit,
		"offset":         offset,
	})
}

// Start starts the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf(":%s", s.config.HTTPPort)

	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	// Handle graceful shutdown
	go func() {
		<-ctx.Done()
		s.logger.Info("Shutting down API server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("Failed to shutdown API server gracefully", "error", err)
		}
	}()

	s.logger.Info("Starting HTTP API server", "addr", addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Helper functions for parsing query parameters

func parseIntOr(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}

func parseUint64Or(s string, defaultVal uint64) uint64 {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return defaultVal
	}
	return val
}

func parseUint8Or(s string, defaultVal uint8) uint8 {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return defaultVal
	}
	return uint8(val)
}
