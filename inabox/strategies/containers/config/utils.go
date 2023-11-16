package config

import (
	"log"
	"os"
)

func ReadFile(name string) []byte {
	data, err := os.ReadFile(name)
	if err != nil {
		log.Panicf("Failed to read file. Error: %s", err)
	}

	return data
}

func WriteFile(name string, data []byte) {
	err := os.WriteFile(name, data, 0644)
	if err != nil {
		log.Panicf("Failed to write file. Err: %s", err)
	}
}
