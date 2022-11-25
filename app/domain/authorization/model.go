package authorization

import "github.com/google/uuid"

type UserDetails struct {
	UID          uuid.UUID
	PasswordHash string
}
