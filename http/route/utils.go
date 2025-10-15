package route

import (
	"strings"
)

// GetWildcards returns the wildcards from the route pattern
//
// Parameters:
//
//   - pattern: The route pattern
//
// Returns:
//
//   - string: The parsed pattern
//   - []string: The wildcards
func GetWildcards(pattern string) (string, []string) {
	var wildcards []string
	for i := 0; i < len(pattern); i++ {
		// Check for '*' wildcard
		if pattern[i] == '*' {
			pattern = pattern[:i] + "{*}" + pattern[i+1:]
			i += 2
			wildcards = append(wildcards, "*")
		} else if pattern[i] == '{' {
			j := strings.IndexByte(pattern[i:], '}')
			if j != -1 {
				pattern = pattern[:i] + "{" + pattern[i+1:i+j] + "}" + pattern[i+j+1:]
				i += j
				wildcards = append(wildcards, pattern[i-j+1:i+j])
			}
		}
	}
	return pattern, wildcards
}

// SplitPattern returns the method and the path from the pattern
//
// Parameters:
//
//   - pattern: The pattern to split
//
// Returns:
//
//   - string: The method
//   - string: The path
//   - error: The error if any
func SplitPattern(pattern string) (string, string, error) {
	// Trim the pattern
	strings.Trim(pattern, " ")

	// Check if the pattern is empty
	if pattern == "" {
		return "", "", ErrEmptyPattern
	}

	// Iterate over the pattern
	var method string
	path := pattern
	for i, char := range pattern {
		// Split the pattern by the first space
		if char == ' ' {
			// Get the method and the path
			method = pattern[:i]
			path = pattern[i+1:]
			break
		}
	}

	// Trim the path
	strings.Trim(path, " ")

	return method, path, nil
}
