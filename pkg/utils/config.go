package utils

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

var (
	instance *Config
	once     sync.Once
)

func SetConfigPath(path string) {
	viper.AddConfigPath(path)
}

// GetConfig returns the configuration instance using once.Do to ensure that the configuration is loaded only once
func GetConfig() *Config {

	once.Do(func() {
		var err error

		log.Info().Msg("Loading config...")

		instance = &Config{}

		viper.AddConfigPath(".")
		viper.SetConfigName("app")
		viper.SetConfigType("env")

		viper.AutomaticEnv()

		err = viper.ReadInConfig()

		if err != nil {
			log.Panic().Err(err).Msg("Error loading config")
		}

		err = viper.Unmarshal(instance)

		if err != nil {
			log.Panic().Err(err).Msg("Error unmarshalling config")
		}
	})

	return instance
}
