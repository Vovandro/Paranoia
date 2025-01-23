package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
)

func TestNewMetrics(t *testing.T) {
	tests := []struct {
		name     string
		cfg      map[string]interface{}
		expected interfaces.IMetrics
	}{
		{
			name:     "Return Prometheus Metric",
			cfg:      map[string]interface{}{"name": "prometheus"},
			expected: &MetricPrometheus{name: "prometheus"},
		},
		{
			name:     "Return Std Metric",
			cfg:      map[string]interface{}{"name": "std"},
			expected: &MetricStd{name: "std"},
		},
		{
			name:     "Return OTLP GRPC Metric",
			cfg:      map[string]interface{}{"name": "oltp_grpc"},
			expected: &MetricOtlpGrpc{name: "oltp_grpc"},
		},
		{
			name:     "Return OTLP HTTP Metric",
			cfg:      map[string]interface{}{"name": "otlp_http"},
			expected: &MetricOtlpHttp{name: "otlp_http"},
		},
		{
			name:     "Return nil for unknown metric type",
			cfg:      map[string]interface{}{"name": "unknown"},
			expected: nil,
		},
		{
			name:     "Return nil for missing name",
			cfg:      map[string]interface{}{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewMetrics(tt.cfg)

			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected.Name(), result.Name())
			}
		})
	}
}

func TestNewTrace(t *testing.T) {
	tests := []struct {
		name string
		cfg  map[string]interface{}
		want interfaces.ITrace
	}{
		{
			name: "std trace",
			cfg:  map[string]interface{}{"name": "std"},
			want: NewTraceStd("std"),
		},
		{
			name: "otlp_grpc trace",
			cfg:  map[string]interface{}{"name": "otlp_grpc"},
			want: NewTraceOtlpGrpc("otlp_grpc"),
		},
		{
			name: "otlp_http trace",
			cfg:  map[string]interface{}{"name": "otlp_http"},
			want: NewTraceOtlpHttp("otlp_http"),
		},
		{
			name: "sentry trace",
			cfg:  map[string]interface{}{"name": "sentry"},
			want: NewTraceSentry("sentry"),
		},
		{
			name: "zipkin trace",
			cfg:  map[string]interface{}{"name": "zipkin"},
			want: NewTraceZipking("zipkin"),
		},
		{
			name: "invalid name",
			cfg:  map[string]interface{}{"name": "unknown"},
			want: nil,
		},
		{
			name: "missing name",
			cfg:  map[string]interface{}{},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewTrace(tt.cfg))
		})
	}
}
