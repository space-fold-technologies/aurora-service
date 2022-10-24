package database

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
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
	m, err := migrate.NewWithSourceInstance("embedded", smh.driver, uri)
	if err != nil {
		return err
	}
	return m.Down()
}

func (smh *SqlMigrationHandler) Migrate(uri string) error {
	logging.GetInstance().Info("MIGRATE - DATABASE")
	m, err := migrate.NewWithSourceInstance("embedded", smh.driver, uri)
	if err != nil {
		return err
	}
	return m.Up()
}
