package dedupeprocessor

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

type traceProcessor struct {
	// TODO: Put Codec dependency here.
}

func newTracesProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg *Config,
	nextConsumer consumer.Traces,
) (processor.Traces, error) {
	tp := &traceProcessor{}
	return processorhelper.NewTraces(
		ctx,
		set,
		cfg,
		nextConsumer,
		tp.processTraces,
		processorhelper.WithCapabilities(consumer.Capabilities{MutatesData: true}))
}

func (tp *traceProcessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		rs := rss.At(i)
		ilss := rs.ScopeSpans()
		for j := 0; j < ilss.Len(); j++ {
			ils := ilss.At(j)
			spans := ils.Spans()
			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)

				// Process the span attributes
				attrs := span.Attributes()

				for key, value := range attrs.All() {
					// Deduplicate the key
					slog.Info("deduplicating key", "key", key)

					// Deduplicate the value
					if value.Type() == pcommon.ValueTypeStr {
						slog.Info("deduplicating value", "value", value.Str())
					}
				}

			}
		}
	}
	return td, nil
}
