module github.com/Azure/go-autorest/tracing/opencensus

go 1.12

require (
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	github.com/Azure/go-autorest/tracing v0.7.0
	go.opencensus.io v0.22.0
)

replace github.com/Azure/go-autorest/tracing => ../
