module gitlab.com/devpro_studio/Paranoia/pkg/external/NetLocker

go 1.23.4

require (
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.4
)

replace gitlab.com/devpro_studio/Paranoia/pkg/client/grpc-client => ../../client/grpc-client

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	gitlab.com/devpro_studio/go_utils v1.1.2 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.59.0 // indirect
	go.opentelemetry.io/otel v1.34.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.34.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
)
