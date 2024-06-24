package mw

import (
	"net/http"
)

// CORSMiddleware sets the CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		origin := r.Header.Get("Origin")
		requestMethod := r.Header.Get("Access-Control-Request-Method")
		requestHeaders := r.Header.Get("Access-Control-Request-Headers")

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", requestMethod)
			w.Header().Set("Access-Control-Allow-Headers", requestHeaders)
			w.Header().Set("Access-Control-Max-Age", "7200")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
