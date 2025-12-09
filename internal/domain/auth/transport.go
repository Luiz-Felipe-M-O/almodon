package auth

import (
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Request struct {
	SIAPE    string `json:"siape"`
	Password string `json:"password"`
}

type Result struct {
	UUID    uuid.UUID `json:"uuid"`
	User    uuid.UUID `json:"user"`
	Expires time.Time `json:"expires"`
}
