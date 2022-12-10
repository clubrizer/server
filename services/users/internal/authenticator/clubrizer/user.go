package clubrizer

import (
	"github.com/clubrizer/services/users/internal/util/enums/approvalstate"
	"github.com/clubrizer/services/users/internal/util/enums/role"
)

// A User represents a Clubrizer user.
type User struct {
	Id         int64
	IsAdmin    bool
	TeamClaims []TeamClaim
}

// A TeamClaim defines which role a user has within a certain team.
type TeamClaim struct {
	Name          string
	ApprovalState approvalstate.ApprovalState
	Role          role.Role
}
