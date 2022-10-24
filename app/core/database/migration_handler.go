package database

type MigrationHandler interface {
	Reset(uri string) error
	Migrate(uri string) error
}
