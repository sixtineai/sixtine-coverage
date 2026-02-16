package service

func CanLogin(username, password string) bool {
	if username == "" || password == "" {
		return false
	}
	return len(password) >= 8
}

func SessionTTLMinutes(role string) int {
	switch role {
	case "admin":
		return 30
	case "support":
		return 45
	default:
		return 20
	}
}
