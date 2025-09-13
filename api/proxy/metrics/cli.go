package metrics

import (
	"errors"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"

	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	"github.com/olekukonko/tablewriter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli/v2"
)

const (
	DeprecatedEnabledFlagName = "metrics.enabled"

	ListenAddrFlagName = "metrics.addr"
	PortFlagName       = "metrics.port"
	defaultListenAddr  = "0.0.0.0"
	defaultListenPort  = 7300

	EnvPrefix = "metrics"
)

var ErrInvalidPort = errors.New("invalid metrics port")

func withEnvPrefix(envPrefix, s string) []string {
	return []string{envPrefix + "_METRICS_" + s}
}

func DefaultConfig() Config {
	return Config{
		Host: defaultListenAddr,
		Port: defaultListenPort,
	}
}

func DeprecatedCLIFlags(envPrefix string, category string) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name: DeprecatedEnabledFlagName,

			Usage:    "Enable the metrics server. On by default, so use --metrics.enabled=false to disable.",
			Category: category,
			Value:    true,
			EnvVars:  withEnvPrefix(envPrefix, "ENABLED"),
			Action: func(*cli.Context, bool) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --apis.enabled with `metrics` to turn on instead",
					DeprecatedEnabledFlagName, withEnvPrefix(envPrefix, "ENABLED"))
			},
			Hidden: true,
		}}
}

func CLIFlags(envPrefix string, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     ListenAddrFlagName,
			Usage:    "Metrics listening address",
			Category: category,
			Value:    defaultListenAddr,
			EnvVars:  withEnvPrefix(envPrefix, "ADDR"),
		},
		&cli.IntFlag{
			Name:     PortFlagName,
			Usage:    "Metrics listening port",
			Category: category,
			Value:    defaultListenPort,
			EnvVars:  withEnvPrefix(envPrefix, "PORT"),
		},
	}
}

func (m Config) Check() error {
	if m.Port < 0 || m.Port > math.MaxUint16 {
		return ErrInvalidPort
	}

	return nil
}

func ReadConfig(ctx *cli.Context) Config {
	return Config{
		Host: ctx.String(ListenAddrFlagName),
		Port: ctx.Int(PortFlagName),
	}
}

// NewSubcommands is used by `doc metrics` to output all supported metrics to
// stdout. For metrics to be included in the output they need to be created
// using the factory defined in `common/metrics.go`, and the metrics interface
// must have a `Document()` func. See interfaces and structs defined in
// `api/clients/v2/metrics` or `api/proxy/metrics/metrics.go` for usage.
func NewSubcommands() cli.Commands {
	return cli.Commands{
		{
			Name:  "metrics",
			Usage: "Dumps a list of supported metrics to stdout",
			Action: func(*cli.Context) error {
				registry := prometheus.NewRegistry()
				supportedMetrics := slices.Concat(
					NewMetrics(registry).Document(),
					metrics.NewAccountantMetrics(registry).Document(),
					metrics.NewDispersalMetrics(registry).Document(),
					metrics.NewRetrievalMetrics(registry).Document(),
				)

				table := tablewriter.NewWriter(os.Stdout)
				table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
				table.SetCenterSeparator("|")
				table.SetAutoWrapText(false)
				table.SetHeader([]string{"Metric", "Description", "Labels", "Type"})
				data := make([][]string, 0, len(supportedMetrics))
				for _, metric := range supportedMetrics {
					labels := strings.Join(metric.Labels, ",")
					data = append(data, []string{metric.Name, metric.Help, labels, metric.Type})
				}
				table.AppendBulk(data)
				table.Render()
				return nil
			},
		},
	}
}
