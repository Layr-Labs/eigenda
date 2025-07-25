package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"

	contractIEigenDADirectory "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDADirectory"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli"
)

var (
	version   = ""
	gitCommit = ""
	gitDate   = ""
)

// ContractInfo represents a contract name and its address
type ContractInfo struct {
	Name    string
	Address string
}

func main() {
	bi, _ := debug.ReadBuildInfo()
	app := cli.NewApp()
	app.Version = bi.Main.Version + bi.Main.Sum + bi.Main.Path
	app.Name = "eigenda-directory"
	app.Usage = "EigenDA Directory Contract Address Discovery Tool"
	app.Description = "Tool for fetching all contract addresses from the EigenDADirectory contract"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "rpc-url",
			Usage:    "Ethereum RPC URL",
			Required: true,
		},
		cli.StringFlag{
			Name:  "directory-address",
			Usage: "EigenDADirectory contract address (required if not using auto-discovery)",
		},
		cli.StringFlag{
			Name:  "output",
			Usage: "Output format: table, csv, json (default: table)",
			Value: "table",
		},
	}

	app.Action = discoverAddresses
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("application failed: %v", err)
	}
}

func discoverAddresses(ctx *cli.Context) error {
	rpcURL := ctx.String("rpc-url")
	directoryAddr := ctx.String("directory-address")
	outputFormat := strings.ToLower(ctx.String("output"))

	// Validate output format
	validFormats := map[string]bool{"table": true, "csv": true, "json": true}
	if !validFormats[outputFormat] {
		return fmt.Errorf("invalid output format: %s. Must be one of: table, csv, json", outputFormat)
	}

	// Simple logging
	logger := log.New(os.Stderr, "[discovery] ", log.LstdFlags)

	// Connect to Ethereum
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum node at %s: %w", rpcURL, err)
	}
	logger.Printf("Connected to Ethereum node at %s", rpcURL)

	// Validate directory address
	if directoryAddr == "" {
		return fmt.Errorf("directory-address is required")
	}
	if !common.IsHexAddress(directoryAddr) {
		return fmt.Errorf("invalid EigenDADirectory address: %s", directoryAddr)
	}

	// Create contract binding
	contractAddress := common.HexToAddress(directoryAddr)
	directory, err := contractIEigenDADirectory.NewContractIEigenDADirectory(contractAddress, client)
	if err != nil {
		return fmt.Errorf("failed to create contract binding: %w", err)
	}

	// Get all contract names
	names, err := directory.GetAllNames(&bind.CallOpts{})
	if err != nil {
		return fmt.Errorf("failed to get contract names: %w", err)
	}

	// Get addresses for each name
	var contracts []ContractInfo

	for _, name := range names {
		addr, err := directory.GetAddress0(&bind.CallOpts{}, name)
		if err != nil {
			logger.Printf("Warning: Failed to get address for %s: %v", name, err)
			continue
		}
		contracts = append(contracts, ContractInfo{
			Name:    name,
			Address: addr.Hex(),
		})
	}

	// Output results
	switch outputFormat {
	case "table":
		printTable(contracts)
	case "csv":
		printCSV(contracts)
	case "json":
		printJSON(contracts)
	}

	return nil
}

func printTable(contracts []ContractInfo) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Contract Name", "Address"})

	for _, c := range contracts {
		t.AppendRow(table.Row{c.Name, c.Address})
	}

	t.Render()
}

func printCSV(contracts []ContractInfo) {
	fmt.Println("Contract Name,Address")
	for _, c := range contracts {
		fmt.Printf("%s,%s\n", c.Name, c.Address)
	}
}

func printJSON(contracts []ContractInfo) {
	fmt.Println("[")
	for i, c := range contracts {
		comma := ","
		if i == len(contracts)-1 {
			comma = ""
		}
		fmt.Printf("  {\"name\": \"%s\", \"address\": \"%s\"}%s\n", c.Name, c.Address, comma)
	}
	fmt.Println("]")
}
