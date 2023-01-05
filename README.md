# OpenTelemetry instrumentation metrics for go-redis v8

## Installation

```shell
go get github.com/panjiang/redisotelv8
```

## Usage


```go
import (
	"time"

	redisotel "github.com/panjiang/redisotelv8"
	"github.com/go-redis/redis/v8"
)

rdb := redis.NewUniversalClient(&redis.UniversalOptions{
    Addrs: []string{"localhost:7000"},
})

if err := redisotel.ExportMetricsForPrometheus(rdb, redisotel.WithSlowDur(time.Second)); err != nil {
    panic(err)
}
```

Or

```go
import (
	"time"

	redisotel "github.com/panjiang/redisotelv8"
	"github.com/go-redis/redis/v8"
)

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
```

## Metrics

```plaintext
# HELP db_client_connections_idle_min The minimum number of idle open connections allowed
# TYPE db_client_connections_idle_min gauge
db_client_connections_idle_min{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000"} 0
db_client_connections_idle_min{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7001"} 0
db_client_connections_idle_min{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7002"} 0
# HELP db_client_connections_max The maximum number of open connections allowed
# TYPE db_client_connections_max gauge
db_client_connections_max{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000"} 40
db_client_connections_max{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7001"} 40
db_client_connections_max{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7002"} 40
# HELP db_client_connections_timeouts The number of connection timeouts that have occurred trying to obtain a connection from the pool
# TYPE db_client_connections_timeouts gauge
db_client_connections_timeouts{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000"} 0
db_client_connections_timeouts{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7001"} 0
db_client_connections_timeouts{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7002"} 0
# HELP db_client_connections_usage The number of connections that are currently in state described by the state attribute
# TYPE db_client_connections_usage gauge
db_client_connections_usage{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",state="idle"} 1
db_client_connections_usage{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",state="used"} 0
db_client_connections_usage{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7001",state="idle"} 1
db_client_connections_usage{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7001",state="used"} 1
db_client_connections_usage{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7002",state="idle"} 2
db_client_connections_usage{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7002",state="used"} 0
# HELP db_client_connections_use_time_milliseconds The time between borrowing a connection and returning it to the pool.
# TYPE db_client_connections_use_time_milliseconds histogram
db_client_connections_use_time_milliseconds_bucket{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command",le="0"} 5
db_client_connections_use_time_milliseconds_bucket{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command",le="5"} 5
db_client_connections_use_time_milliseconds_bucket{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command",le="10"} 5
db_client_connections_use_time_milliseconds_bucket{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command",le="25"} 5
db_client_connections_use_time_milliseconds_bucket{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command",le="50"} 5
db_client_connections_use_time_milliseconds_bucket{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command",le="75"} 5
db_client_connections_use_time_milliseconds_bucket{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command",le="100"} 5
db_client_connections_use_time_milliseconds_bucket{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command",le="250"} 5
db_client_connections_use_time_milliseconds_bucket{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command",le="500"} 5
db_client_connections_use_time_milliseconds_bucket{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command",le="750"} 5
db_client_connections_use_time_milliseconds_bucket{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command",le="1000"} 5
db_client_connections_use_time_milliseconds_bucket{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command",le="2500"} 5
db_client_connections_use_time_milliseconds_bucket{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command",le="5000"} 5
db_client_connections_use_time_milliseconds_bucket{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command",le="7500"} 5
db_client_connections_use_time_milliseconds_bucket{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command",le="10000"} 5
db_client_connections_use_time_milliseconds_bucket{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command",le="+Inf"} 5
db_client_connections_use_time_milliseconds_sum{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command"} 0
db_client_connections_use_time_milliseconds_count{db_system="redis",otel_scope_name="redisotel",otel_scope_version="8.11.5",pool_name="127.0.0.1:7000",status="ok",type="command"} 5
```