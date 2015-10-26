package webhelpers

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

const CommonLogTimeFormat = "06/Jan/2006:15:04:05 -0700"

// LogrusLogger is a middleware that will log each request recieved, along with
// some useful information, to the given logger.
func LogrusLogger(logger *logrus.Logger, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		entry := logger.WithFields(logrus.Fields{
			"request": r.RequestURI,
			"method":  r.Method,
			"remote":  r.RemoteAddr,
		})

		if id := r.Header.Get(RequestIDKey); id != "" {
			entry = entry.WithField("request_id", id)
		}

		// Wrap the writer so we can track data information.
		neww := WrapWriter(w)

		// Dispatch to the underlying handler.
		entry.Info("started handling request")
		h.ServeHTTP(neww, r)

		// Log final information.
		entry.WithFields(logrus.Fields{
			"bytes_written": neww.BytesWritten(),
			"status":        neww.Status(),
			"text_status":   http.StatusText(neww.Status()),
			"took":          time.Since(start),
		}).Info("completed handling request")
	}

	return http.HandlerFunc(fn)
}

// CommonLogger is a middleware that will log each request recieved in Common
// Log Format to the given io.Writer.  See the following URL for more
// information on the Common Log Format:
//
//     https://en.wikipedia.org/wiki/Common_Log_Format
func CommonLogger(out io.Writer, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		neww := WrapWriter(w)
		h.ServeHTTP(neww, r)

		t := time.Now()
		fmt.Fprintf(out, "%s - - [%s] \"%s %s %s\" %d %d\n",
			r.RemoteAddr,
			t.Format(CommonLogTimeFormat),
			r.Method,
			r.URL.Path,
			r.Proto,
			neww.Status(),
			neww.BytesWritten(),
		)
	}

	return http.HandlerFunc(fn)
}
