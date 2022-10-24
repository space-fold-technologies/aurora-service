package clusters

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/space-fold-technologies/aurora-service/app/core/providers"
	"github.com/space-fold-technologies/aurora-service/app/core/security"
	"github.com/space-fold-technologies/aurora-service/app/domain/nodes"
	"github.com/space-fold-technologies/aurora-service/app/domain/teams"
)

type ClusterService struct {
	repository ClusterRepository
	provider   providers.Provider
}

func NewService(provider providers.Provider, repository ClusterRepository) *ClusterService {
	instance := new(ClusterService)
	instance.repository = repository
	instance.provider = provider
	return instance
}

func (cs *ClusterService) Create(order *CreateClusterOrder) error {
	return cs.registeCluster(order.GetName(), order.GetAddress(), order.GetType(), func(token string) error {
		return cs.repository.Create(&ClusterEntry{
			Identifier:  uuid.NewString(),
			Name:        order.GetName(),
			Description: order.GetDescription(),
			Type:        order.GetType().String(),
			Address:     order.GetAddress(),
			Token:       token,
			Teams:       order.GetTeams(),
		})
	})
}

func (cs *ClusterService) Update(order *UpdateClusterOrder) error {
	return cs.repository.Update(order.GetName(), order.GetDescription(), order.GetTeams())
}

func (cs *ClusterService) List(principals *security.Claims) (*Clusters, error) {
	if entries, err := cs.repository.List(principals.Teams); err != nil {
		return nil, err
	} else {
		pack := &Clusters{Entries: make([]*ClusterDetails, 0)}
		for _, entry := range entries {
			pack.Entries = append(pack.Entries, &ClusterDetails{
				Name:      entry.Name,
				Address:   entry.Address,
				Type:      entry.Type,
				Teams:     cs.teamTransform(entry.Teams),
				Namespace: entry.Namespace,
				Nodes:     cs.nodeTransformer(entry.Nodes),
			})
		}
		return pack, nil
	}
}

func (cs *ClusterService) Remove(name string) error {
	if hasNode, err := cs.repository.HasNodes(name); err != nil {
		return err
	} else if hasNode {
		return fmt.Errorf("cannot destroy a cluster with nodes")
	}
	return cs.repository.Remove(name)
}

func (cs *ClusterService) nodeTransformer(entries []*nodes.Node) []*NodeDetails {
	nodes := make([]*NodeDetails, 0)
	for _, entry := range entries {
		nodes = append(nodes, &NodeDetails{
			Name:    entry.Name,
			Type:    entry.Type,
			Address: entry.Address,
		})
	}
	return nodes
}

func (cs *ClusterService) teamTransform(entries []*teams.Team) []string {
	teams := make([]string, 0)
	for _, entry := range entries {
		teams = append(teams, entry.Name)
	}
	return teams
}

func (cs *ClusterService) registeCluster(name, address string, clusterType CreateClusterOrder_Type, callback ClusterCreationCallback) error {
	// This registers a cluster with an implementation of the provider [KUBERNETES, DOCKER-SWARM]
	if clusterType == CreateClusterOrder_KUBERNETES_CLUSTER {
		return errors.New("hi, we are not supporting kubernetes at this time :| sorry")
	} else if clusterType == CreateClusterOrder_DOCKER_SWARM {
		if token, err := cs.provider.Initialize(); err != nil {
			return err
		} else {
			return callback(token)
		}
	}
	return nil
}
