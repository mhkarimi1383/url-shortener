package configuration

import (
	"time"

	"github.com/mhkarimi1383/url-shortener/internal/validator"
)

type Config struct {
	ListenAddress                 string        `validate:"required,hostname_port"`
	Migrate                       bool          `validate:"required_without=RunServer"`
	RunServer                     bool          `validate:""`
	DatabaseEngine                string        `validate:"required,oneof=pgx mysql sqlite"`
	DatabaseConnectionString      string        `validate:"required"`
	DatabaseMaxIdleConnections    int           `validate:"required"`
	DatabaseMaxOpenConnections    int           `validate:"required"`
	DatabaseMaxConnectionLifetime time.Duration `validate:"required"`
	JWTSecret                     string        `validate:"required"`
	AddRefererQueryParam          bool          `validate:""`
	RandomGeneratorMax            int           `validate:"min=10000"`
	RootRedirect                  string        `validate:"required"`
	BaseURI                       string        `validate:""`
	RejectRedirectUrls            bool          `validate:""`
	WhiteListHosts                []string      `validate:""`
	RedisAddress                  string        `validate:"omitempty,hostname_port"`
	RedisUsername                 string        `validate:""`
	RedisPassword                 string        `validate:""`
	RedisDatabase                 int           `validate:"min=0"`
	RedisTLS                      bool          `validate:""`
	RedisCacheTTL                 time.Duration `validate:"gt=0"`
	RedisDialTimeout              time.Duration `validate:"gt=0"`
	RedisReadTimeout              time.Duration `validate:"gt=0"`
	RedisWriteTimeout             time.Duration `validate:"gt=0"`
	VisitFlushInterval            time.Duration `validate:"gt=0"`
	VisitBufferMaxEntries         int           `validate:"gt=0"`
}

var CurrentConfig *Config

func SetConfig(config *Config) error {
	err := validator.Validate.Struct(*config)
	if err != nil {
		return err
	}

	CurrentConfig = config
	return nil
}
