package users

import (
	"github.com/google/uuid"
	"github.com/space-fold-technologies/aurora-service/app/domain/teams"
)

type Permission struct {
	ID   int64  `gorm:"AUTO_INCREMENT;column:id;primary_key"`
	Name string `gorm:"column:name;type:varchar;unique;not null"`
	Idx  int64  `gorm:"column:idx;"`
}

// Override for the table name in `GORM-ORM`
func (Permission) TableName() string {
	return "permission_tb"
}

type User struct {
	ID           int64         `gorm:"AUTO_INCREMENT;column:id;primary_key"`
	Name         string        `gorm:"column:name;type:varchar;unique;not null"`
	NickName     string        `gorm:"column:nick_name;type:varchar;unique;not null"`
	Identifier   uuid.UUID     `gorm:"column:identifier;type:uuid;unique;not null"`
	Email        string        `gorm:"column:email;type:varchar;unique;not null"`
	PasswordHash string        `gorm:"column:password;type:varchar;unique;not null"`
	Permissions  []Permission  `gorm:"many2many:user_permissions"`
	Teams        []*teams.Team `gorm:"many2many:user_teams"`
}

func (User) TableName() string {
	return "user_tb"
}
