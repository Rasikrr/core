package config

import (
	"fmt"
	"time"
)

var (
	errPostgresConfigRequired = fmt.Errorf("postgres config error")
)

type PostgresConfig struct {
	DSN                 string        `yaml:"dsn" env:"POSTGRES_DSN"`
	Required            bool          `yaml:"required"`
	MaxConns            int           `yaml:"max_conns"`
	MinConns            int           `yaml:"min_conns"`
	MaxIdleConnIdleTime time.Duration `yaml:"max_idle_conn_time"`
}

func (c PostgresConfig) Validate() error {
	if !c.Required {
		return nil
	}
	if c.MaxConns == 0 {
		return fmt.Errorf("max_conns is empty: %w", errPostgresConfigRequired)
	}
	if c.MinConns == 0 {
		return fmt.Errorf("min_conns is empty: %w", errPostgresConfigRequired)
	}
	if c.MaxIdleConnIdleTime == 0 {
		return fmt.Errorf("max_idle_conn_time is empty: %w", errPostgresConfigRequired)
	}
	return nil
}
