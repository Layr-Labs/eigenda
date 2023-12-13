package apiserver

import (
	"fmt"
	"strconv"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/urfave/cli"
)

const (
	RegisteredQuorumFlagName        = "auth.registered-quorum"
	TotalUnauthThroughputFlagName   = "auth.total-unauth-throughput"
	PerUserUnauthThroughputFlagName = "auth.per-user-unauth-throughput"
	TotalUnauthBlobRateFlagName     = "auth.total-unauth-blob-rate"
	PerUserUnauthBlobRateFlagName   = "auth.per-user-unauth-blob-rate"
	ClientIPHeaderFlagName          = "auth.client-ip-header"

	blobRateMultiplier = 1e6
)

type QuorumRateInfo struct {
	PerUserUnauthThroughput common.RateParam
	TotalUnauthThroughput   common.RateParam
	PerUserUnauthBlobRate   common.RateParam
	TotalUnauthBlobRate     common.RateParam
}

type RateConfig struct {
	QuorumRateInfos map[core.QuorumID]QuorumRateInfo
	ClientIPHeader  string
}

func CLIFlags(envPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.IntSliceFlag{
			Name:     RegisteredQuorumFlagName,
			Usage:    "The quorum ID for the quorum",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "REGISTERED_QUORUM_ID"),
		},
		cli.IntSliceFlag{
			Name:     TotalUnauthThroughputFlagName,
			Usage:    "Total encoded throughput for unauthenticated requests (Bytes/sec)",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "TOTAL_UNAUTH_THROUGHPUT"),
		},
		cli.IntSliceFlag{
			Name:     PerUserUnauthThroughputFlagName,
			Usage:    "Per-user encoded throughput for unauthenticated requests (Bytes/sec)",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "PER_USER_UNAUTH_THROUGHPUT"),
		},
		cli.StringSliceFlag{
			Name:     TotalUnauthBlobRateFlagName,
			Usage:    "Total blob rate for unauthenticated requests (Blobs/sec)",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "TOTAL_UNAUTH_BLOB_RATE"),
		},
		cli.StringSliceFlag{
			Name:     PerUserUnauthBlobRateFlagName,
			Usage:    "Per-user blob interval for unauthenticated requests (Blobs/sec)",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "PER_USER_UNAUTH_BLOB_RATE"),
		},
		cli.StringFlag{
			Name:     ClientIPHeaderFlagName,
			Usage:    "The name of the header used to get the client IP address. If set to empty string, the IP address will be taken from the connection. The rightmost value of the header will be used. For AWS, this should be set to 'x-forwarded-for'.",
			Required: false,
			Value:    "",
			EnvVar:   common.PrefixEnvVar(envPrefix, "CLIENT_IP_HEADER"),
		},
	}
}

func ReadCLIConfig(c *cli.Context) (RateConfig, error) {

	numQuorums := len(c.IntSlice(RegisteredQuorumFlagName))
	if len(c.StringSlice(TotalUnauthBlobRateFlagName)) != numQuorums {
		return RateConfig{}, fmt.Errorf("number of total unauth blob rates does not match number of quorums")
	}
	if len(c.StringSlice(PerUserUnauthBlobRateFlagName)) != numQuorums {
		return RateConfig{}, fmt.Errorf("number of per user unauth blob intervals does not match number of quorums")
	}
	if len(c.IntSlice(TotalUnauthThroughputFlagName)) != numQuorums {
		return RateConfig{}, fmt.Errorf("number of total unauth throughput does not match number of quorums")
	}
	if len(c.IntSlice(PerUserUnauthThroughputFlagName)) != numQuorums {
		return RateConfig{}, fmt.Errorf("number of per user unauth throughput does not match number of quorums")
	}

	quorumRateInfos := make(map[core.QuorumID]QuorumRateInfo)
	for ind, quorumID := range c.IntSlice(RegisteredQuorumFlagName) {

		totalBlobRate, err := strconv.ParseFloat(c.StringSlice(TotalUnauthBlobRateFlagName)[ind], 64)
		if err != nil {
			return RateConfig{}, err
		}
		accountBlobRate, err := strconv.ParseFloat(c.StringSlice(PerUserUnauthBlobRateFlagName)[ind], 64)
		if err != nil {
			return RateConfig{}, err
		}

		quorumRateInfos[core.QuorumID(quorumID)] = QuorumRateInfo{
			TotalUnauthThroughput:   common.RateParam(c.IntSlice(TotalUnauthThroughputFlagName)[ind]),
			PerUserUnauthThroughput: common.RateParam(c.IntSlice(PerUserUnauthThroughputFlagName)[ind]),
			TotalUnauthBlobRate:     common.RateParam(totalBlobRate * blobRateMultiplier),
			PerUserUnauthBlobRate:   common.RateParam(accountBlobRate * blobRateMultiplier),
		}
	}

	return RateConfig{
		QuorumRateInfos: quorumRateInfos,
		ClientIPHeader:  c.String(ClientIPHeaderFlagName),
	}, nil
}
