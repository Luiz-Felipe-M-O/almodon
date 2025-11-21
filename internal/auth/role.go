package auth

import 	"errors"


// Role represents a user role in the system.
type Role uint8

const (
	// Unlogged represents a user that is not logged in.
	Unlogged Role = iota

	valid_start // exclusive, for validity checks

	// Chief represents a department head/chief user.
	Chief

	// Promoted represents a promoted administrative technician user.
	Promoted

	// Admin represents a standard administrative technician user.
	Admin

	// User represents a standard user.
	User

	valid_end // exclusive, for validity checks

	invalid // sentinel value
)

// Hierarchy defines a partial ordering over the Role type.
//
// A call to an implemeter of [Hierarchy] h can be read as such: if
// h(x, y), then the permissions of x are inherited by y.
type Hierarchy func(Role, Role) bool

// IsValid returns whether the role refers to an actual role, that
// is, not unlogged an on the defined range.
func (l Role) IsValid() bool {
	return valid_start < l && l < valid_end
}

// IsValidOrUnlogged returns whether the role refers to an actual
// role or unlogged, equivalent to:
//
//	l == Unlogged || l.IsValid()
func (l Role) IsValidOrUnlogged() bool {
	return l == Unlogged || valid_start < l && l < valid_end
}

// String returns the string representation of the Role.
func (l Role) String() string {
	return roleStrings[l]
}

func (l Role) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}

var ErrBadJSONString = errors.New("cannot unmarshal non-string into Go value of type auth.Role")

func (l *Role) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(str) < 2 || str[0] != '"' || str[len(str)-1] != '"' {
		return ErrBadJSONString
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
// If [DefaultHierarchy](r0, r1) evaluates to true, then the
// permissions of r0 are inherited by r1.
func DefaultHierarchy(r0, r1 Role) bool {
	if !r0.IsValidOrUnlogged() || !r1.IsValidOrUnlogged() {
		return false
	}

	if r0 == Unlogged {
		return true
	}

	if r1 == Unlogged {
		return false
	}

	return r0 >= r1
}

// FromString returns the Role corresponding to the given string. If
// the string does not correspond to any Role, the second return
// value is false.
func FromString(string string) (Role, bool) {
	role, in := stringRoles[string]
	if !in {
		return Unlogged, false
	}

	return role, true
}

var roleStrings = map[Role]string{
	Chief:    "chief",
	Promoted: "promoted-admin",
	Admin:    "admin",
	User:     "user",

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
