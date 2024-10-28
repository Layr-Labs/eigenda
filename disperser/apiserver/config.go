package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/urfave/cli"
)

const (
	RegisteredQuorumFlagName         = "auth.registered-quorum"
	TotalUnauthThroughputFlagName    = "auth.total-unauth-byte-rate"
	PerUserUnauthThroughputFlagName  = "auth.per-user-unauth-byte-rate"
	TotalUnauthBlobRateFlagName      = "auth.total-unauth-blob-rate"
	PerUserUnauthBlobRateFlagName    = "auth.per-user-unauth-blob-rate"
	ClientIPHeaderFlagName           = "auth.client-ip-header"
	AllowlistFileFlagName            = "auth.allowlist-file"
	AllowlistRefreshIntervalFlagName = "auth.allowlist-refresh-interval"

	RetrievalBlobRateFlagName   = "auth.retrieval-blob-rate"
	RetrievalThroughputFlagName = "auth.retrieval-throughput"

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
	Name       string
	Throughput common.RateParam
	BlobRate   common.RateParam
}

type Allowlist = map[string]map[core.QuorumID]PerUserRateInfo

type AllowlistEntry struct {
	Name     string  `json:"name"`
	Account  string  `json:"account"`
	QuorumID uint8   `json:"quorumID"`
	BlobRate float64 `json:"blobRate"`
	ByteRate float64 `json:"byteRate"`
}

type RateConfig struct {
	QuorumRateInfos map[core.QuorumID]QuorumRateInfo
	ClientIPHeader  string
	Allowlist       Allowlist

	RetrievalBlobRate   common.RateParam
	RetrievalThroughput common.RateParam

	AllowlistFile            string
	AllowlistRefreshInterval time.Duration
}

func AllowlistFileFlag(envPrefix string) cli.Flag {
	return cli.StringFlag{
		Name:     AllowlistFileFlagName,
		Usage:    "Path to a file containing the allowlist of IPs or ethereum addresses (including initial \"0x\") and corresponding blob/byte rates to bypass rate limiting. This file must be in JSON format",
		EnvVar:   common.PrefixEnvVar(envPrefix, "ALLOWLIST_FILE"),
		Required: false,
	}
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
		AllowlistFileFlag(envPrefix),
		cli.DurationFlag{
			Name:     AllowlistRefreshIntervalFlagName,
			Usage:    "The interval at which to refresh the allowlist from the file",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "ALLOWLIST_REFRESH_INTERVAL"),
			Value:    5 * time.Minute,
		},
		cli.IntFlag{
			Name:     RetrievalBlobRateFlagName,
			Usage:    "The blob rate limit for retrieval requests (Blobs/sec)",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "RETRIEVAL_BLOB_RATE"),
		},
		cli.IntFlag{
			Name:     RetrievalThroughputFlagName,
			Usage:    "The throughput rate limit for retrieval requests (Bytes/sec)",
			EnvVar:   common.PrefixEnvVar(envPrefix, "RETRIEVAL_BYTE_RATE"),
			Required: true,
		},
	}
}

func ReadAllowlistFromFile(f string) (Allowlist, error) {
	allowlist := make(Allowlist)
	if f == "" {
		return allowlist, nil
	}

	allowlistFile, err := os.Open(f)
	if err != nil {
		log.Printf("failed to read allowlist file: %s", err)
		return allowlist, err
	}
	defer allowlistFile.Close()
	var allowlistEntries []AllowlistEntry
	content, err := io.ReadAll(allowlistFile)
	if err != nil {
		log.Printf("failed to load allowlist file content: %s", err)
		return allowlist, err
	}
	err = json.Unmarshal(content, &allowlistEntries)
	if err != nil {
		log.Printf("failed to parse allowlist file content: %s", err)
		return allowlist, err
	}

	for _, entry := range allowlistEntries {
		// normalize to lowercase (non-checksummed) address or IP address
		account := strings.ToLower(entry.Account)
		rateInfoByQuorum, ok := allowlist[account]
		if !ok {
			allowlist[account] = map[core.QuorumID]PerUserRateInfo{
				core.QuorumID(entry.QuorumID): {
					Name:       entry.Name,
					Throughput: common.RateParam(entry.ByteRate),
					BlobRate:   common.RateParam(entry.BlobRate * blobRateMultiplier),
				},
			}
		} else {
			rateInfoByQuorum[core.QuorumID(entry.QuorumID)] = PerUserRateInfo{
				Name:       entry.Name,
				Throughput: common.RateParam(entry.ByteRate),
				BlobRate:   common.RateParam(entry.BlobRate * blobRateMultiplier),
			}
		}
	}

	return allowlist, nil
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

	allowlist := make(Allowlist)
	allowlistFileName := c.String(AllowlistFileFlagName)
	if allowlistFileName != "" {
		var err error
		allowlist, err = ReadAllowlistFromFile(allowlistFileName)
		if err != nil {
			return RateConfig{}, fmt.Errorf("failed to read allowlist file %s: %w", allowlistFileName, err)
		}
	}

	return RateConfig{
		QuorumRateInfos:          quorumRateInfos,
		ClientIPHeader:           c.String(ClientIPHeaderFlagName),
		Allowlist:                allowlist,
		RetrievalBlobRate:        common.RateParam(c.Int(RetrievalBlobRateFlagName) * blobRateMultiplier),
		RetrievalThroughput:      common.RateParam(c.Int(RetrievalThroughputFlagName)),
		AllowlistFile:            c.String(AllowlistFileFlagName),
		AllowlistRefreshInterval: c.Duration(AllowlistRefreshIntervalFlagName),
	}, nil
}
