package router

import (
	"fmt"
	"net/http"

	"goji.io"
	"goji.io/pat"

	"github.com/andrew-d/go-webapp-skeleton/handler/api"
	"github.com/andrew-d/go-webapp-skeleton/handler/frontend"
)

func API() *goji.Mux {
	mux := goji.SubMux()

	// We pass the routes as relative to the point where the API router
	// will be mounted.  The super-router will strip any prefix off for us.
	mux.HandleFuncC(pat.Get("/people"), api.ListPeople)
	mux.HandleFuncC(pat.Post("/people"), api.CreatePerson)
	mux.HandleFuncC(pat.Get("/people/:person"), api.GetPerson)
	mux.HandleFuncC(pat.Delete("/people/:person"), api.DeletePerson)

	// Add default 'not found' route that responds with JSON
	mux.HandleFunc(pat.New("/*"), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprint(w, `{"error":"not found"}`)
	})

	return mux
}

func Web() *goji.Mux {
	mux := goji.SubMux()

	mux.HandleFuncC(pat.Get("/people"), frontend.ListPeople)
	mux.HandleFuncC(pat.Get("/people/:person"), frontend.GetPerson)

	return mux
}
