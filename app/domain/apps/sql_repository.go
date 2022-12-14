package apps

import (
	"fmt"
	"strings"
	"time"

	"github.com/space-fold-technologies/aurora-service/app/core/database"
	"gorm.io/gorm"
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
	sql := "UPDATE deployment_tb SET status = ?, report = ?, completed_at = ?, image_uri = ?, service_identifier = ? WHERE identifier = ?"
	tx := sar.dataSource.Connection().Begin()

	if err := tx.Exec(sql, order.Status, order.Report, order.UpdatedAt, order.ImageURI, order.ServiceID, order.Identifier).Error; err != nil {
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
	sql := `INSERT OR REPLACE INTO container_tb(identifier, ip, address_family, application_id, node_id, last_updated_at) VALUES %s`
	entries := []string{}
	parameters := []interface{}{}
	checkedAt := time.Now()
	for _, order := range orders {
		entries = append(entries, "(?, ?, ?, (SELECT application_id FROM deployment_tb WHERE service_identifier = ?), (SELECT id FROM node_tb WHERE identifier = ?), ?)")
		parameters = append(parameters, order.Identifier, order.IP, order.Family, order.ServiceID, order.Node, checkedAt)
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

func (sar *SQLApplicationRepository) FetchContainer(identifier string) (*ContainerDetails, error) {
	details := &ContainerDetails{}
	sql := "SELECT n.address AS node_address, c.identifier AS container_identifier FROM container_tb AS c " +
		"INNER JOIN node_tb AS n ON c.node_id = n.id " +
		"WHERE c.identifier = ?"
	connection := sar.dataSource.Connection()
	if err := connection.Raw(sql, identifier).First(details).Error; err != nil {
		return nil, err
	}
	return details, nil
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
		"WHERE a.name = ? AND d.status = ? AND d.completed_at = a.last_deployment"
	deployment := &LastDeployment{}
	connection := sar.dataSource.Connection()
	if err := connection.Raw(sql, name, "DEPLOYED").First(deployment).Error; err != nil {
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

func (sar *SQLApplicationRepository) Deployments(name string) ([]*DeploymentDetails, error) {
	connection := sar.dataSource.Connection()
	deployments := make([]*DeploymentDetails, 0)
	sql := "SELECT d.identifier, d.image_uri, d.status, d.report, d.completed_at " +
		"FROM deployment_tb AS d " +
		"INNER JOIN application_tb AS a ON a.id = d.application_id " +
		"WHERE a.name = ? " +
		"ORDER BY d.completed_at ASC"
	if err := connection.Raw(sql, name).Find(&deployments).Error; err != nil {
		return nil, err
	}
	return deployments, nil
}

func (sar *SQLApplicationRepository) FetchDeployment(identifier string) (*DeploymentSummary, error) {
	connection := sar.dataSource.Connection()
	summary := &DeploymentSummary{}
	sql := "SELECT d.image_uri, a.name FROM deployment_tb AS d " +
		"INNER JOIN application_tb AS a ON a.id = d.application_id " +
		"WHERE d.identifier = ?"
	if err := connection.Raw(sql, identifier).First(summary).Error; err != nil {
		return nil, err
	}
	return summary, nil
}

func (sar *SQLApplicationRepository) FetchActiveDeployments() ([]*ServiceCheck, error) {
	connection := sar.dataSource.Connection()
	checks := make([]*ServiceCheck, 0)
	sql := "SELECT d.service_identifier FROM deployment_tb AS d " +
		"INNER JOIN application_tb AS a ON d.application_id = a.id " +
		"WHERE d.status = ? AND d.completed_at = a.last_deployment"
	if err := connection.Raw(sql, "DEPLOYED").Find(&checks).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return checks, nil
}

func (sar *SQLApplicationRepository) RemoveContainersOlderThan(targetTime *time.Time) error {
	sql := "DELETE FROM container_tb WHERE last_updated_at < ?"
	tx := sar.dataSource.Connection().Begin()
	if err := tx.Exec(sql, targetTime).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
