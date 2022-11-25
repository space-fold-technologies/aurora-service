package users

import (
	"errors"

	"github.com/google/uuid"
	"github.com/space-fold-technologies/aurora-service/app/core/security"
	"github.com/space-fold-technologies/aurora-service/app/domain/teams"
)

type UserService struct {
	repository      UserRepository
	passwordHandler security.PasswordHandler
}

func NewService(repository UserRepository, passwordHandler security.PasswordHandler) *UserService {
	instance := new(UserService)
	instance.repository = repository
	instance.passwordHandler = passwordHandler
	return instance
}

func (us *UserService) Create(order *CreateUserOrder) error {
	if len(order.GetPassword()) < 8 {
		return errors.New("password is to short or invalid")
	} else if hash, err := us.passwordHandler.CreateHash(order.GetPassword()); err != nil {
		return err
	} else {
		return us.repository.Create(&UserEntry{
			Identifier:   uuid.New(),
			Name:         order.GetName(),
			NickName:     order.GetNickname(),
			Email:        order.GetEmail(),
			PasswordHash: hash,
		})
	}
}

func (us *UserService) List(principal *security.Claims) (*Users, error) {
	if results, err := us.repository.List(); err != nil {
		return &Users{}, err
	} else {
		users := make([]*UserDetails, 0)
		for _, entry := range results {
			users = append(users, &UserDetails{
				Name:        entry.Name,
				Email:       entry.Email,
				Teams:       us.teamsStringList(entry.Teams),
				Permissions: us.permissionsStringList(entry.Permissions),
			})
		}
		return &Users{Users: users}, nil
	}
}

func (us *UserService) Remove(email string) error {
	return us.repository.Remove(email)
}

func (us *UserService) AddTeams(order *UpdateTeams) error {
	if len(order.GetTeams()) == 0 {
		return errors.New("no teams specified")
	}
	return us.repository.AddTeams(order.GetEmail(), order.GetTeams())
}

func (us *UserService) RemoveTeams(order *UpdateTeams) error {
	if len(order.GetTeams()) == 0 {
		return errors.New("no teams specified")
	}
	return us.repository.RemoveTeams(order.GetEmail(), order.GetTeams())
}

func (us *UserService) AddPermissions(order *UpdatePermissions) error {
	if len(order.GetPermissions()) == 0 {
		return errors.New("no permissions specified")
	}
	return us.repository.AddPermissions(order.GetEmail(), order.GetPermissions())
}

func (us *UserService) RemovePermissions(order *UpdatePermissions) error {
	if len(order.GetPermissions()) == 0 {
		return errors.New("no permissions specified")
	}
	return us.repository.RemovePermissions(order.GetEmail(), order.GetPermissions())
}

func (us *UserService) teamsStringList(list []*teams.Team) []string {
	transformed := make([]string, 0)
	for _, entry := range list {
		transformed = append(transformed, entry.Name)
	}
	return transformed
}

func (us *UserService) permissionsStringList(list []Permission) []string {
	transformed := make([]string, 0)
	for _, entry := range list {
		transformed = append(transformed, entry.Name)
	}
	return transformed
}
