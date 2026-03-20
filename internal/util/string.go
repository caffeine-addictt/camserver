package util

import "strings"

// MultilineString provides a
// pleasing way to write multiline strings.
func MultilineString(s ...string) string {
	return strings.Join(s, "\n")
}
