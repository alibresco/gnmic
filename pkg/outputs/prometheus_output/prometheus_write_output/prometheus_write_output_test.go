package prometheus_write_output

import (
	"context"
	"testing"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmic/pkg/outputs"
)

func BenchmarkRemoteWrite(b *testing.B) {
	cfg := map[string]interface{}{
		"url":                       "http://example.com",
		"timeout":                   "100ms",
		"buffer-size":               1000000,
		"max-time-series-per-write": defaultMaxTSPerWrite,
		"num-workers":               4,
		"num-writers":               1,
		"metadata":                  metadata{Include: false},
		"interval":                  "10s",
		// "debug": true,
	}

	output := outputs.Outputs[outputType]()
	output.Init(context.Background(), "testOutput", cfg)

	subResponse := &gnmi.SubscribeResponse{
		Response: &gnmi.SubscribeResponse_Update{
			Update: &gnmi.Notification{Timestamp: 123, Prefix: &gnmi.Path{Origin: "origin"},
				Update: []*gnmi.Update{
					{
						Path: &gnmi.Path{
							Element: []string{"element1", "element2"},
						},
						// Optionally add a value
						Val: &gnmi.TypedValue{
							Value: &gnmi.TypedValue_StringVal{StringVal: "newValue"},
						},
					},
					{
						Path: &gnmi.Path{
							Element: []string{"anotherElement1", "anotherElement2"},
						},
						// Optionally add another value
						Val: &gnmi.TypedValue{
							Value: &gnmi.TypedValue_IntVal{IntVal: 42},
						},
					},
					// Add more gnmi.Update elements here as needed
				},
			},
		},
	}

	// msg, _ := anypb.New(subResponse)

	// Running the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output.Write(context.Background(), subResponse, outputs.Meta{})
	}
	b.StopTimer()

	output.Close()
}
