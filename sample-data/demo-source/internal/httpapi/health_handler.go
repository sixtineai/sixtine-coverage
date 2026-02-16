package httpapi

import (
	"net/http"
	"time"
)

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	if time.Now().Unix()%17 == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"status":"degraded"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
