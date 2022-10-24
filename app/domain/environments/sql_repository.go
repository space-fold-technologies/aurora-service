package environments

import (
	"github.com/space-fold-technologies/aurora-service/app/core/database"
)

type SQLEnvironmentRepository struct {
	dataSource database.DataSource
}

func NewRepository(dataSource database.DataSource) EnvironmentRepository {
	instance := new(SQLEnvironmentRepository)
	instance.dataSource = dataSource
	return instance
}

func (ser *SQLEnvironmentRepository) Add(Scope, Target string, Entries []*Entry) error {
	tx := ser.dataSource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	vars := make([]*EnvVar, 0)
	for _, entry := range Entries {
		vars = append(vars, &EnvVar{
			Scope:  Scope,
			Target: Target,
			Key:    entry.Key,
			Value:  entry.Value,
		})
	}
	if err := tx.CreateInBatches(vars, 10).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func (ser *SQLEnvironmentRepository) List(Scope, Target string) ([]*Entry, error) {
	sql := "SELECT key, value FROM environment_variable_tb WHERE scope = ? AND target = ?"
	connection := ser.dataSource.Connection()
	vars := make([]*Entry, 0)
	if err := connection.Raw(sql, Scope, Target).Find(&vars).Error; err != nil {
		return vars, err
	}
	return vars, nil
}
func (ser *SQLEnvironmentRepository) Remove(Keys []string) error {
	sql := "DELETE FROM environment_variable_tb WHERE key IN(?)"
	tx := ser.dataSource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Exec(sql, Keys).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
