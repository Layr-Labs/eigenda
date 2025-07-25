package store

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
)

type Config struct {
	BackendsToEnable []common.EigenDABackend
	DispersalBackend common.EigenDABackend

	AsyncPutWorkers int
	CacheTargets    []string

	WriteOnCacheMiss bool
}

// Verifies that a backend target slice is constructed correctly
func (cfg *Config) checkCacheTargets(targets []string) error {
	if len(targets) == 0 {
		return nil
	}

	if common.ContainsDuplicates(targets) {
		return fmt.Errorf("duplicate cache targets provided: %+v", targets)
	}

	for _, t := range targets {
		if common.StringToBackendType(t) == common.UnknownBackendType {
			return fmt.Errorf("unknown cache target provided: %s", t)
		}
	}

	return nil
}

// Check ... verifies that configuration values are adequately set
func (cfg *Config) Check() error {

	err := cfg.checkCacheTargets(cfg.CacheTargets)
	if err != nil {
		return err
	}

	// verify that thread counts are sufficiently set
	if cfg.AsyncPutWorkers > 100 {
		return fmt.Errorf("number of secondary write workers can't be greater than 100")
	}

	return nil
}
