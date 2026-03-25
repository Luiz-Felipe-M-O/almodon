package auth

import "github.com/alan-b-lima/almodon/pkg/uuid"

// Actor is a system entity that can access services through it's
// role. It is related the user entity, as a Actor is a shrinked
// version of an user.
type Actor struct {
	User uuid.UUID
	Role Role
}

// NewLogged creates a new actor. This function does not check
// whether the fact is real.
func NewLogged(user uuid.UUID, role Role) Actor {
	return Actor{
		User: user,
		Role: role,
	}
}

// NewUnlogged creates a new unlogged actor. It's also equivalent to
// the zero value of [Actor].
func NewUnlogged() Actor {
	return Actor{Role: Unlogged}
}
