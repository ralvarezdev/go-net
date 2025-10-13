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
