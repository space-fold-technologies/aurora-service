package environments

import "github.com/space-fold-technologies/aurora-service/app/core/security"

type EnvironmentService struct {
	repository EnvironmentRepository
}

func NewService(repository EnvironmentRepository) *EnvironmentService {
	instance := new(EnvironmentService)
	instance.repository = repository
	return instance
}

func (es *EnvironmentService) Create(order *CreateEnvEntryOrder) error {
	return es.repository.Add(
		order.GetScope().String(),
		order.GetTarget(),
		es.from(order.GetEntries()),
	)
}

func (es *EnvironmentService) Remove(order *RemoveEnvEntryOrder) error {
	return es.repository.Remove(order.GetKeys())
}

func (es *EnvironmentService) List(principals *security.Claims, scope, target string) (*EnvSet, error) {
	if list, err := es.repository.List(scope, target); err != nil {
		return nil, err
	} else {
		entries := make([]*EnvEntry, 0)
		for _, entry := range list {
			entries = append(entries, &EnvEntry{Key: entry.Key, Value: entry.Value})
		}
		return &EnvSet{Entries: entries}, nil
	}
}
func (es *EnvironmentService) from(entries []*EnvEntry) []*Entry {
	vars := make([]*Entry, 0)
	for _, entry := range entries {
		vars = append(vars, &Entry{
			Key:   entry.GetKey(),
			Value: entry.GetValue(),
		})
	}
	return vars
}
