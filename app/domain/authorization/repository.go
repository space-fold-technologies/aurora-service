package authorization

import "github.com/space-fold-technologies/aurora-service/app/domain/users"

type AuthorizationRepository interface {
	Fetch(email string) (*users.User, error)
}
