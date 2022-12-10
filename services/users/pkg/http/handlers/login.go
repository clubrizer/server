// Package handlers contains HTTP handlers that are executed on incoming requests.
package handlers

import (
	"context"
	"github.com/clubrizer/services/users/internal/authenticator/clubrizer"
	"github.com/clubrizer/services/users/internal/authenticator/google"
	"github.com/clubrizer/services/users/internal/util/appconfig"
	"net/http"
)

type googleAuthenticator interface {
	GetUserFromContext(ctx context.Context) (*google.User, bool)
}

type authenticator interface {
	Authenticate(user *google.User) (*clubrizer.User, error)
}

// LoginHandler contains HTTP handlers for authenticating users.
type LoginHandler struct {
	jwtConfig appconfig.Jwt
	gAuth     googleAuthenticator
	auth      authenticator
	tokener   tokener
}

// NewLoginHandler creates a new [LoginHandler]
func NewLoginHandler(jwtConfig appconfig.Jwt, gAuth googleAuthenticator, auth authenticator, tokener tokener) *LoginHandler {
	return &LoginHandler{
		jwtConfig: jwtConfig,
		gAuth:     gAuth,
		auth:      auth,
		tokener:   tokener,
	}
}

// Authenticate authenticates/logs a user in and sets and access & refresh token on the response.
func (h LoginHandler) Authenticate() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		googleUser, ok := h.gAuth.GetUserFromContext(r.Context())
		if !ok {
			http.Error(w, "Failed to get Google user claim", http.StatusInternalServerError)
			return
		}

		user, err := h.auth.Authenticate(googleUser)
		if err != nil {
			http.Error(w, "Failed to login or register", http.StatusInternalServerError)
			return
		}

		if addTokensToResponse(h.tokener, h.jwtConfig, user, w) {
			w.WriteHeader(http.StatusOK)
		}
	}
}
