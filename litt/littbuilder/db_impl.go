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
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var tableNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

var _ litt.DB = &db{}

// tableBuilder is a function that creates a new table.
type TableBuilder func(
	ctx context.Context,
	logger logging.Logger,
	timeSource func() time.Time,
	name string,
	ttl time.Duration,
	metrics *metrics.LittDBMetrics) (litt.ManagedTable, error)

// db is an implementation of DB.
type db struct {
	ctx    context.Context
	logger logging.Logger

	// A function that returns the current time.
	timeSource func() time.Time

	// The time-to-live for tables that haven't had their TTL set.
	ttl time.Duration

	// The period between garbage collection runs.
	gcPeriod time.Duration

	// A function that creates new tables.
	tableBuilder TableBuilder

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
func NewDB(config *LittDBConfig) (litt.DB, error) {
	return NewDBWithTableBuilder(config, config.buildTable)
}

// NewDBWithTableBuilder creates a new DB instance with a custom table builder. This is intended for unit test use,
// and should not be considered a stable API.
func NewDBWithTableBuilder(config *LittDBConfig, tableBuilder TableBuilder) (litt.DB, error) {
	logger, err := config.buildLogger()
	if err != nil {
		return nil, fmt.Errorf("error building logger: %w", err)
	}

	var dbMetrics *metrics.LittDBMetrics
	var metricsServer *http.Server
	if config.MetricsEnabled {
		dbMetrics, metricsServer = config.buildMetrics(logger)
	}

	database := &db{
		ctx:           config.CTX,
		logger:        logger,
		timeSource:    config.TimeSource,
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

	logger.Infof("LittDB started, current data size: %d", database.Size())

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
	if name == "" {
		return false
	}
	return tableNameRegex.MatchString(name)
}

func (d *db) GetTable(name string) (litt.Table, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	table, ok := d.tables[name]
	if !ok {
		if !d.isTableNameValid(name) {
			return nil, fmt.Errorf(
				"table name %s is invalid, must be at least one character long and "+
					"contain only letters, numbers, and underscores, and dashes.", name)
		}

		var err error
		table, err = d.tableBuilder(d.ctx, d.logger, d.timeSource, name, d.ttl, d.metrics)
		if err != nil {
			return nil, fmt.Errorf("error creating table: %w", err)
		}

		d.tables[name] = table
	}

	return table, nil
}

func (d *db) DropTable(name string) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	table, ok := d.tables[name]
	if !ok {
		return fmt.Errorf("table %s does not exist", name)
	}

	err := table.Destroy()
	if err != nil {
		return fmt.Errorf("error destroying table: %w", err)
	}
	delete(d.tables, name)

	return nil
}

func (d *db) Stop() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.logger.Infof("Stopping LittDB, estimated data size: %d", d.lockFreeSize())
	d.stopped.Store(true)

	for name, table := range d.tables {
		err := table.Stop()
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

	for !d.stopped.Load() {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			d.metrics.CollectPeriodicMetrics(d, d.tables)
		}
	}
}
