package blobstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common"
	dispersercommon "github.com/Layr-Labs/eigenda/disperser/common"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Ensure PostgresBlobMetadataStore implements MetadataStore
var _ MetadataStore = (*PostgresBlobMetadataStore)(nil)

// PostgresBlobMetadataStore is a blob metadata storage backed by PostgreSQL
type PostgresBlobMetadataStore struct {
	db     *pgxpool.Pool
	logger logging.Logger
}

// GetDB returns the database connection for testing purposes
func (s *PostgresBlobMetadataStore) GetDB() *pgxpool.Pool {
	return s.db
}

// PostgreSQLConfig contains configuration for PostgreSQL connection
type PostgreSQLConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSLMode  string
}

// validateSSLMode validates the SSL mode
func validateSSLMode(mode string) error {
	validModes := map[string]bool{
		"disable":     true,
		"require":     true,
		"verify-ca":   true,
		"verify-full": true,
	}
	if !validModes[mode] {
		return fmt.Errorf("invalid SSL mode: %s. Must be one of: disable, require, verify-ca, verify-full", mode)
	}
	return nil
}

// NewPostgresBlobMetadataStore creates a new PostgresBlobMetadataStore instance
func NewPostgresBlobMetadataStore(config PostgreSQLConfig, logger logging.Logger) (*PostgresBlobMetadataStore, error) {
	// Validate SSL mode
	if err := validateSSLMode(config.SSLMode); err != nil {
		return nil, fmt.Errorf("invalid SSL configuration: %w", err)
	}

	// Create connection string for pgx.ParseConfig
	connStr := fmt.Sprintf(
		"postgres://%s@%s:%d/%s?sslmode=%s",
		config.Username, config.Host, config.Port, config.Database, config.SSLMode,
	)

	// Parse the connection string into a pgx config
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection config: %w", err)
	}

	// Set password from config
	poolConfig.ConnConfig.Password = config.Password

	// Configure connection pool
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 5 * time.Minute

	// Connect to the database
	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Create store
	store := &PostgresBlobMetadataStore{
		db:     pool,
		logger: logger.With("component", "postgresMetadataStore"),
	}

	// Initialize the database tables
	// if err := store.initTables(); err != nil {
	// 	pool.Close() // Close the pool if initialization fails
	// 	return nil, fmt.Errorf("failed to initialize tables: %w", err)
	// }

	return store, nil
}

// initTables creates the necessary tables if they don't exist
func (s *PostgresBlobMetadataStore) initTables() error {
	ctx := context.Background()
	// Create tables for blob metadata
	_, err := s.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS blob_metadata (
			blob_key BYTEA PRIMARY KEY,
			blob_header JSONB NOT NULL,
			blob_size INTEGER NOT NULL,
			signature BYTEA NOT NULL,
			requested_at BIGINT NOT NULL,
			requested_at_bucket BYTEA NOT NULL,
			requested_at_blob_key BYTEA NOT NULL,
			blob_status INTEGER NOT NULL,
			num_retries INTEGER NOT NULL,
			updated_at BIGINT NOT NULL,
			account_id VARCHAR(42) NOT NULL,
			expiry BIGINT NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_blob_metadata_blob_status_updated_at ON blob_metadata (blob_status, updated_at);
		CREATE INDEX IF NOT EXISTS idx_blob_metadata_account_id_requested_at ON blob_metadata (account_id, requested_at);
		CREATE INDEX IF NOT EXISTS idx_blob_metadata_requested_at_bucket_key ON blob_metadata (requested_at_bucket, requested_at_blob_key);
	`)
	if err != nil {
		return fmt.Errorf("failed to create blob_metadata table: %w", err)
	}

	// Create tables for blob certificates
	_, err = s.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS blob_certificates (
			blob_key BYTEA PRIMARY KEY,
			blob_certificate JSONB NOT NULL,
			fragment_info JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create blob_certificates table: %w", err)
	}

	// Create tables for batch headers
	_, err = s.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS batch_headers (
			batch_header_hash BYTEA PRIMARY KEY,
			batch_header JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create batch_headers table: %w", err)
	}

	// Create tables for batches
	_, err = s.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS batches (
			batch_header_hash BYTEA PRIMARY KEY,
			batch_info JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create batches table: %w", err)
	}

	// Create tables for dispersal requests
	_, err = s.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS dispersal_requests (
			batch_header_hash BYTEA NOT NULL,
			operator_id BYTEA NOT NULL,
			dispersal_request JSONB NOT NULL,
			dispersed_at BIGINT NOT NULL,
			PRIMARY KEY (batch_header_hash, operator_id)
		);

		CREATE INDEX IF NOT EXISTS idx_dispersal_requests_operator_dispersed_at ON dispersal_requests (operator_id, dispersed_at);
	`)
	if err != nil {
		return fmt.Errorf("failed to create dispersal_requests table: %w", err)
	}

	// Create tables for dispersal responses
	_, err = s.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS dispersal_responses (
			batch_header_hash BYTEA NOT NULL,
			operator_id BYTEA NOT NULL,
			dispersal_response JSONB NOT NULL,
			responded_at BIGINT NOT NULL,
			PRIMARY KEY (batch_header_hash, operator_id)
		);

		CREATE INDEX IF NOT EXISTS idx_dispersal_responses_operator_responded_at ON dispersal_responses (operator_id, responded_at);
	`)
	if err != nil {
		return fmt.Errorf("failed to create dispersal_responses table: %w", err)
	}

	// Create tables for attestations
	_, err = s.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS attestations (
			batch_header_hash BYTEA PRIMARY KEY,
			attestation JSONB NOT NULL,
			attested_at BIGINT NOT NULL,
			attested_at_bucket VARCHAR(64) NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_attestations_attested_at_bucket_attested_at ON attestations (attested_at_bucket, attested_at);
	`)
	if err != nil {
		return fmt.Errorf("failed to create attestations table: %w", err)
	}

	// Create tables for blob inclusion info
	_, err = s.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS blob_inclusion_info (
			blob_key BYTEA NOT NULL,
			batch_header_hash BYTEA NOT NULL,
			inclusion_info JSONB NOT NULL,
			PRIMARY KEY (blob_key, batch_header_hash)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create blob_inclusion_info table: %w", err)
	}

	return nil
}

// dropTables drops all tables, intended for test environments only
func (s *PostgresBlobMetadataStore) dropTables() error {
	ctx := context.Background()
	_, err := s.db.Exec(ctx, `
		DROP TABLE IF EXISTS blob_metadata CASCADE;
		DROP TABLE IF EXISTS blob_certificates CASCADE;
		DROP TABLE IF EXISTS batch_headers CASCADE;
		DROP TABLE IF EXISTS batches CASCADE;
		DROP TABLE IF EXISTS dispersal_requests CASCADE;
		DROP TABLE IF EXISTS dispersal_responses CASCADE;
		DROP TABLE IF EXISTS attestations CASCADE;
		DROP TABLE IF EXISTS blob_inclusion_info CASCADE;
	`)
	if err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}
	return nil
}

