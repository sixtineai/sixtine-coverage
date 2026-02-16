package validator

import "strings"

func ValidTenantName(name string) bool {
	name = strings.TrimSpace(name)
	return len(name) >= 3 && len(name) <= 64
}
