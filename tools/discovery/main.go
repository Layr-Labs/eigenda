package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"slices"
	"strings"

	proxycmn "github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/common/geth"
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
		Usage: fmt.Sprintf(`The EigenDA network to discover (one of: %s, %s, %s).
Must match the chain-id of the ethereum rpc url provided. Used to select the hardcoded default EigenDADirectory address.
That address can be overridden by providing the --%s flag.`,
			proxycmn.MainnetEigenDANetwork,
			proxycmn.SepoliaTestnetEigenDANetwork,
			proxycmn.HoodiTestnetEigenDANetwork,
			discoverAddressFlag.Name,
		),
		Required: true,
		EnvVars:  []string{"NETWORK"},
		Action: func(ctx *cli.Context, v string) error {
			if v == "" {
				// if no network is provided, we will try to auto-detect it from the chain ID
				return nil
			}
			// try to parse the network from the string.
			// this will validate the network and return an error if it's invalid.
			_, err := proxycmn.EigenDANetworkFromString(v)
			if err != nil {
				return fmt.Errorf("flag validation: %w", err)

			}
			return nil
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
	network, err := proxycmn.EigenDANetworkFromString(ctx.String(networkFlag.Name))
	if err != nil {
		return err
	}

	// Simple logging
	logger := log.New(os.Stderr, "[discovery] ", log.LstdFlags)

	client, err := geth.SafeDial(ctx.Context, rpcURL)
	if err != nil {
		return fmt.Errorf("dial Ethereum node: %w", err)
	}
	sanitizedUrl := geth.SanitizeRpcUrl(rpcURL)
	logger.Printf("Connected to Ethereum node at %s", sanitizedUrl)
	validateNetworkAndEthRpcChainIDMatch(ctx.Context, network, client)

	directoryAddr := ctx.String(discoverAddressFlag.Name)
	if directoryAddr == "" {
		directoryAddr = network.GetEigenDADirectory()
		logger.Printf("No explicit directory address provided, auto-detected EigenDADirectory address %s for network %s", directoryAddr, network)
	}

	// Validate directory address
	if !gethcommon.IsHexAddress(directoryAddr) {
		return fmt.Errorf("invalid EigenDADirectory address: %s", directoryAddr)
	}

	addressMap, err := GetContractAddressMap(
		context.Background(),
		client,
		gethcommon.HexToAddress(directoryAddr))
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
	jsonBytes, err := json.MarshalIndent(addressMap, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonBytes))
}

func validateNetworkAndEthRpcChainIDMatch(ctx context.Context, network proxycmn.EigenDANetwork, client *ethclient.Client) {
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("Failed to get chain ID from Ethereum client: %v", err)
	}
	if chainID == nil {
		log.Fatal("Received nil chain ID from Ethereum client")
	}

	expectedNetwork, err := proxycmn.EigenDANetworksFromChainID(chainID.String())
	if err != nil {
		log.Fatalf("Failed to get expected network from chain ID: %v", err)
	}
	if !slices.Contains(expectedNetwork, network) {
		log.Fatalf("Network mismatch: provided network %s is not part of the networks %v for chain ID %s",
			network, expectedNetwork, chainID.String())
	}
}
