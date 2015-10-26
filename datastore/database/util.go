package database

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

func RebindInsert(db *sqlx.DB, q string) string {
	q = db.Rebind(q)
	q = strings.TrimRight(q, " \t\n;")
	if db.DriverName() == "postgres" {
		q = q + " RETURNING id"
	}

	return q
}
