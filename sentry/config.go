package sentry

import (
	"errors"
)

type Config struct {
	Enabled          bool    `yaml:"enabled"`
	DSN              string  `env:"SENTRY_DSN"`
	SampleRate       float64 `yaml:"sample_rate"`
	Tracing          bool    `yaml:"tracing"`
	TracesSampleRate float64 `yaml:"traces_sample_rate"`
	Debug            bool    `yaml:"debug"`
	EnableLogs       bool    `yaml:"enable_logs"`
}

func (c Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.DSN == "" {
		return errors.New("sentry: missing DSN")
	}
	if c.SampleRate < 0 || c.SampleRate > 1 {
		return errors.New("sentry: invalid sample rate")
	}
	if c.TracesSampleRate < 0 || c.TracesSampleRate > 1 {
		return errors.New("sentry: invalid traces sample rate")
	}
	return nil
}
