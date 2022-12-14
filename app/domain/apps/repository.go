package apps

import "time"

type ApplicationRepository interface {
	Create(entry *ApplicationEntry) error
	Update(Name, Description string, Scale int32) error
	Remove(Name string) error
	List(Cluster string) ([]*ApplicationSummary, error)
	FetchDetails(Name string) (*Application, error)
	RegisterDeploymentEntry(order *DeploymentOrder) error
	UpdateDeploymentEntry(order *DeploymentUpdate) error
	FetchImageURI(identifier string) (string, error)
	FetchEnvVars(name string) ([]*EnvVarEntry, error)
	AddContainers(order []*ContainerOrder) error
	FetchContainer(identifier string) (*ContainerDetails, error)
	RemoveContainer(Identifier string) error
	Deployed(name string) (*LastDeployment, error)
	RemoveContainers(applicationId string) error
	Deployments(name string) ([]*DeploymentDetails, error)
	FetchDeployment(identifier string) (*DeploymentSummary, error)
	FetchActiveDeployments() ([]*ServiceCheck, error)
	RemoveContainersOlderThan(targetTime *time.Time) error
}
