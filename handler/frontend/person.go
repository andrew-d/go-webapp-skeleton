package frontend

import (
	"net/http"
	"strconv"

	"goji.io/pat"
	"golang.org/x/net/context"

	"github.com/andrew-d/go-webapp-skeleton/datastore"
	"github.com/andrew-d/go-webapp-skeleton/handler"
	"github.com/andrew-d/go-webapp-skeleton/log"
)

// ListPeople shows a list of all people
//
//     GET /people
//
func ListPeople(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var (
		limit  = handler.ToLimit(r)
		offset = handler.ToOffset(r)
	)

	people, err := datastore.ListPeople(ctx, limit, offset)
	if err != nil {
		log.FromContext(ctx).WithField("err", err).Error("Error listing people")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	renderTemplate(ctx, w, "person_list.tmpl", M{
		"People": people,
	})
}

// GetPerson accepts a request to retrieve information about a particular person.
//
//     GET /people/:person
//
func GetPerson(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var (
		idStr = pat.Param(ctx, "person")
	)

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	person, err := datastore.GetPerson(ctx, id)
	if err != nil {
		log.FromContext(ctx).WithField("err", err).Error("Error getting person")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	renderTemplate(ctx, w, "person_show.tmpl", M{
		"Person": person,
	})
}
