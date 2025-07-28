package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"

	proxycmn "github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/core/eth"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
)

const (
	ethRpcUrlFlagName        = "eth-rpc-url"
	NetworkFlagName          = "network"
	directoryAddressFlagName = "directory-address"
	outputFlagName           = "output"
)

func main() {
	app := cli.NewApp()
	if bi, ok := debug.ReadBuildInfo(); ok {
		app.Version = bi.Main.Version
	}
	app.Name = "eigenda-directory"
	app.Usage = "EigenDA Directory Contract Address Discovery Tool"
	app.Description = "Tool for fetching all contract addresses from the EigenDADirectory contract on a specified EigenDA network."

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     ethRpcUrlFlagName,
			Usage:    "Ethereum RPC URL",
			EnvVars:  []string{"ETH_RPC_URL"},
			Required: true,
		},
		&cli.StringFlag{
			Name: NetworkFlagName,
			Usage: fmt.Sprintf(`The EigenDA network to discover (%s, %s, %s, %s).
Must match the chain-id of the ethereum rpc url provided. Used to select the hardcoded default EigenDADirectory address.
That address can be overridden by providing the --%s flag.`,
				proxycmn.MainnetEigenDANetwork,
				proxycmn.HoleskyTestnetEigenDANetwork,
				proxycmn.HoleskyPreprodEigenDANetwork,
				proxycmn.SepoliaTestnetEigenDANetwork,
				directoryAddressFlagName,
			),
			EnvVars: []string{"NETWORK"},
			Action: func(ctx *cli.Context, v string) error {
				if v == "" {
					// if no network is provided, we will try to auto-detect it from the chain ID
					return nil
				}
				// try to parse the network from the string.
				// this will validate the network and return an error if it's invalid.
				_, err := proxycmn.EigenDANetworkFromString(v)
				return err
			},
		},
		&cli.StringFlag{
			Name:    directoryAddressFlagName,
			Usage:   "EigenDADirectory contract address (overrides the default network address)",
			EnvVars: []string{"EIGENDA_DIRECTORY_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    outputFlagName,
			Usage:   "Output format: table, csv, json (default: table)",
			Value:   "table",
			EnvVars: []string{"OUTPUT_FORMAT"},
			Action: func(ctx *cli.Context, v string) error {
				validFormats := map[string]bool{"table": true, "csv": true, "json": true}
				if _, ok := validFormats[strings.ToLower(v)]; !ok {
					return fmt.Errorf("invalid output format: %s. Must be one of: table, csv, json", v)
				}
				return nil
			},
		},
	}

	app.Action = discoverAddresses
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("application failed: %v", err)
	}
}

func discoverAddresses(ctx *cli.Context) error {
	outputFormat := strings.ToLower(ctx.String(outputFlagName))
	rpcURL := ctx.String(ethRpcUrlFlagName)

	// Simple logging
	logger := log.New(os.Stderr, "[discovery] ", log.LstdFlags)

	// Connect to Ethereum
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("dial Ethereum node at %s: %w", rpcURL, err)
	}
	logger.Printf("Connected to Ethereum node at %s", rpcURL)

	directoryAddr := ctx.String(directoryAddressFlagName)
	networkName := ctx.String(NetworkFlagName)
	// cases:
	// 1. directory address: use it directly
	// 2. no directory address, but network provided: use the default directory address for the network
	// 3. no directory address and no network: use the default network for the chainID
	if directoryAddr == "" {
		logger.Printf("No directory address provided, attempting to auto-detect.")
		if networkName != "" {
			logger.Printf("network %s provided, attempting to use default directory address.", networkName)
			network, err := proxycmn.EigenDANetworkFromString(ctx.String(NetworkFlagName))
			if err != nil {
				return err
			}
			directoryAddr, err = network.GetEigenDADirectory()
			if err != nil {
				return fmt.Errorf("error getting EigenDADirectory address for network %s: %w", network, err)
			}
			logger.Printf("Auto-detected EigenDADirectory address %s for network %s", directoryAddr, networkName)
		} else {
			logger.Printf("No network provided, attempting to auto-detect EigenDADirectory address from chain ID.")
			chainID, err := client.ChainID(ctx.Context)
			if err != nil {
				return fmt.Errorf("failed to get chain ID from Ethereum client: %w", err)
			}
			if chainID == nil {
				return fmt.Errorf("failed to get chain ID from Ethereum client")
			}
			network, err := proxycmn.DefaultEigenDANetworkFromChainID(chainID.String())
			if err != nil {
				return fmt.Errorf("error determining EigenDA network from chain ID %s: %w", chainID, err)
			}
			directoryAddr, err = network.GetEigenDADirectory()
			if err != nil {
				return fmt.Errorf("error getting EigenDADirectory address for network %s: %w", network, err)
			}
			logger.Printf("Auto-detected EigenDADirectory address %s for chain ID %s (network: %s)", directoryAddr, chainID, network)
		}
	}

	// Validate directory address
	if !gethcommon.IsHexAddress(directoryAddr) {
		return fmt.Errorf("invalid EigenDADirectory address: %s", directoryAddr)
	}

	// Use the directory reader from core/eth package
	directoryReader, err := eth.NewEigenDADirectoryReader(directoryAddr, client)
	if err != nil {
		return fmt.Errorf("failed to create EigenDADirectory reader: %w", err)
	}

	// Get all addresses
	addressMap, err := directoryReader.GetAllAddresses()
	if err != nil {
		return fmt.Errorf("failed to get addresses from directory: %w", err)
	}

	// Output results
	switch outputFormat {
	case "table":
		printTable(addressMap)
	case "csv":
		printCSV(addressMap)
	case "json":
		printJSON(addressMap)
	}

	return nil
}

func printTable(addressMap map[string]gethcommon.Address) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Contract Name", "Address"})

	for name, addr := range addressMap {
		t.AppendRow(table.Row{name, addr.Hex()})
	}

	t.Render()
}

func printCSV(addressMap map[string]gethcommon.Address) {
	fmt.Println("Contract Name,Address")
	for name, addr := range addressMap {
		fmt.Printf("%s,%s\n", name, addr.Hex())
	}
}

func printJSON(addressMap map[string]gethcommon.Address) {
	fmt.Println("[")
	i := 0
	for name, addr := range addressMap {
		comma := ","
		if i == len(addressMap)-1 {
			comma = ""
		}
		fmt.Printf("  {\"name\": \"%s\", \"address\": \"%s\"}%s\n", name, addr.Hex(), comma)
		i++
	}
	fmt.Println("]")
}
