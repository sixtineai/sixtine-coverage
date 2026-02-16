package httpapi

import "net/http"

func tenantsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[ {"id":1,"name":"Acme"}, {"id":2,"name":"Helios"} ]`))
	case http.MethodPost:
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":3,"name":"Orion"}`))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
