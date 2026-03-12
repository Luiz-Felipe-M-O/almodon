// Copyright (C) 2025 Alan Barbosa Lima.
//
// Almodon is licensed under the GNU General Public License
// version 3. You should have received a copy of the
// license, located in LICENSE, at the root of the source
// tree. If not, see <https://www.gnu.org/licenses/>.

// Package auth implements a RBAC (Role Based Access Control)
// low-level framework.
package auth

import (
	"encoding/json"
	"fmt"
)

// Hierarchy is a partial ordering over the set of roles.
//
// Given a [Hierarchy] h, over the set L, x, y in L, h(x, y) is true if, and
// only if the permissions of x are inherited by y.
//
// A partial ordering is defined as a relation h over a set L, x, y, z in L,
// that fulfills:
//   - Reflexivity: h(x, x) is true;
//   - Antisymmetry: if h(x, y) is true, then h(y, x) is false, if x != y; and
//   - Transitivity: if h(x, y) and h(y, z) are true, then h(x, z) is true.
type Hierarchy[R any] func(R, R) bool

// Permission represents an authorization requirement.
type Permission[R any] struct {
	classes   []R
	hierarchy Hierarchy[R]
}

// Allow creates a Permission that authorizes any of the given roles according
// to the provided hierarchy.
func Allow[R any](hierarchy Hierarchy[R], classes ...R) Permission[R] {
	var heads []R

Outer:
	for _, class := range classes {
		for i := 0; i < len(heads); i++ {
			switch head := heads[i]; {
			case hierarchy(class, head):
				continue Outer

			case hierarchy(head, class):
				if i == len(heads)-1 {
					heads = heads[:i]
				} else {
					heads[i] = heads[len(heads)-1]
					heads = heads[:i]
					i--
				}
			}
		}

		heads = append(heads, class)
	}

	return Permission[R]{
		classes:   heads,
		hierarchy: hierarchy,
	}
}

// Authorize returns whether the given role is authorized by the Permission.
func (auth *Permission[R]) Authorize(role R) bool {
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
