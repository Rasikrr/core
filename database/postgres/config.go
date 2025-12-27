package postgres

import (
	"fmt"
	"time"
)

var (
	errConfigRequired = fmt.Errorf("postgres config error")
)

type Config struct {
	DSN                 string        `yaml:"-" env:"POSTGRES_DSN"`
	Required            bool          `yaml:"required"`
	MaxConns            int           `yaml:"max_conns"`
	MinConns            int           `yaml:"min_conns"`
	MaxIdleConnIdleTime time.Duration `yaml:"max_idle_conn_time"`
}

func (c Config) Validate() error {
	if !c.Required {
		return nil
	}
	if c.MaxConns == 0 {
		return fmt.Errorf("max_conns is empty: %w", errConfigRequired)
	}
	if c.MinConns == 0 {
		return fmt.Errorf("min_conns is empty: %w", errConfigRequired)
	}
	if c.MaxIdleConnIdleTime == 0 {
		return fmt.Errorf("max_idle_conn_time is empty: %w", errConfigRequired)
	}
	return nil
}
