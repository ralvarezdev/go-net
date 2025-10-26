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
//   - error: The error if any
func GetWildcards(pattern string) (parsedPattern string, wildcards []string, err error) {
	for i := 0; i < len(pattern); i++ {
		// Check for '*' wildcard
		switch pattern[i] {
		case '*':
			// Return the parsed pattern just before the '*' and the collected wildcards
			return pattern[:i], wildcards, nil
		case '{':
			// Look for the closing '}'
			j := strings.IndexByte(pattern[i:], '}')
			if j == -1 {
				return "", nil, ErrWildcardNotClosed
			}
			
			// Check if the wildcard is empty
			if j == 1 {
				return "", nil, ErrEmptyWildcard
			}
			
			// Check if the wildcard is '*'
			if j == 2 && pattern[i+1] == '*' {
				return pattern[:i], wildcards, nil
			}
	
			// Append the wildcard to the list
			wildcards = append(wildcards, pattern[i+1:i+j])
	
			// Move the index to the end of the wildcard
			i += j
		}
	}
	return pattern, wildcards, nil
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
