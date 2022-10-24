package database

import (
	"database/sql"
	"time"

	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	MAX_IDLE_CONNECTIONS     = 10
	MAX_OPEN_CONNECTIONS     = 25
	MAX_CONNECTION_LIFE_TIME = time.Hour
)

type SqlDataSource struct {
	db                    *gorm.DB
	_innerSqlReference    *sql.DB
	maxConnectionLifeTime time.Duration
	maxIdleConnections    int
	maxOpenConnections    int
	enableLogging         bool
}

func (ds *SqlDataSource) open(uri string) error {
	var err error
	if ds.db, err = ds.connectionResolver(uri); err != nil {
		return err
	}
	if err != nil {
		return err
	}
	ds._innerSqlReference, err = ds.db.DB()
	if err != nil {
		return err
	}
	ds._innerSqlReference.SetMaxIdleConns(ds.maxIdleConnections)
	ds._innerSqlReference.SetConnMaxLifetime(ds.maxConnectionLifeTime)
	ds._innerSqlReference.SetMaxOpenConns(ds.maxOpenConnections)
	return nil
}

func (ds *SqlDataSource) Connection() *gorm.DB {
	return ds.db
}

func (ds *SqlDataSource) Close() error {
	return ds._innerSqlReference.Close()
}

func (ds *SqlDataSource) connectionResolver(path string) (*gorm.DB, error) {

	newLogger := logger.New(
		logging.GetInstance(),
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // Disable color
		},
	)
	configuration := &gorm.Config{Logger: newLogger}

	return gorm.Open(sqlite.Open(path), configuration)

}

type Builder struct {
	path                  string
	maxOpenConnections    int
	maxConnectionLifeTime time.Duration
	maxIdleConnections    int
	enableLogging         bool
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build() *SqlDataSource {
	if b.maxOpenConnections == 0 {
		b.maxOpenConnections = MAX_OPEN_CONNECTIONS
	}
	if b.maxConnectionLifeTime == 0 {
		b.maxConnectionLifeTime = MAX_CONNECTION_LIFE_TIME
	}
	if b.maxIdleConnections == 0 {
		b.maxIdleConnections = MAX_IDLE_CONNECTIONS
	}
	if b.path == "" {
		panic("No database file path has been specified")
	}
	instance := &SqlDataSource{
		maxConnectionLifeTime: b.maxConnectionLifeTime,
		maxIdleConnections:    b.maxIdleConnections,
		maxOpenConnections:    b.maxOpenConnections,
		enableLogging:         b.enableLogging,
	}
	if err := instance.open(b.path); err != nil {
		panic(err)
	}

	return instance
}

func (b *Builder) Path(path string) *Builder {
	b.path = path
	return b
}

func (b *Builder) EnableLogging() *Builder {
	b.enableLogging = true
	return b
}
