package flags

import (
	"log"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const deprecatedUsage = "Deprecated v1 flag. This flag will be ignored"

// Deprecated v1 flags. These flags are no longer functional but are kept
// to avoid breaking users who haven't yet removed them from their configurations.
var (
	DeprecatedDispersalPortFlag = cli.StringFlag{
		Name:   common.PrefixFlag(FlagPrefix, "dispersal-port"),
		Usage:  deprecatedUsage,
		EnvVar: common.PrefixEnvVar(EnvVarPrefix, "DISPERSAL_PORT"),
	}
	DeprecatedRetrievalPortFlag = cli.StringFlag{
		Name:   common.PrefixFlag(FlagPrefix, "retrieval-port"),
		Usage:  deprecatedUsage,
		EnvVar: common.PrefixEnvVar(EnvVarPrefix, "RETRIEVAL_PORT"),
	}
	DeprecatedInternalDispersalPortFlag = cli.StringFlag{
		Name:   common.PrefixFlag(FlagPrefix, "internal-dispersal-port"),
		Usage:  deprecatedUsage,
		EnvVar: common.PrefixEnvVar(EnvVarPrefix, "INTERNAL_DISPERSAL_PORT"),
	}
	DeprecatedInternalRetrievalPortFlag = cli.StringFlag{
		Name:   common.PrefixFlag(FlagPrefix, "internal-retrieval-port"),
		Usage:  deprecatedUsage,
		EnvVar: common.PrefixEnvVar(EnvVarPrefix, "INTERNAL_RETRIEVAL_PORT"),
	}
	DeprecatedRuntimeModeFlag = cli.StringFlag{
		Name:   common.PrefixFlag(FlagPrefix, "runtime-mode"),
		Usage:  deprecatedUsage,
		EnvVar: common.PrefixEnvVar(EnvVarPrefix, "RUNTIME_MODE"),
	}
	DeprecatedDisableDispersalAuthenticationFlag = cli.BoolFlag{
		Name:   common.PrefixFlag(FlagPrefix, "disable-dispersal-authentication"),
		Usage:  deprecatedUsage,
		EnvVar: common.PrefixEnvVar(EnvVarPrefix, "DISABLE_DISPERSAL_AUTHENTICATION"),
	}
	DeprecatedLevelDBDisableSeeksCompactionV1Flag = cli.BoolFlag{
		Name:   common.PrefixFlag(FlagPrefix, "leveldb-disable-seeks-compaction-v1"),
		Usage:  deprecatedUsage,
		EnvVar: common.PrefixEnvVar(EnvVarPrefix, "LEVELDB_DISABLE_SEEKS_COMPACTION_V1"),
	}
	DeprecatedLevelDBEnableSyncWritesV1Flag = cli.BoolFlag{
		Name:   common.PrefixFlag(FlagPrefix, "leveldb-enable-sync-writes-v1"),
		Usage:  deprecatedUsage,
		EnvVar: common.PrefixEnvVar(EnvVarPrefix, "LEVELDB_ENABLE_SYNC_WRITES_V1"),
	}
	DeprecatedEnablePaymentValidationFlag = cli.BoolFlag{
		Name:   common.PrefixFlag(FlagPrefix, "enable-payment-validation"),
		Usage:  deprecatedUsage,
		EnvVar: common.PrefixEnvVar(EnvVarPrefix, "ENABLE_PAYMENT_VALIDATION"),
	}
)

var deprecatedFlags = []cli.Flag{
	DeprecatedDispersalPortFlag,
	DeprecatedRetrievalPortFlag,
	DeprecatedInternalDispersalPortFlag,
	DeprecatedInternalRetrievalPortFlag,
	DeprecatedRuntimeModeFlag,
	DeprecatedDisableDispersalAuthenticationFlag,
	DeprecatedLevelDBDisableSeeksCompactionV1Flag,
	DeprecatedLevelDBEnableSyncWritesV1Flag,
	DeprecatedEnablePaymentValidationFlag,
}

// deprecatedFlagNames contains the CLI names of all deprecated flags for use in CheckDeprecatedCLIFlags.
var deprecatedFlagNames = []string{
	DeprecatedDispersalPortFlag.Name,
	DeprecatedRetrievalPortFlag.Name,
	DeprecatedInternalDispersalPortFlag.Name,
	DeprecatedInternalRetrievalPortFlag.Name,
	DeprecatedRuntimeModeFlag.Name,
	DeprecatedDisableDispersalAuthenticationFlag.Name,
	DeprecatedLevelDBDisableSeeksCompactionV1Flag.Name,
	DeprecatedLevelDBEnableSyncWritesV1Flag.Name,
	DeprecatedEnablePaymentValidationFlag.Name,
}

// CheckDeprecatedCLIFlags logs a warning for each deprecated flag that has been set.
func CheckDeprecatedCLIFlags(ctx *cli.Context) {
	for _, name := range getSetDeprecatedCLIFlags(ctx) {
		log.Printf("WARNING: Flag --%s is deprecated and will be ignored. "+
			"Please remove it from your configuration.", name)
	}
}

// getSetDeprecatedCLIFlags returns the names of deprecated flags that have been explicitly set.
func getSetDeprecatedCLIFlags(ctx *cli.Context) []string {
	var set []string
	for _, name := range deprecatedFlagNames {
		if ctx.GlobalIsSet(name) {
			set = append(set, name)
		}
	}
	return set
}
