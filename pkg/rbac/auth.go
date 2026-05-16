// Copyright (C) 2025 Alan Barbosa Lima.
//
// Almodon is licensed under the GNU General Public License
// version 3. You should have received a copy of the
// license, located in LICENSE, at the root of the source
// tree. If not, see <https://www.gnu.org/licenses/>.

// Package rbac implements a RBAC (Role Based Access Control)
// low-level framework.
package rbac

import (
	"encoding/json"
	"fmt"
)

// Hierarchy is a partial ordering over the set of roles.
//
// Given a [Hierarchy] h, over the set L, x, y in L, h(x, y) is true if, and
// only if the permissions of x are inherited by y.
//
// A partial ordering is defined as a relation H over a set L, x, y, z in L,
// that fulfills:
//   - Reflexivity: h(x, x) is true;
//   - Antisymmetry: if h(x, y) and h(y, x) are true, then x == y; and
//   - Transitivity: if h(x, y) and h(y, z) are true, then h(x, z) is true.
type Hierarchy[R any] func(R, R) bool

// Permission represents an authorization requirement.
type Permission[R any] struct {
	classes   []R
	hierarchy Hierarchy[R]
}

// Allow creates a Permission that authorizes any of the given roles according
// to the provided hierarchy.
func Allow[R any](hierarchy Hierarchy[R], roles ...R) Permission[R] {
	leaves := roles[:0]

Outer:
	for _, role := range roles {
		for i := 0; i < len(leaves); i++ {
			leaf := leaves[i]

			if hierarchy(role, leaf) {
				continue Outer
			}

			if hierarchy(leaf, role) {
				leaves[i] = leaves[len(leaves)-1]
				leaves = leaves[:i]
				i--
			}
		}

		leaves = append(leaves, role)
	}

	return Permission[R]{
		classes:   leaves,
		hierarchy: hierarchy,
	}
}

// Allows reports whether the given role is authorized by the permission.
func (auth Permission[R]) Allows(role R) bool {
	for _, class := range auth.classes {
		if auth.hierarchy(class, role) {
			return true
		}
	}

	return false
}

// MarshalJSON implements the [json.Marshaler] interface.
func (auth Permission[R]) MarshalJSON() ([]byte, error) {
	return json.Marshal(auth.classes)
}

// String returns the string representation of the Permission.
func (auth Permission[R]) String() string {
	return fmt.Sprint(auth.classes)
}
