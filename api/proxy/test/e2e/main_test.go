package e2e

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		fmt.Println("Skipping proxy e2e tests in short mode")
		os.Exit(0)
		return
	}
	code := m.Run()
	os.Exit(code)
}
