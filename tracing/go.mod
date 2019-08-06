module github.com/Azure/go-autorest/tracing

go 1.12

require (
	// use older versions to avoid taking a dependency on protobuf v1.3+
	contrib.go.opencensus.io/exporter/ocagent v0.4.6
	go.opencensus.io v0.18.1-0.20181204023538-aab39bd6a98b
)
