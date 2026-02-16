package validator

func ValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}
