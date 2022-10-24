package users

import "github.com/space-fold-technologies/aurora-service/app/core/database"

type SQLUserRepository struct {
	dataSource database.DataSource
}

func NewRepository(dataSource database.DataSource) UserRepository {
	instance := new(SQLUserRepository)
	instance.dataSource = dataSource
	return instance
}

func (sur *SQLUserRepository) Create(entry *UserEntry) error {
	return nil
}
func (sur *SQLUserRepository) List() ([]*User, error) {
	return nil, nil
}
func (sur *SQLUserRepository) Remove(email string) error {
	return nil
}
