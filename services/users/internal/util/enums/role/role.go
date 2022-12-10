package role

import "strings"

type Role string

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

func FromString(str string) (Role, bool) {
	r, ok := roleMap[strings.ToLower(str)]
	return r, ok
}
