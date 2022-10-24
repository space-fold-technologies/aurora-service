package database

import (
	"gorm.io/gorm"
)

type DataSource interface {
	Connection() *gorm.DB
	Close() error
}
