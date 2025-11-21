package metrics

import (
	"errors"
	"fmt"
)

var (
	errConfigRequired = errors.New("metrics config error")
)

type Config struct {
	Enabled    bool                 `yaml:"enabled"`
	Namespace  string               `yaml:"namespace"`
	Prometheus PrometheusExportConf `yaml:"prometheus"`
}

type PrometheusExportConf struct {
	Port string `yaml:"port" env:"METRICS_PROMETHEUS_PORT" env-default:"9100"`
}

func (m Config) Validate() error {
	if !m.Enabled {
		return nil
	}
	if m.Namespace == "" {
		return fmt.Errorf("namespace is empty: %w", errConfigRequired)
	}
	if m.Prometheus.Port == "" {
		return fmt.Errorf("prometheus port is empty: %w", errConfigRequired)
	}
	return nil
}