// ResetTables drops and recreates all tables, intended for test environments only
func (s *PostgresBlobMetadataStore) ResetTables() error {
	if err := s.dropTables(); err != nil {
		return err
	}
	if err := s.initTables(); err != nil {
		return err
	}
	return nil
}

// CheckBlobExists checks if a blob exists without fetching the entire metadata
func (s *PostgresBlobMetadataStore) CheckBlobExists(ctx context.Context, blobKey corev2.BlobKey) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM blob_metadata WHERE blob_key = $1)"
	err := s.db.QueryRow(ctx, query, blobKey[:]).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check blob existence: %w", err)
	}
	return exists, nil
}

// GetBlobMetadata retrieves blob metadata by its key
func (s *PostgresBlobMetadataStore) GetBlobMetadata(ctx context.Context, blobKey corev2.BlobKey) (*v2.BlobMetadata, error) {
	query := "SELECT blob_header, signature, requested_at, blob_status, updated_at, expiry FROM blob_metadata WHERE blob_key = $1"

	var blobHeader, signature []byte
	var requestedAt, updatedAt, expiry uint64
	var blobStatus int

	err := s.db.QueryRow(ctx, query, blobKey[:]).Scan(
		&blobHeader,
		&signature,
		&requestedAt,
		&blobStatus,
		&updatedAt,
		&expiry,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: metadata not found for key %s", common.ErrMetadataNotFound, blobKey.Hex())
		}
		return nil, fmt.Errorf("failed to get blob metadata: %w", err)
	}

	metadata := &v2.BlobMetadata{
		Signature:   signature,
		RequestedAt: requestedAt,
		BlobStatus:  v2.BlobStatus(blobStatus),
		UpdatedAt:   updatedAt,
		Expiry:      expiry,
	}

	err = json.Unmarshal(blobHeader, &metadata.BlobHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal blob header: %w", err)
	}

	return metadata, nil
}

// PutBlobMetadata stores blob metadata
func (s *PostgresBlobMetadataStore) PutBlobMetadata(ctx context.Context, blobMetadata *v2.BlobMetadata) error {
	s.logger.Debug("store put blob metadata", "blobMetadata", blobMetadata)

	// Marshal BlobHeader to JSON
	blobHeaderJSON, err := json.Marshal(blobMetadata.BlobHeader)
	if err != nil {
		return fmt.Errorf("failed to marshal blob header: %w", err)
	}

	// Get blob key
	blobKey, err := blobMetadata.BlobHeader.BlobKey()
	if err != nil {
		return fmt.Errorf("failed to get blob key: %w", err)
	}

	// Generate additional fields
	requestedAtBucket := computeRequestedAtBucket(blobMetadata.RequestedAt)
	requestedAtBlobKey := encodeBlobFeedCursorKey(blobMetadata.RequestedAt, &blobKey)

	// Marshal FragmentInfo to JSON if present
	var fragmentInfoJSON []byte
	if blobMetadata.FragmentInfo != nil {
		fragmentInfoJSON, err = json.Marshal(blobMetadata.FragmentInfo)
		if err != nil {
			return fmt.Errorf("failed to marshal fragment info: %w", err)
		}
	}

	// Insert into database
	query := `
		INSERT INTO blob_metadata (
			blob_key, blob_header, signature, requested_at, requested_at_bucket, 
			requested_at_blob_key, blob_status, updated_at, account_id, expiry,
			num_retries, blob_size, fragment_info
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		) ON CONFLICT DO NOTHING
	`
	commandTag, err := s.db.Exec(
		ctx, query,
		blobKey[:], blobHeaderJSON, blobMetadata.Signature, blobMetadata.RequestedAt, requestedAtBucket,
		requestedAtBlobKey, int(blobMetadata.BlobStatus), blobMetadata.UpdatedAt,
		blobMetadata.BlobHeader.PaymentMetadata.AccountID.Hex(), blobMetadata.Expiry,
		blobMetadata.NumRetries, blobMetadata.BlobSize, fragmentInfoJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to insert blob metadata: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return dispersercommon.ErrAlreadyExists
	}

	return nil
}

// UpdateBlobStatus updates the status of a blob
func (s *PostgresBlobMetadataStore) UpdateBlobStatus(ctx context.Context, blobKey corev2.BlobKey, status v2.BlobStatus) error {
	validStatuses := statusUpdatePrecondition[status]
	if len(validStatuses) == 0 {
		return fmt.Errorf("%w: invalid status transition to %s", ErrInvalidStateTransition, status.String())
	}

	// Build the WHERE condition for valid statuses
	var placeholders []string
	var args []interface{}
	args = append(args, blobKey[:]) // $1 for the blob key

	for i, validStatus := range validStatuses {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+2))
		args = append(args, int(validStatus))
	}
	statusPlaceholders := strings.Join(placeholders, ", ")

	// Current time for updated_at
	now := time.Now().UnixNano()

	// Update the blob status
	query := fmt.Sprintf(`
		UPDATE blob_metadata 
		SET blob_status = $%d, updated_at = $%d 
		WHERE blob_key = $1 AND blob_status IN (%s)
	`, len(args)+1, len(args)+2, statusPlaceholders)

	args = append(args, int(status), now)

	commandTag, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update blob status: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		// Get current blob status to provide better error messages
		blob, err := s.GetBlobMetadata(ctx, blobKey)
		if err != nil {
			return fmt.Errorf("failed to get blob metadata for key %s: %v", blobKey.Hex(), err)
		}

		if blob.BlobStatus == status {
			return fmt.Errorf("%w: blob already in status %s", dispersercommon.ErrAlreadyExists, status.String())
		}

		return fmt.Errorf("%w: invalid status transition from %s to %s", ErrInvalidStateTransition, blob.BlobStatus.String(), status.String())
	}

	return nil
}

// DeleteBlobMetadata deletes blob metadata by its key (only used in testing)
func (s *PostgresBlobMetadataStore) DeleteBlobMetadata(ctx context.Context, blobKey corev2.BlobKey) error {
	query := "DELETE FROM blob_metadata WHERE blob_key = $1"
	_, err := s.db.Exec(ctx, query, blobKey[:])
	if err != nil {
		return fmt.Errorf("failed to delete blob metadata: %w", err)
	}
	return nil
}

