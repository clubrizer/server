package clubrizer_test

import (
	"errors"
	"github.com/clubrizer/server/pkg/storageutils"
	"github.com/clubrizer/services/users/internal/authenticator/clubrizer"
	"github.com/clubrizer/services/users/internal/authenticator/google"
	mocks "github.com/clubrizer/services/users/internal/mocks/authenticator/clubrizer"
	"github.com/clubrizer/services/users/internal/storage"
	"github.com/clubrizer/services/users/internal/util/appconfig"
	"github.com/clubrizer/services/users/internal/util/enums/approvalstate"
	"github.com/clubrizer/services/users/internal/util/enums/role"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestAuthenticator_Authenticate(t *testing.T) {
	googleUser := google.User{
		Issuer:     "google",
		Id:         "123",
		GivenName:  "John",
		FamilyName: "Doe",
		Email:      "john@doe.com",
	}

	type args struct {
		googleUser *google.User
	}
	tests := []struct {
		name       string
		prepare    func(r *mocks.UserRepository)
		initConfig appconfig.Init
		args       args
		want       *clubrizer.User
		wantErr    bool
	}{
		{
			name: "user exists",
			prepare: func(r *mocks.UserRepository) {
				r.EXPECT().GetFromExternalId(googleUser.Issuer, googleUser.Id).Times(1).Return(&storage.User{Id: 1}, nil)
			},
			initConfig: appconfig.Init{},
			args:       args{googleUser: &googleUser},
			want:       &clubrizer.User{Id: 1},
			wantErr:    false,
		},
		{
			name: "create user",
			prepare: func(r *mocks.UserRepository) {
				r.EXPECT().GetFromExternalId(googleUser.Issuer, googleUser.Id).Times(1).Return(nil, storageutils.ErrNotFound)
				r.EXPECT().Create(&googleUser, false).Times(1).Return(&storage.User{
					Id:         1,
					GivenName:  googleUser.GivenName,
					FamilyName: googleUser.FamilyName,
					Email:      googleUser.Email,
					IsAdmin:    false,
				}, nil)
			},
			initConfig: appconfig.Init{},
			args:       args{googleUser: &googleUser},
			want: &clubrizer.User{
				Id:      1,
				IsAdmin: false,
			},
			wantErr: false,
		},
		{
			name: "create admin",
			prepare: func(r *mocks.UserRepository) {
				r.EXPECT().GetFromExternalId(googleUser.Issuer, googleUser.Id).Times(1).Return(nil, storageutils.ErrNotFound)
				r.EXPECT().Create(&googleUser, true).Times(1).Return(&storage.User{
					Id:         1,
					GivenName:  googleUser.GivenName,
					FamilyName: googleUser.FamilyName,
					Email:      googleUser.Email,
					IsAdmin:    true,
				}, nil)
			},
			initConfig: appconfig.Init{AdminEmail: googleUser.Email},
			args:       args{googleUser: &googleUser},
			want: &clubrizer.User{
				Id:      1,
				IsAdmin: true,
			},
			wantErr: false,
		},
		{
			name: "get user fail",
			prepare: func(r *mocks.UserRepository) {
				r.EXPECT().GetFromExternalId(googleUser.Issuer, googleUser.Id).Times(1).Return(nil, errors.New("something failed :("))
			},
			initConfig: appconfig.Init{},
			args:       args{googleUser: &googleUser},
			want:       nil,
			wantErr:    true,
		},
		{
			name: "create user fail",
			prepare: func(r *mocks.UserRepository) {
				r.EXPECT().GetFromExternalId(googleUser.Issuer, googleUser.Id).Times(1).Return(nil, storageutils.ErrNotFound)
				r.EXPECT().Create(&googleUser, false).Times(1).Return(nil, errors.New("something failed :("))
			},
			initConfig: appconfig.Init{},
			args:       args{googleUser: &googleUser},
			want:       nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mocks.NewUserRepository(t)
			if tt.prepare != nil {
				tt.prepare(r)
			}
			a := clubrizer.NewAuthenticator(tt.initConfig, r)
			got, err := a.Authenticate(tt.args.googleUser)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("Authenticate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthenticator_Get(t *testing.T) {
	type args struct {
		userId int64
	}
	tests := []struct {
		name       string
		prepare    func(r *mocks.UserRepository)
		initConfig appconfig.Init
		args       args
		want       *clubrizer.User
		wantErr    bool
	}{
		{
			name: "existing user",
			prepare: func(r *mocks.UserRepository) {
				r.EXPECT().GetFromId(int64(1)).Times(1).Return(&storage.User{Id: 1}, nil)
			},
			initConfig: appconfig.Init{},
			args:       args{userId: 1},
			want:       &clubrizer.User{Id: 1},
			wantErr:    false,
		},
		{
			name: "not existing user",
			prepare: func(r *mocks.UserRepository) {
				r.EXPECT().GetFromId(int64(1)).Times(1).Return(nil, errors.New("some failure"))
			},
			initConfig: appconfig.Init{},
			args:       args{userId: 1},
			want:       nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mocks.NewUserRepository(t)
			if tt.prepare != nil {
				tt.prepare(r)
			}
			a := clubrizer.NewAuthenticator(tt.initConfig, r)
			got, err := a.Get(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mapToUser(t *testing.T) {
	wantedUser := &storage.User{
		Id:         1,
		GivenName:  "John",
		FamilyName: "Doe",
		Email:      "john@doe.com",
		Picture:    "https://super-fancy-pic.com",
		IsAdmin:    false,
		TeamClaims: []storage.TeamClaim{{
			Team: storage.Team{
				Name:        "test",
				DisplayName: "test team",
			},
			ApprovalState: approvalstate.Pending,
			Role:          role.Member,
		},
		},
	}

	type args struct {
		user *storage.User
	}
	tests := []struct {
		name string
		args args
		want *clubrizer.User
	}{
		{name: "map user",
			args: args{wantedUser},
			want: &clubrizer.User{
				Id:      wantedUser.Id,
				IsAdmin: wantedUser.IsAdmin,
				TeamClaims: []clubrizer.TeamClaim{
					{
						Name:          wantedUser.TeamClaims[0].Team.Name,
						ApprovalState: wantedUser.TeamClaims[0].ApprovalState,
						Role:          wantedUser.TeamClaims[0].Role,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clubrizer.MapToUser(tt.args.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapToUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
