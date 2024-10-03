package ratelimit

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const (
	BucketSizesFlagName       = "bucket-sizes"
	BucketMultipliersFlagName = "bucket-multipliers"
	CountFailedFlagName       = "count-failed"
	BucketStoreSizeFlagName   = "bucket-store-size"
)

type Config struct {
	common.GlobalRateParams
	BucketStoreSize  int
	UniformRateParam common.RateParam
}

func RatelimiterCLIFlags(envPrefix string, flagPrefix string) []cli.Flag {
	bucketSizes := cli.StringSlice([]string{"1s"})
	bucketMultipliers := cli.StringSlice([]string{"1"})

	return []cli.Flag{
		cli.StringSliceFlag{
			Name:   common.PrefixFlag(flagPrefix, BucketSizesFlagName),
			Usage:  "Bucket sizes (duration)",
			Value:  &bucketSizes,
			EnvVar: common.PrefixEnvVar(envPrefix, "BUCKET_SIZES"),
		},
		cli.StringSliceFlag{
			Name:   common.PrefixFlag(flagPrefix, BucketMultipliersFlagName),
			Usage:  "Bucket multipiers (float)",
			Value:  &bucketMultipliers,
			EnvVar: common.PrefixEnvVar(envPrefix, "BUCKET_MULTIPLIERS"),
		},
		cli.BoolFlag{
			Name:   common.PrefixFlag(flagPrefix, CountFailedFlagName),
			Usage:  "Count failed requests",
			EnvVar: common.PrefixEnvVar(envPrefix, "COUNT_FAILED"),
		},
		cli.IntFlag{
			Name:     common.PrefixFlag(flagPrefix, BucketStoreSizeFlagName),
			Usage:    "Bucket store size",
			Value:    1000,
			EnvVar:   common.PrefixEnvVar(envPrefix, "BUCKET_STORE_SIZE"),
			Required: false,
		},
	}
}

func DefaultCLIConfig() Config {
	return Config{}
}

func validateConfig(cfg Config) error {
	if len(cfg.BucketSizes) != len(cfg.Multipliers) {
		return errors.New("number of bucket sizes does not match number of multipliers")
	}
	for _, mult := range cfg.Multipliers {
		if mult <= 0 {
			return errors.New("multiplier must be positive")
		}
	}
	return nil
}

func ReadCLIConfig(ctx *cli.Context, flagPrefix string) (Config, error) {
	cfg := DefaultCLIConfig()

	strings := ctx.StringSlice(common.PrefixFlag(flagPrefix, BucketSizesFlagName))
	sizes := make([]time.Duration, len(strings))
	for i, s := range strings {
		d, err := time.ParseDuration(s)
		if err != nil {
			return Config{}, fmt.Errorf("bucket size failed to parse: %v", err)
		}
		sizes[i] = d
	}
	cfg.BucketSizes = sizes

	strings = ctx.StringSlice(common.PrefixFlag(flagPrefix, BucketMultipliersFlagName))
	multipliers := make([]float32, len(strings))
	for i, s := range strings {
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return Config{}, fmt.Errorf("bucket multiplier failed to parse: %v", err)
		}
		multipliers[i] = float32(f)
	}
	cfg.Multipliers = multipliers
	cfg.GlobalRateParams.CountFailed = ctx.Bool(common.PrefixFlag(flagPrefix, CountFailedFlagName))
	cfg.BucketStoreSize = ctx.Int(common.PrefixFlag(flagPrefix, BucketStoreSizeFlagName))

	err := validateConfig(cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