// GetBlobMetadataByAccountID retrieves blob metadata by account ID within a time range
func (s *PostgresBlobMetadataStore) GetBlobMetadataByAccountID(
	ctx context.Context,
	accountId gethcommon.Address,
	start uint64,
	end uint64,
	limit int,
	ascending bool,
) ([]*v2.BlobMetadata, error) {
	if start+1 > end-1 {
		return nil, fmt.Errorf("no time point in exclusive time range (%d, %d)", start, end)
	}

	// Adjust time range to be exclusive
	adjustedStart, adjustedEnd := start+1, end-1

	// Order by clause based on ascending flag
	orderBy := "ASC"
	if !ascending {
		orderBy = "DESC"
	}

	// Limit clause
	limitClause := ""
	if limit > 0 {
		limitClause = fmt.Sprintf("LIMIT %d", limit)
	}

	query := fmt.Sprintf(`
		SELECT blob_key, blob_header, requested_at, blob_status, updated_at
		FROM blob_metadata
		WHERE account_id = $1 AND requested_at BETWEEN $2 AND $3
		ORDER BY requested_at %s
		%s
	`, orderBy, limitClause)

	rows, err := s.db.Query(ctx, query, accountId.Hex(), adjustedStart, adjustedEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to query blob metadata by account ID: %w", err)
	}
	defer rows.Close()

	var results []*v2.BlobMetadata
	for rows.Next() {
		var blobKey, blobHeader []byte
		var requestedAt, updatedAt uint64
		var blobStatus int

		if err := rows.Scan(&blobKey, &blobHeader, &requestedAt, &blobStatus, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan blob metadata: %w", err)
		}

		metadata := &v2.BlobMetadata{
			RequestedAt: requestedAt,
			BlobStatus:  v2.BlobStatus(blobStatus),
			UpdatedAt:   updatedAt,
		}

		// Unmarshal BlobHeader
		err = json.Unmarshal(blobHeader, &metadata.BlobHeader)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal blob header: %w", err)
		}

		results = append(results, metadata)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return results, nil
}

// GetBlobMetadataByStatus retrieves blob metadata by status updated after a specific timestamp
func (s *PostgresBlobMetadataStore) GetBlobMetadataByStatus(ctx context.Context, status v2.BlobStatus, lastUpdatedAt uint64) ([]*v2.BlobMetadata, error) {
	query := `
		SELECT blob_key, blob_header, signature, requested_at, blob_status, updated_at, expiry
		FROM blob_metadata
		WHERE blob_status = $1 AND updated_at > $2
		ORDER BY updated_at ASC
	`

	rows, err := s.db.Query(ctx, query, int(status), lastUpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to query blob metadata by status: %w", err)
	}
	defer rows.Close()

	var results []*v2.BlobMetadata
	for rows.Next() {
		var blobKey, blobHeader, signature []byte
		var requestedAt, updatedAt, expiry uint64
		var blobStatus int

		if err := rows.Scan(&blobKey, &blobHeader, &signature, &requestedAt, &blobStatus, &updatedAt, &expiry); err != nil {
			return nil, fmt.Errorf("failed to scan blob metadata: %w", err)
		}

		metadata := &v2.BlobMetadata{
			Signature:   signature,
			RequestedAt: requestedAt,
			BlobStatus:  v2.BlobStatus(blobStatus),
			UpdatedAt:   updatedAt,
			Expiry:      expiry,
		}

		// Unmarshal BlobHeader
		err = json.Unmarshal(blobHeader, &metadata.BlobHeader)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal blob header: %w", err)
		}

		results = append(results, metadata)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return results, nil
}

// GetBlobMetadataByStatusPaginated retrieves blob metadata by status with pagination
func (s *PostgresBlobMetadataStore) GetBlobMetadataByStatusPaginated(
	ctx context.Context,
	status v2.BlobStatus,
	exclusiveStartKey *StatusIndexCursor,
	limit int32,
) ([]*v2.BlobMetadata, *StatusIndexCursor, error) {
	var query string
	var args []interface{}

	if exclusiveStartKey != nil && exclusiveStartKey.BlobKey != nil {
		// Continue from the previous cursor
		query = `
			SELECT blob_key, blob_header, signature, requested_at, blob_status, updated_at, expiry, blob_size, num_retries, fragment_info
			FROM blob_metadata
			WHERE blob_status = $1 AND 
				  (updated_at > $2 OR (updated_at = $2 AND blob_key > $3))
			ORDER BY updated_at ASC, blob_key ASC
			LIMIT $4
		`
		args = append(args, int(status), exclusiveStartKey.UpdatedAt, exclusiveStartKey.BlobKey[:], limit)
	} else {
		// Start from the beginning
		query = `
			SELECT blob_key, blob_header, signature, requested_at, blob_status, updated_at, expiry, blob_size, num_retries, fragment_info
			FROM blob_metadata
			WHERE blob_status = $1
			ORDER BY updated_at ASC, blob_key ASC
			LIMIT $2
		`
		args = append(args, int(status), limit)
	}

	s.logger.Info("querying blob metadata by status paginated", "query", query, "args", args)
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query blob metadata by status paginated: %w", err)
	}
	defer rows.Close()

	var results []*v2.BlobMetadata
	var lastBlobKey []byte
	var lastUpdatedAt uint64

	for rows.Next() {
		var blobKey, blobHeader, signature []byte
		var requestedAt, updatedAt, expiry, blobSize uint64
		var numRetries uint
		var blobStatus int
		var fragmentInfoJSON []byte

		if err := rows.Scan(&blobKey, &blobHeader, &signature, &requestedAt, &blobStatus, &updatedAt, &expiry, &blobSize, &numRetries, &fragmentInfoJSON); err != nil {
			return nil, nil, fmt.Errorf("failed to scan blob metadata: %w", err)
		}

		lastBlobKey = blobKey
		lastUpdatedAt = updatedAt

		metadata := &v2.BlobMetadata{
			Signature:   signature,
			RequestedAt: requestedAt,
			BlobStatus:  v2.BlobStatus(blobStatus),
			UpdatedAt:   updatedAt,
			Expiry:      expiry,
			BlobSize:    blobSize,
			NumRetries:  numRetries,
		}

		// Unmarshal BlobHeader
		err = json.Unmarshal(blobHeader, &metadata.BlobHeader)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal blob header: %w", err)
		}

		// Unmarshal FragmentInfo if present
		if fragmentInfoJSON != nil {
			var fragmentInfo encoding.FragmentInfo
			if err := json.Unmarshal(fragmentInfoJSON, &fragmentInfo); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal fragment info: %w", err)
			}
			metadata.FragmentInfo = &fragmentInfo
		}

		results = append(results, metadata)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	// If no results found, return the same cursor
	if len(results) == 0 {
		return results, exclusiveStartKey, nil
	}

	// Check if we've reached the limit (more records might be available)
	if len(results) < int(limit) {
		// No more records, return nil cursor to indicate end
		for _, result := range results {
			s.logger.Info("results is less than limit and no more records",
				"results", fmt.Sprintf("%+v", result),
				"num_results", len(results),
				"limit", limit)
		}
		return results, nil, nil
	}

	// Create next cursor
	var bk corev2.BlobKey
	copy(bk[:], lastBlobKey)
	nextCursor := &StatusIndexCursor{
		BlobKey:   &bk,
		UpdatedAt: lastUpdatedAt,
	}

	for _, result := range results {
		s.logger.Info("results is less than limit and there are more records",
			"results", fmt.Sprintf("%+v", result),
			"num_results", len(results),
			"limit", limit)
	}

	return results, nextCursor, nil
}

