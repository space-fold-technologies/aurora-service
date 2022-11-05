package authorization

import (
	"github.com/space-fold-technologies/aurora-service/app/core/database"
	"github.com/space-fold-technologies/aurora-service/app/domain/users"
)

type SQLAuthorizationRepository struct {
	dataSource database.DataSource
}

func NewRepository(dataSource database.DataSource) AuthorizationRepository {
	instance := new(SQLAuthorizationRepository)
	instance.dataSource = dataSource
	return instance
}

func (ar *SQLAuthorizationRepository) Fetch(email string) (*users.User, error) {
	result := &users.User{}
	connection := ar.dataSource.Connection()
	if err := connection.Preload("Teams").Preload("Permissions").Where("email = ?", email).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
