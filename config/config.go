package config

import (
	"github.com/spf13/viper"
	"log"
)

var EnvConfig *config

func InitConfig() {
	EnvConfig = loadEnvVariables("./config")
}

type config struct {
	HousesUrl     string `mapstructure:"HOUSES_URL"`
	NumPages      int    `mapstructure:"NUM_PAGES"`
	NumPerPage    int    `mapstructure:"NUM_PER_PAGE"`
	PhotosDir     string `mapstructure:"PHOTOS_DIR"`
	ClientRetries int    `mapstructure:"CLIENT_RETRIES"`
}

// loadEnvVariables reads configuration from file or environment variables.
func loadEnvVariables(path string) (cfg *config) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("error reading env app.file", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("unmarshal error", err)
	}

	return cfg
}
