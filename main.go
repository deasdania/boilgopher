package main

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

var (
	version = "devel"
	config  *viper.Viper
)

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
}
