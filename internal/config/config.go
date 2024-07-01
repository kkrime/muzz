package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int64  `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	DBname   string `mapstructure:"DB_NAME"`
}

func ReadConfig() (*DBConfig, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigFile(".env")

	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("error reading file = %+v\n", err)
	}

	config := &DBConfig{}
	err = v.Unmarshal(config)

	return config, err
}
