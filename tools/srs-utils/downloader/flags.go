package downloader

import (
	"github.com/urfave/cli"
)

const (
	flagBlobSize  = "blob-size-bytes"
	flagOutputDir = "output-dir"
	flagBaseURL   = "base-url"
)

// Flags defines command line flags for the download command
var Flags = []cli.Flag{
	cli.Uint64Flag{
		Name:  flagBlobSize,
		Usage: "Size of the blob in bytes",
		Value: 16777216, // Default to 16MB (16 * 1024 * 1024)
	},
	cli.StringFlag{
		Name:  flagOutputDir,
		Usage: "Output directory for downloaded files",
		Value: defaultOutputDir,
	},
	cli.StringFlag{
		Name:  flagBaseURL,
		Usage: "Base URL for downloading SRS files",
		Value: defaultBaseURL,
	},
}

// ReadCLIConfig reads command line flags into a config struct
func ReadCLIConfig(cCtx *cli.Context) (DownloaderConfig, error) {
	return NewDownloaderConfig(
		cCtx.Uint64(flagBlobSize),
		cCtx.String(flagOutputDir),
		cCtx.String(flagBaseURL),
	)
}
