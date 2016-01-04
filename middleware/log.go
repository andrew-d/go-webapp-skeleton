package middleware

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"goji.io"
	"golang.org/x/net/context"
)

// Logger is a middleware that will log each request recieved, along with
// some useful information, using the 'log' package.
func Logger(h goji.Handler) goji.Handler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		// Generate per-request fields
		var fieldsBuf bytes.Buffer
		if id := GetRequestID(ctx); id != "" {
			fmt.Fprintf(&fieldsBuf, " request_id=%q", id)
		}

		fmt.Fprintf(&fieldsBuf,
			" method=%s url=%q remote_addr=%q",
			r.Method,
			r.URL.String(),
			r.RemoteAddr)

		// Print the pre-request log
		log.Printf("request started%s", fieldsBuf.String())

		// Wrap the writer so we can track data written, status, etc.
		wh := WrapWriter(w)

		// Dispatch to the underlying handler.
		start := time.Now()
		h.ServeHTTPC(ctx, wh, r)

		// Ensure that we've started flushing data to the client before
		// we stop the timer.
		if wh.Status() == 0 {
			wh.WriteHeader(http.StatusOK)
		}
		took := time.Since(start)

		// Fill in remainder of the request fields
		fmt.Fprintf(&fieldsBuf,
			" bytes_written=%d status=%d took=%q",
			wh.BytesWritten(),
			wh.Status(),
			took)

		// Log final information.
		log.Printf("request finished%s", fieldsBuf.String())
	}

	return goji.HandlerFunc(fn)
}
