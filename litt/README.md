Note: this document is a work in progress.

# What is LittDB?

Litt: adjective, slang, a synonym for "cool" or "awesome". e.g. "Man, that database is litt, bro!".

LittDB is a highly specialized embedded key-value store that is optimized for the following workload:

- high write throughput
- low read latency
- low memory usage
- write once, never update
- data is only deleted via a TTL (time-to-live) mechanism

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
- TTLs and automatic (lazy) deletion of expired values
- tables with non-overlapping namespaces
- thread safety: all methods are safe to call concurrently, and all modifications are atomic
- multi-drive support (data can be spread across multiple physical volumes)
- incremental backups (both local and remote)

## Planned Features

The following features are planned for future versions of LittDB:

- dynamic multi-drive support Drives can currently only be added/removed with a DB restart.
  It's fast, but not instantaneous. With this feature, drives can be added/removed on the fly.
- full snapshots/backups
- differential snapshots/backups
- read-only mode from an outside process
- CLI utility for managing the DB without the need for custom code
  (e.g. getting info, setting TTLs, adding/removing drives, etc.)

## Anti-Features

These are the features that littDB specifically does not provide, and will never provide:

- mutating existing values (once a value is written, it cannot be changed)
- deleting values (values only leave the DB when they expire via a TTL)
- transactions (individual operations are atomic, but there is no way to group operations atomically)
- fine granularity for TTL (all data in the same table must have the same TTL)
- multi-computer replication (littDB is designed to run on a single machine)

# Semantics & Definitions

TODO

# Configuration Options

TODO

# Architecture

TODO

# File Layout

## Roots

LittDB spreads its data across N root directories. In practice, each root directory will probably be on its
own physical drive, but that's not a hard requirement. Below, the first root is represented as `ROOT-0`, the
second as `ROOT-1`, and so on.

## Tables

LittDB supports multiple tables, each with its own namespace. Each table is stored within its own subdirectories
(each root may have a subdirectory for each table). Below, the first table is represented as `TABLE-0`, the second
as `TABLE-1`, and so on.

## Keymap

TODO

## Segments

TODO explain different types of segment files

## Example Layout

TODO perhaps generate a better view using tree on an example DB

- `ROOT-0`
    - `TABLE-0`
        - `table.metadata`
        - `ldb-keymap`
        - `segments`
            - `0.metadata`
            - `0.keys`
            - `0-0.values`
            - `1.metadata`
            - `1.keys`
            - `1-0.values`
            - ...
            - `N.metadata`
            - `N.keys`
            - `N-0.values`
    - `TABLE-1`
        - `table.metadata`
        - `ldb-keymap`
        - `segments`
            - `0.metadata`
            - `0.keys`
            - `0-0.values`
            - `1.metadata`
            - `1.keys`
            - `1-0.values`
            - ...
            - `N.metadata`
            - `N.keys`
            - `N-0.values`
    - ...
    - `TABLE-N`
        - `table.metadata`
        - `ldb-keymap`
        - `segments`
            - `0.metadata`
            - `0.keys`
            - `0-0.values`
            - `1.metadata`
            - `1.keys`
            - `1-0.values`
            - ...
            - `N.metadata`
            - `N.keys`
            - `N-0.values`
- `ROOT-1`
    - `TABLE-0`
        - `segments`
            - `0-1.values`
            - `1-1.values`
            - ...
            - `N-1.values`
    - `TABLE-1`
        - `segments`
            - `0-1.values`
            - `1-1.values`
            - ...
            - `N-1.values`
    - ...
    - `TABLE-N`
        - `segments`
            - `0-1.values`
            - `1-1.values`
            - ...
            - `N-1.values`
- ...
- `ROOT-N`
    - `TABLE-0`
        - `segments`
            - `0-N.values`
            - `1-N.values`
            - ...
            - `N-N.values`
    - `TABLE-1`
        - `segments`
            - `0-N.values`
            - `1-N.values`
            - ...
            - `N-N.values`
    - `TABLE-N`
        - `segments`
            - `0-N.values`
            - `1-N.values`
            - ...
            - `N-N.values`
