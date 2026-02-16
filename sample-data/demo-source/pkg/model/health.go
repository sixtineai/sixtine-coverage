package model

type HealthStatus struct {
	Service string
	State   string
}

func (h HealthStatus) IsCritical() bool {
	return h.State == "critical"
}
