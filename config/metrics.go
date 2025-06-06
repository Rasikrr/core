package config

import (
	"errors"
	"fmt"
)

var (
	errMetricsConfig = errors.New("metrics config error")
)

type Metrics struct {
	Enabled    bool                 `yaml:"enabled" env:"METRICS_ENABLED" env-default:"false"`
	Namespace  string               `yaml:"namespace" env:"METRICS_NAMESPACE"`
	Prometheus PrometheusExportConf `yaml:"prometheus" env:"METRICS_PROMETHEUS"`
}

type PrometheusExportConf struct {
	Port string `yaml:"port" env:"METRICS_PROMETHEUS_PORT" env-default:"9100"`
}

func (m Metrics) Validate() error {
	if !m.Enabled {
		return nil
	}
	if m.Namespace == "" {
		return fmt.Errorf("namespace is empty: %w", errMetricsConfig)
	}
	if m.Prometheus.Port == "" {
		return fmt.Errorf("prometheus port is empty: %w", errMetricsConfig)
	}
	return nil
}
