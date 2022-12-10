package handlers

import (
	"github.com/clubrizer/services/users/internal/authenticator/clubrizer"
	"github.com/clubrizer/services/users/internal/util/appconfig"
	"net/http"
	"time"
)

func addTokensToResponse(tokener tokener, jwtConfig appconfig.Jwt, user *clubrizer.User, w http.ResponseWriter) bool {
	accessToken, err := tokener.GenerateAccessToken(user)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return false
	}
	refreshToken, refreshTokenExpiresAt, err := tokener.GenerateRefreshToken(user)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return false
	}

	w.Header().Set(jwtConfig.AccessToken.HeaderName, accessToken)
	http.SetCookie(w, getRefreshTokenCookie(refreshToken, refreshTokenExpiresAt, jwtConfig))

	return true
}

func getRefreshTokenCookie(token string, expiresAt time.Time, jwtConfig appconfig.Jwt) *http.Cookie {
	sameSiteMode := http.SameSiteDefaultMode
	sameSiteConfig := jwtConfig.RefreshToken.Cookie.SameSite
	if sameSiteConfig == "none" {
		sameSiteMode = http.SameSiteNoneMode
	} else if sameSiteConfig == "lax" {
		sameSiteMode = http.SameSiteLaxMode
	} else if sameSiteConfig == "strict" {
		sameSiteMode = http.SameSiteStrictMode
	}

	return &http.Cookie{
		Name:     jwtConfig.RefreshToken.Cookie.Name,
		Value:    token,
		Expires:  expiresAt,
		HttpOnly: jwtConfig.RefreshToken.Cookie.HttpOnly,
		Secure:   jwtConfig.RefreshToken.Cookie.Secure,
		SameSite: sameSiteMode,
	}
}
