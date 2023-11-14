package apiserver

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/urfave/cli"
)

const (
	RegisteredQuorumFlagName        = "auth.registered-quorum"
	TotalUnauthThroughputFlagName   = "auth.total-unauth-throughput"
	PerUserUnauthThroughputFlagName = "auth.per-user-unauth-throughput"
	ClientIPHeaderFlagName          = "auth.client-ip-header"
	TrustedProxiesFlagname          = "auth.trusted-proxies"
)

type QuorumRateInfo struct {
	PerUserUnauthThroughput common.RateParam
	TotalUnauthThroughput   common.RateParam
}

type RateConfig struct {
	QuorumRateInfos map[core.QuorumID]QuorumRateInfo
	ClientIPHeader  string
	TruxtedProxies  map[string]struct{}
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
		cli.StringFlag{
			Name:     ClientIPHeaderFlagName,
			Usage:    "The name of the header used to get the client IP address. If set to empty string, the IP address will be taken from the connection. The rightmost value of the header will be used, unless a list of trusted proxies is provided. For AWS, this should be set to 'x-forwarded-for'.",
			Required: false,
			Value:    "",
			EnvVar:   common.PrefixEnvVar(envPrefix, "CLIENT_IP_HEADER"),
		},
		cli.StringSliceFlag{
			Name:     TrustedProxiesFlagname,
			Usage:    "The trusted proxies that the request has been forwarded through. If non-empty, the disperser will pull the request IP from the first non-trusted proxy address in the header, reading from right to left.",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "TRUSTED_PROXIES"),
		},
	}
}

func ReadCLIConfig(c *cli.Context) RateConfig {

	trustedProxies := make(map[string]struct{})
	for _, proxy := range c.StringSlice(TrustedProxiesFlagname) {
		trustedProxies[proxy] = struct{}{}
	}

	quorumRateInfos := make(map[core.QuorumID]QuorumRateInfo)
	for ind, quorumID := range c.IntSlice(RegisteredQuorumFlagName) {

		quorumRateInfos[core.QuorumID(quorumID)] = QuorumRateInfo{
			TotalUnauthThroughput:   common.RateParam(c.IntSlice(TotalUnauthThroughputFlagName)[ind]),
			PerUserUnauthThroughput: common.RateParam(c.IntSlice(PerUserUnauthThroughputFlagName)[ind]),
		}
	}

	return RateConfig{
		QuorumRateInfos: quorumRateInfos,
		ClientIPHeader:  c.String(ClientIPHeaderFlagName),
		TruxtedProxies:  trustedProxies,
	}
}
