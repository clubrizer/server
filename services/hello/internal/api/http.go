package api

import (
	"fmt"
	"github.com/clubrizer/server/pkg/env"
	"github.com/clubrizer/server/pkg/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func Serve() {
	address := fmt.Sprintf(":%s", env.Port())
	log.Info("Listening for HTTP requests on %s", address)

	router := chi.NewRouter()
	router.Use(middleware.Logger) // consider using the Clubrizer log package (logrus)

	router.Get("/", sayHello)

	log.Fatal(
		http.ListenAndServe(address, router),
		"Failed to listen for HTTP requests on %s.", address,
	)
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("Hello world"))
}
