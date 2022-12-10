// Package middleware contains HTTP middleware.
package middleware

import (
	"context"
	"net/http"
)

type googleAuthenticator interface {
	AddUserToContext(ctx context.Context, idToken string) (context.Context, error)
}

// GoogleAuthenticator is responsible for validating Google access tokens. If the access token is valid, the Google
// user is added to the request context.
func GoogleAuthenticator(gAuth googleAuthenticator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, err := gAuth.AddUserToContext(r.Context(), r.Header.Get("Authorization"))
			if err != nil {
				http.Error(w, "Invalid ID Token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
