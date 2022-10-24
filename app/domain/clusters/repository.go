package clusters

type ClusterRepository interface {
	Create(entry *ClusterEntry) error
	Update(Name, Description string, Teams []string) error
	List(teams []string) ([]*Cluster, error)
	HasNodes(name string) (bool, error)
	Remove(name string) error
}
