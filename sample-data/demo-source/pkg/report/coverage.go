package report

type CoverageSummary struct {
	Statements int
	Covered    int
}

func (c CoverageSummary) Ratio() float64 {
	if c.Statements == 0 {
		return 0
	}
	return float64(c.Covered) / float64(c.Statements)
}
