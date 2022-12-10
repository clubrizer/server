// Package appconfig is responsible for handling the apps configuration.
package appconfig

import "github.com/clubrizer/server/pkg/config"

// AppConfig contains the whole configuration for this app/service.
type AppConfig struct {
	Server    config.Server
	Postgres  config.Postgres
	Messaging config.Messaging
	Auth      Auth
	Init      Init
}

// Load loads the config file into a struct and returns it.
func Load() *AppConfig {
	var appConfig AppConfig
	config.Load(&appConfig)
	return &appConfig
}
