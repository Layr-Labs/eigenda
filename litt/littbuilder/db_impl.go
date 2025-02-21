package littbuilder

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/litt"
	"sync"
	"sync/atomic"
	"time"
)

var _ litt.DB = &db{}

// tableBuilder is a function that creates a new table.
type tableBuilder func(timeSource func() time.Time, name string, ttl time.Duration) (litt.ManagedTable, error)

// db is an implementation of DB.
type db struct {
	// A function that returns the current time.
	timeSource func() time.Time

	// The time-to-live for newly created tables.
	ttl time.Duration

	// The period between garbage collection runs.
	gcPeriod time.Duration

	// A function that creates new tables.
	tableBuilder tableBuilder

	// A map of all tables in the database.
	tables map[string]litt.ManagedTable

	// A flag that indicates whether the database is alive (i.e. Stop() has not been called).
	alive atomic.Bool

	// Protects access to tables and ttl.
	lock sync.Mutex
}

func newDB(
	timeSource func() time.Time,
	ttl time.Duration,
	gcPeriod time.Duration,
	tableBuilder tableBuilder) litt.DB {

	return &db{
		timeSource:   timeSource,
		ttl:          ttl,
		gcPeriod:     gcPeriod,
		tableBuilder: tableBuilder,
	}
}

func (d *db) GetTable(name string) (litt.Table, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	table, ok := d.tables[name]
	if !ok {
		var err error
		table, err = d.tableBuilder(d.timeSource, name, d.ttl)
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

func (d *db) Start() {
	d.alive.Store(true)

	ticker := time.NewTicker(d.gcPeriod)
	go func() {
		for d.alive.Load() {
			<-ticker.C
			d.doGarbageCollection()
		}
	}()
}

func (d *db) Stop() {
	d.alive.Store(false)
}

// doGarbageCollection performs garbage collection on all tables in the database.
func (d *db) doGarbageCollection() {
	tables := make([]litt.ManagedTable, 0, len(d.tables))
	d.lock.Lock()
	for _, table := range d.tables {
		tables = append(tables, table)
	}
	d.lock.Unlock()

	for _, table := range tables {
		err := table.DoGarbageCollection()
		if err != nil {
			// TODO log!
			panic(err)
		}
	}
}
