package validator

import "net/http"

func HasJSONContentType(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "application/json"
}
