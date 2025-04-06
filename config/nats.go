package config

import (
	"errors"
)

var (
	errNATSConfigRequired = errors.New("nats config error")
)

type NATSConfig struct {
	Required bool   `yaml:"required"`
	DSN      string `yaml:"dsn"`
	Queue    string `yaml:"queue"`
}

func (c NATSConfig) Validate() error {
	if !c.Required {
		return nil
	}
	if c.DSN == "" {
		return errNATSConfigRequired
	}
	return nil
}
