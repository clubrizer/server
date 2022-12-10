// Package approvalstate contains an enum representing the approval state of a user inside a team.
package approvalstate

import "strings"

// ApprovalState is an enum representing the approval state of a user inside a team.
type ApprovalState string

// Pending, Approved & Declined are the approval states a user can have within a team.
const (
	Pending  ApprovalState = "pending"
	Approved               = "approved"
	Declined               = "declined"
)

var (
	approvalStateMap = map[string]ApprovalState{
		"pending":  Pending,
		"approved": Approved,
		"declined": Declined,
	}
)

// FromString parses a string to an [ApprovalState] enum
func FromString(str string) (ApprovalState, bool) {
	s, ok := approvalStateMap[strings.ToLower(str)]
	return s, ok
}
