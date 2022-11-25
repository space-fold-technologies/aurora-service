package clusters

type ClusterEntry struct {
	Identifier  string
	Name        string
	Description string
	Type        string
	Address     string
	Namespace   string
	Teams       []string
}

type ClusterCreationCallback func(token string) error
