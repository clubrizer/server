package appconfig

import "github.com/clubrizer/server/pkg/config"

type AppConfig struct {
	Server    config.Server
	Postgres  config.Postgres
	Messaging config.Messaging
	Auth      Auth
	Init      Init
}

func Load() *AppConfig {
	var appConfig AppConfig
	config.Load(&appConfig)
	return &appConfig
}
