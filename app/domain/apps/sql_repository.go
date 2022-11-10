package apps

import (
	"fmt"
	"strings"

	"github.com/space-fold-technologies/aurora-service/app/core/database"
)

type SQLApplicationRepository struct {
	dataSource database.DataSource
}

func NewRepository(dataSource database.DataSource) ApplicationRepository {
	instance := new(SQLApplicationRepository)
	instance.dataSource = dataSource
	return instance
}

func (sar *SQLApplicationRepository) Create(entry *ApplicationEntry) error {
	sql := "INSERT INTO application_tb(identifier, name, description, team_id, cluster_id, scale) " +
		"VALUES(?, ?, ?, (SELECT id FROM team_tb WHERE name = ?), (SELECT id FROM cluster_tb WHERE name = ?), ?)"
	tx := sar.dataSource.Connection().Begin()
	if err := tx.Exec(sql, entry.Identifier, entry.Name, entry.Description, entry.Team, entry.Cluster, entry.Scale).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (sar *SQLApplicationRepository) Update(Name, Description string, Scale int32) error {
	sql := "UPDATE application_tb SET description = ?, scale = ? WHERE name = ?"
	tx := sar.dataSource.Connection().Begin()
	if err := tx.Exec(sql, Description, Scale, Name).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (sar *SQLApplicationRepository) Remove(Name string) error {
	sql := "DELETE FROM application_tb WHERE name = ?"
	tx := sar.dataSource.Connection().Begin()
	if err := tx.Exec(sql, Name).Error; err != nil {
		return tx.Rollback().Error
	}
	return tx.Commit().Error
}
func (sar *SQLApplicationRepository) List(Cluster string) ([]*ApplicationSummary, error) {
	apps := make([]*ApplicationSummary, 0)
	sql := "SELECT a.name, a.scale " +
		"FROM application_tb AS a " +
		"INNER JOIN cluster_tb AS c ON a.cluster_id = c.id " +
		"WHERE c.name = ?"

	if err := sar.dataSource.Connection().Raw(sql, Cluster).Find(&apps).Error; err != nil {
		return apps, err
	}
	return apps, nil
}

func (sar *SQLApplicationRepository) FetchDetails(Name string) (*Application, error) {
	application := &Application{}
	conn := sar.dataSource.Connection()
	if err := conn.Preload("Team").Preload("Cluster").Preload("Instances").Where("name = ?", Name).First(application).Error; err != nil {
		return application, err
	}
	return application, nil
}

func (sar *SQLApplicationRepository) RegisterDeploymentEntry(order *DeploymentOrder) error {
	sql := "INSERT INTO deployment_tb(identifier, application_id, image_uri, status, added_at) " +
		"VALUES(?, (SELECT id FROM application_tb WHERE name = ?), ?, ?, ?)"
	tx := sar.dataSource.Connection().Begin()
	if err := tx.Exec(sql, order.Identifier, order.ApplicationName, order.ImageURI, order.Status, order.CreatedAt).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (sar *SQLApplicationRepository) UpdateDeploymentEntry(order *DeploymentUpdate) error {
	sql := "UPDATE deployment_tb SET status = ?, report = ?, completed_at = ?, service_identifier = ? WHERE identifier = ?"
	tx := sar.dataSource.Connection().Begin()

	if err := tx.Exec(sql, order.Status, order.Report, order.UpdatedAt, order.ServiceID, order.Identifier).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (sar *SQLApplicationRepository) FetchImageURI(identifier string) (string, error) {
	sql := "SELECT image_uri FROM deployment_tb WHERE identifier = ?"
	var imageUri string
	conn := sar.dataSource.Connection()
	if err := conn.Raw(sql, identifier).First(&imageUri).Error; err != nil {
		return "", err
	}
	return imageUri, nil
}
func (sar *SQLApplicationRepository) FetchEnvVars(name string) ([]*EnvVarEntry, error) {
	sql := "SELECT e.key, e.value FROM environment_variable_tb AS e, " +
		"(SELECT a.name AS application_name, c.name AS cluster_name, n.name AS node_name FROM application_tb AS a " +
		"INNER JOIN cluster_tb AS c ON a.cluster_id = c.id " +
		"INNER JOIN node_tb AS n ON n.cluster_id = c.id WHERE a.name = ?) AS s " +
		"WHERE (e.scope = ? AND e.target = s.application_name) OR " +
		"(e.scope = ? AND e.target = s.node_name) OR " +
		"(e.scope = ? AND e.target = s.cluster_name)"
	vars := make([]*EnvVarEntry, 0)
	conn := sar.dataSource.Connection()
	if err := conn.Raw(sql, name, "APPLICATION", "NODE", "CLUSTER").Find(&vars).Error; err != nil {
		return vars, err
	}
	return vars, nil
}
func (sar *SQLApplicationRepository) AddContainers(orders []*ContainerOrder) error {
	sql := `INSERT INTO container_tb(identifier, ip, address_family, application_id, node_id) VALUES %s`
	entries := []string{}
	parameters := []interface{}{}
	for _, order := range orders {
		entries = append(entries, "(?, ?, ?, (SELECT application_id FROM deployment_tb WHERE service_identifier = ?), (SELECT id FROM node_tb WHERE identifier = ?))")
		parameters = append(parameters, order.Identifier, order.IP, order.Family, order.ServiceID, order.Node)
	}

	sql = fmt.Sprintf(sql, strings.Join(entries, ","))
	tx := sar.dataSource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Exec(sql, parameters...).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (sar *SQLApplicationRepository) RemoveContainer(Identifier string) error {
	sql := "DELETE FROM container_tb WHERE identifier = ?"
	tx := sar.dataSource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Exec(sql, Identifier).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (sar *SQLApplicationRepository) Deployed(name string) (*LastDeployment, error) {
	sql := "SELECT a.identifier, d.service_identifier FROM deployment_tb AS d " +
		"INNER JOIN application_tb AS a ON d.application_id = a.id " +
		"WHERE a.name = ? AND a.last_deployment = d.completed_at"
	deployment := &LastDeployment{}
	connection := sar.dataSource.Connection()
	if err := connection.Raw(sql, connection).First(deployment).Error; err != nil {
		return nil, err
	}
	return deployment, nil
}

func (sar *SQLApplicationRepository) RemoveContainers(applicationId string) error {
	sql := "DELETE FROM container_tb WHERE " +
		"application_id = (SELECT a.id FROM application_tb AS a WHERE a.identifier = ?"
	tx := sar.dataSource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Exec(sql, applicationId).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