// GetBlobMetadataCountByStatus counts the number of blobs with a given status
func (s *PostgresBlobMetadataStore) GetBlobMetadataCountByStatus(ctx context.Context, status v2.BlobStatus) (int32, error) {
	var count int32
	query := "SELECT COUNT(*) FROM blob_metadata WHERE blob_status = $1"
	err := s.db.QueryRow(ctx, query, int(status)).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count blob metadata by status: %w", err)
	}
	return count, nil
}

// queryBucketBlobMetadata retrieves blob metadata from a specific bucket
func (s *PostgresBlobMetadataStore) queryBucketBlobMetadata(
	ctx context.Context,
	bucket uint64,
	ascending bool,
	after BlobFeedCursor,
	before BlobFeedCursor,
	startKey string,
	endKey string,
	limit int,
	result []*v2.BlobMetadata,
	lastProcessedCursor **BlobFeedCursor,
) ([]*v2.BlobMetadata, error) {
	// Order by clause based on ascending flag
	orderBy := "ASC"
	if !ascending {
		orderBy = "DESC"
	}

	// Limit clause
	limitClause := ""
	if limit > 0 {
		remainingLimit := limit - len(result)
		if remainingLimit <= 0 {
			return result, nil
		}
		limitClause = fmt.Sprintf("LIMIT %d", remainingLimit)
	}

	query := fmt.Sprintf(`
		SELECT blob_key, blob_header, requested_at, blob_status, updated_at
		FROM blob_metadata
		WHERE requested_at_bucket = $1 
		  AND requested_at_blob_key  > $2
		  AND requested_at_blob_key  < $3
		ORDER BY requested_at_blob_key %s
		%s
	`, orderBy, limitClause)

	rows, err := s.db.Query(ctx, query, fmt.Sprintf("%d", bucket), startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("failed to query blob metadata from bucket: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var blobKey, blobHeader []byte
		var requestedAt, updatedAt uint64
		var blobStatus int

		if err := rows.Scan(&blobKey, &blobHeader, &requestedAt, &blobStatus, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan blob metadata: %w", err)
		}

		metadata := &v2.BlobMetadata{
			RequestedAt: requestedAt,
			BlobStatus:  v2.BlobStatus(blobStatus),
			UpdatedAt:   updatedAt,
		}

		// Unmarshal BlobHeader
		err = json.Unmarshal(blobHeader, &metadata.BlobHeader)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal blob header: %w", err)
		}

		// Get blob key for filtering
		blobKeyValue, err := metadata.BlobHeader.BlobKey()
		if err != nil {
			return nil, fmt.Errorf("failed to get blob key: %w", err)
		}

		// Skip blobs at the endpoints (exclusive bounds)
		if after.Equal(metadata.RequestedAt, &blobKeyValue) || before.Equal(metadata.RequestedAt, &blobKeyValue) {
			continue
		}

		// Add to result
		result = append(result, metadata)

		// Update last processed cursor
		*lastProcessedCursor = &BlobFeedCursor{
			RequestedAt: metadata.RequestedAt,
			BlobKey:     &blobKeyValue,
		}

		// Check limit
		if limit > 0 && len(result) >= limit {
			return result, nil
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return result, nil
}

// GetBlobMetadataByRequestedAtForward retrieves blob metadata ordered by requested time in ascending order
func (s *PostgresBlobMetadataStore) GetBlobMetadataByRequestedAtForward(
	ctx context.Context,
	after BlobFeedCursor,
	before BlobFeedCursor,
	limit int,
) ([]*v2.BlobMetadata, *BlobFeedCursor, error) {
	if !after.LessThan(&before) {
		return nil, nil, errors.New("after cursor must be less than before cursor")
	}

	startBucket, endBucket := GetRequestedAtBucketIDRange(after.RequestedAt, before.RequestedAt)
	startKey := after.ToCursorKey()
	endKey := before.ToCursorKey()

	result := make([]*v2.BlobMetadata, 0)
	var lastProcessedCursor *BlobFeedCursor

	for bucket := startBucket; bucket <= endBucket; bucket++ {
		var err error
		result, err = s.queryBucketBlobMetadata(
			ctx, bucket, true, after, before, startKey, endKey, limit, result, &lastProcessedCursor,
		)
		if err != nil {
			return nil, nil, err
		}

		if limit > 0 && len(result) >= limit {
			break
		}
	}

	return result, lastProcessedCursor, nil
}

// GetBlobMetadataByRequestedAtBackward retrieves blob metadata ordered by requested time in descending order
func (s *PostgresBlobMetadataStore) GetBlobMetadataByRequestedAtBackward(
	ctx context.Context,
	before BlobFeedCursor,
	after BlobFeedCursor,
	limit int,
) ([]*v2.BlobMetadata, *BlobFeedCursor, error) {
	if !after.LessThan(&before) {
		return nil, nil, errors.New("after cursor must be less than before cursor")
	}

	startBucket, endBucket := GetRequestedAtBucketIDRange(after.RequestedAt, before.RequestedAt)
	startKey := after.ToCursorKey()
	endKey := before.ToCursorKey()
	result := make([]*v2.BlobMetadata, 0)
	var lastProcessedCursor *BlobFeedCursor

	// Traverse buckets in reverse order
	for bucket := endBucket; bucket >= startBucket; bucket-- {
		var err error
		result, err = s.queryBucketBlobMetadata(
			ctx, bucket, false, after, before, startKey, endKey, limit, result, &lastProcessedCursor,
		)
		if err != nil {
			return nil, nil, err
		}

		if limit > 0 && len(result) >= limit {
			break
		}
	}
	return result, lastProcessedCursor, nil
}

// PutBlobCertificate stores a blob certificate
func (s *PostgresBlobMetadataStore) PutBlobCertificate(ctx context.Context, blobCert *corev2.BlobCertificate, fragmentInfo *encoding.FragmentInfo) error {
	// Marshal to JSON
	blobCertJSON, err := json.Marshal(blobCert)
	if err != nil {
		return fmt.Errorf("failed to marshal blob certificate: %w", err)
	}

	fragmentInfoJSON, err := json.Marshal(fragmentInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal fragment info: %w", err)
	}

	// Get blob key
	blobKey, err := blobCert.BlobHeader.BlobKey()
	if err != nil {
		return fmt.Errorf("failed to get blob key: %w", err)
	}

	// Insert into database
	query := `
		INSERT INTO blob_certificates (blob_key, blob_certificate, fragment_info)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`
	commandTag, err := s.db.Exec(ctx, query, blobKey[:], blobCertJSON, fragmentInfoJSON)
	if err != nil {
		return fmt.Errorf("failed to insert blob certificate: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return dispersercommon.ErrAlreadyExists
	}

	return nil
}

// DeleteBlobCertificate deletes a blob certificate
func (s *PostgresBlobMetadataStore) DeleteBlobCertificate(ctx context.Context, blobKey corev2.BlobKey) error {
	query := "DELETE FROM blob_certificates WHERE blob_key = $1"
	_, err := s.db.Exec(ctx, query, blobKey[:])
	if err != nil {
		return fmt.Errorf("failed to delete blob certificate: %w", err)
	}
	return nil
}

// GetBlobCertificate retrieves a blob certificate
func (s *PostgresBlobMetadataStore) GetBlobCertificate(ctx context.Context, blobKey corev2.BlobKey) (*corev2.BlobCertificate, *encoding.FragmentInfo, error) {
	query := "SELECT blob_certificate, fragment_info FROM blob_certificates WHERE blob_key = $1"

	var blobCertificateBytes, fragmentInfoBytes []byte
	err := s.db.QueryRow(ctx, query, blobKey[:]).Scan(&blobCertificateBytes, &fragmentInfoBytes)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, fmt.Errorf("%w: certificate not found for key %s", common.ErrMetadataNotFound, blobKey.Hex())
		}
		return nil, nil, fmt.Errorf("failed to get blob certificate: %w", err)
	}

	var cert corev2.BlobCertificate
	if err := json.Unmarshal(blobCertificateBytes, &cert); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal blob certificate: %w", err)
	}

	var fragmentInfo encoding.FragmentInfo
	if err := json.Unmarshal(fragmentInfoBytes, &fragmentInfo); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal fragment info: %w", err)
	}

	return &cert, &fragmentInfo, nil
}

