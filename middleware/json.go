package middleware

import (
	"net/http"
)

// JSON sets the Content-Type to 'application/json' and sets the
// 'Access-Control-Allow-Origin' to allow all requests.
func JSON(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
