package session

import (
	"math"
	"time"

	"github.com/alan-b-lima/almodon/internal/support/entity"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/problem"
)

const MaxRenews = math.MaxInt
const MaxAgeMax = 3 * 24 * time.Hour

type Session struct {
	UUID    uuid.UUID
	User    uuid.UUID
	Renewed int
	Expires time.Time
}

func New(user uuid.UUID, max_age time.Duration) (Session, error) {
	session := Session{
		User:    user,
		Renewed: 0,
	}

	err := problem.Join(
		entity.Set(&session.Expires, max_age, ProcessMaxAge),
	)
	if err != nil {
		return Session{}, err
	}

	session.UUID = uuid.NewUUIDv7()
	return session, nil
}

func ProcessRenewed(renewed int) (int, error) {
	renewed = max(renewed+1, 0)
	if renewed > MaxRenews {
		return 0, ErrUnrenewable
	}

	return renewed, nil
}

func ProcessMaxAge(max_age time.Duration) (time.Time, error) {
	if max_age > MaxAgeMax {
		return time.Time{}, ErrTooLong
	}

	expires := time.Now().Add(max_age)
	return expires, nil
}
