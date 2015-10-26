package webhelpers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test that the request ID is added and non-empty.
func TestRequestID(t *testing.T) {
	var run bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		run = true
		assert.True(t, len(r.Header.Get(RequestIDKey)) > 0)
	})

	mux := http.NewServeMux()
	mux.Handle("/", RequestID(handler))

	var w http.ResponseWriter = httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	mux.ServeHTTP(w, r)

	assert.True(t, run)
}
