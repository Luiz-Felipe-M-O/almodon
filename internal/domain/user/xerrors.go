package user

import "github.com/alan-b-lima/almodon/pkg/errors"

var (
	ErrUserCreate = errors.Imp(errors.InvalidInput, "user-create").Message("could not create user")
	ErrUserUpdate = errors.Imp(errors.InvalidInput, "user-update").Message("could not update user")
	ErrUserNotFound = errors.New(errors.NotFound, "user-not-found", "user not found",nil,nil)

	ErrNameEmpty = errors.New(errors.InvalidInput, "name-empty", "name must not be empty", nil, nil)

	ErrSiapeTaken=errors.New(errors.Conflict, "siape-in-use", "siape already taken",nil,nil)

	ErrEmailInvalid = errors.New(errors.InvalidInput, "email-invalid", "email invalid", nil, map[string]any{"pattern": reEmail})

	ErrPasswordTooShort              = errors.New(errors.InvalidInput, "password-too-short", "password too short", nil, map[string]any{"max": PasswordMinLen})
	ErrPasswordTooLong               = errors.New(errors.InvalidInput, "password-too-long", "password too long", nil, map[string]any{"max": PasswordMaxLen})
	ErrPasswordLeadOrTrailWhitespace = errors.New(errors.InvalidInput, "password-edge-whitespace", "password must not start or end in whitespace", nil, nil)
	ErrPasswordIllegalCharacters     = errors.New(errors.InvalidInput, "password-illegal-chars", "password must not include unprintable or illegal UTF-8 characters", nil, nil)
	ErrPasswordFailedToHash          = errors.Imp(errors.Internal, "password-hash-failed").Message("unexpected error ocurred while hashing password")
	ErrPasswordIncorrect             = errors.New(errors.Internal, "password-incorrect", "incorrect password", nil, nil)

	ErrRoleInvalid = errors.Imp(errors.InvalidInput, "role-invalid").Metadata(map[string]any{"accept_roles": acceptRoles}).Format("role must be one of %v").Make(acceptRoles)
	ErrNotEnoughChiefs=errors.New(errors.Conflict, "not-enough-chiefs", "there must be at least one chief at any moment",nil,nil)
)
