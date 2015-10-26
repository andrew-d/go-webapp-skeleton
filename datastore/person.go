package datastore

import (
	"golang.org/x/net/context"

	"github.com/andrew-d/go-webapp-skeleton/model"
)

type PeopleStore interface {
	// ListPeople retrieves all people from the database, possibly with an
	// offset or limit provided.
	ListPeople(limit, offset int) ([]*model.Person, error)

	// GetPerson retrieves a person from the datastore for the given ID.
	GetPerson(id int64) (*model.Person, error)

	// CreatePerson saves a new person in the datastore.
	CreatePerson(person *model.Person) error

	// DeletePerson removes a person from the datastore.
	DeletePerson(id int64) error
}

func ListPeople(c context.Context, limit, offset int) ([]*model.Person, error) {
	return FromContext(c).ListPeople(limit, offset)
}

func GetPerson(c context.Context, id int64) (*model.Person, error) {
	return FromContext(c).GetPerson(id)
}

func CreatePerson(c context.Context, person *model.Person) error {
	return FromContext(c).CreatePerson(person)
}

func DeletePerson(c context.Context, id int64) error {
	return FromContext(c).DeletePerson(id)
}
