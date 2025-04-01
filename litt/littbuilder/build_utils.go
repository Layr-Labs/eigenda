package littbuilder

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/cache"
	"github.com/Layr-Labs/eigenda/litt"
	tablecache "github.com/Layr-Labs/eigenda/litt/cache"
	"github.com/Layr-Labs/eigenda/litt/disktable"
	"github.com/Layr-Labs/eigenda/litt/disktable/keymap"
	"github.com/Layr-Labs/eigenda/litt/metrics"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// cacheWeight is a function that calculates the weight of a cache entry.
func cacheWeight(key string, value []byte) uint64 {
	return uint64(len(key) + len(value))
}

// buildKeymap creates a new keymap based on the configuration.
func buildKeymap(
	config *LittDBConfig,
	logger logging.Logger,
	tableName string,
) (kmap keymap.Keymap, keymapPath string, keymapTypeFile *keymap.KeymapTypeFile, requiresReload bool, err error) {

	builderForConfiguredType, ok := keymapBuilders[config.KeymapType]
	if !ok {
		return nil, "", nil, false,
			fmt.Errorf("unsupported keymap type: %v", config.KeymapType)
	}

	keymapDirectories := make([]string, len(config.Paths))
	for i, p := range config.Paths {
		keymapDirectories[i] = path.Join(p, tableName, keymap.KeymapDirectoryName)
	}

	var keymapDirectory string
	for _, directory := range keymapDirectories {
		exists, err := keymap.KeymapFileExists(directory)
		if err != nil {
			return nil, "", nil, false,
				fmt.Errorf("error checking for keymap type file: %w", err)
		}
		if exists {
			keymapDirectory = directory
			keymapTypeFile, err = keymap.LoadKeymapTypeFile(directory)
			if err != nil {
				return nil, "", nil, false,
					fmt.Errorf("error loading keymap type file: %w", err)
			}
			break
		}
	}

	newKeymap := false
	if keymapTypeFile == nil {
		// No previous keymap exists. Either we are starting fresh or the keymap was deleted manually.
		newKeymap = true

		// by convention, always select the first path as the keymap directory
		keymapDirectory = keymapDirectories[0]
		keymapTypeFile = keymap.NewKeymapTypeFile(keymapDirectory, config.KeymapType)

		// create the keymap directory
		err := os.MkdirAll(keymapDirectory, 0755)
		if err != nil {
			return nil, "", nil, false,
				fmt.Errorf("error creating keymap directory: %w", err)
		}

		// write the keymap type file
		err = keymapTypeFile.Write()
		if err != nil {
			return nil, "", nil, false,
				fmt.Errorf("error writing keymap type file: %w", err)
		}

	} else {
		// A previous keymap exists. Check if the keymap type has changed.

		builderForTypeOnDisk, ok := keymapBuilders[keymapTypeFile.Type()]
		if !ok {
			return nil, "", nil, false,
				fmt.Errorf("unsupported keymap type: %v", keymapTypeFile.Type())
		}

		if config.KeymapType != keymapTypeFile.Type() {
			// The previously used keymap type is different from the one in the configuration.

			// delete the old keymap
			err := builderForTypeOnDisk.DeleteFiles(logger, keymapDirectory)
			if err != nil {
				return nil, "", nil, false,
					fmt.Errorf("error deleting keymap files: %w", err)
			}

			// delete the keymap type file
			err = keymapTypeFile.Delete()
			if err != nil {
				return nil, "", nil, false,
					fmt.Errorf("error deleting keymap type file: %w", err)
			}

			// finally, delete the keymap directory
			_, err = os.Stat(keymapDirectory)
			if err == nil {
				err = os.Remove(keymapDirectory)
				if err != nil {
					return nil, "", nil, false,
						fmt.Errorf("error deleting keymap directory: %w", err)
				}
			} else if !os.IsNotExist(err) {
				return nil, "", nil, false,
					fmt.Errorf("error checking for keymap directory: %w", err)
			}
		}
	}

	keymapDataDirectory := path.Join(keymapDirectory, keymap.KeymapDataDirectoryName)
	kmap, requiresReload, err = builderForConfiguredType.Build(logger, keymapDataDirectory, config.DoubleWriteProtection)
	if err != nil {
		return nil, "", nil, false,
			fmt.Errorf("error building keymap: %w", err)
	}

	return kmap, keymapDirectory, keymapTypeFile, requiresReload || newKeymap, nil
}

// buildTable creates a new table based on the configuration.
func buildTable(
	config *LittDBConfig,
	ctx context.Context,
	logger logging.Logger,
	timeSource func() time.Time,
	name string,
	ttl time.Duration,
	metrics *metrics.LittDBMetrics) (litt.ManagedTable, error) {

	var table litt.ManagedTable

	if config.ShardingFactor < 1 {
		return nil, fmt.Errorf("sharding factor must be at least 1")
	}

	kmap, keymapDirectory, keymapTypeFile, requiresReload, err := buildKeymap(config, logger, name)
	if err != nil {
		return nil, fmt.Errorf("error creating keymap: %w", err)
	}

	tableRoots := make([]string, len(config.Paths))
	for i, p := range config.Paths {
		tableRoots[i] = path.Join(p, name)
	}

	table, err = disktable.NewDiskTable(
		ctx,
		logger,
		timeSource,
		name,
		kmap,
		keymapDirectory,
		keymapTypeFile,
		tableRoots,
		config.TargetSegmentFileSize,
		config.ControlChannelSize,
		config.ShardingFactor,
		config.SaltShaker,
		ttl,
		config.GCPeriod,
		requiresReload,
		config.Fsync,
		metrics)

	if err != nil {
		return nil, fmt.Errorf("error creating table: %w", err)
	}

	tableCache := cache.NewFIFOCache[string, []byte](config.CacheSize, cacheWeight)
	tableCache = cache.NewThreadSafeCache(tableCache)
	cachedTable := tablecache.NewCachedTable(table, tableCache)

	return cachedTable, nil
}

// buildLogger creates a new logger based on the configuration.
func buildLogger(config *LittDBConfig) (logging.Logger, error) {
	if config.Logger != nil {
		return config.Logger, nil
	}

	return common.NewLogger(config.LoggerConfig)
}

// buildMetrics creates a new metrics object based on the configuration. If the returned server is not nil,
// then it is the responsibility of the caller to eventually call server.Shutdown().
func buildMetrics(config *LittDBConfig, logger logging.Logger) (*metrics.LittDBMetrics, *http.Server) {
	if !config.MetricsEnabled {
		return nil, nil
	}

	var registry *prometheus.Registry
	var server *http.Server

	if config.MetricsRegistry != nil {
		registry = prometheus.NewRegistry()
		registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		registry.MustRegister(collectors.NewGoCollector())

		logger.Infof("Starting metrics server at port %d", config.MetricsPort)
		addr := fmt.Sprintf(":%d", config.MetricsPort)
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.HandlerFor(
			registry,
			promhttp.HandlerOpts{},
		))
		server = &http.Server{
			Addr:    addr,
			Handler: mux,
		}

		go func() {
			err := server.ListenAndServe()
			if err != nil && !strings.Contains(err.Error(), "http: Server closed") {
				logger.Errorf("metrics server error: %v", err)
			}
		}()
	}

	return metrics.NewLittDBMetrics(registry, config.MetricsNamespace), server
}
