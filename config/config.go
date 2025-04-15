package config

import (
	"errors"
	"fmt"
	"github.com/Rasikrr/core/enum"
	"github.com/Rasikrr/core/interfaces"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
)

const (
	configPathEnv = "CONFIG_PATH"
)

var (
	errConfigNotFound = errors.New("config not found")
)

type Config struct {
	AppName     string           `yaml:"name"`
	Environment enum.Environment `yaml:"environment"`
	HTTP        HTTPConfig       `yaml:"http"`
	GRPC        GRPCConfig       `yaml:"grpc"`
	Postgres    PostgresConfig   `yaml:"postgres"`
	Redis       RedisConfig      `yaml:"redis"`
	NATS        NATSConfig       `yaml:"nats"`
	Variables   Variables        `yaml:"env"`
}

func Parse() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, err
	}

	configFile, ok := os.LookupEnv(configPathEnv)
	if !ok || configFile == "" {
		return Config{}, errConfigNotFound
	}

	var config Config
	if err := cleanenv.ReadConfig(configFile, &config); err != nil {
		return Config{}, err
	}

	if err := config.validate(); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c *Config) validate() error {
	for _, v := range []interfaces.Validatable{
		c.HTTP,
		c.GRPC,
		c.Postgres,
		c.Redis,
		c.NATS,
		c.Variables,
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

func (c *Config) GRPCConfig() GRPCConfig {
	return c.GRPC
}

func (c *Config) HTTPConfig() HTTPConfig {
	return c.HTTP
}

func (c *Config) NATSConfig() NATSConfig {
	return c.NATS
}

func (c *Config) PostgresConfig() PostgresConfig {
	return c.Postgres
}

func (c *Config) RedisConfig() RedisConfig {
	return c.Redis
}
