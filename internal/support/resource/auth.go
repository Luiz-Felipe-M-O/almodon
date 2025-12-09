package resource

import (
	"net/http"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

const SessionCookieName = "session"

func Session(rc auth.Identifier, r *http.Request) (auth.Actor, error) {
	session, err := SessionCookie(r)
	if err != nil {
		return auth.NewUnlogged(), nil
	}

	actor, err := rc.Actor(session)
	if err, ok := errors.AsType[*errors.Error](err); ok && err.Kind.IsClient() {
		return auth.NewUnlogged(), nil
	}
	if err != nil {
		return auth.NewUnlogged(), err
	}

	return actor, err
}

func SessionCookie(r *http.Request) (uuid.UUID, error) {
	s, err := r.Cookie(SessionCookieName)
	if err != nil {
		return uuid.UUID{}, auth.ErrUnauthenticated.Make()
	}

	session, err := uuid.FromString(s.Value)
	if err != nil {
		return uuid.UUID{}, ErrBadUUID
	}

	return session, nil
}

func SetSessionCookie(w http.ResponseWriter, uuid uuid.UUID, expires time.Time) {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    uuid.String(),
		Expires:  expires,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
}

func DeleteSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
}
