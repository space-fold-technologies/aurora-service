package users

import "github.com/google/uuid"

type UserEntry struct {
	Name         string
	NickName     string
	Identifier   uuid.UUID
	Email        string
	PasswordHash string
}

type UpdateEntry struct {
	Email    string
	NickName string
}
