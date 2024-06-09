package e2e_test

import (
	"os"
	"testing"
)

var (
	runOptimismIntegrationTests bool
	runTestnetIntegrationTests  bool
)

func ParseEnv() {
	runOptimismIntegrationTests = os.Getenv("OPTIMISM") == "true" || os.Getenv("OPTIMISM") == "1"
	runTestnetIntegrationTests = os.Getenv("TESTNET") == "true" || os.Getenv("TESTNET") == "1"
}

func TestMain(m *testing.M) {
	ParseEnv()
	code := m.Run()
	os.Exit(code)
}
