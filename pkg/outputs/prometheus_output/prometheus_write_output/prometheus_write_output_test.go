package prometheus_write_output

import (
	"context"
	"sync"
	"testing"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmic/pkg/outputs"
)

// 10000 buffer, 1 writer
// BenchmarkRemoteWrite-10          1000000              1389 ns/op            1977 B/op         37 allocs/op

// 1000 buffer, 1 writer
// BenchmarkRemoteWrite-10           783300              1523 ns/op            1992 B/op         37 allocs/op

// 256 buffer, 1 writer
// BenchmarkRemoteWrite-10           842360              1355 ns/op            1993 B/op         37 allocs/op

// 1 buffer, 1 writer
// BenchmarkRemoteWrite-10           427292              2632 ns/op            1994 B/op         38 allocs/op

// 0 buffer, 1 writer
// BenchmarkRemoteWrite-10           494935              2125 ns/op            1994 B/op         38 allocs/op

// 10000 buffer, 60 writer
// BenchmarkRemoteWrite-10          1000000              2757 ns/op            1978 B/op         37 allocs/op

// 1000 buffer, 60 writer
// BenchmarkRemoteWrite-10           418759              2451 ns/op            1991 B/op         37 allocs/op

// 256 buffer, 60 writer
// BenchmarkRemoteWrite-10           387818              3391 ns/op            1993 B/op         37 allocs/op

// 1 buffer, 60 writer
// BenchmarkRemoteWrite-10           380790              3164 ns/op            1993 B/op         38 allocs/op

// 0 buffer, 60 writer
// BenchmarkRemoteWrite-10           281700              3637 ns/op            1996 B/op         38 allocs/op

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
	wg := new(sync.WaitGroup)
	writers := 60
	for j := 0; j < writers; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < b.N/writers+1; i++ {
				output.Write(context.Background(), subResponse, outputs.Meta{})
			}
			wg.Done()
		}()
	}
	wg.Wait()

	b.StopTimer()

	output.Close()
}
