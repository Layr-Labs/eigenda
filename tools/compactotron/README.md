# Compactotron

A tool for compacting LevelDB databases. Does not modify the original database, but creates a new database in a new
location that contains only the data that is reachable.

## Build

Clone the EigenDA repo:

```bash
git clone https://github.com/Layr-Labs/eigenda.git
```

You will need to install Go 1.24 or later to build this tool. The instructions for doing this are OS specific, but
easily available on google. (Hint: modern LLMs are surprisingly adept at finding the right instructions for your OS.)

Once you have Go installed, you can build the tool by running:

```bash
cd eigenda/tools/compactotron
make build
```

A binary will be created at `eigenda/tools/compactotron/bin/compactotron`.

## Usage

```bash
eigenda/tools/compactotron/bin/compactotron <source_path> <destination_path>
```

**Arguments:**
- `source_path`: Path to the existing LevelDB database to compact. If this is for a validator, this path should be
                 `$NODE_DB_PATH/chunks`. This path will not be modified by this tool.
- `destination_path`: Path where the compacted database will be written.

Once this tool completes successfully and terminates, you can replace the original database with the compacted one.

IMPORTANT: if you are using this tool on a validator, the validator MUST be stopped before running this tool. Data
corruption is likely if you do not stop the validator first.

## In Case of Failure

Do not attempt to use the compacted database if this utility throws any errors during execution. Delete whatever
files created during the failed run and try again.

This tool does not modify the original database, so it is always safe to go back to the original database if this tool
has problems or takes too long to complete.