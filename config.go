package dedupeprocessor

import (
	"errors"

	"go.opentelemetry.io/collector/component"
	"slices"
)

type TableType string

type Config struct {
	TableType TableType `mapstructure:"table_type"`
}

const (
	REDIS_TABLE TableType = "redis"
)

var (
	validTables = []TableType{REDIS_TABLE}
)

func createDefaultConfig() component.Config {
	return &Config{
		TableType: REDIS_TABLE,
	}
}

func (c *Config) Validate() error {
	if slices.Contains(validTables, c.TableType) {
		return nil
	}
	return errors.New("invalid table type: " + string(c.TableType))
}
