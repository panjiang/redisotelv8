package redisotel_test

import (
	"time"

	"git.shiyou.kingsoft.com/gem/pkg/redisotel"
	"github.com/go-redis/redis/v8"
)

func ExampleExportMetricsForPrometheus() {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"localhost:7000"},
	})

	if err := redisotel.ExportMetricsForPrometheus(rdb, redisotel.WithSlowDur(time.Second)); err != nil {
		panic(err)
	}

	// Output:
}

func ExampleInstrumentMetrics() {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"localhost:7000"},
	})

	mp, err := redisotel.NewPrometheusMeterProvider()
	if err != nil {
		panic(err)
	}

	if err := redisotel.InstrumentMetrics(rdb, redisotel.WithMeterProvider(mp)); err != nil {
		panic(err)
	}

	// Output:
}
