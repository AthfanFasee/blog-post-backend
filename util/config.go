// Package util provides various utility functions and helpers for the application.
package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort         int    `mapstructure:"SERVER_PORT"`
	DSN                string `mapstructure:"DB_DSN"`
	SmtpHost           string `mapstructure:"SMTP_HOST"`
	SmtpPort           int    `mapstructure:"SMTP_PORT"`
	SmtpUsername       string `mapstructure:"SMTP_USERNAME"`
	SmtpPassword       string `mapstructure:"SMPT_PASSWORD"`
	SmtpSender         string `mapstructure:"SMTP_SENDER"`
	CorsTrustedOrigins string `mapstructure:"CORS_TRUSTED_ORIGINS"`
}

var (
	configName = "app"
	configType = "env"
	configPath = "."
)

// LoadEnv reads configuration from environemnt variables.
func LoadEnv() (config Config, err error) {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	viper.AddConfigPath(configPath)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
