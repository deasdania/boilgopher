package main

import (
	"boilgopher/storage/postgres"
	"log"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	version = "devel"
	config  *viper.Viper
)

const service = "boilgopher"

func main() {
	config = viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}
	newLogger := logrus.New()
	logger := newLogger.WithFields(logrus.Fields{
		"service": service,
		"version": version,
	})
	logger.Println("starting service...")

	// use it when we need
	_, err := postgres.New(config)
	if err != nil {
		log.Fatalf("failed creating postgres storage: %v", err)
	}
}
