package littbuilder

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/metrics"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// All table names must match this regex. Permits uppercase and lowercase letters, numbers, underscores, and dashes.
// Must be at least one character long.
var tableNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

var _ litt.DB = &db{}

// TableBuilderFunc is a function that creates a new table.
type TableBuilderFunc func(
	ctx context.Context,
	logger logging.Logger,
	name string,
	metrics *metrics.LittDBMetrics) (litt.ManagedTable, error)

// db is an implementation of DB.
type db struct {
	ctx    context.Context
	logger logging.Logger

	// A function that returns the current time.
	clock func() time.Time

	// The default time-to-live for new tables. Once created, the TTL for a table can be changed.
	ttl time.Duration

	// The period between garbage collection runs.
	gcPeriod time.Duration

	// A function that creates new tables.
	tableBuilder TableBuilderFunc

	// A map of all tables in the database.
	tables map[string]litt.ManagedTable

	// Protects access to tables and ttl.
	lock sync.Mutex

	// True if the database has been stopped.
	stopped atomic.Bool

	// Metrics for the database.
	metrics *metrics.LittDBMetrics

	// The HTTP server for metrics. nil if metrics are disabled or if an external party is managing the server.
	metricsServer *http.Server
}

// NewDB creates a new DB instance. After this method is called, the config object should not be modified.
func NewDB(config *litt.Config) (litt.DB, error) {
	var err error

	if config.Logger == nil {
		config.Logger, err = buildLogger(config)
		if err != nil {
			return nil, fmt.Errorf("error building logger: %w", err)
		}
	}

	err = config.SanityCheck()
	if err != nil {
		return nil, fmt.Errorf("error checking config: %w", err)
	}

	err = config.SanitizePaths()
	if err != nil {
		return nil, fmt.Errorf("error expanding tildes in config: %w", err)
	}

	if !config.Fsync {
		config.Logger.Warnf(
			"Fsync is disabled. Ok for unit tests that need to run fast, NOT OK FOR PRODUCTION USE.")
	}

	tableBuilder := func(
		ctx context.Context,
		logger logging.Logger,
		name string,
		metrics *metrics.LittDBMetrics) (litt.ManagedTable, error) {

		return buildTable(config, logger, name, metrics)
	}

	return NewDBUnsafe(config, tableBuilder)
}

// NewDBUnsafe creates a new DB instance with a custom table builder. This is intended for unit test use,
// and should not be considered a stable API.
func NewDBUnsafe(config *litt.Config, tableBuilder TableBuilderFunc) (litt.DB, error) {
	var err error

	if config.Logger == nil {
		config.Logger, err = buildLogger(config)
		if err != nil {
			return nil, fmt.Errorf("error building logger: %w", err)
		}
	}

	var dbMetrics *metrics.LittDBMetrics
	var metricsServer *http.Server
	if config.MetricsEnabled {
		dbMetrics, metricsServer = buildMetrics(config, config.Logger)
	}

	database := &db{
		ctx:           config.CTX,
		logger:        config.Logger,
		clock:         config.Clock,
		ttl:           config.TTL,
		gcPeriod:      config.GCPeriod,
		tableBuilder:  tableBuilder,
		tables:        make(map[string]litt.ManagedTable),
		metrics:       dbMetrics,
		metricsServer: metricsServer,
	}

	if config.MetricsEnabled {
		go database.gatherMetrics(config.MetricsUpdateInterval)
	}

	return database, nil
}

func (d *db) KeyCount() uint64 {
	d.lock.Lock()
	defer d.lock.Unlock()

	count := uint64(0)
	for _, table := range d.tables {
		count += table.KeyCount()
	}

	return count
}

func (d *db) Size() uint64 {
	d.lock.Lock()
	defer d.lock.Unlock()

	return d.lockFreeSize()
}

func (d *db) lockFreeSize() uint64 {
	size := uint64(0)
	for _, table := range d.tables {
		size += table.Size()
	}

	return size
}

// isTableNameValid returns true if the table name is valid.
func (d *db) isTableNameValid(name string) bool {
	return tableNameRegex.MatchString(name)
}

func (d *db) GetTable(name string) (litt.Table, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	table, ok := d.tables[name]
	if !ok {
		if !d.isTableNameValid(name) {
			return nil, fmt.Errorf(
				"Table name '%s' is invalid, must be at least one character long and "+
					"contain only letters, numbers, and underscores, and dashes.", name)
		}

		var err error
		table, err = d.tableBuilder(d.ctx, d.logger, name, d.metrics)
		if err != nil {
			return nil, fmt.Errorf("error creating table: %w", err)
		}
		d.logger.Infof(
			"Table '%s' initialized, table contains %d key-value pairs and has a size of %s.",
			name, table.KeyCount(), util.PrettyPrintBytes(table.Size()))

		d.tables[name] = table
	}

	return table, nil
}

func (d *db) DropTable(name string) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	table, ok := d.tables[name]
	if !ok {
		// Table does not exist, nothing to do.
		d.logger.Infof("table %s does not exist, cannot drop", name)
		return nil
	}

	d.logger.Infof("dropping table %s", name)
	err := table.Destroy()
	if err != nil {
		return fmt.Errorf("error destroying table: %w", err)
	}
	delete(d.tables, name)

	return nil
}

func (d *db) Close() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.logger.Infof("Stopping LittDB, estimated data size: %d", d.lockFreeSize())
	d.stopped.Store(true)

	for name, table := range d.tables {
		err := table.Close()
		if err != nil {
			return fmt.Errorf("error stopping table %s: %w", name, err)
		}
	}

	return nil
}

func (d *db) Destroy() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.stopped.Store(true)

	for name, table := range d.tables {
		err := table.Destroy()
		if err != nil {
			return fmt.Errorf("error destroying table %s: %w", name, err)
		}
	}

	return nil
}

// gatherMetrics is a method that periodically collects metrics.
func (d *db) gatherMetrics(interval time.Duration) {
	if d.metricsServer != nil {
		defer func() {
			err := d.metricsServer.Close()
			if err != nil {
				d.logger.Errorf("error closing metrics server: %v", err)
			}
		}()
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for !d.stopped.Load() {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			d.lock.Lock()
			tablesCopy := make(map[string]litt.ManagedTable, len(d.tables))
			for name, table := range d.tables {
				tablesCopy[name] = table
			}
			d.lock.Unlock()

			d.metrics.CollectPeriodicMetrics(tablesCopy)
		}
	}
}
