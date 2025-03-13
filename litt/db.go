package litt

// DB is a highly specialized key-value store. It is intentionally very feature poor, sacrificing
// unnecessary features for simplicity, high performance, and low memory usage.
//
// Litt: slang, a synonym for "cool" or "awesome". e.g. "Man, that database is litt, bro!".
//
// Supported features:
// - writing values
// - reading values
// - TTLs and automatic (lazy) deletion of expired values
// - tables with non-overlapping namespaces
// - thread safety: all methods are safe to call concurrently, and all modifications are atomic
//
// Unsupported features:
// - mutating existing values (once a value is written, it cannot be changed)
// - deleting values (values only leave the DB when they expire via a TTL)
// - transactions (individual operations are atomic, but there is no way to group operations atomically)
// - fine granularity for TTL (all data in the same table must have the same TTL)
type DB interface {
	// GetTable gets a table by name, creating one if it does not exist.
	//
	// Table names appear as directories on the file system, and so table names are restricted to be
	// ASCII alphanumeric characters, dashes, and underscores. The name must be at least one character long.
	//
	// The first time a table is fetched (either a new table or an existing one loaded from disk), its TTL is always
	// set to 0 (i.e. it has no TTL). If you want to set a TTL, you must call Table.SetTTL() to do so. This is
	// necessary after each time the database is started/restarted.
	GetTable(name string) (Table, error)

	// DropTable deletes a table and all of its data.
	//
	// Note that it is NOT thread safe to drop a table concurrently with any operation that accesses the table.
	// The table returned by GetTable() before DropTable() is called must not be used once DropTable() is called.
	DropTable(name string) error

	// Stop stops the database. This method must be called when the database is no longer needed.
	// Stop ensures that all non-flushed data is crash durable on disk before returning. Calls to
	// Put() concurrent with Stop() may not be crash durable after Stop() returns.
	Stop() error

	// Destroy deletes all data in the database.
	Destroy() error

	// UpdateTopology changes the on-disk topology of the data contained within the database. It returns a
	// channel that produces a value when the update is complete. The update is considered to be complete when
	// all data at removed paths has been fully migrated, and those paths are no longer being used by the DB
	// (and are therefore safe to repurpose or delete).
	//
	// shardingFactor is the number of DB shards to use. If a value of 0 is passed into this method, then the
	// previous sharding factor is used. Setting a sharding factor of 0 will keep the previous sharding factor.
	//
	// Setting the sharding factor to be smaller than the number of paths will mean that not each path will
	// have data being actively written to it at a particular point in time. Although this is unlikely to be
	// a useful way of configuring the DB in most cases, it may be desirable to use such a configuration if
	// the intent is to spread data across a large fleet of disks.
	//
	// paths is a list of directories where the database will store its data. If this list is nil or empty,
	// then the previous paths are used.
	//
	// If the paths list is missing paths that were previously used by the database, then the missing will stop being
	// used by the DB. In the background, the DB will relocate all files in the deleted paths to the remaining paths.
	// When it is finished relocating the data, there will be no files on the deleted paths that were previously used
	// by the DB.
	//
	// The first path in this list is called the "primary path". The primary path will contain keymap data if using a
	// keymap that writes data to disk. Note that changing the primary path for a keymap that stores data on disk
	// may be disruptive from a performance point of view. This is because the DB will need to freeze all ongoing
	// operations while it synchronously moves the keymap data to the new primary path. During this move, all read
	// and write operations will be blocked.
	//
	// If new paths are added to the list, then the DB will start writing data to the new paths. The DB will not
	// immediately move data from the old paths to the new paths (unless there are deleted old paths).
	UpdateTopology(shardingFactor uint32, paths []string) (chan struct{}, error)

	// SetGlobalReservedDriveCapacity sets the amount of disk space for all drives that the DB is not allowed to
	// fill up. If set to 100gb, then the DB will stop writing to any drive that has less than 100gb of free space.
	//
	// The default value is 64gb. If the reserved capacity is set to 0, then the DB will happily continue writing
	// until the drive is full to capacity.
	//
	// Note that space reservations may not protect the DB from filling up a disk if some other process rapidly
	// consumes disk space as well. The DB makes a best-effort attempt to avoid permitting a disk to fill up, but
	// can only guarantee this if it is the only process writing to the disk. Note that in general, it's safe to
	// allow other writers to the same disk as long as the reserved capacity is sufficiently large and the other
	// writers are not writing data in extreme quantities. If other writers limit their data size to less than the
	// reserved capacity, then the DB will always stop writing before the disk is 100% full.
	SetGlobalReservedDriveCapacity(reservedBytes uint64) error

	// SetReservedDriveCapacity sets the amount of disk space for the given path that the DB is not allowed to fill up.
	// Where SetGlobalReservedDriveCapacity sets the reserved capacity for all drives, this method allows for setting
	// the reserved capacity for individual drives, overriding the global setting.
	SetReservedDriveCapacity(path string, reservedBytes uint64) error

	// HardlinkBackup creates a backup of the database at the given path. The backup is a hardlink copy of the
	// database, and so it is very fast and uses very little disk space. Due to the nature of hardlinking, the
	// target path must be on the same volume as the database, and this method will return an error if it is not.
	// Furthermore, this method is fully incompatible with deployments that utilize multiple volumes.
	//
	// If this method is called targeting a path that already has a backup, then this operation is incremental.
	// This means that only new data will be copied to the backup path.
	//
	// If the copy of the database is later opened, the TTL for all tables will be set to 0 (i.e. TTL is infinite).
	// This is because it's plausible that the data in the backup could be quite old, and a backup is not very useful
	// if the first thing the DB does when it reads the backup is to delete all the data.
	//
	// Backups are atomic w.r.t. individual key-value pairs in the DB, but are not atomic as a whole. That is, if a
	// backup is interrupted mid-backup, some keys may be copied to the backup, while others may not. If a backup
	// is interrupted, running another backup against the same target path will resume the backup from where it
	// left off.
	HardlinkBackup(path string) error

	// LocalBackup creates a backup of the database at the given path(s). The backup is a full copy of the database.
	// This backup will not copy data faster than the given maxBytesPerSecond. If maxBytesPerSecond is 0, then there
	// is no limit to the speed of the backup.
	//
	// If this method is called targeting path(s) that already have a backup, then this operation is incremental. This
	// means that only new data will be copied to the backup path(s).
	//
	// If the copy of the database is later opened, the TTL for all tables will be set to 0 (i.e. TTL is infinite).
	// This is because it's plausible that the data in the backup could be quite old, and a backup is not very useful
	// if the first thing the DB does when it reads the backup is to delete all the data.
	//
	// Backups are atomic w.r.t. individual key-value pairs in the DB, but are not atomic as a whole. That is, if a
	// backup is interrupted mid-backup, some keys may be copied to the backup, while others may not. If a backup
	// is interrupted, running another backup against the same target path will resume the backup from where it
	// left off.
	LocalBackup(paths []string, maxBytesPerSecond uint64) error

	// RemoteBackup creates a backup of the database at the given socket. This backup will not copy data faster than
	// the given maxBytesPerSecond. If maxBytesPerSecond is 0, then there is no limit to the speed of the backup.
	//
	// If this method is called targeting a socket that already has a backup, then this operation is incremental. This
	// means that only new data will be copied to the backup socket.
	//
	// If the copy of the database is later opened, the TTL for all tables will be set to 0 (i.e. TTL is infinite).
	// This is because it's plausible that the data in the backup could be quite old, and a backup is not very useful
	// if the first thing the DB does when it reads the backup is to delete all the data.
	//
	// Backups are atomic w.r.t. individual key-value pairs in the DB, but are not atomic as a whole. That is, if a
	// backup is interrupted mid-backup, some keys may be copied to the backup, while others may not. If a backup
	// is interrupted, running another backup against the same target path will resume the backup from where it
	// left off.
	RemoteBackup(socket string, maxBytesPerSecond uint64) error
}
