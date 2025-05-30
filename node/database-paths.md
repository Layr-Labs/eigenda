# Configuring Validator Storage Paths

A validator node is responsible for storing chunk data on disk, and for making that data available when requested.
Until the V1 protocol is fully deprecated (in favor of the V2 protocol introduced in the `Blazar` release), a validator
node will store chunk data for both V1 and V2 protocols. The way that data is managed is different between the V1
protocol and the V2 protocol.

The location on disk where this data is stored is configured by the following two flags:

- `NODE_DB_PATH`: This flag specifies the path where the V1 protocol chunk data is stored. This flag should
  contain a fully qualified path to a directory where the V1 protocol should store its chunk data.
- `NODE_LITT_DB_STORAGE_PATHS`: This flag specifies the path where the V2 protocol chunk data is stored.
  Unlike V1, the V2 data storage engine (LittDB) is capable of spreading data across multiple directories.
  These directories do not need to be on the same filesystem (e.g. if you want to use multiple disks).
  To pass in multiple directories, provide a comma-separated list. Each directory should be a fully qualified path.

Until the V1 protocol is fully deprecated, `NODE_DB_PATH` must be set. 

Technically, the new flag `NODE_LITT_DB_STORAGE_PATHS` is optional, since if it is not set then the validator 
software will store its data in the location specified by `NODE_DB_PATH`. This is not recommended. Eventually,
the V1 protocol will be disabled, and the `NODE_DB_PATH` flag will be removed along with it. In order to be future
proof, it is highly recommended to set the `NODE_LITT_DB_STORAGE_PATHS` flag.

# File System Layout

## V1 Protocol

The V1 protocol's disk footprint looks something like this:

```
${NODE_DB_PATH}
├── chunk
│   ├── 000001.log
│   ├── CURRENT
│   ├── LOCK
│   ├── LOG
│   └── MANIFEST-000000
```

The `chunk` directory is created by the V1 software inside the directory specified by `NODE_DB_PATH`. Inside
the `chunk` directory are files maintained by the V1 data storage engine (i.e. `LevelDB`).

## V2 Protocol

The V2 protocol's disk footprint depends on how it is configured.

### Deprecated Configuration: only `NODE_DB_PATH` set

If only `NODE_DB_PATH` is set and `NODE_LITT_DB_STORAGE_PATHS` is not set (not recommended!), then the V2 protocol
will store its data like this:

```
${NODE_DB_PATH}
├── chunk_v2_litt
│   └── chunks
│       ├── keymap
│       │   ├── data
│       │   │   ├── 000001.log
│       │   │   ├── CURRENT
│       │   │   ├── LOCK
│       │   │   ├── LOG
│       │   │   └── MANIFEST-000000
│       │   ├── initialized
│       │   └── keymap-type.txt
│       ├── segments
│       │   ├── 0.keys
│       │   └── 0.metadata
│       └── table.metadata
```

The `chunk_v2_litt` directory is created by the V2 software inside the directory specified by `NODE_DB_PATH`.
The `chunks` directory is created and maintained by the V2 data storage engine (i.e. `LittDB`).

### Recommended Configuration: `NODE_LITT_DB_STORAGE_PATHS` set

Suppose `NODE_LITT_DB_STORAGE_PATHS` is provided 3 paths: `${volume1}`, `${volume2}`, and `${volume3}`.

```
${volume1}
   └── chunks
       ├── keymap
       │   ├── data
       │   │   ├── 000001.log
       │   │   ├── CURRENT
       │   │   ├── LOCK
       │   │   ├── LOG
       │   │   └── MANIFEST-000000
       │   ├── initialized
       │   └── keymap-type.txt
       ├── segments
       │   ├── 0-2.values
       │   ├── 0.keys
       │   └── 0.metadata
       └── table.metadata

${volume2}
   └── chunks
       └── segments
           └── 0-0.values

${volume3}
   └── chunks
       └── segments
           └── 0-1.values
```

In each of the directories specified by `NODE_LITT_DB_STORAGE_PATHS`, a `chunks` directory is created and maintained
by the V2 data storage engine (i.e. `LittDB`).

Notice that the first volume has more files than the other two volumes. LittDB selects one of the volumes to store
metadata files. In the other volumes, it only stores values files (i.e. the `*.values` files). 99.99% of the 
data written to disk is stored in the `*.values` files, so disk utilization across volumes is fairly even.

# Changing `NODE_DB_PATH`

It's possible to change the `NODE_DB_PATH` after it has been set with the following manual steps:

- Stop the validator node.
- Copy/move the contents of the old `NODE_DB_PATH` to the new intended `NODE_DB_PATH`, e.g. `mv /old/path/ /new/path/`.
- Update the `NODE_DB_PATH` environment variable to point to the new path.
- Restart the validator node.

# Changing `NODE_LITT_DB_STORAGE_PATHS`

## Adding a Path

It's possible to add additional paths to `NODE_LITT_DB_STORAGE_PATHS`. This might be useful if want to add
additional storage space by adding additional disks. To do this, do the following:

- Stop the validator node.
- Update the `NODE_LITT_DB_STORAGE_PATHS` environment variable to include the new path(s). This flag
  accepts a comma-separated list of paths.
- Restart the validator node.

In the future, the data storage engine will get an upgrade that allows it to write to new paths without restarting
the validator software. Stay tuned for more info!

## Removing a Path

Removing a path from `NODE_LITT_DB_STORAGE_PATHS` is more involved, but still possible. In order to remove a path,
it is necessary to move all data from the path you want to remove into a path that you want to keep. The contents
of the `chunks` directories must be merged. The data storage engine (LittDB) always uses unique file names across
all paths, so there will be no file name conflicts.

- Stop the validator node.
- Move the data out of the path you want to remove into one of the paths you want to keep. Merge the contents
  of the `chunks` directories.
- Update the `NODE_LITT_DB_STORAGE_PATHS` environment variable to remove the path you want to remove.
- Restart the validator node.

In the future, the data storage engine will get an upgrade that allows it to remove paths without restarting
the validator software. This update will also streamline this process and will remove the need to manually
merge the contents of the `chunks` directories. Stay tuned for more info!

## Oops! I didn't initially set `NODE_LITT_DB_STORAGE_PATHS`, how do I fix this?

If you initially run a validator node without setting `NODE_LITT_DB_STORAGE_PATHS`, the V2 protocol will
store its data in the same location as the V1 protocol, i.e. in the directory specified by `NODE_DB_PATH`.
If you later decide to set `NODE_LITT_DB_STORAGE_PATHS`, manual steps are required or else the validator node
may lose data. (This is why it's highly recommended to set `NODE_LITT_DB_STORAGE_PATHS` from the start!)

If using the legacy fallback `NODE_DB_PATH` for V2 data storage, the validator software stores its data at
`${NODE_DB_PATH}/chunk_v2_litt/`. The `chunk_v2_litt` directory is hard coded and always added to the path.
But if you are using the new `NODE_LITT_DB_STORAGE_PATHS` flag, the `chunk_v2_litt` directory is NOT added to the path.

In order to remedy this, you will need to move the contents of `${NODE_DB_PATH}/chunk_v2_litt/` to the location where
you want to store the V2 data. For example, if you want to store the V2 data in `${volume1}`, then you would
do `mv ${NODE_DB_PATH}/chunk_v2_litt/ ${volume1}/`.

- Stop the validator node.
- Move the contents of `${NODE_DB_PATH}/chunk_v2_litt/` to the new location.
- Update the `NODE_LITT_DB_STORAGE_PATHS` environment variable to point to the new location.
- Restart the validator node.