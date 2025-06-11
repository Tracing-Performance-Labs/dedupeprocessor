package dedupeprocessor

import (
	"errors"

	"go.opentelemetry.io/collector/component"
	"slices"
)

type TableType string

type Config struct {
	tableType TableType `mapstructure:"table-type"`
}

const (
	REDIS_TABLE TableType = "redis"
)

var (
	validTables = []TableType{REDIS_TABLE}
)

func createDefaultConfig() component.Config {
	return &Config{
		tableType: REDIS_TABLE,
	}
}

func (c *Config) Validate() error {
	if slices.Contains(validTables, c.tableType) {
		return nil
	}
	return errors.New("invalid table type: " + string(c.tableType))
}
