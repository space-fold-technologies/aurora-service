package nodes

type NodeEntry struct {
	Name        string
	Identifier  string
	Description string
	Type        string
	Address     string
	Cluster     string
}

type NodeSummary struct {
	Identifier string `gorm:"identifier"`
	Type       string `gorm:"type"`
	Address    string `gorm:"address"`
}
