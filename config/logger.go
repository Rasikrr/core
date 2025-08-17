package config

import "github.com/Rasikrr/core/enum"

type LoggerConfig struct {
	Level     enum.LogLevel `yaml:"level"`
	AddSource bool          `yaml:"add_source"`
}
