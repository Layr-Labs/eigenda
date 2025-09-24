package downloader

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

const (
	flagDimension       = "dimension"
	flagTablesOutputDir = "output-dir"
	flagTablesBaseURL   = "base-url"
	flagCosetSizes      = "coset-sizes"
)

// TablesFlags defines command line flags for the download-tables command
var TablesFlags = []cli.Flag{
	cli.StringFlag{
		Name:  flagDimension,
		Usage: "Dimension name (e.g., dimE8192)",
		Value: defaultDimension,
	},
	cli.StringFlag{
		Name:  flagTablesOutputDir,
		Usage: "Output directory for downloaded SRS table files",
		Value: defaultTablesOutputDir,
	},
	cli.StringFlag{
		Name:  flagTablesBaseURL,
		Usage: "Base URL for downloading SRS table files",
		Value: defaultTablesBaseURL,
	},
	cli.StringFlag{
		Name:  flagCosetSizes,
		Usage: "Comma-separated list of coset sizes to download (e.g., 4,8,16,32,64,128,256,512,1024)",
		Value: "4,8,16,32,64,128,256,512,1024",
	},
}

// ReadTablesConfig reads command line flags into a config struct
func ReadTablesConfig(cCtx *cli.Context) (TablesDownloaderConfig, error) {
	cosetSizesStr := cCtx.String(flagCosetSizes)
	var cosetSizes []int
	if cosetSizesStr != "" {
		parts := strings.Split(cosetSizesStr, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			var size int
			if _, err := fmt.Sscanf(part, "%d", &size); err == nil {
				cosetSizes = append(cosetSizes, size)
			}
		}
	}

	return NewTablesDownloaderConfig(
		cCtx.String(flagDimension),
		cCtx.String(flagTablesOutputDir),
		cCtx.String(flagTablesBaseURL),
		cosetSizes,
	)
}
