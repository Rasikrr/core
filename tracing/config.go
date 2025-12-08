package tracing

import "github.com/cockroachdb/errors"

type Config struct {
	Enabled bool   `yaml:"enabled"`
	DSN     string `yaml:"-" env:"TRACING_EXPORTER_DSN"`
}

func (c *Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.DSN == "" {
		return errors.New("dsn is required")
	}
	return nil
}
