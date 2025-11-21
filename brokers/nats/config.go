package nats

import (
	"errors"
)

var (
	errNATSConfigRequired = errors.New("nats config error")
)

type Config struct {
	Required bool   `yaml:"required"`
	DSN      string `yaml:"-" env:"NATS_DSN"`
	Queue    string `yaml:"queue"`
}

func (c Config) Validate() error {
	if !c.Required {
		return nil
	}
	if c.DSN == "" {
		return errNATSConfigRequired
	}
	return nil
}
