package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConfig ... user configurable
type RedisConfig struct {
	Endpoint string
	Password string
	DB       int
	Eviction time.Duration
	Profile  bool
}

// RedStore ... Redis storage backend implementation (This not safe for concurrent usage)
type RedStore struct {
	eviction time.Duration

	client *redis.Client

	profile bool
	reads   int
	entries int
}

var _ PrecomputedKeyStore = (*RedStore)(nil)

// NewRedisStore ... constructor
func NewRedisStore(cfg *RedisConfig) (*RedStore, error) {
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

	return &RedStore{
		eviction: cfg.Eviction,
		client:   client,
		profile:  cfg.Profile,
		reads:    0,
	}, nil
}

// Get ... retrieves a value from the Redis store. Returns nil if the key is not found vs. an error
// if the key is found but the value is not retrievable.
func (r *RedStore) Get(ctx context.Context, key []byte) ([]byte, error) {
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
func (r *RedStore) Put(ctx context.Context, key []byte, value []byte) error {
	err := r.client.Set(ctx, string(key), string(value), r.eviction).Err()
	if err == nil && r.profile {
		r.entries++
	}

	return err
}

func (r *RedStore) Verify(_ []byte, _ []byte) error {
	return nil
}

func (r *RedStore) BackendType() BackendType {
	return Redis
}

func (r *RedStore) Stats() *Stats {
	return &Stats{
		Entries: r.entries,
		Reads:   r.reads,
	}
}
