package config

import "strings"

// Converts a string in CamelCase to SCREAMING_SNAKE_CASE.
// Rules:
//
//  1. Insert underscore before any uppercase letter that follows a non-uppercase letter
//     Examples: "myField" -> "MY_FIELD", "field123Name" -> "FIELD123_NAME"
//
//  2. When N consecutive uppercase letters are followed by a lowercase letter:
//     - If only a single lowercase letter follows, keep it grouped with the uppercase letters
//     (This handles common pluralization patterns like "URLs", "IDs", etc. Without this exception,
//     "URLs" would become "UR_LS" instead of "URLS", which breaks the semantic meaning of the acronym)
//     - If multiple lowercase letters follow, split before the last uppercase letter
//     Examples: "IPAddress" -> "IP_ADDRESS", "URLs" -> "URLS", "IDs" -> "IDS"
//
//  3. When N consecutive uppercase letters are at the end (not followed by lowercase):
//     - Group all uppercase letters together (no split)
//     Examples: "NodeID" -> "NODE_ID", "ServerHTTP" -> "SERVER_HTTP"
//
//  4. Strings that are already all uppercase remain unchanged
//     Examples: "FIELD" -> "FIELD", "HTTP" -> "HTTP"
func toScreamingSnakeCase(s string) string {
	if s == "" {
		return ""
	}

	var result strings.Builder

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if i == 0 {
			// First character, don't prepend underscore
			result.WriteRune(r)
			continue
		}

		prev := runes[i-1]
		isCurrentUpper := r >= 'A' && r <= 'Z'
		isPrevUpper := prev >= 'A' && prev <= 'Z'
		isCurrentLower := r >= 'a' && r <= 'z'

		// Insert underscore if:
		// 1. Current is uppercase, previous is not uppercase (camelCase boundary)
		// 2. Current is lowercase, previous is uppercase, and there are multiple consecutive uppercase before
		//    This handles the transition from consecutive uppercase to lowercase
		//    e.g., "YAMLParser" -> at 'a', we need underscore before 'P'

		if isCurrentUpper && !isPrevUpper {
			// Transition from lowercase/other to uppercase: "myField" -> "my_Field"
			result.WriteRune('_')
		} else if isCurrentLower && isPrevUpper && i >= 2 {
			// We're at a lowercase letter after uppercase(s)
			// Check if there were multiple consecutive uppercase letters before this
			prevPrev := runes[i-2]
			isPrevPrevUpper := prevPrev >= 'A' && prevPrev <= 'Z'

			if isPrevPrevUpper {
				// Multiple uppercase letters followed by lowercase
				// Check if this is a single lowercase letter or if multiple lowercase letters follow
				// nolint:staticcheck
				isSingleLowercase := i == len(runes)-1 || !(runes[i+1] >= 'a' && runes[i+1] <= 'z')

				if !isSingleLowercase {
					// Multiple lowercase letters follow, so split before the last uppercase letter
					// e.g., "YAMLParser" at 'a': need underscore before 'P'
					// Remove the last character we wrote (the last uppercase letter)
					resultStr := result.String()
					result.Reset()
					result.WriteString(resultStr[:len(resultStr)-1])
					result.WriteRune('_')
					result.WriteRune(prev)
				}
				// If single lowercase, keep it grouped with the uppercase letters (no split)
				// e.g., "URLs" -> "URLS", "IDs" -> "IDS"
			}
		}

		result.WriteRune(r)
	}

	return strings.ToUpper(result.String())
}
