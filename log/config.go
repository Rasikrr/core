package log

import "github.com/Rasikrr/core/enum"

type Config struct {
	Level     enum.LogLevel `yaml:"level"`
	AddSource bool          `yaml:"add_source"`
}
