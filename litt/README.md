![](resources/littdb-logo.png)

# Contents

- [What is LittDB?](#what-is-littdb)
    - [Features](#features)
    - [Consistency Guarantees](#consistency-guarantees)
    - [Planned/Possible Features](#plannedpossible-features)
    - [Anti-Features](#anti-features)
- [API](#api)
    - [Overview](#overview)
    - [Getting Started](#getting-started)
    - [Configuration Options](#configuration-options)
- [Definitions](#definitions)
- [Architecture](#architecture)
    - [Big Picture Diagram](#putting-it-all-together-littdb) 
- [File Layout](#file-layout)

# What is LittDB?

LittDB is a highly specialized embedded key-value store that is optimized for the following workload:

- high write throughput
- low read latency
- low memory usage
- write once, never update
- data is only deleted via a [TTL](#ttl) (time-to-live) mechanism

In order to achieve these goals, LittDB provides an intentionally limited feature set. For workloads
that are capable of being handled with this limited feature set, LittDB is going to be more performant
than just about any other key-value store on the market. For workloads that require more advanced
features, "sorry, not sorry". LittDB is able to do what it does precisely because it doesn't provide
a lot of the features that a more general-purpose key-value store would provide, and adding those
can only be done by sacrificing the performance that LittDB is designed to provide.

## Features

The following features are currently supported by LittDB:

- writing values (once)
- reading values
- [TTLs](#ttl) and automatic (lazy) deletion of expired values
- [tables](#table) with non-overlapping namespaces
- multi-drive support (data can be spread across multiple physical volumes)
- incremental backups (both local and remote)

## Consistency Guarantees

The consistency guarantees provided by LittDB are more limited than those provided by typical general-purpose
transactional databases. This is intentional, as the intended use cases of LittDB do not require higher order
consistency guarantees.

- thread safety
- [read-your-writes consistency](#read-your-writes-consistency)
- crash [durability](#durability) for data that has been [flushed](#flushing)
- [atomic](#atomicity) writes
    - Although [batched writes](#batched-writes) are supported (for performance), batches are not [atomic](#atomicity).
      Each individual write within a batch is [atomic](#atomicity), but the batch as a whole is not. That is to say,
      if the computer crashes after a [batch](#batched-writes) has been written but before [flushing](#flushing),
      some of the writes in the [batch](#batched-writes) may be [durable](#durability) on disk, while others may
      not be.

## Planned/Possible Features

The following features are planned for future versions of LittDB, or are technically feasible if a strong
enough need is demonstrated:

- dynamic multi-drive support: Drives can currently only be added/removed with a DB restart.
  It's currently fast, but not instantaneous. With this feature, drives can be added/removed on the fly.
- full snapshots/backups
- differential snapshots/backups
- read-only mode from an outside process
- CLI utility for managing the DB without the need for custom code
  (e.g. getting info, setting TTLs, adding/removing drives, etc.)
- DB iteration (this is plausible to implement without high overhead, but we don't currently have
  a good use case to justify the implementation effort)
- more keymap implementations (e.g. badgerDB, a custom solution, etc.)
- data check-summing and verification (to protect/detect things like disk corruption)

## Anti-Features

These are the features that LittDB specifically does not provide, and will never provide. This is
not done because we're lazy, but because these features would significantly impact the performance
of the database, and because they are simply not needed for the intended use cases of LittDB. LittDB
is a highly specialized tool for a very specific task, and it is not intended to be a general-purpose
key-value store.

- mutating existing values (once a value is written, it cannot be changed)
- deleting values (values only leave the DB when they expire via a TTL)
- transactions (individual operations are atomic, but there is no way to group operations atomically)
- fine granularity for [TTL](#ttl) (all data in the same table must have the same TTL)
- multi-computer replication (LittDB is designed to run on a single machine)
- data encryption
- data compression
- any sort of query language other than "get me the value associated with this key"
- ordered data iteration

# API

## Overview

Below is a high level overview of the LittDB API. For more detailed information, see the inline documentation in the
interface files.

Source: [db.go](db.go)

```go
type DB interface {
    GetTable(name string) (Table, error)
    DropTable(name string) error
    Stop() error
    Destroy() error
}
```

Source: [table.go](table.go)

```go
type Table interface {
    Name() string
    Put(key []byte, value []byte) error
    PutBatch(batch []*types.KVPair) error
    Get(key []byte) ([]byte, bool, error)
    Flush() error
    SetTTL(ttl time.Duration) error
    SetCacheSize(size uint64) error
}
```

Source: [kv_pair.go](types/kv_pair.go)

```
type KVPair struct {
	Key []byte
	Value []byte
}
```

## Getting Started

Below is a functional example showing how to use LittDB.

```go
// Configure and build the database.
config, err := littbuilder.DefaultConfig("path/to/where/data/is/stored")
if err != nil {
    return err
}

db, err := config.Build(context.Background())
if err != nil {
    return err
}

myTable, err := db.GetTable("my-table") // this code works if the table is new or if the table already exists
if err != nil {
    return err
}

// Write a key-value pair to the table.
key := []byte("this is a key")
value := []byte("this is a value")

err = myTable.Put(key, value)
if err != nil {
    return err
}

// Flush the data to disk.
err = myTable.Flush()
if err != nil {
    return err
}

// Congratulations! Your data is now durable on disk.

// Read the value back. This works before or after a flush.
val, ok, err := myTable.Get(key)
if err != nil {
    return err
}
```

## Configuration Options

The "source of truth" for LittDB configuration documentation is the `Config` struct in
[littdb_config.go](littbuilder/littdb_config.go), although an overview is provided here.

Options marked with a `*` are options that are safe to ignore in most cases. In many cases, these options
are present to support testing.

- `Paths`: a list of directories. LittDB will do its best to spread data across these directories.
  Directories may or may not be on the same physical drive.
- `LoggerConfig`: a struct containing configuration options for the logger. A sane default is provided.
- `KeymapType`: the type of [keymap](#keymap) to use. The default is `keymap.LevelDBKeymapType`. An in-memory
  is also supported: `keymap.MemKeymapType`. It will be faster, but may have longer startup times and higher
  memory usage.
- `TTL`: the [time-to-live](#ttl) for data in the database. If set to `0`, data will never expire. The default
  is `0`. Changing the [TTL](#ttl) for a table effects all data currently in the table, as well as all data
  written to the table in the future. Lowering the [TTL](#ttl) may cause some data to immediately become eligible
  for deletion.
- `ControlChannelSize`*: the size of an internal channel used for controlling the database. The default is `64`.
- `TargetSegmentFileSize`*: the target size for segment files. The default is `2^32` (4GB).
- `GCPeriod`*: the frequency at which the garbage collector runs. The default is `1m`.
- `ShardingFactor`: the number of shards to use for each segment. The default is `1`. If more than one path is provided
  in the `Paths` field, then shards will be spread out across the available paths. If there are more shards than paths,
  some paths will have more than one shard. If there are more paths than shards, some paths will have no shards.
- `SaltShaker`*: a random number generator used to generate [sharding salt](#sharding-salt). Default is a standard
  PRNG seeded with the current time.
- `CacheSize`: the size of the [in-memory cache](#cache), in bytes. The default 1GB.
- `TimeSource`*: a function that returns the current time. The default is `time.Now()`.
- `Fsync`: if true, then each flush operation performs an fsync operation. Ensures the data is durable even if the
  OS crashes. Otherwise, there may be data in the OS's internal buffers. This may cause significant performance
  overhead. The default is `false`.

# Definitions

This section contains an alphabetized list of technical definitions for a number of terms used by LittDB. This
list is not intended to be read in order, but rather to be used as a reference when reading other parts of the
documentation.

## Address

An address partially describes the location on disk where a [value](#value) is stored. Together with a [key](#key),
the [value](#value) associated with a [key](#key) can be retrieved from disk.

An address is encoded in a 64-bit integer. It contains two pieces of information:

- the [segment](#segment) [index](#segment-index) where the [value](#value) is stored
- the offset within the [value file](#segment-value-files) where the first byte of
  the [value](#value) is stored

This information is not enough by itself to retrieve the [value](#value) from disk if there is more than one
[shard](#shard) in the [table](#table). When there is more than one [shard](#shard), the following information
must also be known in order to retrieve the [value](#value) (i.e. to figure out which [shard](#shard) to look in):

- the [sharding factor](#sharding-factor) for the [segment](#segment) where the [value](#value) is stored
  (stored in the [segment metadata file](#segment-metadata-file))
- the [sharding salt](#sharding-salt) for the [table](#table) where the [value](#value) is stored
  (stored in the [table metadata file](#table-metadata-file))
- the [key](#key) that the [value](#value) is associated with

## Atomicity

In the context of this document, atomicity means that an operation is either done completely or not at all. That is
to say, if there is a crash while an operation is in progress, the operation will either be completed when the
database is restarted, or it will not be completed at all.

As a specific example, if writing a [value](#value) and there is a crash, either the entire [value](#value) will be
written to disk and available when the database is restarted, or the [value](#value) will be completely absent.
It will never be the case that only part of the [value](#value) is written to disk.

## Cache

LittDB maintains an in-memory cache of [key](#key)-[value](#value) pairs. Data is stored in this cache when a value
is first written, as well as when it is read from disk. This is not needed for correctness, but is rather a performance
optimization. The cache is not persistent, and is lost when the database is restarted. The size of the cache is
configurable.

## Batched Writes

LittDB supports bated write operations. Multiple write operations can be grouped together and passed to the database
as a single operation. This may have positive performance implications, but is semantically equivalent to writing each
value individually. A batch of writes is not [atomic](#atomicity) as a whole, but each individual write within the
batch is [atomic](#atomicity). That is to say, if there is a crash after a batch of writes has been written but before
it has been [flushed](#flushing), some of the writes in the batch may be [durable](#durability) on disk, while others
may not be.

## Durability

In this context, the term "durable" is used to mean that data is stored on disk in such a way that it will not be lost
in the event of a crash. Data that has been [flushed](#flushing) is considered durable. Data that has not been flushed
is not considered durable. That doesn't mean that the data will be lost in the event of a crash, but rather that it
is not guaranteed to be present after a crash.

There are some limits to the strength of the durability guarantee provided by LittDB. For example, some drives buffer
data in internal buffers before writing it to disk, and do not necessarily write data to disk immediately. LittDB is
only as robust as the OS/hardware it is running on. This is true for any database, but it is worth mentioning here
for the sake of completeness.

## Flushing

Calling `Flush()` causes all data previously written to be written [durably](#durability) to disk. A call to `Flush()`
blocks until all data that was written prior to the call to `Flush()` has been written to disk.

It is ok to never call `Flush()`. As internal buffers fill, data is written to disk automatically. However, calling
`Flush()` can be useful in some cases, such as when you want to ensure that data is written to disk before proceeding
with other operations.

`Flush()` makes no guarantees about the [durability](#durability) of data written concurrently with the call to
`Flush()` or after the call to `Flush()` has returned. It's not harmful to write data concurrently with a call to
`Flush()` as long as it is understood that this data may or may not be [durable](#durability) on disk when the call
to `Flush()` returns.

## Key

A key in a key-[value](#value) store. A key is a byte slice that is used to look up a [value](#value) in the database.

LittDB is agnostic to the contents of the key, other than requiring that keys be unique within a [table](#table).
Although large keys are supported, performance has been tuned under the assumption that keys are generally small
compared to [values](#value). The use case LittDB was originally intended for uses 32-byte keys.

## Keymap

At a conceptual level, a keymap is a mapping from [keys](#key) to [addresses](#address). In order to look up a
[value](#value) in the database one needs to know two things: the [key](#key) and the [address](#address). The keymap
is therefor necessary to lookup data given a specific [key](#key).

There are currently two implementations of the keymap in LittDB: an in-memory keymap and a keymap that uses levelDB.
There are tradeoffs to each implementation. The in-memory keymap is faster, but has higher memory usage and longer
startup times (it has to be rebuilt at boot time). The levelDB keymap is slower, but has a lower memory footprint and
faster startup times.

From a thread safety point of view, if a mapping is present in the keymap, the [value](#value) associated with the
entry is guaranteed to be present on disk.

- When writing a new [value](#value), it is first written to disk, and when that is complete the [key](#key) and
  [address](#address) are written to the keymap.
- When deleting a [value](#value), the [key](#key) and [address](#address) are first removed from the keymap, and
  then the [value](#value) is deleted from disk.

LittDB supports reading [values](#value) immediately after they are written, and during that period there may not
be a corresponding entry in the keymap. For more information on how this edge case is handled, information about the
[unflushed data map](#unflushed-data-map).

## Read-Your-Writes Consistency

The definition of read-your-writes consistency is well summarized by its name. If a thread writes a [value](#value)
to the database and then turns around and attempts to read that [value](#value) back, it will either

1. read the [value](#value) that was just written, or
2. read an updated [value](#value) that was written AFTER the [value](#value) that was just written

Note that in LittDB, values are never updated. In this context, an "updated" value the absence of a value
when it eventually outlives its [TTL](#ttl) and is deleted by the garbage collector.

An "eventual consistent" database does not necessarily provide read-your-writes consistency. In the author's experience,
such systems can be very difficult to reason about, and can lead to subtle bugs that are difficult to track down.
Read-your-writes consistency is simple, yet powerful and intuitive. Since providing this level of consistency
does not hurt performance, and so the complexity of its implementation is justified.

## Segment

Data in LittDB [table](#table) can be visualized as a linked list. Each element in that linked is called a "segment". 
A segment can hold many individual [values](#value). Old data is near the beginning of the list, and new data is near 
the end. Old, [expired](#ttl) data is always deleted from the first segment currently in the list. New data is always 
written to the last segment currently in the list.

Segments are deleted as a whole. That is, when a segment is deleted, all data in that segment is deleted at the same
time. Segments are only deleted when all data contained within them has [expired](#ttl).

Segments have a maximum data size. When a segment is full, that segment is made immutable, and a new segment is created
and added to the end of the list.

Each segment may split its data into multiple [shards](#shard). The number of shards in a segment is called the
[sharding factor](#sharding-factor). The [sharding factor](#sharding-factor) is configurable, and different segments
may use different [sharding factors](#sharding-factor).

There are three types of files that contain data for a segment:

- [metadata](#segment-metadata-file)
- [keys](#segment-key-file)
- [values](#segment-value-files)

### Segment Index

Each segment has a serial number called a "segment index". The first segment ever created with index `0`, the next
segment created has index `1`, and so on. Segment `N` is always deleted before segment `N+1`, meaning there will
never be a gap in the segment indices currently in use.

### Segment Key File

A segment key file contains the [keys](#key) and [addresses](#address) for all the [values](#value) stored the segment.
At runtime, [keys](#key)-[address](#address) pairs are appended to the key file. It is not read except during the
following circumstances:

- when a [segment](#segment) is deleted, the file is iterated to delete entries from the [keymap](#keymap)
- when the DB is loaded from disk, the data is used to rebuild the [keymap](#keymap). This may not be needed
  in situations where the keymap has durably stored data, and does not need to be rebuilt.

The file name of a key file is `X.keys`, where `X` is the [segment index](#segment-index).

### Segment Metadata File

This file contains metadata about the segment. This metadata is small, and so it can be kept in memory. The file is
read at startup to rebuild the in-memory representation of the segment.

Each metadata contains the following information:

- the [segment index](#segment-index)
- serialization version (in case the format changes in the future)
- the [sharding factor](#sharding-factor) for the segment
- the [salt](#sharding-salt) used for the segment
- the [timestamp](#segment-timestamp) of the last element written in the segment.
  the [TTL](#ttl) of any data contained within it.
- whether or not the segment is [immutable](#segment-mutability)

The file name of a metadata file is `X.metadata`, where `X` is the [segment index](#segment-index).

### Segment Mutability

Only the last segment in the "linked list" is mutable. All other segments are immutable.

### Segment Timestamp

The timestamp of the last element written to the segment. This is used to determine when it is safe to delete a
segment without violating the [TTL](#ttl) of any data contained within it. This value is unset for the last segment
in the list, as it is still being written to.

### Segment Value Files

Each segment has one value file for each [shard](#shard) in the segment. Values are appended to the value files.
The [address](#address) of a [value](#value) is the offset within the value file where the [value](#value) begins.

The file name of a value file is `X-Y.values`, where `X` is the [segment index](#segment-index) and `Y` is the
[shard](#shard) index.

## Shard

LittDB supports sharding. That is to say, it can break the data into smaller pieces and spread those pieces across
multiple locations.

In order to determine the shard that a particular [key](#key) is in, a hash function is used. The data that goes
into the hash function is the [key](#key) itself, as well as a [sharding salt](#sharding-salt) that is unique to
each [segment](#segment). 

The [sharding salt](#sharding-salt) is chosen randomly. Its purpose is to make the mapping between [keys](#key) and 
shards unpredictable to an outside attacker. Without this sort of randomness, an attacker could intentionally craft 
keys that all map to the same shard, causing a hot spot in the database and potentially degrading performance.

### Sharding Factor

The number of [shards](#shard) in a [segment](#segment) is called the "sharding factor". The sharding factor must be
a positive, non-zero integer. The sharding factor can be changed at runtime without restarting the database or 
performing a data migration.

### Sharding Salt

A random number chosen to make the [shard](#shard) hash function unpredictable to an outside attacker. This number
does not need to be chosen via a cryptographically secure random number generator, as long as it is not publicly
known.

## Table

A table in LittDB is a unique namespace. Two [keys](#key) with identical values do not conflict with each other as
long as they are in different tables.

Each table has its own [TTL](#ttl), and all data in the table is subject to that [TTL](#ttl). Each table has its 
own [keymap](#keymap) and its own set of [segments](#segment). [Flushing](#flushing) one table does not affect 
any other table. Aside from hardware, tables do not share any resources.

In many ways, a table is a stand-alone database. The higher level [API](#api) that works with multiple tables is 
provided as a convenience, but does not enhance the performance of the DB in any way.

### Table Metadata File

A [table](#table) metadata file contains configuration for the table. It is intended to preserve high level
configuration between restarts.

## TTL

TTL stands for "time-to-live". If data is configured to have a TTL of X hours, the data is automatically deleted
approximately X hours after it is written.

Note that TTL is the only way littDB supports removing data from the database. Although it is legal to configure
a table with a TTL of 0 (i.e. where data never expires), such a table will never be able to remove data.

## Unflushed Data Map

An in-memory map that contains [keys](#key)-[values](#value) pairs that are not yet [durable](#durability) on disk.
Entries are added to the map when a [value](#value) is written, and removed when the [value](#value) is fully
written to both the [keymap](#keymap) and the [segment](#segment) files.

This data structure is not to be confused with the [cache](#cache). Its purpose is not to improve performance, but
rather to provide [read-your-writes consistency](#read-your-writes-consistency).

## Value

The value in a key-[value](#value) store. A value is a byte slice that is associated with a [key](#key) in the database.
LittDB is optimized to support large values, although small values are perfectly fine as well. Writing the X bytes
of data as a single large value is more efficient than writing X bytes of data as Y smaller values.

# Architecture

This section explains the high level architecture of LittDB. It starts out by describing a simple (but inefficient) 
storage solution, and incrementally adds complexity in order to solve various problems. For the full picture, skip to 
[Putting it all together: LittDB](#putting-it-all-together-littdb).

For each iteration, the database is must fulfill the following requirements:
- must support `put(key, value)`/`get(key)` operations
- must be thread safe
- must support a TTL
- must be crash durable

## Iteration 1: Appending data to a file

Let's implement the simplest possible key-value store that satisfies the requirements above. It's going to be super
slow. Ok, fine. We want simple.

![](./resources/iteration1.png)

When the user writes a key-value pair to the database, append the key and the value to the end of the file, along
with a timestamp. When the user reads a key, scan the file from the beginning until you find the key and 
return the value.

Periodically, scan the data in the file to check for expired data. If a key has expired, remove it from the file 
(will require the file to be rewritten).

This needs to be thread safe. Keep a global read-write lock around the file. When a write or GC operation is in 
progress, no reads are allowed. GC operations and writes are not permitted to happen in parallel. Allow multiple 
reads to happen concurrently.

In order to provide durability, ensure the file is fully flushed to disk before releasing a write lock.

Congratulations! You've written your very own database!

![](./resources/iDidIt.png)

## Iteration 2: Add a cache

Reads against the database in 1 are slow. If there is any way we could reduce the number of times we have to iterate 
over the file, that would be great. Let's add an in-memory cache.

![](./resources/iteration2.png)

Let's assume we are using a thread safe map to implement the cache.

When reading data, first check to see if the data is in the cache. If it is, return it. If it is not, acquire a read
lock and scan the file. Be sure to update the cache with the data you read.

When writing data, write the data to the file, and then update the cache. Data that is recently written is often
read immediately shortly after, for many workloads.

When deleting data, remove the data from the file, and then remove the data from the cache.

## Iteration 3: Add an index

Reading recent values is a lot faster now. But if you miss the cache, things start getting slow. `O(n)` isn't fun
when you database holds 100TB. To address this, let's add an index that allows us to jump straight to the data we
are looking for. For the sake of consistency with other parts of this document, let's call this index a "keymap".

![](./resources/iteration3.png)

Inside the keymap, maintain a mapping from each key to the offset in the file where the first byte of the value is
stored.

When writing a value, take note of the offset in the file where the value was written. Store the key and the offset
in the keymap.

When reading a value and there is a cache miss, look up the key inside the keymap. If the key is present, jump to
start of the value in the file and read the value. If the key is not present, tell the user that the key is not
present in the database.

When deleting a value, remove the key from the keymap in addition to removing the value from the file.

At startup time, we will have to rebuild the keymap, since we are only storing it in memory. In order to do so,
iterate over the file and reconstruct the keymap. If this is too slow, consider storing the keymap on disk (perhaps
using an off-the-shelf key-value store like levelDB).

The database needs to do a little extra bookkeeping when it deletes data from the file. If it deletes X bytes from
the beginning of the file, then the offsets recorded in the keymap are off by a factor of X. The key map doesn't
need to be rebuilt in order to fix this. Rather, the database can simply subtract X from all the offsets in the
keymap to find the actual location of the data in the file. Additionally, it must add X to the offset when computing
the "offset" of new data that is written to the file.

## Iteration 4: Unflushed data map

In order to be thread safe, the solution above uses a global lock. While one thread is writing, readers must wait
unless they get lucky and find their data in the cache. It would be really nice if we could permit reads to continue
uninterrupted while writes are happening in the background.

![](./resources/iteration4.png)

Create another key->value map called the "unflushed data map". Use a thread safe map implementation.

When the user writes data to the database, immediately add it to the unflushed data map, but not the key map. 
After that is completed, write it to file. The write doesn't need to be synchronous. For example, you can use file 
stream APIs that buffer data in memory before writing it to disk in larger chunks. The write operation doesn't need 
to block until the data is written to disk, it can return as soon as the data is in the unflushed data map and written 
to the buffer.

Expose a new method in the database called `Flush()`. When `Flush()` is called, first flush all data in buffers to disk,
then empty out the unflushed data map. Before each entry is removed, write the key-address pair to the key map.
This flush operation should block until all of this work is done.

When reading data, look for it in the following places, in order:
- the cache
- the unflushed data map
- on disk (via the keymap and data file)

Unlike previous iterations, write no longer needs to hold a lock that blocks readers. This is thread safe, and it
provides read-your-writes consistency.

If a reader is attempting to read data that is currently in the process of being written to disk, then the data will
be present in the unflushed data map. If the reader finds an entry in the key map, this means that the data has already
been written out to disk, and is therefore safe to read from the file. Even if the writer is writing later in the file,
the bytes the reader wants to read will be immutable.

Although the strategy described above allows read operations to execute concurrently with write operations, it does
not solve the problem for deletions of values that have exceeded their TTL. This operation will still require a global
lock that blocks all reads and writes.

## Iteration 5: Break the file into segments

One of the biggest inefficiencies in design, to this point, is that the deleting values is exceptionally slow. The
entire file must be rewritten in order to trim bytes from the beginning. And to make matters worse, we need to hold
a global lock while we do it. To fix this, let's break apart the data file into multiple data files. We'll call each
data file a "segment".

![](./resources/iteration5.png)

Decide on a maximum file size for each segment. Whenever a file gets "full", close it and open a new one. Let's assign
each of these files a serial number starting with `0` and increasing monotonically. We'll call this serial number the
"segment index".

Previously, the address stored in the key map told us the offset in the file where the value was stored. Now, the
address will also need to keep track of the segment index, as well as the offset.

Deletion of data is now super easy. When all data in the oldest segment file exceeds its TTL, we can delete just that
segment without modifying any of the other segment files. Iterate over the segment file to delete values from the key 
map.

In order to avoid the race condition where a reader is reading data from a segment that is in the process of being 
deleted, use reference counters for each segment. When a reader goes to read data, it first finds the address in the
keymap, than increments the reference counter for the segment. When the reader is done reading, it decrements the
reference counter. When the garbage collector goes to delete a segment, it waits to actually delete the file on disk
until the reference counter is zero. As a result of this strategy, there is now no longer a need for garbage collection
to hold a global lock.

## Iteration 6: Metadata files

In the previous iteration, we do garbage collection by deleting a segment once all data contained within that segment
has expired. But how do we figure out when that actually is? In the previous iteration, the only way to do that is to
iterate over the entire segment file and read the timestamp stored with the last value. Let's do better.

For each segment, let's create a metadata file. We'll put the timestamp of the last value written to the segment into
this file. As a result, we will no longer need to store timestamp information inside the value files, which will
save us a few bytes per entry.

![](./resources/iteration6.png)

Now, all the garbage collector needs to read to decide when it is time to delete a segment is the metadata file for
that segment.

## Iteration 7: Key files

Storing timestamp information in a metadata file is a good start, but we still need to scan the value file. When a
segment is deleted, we need to clean up the key map. In order to do this, we need to know which keys are stored in the
segment. Additionally, when we start up the database, we need to rebuild the key map. This requires us to scan each
segment file to find the keys.

From an optimization point of view, we can assume that in general keys will be much smaller than values. During the
operations described above, we don't care about the values, only the keys. So lets separate the keys from the values
to avoid having to read the values when we don't need them.

![](./resources/iteration7.png)

Everything works the same way as before. But instead of iterating huge segment files when deleting a segment
or rebuilding the key map at startup, we only have to iterate over the key file. The key file is going to be 
significantly smaller than the value file (for sane key-value size ratios), and so this will be much faster.

## Iteration 8: Sharding

A highly desirable property for this database is the capability to spread its data across multiple physical drives.
In order to do this, we need to shard the data. That is to say, we need to break the data into smaller pieces and
spread those pieces across multiple locations.

![iteration8](./resources/iteration8.png)

Key files and metadata files are small. For the sake of simplicity, let's not bother sharding those. Value files
are big. Break apart value files, and have one value file per shard.

When writing data, the first thing to do will be to figure out which shard the data belongs in. Do this by taking a
hash of the key modulo the number of shards.

When reading data, we need to do the reverse. Take a hash of the key modulo the number of shards to figure out which
shard to look in. As a consequence, the address alone is no longer enough information to find the data. We also need
to know the key when looking up data. But this isn't a problem, since we always have access to the key when we are
looking up data.

From a security perspective, sharding with a predictable hash is dangerous. An attacker could, in theory, craft keys
that all map to the same shard, causing a hot spot in the database. To prevent this, the database chooses a random
"salt" value that it includes in the hash function. As long as an attacker does not know the salt value, they cannot
predict which shard a key will map to.

We already have a metadata file for each segment. We can go ahead and safe the sharding factor and salt in the metadata
file. This will give us enough information to find data contained within the segment.

## Iteration 9: Multi-table support

A nice-to-have feature would be the ability to support multiple tables. Each table would have its own namespace, and
data in one table would not conflict with data in another table.

This is simple! Let's just run a different DB instance for each table.

![](./resources/iteration9.png)

Since each table might want to have its own configuration, we can store that configuration in a metadata file for each
table.

## Putting it all together: LittDB

![littdb](./resources/littdb-big-picture.png)

# File Layout

This section provides an overview of how LittDB stores data on disk.

## Root Directories

LittDB spreads its data across N root directories. In practice, each root directory will probably be on its
own physical drive, but that's not a hard requirement.

In the example below, the root directories are named `root0`, `root1`, and `root2`.

## Table Directories

LittDB supports multiple [tables](#table), each with its own namespace. Each table is stored within its own subdirectory.

In the example below, there are three tables: `tableA`, `tableB`, and `tableC`.

## Keymap Directory

All keymap data appears in the directory named `keymap`. The file `keymap-type.txt` contains the name of the
keymap implementation. If the keymap writes data to disk (e.g. levelDB, as pictured below), then the data will
be stored in the `keymap/data` directory.

## Segment Files

There are three types of files that contain data for a [segment](#segment):

- [metadata](#segment-metadata-file)
- [keys](#segment-key-file)
- [values](#segment-value-files)

## Example Layout

The following is an example file tree for a simple LittDB instance.
(This example file tree was generated using generate_example_tree_test.go.)

There are three directories into which data is written. In theory, these could be located on three separate
physical drives. Those directories are

- `root/root0`
- `root/root1`
- `root/root2`

The table is configured to have four shards. That's one more shard than root directory, meaning that one of the
root directories will have two shards, and all the others will have one shard.

There are three tables, each with its own namespace. The tables are

- `tableA`
- `tableB`
- `tableC`

A little data has been written to the DB.

- `tableA` has enough data to have three segments
- `tableB` has enough data to have two segments
- `tableC` has enough data to have one segment

The keymap is implemented using levelDB.

```text
root
├── root0
│   ├── tableA
│   │   ├── keymap
│   │   │   ├── data
│   │   │   │   ├── 000001.log
│   │   │   │   ├── CURRENT
│   │   │   │   ├── LOCK
│   │   │   │   ├── LOG
│   │   │   │   └── MANIFEST-000000
│   │   │   └── keymap-type.txt
│   │   ├── segments
│   │   │   ├── 0-2.values
│   │   │   ├── 0.keys
│   │   │   ├── 0.metadata
│   │   │   ├── 1-2.values
│   │   │   ├── 1.keys
│   │   │   ├── 1.metadata
│   │   │   ├── 2-2.values
│   │   │   ├── 2.keys
│   │   │   ├── 2.metadata
│   │   │   ├── 3-2.values
│   │   │   ├── 3.keys
│   │   │   └── 3.metadata
│   │   └── table.metadata
│   ├── tableB
│   │   ├── keymap
│   │   │   ├── data
│   │   │   │   ├── 000001.log
│   │   │   │   ├── CURRENT
│   │   │   │   ├── LOCK
│   │   │   │   ├── LOG
│   │   │   │   └── MANIFEST-000000
│   │   │   └── keymap-type.txt
│   │   ├── segments
│   │   │   ├── 0-2.values
│   │   │   ├── 0.keys
│   │   │   ├── 0.metadata
│   │   │   ├── 1-2.values
│   │   │   ├── 1.keys
│   │   │   ├── 1.metadata
│   │   │   ├── 2-2.values
│   │   │   ├── 2.keys
│   │   │   └── 2.metadata
│   │   └── table.metadata
│   └── tableC
│       ├── keymap
│       │   ├── data
│       │   │   ├── 000001.log
│       │   │   ├── CURRENT
│       │   │   ├── LOCK
│       │   │   ├── LOG
│       │   │   └── MANIFEST-000000
│       │   └── keymap-type.txt
│       ├── segments
│       │   ├── 0-2.values
│       │   ├── 0.keys
│       │   └── 0.metadata
│       └── table.metadata
├── root1
│   ├── tableA
│   │   └── segments
│   │       ├── 0-0.values
│   │       ├── 0-3.values
│   │       ├── 1-0.values
│   │       ├── 1-3.values
│   │       ├── 2-0.values
│   │       ├── 2-3.values
│   │       ├── 3-0.values
│   │       └── 3-3.values
│   ├── tableB
│   │   └── segments
│   │       ├── 0-0.values
│   │       ├── 0-3.values
│   │       ├── 1-0.values
│   │       ├── 1-3.values
│   │       ├── 2-0.values
│   │       └── 2-3.values
│   └── tableC
│       └── segments
│           ├── 0-0.values
│           └── 0-3.values
└── root2
    ├── tableA
    │   └── segments
    │       ├── 0-1.values
    │       ├── 1-1.values
    │       ├── 2-1.values
    │       └── 3-1.values
    ├── tableB
    │   └── segments
    │       ├── 0-1.values
    │       ├── 1-1.values
    │       └── 2-1.values
    └── tableC
        └── segments
            └── 0-1.values
```
