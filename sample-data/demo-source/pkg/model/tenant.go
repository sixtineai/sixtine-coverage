package model

type Tenant struct {
	ID      int
	Name    string
	Country string
}

func (t Tenant) DisplayName() string {
	if t.Country == "" {
		return t.Name
	}
	return t.Name + " (" + t.Country + ")"
}
