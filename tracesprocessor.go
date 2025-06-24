package dedupeprocessor

import (
	"context"
	"sync"

	"github.com/Tracing-Performance-Labs/go-dedupe"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

type traceProcessor struct {
	codec *dedupe.Codec
}

func newTracesProcessor(
	codec *dedupe.Codec,
	ctx context.Context,
	set processor.Settings,
	cfg *Config,
	nextConsumer consumer.Traces,
) (processor.Traces, error) {
	tp := &traceProcessor{codec: codec}
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
	for i := range rss.Len() {
		rs := rss.At(i)
		ilss := rs.ScopeSpans()
		for j := range ilss.Len() {
			ils := ilss.At(j)
			spans := ils.Spans()

			var wg sync.WaitGroup
			wg.Add(spans.Len())

			for k := range spans.Len() {
				span := spans.At(k)
				attrs := span.Attributes()

				for _, value := range attrs.All() {
					if value.Type() == pcommon.ValueTypeStr {
						newValue := pcommon.NewValueStr(tp.codec.Encode(value.Str()))
						newValue.MoveTo(value)
					}
				}
			}
		}
	}
	return td, nil
}
