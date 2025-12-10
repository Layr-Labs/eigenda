package nameremapping

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Loads a name remapping from a YAML file.
//
// Expected YAML format:
//
//	"0xFfFfFfFfFfFfFfFfFfFfFfFfFfFfFfFfFfFfFfFf": "Traffic Generator"
//	"0x1234567890AbcdEF1234567890aBcdef12345678": "User1"
//	"0xAbCdEf1234567890aBcDeF1234567890AbCdEf12": "User2"
func LoadNameRemapping(path string) (map[string]string, error) {
	if path == "" {
		return nil, fmt.Errorf("remapping file path is empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read remapping file %q: %w", path, err)
	}

	var remapping map[string]string
	if err := yaml.Unmarshal(data, &remapping); err != nil {
		return nil, fmt.Errorf("parse remapping file %q: %w", path, err)
	}

	return remapping, nil
}

// Returns the appropriate label for an account based on remapping and cardinality settings.
// Remapped names are formatted as "Name (0x123456)" with the account ID truncated to 8 characters.
func GetAccountLabel(accountId string, remappedNames map[string]string, highCardinalityNames bool) string {
	if remappedName, found := remappedNames[accountId]; found && remappedName != "" {
		truncatedId := accountId
		if len(accountId) >= 8 {
			truncatedId = accountId[:8]
		}
		return fmt.Sprintf("%s (%s)", remappedName, truncatedId)
	}

	if highCardinalityNames {
		return accountId
	}

	return "0x0"
}
