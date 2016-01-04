package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"goji.io/pat"
	"golang.org/x/net/context"

	"github.com/andrew-d/go-webapp-skeleton/datastore"
	"github.com/andrew-d/go-webapp-skeleton/handler"
	"github.com/andrew-d/go-webapp-skeleton/model"
)

// ListPeople accepts a request to retrieve a list of people.
//
//     GET /api/people
//
func ListPeople(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var (
		limit  = handler.ToLimit(r)
		offset = handler.ToOffset(r)
	)

	people, err := datastore.ListPeople(ctx, limit, offset)
	if err != nil {
		log.Printf("error: error listing people err=%q", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(people)
}

// GetPerson accepts a request to retrieve information about a particular person.
//
//     GET /api/people/:person
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
		log.Printf("error: error getting person err=%q", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(person)
}

// DeletePerson accepts a request to delete a person.
//
//     DELETE /api/people/:person
//
func DeletePerson(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var (
		idStr = pat.Param(ctx, "person")
	)

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = datastore.DeletePerson(ctx, id)
	if err != nil {
		log.Printf("error: error deleting person err=%q", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreatePerson accepts a request to add a new person.
//
//     POST /api/people
//
func CreatePerson(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// Unmarshal the person from the payload
	defer r.Body.Close()
	in := struct {
		Name string `json:"name"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate input
	if len(in.Name) < 1 {
		http.Error(w, "no name given", http.StatusBadRequest)
		return
	}

	// Create our 'normal' model.
	person := &model.Person{Name: in.Name}
	err := datastore.CreatePerson(ctx, person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(person)
}
