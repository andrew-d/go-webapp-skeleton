package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/tylerb/graceful"
	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"

	"github.com/andrew-d/go-webapp-skeleton/conf"
	"github.com/andrew-d/go-webapp-skeleton/datastore"
	"github.com/andrew-d/go-webapp-skeleton/datastore/database"
	"github.com/andrew-d/go-webapp-skeleton/middleware"
	"github.com/andrew-d/go-webapp-skeleton/router"
	"github.com/andrew-d/go-webapp-skeleton/static"
)

func main() {
	log.Printf("initializing project_name=%q version=%q revision=%q",
		conf.ProjectName,
		conf.Version,
		conf.Revision)

	// Connect to the database.
	db, err := database.Connect(conf.C.DbType, conf.C.DbConn)
	if err != nil {
		log.Printf("could not connect to database err=%q db_type=%q db_conn=%q",
			err,
			conf.C.DbType,
			conf.C.DbConn)
		return
	}

	// Create datastore.
	ds := database.NewDatastore(db)

	// Create API router and add middleware.
	apiMux := router.API()
	apiMux.Use(middleware.Options)
	apiMux.Use(middleware.JSON)

	// Create web router.
	webMux := router.Web()

	// Create root mux and add common middleware.
	rootMux := goji.NewMux()
	rootMux.UseC(middleware.RequestID)
	rootMux.UseC(middleware.Logger)
	rootMux.UseC(middleware.Recoverer)
	rootMux.Use(middleware.SetHeaders)

	// Serve all static assets from the root.
	serveAssetAt := func(asset, path string) {
		info, _ := static.AssetInfo(asset)
		modTime := info.ModTime()
		data := static.MustAsset(asset)

		webMux.Handle(pat.Get(path), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("debug: serving asset: %s", asset)
			http.ServeContent(w, r, asset, modTime, bytes.NewReader(data))
		}))
	}
	for _, asset := range static.AssetNames() {
		log.Printf("debug: adding route for asset: %s", asset)
		serveAssetAt(asset, "/static/"+asset)
	}

	// Special case a bunch of assets that should be served from the root.
	for _, asset := range []string{
		"clientaccesspolicy.xml",
		"crossdomain.xml",
		"favicon.ico",
		"humans.txt",
		"robots.txt",
	} {
		// Note: only serve if we have this asset.
		if _, err := static.Asset(asset); err == nil {
			log.Printf("debug: adding special route for asset: %s", asset)
			serveAssetAt(asset, "/"+asset)
		}
	}

	// Serve the index page if we have one.
	for _, asset := range []string{"index.html", "index.htm"} {
		// Note: only serve if we have this asset, and only serve the first
		// option.
		if _, err := static.Asset(asset); err == nil {
			log.Printf("debug: adding index route for asset: %s", asset)
			serveAssetAt(asset, "/")
			break
		}
	}

	// Mount the API/Web muxes last (since order matters).
	rootMux.HandleC(pat.New("/api/*"), apiMux)
	rootMux.HandleC(pat.New("/*"), webMux)

	// Create a top-level wrapper that implements ServeHTTP, so we can inject
	// the root (Background) context and any other contexts that we wish to
	// inject.
	outer := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		ctx = datastore.NewContext(ctx, ds)
		rootMux.ServeHTTPC(ctx, w, r)
	})

	// Start serving
	log.Printf("starting server on: %s", conf.C.HostString())
	graceful.Run(conf.C.HostString(), 10*time.Second, outer)
	log.Printf("server finished")
}
