package migrate

import (
	"github.com/BurntSushi/migration"
	"github.com/jmoiron/sqlx"
)

type Migrator struct {
	DbType string
}

func (m Migrator) rebind(s string) string {
	return sqlx.Rebind(sqlx.BindType(m.DbType), s)
}

// Setup will create all necessary tables and indexes in the database.
func (m Migrator) Setup(tx migration.LimitedTx) error {
	stmts := []string{
		peopleTable,
	}

	for _, stmt := range stmts {
		if _, err := tx.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}

// CreateDefaultPerson will insert a default (empty) list into the database.
func (m Migrator) CreateDefaultPerson(tx migration.LimitedTx) error {
	_, err := tx.Exec(createDefaultPerson)
	return err
}

const peopleTable = `
CREATE TABLE IF NOT EXISTS people (
	 id   INTEGER PRIMARY KEY AUTOINCREMENT
	,name TEXT NOT NULL
)
`

const createDefaultPerson = `
INSERT INTO people(name)
VALUES ("Joe Smith")
`
