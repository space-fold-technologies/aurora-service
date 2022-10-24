package teams

import "github.com/space-fold-technologies/aurora-service/app/core/security"

type TeamService struct {
	repository TeamRepository
}

func NewService(repository TeamRepository) *TeamService {
	instance := new(TeamService)
	instance.repository = repository
	return instance
}

func (ts *TeamService) Create(order *CreateTeamOrder) error {
	return ts.repository.Create(&Team{
		Name:        order.GetName(),
		Description: order.GetDescription(),
	})
}

func (ts *TeamService) List(principals *security.Claims) (*Teams, error) {
	if result, err := ts.repository.List(); err != nil {
		return &Teams{}, nil
	} else {
		teams := &Teams{Entries: make([]*TeamEntry, 0)}
		for _, entry := range result {
			teams.Entries = append(teams.Entries, &TeamEntry{
				Name:        entry.Name,
				Description: entry.Description,
			})
		}
		return teams, nil
	}
}

func (ts *TeamService) Remove(name string) error {
	return ts.repository.Remove(name)
}
