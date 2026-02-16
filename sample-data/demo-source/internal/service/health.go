package service

import "strings"

func HealthSummary(statuses []string) string {
	if len(statuses) == 0 {
		return "unknown"
	}
	return strings.Join(statuses, ",")
}

func IsHealthy(code int) bool {
	return code >= 200 && code < 400
}

func Classify(code int) string {
	if code >= 500 {
		return "critical"
	}
	if code >= 400 {
		return "degraded"
	}
	return "healthy"
}
