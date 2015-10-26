package main

import (
	"fmt"
	"html"
	"net/http"
	"os"

	"github.com/andrew-d/webhelpers"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	mux := http.NewServeMux()
	mux.Handle("/foo", handler)

	// Note: we wrap this "outside" the ServeMux, since we would like to log
	// requests that do not match any route.  If this middleware was attached to
	// each handler function, requests that did not match a handler would not be
	// logged.
	http.ListenAndServe(":8080", webhelpers.CommonLogger(os.Stdout, mux))
}
