package users

import "github.com/space-fold-technologies/aurora-service/app/core/security"

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
	if hash, err := us.passwordHandler.CreateHash(order.Password); err != nil {
		return err
	} else {
		return us.repository.Create(&UserEntry{
			Name:     order.Name,
			Email:    order.Email,
			Password: hash,
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
				Name:  entry.Name,
				Email: entry.Email,
				Teams: entry.Teams,
			})
		}
		return &Users{Users: users}, nil
	}
}

func (us *UserService) Remove(email string) error {
	return us.repository.Remove(email)
}
