package littbuilder

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var tableNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

var _ litt.DB = &db{}

// tableBuilder is a function that creates a new table.
type tableBuilder func(
	ctx context.Context,
	logger logging.Logger,
	timeSource func() time.Time,
	name string,
	ttl time.Duration) (litt.ManagedTable, error)

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
	tableBuilder tableBuilder

	// A map of all tables in the database.
	tables map[string]litt.ManagedTable

	// Protects access to tables and ttl.
	lock sync.Mutex
}

// NewDB creates a new DB instance. In general, this should not be used directly. Instead, use LittDBConfig.Build()
func NewDB(
	ctx context.Context,
	logger logging.Logger,
	timeSource func() time.Time,
	ttl time.Duration,
	gcPeriod time.Duration,
	tableBuilder tableBuilder) litt.DB {

	return &db{
		ctx:          ctx,
		logger:       logger,
		timeSource:   timeSource,
		ttl:          ttl,
		gcPeriod:     gcPeriod,
		tableBuilder: tableBuilder,
		tables:       make(map[string]litt.ManagedTable),
	}
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
		table, err = d.tableBuilder(d.ctx, d.logger, d.timeSource, name, d.ttl)
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

	for name, table := range d.tables {
		err := table.Destroy()
		if err != nil {
			return fmt.Errorf("error destroying table %s: %w", name, err)
		}
	}

	return nil
}
