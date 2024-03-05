package apiserver

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/urfave/cli"
)

const (
	RegisteredQuorumFlagName        = "auth.registered-quorum"
	TotalUnauthThroughputFlagName   = "auth.total-unauth-byte-rate"
	PerUserUnauthThroughputFlagName = "auth.per-user-unauth-byte-rate"
	TotalUnauthBlobRateFlagName     = "auth.total-unauth-blob-rate"
	PerUserUnauthBlobRateFlagName   = "auth.per-user-unauth-blob-rate"
	ClientIPHeaderFlagName          = "auth.client-ip-header"
	AllowlistFlagName               = "auth.allowlist"

	// We allow the user to specify the blob rate in blobs/sec, but internally we use blobs/sec * 1e6 (i.e. blobs/microsec).
	// This is because the rate limiter takes an integer rate.
	blobRateMultiplier = 1e6
)

type QuorumRateInfo struct {
	PerUserUnauthThroughput common.RateParam
	TotalUnauthThroughput   common.RateParam
	PerUserUnauthBlobRate   common.RateParam
	TotalUnauthBlobRate     common.RateParam
}

type PerUserRateInfo struct {
	Throughput common.RateParam
	BlobRate   common.RateParam
}

type Allowlist = map[string]map[core.QuorumID]PerUserRateInfo

type RateConfig struct {
	QuorumRateInfos map[core.QuorumID]QuorumRateInfo
	ClientIPHeader  string
	Allowlist       Allowlist
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
			EnvVar:   common.PrefixEnvVar(envPrefix, "TOTAL_UNAUTH_BYTE_RATE"),
		},
		cli.IntSliceFlag{
			Name:     PerUserUnauthThroughputFlagName,
			Usage:    "Per-user encoded throughput for unauthenticated requests (Bytes/sec)",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "PER_USER_UNAUTH_BYTE_RATE"),
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
		cli.StringSliceFlag{
			Name:     AllowlistFlagName,
			Usage:    "Allowlist of IPs or ethereum addresses (including initial \"0x\") and corresponding blob/byte rates to bypass rate limiting. Format: [<IP>||<ETH ADDRESS>]/<quorum ID>/<blob rate>/<byte rate>. Example: 127.0.0.1/0/10/10485760",
			EnvVar:   common.PrefixEnvVar(envPrefix, "ALLOWLIST"),
			Required: false,
			Value:    &cli.StringSlice{},
		},
	}
}

func ReadCLIConfig(c *cli.Context) (RateConfig, error) {

	numQuorums := len(c.IntSlice(RegisteredQuorumFlagName))
	if len(c.StringSlice(TotalUnauthBlobRateFlagName)) != numQuorums {
		return RateConfig{}, errors.New("number of total unauth blob rates does not match number of quorums")
	}
	if len(c.StringSlice(PerUserUnauthBlobRateFlagName)) != numQuorums {
		return RateConfig{}, errors.New("number of per user unauth blob intervals does not match number of quorums")
	}
	if len(c.IntSlice(TotalUnauthThroughputFlagName)) != numQuorums {
		return RateConfig{}, errors.New("number of total unauth throughput does not match number of quorums")
	}
	if len(c.IntSlice(PerUserUnauthThroughputFlagName)) != numQuorums {
		return RateConfig{}, errors.New("number of per user unauth throughput does not match number of quorums")
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

	// Parse allowlist
	allowlist := make(Allowlist)
	for _, allowlistEntry := range c.StringSlice(AllowlistFlagName) {
		allowlistEntrySplit := strings.Split(allowlistEntry, "/")
		if len(allowlistEntrySplit) != 4 {
			log.Printf("invalid allowlist entry: entry should contain exactly 4 elements: %s", allowlistEntry)
			continue
		}
		ip := allowlistEntrySplit[0]
		quorumID, err := strconv.Atoi(allowlistEntrySplit[1])
		if err != nil {
			log.Printf("invalid allowlist entry: failed to convert quorum ID from string: %s", allowlistEntry)
			continue
		}
		blobRate, err := strconv.ParseFloat(allowlistEntrySplit[2], 64)
		if err != nil {
			log.Printf("invalid allowlist entry: failed to convert blob rate from string: %s", allowlistEntry)
			continue
		}
		byteRate, err := strconv.ParseFloat(allowlistEntrySplit[3], 64)
		if err != nil {
			log.Printf("invalid allowlist entry: failed to convert throughput from string: %s", allowlistEntry)
			continue
		}
		rateInfoByQuorum, ok := allowlist[ip]
		if !ok {
			allowlist[ip] = map[core.QuorumID]PerUserRateInfo{
				core.QuorumID(quorumID): {
					Throughput: common.RateParam(byteRate),
					BlobRate:   common.RateParam(blobRate * blobRateMultiplier),
				},
			}
		} else {
			rateInfoByQuorum[core.QuorumID(quorumID)] = PerUserRateInfo{
				Throughput: common.RateParam(byteRate),
				BlobRate:   common.RateParam(blobRate * blobRateMultiplier),
			}
		}
	}

	return RateConfig{
		QuorumRateInfos: quorumRateInfos,
		ClientIPHeader:  c.String(ClientIPHeaderFlagName),
		Allowlist:       allowlist,
	}, nil
}
