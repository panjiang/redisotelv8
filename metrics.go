package redisotel

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.uber.org/zap"
)

const (
	instrumName = "redisotel"
)

func InstrumentMetrics(rdb redis.UniversalClient, opts ...MetricsOption) error {
	baseOpts := make([]baseOption, len(opts))
	for i, opt := range opts {
		baseOpts[i] = opt
	}
	conf := newConfig(baseOpts...)

	if conf.meter == nil {
		conf.meter = conf.mp.Meter(
			instrumName,
			metric.WithInstrumentationVersion(redis.Version()),
		)
	}

	if conf.slowDur == 0 {
		conf.slowDur = time.Second
	}

	switch rdb := rdb.(type) {
	case *redis.Client:
		if conf.poolName == "" {
			opt := rdb.Options()
			conf.poolName = opt.Addr
		}
		conf.attrs = append(conf.attrs, attribute.String("pool.name", conf.poolName))

		if err := reportPoolStats(rdb, conf); err != nil {
			return err
		}
		if err := addMetricsHook(rdb, conf); err != nil {
			return err
		}
		return nil
	case *redis.ClusterClient:
		return rdb.ForEachMaster(context.Background(), func(ctx context.Context, rdb *redis.Client) error {
			confCopy := &config{}
			*confCopy = *conf
			if confCopy.poolName == "" {
				opt := rdb.Options()
				confCopy.poolName = opt.Addr
			}
			confCopy.attrs = append(confCopy.attrs, attribute.String("pool.name", confCopy.poolName))

			if err := reportPoolStats(rdb, confCopy); err != nil {
				otel.Handle(err)
			}
			if err := addMetricsHook(rdb, confCopy); err != nil {
				otel.Handle(err)
			}
			return nil
		})
	case *redis.Ring:
		return rdb.ForEachShard(context.Background(), func(ctx context.Context, rdb *redis.Client) error {
			confCopy := &config{}
			*confCopy = *conf
			if confCopy.poolName == "" {
				opt := rdb.Options()
				confCopy.poolName = opt.Addr
			}
			confCopy.attrs = append(confCopy.attrs, attribute.String("pool.name", confCopy.poolName))

			if err := reportPoolStats(rdb, confCopy); err != nil {
				otel.Handle(err)
			}
			if err := addMetricsHook(rdb, confCopy); err != nil {
				otel.Handle(err)
			}
			return nil
		})
	default:
		return fmt.Errorf("redisotel: %T not supported", rdb)
	}
}

func reportPoolStats(rdb *redis.Client, conf *config) error {
	labels := conf.attrs
	idleAttrs := append(labels, attribute.String("state", "idle"))
	usedAttrs := append(labels, attribute.String("state", "used"))

	idleMax, err := conf.meter.AsyncInt64().UpDownCounter(
		"db.client.connections.idle.max",
		instrument.WithDescription("The maximum number of idle open connections allowed"),
	)
	if err != nil {
		return err
	}

	idleMin, err := conf.meter.AsyncInt64().UpDownCounter(
		"db.client.connections.idle.min",
		instrument.WithDescription("The minimum number of idle open connections allowed"),
	)
	if err != nil {
		return err
	}

	connsMax, err := conf.meter.AsyncInt64().UpDownCounter(
		"db.client.connections.max",
		instrument.WithDescription("The maximum number of open connections allowed"),
	)
	if err != nil {
		return err
	}

	usage, err := conf.meter.AsyncInt64().UpDownCounter(
		"db.client.connections.usage",
		instrument.WithDescription("The number of connections that are currently in state described by the state attribute"),
	)
	if err != nil {
		return err
	}

	timeouts, err := conf.meter.AsyncInt64().UpDownCounter(
		"db.client.connections.timeouts",
		instrument.WithDescription("The number of connection timeouts that have occurred trying to obtain a connection from the pool"),
	)
	if err != nil {
		return err
	}

	redisConf := rdb.Options()
	return conf.meter.RegisterCallback(
		[]instrument.Asynchronous{
			idleMax,
			idleMin,
			connsMax,
			usage,
			timeouts,
		},
		func(ctx context.Context) {
			stats := rdb.PoolStats()

			idleMin.Observe(ctx, int64(redisConf.MinIdleConns), labels...)
			connsMax.Observe(ctx, int64(redisConf.PoolSize), labels...)

			usage.Observe(ctx, int64(stats.IdleConns), idleAttrs...)
			usage.Observe(ctx, int64(stats.TotalConns-stats.IdleConns), usedAttrs...)

			timeouts.Observe(ctx, int64(stats.Timeouts), labels...)
		},
	)
}

func addMetricsHook(rdb *redis.Client, conf *config) error {
	useTime, err := conf.meter.SyncInt64().Histogram(
		"db.client.connections.use_time",
		instrument.WithDescription("The time between borrowing a connection and returning it to the pool."),
		instrument.WithUnit("ms"),
	)
	if err != nil {
		return err
	}

	rdb.AddHook(&metricsHook{
		logger:       zap.L(),
		useTime:      useTime,
		attrs:        conf.attrs,
		slowDuration: conf.slowDur,
	})
	return nil
}

type (
	contextKey struct{}
)

var (
	startKey contextKey
)

type metricsHook struct {
	useTime      syncint64.Histogram
	attrs        []attribute.KeyValue
	logger       *zap.Logger
	slowDuration time.Duration
}

var _ redis.Hook = (*metricsHook)(nil)

func (mh *metricsHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startKey, time.Now()), nil
}

func (mh *metricsHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if start, ok := ctx.Value(startKey).(time.Time); ok {
		dur := time.Since(start)

		attrs := make([]attribute.KeyValue, 0, len(mh.attrs)+2)
		attrs = append(attrs, mh.attrs...)
		attrs = append(attrs, attribute.String("type", "command"))

		statusAttr, ok := statusOkAttr(cmd.Err())
		if !ok {
			mh.logger.Error("Command error", zap.String("cmd", cmd.String()))
		}
		attrs = append(attrs, statusAttr)
		mh.useTime.Record(ctx, dur.Milliseconds(), attrs...)

		if dur >= mh.slowDuration {
			mh.logger.Warn("Command slow", zap.Duration("dur", dur), zap.String("cmd", cmd.String()))
		}
	}
	return nil
}

func (mh *metricsHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startKey, time.Now()), nil
}

func (mh *metricsHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if start, ok := ctx.Value(startKey).(time.Time); ok {
		dur := time.Since(start)

		attrs := make([]attribute.KeyValue, 0, len(mh.attrs)+2)
		attrs = append(attrs, mh.attrs...)
		attrs = append(attrs, attribute.String("type", "pipeline"))

		firstCmd := cmds[0]
		statusAttr, ok := statusOkAttr(firstCmd.Err())
		if !ok {
			mh.logger.Error("Pipeline error", zap.String("cmd", firstCmd.String()))
		}
		attrs = append(attrs, statusAttr)

		if dur >= mh.slowDuration {
			mh.logger.Warn("Pipeline slow", zap.Duration("dur", dur), zap.String("cmd", firstCmd.String()))
		}

		mh.useTime.Record(ctx, dur.Milliseconds(), attrs...)
	}
	return nil
}

func milliseconds(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}

func statusOkAttr(err error) (attribute.KeyValue, bool) {
	if err != nil && err != redis.Nil {
		return attribute.String("status", "error"), false
	}
	return attribute.String("status", "ok"), true
}
