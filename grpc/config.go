package grpc

import "fmt"

var (
	errConfigRequired = fmt.Errorf("grpc config error")
)

// Config содержит настройки для gRPC сервера
type Config struct {
	Host     string `yaml:"host" env:"GRPC_HOST" env-default:"0.0.0.0"`
	Port     int    `yaml:"port" env:"GRPC_PORT" env-default:"3000"`
	Required bool   `yaml:"required" env:"GRPC_REQUIRED" env-default:"false"`
}

// Validate проверяет корректность конфигурации
func (c Config) Validate() error {
	if !c.Required {
		return nil
	}
	if c.Port == 0 {
		return fmt.Errorf("port is empty: %w", errConfigRequired)
	}
	return nil
}
