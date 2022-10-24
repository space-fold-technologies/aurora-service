package users

type UserRepository interface {
	Create(entry *UserEntry) error
	List() ([]*User, error)
	Remove(email string) error
}
