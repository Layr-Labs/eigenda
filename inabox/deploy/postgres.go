package deploy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// StartDockertestWithPostgresContainer starts a Postgres container for testing
func StartDockertestWithPostgresContainer(postgresPort string) (*dockertest.Pool, *dockertest.Resource, error) {
	fmt.Println("Starting Postgres container")
	pool, err := dockertest.NewPool("")
	if err != nil {
		fmt.Println("Could not construct pool: %w", err)
		return nil, nil, err
	}

	err = pool.Client.Ping()
	if err != nil {
		fmt.Println("Could not connect to Docker: %w", err)
		return nil, nil, err
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "postgres",
		Tag:          "14",
		Name:         "postgres-eigenda-test",
		ExposedPorts: []string{postgresPort},
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port("5432/tcp"): {
				{HostIP: "0.0.0.0", HostPort: postgresPort},
			},
		},
		Env: []string{
			fmt.Sprintf("POSTGRES_PORT=%s", postgresPort),
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_DB=eigenda",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		fmt.Println("Could not start resource: %w", err)
		return nil, nil, err
	}

	pool.MaxWait = 60 * time.Second
	if err := pool.Retry(func() error {
		fmt.Println("Waiting for Postgres to start")
		// Using the exec function to check if PostgreSQL is up
		exitCode, err := resource.Exec([]string{"pg_isready", "-U", "postgres"}, dockertest.ExecOptions{})
		if err != nil || exitCode != 0 {
			fmt.Println("PostgreSQL is not ready yet...")
			return fmt.Errorf("postgres is not ready yet")
		}

		fmt.Println("PostgreSQL is running and responding!")
		return nil

	}); err != nil {
		fmt.Println("Could not connect to PostgreSQL:", err)
		return nil, nil, err
	}

	log.Printf("PostgreSQL started successfully! Available at localhost:%s", postgresPort)

	// Initialize the database
	if err := InitializePostgresDatabase(postgresPort); err != nil {
		fmt.Println("Could not initialize PostgreSQL database:", err)
		return nil, nil, err
	}

	return pool, resource, nil
}

// InitializePostgresDatabase initializes the database with schemas required for eigenda
func InitializePostgresDatabase(postgresPort string) error {
	fmt.Println("Initializing PostgreSQL database schemas")
	db, err := pgxpool.New(context.Background(), fmt.Sprintf("host=localhost port=%s user=postgres password=postgres dbname=eigenda sslmode=disable", postgresPort))
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	if err := initTables(db); err != nil {
		return fmt.Errorf("failed to initialize tables: %w", err)
	}

	fmt.Println("PostgreSQL database schemas initialized successfully!")
	return nil
}

// initTables creates the necessary tables if they don't exist
func initTables(db *pgxpool.Pool) error {
	ctx := context.Background()
	// Create tables for blob metadata
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS blob_metadata (
			blob_key BYTEA PRIMARY KEY,
			blob_header JSONB NOT NULL,
			blob_size INTEGER NOT NULL,
			num_retries INTEGER NOT NULL,
			signature BYTEA NOT NULL,
			requested_at BIGINT NOT NULL,
			requested_at_bucket BYTEA NOT NULL,
			requested_at_blob_key BYTEA NOT NULL,
			blob_status INTEGER NOT NULL,
			updated_at BIGINT NOT NULL,
			account_id VARCHAR(42) NOT NULL,
			expiry BIGINT NOT NULL,
			fragment_info JSONB
		);

		CREATE INDEX IF NOT EXISTS idx_blob_metadata_blob_status_updated_at ON blob_metadata (blob_status, updated_at);
		CREATE INDEX IF NOT EXISTS idx_blob_metadata_account_id_requested_at ON blob_metadata (account_id, requested_at);
		CREATE INDEX IF NOT EXISTS idx_blob_metadata_requested_at_bucket_key ON blob_metadata (requested_at_bucket, requested_at_blob_key);
	`)
	if err != nil {
		return fmt.Errorf("failed to create blob_metadata table: %w", err)
	}

	// Create tables for blob certificates
	_, err = db.Exec(ctx, `
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
	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS batch_headers (
			batch_header_hash BYTEA PRIMARY KEY,
			batch_header JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create batch_headers table: %w", err)
	}

	// Create tables for batches
	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS batches (
			batch_header_hash BYTEA PRIMARY KEY,
			batch_info JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create batches table: %w", err)
	}

	// Create tables for dispersal requests
	_, err = db.Exec(ctx, `
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
	_, err = db.Exec(ctx, `
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
	_, err = db.Exec(ctx, `
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
	_, err = db.Exec(ctx, `
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
