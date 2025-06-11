package dedupeprocessor

import (
	"context"
	"errors"

	"github.com/Tracing-Performance-Labs/go-dedupe"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
)

func NewFactory() processor.Factory {
	return processor.NewFactory(
		component.MustNewType("dedupeprocessor"),
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

	var codec *dedupe.Codec

	switch oCfg.TableType {
	case REDIS_TABLE:
		codec = dedupe.NewCodec(dedupe.WithRedisTable())
	default:
		return nil, errors.New("invalid table type: " + string(oCfg.TableType))
	}

	return newTracesProcessor(
		codec,
		ctx,
		set,
		cfg.(*Config),
		nextConsumer)
}
