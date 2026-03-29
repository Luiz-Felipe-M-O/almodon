package promotion

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/support/entity"
	"github.com/alan-b-lima/almodon/pkg/uuid"

	"github.com/alan-b-lima/pkg/problem"
)

const MaxAgeMax = 3 * 24 * time.Hour

type Promotion struct {
	UUID    uuid.UUID
	User    uuid.UUID
	Expires time.Time
}

func New(user uuid.UUID, max_age time.Duration) (Promotion, error) {
	session := Promotion{
		User: user,
	}

	err := problem.Join(
		entity.Set(&session.Expires, max_age, ProcessMaxAge),
	)
	if err != nil {
		return Promotion{}, err
	}

	session.UUID = uuid.NewUUIDv7()
	return session, nil
}

func ProcessMaxAge(max_age time.Duration) (time.Time, error) {
	if max_age > MaxAgeMax {
		return time.Time{}, ErrTooLong
	}

	return time.Now().Add(max_age), nil
}
