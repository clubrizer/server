package google_test

import (
	"context"
	"errors"
	"github.com/clubrizer/services/users/internal/authenticator/google"
	mocks "github.com/clubrizer/services/users/internal/mocks/authenticator/google"
	"github.com/clubrizer/services/users/internal/util/appconfig"
	"google.golang.org/api/idtoken"
	"reflect"
	"testing"
)

func TestAuthenticator_AddAndGetUser(t *testing.T) {
	ctx := context.Background()
	idToken := "google id token"
	user := &google.User{
		Issuer:     "google",
		ID:         "123",
		GivenName:  "John",
		FamilyName: "Doe",
		Email:      "hi@john.doe",
	}
	payload := &idtoken.Payload{
		Issuer:  user.Issuer,
		Subject: user.ID,
		Claims: map[string]interface{}{
			"given_name":  user.GivenName,
			"family_name": user.FamilyName,
			"email":       user.Email,
			"picture":     user.Picture,
		},
	}

	type prepareFields struct {
		validator func(v *mocks.IdTokenValidator)
	}
	type args struct {
		idToken string
		ctx     context.Context
	}
	tests := []struct {
		name          string
		prepareFields prepareFields
		args          args
		wantSetErr    bool
		wantGetOk     bool
		wantUser      *google.User
	}{
		{
			name: "ok",
			prepareFields: prepareFields{
				validator: func(v *mocks.IdTokenValidator) {
					v.EXPECT().Validate(ctx, idToken, "").Once().Return(payload, nil)
				},
			},
			args: args{
				idToken: idToken,
				ctx:     context.Background(),
			},
			wantSetErr: false,
			wantGetOk:  true,
			wantUser:   user,
		},
		{
			name: "token validation failed",
			prepareFields: prepareFields{
				validator: func(v *mocks.IdTokenValidator) {
					v.EXPECT().Validate(ctx, idToken, "").Once().Return(nil, errors.New("err"))
				},
			},
			args: args{
				idToken: idToken,
				ctx:     context.Background(),
			},
			wantSetErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := mocks.NewIdTokenValidator(t)
			if tt.prepareFields.validator != nil {
				tt.prepareFields.validator(v)
			}

			a := google.NewTestAuthenticator(appconfig.Auth{}, v)

			ctx, err := a.AddUserToContext(tt.args.ctx, tt.args.idToken)
			if (err != nil) != tt.wantSetErr {
				t.Errorf("AddUserToContext() error = %v, wantErr %v", err, tt.wantSetErr)
				return
			}
			if err != nil {
				return
			}
			gotUser, ok := a.GetUserFromContext(ctx)
			if ok != tt.wantGetOk {
				t.Errorf("GetUserFromContext() ok = %v, wantOk %v", err, tt.wantGetOk)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("AddUserToContext() & GetUserFromContext() got = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}
