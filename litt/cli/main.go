package main

import (
	"fmt"
	"os"
)

// main is the entry point for the LittDB cli.
func main() {
	err := buildCLIParser().Run(os.Args)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
