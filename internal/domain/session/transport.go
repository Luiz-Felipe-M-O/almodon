package session

import (
	"time"

	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type (
	Create struct {
		User   uuid.UUID              `json:"user"`
		MaxAge opt.Opt[time.Duration] `json:"max_age"`
	}

	Update struct {
		MaxAge opt.Opt[time.Duration] `json:"max_age"`
	}
)

type (
	Response struct {
		User    uuid.UUID `json:"user"`
		Expires time.Time `json:"expires"`
	}
)
