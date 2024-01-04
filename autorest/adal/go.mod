module github.com/Azure/go-autorest/autorest/adal

go 1.15

require (
	github.com/Azure/go-autorest v14.2.0+incompatible
	github.com/Azure/go-autorest/autorest/date v0.3.0
	github.com/Azure/go-autorest/autorest/mocks v0.4.1
	github.com/Azure/go-autorest/logger v0.2.1
	github.com/Azure/go-autorest/tracing v0.6.0
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/stretchr/testify v1.8.2
	golang.org/x/crypto v0.17.0
)

retract [v0.9.5, v0.9.19] // retracted due to token refresh errors