func (s *PostgresBlobMetadataStore) GetBlobCertificates(
	ctx context.Context,
	blobKeys []corev2.BlobKey,
) ([]*corev2.BlobCertificate, []*encoding.FragmentInfo, error) {
	// 1) nothing to do
	if len(blobKeys) == 0 {
		return nil, nil, nil
	}

	// 2) build a [][]byte of raw keys
	keys := make([][]byte, len(blobKeys))
	for i, k := range blobKeys {
		keys[i] = append([]byte(nil), k[:]...)
	}

	// 3) fetch rows
	query := `
		SELECT blob_certificate, fragment_info
		FROM blob_certificates
		WHERE blob_key = ANY($1)
	`
	rows, err := s.db.Query(ctx, query, keys)
	if err != nil {
		return nil, nil, fmt.Errorf("query blob_certificates: %w", err)
	}
	defer rows.Close()

	// 4) unmarshal each JSONB column into your structs
	var certs []*corev2.BlobCertificate
	var fragments []*encoding.FragmentInfo
	for rows.Next() {
		var certJSON, fragJSON []byte
		if err := rows.Scan(&certJSON, &fragJSON); err != nil {
			return nil, nil, fmt.Errorf("scan row: %w", err)
		}

		var cert corev2.BlobCertificate
		if err := json.Unmarshal(certJSON, &cert); err != nil {
			return nil, nil, fmt.Errorf("unmarshal BlobCertificate: %w", err)
		}

		var fi encoding.FragmentInfo
		if err := json.Unmarshal(fragJSON, &fi); err != nil {
			return nil, nil, fmt.Errorf("unmarshal FragmentInfo: %w", err)
		}

		certs = append(certs, &cert)
		fragments = append(fragments, &fi)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("row iteration error: %w", err)
	}

	return certs, fragments, nil
}

// PutBatch stores a batch
func (s *PostgresBlobMetadataStore) PutBatch(ctx context.Context, batch *corev2.Batch) error {
	// Marshal to JSON
	batchJSON, err := json.Marshal(batch)
	if err != nil {
		return fmt.Errorf("failed to marshal batch: %w", err)
	}

	// Get batch header hash
	batchHeaderHash, err := batch.BatchHeader.Hash()
	if err != nil {
		return fmt.Errorf("failed to hash batch header: %w", err)
	}

	// Insert into database
	query := `
		INSERT INTO batches (batch_header_hash, batch_info)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	commandTag, err := s.db.Exec(ctx, query, batchHeaderHash[:], batchJSON)
	if err != nil {
		return fmt.Errorf("failed to insert batch: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return dispersercommon.ErrAlreadyExists
	}

	return nil
}

// GetBatch retrieves a batch by its header hash
func (s *PostgresBlobMetadataStore) GetBatch(ctx context.Context, batchHeaderHash [32]byte) (*corev2.Batch, error) {
	query := "SELECT batch_info FROM batches WHERE batch_header_hash = $1"

	var batchJSON []byte
	err := s.db.QueryRow(ctx, query, batchHeaderHash[:]).Scan(&batchJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: batch info not found for hash %x", common.ErrMetadataNotFound, batchHeaderHash)
		}
		return nil, fmt.Errorf("failed to get batch: %w", err)
	}

	var batch corev2.Batch
	if err := json.Unmarshal(batchJSON, &batch); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch: %w", err)
	}

	return &batch, nil
}

// PutBatchHeader stores a batch header
func (s *PostgresBlobMetadataStore) PutBatchHeader(ctx context.Context, batchHeader *corev2.BatchHeader) error {
	// Marshal to JSON
	batchHeaderJSON, err := json.Marshal(batchHeader)
	if err != nil {
		return fmt.Errorf("failed to marshal batch header: %w", err)
	}

	// Get batch header hash
	batchHeaderHash, err := batchHeader.Hash()
	if err != nil {
		return fmt.Errorf("failed to hash batch header: %w", err)
	}

	// Insert into database
	query := `
		INSERT INTO batch_headers (batch_header_hash, batch_header)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	commandTag, err := s.db.Exec(ctx, query, batchHeaderHash[:], batchHeaderJSON)
	if err != nil {
		return fmt.Errorf("failed to insert batch header: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return dispersercommon.ErrAlreadyExists
	}

	return nil
}

// DeleteBatchHeader deletes a batch header
func (s *PostgresBlobMetadataStore) DeleteBatchHeader(ctx context.Context, batchHeaderHash [32]byte) error {
	query := "DELETE FROM batch_headers WHERE batch_header_hash = $1"
	_, err := s.db.Exec(ctx, query, batchHeaderHash[:])
	if err != nil {
		return fmt.Errorf("failed to delete batch header: %w", err)
	}
	return nil
}

// GetBatchHeader retrieves a batch header by its hash
func (s *PostgresBlobMetadataStore) GetBatchHeader(ctx context.Context, batchHeaderHash [32]byte) (*corev2.BatchHeader, error) {
	query := "SELECT batch_header FROM batch_headers WHERE batch_header_hash = $1"

	var batchHeaderJSON []byte
	err := s.db.QueryRow(ctx, query, batchHeaderHash[:]).Scan(&batchHeaderJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: batch header not found for hash %x", common.ErrMetadataNotFound, batchHeaderHash)
		}
		return nil, fmt.Errorf("failed to get batch header: %w", err)
	}

	var batchHeader corev2.BatchHeader
	if err := json.Unmarshal(batchHeaderJSON, &batchHeader); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch header: %w", err)
	}

	return &batchHeader, nil
}

