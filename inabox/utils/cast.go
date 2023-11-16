package utils

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

// Converts a private key to an address.
func GetAddress(privateKey string) string {
	cmd := exec.Command("cast", "wallet", "address", "--private-key", privateKey)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Print(fmt.Sprint(err) + ": " + stderr.String())
		log.Panicf("Failed to execute cast wallet command. Err: %s", err)
	}

	//log.Print("Cast wallet command ran succesfully")
	return strings.Trim(out.String(), "\n")
}

// From the Foundry book: "Perform a call on an account without publishing a transaction."
func GetLatestBlockNumber(rpcUrl string) int {
	cmd := exec.Command("cast", "bn", "--rpc-url", rpcUrl)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Print(fmt.Sprint(err) + ": " + stderr.String())
		log.Panicf("Failed to execute cast wallet command. Err: %s", err)
	}

	log.Print("Cast bn command ran succesfully")
	blockNum, err := strconv.ParseInt(strings.Trim(out.String(), "\n"), 10, 0)
	if err != nil {
		log.Print(fmt.Sprint(err) + ": " + stderr.String())
		log.Panicf("Failed parse integer from blocknum string. Err: %s", err)
	}
	return int(blockNum)
}

// From the Foundry book: "Perform a call on an account without publishing a transaction."
func CallContract(destination string, signature string, rpcUrl string) string {
	cmd := exec.Command(
		"cast", "call", destination, signature,
		"--rpc-url", rpcUrl)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Print(fmt.Sprint(err) + ": " + stderr.String())
		log.Panicf("Failed to execute cast wallet command. Err: %s", err)
	}

	log.Print("Cast call command ran succesfully")
	return strings.Trim(out.String(), "\n")
}
