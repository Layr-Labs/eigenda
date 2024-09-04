package configuration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// ParseJsonFiles parses json files and stores the configuration in the provided struct.
// Later json files override earlier ones.
// Values in the struct that do not appear in any json file will be left unchanged.
// Values in any json file that do not appear in the struct will cause this method to return an error.
func ParseJsonFiles[T interface{}](t *T, configFiles ...string) error {
	for _, configFile := range configFiles {
		err := ParseJsonFile(t, configFile)
		if err != nil {
			return err
		}
	}

	return nil
}

// ParseJsonFile parses a json file and stores the configuration in the provided struct.
// Values in the struct that do not appear in the json file will be left unchanged.
// Values in the json file that do not appear in the struct will cause this method to return an error.
func ParseJsonFile[T interface{}](t *T, configFile string) error {
	fileBytes, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = ParseJsonString(t, string(fileBytes))
	if err != nil {
		return fmt.Errorf("error parsing json file %s: %w", configFile, err)
	}

	return nil
}

// ParseJsonStrings parses json strings and stores the configuration in the provided struct.
// Later json strings override earlier ones.
// Values in the struct that do not appear in any json string will be left unchanged.
// Values in any json string that do not appear in the struct will cause this method to return an error.
func ParseJsonStrings[T interface{}](t *T, configStrings ...string) error {
	for _, configString := range configStrings {
		err := ParseJsonString(t, configString)
		if err != nil {
			return err
		}
	}

	return nil
}

// ParseJsonString parses a json string and stores the configuration in the provided struct.
// Values in the struct that do not appear in the json string will be left unchanged.
// Values in the json string that do not appear in the struct will cause this method to return an error.
func ParseJsonString[T interface{}](t *T, configString string) error {
	reader := bytes.NewReader([]byte(configString))
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(t)
	if err != nil {
		// Add the line number and column number to the error message.
		var syntaxErr *json.SyntaxError
		ok := errors.As(err, &syntaxErr)
		if !ok {
			// The json parser didn't tell us where the error was. Nothing we can do to extract line number.
			return err
		}

		lineNumber, columnNumber := determineErrorPosition(configString, syntaxErr.Offset)
		return fmt.Errorf("error parsing json string at line %d, column %d: %w", lineNumber, columnNumber, err)
	}

	return nil
}

// determineErrorPosition determines the line and column number of the error in the json string.
// Useful for debugging malformed json strings/files.
func determineErrorPosition(jsonString string, offset int64) (lineNumber int64, columnNumber int64) {
	// Lines are numbered starting from 1.
	lineNumber = 1

	for i, char := range jsonString {
		if int64(i) == offset {
			break
		}
		if char == '\n' {
			lineNumber++
			columnNumber = 0
		} else {
			columnNumber++
		}
	}
	return
}
