package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/go-redis/redis/v8"
)

// Config ... user configurable
type Config struct {
	Endpoint string
	Password string
	DB       int
	Eviction time.Duration
	Profile  bool
}

// Store ... Redis storage backend implementation (This not safe for concurrent usage)
type Store struct {
	eviction time.Duration

	client *redis.Client

	profile bool
	reads   int
	entries int
}

var _ store.PrecomputedKeyStore = (*Store)(nil)

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
		profile:  cfg.Profile,
		reads:    0,
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

	if r.profile {
		r.reads++
	}

	// cast value to byte slice
	return []byte(value), nil
}

// Put ... inserts a value into the Redis store
func (r *Store) Put(ctx context.Context, key []byte, value []byte) error {
	err := r.client.Set(ctx, string(key), string(value), r.eviction).Err()
	if err == nil && r.profile {
		r.entries++
	}

	return err
}

func (r *Store) Verify(_ []byte, _ []byte) error {
	return nil
}

func (r *Store) BackendType() store.BackendType {
	return store.RedisBackendType
}

func (r *Store) Stats() *store.Stats {
	return &store.Stats{
		Entries: r.entries,
		Reads:   r.reads,
	}
}
