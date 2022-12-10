// Package role contains an enum representing the role of a user inside a team.
package role

import "strings"

// Role is an enum representing the role of a user inside a team.
type Role string

// Member & Admin are the roles a user can have within a team.
const (
	Member Role = "member"
	Admin       = "admin"
)

var (
	roleMap = map[string]Role{
		"member": Member,
		"admin":  Admin,
	}
)

// FromString parses a string to an [Role] enum
func FromString(str string) (Role, bool) {
	r, ok := roleMap[strings.ToLower(str)]
	return r, ok
}
