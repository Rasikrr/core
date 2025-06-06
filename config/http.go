package config

import "fmt"

var (
	errHTTPConfigRequired = fmt.Errorf("http config error")
)

type HTTPConfig struct {
	Name     string `yaml:"name" env-default:"" example:"metrics or something"`
	Host     string `yaml:"host" env:"HTTP_HOST" env-default:"0.0.0.0"`
	Port     string `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	Required bool   `yaml:"required" env:"HTTP_REQUIRED" env-default:"false"`
}

func (c HTTPConfig) Validate() error {
	if !c.Required {
		return nil
	}
	if c.Host == "" {
		return fmt.Errorf("host is empty: %w", errHTTPConfigRequired)
	}
	if c.Port == "" {
		return fmt.Errorf("port is empty: %w", errHTTPConfigRequired)
	}
	return nil
}
