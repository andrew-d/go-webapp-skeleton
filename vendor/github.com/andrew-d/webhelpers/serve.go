package webhelpers

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

func serveContent(w http.ResponseWriter, r *http.Request, f http.File, fname string) {
	var modtime time.Time
	if fi, err := f.Stat(); err != nil {
		modtime = fi.ModTime()
	}

	http.ServeContent(w, r, fname, modtime, f)
}

// ServeFile creates a net/http-style Handler that serves the file at the given
// path.
func ServeFile(fpath string) http.Handler {
	// We should be able to open the file with no errors
	f, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}
	f.Close()

	// Saved for below.
	fname := path.Base(fpath)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "HEAD" {
			return
		}

		f, err := os.Open(fpath)
		if err != nil {
			http.Error(w, fmt.Sprintf("internal error opening file: %s", err),
				http.StatusInternalServerError)
			return
		}
		defer f.Close()

		serveContent(w, r, f, fname)
	})
}

// ServeDirectory creates a net/http-style Handler that serves the directory
// given.  It will properly prevent directory traversal attacks.
func ServeDirectory(dpath string) http.Handler {
	if !filepath.IsAbs(dpath) {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		dpath = filepath.Join(wd, dpath)
	}

	// Validate that this is in fact a directory.
	f, err := os.Open(dpath)
	if err != nil {
		panic(err)
	}
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	if !fi.IsDir() {
		panic(fmt.Sprintf("%s is not a directory", dpath))
	}
	f.Close()

	dir := http.Dir(dpath)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "HEAD" {
			return
		}

		f, err := dir.Open(r.URL.Path)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		defer f.Close()

		serveContent(w, r, f, r.URL.Path)
	})
}
