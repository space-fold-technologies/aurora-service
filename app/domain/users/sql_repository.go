package users

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/space-fold-technologies/aurora-service/app/core/database"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"gorm.io/gorm"
)

type SQLUserRepository struct {
	dataSource database.DataSource
	log        *logrus.Logger
}

func NewRepository(dataSource database.DataSource) UserRepository {
	instance := new(SQLUserRepository)
	instance.dataSource = dataSource
	instance.log = logging.GetInstance()
	return instance
}

func (sur *SQLUserRepository) Create(entry *UserEntry) error {
	sql := "INSERT INTO user_tb(identifier, name, nick_name, email, password) VALUES(?, ?, ?, ?, ?)"
	tx := sur.dataSource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			sur.log.Error("Failed transaction")
			tx.Rollback()
		}
	}()
	if err := tx.Exec(sql, entry.Identifier, entry.Name, entry.NickName, entry.Email, entry.PasswordHash).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func (sur *SQLUserRepository) Update(order *UpdateEntry) error {
	tx := sur.dataSource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			sur.log.Error("Failed transaction")
			tx.Rollback()
		}
	}()
	sql := "UPDATE user_tb SET nick_name = ? WHERE email = ?"
	if err := tx.Exec(sql, order.NickName, order.Email).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func (sur *SQLUserRepository) AddTeams(email string, teams []string) error {
	entries := []string{}
	arguments := []interface{}{}

	tx := sur.dataSource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			sur.log.Error("Failed transaction")
			tx.Rollback()
		}
	}()
	if teamIds, err := sur.findTeamIds(tx, teams); err != nil {
		return err
	} else if id, err := sur.findIdByEmail(tx, email); err != nil {
		return err
	} else {
		for _, teamId := range teamIds {
			entries = append(entries, "(? , ?)")
			arguments = append(arguments, id, teamId)
		}
		var sql = `INSERT INTO user_teams(user_id, team_id) VALUES %s`
		sql = fmt.Sprintf(sql, strings.Join(entries, ","))
		if err := tx.Exec(sql, arguments...).Error; err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit().Error
	}
}
func (sur *SQLUserRepository) RemoveTeams(email string, teams []string) error {
	tx := sur.dataSource.Connection().Begin()
	sql := "DELETE FROM user_teams WHERE user_id = ? AND team_id IN (?)"
	defer func() {
		if r := recover(); r != nil {
			sur.log.Error("Failed transaction")
			tx.Rollback()
		}
	}()
	if teamIds, err := sur.findTeamIds(tx, teams); err != nil {
		return err
	} else if id, err := sur.findIdByEmail(tx, email); err != nil {
		return err
	} else if err := tx.Exec(sql, id, teamIds).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func (sur *SQLUserRepository) AddPermissions(email string, permissions []string) error {
	entries := []string{}
	arguments := []interface{}{}
	tx := sur.dataSource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			sur.log.Error("Failed transaction")
			tx.Rollback()
		}
	}()
	if permissionIds, err := sur.findPermissionIds(tx, permissions); err != nil {
		return err
	} else if id, err := sur.findIdByEmail(tx, email); err != nil {
		return err
	} else {
		for _, permissionId := range permissionIds {
			entries = append(entries, "(? , ?)")
			arguments = append(arguments, id, permissionId)
		}
		var sql = `INSERT INTO user_permissions(user_id, permission_id) VALUES %s`
		sql = fmt.Sprintf(sql, strings.Join(entries, ","))
		if err := tx.Exec(sql, arguments...).Error; err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit().Error
	}
}
func (sur *SQLUserRepository) RemovePermissions(email string, permissions []string) error {
	sql := "DELETE FROM user_permissions WHERE user_id = ? AND permission_id IN (?)"
	tx := sur.dataSource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			sur.log.Error("Failed transaction")
			tx.Rollback()
		}
	}()
	defer func() {
		if r := recover(); r != nil {
			sur.log.Error("Failed transaction")
			tx.Rollback()
		}
	}()
	if permissionIds, err := sur.findPermissionIds(tx, permissions); err != nil {
		return err
	} else if id, err := sur.findIdByEmail(tx, email); err != nil {
		return err
	} else if err := tx.Exec(sql, id, permissionIds).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func (sur *SQLUserRepository) List() ([]*User, error) {
	list := make([]*User, 0)
	connection := sur.dataSource.Connection()
	if err := connection.Preload("Teams").Preload("Permissions").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
func (sur *SQLUserRepository) Remove(email string) error {
	tx := sur.dataSource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			sur.log.Error("Failed transaction")
			tx.Rollback()
		}
	}()
	if err := tx.Exec("DELETE FROM user_tb WHERE email = ?", email).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (sur *SQLUserRepository) FindById(identifier uuid.UUID) (*User, error) {
	result := &User{}
	connection := sur.dataSource.Connection()
	if err := connection.Preload("Teams").Preload("Permissions").Where("identifier = ?", identifier).First(result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (sur *SQLUserRepository) findIdByEmail(transaction *gorm.DB, email string) (int64, error) {
	var id int64
	connection := sur.dataSource.Connection()
	if err := connection.Raw("SELECT id FROM user_tb WHERE email = ?", email).Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (sur *SQLUserRepository) findPermissionIds(transaction *gorm.DB, permissions []string) ([]int64, error) {
	ids := make([]int64, 0)
	connection := sur.dataSource.Connection()
	if err := connection.Raw("SELECT id FROM permission_tb WHERE name IN (?)", permissions).Scan(&ids).Error; err != nil {
		return ids, err
	}
	return ids, nil
}

func (sur *SQLUserRepository) findTeamIds(transaction *gorm.DB, teams []string) ([]int64, error) {
	ids := make([]int64, 0)
	connection := sur.dataSource.Connection()
	if err := connection.Raw("SELECT id FROM team_tb WHERE name IN (?)", teams).Scan(&ids).Error; err != nil {
		return ids, err
	}
	return ids, nil
}
