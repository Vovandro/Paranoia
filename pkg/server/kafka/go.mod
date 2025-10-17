module gitlab.com/devpro_studio/Paranoia/pkg/server/kafka

go 1.24.0

require (
	github.com/confluentinc/confluent-kafka-go/v2 v2.12.0
	github.com/jurabek/otelkafka v1.0.1
	gitlab.com/devpro_studio/Paranoia v0.0.0-00010101000000-000000000000
	gitlab.com/devpro_studio/go_utils v1.1.5
	go.opentelemetry.io/otel v1.38.0
	go.opentelemetry.io/otel/metric v1.38.0
)

replace gitlab.com/devpro_studio/Paranoia => ../../../

require (
	github.com/docker/docker v27.4.1+incompatible // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.59.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.59.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/oauth2 v0.31.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
)
