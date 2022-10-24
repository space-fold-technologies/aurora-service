package teams

import "github.com/space-fold-technologies/aurora-service/app/core/database"

type SQLTeamRepository struct {
	dataSource database.DataSource
}

func NewRepository(dataSource database.DataSource) TeamRepository {
	instance := new(SQLTeamRepository)
	instance.dataSource = dataSource
	return instance
}

func (str *SQLTeamRepository) Create(team *Team) error {
	tx := str.dataSource.Connection().Begin()
	if err := tx.Create(team).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func (str *SQLTeamRepository) List() ([]*Team, error) {
	teams := make([]*Team, 0)
	sql := "SELECT id,  name, description FROM team_tb"
	connection := str.dataSource.Connection()
	if err := connection.Raw(sql).Find(&teams).Error; err != nil {
		return teams, err
	}
	return teams, nil
}
func (str *SQLTeamRepository) Remove(name string) error {
	sql := "DELETE FROM team_tb WHERE name = ?"
	tx := str.dataSource.Connection().Begin()
	if err := tx.Exec(sql, name).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
