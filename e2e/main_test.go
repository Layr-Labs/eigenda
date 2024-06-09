package e2e_test

import (
	"flag"
	"os"
	"testing"
)

var (
	optimism                   bool
	runTestnetIntegrationTests bool
)

func init() {
	flag.BoolVar(&optimism, "optimism", false, "Run OP Stack integration tests")
	flag.BoolVar(&runTestnetIntegrationTests, "testnet-integration", false, "Run testnet-based integration tests")

}

func TestMain(m *testing.M) {
	println("Parsing flags")
	flag.Parse()
	println("Flags parsed")

	code := m.Run()
	os.Exit(code)
}
