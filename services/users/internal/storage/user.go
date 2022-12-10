package storage

import (
	"github.com/clubrizer/services/users/internal/util/enums/approvalstate"
	"github.com/clubrizer/services/users/internal/util/enums/role"
)

// A User represents a user with team claims as it is in the database.
type User struct {
	Id         int64
	GivenName  string
	FamilyName string
	Email      string
	Picture    string
	IsAdmin    bool
	TeamClaims []TeamClaim
}

// A TeamClaim defines the role of a user inside a team.
type TeamClaim struct {
	Team          Team
	ApprovalState approvalstate.ApprovalState
	Role          role.Role
}

// A Team represents one club/team that a user can join.
type Team struct {
	Name        string
	DisplayName string
	Description string
	Logo        string
}
