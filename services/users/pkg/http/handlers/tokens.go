package handlers

import (
	"fmt"
	"github.com/clubrizer/services/users/internal/authenticator/clubrizer"
	"github.com/clubrizer/services/users/internal/util/appconfig"
	"net/http"
	"time"
)

type tokener interface {
	GenerateAccessToken(user *clubrizer.User) (string, error)
	GenerateRefreshToken(user *clubrizer.User) (string, time.Time, error)
	ValidateAccessToken(tokenString string) error
	ValidateRefreshTokenAndGetUserId(tokenString string) (int64, error)
}

type tokenAuthenticator interface {
	Get(userId int64) (*clubrizer.User, error)
}

// TokenHandler contains HTTP handlers for generating and validating JWTs.
type TokenHandler struct {
	jwtConfig appconfig.Jwt
	auth      tokenAuthenticator
	tokener   tokener
}

// NewTokenHandler creates a new [TokenHandler].
func NewTokenHandler(jwtConfig appconfig.Jwt, auth tokenAuthenticator, tokener tokener) *TokenHandler {
	return &TokenHandler{
		jwtConfig: jwtConfig,
		auth:      auth,
		tokener:   tokener,
	}
}

// ValidateAccessToken checks if the access token provided in the request is valid.
func (h TokenHandler) ValidateAccessToken() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get(h.jwtConfig.AccessToken.HeaderName)
		if accessToken == "" {
			http.Error(w, "No JWT access token set", http.StatusUnauthorized)
			return
		}

		err := h.tokener.ValidateAccessToken(accessToken)
		if err != nil {
			http.Error(w, "Invalid JWT access token", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// RefreshTokens refreshes both, the JWT access & refresh token.
func (h TokenHandler) RefreshTokens() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(h.jwtConfig.RefreshToken.Cookie.Name)
		if err != nil {
			http.Error(w, "No JWT refresh token cookie set", http.StatusUnauthorized)
			return
		}

		userId, err := h.tokener.ValidateRefreshTokenAndGetUserId(cookie.Value)
		if err != nil {
			http.Error(w, "Invalid JWT refresh token", http.StatusUnauthorized)
			return
		}
		user, err := h.auth.Get(userId)
		if err != nil {
			http.Error(w, fmt.Sprintf("No user found for the id %d", userId), http.StatusNotFound)
			return
		}

		if addTokensToResponse(h.tokener, h.jwtConfig, user, w) {
			w.WriteHeader(http.StatusOK)
		}
	}
}
