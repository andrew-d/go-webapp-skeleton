package router

import (
	"github.com/zenazn/goji/web"

	"github.com/andrew-d/go-webapp-skeleton/handler/api"
	"github.com/andrew-d/go-webapp-skeleton/handler/frontend"
)

func New() *web.Mux {
	mux := web.New()

	mux.Get("/api/people", api.ListPeople)
	mux.Post("/api/people", api.CreatePerson)
	mux.Get("/api/people/:person", api.GetPerson)
	mux.Delete("/api/people/:list", api.DeletePerson)

	mux.Get("/people", frontend.ListPeople)
	mux.Get("/people/:person", frontend.GetPerson)

	return mux
}
