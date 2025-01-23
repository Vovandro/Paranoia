package telemetry

import "gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"

func NewMetrics(cfg map[string]interface{}) interfaces.IMetrics {
	t, ok := cfg["name"].(string)

	if ok {
		switch t {
		case "prometheus":
			return NewMetricPrometheus(t)

		case "std":
			return NewMetricStd(t)

		case "oltp_grpc":
			return NewMetricOtlpGrpc(t)

		case "otlp_http":
			return NewMetricOtlpHttp(t)
		}
	}

	return nil
}

func NewTrace(cfg map[string]interface{}) interfaces.ITrace {
	t, ok := cfg["name"].(string)
	if ok {
		switch t {
		case "std":
			return NewTraceStd(t)

		case "otlp_grpc":
			return NewTraceOtlpGrpc(t)

		case "otlp_http":
			return NewTraceOtlpHttp(t)

		case "sentry":
			return NewTraceSentry(t)

		case "zipkin":
			return NewTraceZipking(t)
		}
	}

	return nil
}