// PutDispersalRequest stores a dispersal request
func (s *PostgresBlobMetadataStore) PutDispersalRequest(ctx context.Context, req *corev2.DispersalRequest) error {
	// Marshal to JSON
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal dispersal request: %w", err)
	}

	// Get batch header hash
	batchHeaderHash, err := req.BatchHeader.Hash()
	if err != nil {
		return fmt.Errorf("failed to hash batch header: %w", err)
	}

	// Insert into database
	query := `
		INSERT INTO dispersal_requests (batch_header_hash, operator_id, dispersal_request, dispersed_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT DO NOTHING
	`
	commandTag, err := s.db.Exec(ctx, query, batchHeaderHash[:], req.OperatorID[:], reqJSON, req.DispersedAt)
	if err != nil {
		return fmt.Errorf("failed to insert dispersal request: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return dispersercommon.ErrAlreadyExists
	}

	return nil
}

// GetDispersalRequest retrieves a dispersal request
func (s *PostgresBlobMetadataStore) GetDispersalRequest(ctx context.Context, batchHeaderHash [32]byte, operatorID core.OperatorID) (*corev2.DispersalRequest, error) {
	query := "SELECT dispersal_request FROM dispersal_requests WHERE batch_header_hash = $1 AND operator_id = $2"

	var reqJSON []byte
	err := s.db.QueryRow(ctx, query, batchHeaderHash[:], operatorID[:]).Scan(&reqJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: dispersal request not found for batch header hash %x and operator %s", common.ErrMetadataNotFound, batchHeaderHash, operatorID.Hex())
		}
		return nil, fmt.Errorf("failed to get dispersal request: %w", err)
	}

	var req corev2.DispersalRequest
	if err := json.Unmarshal(reqJSON, &req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dispersal request: %w", err)
	}

	return &req, nil
}

// PutDispersalResponse stores a dispersal response
func (s *PostgresBlobMetadataStore) PutDispersalResponse(ctx context.Context, res *corev2.DispersalResponse) error {
	// Marshal to JSON
	resJSON, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal dispersal response: %w", err)
	}

	// Get batch header hash
	batchHeaderHash, err := res.BatchHeader.Hash()
	if err != nil {
		return fmt.Errorf("failed to hash batch header: %w", err)
	}

	// Insert into database
	query := `
		INSERT INTO dispersal_responses (batch_header_hash, operator_id, dispersal_response, responded_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT DO NOTHING
	`
	commandTag, err := s.db.Exec(ctx, query, batchHeaderHash[:], res.OperatorID[:], resJSON, res.RespondedAt)
	if err != nil {
		return fmt.Errorf("failed to insert dispersal response: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return dispersercommon.ErrAlreadyExists
	}

	return nil
}

// GetDispersalResponse retrieves a dispersal response
func (s *PostgresBlobMetadataStore) GetDispersalResponse(ctx context.Context, batchHeaderHash [32]byte, operatorID core.OperatorID) (*corev2.DispersalResponse, error) {
	query := "SELECT dispersal_response FROM dispersal_responses WHERE batch_header_hash = $1 AND operator_id = $2"

	var resJSON []byte
	err := s.db.QueryRow(ctx, query, batchHeaderHash[:], operatorID[:]).Scan(&resJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: dispersal response not found for batch header hash %x and operator %s", common.ErrMetadataNotFound, batchHeaderHash, operatorID.Hex())
		}
		return nil, fmt.Errorf("failed to get dispersal response: %w", err)
	}

	var res corev2.DispersalResponse
	if err := json.Unmarshal(resJSON, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dispersal response: %w", err)
	}

	return &res, nil
}

