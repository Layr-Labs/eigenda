# Filesystem Layout

This document provides an overview of how LittDB stores data on disk.

TODO: Talk about the following
- lock files
- snapshot files (both hard and soft links)

## Root Directories

LittDB spreads its data across N root directories. In practice, each root directory will probably be on its
own physical drive, but that's not a hard requirement.

In the example below, the root directories are named `root0`, `root1`, and `root2`.

## Table Directories

LittDB supports multiple [tables](#table), each with its own namespace. Each table is stored within its own
subdirectory.

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
├── rolling_snapshot
│   ├── tableA
│   │   └── segments
│   │       ├── 0-0.values -> root/root1/tableA/snapshot/0-0.values
│   │       ├── 0-1.values -> root/root2/tableA/snapshot/0-1.values
│   │       ├── 0-2.values -> root/root0/tableA/snapshot/0-2.values
│   │       ├── 0-3.values -> root/root1/tableA/snapshot/0-3.values
│   │       ├── 0.keys -> root/root0/tableA/snapshot/0.keys
│   │       ├── 0.metadata -> root/root0/tableA/snapshot/0.metadata
│   │       ├── 1-0.values -> root/root1/tableA/snapshot/1-0.values
│   │       ├── 1-1.values -> root/root2/tableA/snapshot/1-1.values
│   │       ├── 1-2.values -> root/root0/tableA/snapshot/1-2.values
│   │       ├── 1-3.values -> root/root1/tableA/snapshot/1-3.values
│   │       ├── 1.keys -> root/root0/tableA/snapshot/1.keys
│   │       ├── 1.metadata -> root/root0/tableA/snapshot/1.metadata
│   │       ├── 2-0.values -> root/root1/tableA/snapshot/2-0.values
│   │       ├── 2-1.values -> root/root2/tableA/snapshot/2-1.values
│   │       ├── 2-2.values -> root/root0/tableA/snapshot/2-2.values
│   │       ├── 2-3.values -> root/root1/tableA/snapshot/2-3.values
│   │       ├── 2.keys -> root/root0/tableA/snapshot/2.keys
│   │       └── 2.metadata -> root/root0/tableA/snapshot/2.metadata
│   ├── tableB
│   │   └── segments
│   │       ├── 0-0.values -> root/root1/tableB/snapshot/0-0.values
│   │       ├── 0-1.values -> root/root2/tableB/snapshot/0-1.values
│   │       ├── 0-2.values -> root/root0/tableB/snapshot/0-2.values
│   │       ├── 0-3.values -> root/root1/tableB/snapshot/0-3.values
│   │       ├── 0.keys -> root/root0/tableB/snapshot/0.keys
│   │       ├── 0.metadata -> root/root0/tableB/snapshot/0.metadata
│   │       ├── 1-0.values -> root/root1/tableB/snapshot/1-0.values
│   │       ├── 1-1.values -> root/root2/tableB/snapshot/1-1.values
│   │       ├── 1-2.values -> root/root0/tableB/snapshot/1-2.values
│   │       ├── 1-3.values -> root/root1/tableB/snapshot/1-3.values
│   │       ├── 1.keys -> root/root0/tableB/snapshot/1.keys
│   │       └── 1.metadata -> root/root0/tableB/snapshot/1.metadata
│   └── tableC
│       └── segments
├── root0
│   ├── litt.lock
│   ├── tableA
│   │   ├── keymap
│   │   │   ├── data
│   │   │   │   ├── 000001.log
│   │   │   │   ├── CURRENT
│   │   │   │   ├── LOCK
│   │   │   │   ├── LOG
│   │   │   │   └── MANIFEST-000000
│   │   │   ├── initialized
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
│   │   ├── snapshot
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
│   ├── tableB
│   │   ├── keymap
│   │   │   ├── data
│   │   │   │   ├── 000001.log
│   │   │   │   ├── CURRENT
│   │   │   │   ├── LOCK
│   │   │   │   ├── LOG
│   │   │   │   └── MANIFEST-000000
│   │   │   ├── initialized
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
│   │   ├── snapshot
│   │   │   ├── 0-2.values
│   │   │   ├── 0.keys
│   │   │   ├── 0.metadata
│   │   │   ├── 1-2.values
│   │   │   ├── 1.keys
│   │   │   └── 1.metadata
│   │   └── table.metadata
│   └── tableC
│       ├── keymap
│       │   ├── data
│       │   │   ├── 000001.log
│       │   │   ├── CURRENT
│       │   │   ├── LOCK
│       │   │   ├── LOG
│       │   │   └── MANIFEST-000000
│       │   ├── initialized
│       │   └── keymap-type.txt
│       ├── segments
│       │   ├── 0-2.values
│       │   ├── 0.keys
│       │   └── 0.metadata
│       ├── snapshot
│       └── table.metadata
├── root1
│   ├── litt.lock
│   ├── tableA
│   │   ├── segments
│   │   │   ├── 0-0.values
│   │   │   ├── 0-3.values
│   │   │   ├── 1-0.values
│   │   │   ├── 1-3.values
│   │   │   ├── 2-0.values
│   │   │   ├── 2-3.values
│   │   │   ├── 3-0.values
│   │   │   └── 3-3.values
│   │   └── snapshot
│   │       ├── 0-0.values
│   │       ├── 0-3.values
│   │       ├── 1-0.values
│   │       ├── 1-3.values
│   │       ├── 2-0.values
│   │       └── 2-3.values
│   ├── tableB
│   │   ├── segments
│   │   │   ├── 0-0.values
│   │   │   ├── 0-3.values
│   │   │   ├── 1-0.values
│   │   │   ├── 1-3.values
│   │   │   ├── 2-0.values
│   │   │   └── 2-3.values
│   │   └── snapshot
│   │       ├── 0-0.values
│   │       ├── 0-3.values
│   │       ├── 1-0.values
│   │       └── 1-3.values
│   └── tableC
│       ├── segments
│       │   ├── 0-0.values
│       │   └── 0-3.values
│       └── snapshot
└── root2
    ├── litt.lock
    ├── tableA
    │   ├── segments
    │   │   ├── 0-1.values
    │   │   ├── 1-1.values
    │   │   ├── 2-1.values
    │   │   └── 3-1.values
    │   └── snapshot
    │       ├── 0-1.values
    │       ├── 1-1.values
    │       └── 2-1.values
    ├── tableB
    │   ├── segments
    │   │   ├── 0-1.values
    │   │   ├── 1-1.values
    │   │   └── 2-1.values
    │   └── snapshot
    │       ├── 0-1.values
    │       └── 1-1.values
    └── tableC
        ├── segments
        │   └── 0-1.values
        └── snapshot
```
