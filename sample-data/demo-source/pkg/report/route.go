package report

type RouteStatusCount struct {
	Status int
	Count  int
}

func TotalRequests(counts []RouteStatusCount) int {
	total := 0
	for _, c := range counts {
		total += c.Count
	}
	return total
}
