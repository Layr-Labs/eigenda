package store

import (
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/common"
)

type Config struct {
	BackendsToEnable []common.EigenDABackend
	DispersalBackend common.EigenDABackend

	AsyncPutWorkers int
	FallbackTargets []string
	CacheTargets    []string
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
