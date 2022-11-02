package config

import (
	"flag"
	"github.com/clubrizer/server/pkg/log"
	"github.com/spf13/viper"
	"os"
)

type Server struct {
	Port string
}

type Postgres struct {
	Url      string
	User     string
	Password string
}

type Messaging struct {
	ProjectId       string
	CredentialsFile string
}

type SimpleAppConfig struct {
	Server    Server
	Postgres  Postgres
	Messaging Messaging
}

// viper also supports watching and re-reading appconfig files: https://github.com/spf13/viper#watching-and-re-reading-config-files

func Load(out interface{}) {
	configFile := flag.String("config", "config.yaml", "relative path to the configuration file (must be in yaml format)")
	flag.Parse()

	viper.SetConfigFile(*configFile)
	viper.AutomaticEnv()

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