// GetDispersalResponses retrieves all dispersal responses for a batch
func (s *PostgresBlobMetadataStore) GetDispersalResponses(ctx context.Context, batchHeaderHash [32]byte) ([]*corev2.DispersalResponse, error) {
	query := "SELECT dispersal_response FROM dispersal_responses WHERE batch_header_hash = $1"
	rows, err := s.db.Query(ctx, query, batchHeaderHash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to query dispersal responses: %w", err)
	}
	defer rows.Close()

	var responses []*corev2.DispersalResponse
	for rows.Next() {
		var resJSON []byte
		if err := rows.Scan(&resJSON); err != nil {
			return nil, fmt.Errorf("failed to scan dispersal response: %w", err)
		}

		var res corev2.DispersalResponse
		if err := json.Unmarshal(resJSON, &res); err != nil {
			return nil, fmt.Errorf("failed to unmarshal dispersal response: %w", err)
		}

		responses = append(responses, &res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	if len(responses) == 0 {
		return nil, fmt.Errorf("%w: dispersal responses not found for batch header hash %x", common.ErrMetadataNotFound, batchHeaderHash)
	}

	return responses, nil
}

// GetDispersalsByRespondedAt retrieves dispersal responses within a time range
func (s *PostgresBlobMetadataStore) GetDispersalsByRespondedAt(
	ctx context.Context,
	operatorId core.OperatorID,
	start uint64,
	end uint64,
	limit int,
	ascending bool,
) ([]*corev2.DispersalResponse, error) {
	if start+1 > end-1 {
		return nil, fmt.Errorf("no time point in exclusive time range (%d, %d)", start, end)
	}

	// Adjust time range to be exclusive
	adjustedStart, adjustedEnd := start+1, end-1

	// Order by clause based on ascending flag
	orderBy := "ASC"
	if !ascending {
		orderBy = "DESC"
	}

	// Limit clause
	limitClause := ""
	if limit > 0 {
		limitClause = fmt.Sprintf("LIMIT %d", limit)
	}

	query := fmt.Sprintf(`
		SELECT dispersal_response
		FROM dispersal_responses
		WHERE operator_id = $1 AND responded_at BETWEEN $2 AND $3
		ORDER BY responded_at %s
		%s
	`, orderBy, limitClause)

	rows, err := s.db.Query(ctx, query, operatorId[:], adjustedStart, adjustedEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to query dispersal responses by responded at: %w", err)
	}
	defer rows.Close()

	var responses []*corev2.DispersalResponse
	for rows.Next() {
		var resJSON []byte
		if err := rows.Scan(&resJSON); err != nil {
			return nil, fmt.Errorf("failed to scan dispersal response: %w", err)
		}

		var res corev2.DispersalResponse
		if err := json.Unmarshal(resJSON, &res); err != nil {
			return nil, fmt.Errorf("failed to unmarshal dispersal response: %w", err)
		}

		responses = append(responses, &res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return responses, nil
}

// PutAttestation stores an attestation
func (s *PostgresBlobMetadataStore) PutAttestation(ctx context.Context, attestation *corev2.Attestation) error {
	// Marshal to JSON
	attestationJSON, err := json.Marshal(attestation)
	if err != nil {
		return fmt.Errorf("failed to marshal attestation: %w", err)
	}

	// Get batch header hash
	batchHeaderHash, err := attestation.BatchHeader.Hash()
	if err != nil {
		return fmt.Errorf("failed to hash batch header: %w", err)
	}

	// Insert into database - allow overwrite
	query := `
		INSERT INTO attestations (batch_header_hash, attestation, attested_at, attested_at_bucket)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (batch_header_hash) DO UPDATE
		SET attestation = $2, attested_at = $3, attested_at_bucket = $4
	`
	attestedAtBucket := computeAttestedAtBucket(attestation.AttestedAt)
	_, err = s.db.Exec(ctx, query, batchHeaderHash[:], attestationJSON, attestation.AttestedAt, attestedAtBucket)
	if err != nil {
		return fmt.Errorf("failed to insert attestation: %w", err)
	}

	return nil
}

// GetAttestation retrieves an attestation
func (s *PostgresBlobMetadataStore) GetAttestation(ctx context.Context, batchHeaderHash [32]byte) (*corev2.Attestation, error) {
	query := "SELECT attestation FROM attestations WHERE batch_header_hash = $1"

	var attestationJSON []byte
	err := s.db.QueryRow(ctx, query, batchHeaderHash[:]).Scan(&attestationJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: attestation not found for hash %x", common.ErrMetadataNotFound, batchHeaderHash)
		}
		return nil, fmt.Errorf("failed to get attestation: %w", err)
	}

	var attestation corev2.Attestation
	if err := json.Unmarshal(attestationJSON, &attestation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal attestation: %w", err)
	}

	return &attestation, nil
}

// queryBucketAttestation retrieves attestations from a specific bucket
func (s *PostgresBlobMetadataStore) queryBucketAttestation(
	ctx context.Context,
	bucket, start, end uint64,
	numToReturn int,
	ascending bool,
) ([]*corev2.Attestation, error) {
	// Order by clause based on ascending flag
	orderBy := "ASC"
	if !ascending {
		orderBy = "DESC"
	}

	// Limit clause
	limitClause := ""
	if numToReturn > 0 {
		limitClause = fmt.Sprintf("LIMIT %d", numToReturn)
	}

	query := fmt.Sprintf(`
		SELECT attestation
		FROM attestations
		WHERE attested_at_bucket = $1 AND attested_at BETWEEN $2 AND $3
		ORDER BY attested_at %s
		%s
	`, orderBy, limitClause)

	rows, err := s.db.Query(ctx, query, fmt.Sprintf("%d", bucket), start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query bucket attestations: %w", err)
	}
	defer rows.Close()

	var attestations []*corev2.Attestation
	for rows.Next() {
		var attestationJSON []byte
		if err := rows.Scan(&attestationJSON); err != nil {
			return nil, fmt.Errorf("failed to scan attestation: %w", err)
		}

		var attestation corev2.Attestation
		if err := json.Unmarshal(attestationJSON, &attestation); err != nil {
			return nil, fmt.Errorf("failed to unmarshal attestation: %w", err)
		}

		attestations = append(attestations, &attestation)

		// Check limit
		if numToReturn > 0 && len(attestations) >= numToReturn {
			break
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return attestations, nil
}

// GetAttestationByAttestedAtForward retrieves attestations ordered by attested time in ascending order
func (s *PostgresBlobMetadataStore) GetAttestationByAttestedAtForward(
	ctx context.Context,
	after uint64,
	before uint64,
	limit int,
) ([]*corev2.Attestation, error) {
	if after+1 > before-1 {
		return nil, fmt.Errorf("no time point in exclusive time range (%d, %d)", after, before)
	}
	startBucket, endBucket := GetAttestedAtBucketIDRange(after, before)
	result := make([]*corev2.Attestation, 0)

	// Traverse buckets in forward order
	for bucket := startBucket; bucket <= endBucket; bucket++ {
		if limit > 0 && len(result) >= limit {
			break
		}
		remaining := math.MaxInt
		if limit > 0 {
			remaining = limit - len(result)
		}
		// Query bucket in ascending order
		bucketAttestation, err := s.queryBucketAttestation(ctx, bucket, after+1, before-1, remaining, true)
		if err != nil {
			return nil, err
		}
		for _, ba := range bucketAttestation {
			result = append(result, ba)
			if limit > 0 && len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

// GetAttestationByAttestedAtBackward retrieves attestations ordered by attested time in descending order
func (s *PostgresBlobMetadataStore) GetAttestationByAttestedAtBackward(
	ctx context.Context,
	before uint64,
	after uint64,
	limit int,
) ([]*corev2.Attestation, error) {
	if after+1 > before-1 {
		return nil, fmt.Errorf("no time point in exclusive time range (%d, %d)", after, before)
	}
	// Note: we traverse buckets in reverse order for backward query
	startBucket, endBucket := GetAttestedAtBucketIDRange(after, before)
	result := make([]*corev2.Attestation, 0)

	// Traverse buckets in reverse order
	for bucket := endBucket; bucket >= startBucket; bucket-- {
		if limit > 0 && len(result) >= limit {
			break
		}
		remaining := math.MaxInt
		if limit > 0 {
			remaining = limit - len(result)
		}
		// Query bucket in descending order
		bucketAttestation, err := s.queryBucketAttestation(ctx, bucket, after+1, before-1, remaining, false)
		if err != nil {
			return nil, err
		}
		for _, ba := range bucketAttestation {
			result = append(result, ba)
			if limit > 0 && len(result) >= limit {
				break
			}
		}
	}

	return result, nil
}

// PutBlobInclusionInfo stores blob inclusion information
func (s *PostgresBlobMetadataStore) PutBlobInclusionInfo(ctx context.Context, inclusionInfo *corev2.BlobInclusionInfo) error {
	// Marshal to JSON
	inclusionInfoJSON, err := json.Marshal(inclusionInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal blob inclusion info: %w", err)
	}

	// Get batch header hash
	batchHeaderHash, err := inclusionInfo.BatchHeader.Hash()
	if err != nil {
		return fmt.Errorf("failed to hash batch header: %w", err)
	}

	// Insert into database
	query := `
		INSERT INTO blob_inclusion_info (blob_key, batch_header_hash, inclusion_info)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`
	commandTag, err := s.db.Exec(ctx, query, inclusionInfo.BlobKey[:], batchHeaderHash[:], inclusionInfoJSON)
	if err != nil {
		return fmt.Errorf("failed to insert blob inclusion info: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return dispersercommon.ErrAlreadyExists
	}

	return nil
}

// PutBlobInclusionInfos stores multiple blob inclusion information entries
func (s *PostgresBlobMetadataStore) PutBlobInclusionInfos(ctx context.Context, inclusionInfos []*corev2.BlobInclusionInfo) error {
	if len(inclusionInfos) == 0 {
		return nil
	}

	// Start a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// For pgx, we won't use prepared statements directly as pgx handles this efficiently
	for _, info := range inclusionInfos {
		// Marshal to JSON
		inclusionInfoJSON, err := json.Marshal(info)
		if err != nil {
			return fmt.Errorf("failed to marshal blob inclusion info: %w", err)
		}

		// Get batch header hash
		batchHeaderHash, err := info.BatchHeader.Hash()
		if err != nil {
			return fmt.Errorf("failed to hash batch header: %w", err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO blob_inclusion_info (blob_key, batch_header_hash, inclusion_info)
			VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING
		`, info.BlobKey[:], batchHeaderHash[:], inclusionInfoJSON)

		if err != nil {
			return fmt.Errorf("failed to insert blob inclusion info: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetBlobInclusionInfo retrieves blob inclusion information
func (s *PostgresBlobMetadataStore) GetBlobInclusionInfo(ctx context.Context, blobKey corev2.BlobKey, batchHeaderHash [32]byte) (*corev2.BlobInclusionInfo, error) {
	query := "SELECT inclusion_info FROM blob_inclusion_info WHERE blob_key = $1 AND batch_header_hash = $2"

	var inclusionInfoJSON []byte
	err := s.db.QueryRow(ctx, query, blobKey[:], batchHeaderHash[:]).Scan(&inclusionInfoJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: inclusion info not found for key %s", common.ErrMetadataNotFound, blobKey.Hex())
		}
		return nil, fmt.Errorf("failed to get blob inclusion info: %w", err)
	}

	var inclusionInfo corev2.BlobInclusionInfo
	if err := json.Unmarshal(inclusionInfoJSON, &inclusionInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal blob inclusion info: %w", err)
	}

	return &inclusionInfo, nil
}

// GetBlobInclusionInfos retrieves all inclusion information for a blob
func (s *PostgresBlobMetadataStore) GetBlobInclusionInfos(ctx context.Context, blobKey corev2.BlobKey) ([]*corev2.BlobInclusionInfo, error) {
	query := "SELECT inclusion_info FROM blob_inclusion_info WHERE blob_key = $1"
	rows, err := s.db.Query(ctx, query, blobKey[:])
	if err != nil {
		return nil, fmt.Errorf("failed to query blob inclusion infos: %w", err)
	}
	defer rows.Close()

	var inclusionInfos []*corev2.BlobInclusionInfo
	for rows.Next() {
		var inclusionInfoJSON []byte
		if err := rows.Scan(&inclusionInfoJSON); err != nil {
			return nil, fmt.Errorf("failed to scan blob inclusion info: %w", err)
		}

		var inclusionInfo corev2.BlobInclusionInfo
		if err := json.Unmarshal(inclusionInfoJSON, &inclusionInfo); err != nil {
			return nil, fmt.Errorf("failed to unmarshal blob inclusion info: %w", err)
		}

		inclusionInfos = append(inclusionInfos, &inclusionInfo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	if len(inclusionInfos) == 0 {
		return nil, fmt.Errorf("%w: inclusion info not found for key %s", common.ErrMetadataNotFound, blobKey.Hex())
	}

	return inclusionInfos, nil
}

// GetBlobAttestationInfo retrieves blob attestation information
func (s *PostgresBlobMetadataStore) GetBlobAttestationInfo(ctx context.Context, blobKey corev2.BlobKey) (*v2.BlobAttestationInfo, error) {
	blobInclusionInfos, err := s.GetBlobInclusionInfos(ctx, blobKey)
	if err != nil {
		s.logger.Error("failed to get blob inclusion info for blob", "err", err, "blobKey", blobKey.Hex())
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to get blob inclusion info: %s", err.Error()))
	}

	if len(blobInclusionInfos) == 0 {
		s.logger.Error("no blob inclusion info found for blob", "blobKey", blobKey.Hex())
		return nil, api.NewErrorInternal("no blob inclusion info found")
	}

	if len(blobInclusionInfos) > 1 {
		s.logger.Warn("multiple inclusion info found for blob", "blobKey", blobKey.Hex())
	}

	for _, inclusionInfo := range blobInclusionInfos {
		// get the signed batch from this inclusion info
		batchHeaderHash, err := inclusionInfo.BatchHeader.Hash()
		if err != nil {
			s.logger.Error("failed to get batch header hash from blob inclusion info", "err", err, "blobKey", blobKey.Hex())
			continue
		}
		_, attestation, err := s.GetSignedBatch(ctx, batchHeaderHash)
		if err != nil {
			s.logger.Error("failed to get signed batch", "err", err, "blobKey", blobKey.Hex())
			continue
		}

		return &v2.BlobAttestationInfo{
			InclusionInfo: inclusionInfo,
			Attestation:   attestation,
		}, nil
	}

	return nil, fmt.Errorf("no attestation info found for blobkey: %s", blobKey.Hex())
}

// GetSignedBatch retrieves a batch header and its attestation
func (s *PostgresBlobMetadataStore) GetSignedBatch(ctx context.Context, batchHeaderHash [32]byte) (*corev2.BatchHeader, *corev2.Attestation, error) {
	// Get batch header
	batchHeader, err := s.GetBatchHeader(ctx, batchHeaderHash)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get batch header: %w", err)
	}

	// Get attestation
	attestation, err := s.GetAttestation(ctx, batchHeaderHash)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get attestation: %w", err)
	}

	return batchHeader, attestation, nil
}

// Close closes the database connection
func (s *PostgresBlobMetadataStore) Close() error {
	s.db.Close()
	return nil
}
