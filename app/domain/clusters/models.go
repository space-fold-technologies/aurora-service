package clusters

import (
	"github.com/space-fold-technologies/aurora-service/app/domain/nodes"
	"github.com/space-fold-technologies/aurora-service/app/domain/teams"
)

type Cluster struct {
	ID          uint          `gorm:"AUTO_INCREMENT;column:id;primary_key"`
	Identifier  string        `gorm:"column:identifier;type:varchar;not null"`
	Name        string        `gorm:"column:name;type:varchar;not null;unique"`
	Description string        `gorm:"column:description;type:varchar;not null"`
	Token       string        `gorm:"column:token;type:varchar;not null"`
	Type        string        `gorm:"column:type;type:varchar;not null"`           // KUBERNETES|DOCKER-SWARM|PODMAN-SWARM
	Address     string        `gorm:"column:address;type:varchar;not null;unique"` // CONTROL NODE <<http://>> might consider grpc or other impl
	Namespace   string        `gorm:"column:namespace;type:varchar;not null"`      // Only required for kubernetes namespace access
	Teams       []*teams.Team `gorm:"many2many:cluster_teams"`                     // team owning the cluster
	Nodes       []*nodes.Node `gorm:"foreignKey:cluster_id"`
}

func (Cluster) TableName() string {
	return "cluster_tb"
}
