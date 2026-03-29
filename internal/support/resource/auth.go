package resource

import (
	"context"
	"net/http"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

const SessionCookieName = "session"

func Session(ctx context.Context, r *http.Request) (context.Context, error) {
	session, err := SessionCookie(r)
	if err != nil {
		if err == http.ErrNoCookie {
			return ctx, nil
		}

		return nil, auth.ErrUnauthenticated.Cause(err).Make()
	}

	return context.WithValue(ctx, "session", session), nil
}

func SessionCookie(r *http.Request) (uuid.UUID, error) {
	s, err := r.Cookie(SessionCookieName)
	if err != nil {
		return uuid.UUID{}, http.ErrNoCookie
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
