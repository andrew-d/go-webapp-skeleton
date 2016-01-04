package main

import (
	"bytes"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/andrew-d/webhelpers"
	"github.com/jmoiron/sqlx"
	"github.com/tylerb/graceful"
	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"

	"github.com/andrew-d/go-webapp-skeleton/conf"
	"github.com/andrew-d/go-webapp-skeleton/datastore"
	"github.com/andrew-d/go-webapp-skeleton/datastore/database"
	"github.com/andrew-d/go-webapp-skeleton/log"
	"github.com/andrew-d/go-webapp-skeleton/middleware"
	"github.com/andrew-d/go-webapp-skeleton/router"
	"github.com/andrew-d/go-webapp-skeleton/static"
)

// Generic structure that holds all created things - database connection,
// datastore, etc.
type Vars struct {
	db  *sqlx.DB
	ds  datastore.Datastore
	log *logrus.Logger
}

func main() {
	var vars Vars

	// Create logger.
	vars.log = log.NewLogger()
	vars.log.WithFields(logrus.Fields{
		"project_name": conf.ProjectName,
		"version":      conf.Version,
		"revision":     conf.Revision,
	}).Info("initializing...")

	// Connect to the database.
	db, err := database.Connect(conf.C.DbType, conf.C.DbConn)
	if err != nil {
		vars.log.WithFields(logrus.Fields{
			"err":     err,
			"db_type": conf.C.DbType,
			"db_conn": conf.C.DbConn,
		}).Error("Could not connect to database")
		return
	}
	vars.db = db

	// Create datastore.
	vars.ds = database.NewDatastore(db)

	// Create API router and add middleware.
	apiMux := router.API()
	apiMux.Use(middleware.Options)
	apiMux.Use(middleware.JSON)

	// Create web router and add middleware.
	webMux := router.Web()
	webMux.Use(webhelpers.Recoverer)
	webMux.UseC(ContextMiddleware(&vars))
	webMux.Use(middleware.SetHeaders)

	// "Mount" the API mux under the web mux, on the "/api" prefix.
	webMux.HandleC(pat.New("/api/*"), apiMux)

	// Serve all static assets.
	serveAssetAt := func(asset, path string) {
		info, _ := static.AssetInfo(asset)
		modTime := info.ModTime()
		data := static.MustAsset(asset)

		webMux.Handle(pat.Get(path), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars.log.Debugf("serving asset: %s", asset)
			http.ServeContent(w, r, asset, modTime, bytes.NewReader(data))
		}))
	}
	for _, asset := range static.AssetNames() {
		vars.log.Debugf("adding route for asset: %s", asset)
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
			vars.log.Debugf("adding special route for asset: %s", asset)
			serveAssetAt(asset, "/"+asset)
		}
	}

	// Serve the index page if we have one.
	for _, asset := range []string{"index.html", "index.htm"} {
		// Note: only serve if we have this asset, and only serve the first
		// option.
		if _, err := static.Asset(asset); err == nil {
			vars.log.Debugf("adding index route for asset: %s", asset)
			serveAssetAt(asset, "/")
			break
		}
	}

	// We wrap the Request ID middleware and our logger 'outside' the mux, so
	// all requests (including ones that aren't matched by the router) get
	// logged.
	var handler http.Handler = webMux
	handler = webhelpers.LogrusLogger(vars.log, handler)
	handler = webhelpers.RequestID(handler)

	// Start serving
	vars.log.Infof("starting server on: %s", conf.C.HostString())
	graceful.Run(conf.C.HostString(), 10*time.Second, handler)
	vars.log.Info("server finished")
}

// ContextMiddleware will add our variables to the per-request context.
func ContextMiddleware(vars *Vars) func(goji.Handler) goji.Handler {
	mfn := func(h goji.Handler) goji.Handler {
		fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			ctx = datastore.NewContext(ctx, vars.ds)
			ctx = log.NewContext(ctx, vars.log)
			h.ServeHTTPC(ctx, w, r)
		}

		return goji.HandlerFunc(fn)
	}

	return mfn
}
