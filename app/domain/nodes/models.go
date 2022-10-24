package nodes

type Container struct {
	ID            uint   `gorm:"AUTO_INCREMENT;column:id;primary_key"`
	Identifier    string `gorm:"column:identifier;type:varchar;not null"`
	IP            string `gorm:"column:ip;type:varchar;not null"`
	Family        int32  `gorm:"column:address_family;type:integer;not null"`
	Node          uint   `gorm:"column:node_id;type:integer"`
	ApplicationID uint   `gorm:"column:application_id;integer;not null"`
}

func (Container) TableName() string {
	return "container_tb"
}

type Node struct {
	ID          uint         `gorm:"AUTO_INCREMENT;column:id;primary_key"`
	Identifier  string       `gorm:"identifier;type:varchar;not null"`
	Type        string       `gorm:"type;type:varchar;not null"`
	Name        string       `gorm:"name;type:varchar;not null;unique"`
	Description string       `gorm:"description;type:varchar;not null"`
	Address     string       `gorm:"address;type:varchar;not null;unique"`
	ClusterID   uint         `gorm:"cluster_id;type:integer;not null"`
	Containers  []*Container `gorm:"foreignKey:application_id"`
}

func (Node) TableName() string {
	return "node_tb"
}

type ClusterInfo struct {
	Address string `gorm:"address"`
	Token   string `gorm:"token"`
}
