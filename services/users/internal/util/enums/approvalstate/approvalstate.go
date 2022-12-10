package approvalstate

import "strings"

type ApprovalState string

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

func FromString(str string) (ApprovalState, bool) {
	s, ok := approvalStateMap[strings.ToLower(str)]
	return s, ok
}
