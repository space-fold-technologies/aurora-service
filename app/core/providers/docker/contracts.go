package docker

type ProgressChunk struct {
	Status string `json:"status"`
}

type ContainerDetails struct {
	ID            string
	NodeID        string
	TaskID        string
	ServiceID     string
	IPAddress     string
	AddressFamily uint
}
