package database

import (
	"github.com/jmoiron/sqlx"

	"github.com/andrew-d/go-webapp-skeleton/model"
)

type PeopleStore struct {
	db *sqlx.DB
}

func NewPeopleStore(db *sqlx.DB) *PeopleStore {
	return &PeopleStore{db}
}

func (s *PeopleStore) ListPeople(limit, offset int) ([]*model.Person, error) {
	people := []*model.Person{}
	err := s.db.Select(&people, s.db.Rebind(personListQuery), limit, offset)
	return people, err
}

func (s *PeopleStore) GetPerson(id int64) (*model.Person, error) {
	person := &model.Person{}
	err := s.db.Get(person, s.db.Rebind(personGetQuery), id)
	return person, err
}

func (s *PeopleStore) CreatePerson(person *model.Person) error {
	ret, err := s.db.Exec(RebindInsert(s.db, personInsertQuery), person.Name)
	if err != nil {
		return err
	}

	person.ID, _ = ret.LastInsertId()
	return nil
}

func (s *PeopleStore) DeletePerson(id int64) (err error) {
	var tx *sqlx.Tx

	tx, err = s.db.Beginx()
	if err != nil {
		return
	}

	// Automatically rollback/commit if there's an error.
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Remove the given Person
	if _, err = tx.Exec(s.db.Rebind(personDeleteQuery), id); err != nil {
		return
	}

	// Done!
	return nil
}

const personListQuery = `
SELECT *
FROM people
ORDER BY id DESC
LIMIT ? OFFSET ?
`

const personGetQuery = `
SELECT *
FROM people
WHERE id = ?
`

const personInsertQuery = `
INSERT
INTO people (
     name
)
VALUES (?)
`

const personDeleteQuery = `
DELETE
FROM people
WHERE id = ?
`
