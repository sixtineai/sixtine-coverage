package model

type Session struct {
	Token     string
	TenantID  int
	ExpiresIn int
}

func (s Session) IsValid() bool {
	return s.Token != "" && s.ExpiresIn > 0
}
