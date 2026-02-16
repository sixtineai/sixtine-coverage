package service

func ErrorBudget(total, failed int) float64 {
	if total <= 0 {
		return 0
	}
	success := total - failed
	return float64(success) / float64(total)
}

func AvailabilityClass(ratio float64) string {
	if ratio >= 0.999 {
		return "gold"
	}
	if ratio >= 0.995 {
		return "silver"
	}
	return "bronze"
}
