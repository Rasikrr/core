package config

import "fmt"

var (
	errGRPCConfigRequired = fmt.Errorf("grpc config error")
)

type GRPCConfig struct {
	Host     string `yaml:"host" env:"GRPC_HOST" env-default:"0.0.0.0"`
	Port     int    `yaml:"port" env:"GRPC_PORT" env-default:"3000"`
	Required bool   `yaml:"required" env:"GRPC_REQUIRED" env-default:"false"`
}

func (c GRPCConfig) Validate() error {
	if !c.Required {
		return nil
	}
	if c.Port == 0 {
		return fmt.Errorf("port is empty: %w", errGRPCConfigRequired)
	}
	return nil
}
