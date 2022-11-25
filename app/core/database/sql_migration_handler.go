package database

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	//_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/mattn/go-sqlite3"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
)

type SqlMigrationHandler struct {
	driver source.Driver
}

func NewMigrationHandler() MigrationHandler {
	instance := &SqlMigrationHandler{}
	instance.initialize()
	return instance
}

func (smh *SqlMigrationHandler) initialize() {
	var err error
	smh.driver, err = WithInstance("resources/migrations")
	if err != nil {
		logging.GetInstance().Error(err)
	}
}
func (smh *SqlMigrationHandler) Reset(uri string) error {
	logging.GetInstance().Info("RESET - DATABASE")
	db, err := sql.Open("sqlite3", uri)
	if err != nil {
		return err
	}
	defer db.Close()
	if instance, err := sqlite3.WithInstance(db, &sqlite3.Config{}); err != nil {
		return err
	} else if driver, err := WithInstance("resources/migrations"); err != nil {
		return err
	} else if m, err := migrate.NewWithInstance("embedded", driver, "sqlite3", instance); err != nil {
		return err
	} else {
		return m.Down()
	}
}

func (smh *SqlMigrationHandler) Migrate(uri string) error {
	logging.GetInstance().Info("MIGRATE - DATABASE")
	db, err := sql.Open("sqlite3", uri)
	if err != nil {
		return err
	}
	defer db.Close()
	if instance, err := sqlite3.WithInstance(db, &sqlite3.Config{}); err != nil {
		return err
	} else if driver, err := WithInstance("resources/migrations"); err != nil {
		return err
	} else if m, err := migrate.NewWithInstance("embedded", driver, "sqlite3", instance); err != nil {
		return err
	} else {
		return m.Up()
	}
}
