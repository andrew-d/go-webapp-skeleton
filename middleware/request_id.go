package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"goji.io"
	"golang.org/x/net/context"
)

type private struct{}

var requestIdKey private

var (
	prefix string
	reqid  uint64
)

func init() {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}

	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}

	prefix = fmt.Sprintf("%s/%s", hostname, b64[0:10])
}

// RequestID is a middleware that injects a request ID into the headers of each
// request. A request ID is a string of the form "host.example.com/random-0001",
// where "random" is a base62 random string that uniquely identifies this go
// process, and where the last number is an atomically incremented request
// counter.
//
// Note: this middleware is adapted from goji:
//	https://github.com/zenazn/goji/blob/master/web/middleware/request_id.go
func RequestID(h goji.Handler) goji.Handler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctr := atomic.AddUint64(&reqid, 1)
		id := fmt.Sprintf("%s-%06d", prefix, ctr)

		ctx = context.WithValue(ctx, requestIdKey, id)
		h.ServeHTTPC(ctx, w, r)
	}

	return goji.HandlerFunc(fn)
}

// GetRequestID retrieves the request ID (if any) from the given context.  It
// will return the empty string ("") if none was set.
func GetRequestID(ctx context.Context) string {
	val := ctx.Value(requestIdKey)
	if val != nil {
		return val.(string)
	}

	return ""
}
