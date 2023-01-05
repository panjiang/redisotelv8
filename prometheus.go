package redisotel

import (
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

func NewPrometheusMeterProvider(opts ...prometheus.Option) (*metric.MeterProvider, error) {
	exporter, err := prometheus.New(opts...)
	if err != nil {
		return nil, err
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	return provider, nil
}

// Export redis metrics with otel meter for default prometheus register.
func ExportMetricsForPrometheus(rdb redis.UniversalClient, opts ...MetricsOption) error {
	mp, err := NewPrometheusMeterProvider()
	if err != nil {
		return err
	}

	opts = append([]MetricsOption{WithMeterProvider(mp)}, opts...)
	return InstrumentMetrics(rdb, opts...)
}
