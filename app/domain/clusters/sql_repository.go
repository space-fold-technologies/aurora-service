package clusters

import (
	"fmt"
	"strings"

	"github.com/space-fold-technologies/aurora-service/app/core/database"
	"github.com/space-fold-technologies/aurora-service/app/domain/teams"
	"gorm.io/gorm"
)

type SQLClusterRepository struct {
	dataSource database.DataSource
}

func NewRepository(dataSource database.DataSource) ClusterRepository {
	instance := new(SQLClusterRepository)
	instance.dataSource = dataSource
	return instance
}

func (scr *SQLClusterRepository) Create(entry *ClusterEntry) error {
	tx := scr.dataSource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	cluster := &Cluster{
		Identifier:  entry.Identifier,
		Name:        entry.Name,
		Description: entry.Description,
		Type:        entry.Type,
		Address:     entry.Address,
		Namespace:   entry.Namespace,
		Teams:       nil,
		Nodes:       nil,
	}
	teams := make([]*teams.Team, 0)
	if err := tx.Where("name IN (?)", entry.Teams).Find(&teams).Error; err != nil {
		return err
	} else if len(teams) == 0 {
		return fmt.Errorf("no teams were found to match")
	} else if err := tx.Create(cluster).Error; err != nil {
		tx.Rollback()
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return fmt.Errorf("record with duplicate unique value found")
		}
		return err
	} else if err := tx.Model(cluster).Association("Teams").Append(teams); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
func (scr *SQLClusterRepository) Update(Name, Description string, Teams []string) error {
	tx := scr.dataSource.Connection().Begin()
	if cluster, err := scr.fetch(tx, Name); err != nil {
		return err
	} else {
		if cluster.Description != Description {
			if err := tx.Model(cluster).Updates(map[string]interface{}{"description": Description}).Error; err != nil {
				tx.Rollback()
				if strings.Contains(strings.ToLower(err.Error()), "unique") {
					return fmt.Errorf("record with duplicate unique value found")
				}
				return err
			}
		}
		if len(Teams) > 0 {
			teams := make([]*teams.Team, 0)
			if err := tx.Where("name IN (?)", Teams).Find(teams).Error; err != nil {
				tx.Rollback()
				return err
			} else if err := tx.Model(cluster).Association("Teams").Append(teams); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}
func (scr *SQLClusterRepository) List(teams []string) ([]*Cluster, error) {
	clusters := make([]*Cluster, 0)
	connection := scr.dataSource.Connection()
	if err := connection.Preload("Teams").
		Preload("Nodes").
		Joins("JOIN cluster_teams AS ct ON ct.cluster_id = cluster_tb.id").
		Joins("JOIN team_tb AS t ON t.id = ct.team_id").
		Where("t.name IN(?)", teams).
		Find(&clusters).Error; err != nil {
		return clusters, err
	}
	return clusters, nil
}
func (scr *SQLClusterRepository) HasNodes(name string) (bool, error) {
	sql := "SELECT EXISTS(SELECT 1 FROM node_tb AS n INNER JOIN cluster_tb AS c ON n.cluster_id = c.id WHERE c.name = ?) AS exist"
	var exist int64
	connection := scr.dataSource.Connection()
	if err := connection.Raw(sql, name).Scan(&exist).Error; err != nil {
		return false, err
	}
	return exist > 0, nil
}
func (scr *SQLClusterRepository) Remove(name string) error {
	sql := "DELETE FROM cluster_tb WHERE " +
		"NOT EXISTS(SELECT 1 FROM node_tb AS n INNER JOIN cluster_tb AS c ON n.cluster_id = c.id WHERE c.name = ?) " +
		"AND name = ?"
	tx := scr.dataSource.Connection().Begin()
	if err := tx.Exec(sql, name, name).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (scr *SQLClusterRepository) fetch(transaction *gorm.DB, name string) (*Cluster, error) {
	cluster := &Cluster{}
	if err := transaction.Where("name = ?", name).First(cluster).Error; err != nil {
		return cluster, err
	}
	return cluster, nil
}
