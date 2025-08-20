module gitlab.com/devpro_studio/Paranoia/pkg/external/NetLocker

go 1.23.4

require (
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.4
)

replace gitlab.com/devpro_studio/Paranoia/pkg/client/grpc-client => ../../client/grpc-client

require (
	go.opentelemetry.io/otel v1.34.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.34.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250127172529-29210b9bc287 // indirect
)
