package approvalstate_test

import (
	"github.com/clubrizer/services/users/internal/util/enums/approvalstate"
	"testing"
)

func TestFromString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name   string
		args   args
		want   approvalstate.ApprovalState
		wantOk bool
	}{
		{
			name:   "pending",
			args:   args{"pending"},
			want:   approvalstate.Pending,
			wantOk: true,
		},
		{
			name:   "approved",
			args:   args{"approved"},
			want:   approvalstate.Approved,
			wantOk: true,
		},
		{
			name:   "declined",
			args:   args{"declined"},
			want:   approvalstate.Declined,
			wantOk: true,
		},
		{
			name:   "invalid",
			args:   args{"invalid"},
			want:   "",
			wantOk: false,
		},
		{
			name:   "empty",
			args:   args{""},
			want:   "",
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := approvalstate.FromString(tt.args.str)
			if got != tt.want {
				t.Errorf("FromString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantOk {
				t.Errorf("FromString() got1 = %v, want %v", got1, tt.wantOk)
			}
		})
	}
}
