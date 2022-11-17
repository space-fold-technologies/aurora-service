package nodes

import (
	"fmt"

	"github.com/space-fold-technologies/aurora-service/app/core/providers"
)

type NodeService struct {
	repository NodeRepository
	provider   providers.Provider
}

func NewService(provider providers.Provider, repository NodeRepository) *NodeService {
	instance := new(NodeService)
	instance.provider = provider
	instance.repository = repository
	return instance
}

func (ns *NodeService) Create(order *CreateNodeOrder, token string) error {
	if info, err := ns.repository.FetchClusterInfo(order.GetCluster()); err != nil {
		return err
	} else if details, err := ns.provider.Join(&providers.JoinOrder{
		Name:           order.Name,
		CaptainAddress: info.Address,
		WorkerAddress:  order.GetAddress(),
		Token:          token,
	}); err != nil {
		return err
	} else {
		return ns.repository.Create(&NodeEntry{
			Name:        order.GetName(),
			Identifier:  details.ID,
			Type:        order.GetType(),
			Description: order.GetDescription(),
			Address:     order.GetAddress(),
			Cluster:     order.GetCluster(),
		})
	}

}

func (ns *NodeService) Update(order *UpdateNodeOrder) error {
	return ns.repository.Update(order.GetName(), order.GetDescription())
}

func (ns *NodeService) List(cluster string) (*Nodes, error) {
	if results, err := ns.repository.List(cluster); err != nil {
		return nil, err
	} else {
		nodes := &Nodes{Entries: make([]*Entry, 0)}
		for _, entry := range results {
			nodes.Type = entry.Type
			nodes.Cluster = cluster
			nodes.Entries = append(nodes.Entries, &Entry{
				Name:        entry.Name,
				Description: entry.Description,
				Address:     entry.Address,
				Containers:  int32(len(entry.Containers)),
			})
		}
		return nodes, nil
	}
}

func (ns *NodeService) Remove(name, token string) error {
	if hasContainers, err := ns.repository.HasContainers(name); err != nil {
		return err
	} else if hasContainers {
		return fmt.Errorf("cannot remove node with running containers")
	} else if node, err := ns.repository.FetchSummary(name); err != nil {
		return err
	} else if err := ns.provider.Leave(&providers.LeaveOrder{
		NodeID:  node.Identifier,
		Address: node.Address,
		Token:   token}); err != nil {
		return err
	}
	return ns.repository.Remove(name)
}
