package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Rasikrr/core/log"
)

type Variable struct {
	Type     string `yaml:"type"`
	Name     string `yaml:"name"`
	EnvName  string `yaml:"env_name"`
	Value    any    `yaml:"value"`
	Required bool   `yaml:"required"`
}

// nolint: funlen
func (v *Variable) Validate() error {
	ctx := context.Background()
	if v.EnvName != "" {
		if val, ok := os.LookupEnv(v.EnvName); ok {
			v.Value = val
		}
	}

	if !v.Required {
		log.Debugf(ctx, "variable %s is not required", v.Name)
		return nil
	}

	if v.Required && v.Value == nil {
		return fmt.Errorf("variable %s is required but missing", v.Name)
	}

	switch v.Type {
	case "string":
		_, ok := v.Value.(string)
		if !ok {
			return fmt.Errorf("variable %s should be string", v.Name)
		}
	case "int":
		switch val := v.Value.(type) {
		case int:
		case string:
			parsed, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("invalid int for %s: %v", v.Name, err)
			}
			v.Value = parsed
		default:
			return fmt.Errorf("invalid type for int: %s", v.Name)
		}
	case "float":
		switch val := v.Value.(type) {
		case float64:
		case string:
			parsed, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return fmt.Errorf("invalid float for %s: %v", v.Name, err)
			}
			v.Value = parsed
		default:
			return fmt.Errorf("invalid type for float: %s", v.Name)
		}
	case "bool":
		switch val := v.Value.(type) {
		case bool:
		case string:
			parsed, err := strconv.ParseBool(val)
			if err != nil {
				return fmt.Errorf("invalid bool for %s: %v", v.Name, err)
			}
			v.Value = parsed
		default:
			return fmt.Errorf("invalid type for bool: %s", v.Name)
		}
	case "duration":
		switch val := v.Value.(type) {
		case time.Duration:
		case string:
			dur, err := time.ParseDuration(val)
			if err != nil {
				return fmt.Errorf("invalid duration for %s: %v", v.Name, err)
			}
			v.Value = dur
		default:
			return fmt.Errorf("invalid type for duration: %s", v.Name)
		}
	default:
		return fmt.Errorf("unknown type: %s", v.Type)
	}

	return nil
}

type Variables map[string]*Variable

func (v Variables) Validate() error {
	for name, val := range v {
		if err := val.Validate(); err != nil {
			return fmt.Errorf("validation error for %s: %w", name, err)
		}
	}
	return nil
}

func (v Variables) GetString(name string) string {
	val, ok := v[name]
	if !ok {
		return ""
	}
	str, ok := val.Value.(string)
	if !ok {
		return ""
	}
	return str
}

func (v Variables) GetInt(name string) int {
	val, ok := v[name]
	if !ok {
		return 0
	}
	parsed, ok := val.Value.(int)
	if !ok {
		return 0
	}
	return parsed
}

func (v Variables) GetFloat(name string) float64 {
	val, ok := v[name]
	if !ok {
		return 0
	}
	parsed, ok := val.Value.(float64)
	if !ok {
		return 0
	}
	return parsed
}

func (v Variables) GetDuration(name string) time.Duration {
	val, ok := v[name]
	if !ok {
		return 0
	}
	dur, ok := val.Value.(time.Duration)
	if !ok {
		return 0
	}
	return dur
}

func (v Variables) GetBool(name string) bool {
	val, ok := v[name]
	if !ok {
		return false
	}

	if !val.Required {
		return false
	}

	parsed, ok := val.Value.(bool)
	if !ok {
		return false
	}
	return parsed
}
