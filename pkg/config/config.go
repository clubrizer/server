// Package config loads the config for services from config files. Further, it provides some shared config structs that
// are the same in all services.
package config

import (
	"flag"
	"github.com/clubrizer/server/pkg/log"
	"github.com/spf13/viper"
	"os"
)

// Cors represents the CORS config block.
type Cors struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// Server represents the Server config block.
type Server struct {
	Port string
	Cors Cors
}

// Postgres represents the Postgres config block.
type Postgres struct {
	URL      string
	User     string
	Password string
}

// Messaging represents the Messaging config block.
type Messaging struct {
	ProjectID       string
	CredentialsFile string
}

// SimpleAppConfig represents the most basic app config for a service.
type SimpleAppConfig struct {
	Server    Server
	Postgres  Postgres
	Messaging Messaging
}

// Load loads the app config into the given struct.
func Load(out interface{}) {
	configFile := flag.String("config", "config.yaml", "relative path to the configuration file (must be in yaml format)")
	flag.Parse()

	viper.SetConfigFile(*configFile)
	viper.AutomaticEnv()

	// viper also supports watching and re-reading appconfig files: https://github.com/spf13/viper#watching-and-re-reading-config-files
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err, "Failed to read application config from ./%s", *configFile)
	}

	// Unfortunately, this is necessary because of https://github.com/spf13/viper/issues/761
	for _, k := range viper.AllKeys() {
		v := viper.GetString(k)
		viper.Set(k, os.ExpandEnv(v))
	}

	err = viper.Unmarshal(out)
	if err != nil {
		log.Fatal(err, "Failed to unmarshal application config from ./%s", *configFile)
	}

	log.Info("Loaded the application config from ./%s", *configFile)
}
