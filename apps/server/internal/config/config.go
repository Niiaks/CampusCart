package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"github.com/rs/zerolog"
)

type Config struct {
	Primary     PrimaryConfig     `koanf:"primary" validate:"required"`
	Server      ServerConfig      `koanf:"server" validate:"required"`
	Database    DatabaseConfig    `koanf:"database" validate:"required"`
	Redis       RedisConfig       `koanf:"redis" validate:"required"`
	Integration IntegrationConfig `koanf:"integration" validate:"required"`
	Auth        AuthConfig        `koanf:"auth_config" validate:"required"`
}

type PrimaryConfig struct {
	Env string `koanf:"env" validate:"required"`
}

type ServerConfig struct {
	Port               string   `koanf:"port" validate:"required"`
	ReadTimeout        int      `koanf:"read_timeout" validate:"required"`
	WriteTimeout       int      `koanf:"write_timeout" validate:"required"`
	IdleTimeout        int      `koanf:"idle_timeout" validate:"required"`
	CorsAllowedOrigins []string `koanf:"cors_allowed_origins" validate:"required"`
}

type DatabaseConfig struct {
	Host            string `koanf:"host" validate:"required"`
	Port            int    `koanf:"port" validate:"required"`
	User            string `koanf:"user" validate:"required"`
	Password        string `koanf:"password"`
	Name            string `koanf:"name" validate:"required"`
	SSLMode         string `koanf:"ssl_mode" validate:"required"`
	MaxOpenConns    int    `koanf:"max_open_conns" validate:"required"`
	MaxIdleConns    int    `koanf:"max_idle_conns" validate:"required"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime int    `koanf:"conn_max_idle_time" validate:"required"`
}

type RedisConfig struct {
	Address string `koanf:"address" validate:"required"`
}

type IntegrationConfig struct {
	ResendApiKey string `koanf:"resend_api_key" validate:"required"`
	SentryDsn    string `koanf:"sentry_dsn" validate:"required"`
}

type AuthConfig struct {
	SecretKey string `koanf:"secret_key" validate:"required"`
}

func LoadConfig() (*Config, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	var k = koanf.New(".")

	err := k.Load(env.Provider("CAMPUS_CART_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "CAMPUS_CART_"))
	}), nil)

	if err != nil {
		logger.Fatal().Err(err).Msg("could not initialize environment variables")
	}

	mainConfig := &Config{}

	err = k.Unmarshal(".", mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not unmarshal config")
	}

	validate := validator.New()

	err = validate.Struct(mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not validate config")
	}

	//initialize sentry
	sentryConfig(mainConfig.Integration.SentryDsn)

	return mainConfig, nil

}

const serverName = "campusCart"

func sentryConfig(dsn string) {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		ServerName:       serverName,
		// Enable printing of sdk debug messages (*remove)
		Debug: true,
		// Adds request headers and ip
		SendDefaultPII: true,
	}); err != nil {
		fmt.Printf("sentry initialization failed")
		os.Exit(1)
	}
}
