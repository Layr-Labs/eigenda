package flags

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

// runWithFlags runs a cli.App with the given flags and args, invoking fn inside the action.
func runWithFlags(t *testing.T, cliFlags []cli.Flag, args []string, fn func(ctx *cli.Context)) {
	t.Helper()
	app := cli.NewApp()
	app.Flags = cliFlags
	app.Action = func(ctx *cli.Context) error {
		fn(ctx)
		return nil
	}
	err := app.Run(args)
	assert.NoError(t, err)
}

func TestDeprecatedFlags_NoneSet(t *testing.T) {
	runWithFlags(t, deprecatedFlags, []string{"test"}, func(ctx *cli.Context) {
		set := getSetDeprecatedCLIFlags(ctx)
		assert.Empty(t, set)
	})
}

func TestDeprecatedFlags_StringFlagSetViaCLI(t *testing.T) {
	runWithFlags(t, deprecatedFlags, []string{"test", "--node.dispersal-port", "9000"}, func(ctx *cli.Context) {
		set := getSetDeprecatedCLIFlags(ctx)
		assert.Contains(t, set, "node.dispersal-port")
		assert.Len(t, set, 1)
	})
}

func TestDeprecatedFlags_BoolFlagSetViaCLI(t *testing.T) {
	runWithFlags(t, deprecatedFlags, []string{"test", "--node.disable-dispersal-authentication"}, func(ctx *cli.Context) {
		set := getSetDeprecatedCLIFlags(ctx)
		assert.Contains(t, set, "node.disable-dispersal-authentication")
		assert.Len(t, set, 1)
	})
}

func TestDeprecatedFlags_MultipleFlagsSet(t *testing.T) {
	args := []string{
		"test",
		"--node.dispersal-port", "9000",
		"--node.retrieval-port", "9001",
		"--node.runtime-mode", "v1-only",
		"--node.leveldb-disable-seeks-compaction-v1",
	}
	runWithFlags(t, deprecatedFlags, args, func(ctx *cli.Context) {
		set := getSetDeprecatedCLIFlags(ctx)
		assert.Len(t, set, 4)
		assert.Contains(t, set, "node.dispersal-port")
		assert.Contains(t, set, "node.retrieval-port")
		assert.Contains(t, set, "node.runtime-mode")
		assert.Contains(t, set, "node.leveldb-disable-seeks-compaction-v1")
	})
}

func TestDeprecatedFlags_SetViaEnvVar(t *testing.T) {
	t.Setenv("NODE_DISPERSAL_PORT", "9000")
	runWithFlags(t, deprecatedFlags, []string{"test"}, func(ctx *cli.Context) {
		set := getSetDeprecatedCLIFlags(ctx)
		assert.Contains(t, set, "node.dispersal-port")
		assert.Len(t, set, 1)
	})
}

func TestDeprecatedFlags_AllFlagsSetViaCLI(t *testing.T) {
	args := []string{
		"test",
		"--node.dispersal-port", "9000",
		"--node.retrieval-port", "9001",
		"--node.internal-dispersal-port", "9002",
		"--node.internal-retrieval-port", "9003",
		"--node.runtime-mode", "v1-and-v2",
		"--node.disable-dispersal-authentication",
		"--node.leveldb-disable-seeks-compaction-v1",
		"--node.leveldb-enable-sync-writes-v1",
		"--node.enable-payment-validation",
	}
	runWithFlags(t, deprecatedFlags, args, func(ctx *cli.Context) {
		set := getSetDeprecatedCLIFlags(ctx)
		assert.Len(t, set, len(deprecatedFlagNames))
		for _, name := range deprecatedFlagNames {
			assert.Contains(t, set, name)
		}
	})
}

func TestDeprecatedFlags_UsageText(t *testing.T) {
	for _, f := range deprecatedFlags {
		switch flag := f.(type) {
		case cli.StringFlag:
			assert.Equal(t, deprecatedUsage, flag.Usage, "flag %s should have deprecated usage", flag.Name)
		case cli.BoolFlag:
			assert.Equal(t, deprecatedUsage, flag.Usage, "flag %s should have deprecated usage", flag.Name)
		case cli.BoolTFlag:
			assert.Equal(t, deprecatedUsage, flag.Usage, "flag %s should have deprecated usage", flag.Name)
		default:
			t.Errorf("unexpected flag type for %v", f)
		}
	}
}

func TestDeprecatedFlags_IncludedInGlobalFlags(t *testing.T) {
	flagNames := make(map[string]bool)
	for _, f := range Flags {
		flagNames[f.GetName()] = true
	}
	for _, f := range deprecatedFlags {
		assert.True(t, flagNames[f.GetName()], "deprecated flag %s should be in global Flags", f.GetName())
	}
}

func TestDeprecatedFlags_DoNotBreakApp(t *testing.T) {
	// Verify that the app does not error when deprecated flags are passed alongside real flags.
	allFlags := append([]cli.Flag{}, Flags...)
	args := []string{
		"test",
		"--node.dispersal-port", "9000",
		"--node.runtime-mode", "v1-only",
		"--node.disable-dispersal-authentication",
	}

	// Set required flags via env so the app can parse without errors.
	requiredEnvs := map[string]string{
		"NODE_HOSTNAME":                   "localhost",
		"NODE_ENABLE_NODE_API":            "true",
		"NODE_ENABLE_METRICS":             "true",
		"NODE_TIMEOUT":                    "1s",
		"NODE_QUORUM_ID_LIST":             "0",
		"NODE_DB_PATH":                    "/tmp/test",
		"NODE_EIGENDA_DIRECTORY":          "0x0000000000000000000000000000000000000000",
		"NODE_CHURNER_URL":                "http://localhost:1234",
		"NODE_PUBLIC_IP_PROVIDER":         "ipify",
		"NODE_PUBLIC_IP_CHECK_INTERVAL":   "0s",
		"NODE_CHAIN_RPC":                  "http://localhost:8545",
		"NODE_PRIVATE_KEY":                "0x00",
		"NODE_G1_PATH":                    "/tmp/g1.point",
		"NODE_CACHE_PATH":                 "/tmp/srs",
		"NODE_SRS_ORDER":                  "1",
		"NODE_SRS_LOAD":                   "1",
		"NODE_V2_DISPERSAL_PORT":          "32005",
		"NODE_V2_RETRIEVAL_PORT":          "32004",
		"NODE_INTERNAL_V2_DISPERSAL_PORT": "32007",
		"NODE_INTERNAL_V2_RETRIEVAL_PORT": "32006",
	}
	for k, v := range requiredEnvs {
		t.Setenv(k, v)
	}
	// Clear any stale env vars for the deprecated flags being tested via CLI.
	os.Unsetenv("NODE_DISPERSAL_PORT")
	os.Unsetenv("NODE_RUNTIME_MODE")
	os.Unsetenv("NODE_DISABLE_DISPERSAL_AUTHENTICATION")

	app := cli.NewApp()
	app.Flags = allFlags
	var actionCalled bool
	app.Action = func(ctx *cli.Context) error {
		actionCalled = true
		set := getSetDeprecatedCLIFlags(ctx)
		assert.Len(t, set, 3)
		return nil
	}
	err := app.Run(args)
	assert.NoError(t, err)
	assert.True(t, actionCalled, "app action should have been called")
}
