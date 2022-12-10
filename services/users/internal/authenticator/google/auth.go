// Package google handles authentication of users that login via Google.
package google

import (
	"context"
	"github.com/clubrizer/server/pkg/log"
	"github.com/clubrizer/services/users/internal/util/appconfig"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

type userKey string

const (
	googleUserKey userKey = "google-user"
)

type idTokenValidator interface {
	Validate(ctx context.Context, idToken string, audience string) (*idtoken.Payload, error)
}

// Authenticator allows users to authenticate against Google.
type Authenticator struct {
	authConfig appconfig.Auth
	validator  idTokenValidator
}

// NewAuthenticator creates a new [Authenticator] with the given parameters.
func NewAuthenticator(authConfig appconfig.Auth) *Authenticator {
	v, err := idtoken.NewValidator(context.Background(), option.WithoutAuthentication())
	if err != nil {
		log.Fatal(err, "Failed to create Google token validator")
	}
	return &Authenticator{authConfig, v}
}

// AddUserToContext gets and validates the Google user that is about to login and adds this user to the given context.
func (a Authenticator) AddUserToContext(ctx context.Context, idToken string) (context.Context, error) {
	tokenPayload, err := a.validator.Validate(context.Background(), idToken, a.authConfig.GoogleClientID)
	if err != nil {
		return nil, err
	}

	user := &User{
		Issuer:     tokenPayload.Issuer,
		ID:         tokenPayload.Subject,
		GivenName:  tokenPayload.Claims["given_name"].(string),
		FamilyName: tokenPayload.Claims["family_name"].(string),
		Email:      tokenPayload.Claims["email"].(string),
		Picture:    tokenPayload.Claims["picture"].(string),
	}

	ctxWithValue := context.WithValue(ctx, googleUserKey, user)

	return ctxWithValue, nil
}

// GetUserFromContext gets the Google user that is set on the given context.
func (a Authenticator) GetUserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(googleUserKey).(*User)
	return user, ok
}
