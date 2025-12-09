package auth

import "github.com/alan-b-lima/almodon/pkg/uuid"

type Service interface {
	Login(siape string, password string) (Result, error)
	Logout(session uuid.UUID) error

	Actor(session uuid.UUID) (Actor, error)
}

