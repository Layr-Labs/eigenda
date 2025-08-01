package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"slices"
	"strings"

	proxycmn "github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
)

var (
	ethRpcUrlFlag = &cli.StringFlag{
		Name:     "eth-rpc-url",
		Usage:    "Ethereum RPC URL",
		EnvVars:  []string{"ETH_RPC_URL"},
		Required: true,
	}
	networkFlag = &cli.StringFlag{
		Name: "network",
		Usage: fmt.Sprintf(`The EigenDA network to discover (one of: %s, %s, %s, %s).
Must match the chain-id of the ethereum rpc url provided. Used to select the hardcoded default EigenDADirectory address.
That address can be overridden by providing the --%s flag.`,
			proxycmn.MainnetEigenDANetwork,
			proxycmn.HoleskyTestnetEigenDANetwork,
			proxycmn.HoleskyPreprodEigenDANetwork,
			proxycmn.SepoliaTestnetEigenDANetwork,
			discoverAddressFlag.Name,
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
	}
	discoverAddressFlag = &cli.StringFlag{
		Name:    "directory-address",
		Usage:   "EigenDADirectory contract address (overrides the default network address)",
		EnvVars: []string{"EIGENDA_DIRECTORY_ADDRESS"},
	}
	validOutputFormats = []string{"table", "csv", "json"}
	outputFormatFlag   = &cli.StringFlag{
		Name:    "output-format",
		Usage:   fmt.Sprintf("Output format. Must be one of: %v", validOutputFormats),
		Value:   "table",
		EnvVars: []string{"OUTPUT_FORMAT"},
		Action: func(ctx *cli.Context, v string) error {
			if !slices.Contains(validOutputFormats, strings.ToLower(v)) {
				return fmt.Errorf("invalid output format: %s. Must be one of: %v", v, validOutputFormats)
			}
			return nil
		},
	}
)

func main() {
	app := cli.NewApp()
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		app.Version = buildInfo.Main.Version
	}
	app.Name = "eigenda-directory"
	app.Usage = "EigenDA Directory Contract Address Discovery Tool"
	app.Description = "Tool for fetching all contract addresses from the EigenDADirectory contract on a specified EigenDA network."
	app.Flags = []cli.Flag{
		ethRpcUrlFlag,
		networkFlag,
		discoverAddressFlag,
		outputFormatFlag,
	}
	app.Action = discoverAddresses
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("application failed: %v", err)
	}
}

func discoverAddresses(ctx *cli.Context) error {
	outputFormat := strings.ToLower(ctx.String(outputFormatFlag.Name))
	rpcURL := ctx.String(ethRpcUrlFlag.Name)

	// Simple logging
	logger := log.New(os.Stderr, "[discovery] ", log.LstdFlags)

	// Connect to Ethereum
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("dial Ethereum node at %s: %w", rpcURL, err)
	}
	logger.Printf("Connected to Ethereum node at %s", rpcURL)

	directoryAddr := ctx.String(discoverAddressFlag.Name)
	// cases:
	// 1. directory address: use it directly
	// 2. no directory address, but network provided: use the default directory address for the network
	// 3. no directory address and no network: use the default network for the chainID
	if directoryAddr == "" {
		logger.Printf("No directory address provided, attempting to auto-detect.")
		networkName := ctx.String(networkFlag.Name)
		var network proxycmn.EigenDANetwork
		if networkName != "" {
			logger.Printf("network %s provided, attempting to use default directory address.", networkName)
			network, err = proxycmn.EigenDANetworkFromString(ctx.String(networkFlag.Name))
			if err != nil {
				return err
			}
		} else {
			logger.Printf("No network provided, attempting to auto-detect EigenDADirectory address from chain ID.")
			chainID, err := client.ChainID(ctx.Context)
			if err != nil {
				return fmt.Errorf("chainID call: %w", err)
			}
			if chainID == nil {
				return fmt.Errorf("received nil chainID from Ethereum client")
			}
			networks, err := proxycmn.EigenDANetworksFromChainID(chainID.String())
			if err != nil {
				return fmt.Errorf("EigenDANetworksFromChainID: %w", err)
			}
			network = networks[0]
			if len(networks) > 1 {
				logger.Printf("Multiple EigenDA networks found for chain ID %s: %v. Using the first one: %s", chainID, networks, network)
			} else {
				logger.Printf("Auto-detected EigenDA network from chain ID %s: %s", chainID, network)
			}
		}
		directoryAddr, err = network.GetEigenDADirectory()
		if err != nil {
			return fmt.Errorf("GetEigenDADirectory: %w", err)
		}
		logger.Printf("Auto-detected EigenDADirectory address %s for network %s", directoryAddr, network)
	}

	// Validate directory address
	if !gethcommon.IsHexAddress(directoryAddr) {
		return fmt.Errorf("invalid EigenDADirectory address: %s", directoryAddr)
	}

	// Use the directory reader from core/eth package
	directoryReader, err := eth.NewEigenDADirectoryReader(directoryAddr, client)
	if err != nil {
		return fmt.Errorf("NewEigenDADirectoryReader: %w", err)
	}

	addressMap, err := directoryReader.GetAllAddresses(&bind.CallOpts{Context: ctx.Context})
	if err != nil {
		return fmt.Errorf("GetAllAddresses from directory: %w", err)
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
		fmt.Printf("  {\"contract_name\": \"%s\", \"address\": \"%s\"}%s\n", name, addr.Hex(), comma)
		i++
	}
	fmt.Println("]")
}
