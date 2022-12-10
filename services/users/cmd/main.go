/*
The user service handles authentication & authorization of users in Clubrizer.

Other services and the frontend can interact with the user service via it's REST API.
*/
package main

import (
	"github.com/clubrizer/services/users/internal/authenticator/clubrizer"
	"github.com/clubrizer/services/users/internal/authenticator/google"
	"github.com/clubrizer/services/users/internal/storage"

	//"github.com/clubrizer/services/users/internal/storage"
	"github.com/clubrizer/services/users/internal/tokener"
	"github.com/clubrizer/services/users/internal/util/appconfig"
	"github.com/clubrizer/services/users/pkg/http"
	"github.com/clubrizer/services/users/pkg/http/handlers"
)

func main() {
	config := appconfig.Load()

	userStore := storage.NewUserEditor(config.Postgres)

	gAuth := google.NewAuthenticator(config.Auth)
	auth := clubrizer.NewAuthenticator(config.Init, userStore)
	tokenGenerator := tokener.NewGenerator(config.Auth.Jwt)

	loginHandler := handlers.NewLoginHandler(config.Auth.Jwt, gAuth, auth, tokenGenerator)
	tokenHandler := handlers.NewTokenHandler(config.Auth.Jwt, auth, tokenGenerator)

	r := http.NewRouter(config.Server, gAuth, loginHandler, tokenHandler)
	r.Listen()
}
