// Package http provides an HTTP API to interact with.
package http

import (
	"context"
	"fmt"
	"github.com/clubrizer/server/pkg/config"
	"github.com/clubrizer/server/pkg/log"
	"github.com/clubrizer/services/users/pkg/http/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

type loginHandler interface {
	Authenticate() func(w http.ResponseWriter, r *http.Request)
}

type tokenHandler interface {
	ValidateAccessToken() func(w http.ResponseWriter, r *http.Request)
	RefreshTokens() func(w http.ResponseWriter, r *http.Request)
}

type googleAuthenticator interface {
	AddUserToContext(ctx context.Context, idToken string) (context.Context, error)
}

// A Router is responsible for handling HTTP requests.
type Router struct {
	serverConfig        config.Server
	googleAuthenticator googleAuthenticator
	loginHandler        loginHandler
	tokenHandler        tokenHandler
}

// NewRouter generates a new [Router].
func NewRouter(
	serverConfig config.Server,
	googleAuthenticator googleAuthenticator,
	loginHandler loginHandler,
	tokenHandler tokenHandler,
) *Router {
	return &Router{
		serverConfig, googleAuthenticator,
		loginHandler, tokenHandler,
	}
}

// Listen listens for HTTP requests on the port provided in the server config.
func (cr Router) Listen() {
	address := fmt.Sprintf(":%s", cr.serverConfig.Port)
	log.Info("Listening for HTTP requests on %s", address)

	r := chi.NewRouter()
	configureRouter(r, cr.serverConfig.Cors)

	r.
		With(middleware.GoogleAuthenticator(cr.googleAuthenticator)).
		Post("/login", cr.loginHandler.Authenticate())

	r.Route("/token", func(r chi.Router) {
		r.Post("/validate", cr.tokenHandler.ValidateAccessToken())
		r.Post("/refresh", cr.tokenHandler.RefreshTokens())
	})

	log.Fatal(
		http.ListenAndServe(address, r),
		"Failed to listen for HTTP requests on %s.", address,
	)
}

func configureRouter(router *chi.Mux, corsConfig config.Cors) {
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   corsConfig.AllowedOrigins,
		AllowedMethods:   corsConfig.AllowedMethods,
		AllowedHeaders:   corsConfig.AllowedHeaders,
		ExposedHeaders:   corsConfig.ExposedHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
		MaxAge:           corsConfig.MaxAge,
	}))
	router.Use(chiMiddleware.Logger) // consider using the Clubrizer log package (logrus)
}
