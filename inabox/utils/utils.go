package utils

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func RunCommand(name string, args ...string) error {
	log.Printf("Running command: %s\n", strings.Join(append([]string{name}, args...), " "))
	cmd := exec.Command(name, args...)

	// Set the output to the corresponding os.Stdout and os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command and wait for it to finish
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

func StartCommand(name string, args ...string) (*exec.Cmd, error) {
	log.Printf("Running command: %s\n", strings.Join(append([]string{name}, args...), " "))
	cmd := exec.Command(name, args...)

	// Set the output to the corresponding os.Stdout and os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command and wait for it to finish
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}
