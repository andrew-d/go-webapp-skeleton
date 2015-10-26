package database

import (
	"github.com/BurntSushi/migration"
	"github.com/jmoiron/sqlx"

	"github.com/andrew-d/go-webapp-skeleton/datastore"
	"github.com/andrew-d/go-webapp-skeleton/datastore/migrate"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func Connect(driver, conn string) (*sqlx.DB, error) {
	migrator := migrate.Migrator{driver}

	migration.DefaultGetVersion = migrator.GetVersion
	migration.DefaultSetVersion = migrator.SetVersion

	migrations := []migration.Migrator{
		migrator.Setup,
		migrator.CreateDefaultPerson,
	}

	db, err := migration.Open(driver, conn, migrations)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return sqlx.NewDb(db, driver), nil
}

func MustConnect(driver, conn string) *sqlx.DB {
	db, err := Connect(driver, conn)
	if err != nil {
		panic(err)
	}
	return db
}

func NewDatastore(db *sqlx.DB) datastore.Datastore {
	return struct {
		*PeopleStore
	}{
		NewPeopleStore(db),
	}
}
