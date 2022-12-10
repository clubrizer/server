package google

import (
	"github.com/clubrizer/services/users/internal/util/appconfig"
)

func NewTestAuthenticator(config appconfig.Auth, v idTokenValidator) *Authenticator {
	return &Authenticator{config, v}
}
