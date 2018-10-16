package tracing

// Copyright 2018 Microsoft Corporation
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

import (
	"context"
	"net/http"
	"os"

	"contrib.go.opencensus.io/exporter/ocagent"
	"go.opencensus.io/stats/view"

	"go.opencensus.io/trace"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
)

var (
	// Transport is the default tracing RoundTripper. If tracing is disabled
	// the default RoundTripper is used which provides no instrumentation.
	Transport = http.DefaultTransport

	// Enabled is the flag for marking if tracing is enabled.
	Enabled = false

	// Sampler is the tracing sampler. If tracing is disabled it will never sample. Otherwise
	// it will be using the parent sampler or the default.
	sampler = trace.NeverSample()
)

func init() {
	enableFromEnv()
}

func enableFromEnv() {
	_, ok := os.LookupEnv("AZURE_SDK_TRACING_ENABELD")
	if ok {
		agentEndpoint, ok := os.LookupEnv("OCAGENT_TRACE_EXPORTER_ENDPOINT")

		if ok {
			EnableWithAIForwarding(agentEndpoint)
		} else {
			Enable()
		}
	}
}

// Enable will start instrumentation for metrics and traces.
func Enable() (err error) {
	Enabled = true
	sampler = nil
	Transport = &ochttp.Transport{Propagation: &tracecontext.HTTPFormat{}}

	if err != nil {
		return
	}
	err = initStats()
	return
}

// EnableWithAIForwarding will start instrumentation and will connect to app insights forwarder
// exporter making the metrics and traces available in app insights.
func EnableWithAIForwarding(agentEndpoint string) (err error) {
	err = Enable()
	if err != nil {
		return err
	}

	exporter, err := ocagent.NewExporter(ocagent.WithInsecure(), ocagent.WithAddress(agentEndpoint))
	if err != nil {
		return err
	}
	trace.RegisterExporter(exporter)
	return
}

// initStats registers the views for the http metrics
func initStats() (err error) {
	clientViews := []*view.View{
		ochttp.ClientCompletedCount,
		ochttp.ClientRoundtripLatencyDistribution,
		ochttp.ClientReceivedBytesDistribution,
		ochttp.ClientSentBytesDistribution,
	}
	if err = view.Register(clientViews...); err != nil {
		return err
	}
	return
}

// StartSpan starts a trace span
func StartSpan(ctx context.Context, name string) context.Context {
	ctx, _ = trace.StartSpan(ctx, name, trace.WithSampler(sampler))
	return ctx
}

// EndSpan ends a previously started span stored in the context
func EndSpan(ctx context.Context, statusCode int, err error) {
	span := trace.FromContext(ctx)

	if span == nil {
		return
	}

	if err != nil {
		span.SetStatus(trace.Status{Message: err.Error(), Code: int32(statusCode)})
	}
	span.End()
}
