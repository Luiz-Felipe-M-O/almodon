package user

import (
	"regexp"
	"slices"
	"unicode/utf8"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/pkg/errors"
	hashpkg "github.com/alan-b-lima/almodon/pkg/hash"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

const (
	PasswordMinLen = 64
	PasswordMaxLen = 64
)

var (
	reEmail     = regexp.MustCompile(`^[0-9A-Za-z_%+-]+(\.[0-9A-Za-z_%+-]+)*@[0-9A-Za-z-]+(\.[0-9A-Za-zA-Z-]+)*\.[A-Za-z]{2,}$`)
	acceptRoles = [...]auth.Role{auth.User, auth.Admin, auth.Chief}
)

type User struct {
	uuid     uuid.UUID
	siape    string
	name     string
	email    string
	password [60]byte
	role     auth.Role
}

func New(siape, name, email, password string, role auth.Role) (User, error) {
	var u User

	errpwd := u.SetPassword(password)
	if err, ok := errors.AsType[*errors.Error](errpwd); ok && err.IsInternal() {
		return User{}, err
	}

	err := errors.Join(
		u.SetSIAPE(siape),
		u.SetName(name),
		u.SetEmail(email),
		errpwd,
		u.SetRole(role),
	)
	if err != nil {
		return User{}, ErrUserCreate.Cause(err).Make()
	}

	u.uuid = uuid.NewUUIDv7()
	return u, nil
}

func (u *User) UUID() uuid.UUID    { return u.uuid }
func (u *User) SIAPE() string      { return u.siape }
func (u *User) Name() string       { return u.name }
func (u *User) Email() string      { return u.email }
func (u *User) Password() [60]byte { return u.password }
func (u *User) Role() auth.Role    { return u.role }

func (u *User) SetSIAPE(siape string) error       { return set(&u.siape, siape, ProcessSiape) }
func (u *User) SetName(name string) error         { return set(&u.name, name, ProcessName) }
func (u *User) SetEmail(email string) error       { return set(&u.email, email, ProcessEmail) }
func (u *User) SetPassword(password string) error { return set(&u.password, password, ProcessPassword) }
func (u *User) SetRole(role auth.Role) error      { return set(&u.role, role, ProcessRole) }

func ComparePassword(hash, password []byte) error {
	if hashpkg.Compare(hash, password) {
		return nil
	}

	return ErrPasswordIncorrect
}

func ProcessSiape(siape string) (string, error) {
	return siape, nil
}

func ProcessName(name string) (string, error) {
	if name == "" {
		return "", ErrNameEmpty
	}

	return name, nil
}

func ProcessEmail(email string) (string, error) {
	if !reEmail.MatchString(email) {
		return "", ErrEmailInvalid
	}

	return email, nil
}

func ProcessPassword(password string) ([60]byte, error) {
	if len(password) < PasswordMinLen {
		return [60]byte{}, ErrPasswordTooShort
	}

	if len(password) > PasswordMaxLen {
		return [60]byte{}, ErrPasswordTooLong
	}

	switch password[0] {
	case ' ', '\t', '\n', '\r':
		return [60]byte{}, ErrPasswordLeadOrTrailWhitespace
	}

	switch password[len(password)-1] {
	case ' ', '\t', '\n', '\r':
		return [60]byte{}, ErrPasswordLeadOrTrailWhitespace
	}

	for _, rune := range password {
		if rune < ' ' || !utf8.ValidRune(rune) {
			return [60]byte{}, ErrPasswordIllegalCharacters
		}
	}

	hash, err := hashpkg.Hash([]byte(password))
	if err != nil {
		return [60]byte{}, ErrPasswordFailedToHash.Cause(err).Make()
	}

	return hash, nil
}

func ProcessRole(role auth.Role) (auth.Role, error) {
	if !slices.Contains(acceptRoles[:], role) {
		return 0, ErrRoleInvalid
	}

	return role, nil
}

func set[D, S any](dst *D, src S, proc func(S) (D, error)) error {
	val, err := proc(src)
	if err != nil {
		return err
	}

	*dst = val
	return nil
}
