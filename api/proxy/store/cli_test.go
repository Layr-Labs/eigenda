package store

import (
	"testing"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func runWithStorageFlags(t *testing.T, args []string, fn func(ctx *cli.Context)) {
	t.Helper()

	app := cli.NewApp()
	app.Flags = CLIFlags("EIGENDA_PROXY", "Storage")
	app.Action = func(ctx *cli.Context) error {
		fn(ctx)
		return nil
	}

	err := app.Run(args)
	require.NoError(t, err)
}

func TestReadConfigDefaultsToV2Backend(t *testing.T) {
	runWithStorageFlags(t, []string{"test"}, func(ctx *cli.Context) {
		cfg, err := ReadConfig(ctx)

		require.NoError(t, err)
		require.Equal(t, []common.EigenDABackend{common.V2EigenDABackend}, cfg.BackendsToEnable)
		require.Equal(t, common.V2EigenDABackend, cfg.DispersalBackend)
	})
}

func TestReadConfigFiltersEmptySecondaryTargets(t *testing.T) {
	t.Setenv("EIGENDA_PROXY_STORAGE_CACHE_TARGETS", "")
	t.Setenv("EIGENDA_PROXY_STORAGE_FALLBACK_TARGETS", "")

	runWithStorageFlags(t, []string{"test"}, func(ctx *cli.Context) {
		cfg, err := ReadConfig(ctx)

		require.NoError(t, err)
		require.Empty(t, cfg.CacheTargets)
		require.Empty(t, cfg.FallbackTargets)
	})
}

func TestReadConfigRejectsEmptyBackends(t *testing.T) {
	runWithStorageFlags(t, []string{"test", "--storage.backends-to-enable", ""}, func(ctx *cli.Context) {
		_, err := ReadConfig(ctx)

		require.ErrorContains(t, err, "backends must not be empty")
	})
}
