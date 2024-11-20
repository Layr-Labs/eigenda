package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/go-redis/redis/v8"
)

// Config ... user configurable
type Config struct {
	Endpoint string
	Password string
	DB       int
	Eviction time.Duration
}

// Custom MarshalJSON function to control what gets included in the JSON output.
// TODO: Probably best would be to separate config from secrets everywhere.
// Then we could just log the config and not worry about secrets.
func (c Config) MarshalJSON() ([]byte, error) {
	type Alias Config // Use an alias to avoid recursion with MarshalJSON
	aux := (Alias)(c)
	// Conditionally include a masked password if it is set
	if aux.Password != "" {
		aux.Password = "*****"
	}
	return json.Marshal(aux)
}

// Store ... Redis storage backend implementation
// go-redis client is safe for concurrent usage: https://github.com/redis/go-redis/blob/v8.11.5/redis.go#L535-L544
type Store struct {
	eviction time.Duration

	client *redis.Client
}

var _ common.PrecomputedKeyStore = (*Store)(nil)

// NewStore ... constructor
func NewStore(cfg *Config) (*Store, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Endpoint,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// ensure server can be pinged using potential client connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := client.Ping(ctx)
	if cmd.Err() != nil {
		return nil, fmt.Errorf("failed to ping redis server: %w", cmd.Err())
	}

	return &Store{
		eviction: cfg.Eviction,
		client:   client,
	}, nil
}

// Get ... retrieves a value from the Redis store. Returns nil if the key is not found vs. an error
// if the key is found but the value is not retrievable.
func (r *Store) Get(ctx context.Context, key []byte) ([]byte, error) {
	value, err := r.client.Get(ctx, string(key)).Result()
	if errors.Is(err, redis.Nil) { // key DNE
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// cast value to byte slice
	return []byte(value), nil
}

// Put ... inserts a value into the Redis store
func (r *Store) Put(ctx context.Context, key []byte, value []byte) error {
	return r.client.Set(ctx, string(key), string(value), r.eviction).Err()
}

func (r *Store) Verify(_ context.Context, _, _ []byte) error {
	return nil
}

func (r *Store) BackendType() common.BackendType {
	return common.RedisBackendType
}
