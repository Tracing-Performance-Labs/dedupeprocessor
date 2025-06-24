package dedupeprocessor

import (
	"context"
	"testing"

	"fmt"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor/processortest"
	conventions "go.opentelemetry.io/otel/semconv/v1.6.1"
)

type testCase struct {
	name            string
	serviceName     string
	inputAttributes map[string]any
	nSpans          int
}

func TestAllKeysArePreserved(t *testing.T) {
	testCases := make([]testCase, 5)

	for i := range len(testCases) {
		testCases[i].inputAttributes = make(map[string]any, i*10)
		testCases[i].nSpans = 1
		addAttributes(&testCases[i], i*10)
	}

	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	tp, err := factory.CreateTraces(context.Background(), processortest.NewNopSettings(Type), cfg, consumertest.NewNop())
	if err != nil {
		t.Fatalf("Failed to create traces processor: %v", err)
	}

	for _, tt := range testCases {
		traces := generateTraceData(tt.serviceName, tt.inputAttributes, tt.nSpans)
		err := tp.ConsumeTraces(context.Background(), traces)
		if err != nil {
			t.Fatalf("Failed to consume traces: %v", err)
		}

		for _, rs := range traces.ResourceSpans().All() {
			for _, ss := range rs.ScopeSpans().All() {
				for _, span := range ss.Spans().All() {
					attrs := span.Attributes()
					for k := range tt.inputAttributes {
						_, exists := attrs.Get(k)
						if !exists {
							t.Fatalf("Attribute %s not found in span attributes after processing", k)
						}
					}
				}
			}
		}
	}
}

func generateTraceData(serviceName string, attrs map[string]any, nSpans int) ptrace.Traces {
	td := ptrace.NewTraces()
	rs := td.ResourceSpans().AppendEmpty()
	if serviceName != "" {
		rs.Resource().Attributes().PutStr(string(conventions.ServiceNameKey), serviceName)
	}

	for range nSpans {
		span := rs.ScopeSpans().AppendEmpty().Spans().AppendEmpty()
		span.Attributes().FromRaw(attrs)
	}

	return td
}

func addAttributes(testCase *testCase, n int) {
	for i := range n {
		testCase.inputAttributes["key_"+fmt.Sprint(i)] = "value_" + fmt.Sprint(i) + "_for_" + testCase.name
	}
}

func BenchmarkTracesProcessor_ForDifferentSpanQuantities(b *testing.B) {
	testCases := []testCase{
		{
			name:            "apply_to_trace_with_no_spans",
			inputAttributes: map[string]any{},
			nSpans:          0,
		},
		{
			name:            "apply_to_trace_with_one_span",
			inputAttributes: map[string]any{},
			nSpans:          1,
		},
		{
			name:            "apply_to_trace_with_ten_spans",
			inputAttributes: map[string]any{},
			nSpans:          10,
		},
		{
			name:            "apply_to_trace_with_fifty_spans",
			inputAttributes: map[string]any{},
			nSpans:          50,
		},
		{
			name:            "apply_to_trace_with_one_hundred_spans",
			inputAttributes: map[string]any{},
			nSpans:          100,
		},
	}

	// Make sure the processor does some work.
	for i := range testCases {
		addAttributes(&testCases[i], 20)
	}

	// TODO: Configure with Redis.

	withRedis(b, func() {
		factory := NewFactory()
		cfg := factory.CreateDefaultConfig()
		oCfg := cfg.(*Config)
		oCfg.TableType = REDIS_TABLE

		tp, err := factory.CreateTraces(context.Background(), processortest.NewNopSettings(Type), cfg, consumertest.NewNop())
		if err != nil {
			b.Fatalf("Failed to create traces processor: %v", err)
		}

		for _, tt := range testCases {
			td := generateTraceData(tt.serviceName, tt.inputAttributes, tt.nSpans)

			b.Run(tt.name, func(b *testing.B) {
				for b.Loop() {
					err = tp.ConsumeTraces(context.Background(), td)
					if err != nil {
						b.Fatalf("Failed to consume traces: %v", err)
					}
				}
			})
		}
	})
}
