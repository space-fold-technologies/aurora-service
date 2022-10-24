package nodes

type NodeRepository interface {
	Create(entry *NodeEntry) error
	Update(Name, Description string) error
	FetchSummary(Name string) (*NodeSummary, error)
	List(Cluster string) ([]*Node, error)
	HasContainers(Name string) (bool, error)
	Remove(Name string) error
	FetchClusterInfo(Cluster string) (*ClusterInfo, error)
}
