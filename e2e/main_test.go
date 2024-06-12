package e2e_test

import (
	"os"
	"testing"
)

var (
	runTestnetIntegrationTests bool
)

func ParseEnv() {
	runTestnetIntegrationTests = os.Getenv("TESTNET") == "true" || os.Getenv("TESTNET") == "1"
}

func TestMain(m *testing.M) {
	ParseEnv()
	println("runTestnetIntegrationTests:", runTestnetIntegrationTests)
	code := m.Run()
	os.Exit(code)
}
