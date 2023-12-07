package utils

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

// Creates a directory if it doesn't exist.
func MustMkdirall(name string) {
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(name, os.ModePerm)
		if err != nil {
			log.Panicf("Failed to create directory. Error: %s", err)
		}
	}
}

// Changes current working directory.
func MustGetwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Panicf("Failed to get working directory. Error: %s", err)
	}
	return dir
}

// Changes current working directory.
func MustChdir(path string) {
	err := os.Chdir(path)
	if err != nil {
		log.Panicf("Failed to change directories. Error: %s", err)
	}

	dir := MustGetwd()
	log.Printf("Current Working Directory: %s\n", dir)
}

func RunCommandAndCaptureOutput(name string, args ...string) (*bytes.Buffer, error) {
	log.Printf("Running command: %s\n", strings.Join(append([]string{name}, args...), " "))
	cmd := exec.Command(name, args...)

	var stdout bytes.Buffer
	// Set the output to the corresponding os.Stdout and os.Stderr
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	// Start the command and wait for it to finish
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &stdout, cmd.Wait()
}

func ExecCmd(name string, args []string, envVars []string) error {
	cmd := exec.Command(name, args...)
	if len(envVars) > 0 {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, envVars...)
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	// TODO: When these are uncommented, the deployer sometimes fails to start anvil
	// cmd.Stdout = &out
	// cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), stderr.String())
	}
	fmt.Print(out.String())
	return nil
}

func MustReadFile(name string) []byte {
	data, err := os.ReadFile(name)
	if err != nil {
		log.Panicf("Failed to read file. Error: %s", err)
	}

	return data
}

func MustWriteFile(name string, data []byte) {
	err := os.WriteFile(name, data, 0644)
	if err != nil {
		log.Panicf("Failed to write file. Err: %s", err)
	}
}

func MustMarshalYaml(obj interface{}) []byte {
	bz, err := yaml.Marshal(obj)
	if err != nil {
		log.Panicf("Yaml serialization of config.lock failed: %v", err)
	}
	return bz
}

func MustUnmarshalYaml[T any](bz []byte, obj *T) {
	err := yaml.Unmarshal(bz, obj)
	if err != nil {
		log.Panicf(err.Error())
	}
}

func MustWriteObjectToFile[T any](name string, obj T) {
	MustWriteFile(name, MustMarshalYaml(obj))
}
