package shell

import "strings"

// Escape safely quotes values for POSIX shells.
func Escape(value string) string {
	escaped := strings.ReplaceAll(value, `'`, `'\''`)
	return "'" + escaped + "'"
}
