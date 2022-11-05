package users

import "github.com/google/uuid"

type UserRepository interface {
	Create(entry *UserEntry) error
	Update(order *UpdateEntry) error
	AddTeams(email string, teams []string) error
	RemoveTeams(email string, teams []string) error
	AddPermissions(email string, permissions []string) error
	RemovePermissions(email string, permissions []string) error
	List() ([]*User, error)
	Remove(email string) error
	FindById(identifier uuid.UUID) (*User, error)
}
