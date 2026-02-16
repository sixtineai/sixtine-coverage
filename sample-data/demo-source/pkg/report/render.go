package report

import "fmt"

func RenderCoverageLabel(percent float64) string {
	return fmt.Sprintf("coverage %.1f%%", percent)
}
