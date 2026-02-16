package utils

func IsSuccess(code int) bool {
	return code >= 200 && code < 300
}

func IsServerError(code int) bool {
	return code >= 500
}
