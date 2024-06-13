package e2e_test

import (
	"os"
	"testing"
)

var (
	runTestnetIntegrationTests bool
	runIntegrationTests        bool
)

func ParseEnv() {
	runIntegrationTests = os.Getenv("INTEGRATION") == "true" || os.Getenv("INTEGRATION") == "1"
	runTestnetIntegrationTests = os.Getenv("TESTNET") == "true" || os.Getenv("TESTNET") == "1"
}

func TestMain(m *testing.M) {
	ParseEnv()
	code := m.Run()
	os.Exit(code)
}
