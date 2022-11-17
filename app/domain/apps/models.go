package apps

import (
	"time"

	"github.com/space-fold-technologies/aurora-service/app/domain/clusters"
	"github.com/space-fold-technologies/aurora-service/app/domain/nodes"
	"github.com/space-fold-technologies/aurora-service/app/domain/teams"
)

type EnvVarEntry struct {
	Key string `gorm:"column:key;type:varchar"`
	Val string `gorm:"column:value;type:varchar"`
}

type ApplicationSummary struct {
	Name  string `gorm:"column:name;type:varchar"`
	Scale uint   `gorm:"column:scale;type:varchar"`
}

type Application struct {
	ID             uint               `gorm:"AUTO_INCREMENT;column:id;primary_key"`
	Identifier     string             `gorm:"column:identifier;type:varchar;not null"`
	Name           string             `gorm:"column:name;type:varchar;not null;unique"`
	Description    string             `gorm:"column:description;type:varchar;not null"`
	TeamID         uint               `gorm:"column:team_id;type:integer;not null"`
	ClusterID      uint               `gorm:"column:cluster_id;type:integer;not null"`
	Team           *teams.Team        `gorm:"foreignkey:team_id"`
	Cluster        *clusters.Cluster  `gorm:"foreignkey:cluster_id"`
	Scale          int32              `gorm:"column:scale;type:integer;not null"`
	Instances      []*nodes.Container `gorm:"foreignKey:application_id"`
	LastDeployment *time.Time         `gorm:"column:last_deployment;type:timestamp"`
}

func (Application) TableName() string {
	return "application_tb"
}

type Deployment struct {
	Identifier    string     `gorm:"column:identifier;type:uuid"`
	ApplicationID uint       `gorm:"column:application_id;type:integer;not null"`
	ImageURI      string     `gorm:"column:image_uri;type:varchar;not null"`
	AddedAt       *time.Time `gorm:"column:added_at;type:timestamp;not null"`
	CompletedAt   *time.Time `gorm:"column:completed_at;type:timestamp"`
	Status        string     `gorm:"column:status;type:varchar"`
	Report        string     `gorm:"column:report;type:varchar"`
	ServiceID     string     `gorm:"column:service_identifier;type:varchar"`
}

func (Deployment) TableName() string {
	return "deployment_tb"
}

type TargetDeployment struct {
	ImageURI string `gorm:"column:image_uri;type:varchar"`
}

type LastDeployment struct {
	ServiceID     string `gorm:"column:service_identifier;type:varchar"`
	ApplicationID string `gorm:"column:identifier;type:varchar"`
}

type DeploymentDetails struct {
	Identifier  string     `gorm:"column:identifier;type:uuid"`
	ImageURI    string     `gorm:"column:image_uri;type:varchar;not null"`
	Status      string     `gorm:"column:status;type:varchar"`
	Report      string     `gorm:"column:report;type:varchar"`
	CompletedAt *time.Time `gorm:"column:completed_at;type:timestamp"`
}

type DeploymentSummary struct {
	ImageURI string `gorm:"column:image_uri;type:varchar;not null"`
	Name     string `gorm:"column:name;type:varchar;not null"`
}
