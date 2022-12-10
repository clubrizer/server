package role_test

import (
	"github.com/clubrizer/services/users/internal/util/enums/role"
	"testing"
)

func TestFromString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name   string
		args   args
		want   role.Role
		wantOk bool
	}{
		{
			name:   "member",
			args:   args{"member"},
			want:   role.Member,
			wantOk: true,
		},
		{
			name:   "admin",
			args:   args{"admin"},
			want:   role.Admin,
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
			got, got1 := role.FromString(tt.args.str)
			if got != tt.want {
				t.Errorf("FromString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantOk {
				t.Errorf("FromString() got1 = %v, want %v", got1, tt.wantOk)
			}
		})
	}
}
