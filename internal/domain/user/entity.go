package user

import (
	"errors"
	"regexp"
	"slices"
	"unicode/utf8"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/support/entity"
	"github.com/alan-b-lima/almodon/pkg/uuid"

	"github.com/alan-b-lima/pkg/problem"

	"golang.org/x/crypto/bcrypt"
)

const (
	PasswordMinLen = 8
	PasswordMaxLen = 64
)

var (
	re_siape  = regexp.MustCompile(`^\d{7}$`)
	re_email = regexp.MustCompile(`^[0-9A-Za-z_%+-]+(\.[0-9A-Za-z_%+-]+)*@[0-9A-Za-z-]+(\.[0-9A-Za-zA-Z-]+)*\.[A-Za-z]{2,}$`)

	accept_roles = [...]auth.Role{auth.User, auth.Admin, auth.Chief, auth.Maintainer}
)

type User struct {
	UUID     uuid.UUID
	SIAPE    string
	Name     string
	Email    string
	Password []byte
	Role     auth.Role
}

func New(siape, name, email, password string, role auth.Role) (User, error) {
	var u User

	errpwd := entity.Set(&u.Password, password, ProcessPassword)
	if err, ok := errors.AsType[*problem.Error](errpwd); ok && err.IsInternal() {
		return User{}, err
	}

	err := problem.Join(
		entity.Set(&u.SIAPE, siape, ProcessSIAPE),
		entity.Set(&u.Name, name, ProcessName),
		entity.Set(&u.Email, email, ProcessEmail),
		errpwd,
		entity.Set(&u.Role, role, ProcessRole),
	)
	if err != nil {
		return User{}, ErrCreate.Cause(err).Make()
	}

	u.UUID = uuid.NewUUIDv7()
	return u, nil
}

func ComparePassword(hash []byte, password string) error {
	if bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil {
		return nil
	}

	return ErrPasswordIncorrect
}

func ProcessSIAPE(siape string) (string, error) {
	if !re_siape.MatchString(siape) {
		return "", ErrSiapeInvalid
	}

	return siape, nil
}

func ProcessName(name string) (string, error) {
	if name == "" {
		return "", ErrNameEmpty
	}

	return name, nil
}

func ProcessEmail(email string) (string, error) {
	if !re_email.MatchString(email) {
		return "", ErrEmailInvalid
	}

	return email, nil
}

func ProcessPassword(password string) ([]byte, error) {
	if len(password) < PasswordMinLen {
		return nil, ErrPasswordTooShort
	}

	if len(password) > PasswordMaxLen {
		return nil, ErrPasswordTooLong
	}

	switch password[0] {
	case ' ', '\t', '\n', '\r':
		return nil, ErrPasswordLeadOrTrailWhitespace
	}

	switch password[len(password)-1] {
	case ' ', '\t', '\n', '\r':
		return nil, ErrPasswordLeadOrTrailWhitespace
	}

	for _, rune := range password {
		if rune < ' ' || !utf8.ValidRune(rune) {
			return nil, ErrPasswordIllegalCharacters
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrPasswordFailedToHash.Cause(err).Make()
	}

	return hash, nil
}

func ProcessRole(role auth.Role) (auth.Role, error) {
	if !slices.Contains(accept_roles[:], role) {
		return 0, ErrRoleInvalid
	}

	return role, nil
}
