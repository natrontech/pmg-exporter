package auth

import "fmt"

type UserRole int

const (
	RoleAnonymous UserRole = 0
	RoleUser      UserRole = 1
)

func NewRole(role string) (UserRole, error) {
	switch role {
	case "ANONYMOUS":
		return RoleAnonymous, nil
	case "USER":
		return RoleUser, nil
	}

	return 0, fmt.Errorf("invalid role specified: %s", role)
}

func (r *UserRole) IsAuthorized(minimumRole UserRole) bool {
	return *r >= minimumRole
}
