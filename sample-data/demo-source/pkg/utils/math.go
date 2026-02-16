package utils

func Percent(part, total float64) float64 {
	if total <= 0 {
		return 0
	}
	return (part / total) * 100
}
