module github.com/Azure/go-autorest/autorest

go 1.12

replace github.com/Azure/go-autorest/autorest => ../autorest

replace github.com/Azure/go-autorest/logger => ../logger

replace github.com/Azure/go-autorest/tracing => ../tracing

require (
	github.com/Azure/go-autorest/autorest/adal v0.1.0
	github.com/Azure/go-autorest/autorest/mocks v0.1.0
	github.com/Azure/go-autorest/logger v0.0.0-00010101000000-000000000000
	github.com/Azure/go-autorest/tracing v0.1.0
	go.opencensus.io v0.20.2
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
)
