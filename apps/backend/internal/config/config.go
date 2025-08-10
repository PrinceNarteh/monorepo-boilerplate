// Package config contains all the configuration for the application
package config

import (
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload" // Load .env file automatically
	env "github.com/knadh/koanf/providers/env/v2"
	koanf "github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"

	"github.com/PrinceNarteh/go-boilerplate/internal/libs"
)

// Config for the application
type Config struct {
	Auth     AuthConfig     `koanf:"auth"     validate:"required"`
	Core     CoreConfig     `koanf:"core"     validate:"required"`
	Database DatabaseConfig `koanf:"database" validate:"required"`
	Redis    RedisConfig    `koanf:"redis"    validate:"required"`
	Server   ServerConfig   `koanf:"server"   validate:"required"`
}

// CoreConfig contains core configuration for the application
type CoreConfig struct {
	Env string `koanf:"env" validate:"required"`
}

// ServerConfig contains configuration for the server
type ServerConfig struct {
	Port               string   `koanf:"port"                 validate:"required"`
	ReadTimeout        int      `koanf:"read_timeout"         validate:"required"`
	WriteTimeout       int      `koanf:"write_timeout"        validate:"required"`
	IdleTimeout        int      `koanf:"idle_timeout"         validate:"required"`
	CORSAllowedOrigins []string `koanf:"cors_allowed_origins" validate:"required"`
}

// RedisConfig contains configuration for Redis
type RedisConfig struct {
	Address string `koanf:"address" validate:"required"`
}

// DatabaseConfig contains configuration for database
type DatabaseConfig struct {
	Host            string `koanf:"host"              validate:"required"`
	Port            string `koanf:"port"              validate:"required"`
	User            string `koanf:"user"              validate:"required"`
	Password        string `koanf:"password"          validate:"required"`
	Name            string `koanf:"name"              validate:"required"`
	SSLMode         string `koanf:"ssl_mode"          validate:"required"`
	MaxOpenConns    string `koanf:"max_open_conns"    validate:"required"`
	MaxIdleConns    string `koanf:"max_idle_conns"    validate:"required"`
	ConnMaxLifetime string `koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdletime string `koanf:"conn_max_idletime" validate:"required"`
}

// AuthConfig contains configuration for authentication
type AuthConfig struct {
	SecretKey string `koanf:"secret_key" validate:"required"`
}

// LoadConfig loads the configuration from a file or environment variables
// and returns a Config instance. It uses the koanf library for configuration management.
func LoadConfig() (*Config, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	k := koanf.New(".")

	if err := k.Load(env.Provider("API_", env.Opt{
		Prefix: "API_",
		TransformFunc: func(k, v string) (string, any) {
			// Transform the key.
			k = strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(k, "MYVAR_")), "_", ".")

			// Transform the value into slices, if they contain spaces.
			// Eg: MYVAR_TAGS="foo bar baz" -> tags: ["foo", "bar", "baz"]
			// This is to demonstrate that string values can be transformed to any type
			// where necessary.
			if strings.Contains(v, " ") {
				return k, strings.Split(v, " ")
			}

			return k, v
		},
	}), nil); err != nil {
		logger.Fatal().Err(err).Msg("could not load initial env variables")
	}

	mainConfig := &Config{}
	if err := k.Unmarshal("", mainConfig); err != nil {
		logger.Fatal().Err(err).Msg("could not unmarshal main config")
	}

	if err := libs.ValidateStruct(mainConfig); err != nil {
		logger.Fatal().Err(nil).Fields(err)
	}

	return mainConfig, nil
}
