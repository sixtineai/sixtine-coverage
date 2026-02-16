package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/rest/v1/health", healthHandler)
	mux.HandleFunc("/api/rest/v1/tenants", tenantsHandler)

	addr := ":8080"
	fmt.Println("demo api listening on", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func tenantsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":1,"name":"Acme"}]`))
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
