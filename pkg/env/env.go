package env

import (
	"flag"
	"fmt"
	"github.com/clubrizer/log"
	"github.com/joho/godotenv"
	"os"
)

func init() {
	devModeFlag := flag.Bool("dev", false, "set for development mode")
	flag.Parse()
	if *devModeFlag {
		loadDotEnv()
	}
}

func Port() string {
	return os.Getenv("PORT")
}

func loadDotEnv() {
	log.Info("Loading .env file")
	workingDir, err := os.Getwd()
	if err != nil {
		log.Error(err, "Error loading .env file")
	}
	err = godotenv.Load(fmt.Sprintf("%s/.env", workingDir))
	if err != nil {
		log.Error(err, "Error loading .env file")
	}
	log.Debug(".env file loaded")
}
