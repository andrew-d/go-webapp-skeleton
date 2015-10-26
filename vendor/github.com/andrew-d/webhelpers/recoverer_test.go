package webhelpers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test that the Recoverer will recover from panics
func TestRecoverer(t *testing.T) {
	var run bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		run = true
		panic("foo bar")
	})

	mux := http.NewServeMux()
	mux.Handle("/", Recoverer(handler))

	recorder := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	mux.ServeHTTP(recorder, r)

	// The recoverer should have caught the panic (i.e. this code will run), and
	// return a 500 error.
	assert.True(t, run)
	assert.Equal(t, 500, recorder.Code)
}

// Test that the CustomRecoverer will pass appropriate information to the callback.
func TestCustomRecoverer(t *testing.T) {
	var (
		info  RecoverInformation
		cbRun bool
	)
	recoverer := CustomRecoverer(func(w http.ResponseWriter, r *http.Request, i RecoverInformation) {
		cbRun = true
		info = i
	})

	// Dummy request-id function
	requestID := func(h http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set(RequestIDKey, "the_id")
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(f)
	}

	var run bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		run = true
		panic("foo bar")
	})

	mux := http.NewServeMux()
	mux.Handle("/", requestID(recoverer(handler)))

	recorder := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	mux.ServeHTTP(recorder, r)

	// Both the handler and the callback should have run
	assert.True(t, cbRun)
	assert.True(t, run)

	// The information should be valid - proper error, request ID, etc.
	assert.Equal(t, "the_id", info.RequestID)
	assert.NotEmpty(t, info.Stack)

	assert.Equal(t, "foo bar", info.Error)
}
