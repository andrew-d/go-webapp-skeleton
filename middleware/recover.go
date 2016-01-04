package middleware

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"runtime"

	"goji.io"
	"golang.org/x/net/context"

	"github.com/andrew-d/go-webapp-skeleton/conf"
)

func Recoverer(h goji.Handler) goji.Handler {
	f := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		id := GetRequestID(ctx)

		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				// Get the stack (from here, so we don't have
				// an extraneous call)
				stack := make([]byte, 8*1024)
				stack = stack[:runtime.Stack(stack, false)]

				// Handle the panic
				handlePanic(id, err, stack)
			}
		}()

		h.ServeHTTPC(ctx, w, r)
	}
	return goji.HandlerFunc(f)
}

func handlePanic(requestId string, err interface{}, stack []byte) {
	log.Printf(
		"error: recovered from panic request_id=%q err=%q",
		requestId,
		err)

	if conf.C.IsDebug() {
		// Split the stack by newlines, prepend a tab to each line, and then re-join
		deSpaced := bytes.TrimRight(stack, "\r\n ")
		lines := bytes.Split(deSpaced, []byte{'\n'})
		prettyStack := bytes.Join(lines, []byte{'\n', '\t'})
		prettyStack = append([]byte{'\t'}, prettyStack...)
		prettyStack = append(prettyStack, '\n')

		os.Stderr.Write(prettyStack)
	}
}
