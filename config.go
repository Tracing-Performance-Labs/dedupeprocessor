package dedupeprocessor

import (
	"errors"

	"go.opentelemetry.io/collector/component"
	"slices"
)

type TableType string

type ReprType string

type Config struct {
	TableType TableType `mapstructure:"table_type"`
	ReprType  ReprType  `mapstructure:"repr_type"`
}

const (
	REDIS_TABLE TableType = "redis"
)

const (
	DEFAULT_REPR ReprType = "default"
	MURMUR_REPR           = "murmur"
)

var (
	validTables = []TableType{REDIS_TABLE}
	validReprs  = []ReprType{DEFAULT_REPR, MURMUR_REPR}
)

func createDefaultConfig() component.Config {
	return &Config{
		TableType: REDIS_TABLE,
		ReprType:  DEFAULT_REPR,
	}
}

func (c *Config) Validate() error {
	if slices.Contains(validTables, c.TableType) && slices.Contains(validReprs, c.ReprType) {
		return nil
	}

	if !slices.Contains(validReprs, c.ReprType) {
		return errors.New("invalid representation type: " + string(c.ReprType))
	}

	return errors.New("invalid table type: " + string(c.TableType))
}
