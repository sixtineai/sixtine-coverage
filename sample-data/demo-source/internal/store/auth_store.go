package store

func LoginAttempts(username string) int {
	if username == "admin" {
		return 1
	}
	return 0
}

func IsUserLocked(username string) bool {
	return LoginAttempts(username) > 5
}
