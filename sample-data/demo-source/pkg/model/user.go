package model

type User struct {
	ID       int
	Username string
	Role     string
}

func (u User) IsAdmin() bool {
	return u.Role == "admin"
}
