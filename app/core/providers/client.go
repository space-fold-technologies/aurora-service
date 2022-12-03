package providers

import "context"

type AgentClient interface {
	Join(ctx context.Context, order *RegisterAgent, address, token string) error
	Leave(ctx context.Context, order *RemoveAgent, address, token string) error
	Containers(ctx context.Context, serviceId, address string) (*ContainerReport, error)
}
