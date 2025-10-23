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
func GetWildcards(pattern string) (parsedPattern string, wildcards []string) {
	for i := 0; i < len(pattern); i++ {
		// Check for '*' wildcard
		switch pattern[i] {
		case '*':
			pattern = pattern[:i] + "{*}" + pattern[i+1:]
			i += 2
			wildcards = append(wildcards, "*")
		case '{':
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
func SplitPattern(pattern string) (method, path string, err error) {
	// Trim the pattern
	pattern = strings.Trim(pattern, " ")

	// Check if the pattern is empty
	if pattern == "" {
		return "", "", ErrEmptyPattern
	}

	// Split the pattern by space
	parts := strings.SplitN(pattern, " ", 2)

	// Get the method
	method = strings.ToUpper(strings.Trim(parts[0], " "))

	// Get the path
	path = "/"
	if len(parts) > 1 {
		path = strings.Trim(parts[1], " ")
	}

	return method, path, nil
}
