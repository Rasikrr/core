package config

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Rasikrr/core/brokers/nats"
	"github.com/Rasikrr/core/cache/redis"
	"github.com/Rasikrr/core/config/appenv"
	"github.com/Rasikrr/core/database"
	"github.com/Rasikrr/core/enum"
	"github.com/Rasikrr/core/grpc"
	"github.com/Rasikrr/core/http"
	"github.com/Rasikrr/core/interfaces"
	"github.com/Rasikrr/core/log"
	"github.com/Rasikrr/core/metrics"
	"github.com/Rasikrr/core/sentry"
	"github.com/Rasikrr/core/tracing"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var (
	errConfigNotFound = errors.New("config not found")
)

type Config struct {
	AppName     string           `yaml:"name"`
	Environment enum.Environment `env:"ENVIRONMENT"`
	Version     string           `desc:"git tag -> commit hash -> unknown"`
	Variables   Variables        `yaml:"env"`

	Logger   log.Config      `yaml:"log"`
	HTTP     http.Config     `yaml:"http"`
	GRPC     grpc.Config     `yaml:"grpc"`
	Postgres database.Config `yaml:"postgres"`
	Redis    redis.Config    `yaml:"redis"`
	NATS     nats.Config     `yaml:"nats"`
	Metrics  metrics.Config  `yaml:"metrics"`
	Tracing  tracing.Config  `yaml:"tracing"`
	Sentry   sentry.Config   `yaml:"sentry"`
}

func Parse() (Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Warnf(context.Background(), "failed to load .env file: %v", err)
	}

	configFile, ok := os.LookupEnv(appenv.ConfigPathEnv)
	if !ok || configFile == "" {
		return Config{}, errConfigNotFound
	}

	var config Config
	if err := cleanenv.ReadConfig(configFile, &config); err != nil {
		return Config{}, err
	}

	config.setAppVersion()

	if err := config.validate(); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c *Config) validate() error {
	for _, v := range []interfaces.Validatable{
		c.Sentry,
		c.HTTP,
		c.GRPC,
		c.Postgres,
		c.Redis,
		c.NATS,
		c.Variables,
		c.Metrics,
	} {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("error while validating config: %w", err)
		}
	}
	return nil
}

func (c *Config) Env() enum.Environment {
	return c.Environment
}

func (c *Config) Name() string {
	return c.AppName
}

func (c *Config) AppVersion() string {
	return c.Version
}

// setAppVersion устанавливает версию приложения из Git
// Приоритет: Git tag → Git commit hash → "unknown"
func (c *Config) setAppVersion() {
	c.Version = buildVersion()
}
