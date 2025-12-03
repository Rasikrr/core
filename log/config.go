package log

import "github.com/Rasikrr/core/enum"

type Config struct {
	Level     enum.LogLevel  `yaml:"level"`
	Format    enum.LogFormat `yaml:"format"`
	AddSource bool           `yaml:"add_source"`
}
