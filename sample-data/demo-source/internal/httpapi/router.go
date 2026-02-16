package httpapi

import "net/http"

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/rest/v1/health", healthHandler)
	mux.HandleFunc("/api/rest/v1/tenants", tenantsHandler)
	mux.HandleFunc("/api/rest/v1/authentication/login", loginHandler)
	return mux
}
