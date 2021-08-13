module github.com/Azure/go-autorest/autorest/adal

go 1.15

require (
	github.com/Azure/go-autorest v14.2.0+incompatible
	github.com/Azure/go-autorest/autorest/date v0.3.0
	github.com/Azure/go-autorest/autorest/mocks v0.4.1
	github.com/Azure/go-autorest/logger v0.2.1
	github.com/Azure/go-autorest/tracing v0.6.0
	// NOTE: cannot move to github.com/golang-jwt/jwt/v4 yet, because the main
	//       github.com/Azure/go-autorest code uses go dep, which is not
	//       compatible with v2+ go modules that do not use a v<version> directory.
	github.com/golang-jwt/jwt v3.2.2+incompatible
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0
)
