package nodes

import "github.com/space-fold-technologies/aurora-service/app/core/database"

type SQLNodeRepository struct {
	dataSource database.DataSource
}

func NewRepository(dataSource database.DataSource) NodeRepository {
	instance := new(SQLNodeRepository)
	instance.dataSource = dataSource
	return instance
}

func (nr *SQLNodeRepository) Create(entry *NodeEntry) error {
	sql := "INSERT INTO node_tb(identifier, name, type, description, address, cluster_id) " +
		"VALUES(?, ?, ?, ?, ?, (SELECT c.id FROM cluster_tb AS c WHERE c.name = ?))"
	tx := nr.dataSource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Exec(sql, entry.Identifier, entry.Name, entry.Type, entry.Description, entry.Address, entry.Cluster).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func (nr *SQLNodeRepository) Update(Name, Description string) error {
	sql := "UPDATE node_tb SET description = ? WHERE name = ?"
	tx := nr.dataSource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Exec(sql, Description, Name).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (nr *SQLNodeRepository) FetchSummary(Name string) (*NodeSummary, error) {
	sql := "SELECT identifier, type, address FROM node_tb WHERE name = ?"
	summary := &NodeSummary{}
	connection := nr.dataSource.Connection()
	if err := connection.Raw(sql, Name).First(summary).Error; err != nil {
		return nil, err
	}
	return summary, nil
}

func (nr *SQLNodeRepository) List(Cluster string) ([]*Node, error) {
	nodes := make([]*Node, 0)
	connection := nr.dataSource.Connection()
	if err := connection.
		Preload("Containers").
		Joins("JOIN cluster_tb AS ct ON node_tb.cluster_id = ct.id").
		Where("ct.name = ?", Cluster).
		Find(&nodes).Error; err != nil {
		return nodes, err
	}
	return nodes, nil
}
func (nr *SQLNodeRepository) HasContainers(Name string) (bool, error) {
	sql := "SELECT EXISTS(SELECT 1 FROM container_tb AS cn " +
		"INNER JOIN node_tb AS n ON n.cluster_id = cn.id WHERE n.name = ?) AS exist"
	connection := nr.dataSource.Connection()
	var exist int
	if err := connection.Raw(sql, Name).Find(&exist).Error; err != nil {
		return false, err
	}
	return exist > 0, nil
}

func (nr *SQLNodeRepository) Remove(Name string) error {
	sql := "DELETE FROM node_tb WHERE name = ?"
	tx := nr.dataSource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Exec(sql, Name, Name).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (nr *SQLNodeRepository) FetchClusterInfo(Cluster string) (*ClusterInfo, error) {
	info := &ClusterInfo{}
	sql := "SELECT address, token FROM cluster_tb WHERE name = ?"
	connection := nr.dataSource.Connection()
	if err := connection.Raw(sql, Cluster).First(info).Error; err != nil {
		return nil, err
	}
	return info, nil
}
