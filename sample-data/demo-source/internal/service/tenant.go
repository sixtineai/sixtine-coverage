package service

import "strings"

type Tenant struct {
	ID   int
	Name string
}

func NormalizeTenantName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "Unnamed"
	}
	return strings.Title(strings.ToLower(name))
}

func IsReservedTenant(name string) bool {
	reserved := map[string]bool{"system": true, "root": true}
	return reserved[strings.ToLower(strings.TrimSpace(name))]
}
