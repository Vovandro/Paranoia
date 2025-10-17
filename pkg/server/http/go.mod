module gitlab.com/devpro_studio/Paranoia/pkg/server/http

go 1.24.0

require (
	github.com/golang-jwt/jwt/v5 v5.3.0
	gitlab.com/devpro_studio/go_utils v1.1.5
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.63.0
	go.opentelemetry.io/otel v1.38.0
	go.opentelemetry.io/otel/metric v1.38.0
)

replace gitlab.com/devpro_studio/Paranoia => ../../../

require (
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
)
