module gitlab.com/devpro_studio/Paranoia/pkg/server/rabbitmq

go 1.24.0

require (
	github.com/rabbitmq/amqp091-go v1.10.0
	gitlab.com/devpro_studio/Paranoia v0.0.0-00010101000000-000000000000
	gitlab.com/devpro_studio/go_utils v1.1.5
	go.opentelemetry.io/otel v1.38.0
	go.opentelemetry.io/otel/metric v1.38.0
)

replace gitlab.com/devpro_studio/Paranoia => ../../../

require (
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
)
