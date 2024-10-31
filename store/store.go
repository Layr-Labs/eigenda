package store

import (
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/redis"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/s3"
)

type Config struct {
	AsyncPutWorkers int
	FallbackTargets []string
	CacheTargets    []string

	// secondary storage cfgs
	RedisConfig redis.Config
	S3Config    s3.Config
}

// checkTargets ... verifies that a backend target slice is constructed correctly
func (cfg *Config) checkTargets(targets []string) error {
	if len(targets) == 0 {
		return nil
	}

	if common.ContainsDuplicates(targets) {
		return fmt.Errorf("duplicate targets provided: %+v", targets)
	}

	for _, t := range targets {
		if common.StringToBackendType(t) == common.UnknownBackendType {
			return fmt.Errorf("unknown fallback target provided: %s", t)
		}
	}

	return nil
}

// Check ... verifies that configuration values are adequately set
func (cfg *Config) Check() error {
	if cfg.S3Config.CredentialType == s3.CredentialTypeUnknown && cfg.S3Config.Endpoint != "" {
		return fmt.Errorf("s3 credential type must be set")
	}
	if cfg.S3Config.CredentialType == s3.CredentialTypeStatic {
		if cfg.S3Config.Endpoint != "" && (cfg.S3Config.AccessKeyID == "" || cfg.S3Config.AccessKeySecret == "") {
			return fmt.Errorf("s3 endpoint is set, but access key id or access key secret is not set")
		}
	}

	if cfg.RedisConfig.Endpoint == "" && cfg.RedisConfig.Password != "" {
		return fmt.Errorf("redis password is set, but endpoint is not")
	}

	err := cfg.checkTargets(cfg.FallbackTargets)
	if err != nil {
		return err
	}

	err = cfg.checkTargets(cfg.CacheTargets)
	if err != nil {
		return err
	}

	// verify that same target is not in both fallback and cache targets
	for _, t := range cfg.FallbackTargets {
		if common.Contains(cfg.CacheTargets, t) {
			return fmt.Errorf("target %s is in both fallback and cache targets", t)
		}
	}

	// verify that thread counts are sufficiently set
	if cfg.AsyncPutWorkers >= 100 {
		return fmt.Errorf("number of secondary write workers can't be greater than 100")
	}

	return nil
}
