package e2e_test

import (
	"os"
	"testing"
)

// Integration tests are run against memstore whereas.
// Testnetintegration tests are run against eigenda backend talking to testnet disperser.
// Some of the assertions in the tests are different based on the backend as well.
// e.g, in TestProxyServerCaching we only assert to read metrics with EigenDA
// when referencing memstore since we don't profile the eigenDAClient interactions
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
