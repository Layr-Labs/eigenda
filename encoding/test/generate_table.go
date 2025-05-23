package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
)

func main() {
	generateTable()
}

func generateTable() {
	fmt.Printf("* Generating table with DimE=16, CosetSize=1048576\n")

	kzgConfig := &kzg.KzgConfig{
		G1Path:          "g1.point",
		G2Path:          "g2.point",
		CacheDir:        "SRSTables",
		SRSOrder:        268435456,
		SRSNumberToLoad: 16777216,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    true,
	}

	// Create prover
	p, err := prover.NewProver(kzgConfig, nil)
	if err != nil {
		log.Fatalf("Failed to create prover: %v", err)
	}

	// Generate table with specified parameters
	dimE := uint64(16)
	cosetSize := uint64(1048576)

	fmt.Printf("    DimE: %v\n", dimE)
	fmt.Printf("    CosetSize: %v\n", cosetSize)

	// Create SRS table and generate the precomputed table
	subTable, err := prover.NewSRSTable(p.KzgConfig.CacheDir, p.Srs.G1, p.KzgConfig.NumWorker)
	if err != nil {
		log.Fatalf("Failed to create SRS table: %v", err)
	}

	// Call GetSubTables which will generate the table if it doesn't exist
	_, err = subTable.GetSubTables(dimE, cosetSize)
	if err != nil {
		log.Fatalf("Failed to generate table: %v", err)
	}

	fmt.Printf("* Table generation completed successfully\n")
}
