package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int64  `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	DBname   string `mapstructure:"DB_NAME"`
	TokenKey string `mapstructure:"TOKEN_KEY"`
}

func ReadConfig() (*Config, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigFile(".env")

	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("error reading file = %+v\n", err)
	}

	config := &Config{}
	err = v.Unmarshal(config)

	return config, err
}
