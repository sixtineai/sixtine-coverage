package utils

import "strings"

func Slug(input string) string {
	s := strings.TrimSpace(strings.ToLower(input))
	s = strings.ReplaceAll(s, " ", "-")
	return strings.ReplaceAll(s, "_", "-")
}
