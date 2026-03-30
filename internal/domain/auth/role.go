package auth

import (
	"errors"

	"github.com/alan-b-lima/almodon/pkg/rbac"
)

// Role represents a user role in the system.
type Role uint8

const (
	// Unlogged represents a user that is not logged in.
	Unlogged Role = iota

	Maintainer // Maintainer represents a system maintainer.

	Chief    // Chief represents a department head/chief user.
	Promoted // Promoted represents a promoted administrative technician user.
	Admin    // Admin represents a standard administrative technician user.
	User     // User represents a standard user.

	invalid // sentinel value
)

func Allow(classes ...Role) rbac.Permission[Role] {
	return rbac.Allow(DefaultHierarchy, classes...)
}

// IsValid returns whether the role refers to an actual role, that is, not
// unlogged an on the defined range.
func (l Role) IsValid() bool {
	return l != Unlogged && l < invalid
}

// String returns the string representation of the Role.
func (l Role) String() string {
	return roleStrings[l]
}

func (l Role) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}

var errBadJSONString = errors.New("cannot unmarshal non-string into Go value of type auth.Role")

func (l *Role) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(str) < 2 || str[0] != '"' || str[len(str)-1] != '"' {
		return errBadJSONString
	}

	role, ok := FromString(str[1 : len(str)-1])
	if !ok {
		*l = invalid
		return nil
	}

	*l = role
	return nil
}

// DefaultHierarchy defines a partial ordering in the Role type.
//
// If [DefaultHierarchy](x, y) evaluates to true, then the permissions of x are
// inherited by y.
//
// Considere [DefaultHierarchy](x, y) iff x < y, then the hierarchy is defined
// as follows:
//
//   - [Unlogged] < [User] < [Admin] < [Promoted] < [Chief] < [Maintainer]
func DefaultHierarchy(x, y Role) bool {
	if x == Unlogged {
		return true
	}

	if y == Unlogged {
		return false
	}

	if !x.IsValid() || !y.IsValid() {
		return false
	}

	return x >= y
}

// FromString returns the Role corresponding to the given string. If the string
// does not correspond to any Role, it returns false.
func FromString(string string) (Role, bool) {
	role, in := stringRoles[string]
	if !in {
		return Unlogged, false
	}

	return role, true
}

var roleStrings = map[Role]string{
	Maintainer: "maintainer",
	Chief:      "chief",
	Promoted:   "promoted-admin",
	Admin:      "admin",
	User:       "user",

	Unlogged: "unlogged",
}

var stringRoles = mirror(roleStrings)

func mirror[K, V comparable](m map[K]V) map[V]K {
	nm := make(map[V]K, len(m))
	for k, v := range m {
		nm[v] = k
	}
	return nm
}
