# Service 1: Crud
# Service 2: Monitor
ES + CQRS


```dockerfile
docker-compose up -d
```

equipment-api

```dockerfile
LOG:
DONT use it
https://github.com/Melanjnk/equipment-monitor

go mod init github.com/Melanjnk/equipment-monitor 


go get github.com/grpc-ecosystem/grpc-gateway/v2@v2.6.0


list of pgk:
github.com/jmoiron/sqlx

github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.6.0
	github.com/jackc/pgx/v4 v4.13.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/lib/pq v1.10.3
	github.com/opentracing/opentracing-go v1.2.0
	github.com/ozonmp/omp-template-api/pkg/omp-template-api v0.0.0-00010101000000-000000000000
	github.com/pressly/goose/v3 v3.1.0
	github.com/prometheus/client_golang v1.11.0
	github.com/rs/zerolog v1.24.0
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	google.golang.org/grpc v1.41.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b

```

pgcli -h 0.0.0.0 -p 54327 -u postgres
