package configuration

import (
	"time"

	"github.com/mhkarimi1383/url-shortener/internal/validator"
)

type Config struct {
	ListenAddress                 string `validate:"hostname_port"`
	Migrate                       bool   `validate:"required_without=RunServer"`
	RunServer                     bool
	DatabaseEngine                string `validate:"oneof=pgx mysql sqlite"`
	DatabaseConnectionString      string
	DatabaseMaxIdleConnections    int
	DatabaseMaxOpenConnections    int
	DatabaseMaxConnectionLifetime time.Duration
}

var currentConfig Config

func SetConfig(config *Config) error {
	if err := validator.Validate.Struct(*config); err != nil {
		return err
	}

	currentConfig = *config
	return nil
}

func GetConfig() *Config {
	return &currentConfig
}
