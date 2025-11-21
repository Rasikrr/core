package redis

import (
	"fmt"
	"time"
)

var (
	errConfigRequired = fmt.Errorf("redis config error")
)

// Config содержит настройки для подключения к Redis
type Config struct {
	Host        string        `yaml:"-" env:"REDIS_HOST"`
	Port        int           `yaml:"-" env:"REDIS_PORT"`
	User        string        `yaml:"-" env:"REDIS_USER"`
	Password    string        `yaml:"-" env:"REDIS_PASSWORD"`
	DB          int           `yaml:"-" env:"REDIS_DB"`
	Required    bool          `yaml:"required"`
	PoolSize    int           `yaml:"pool_size"`
	MinIdle     int           `yaml:"min_idle_conns"`
	MaxIdle     int           `yaml:"max_idle_conns"`
	ReadTimeout time.Duration `yaml:"read_timeout"`
	PrefixKey   string        `yaml:"prefix_key"`
}

// Validate проверяет корректность конфигурации
func (c Config) Validate() error {
	if !c.Required {
		return nil
	}
	if c.PoolSize == 0 {
		return fmt.Errorf("pool_size is empty: %w", errConfigRequired)
	}
	if c.MinIdle == 0 {
		return fmt.Errorf("min_idle_conns is empty: %w", errConfigRequired)
	}
	if c.MaxIdle == 0 {
		return fmt.Errorf("max_idle_conns is empty: %w", errConfigRequired)
	}
	if c.ReadTimeout == 0 {
		return fmt.Errorf("read_timeout is empty: %w", errConfigRequired)
	}
	return nil
}
