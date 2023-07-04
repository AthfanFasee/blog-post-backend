package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort   int    `mapstructure:"SERVER_PORT"`
	DSN          string `mapstructure:"DB_DSN"`
	SmtpHost     string `mapstructure:"SMTP_HOST"`
	SmtpPort     int    `mapstructure:"SMTP_PORT"`
	SmtpUsername string `mapstructure:"SMTP_USERNAME"`
	SmtpPassword string `mapstructure:"SMPT_PASSWORD"`
	SmtpSender   string `mapstructure:"SMTP_SENDER"`
}

var (
	configName = "app"
	configType = "env"
	configPath = "."
)

// Read configuration from environemnt variables.
// Environemnt variables are replaced with correct values in prod environment.
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
