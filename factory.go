package dedupeprocessor

import (
	"context"

	"github.com/Tracing-Performance-Labs/go-dedupe"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
)

var Type = component.MustNewType("dedupeprocessor")

func NewFactory() processor.Factory {
	return processor.NewFactory(
		Type,
		createDefaultConfig,
		processor.WithTraces(createTracesProcessor, component.StabilityLevelAlpha))
}

func createTracesProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer consumer.Traces,
) (processor.Traces, error) {
	oCfg := cfg.(*Config)
	codec := dedupe.NewCodec(withTableType(oCfg), withReprType(oCfg))
	return newTracesProcessor(
		codec,
		ctx,
		set,
		cfg.(*Config),
		nextConsumer)
}

func withTableType(oCfg *Config) dedupe.CodecOption {
	switch oCfg.TableType {
	case REDIS_TABLE:
		return dedupe.WithRedisTable()
	case MEMORY_TABLE:
		return dedupe.WithMemoryTable()
	default:
		return nil
	}
}

func withReprType(oCfg *Config) dedupe.CodecOption {
	switch oCfg.ReprType {
	case MURMUR_REPR:
		return dedupe.WithMurmurRepr()
	default:
		return dedupe.WithDefaultObjectRepr()
	}
}
