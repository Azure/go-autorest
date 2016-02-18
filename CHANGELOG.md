# CHANGELOG

## v5.0.0
- Added new RespondDecorators unmarshalling primitive types
- Corrected application of inspection and authorization PrependDecorators

## v4.0.0
- Added support for Azure long-running operations.
- Added cancelation support to all decorators and functions that may delay.
- Breaking: `DelayForBackoff` now accepts a channel, which may be nil.

## v3.1.0
- Add support for OAuth Device Flow authorization.
- Add support for ServicePrincipalTokens that are backed by an existing token, rather than other secret material.
- Add helpers for persisting and restoring Tokens.
- Increased code coverage in the github.com/Azure/autorest/azure package

## v3.0.0
- Breaking: `NewErrorWithError` no longer takes `statusCode int`.
- Breaking: `NewErrorWithStatusCode` is replaced with `NewErrorWithResponse`.
- Breaking: `Client#Send()` no longer takes `codes ...int` argument.
- Add: XML unmarshaling support with `ByUnmarshallingXML()`
- Stopped vending dependencies locally and switched to [Glide](https://github.com/Masterminds/glide).
  Applications using this library should either use Glide or vendor dependencies locally some other way.
- Add: `azure.WithErrorUnlessStatusCode()` decorator to handle Azure errors.
- Fix: use `net/http.DefaultClient` as base client.
- Fix: Missing inspection for polling responses added.
- Add: CopyAndDecode helpers.
- Improved `./autorest/to` with `[]string` helpers.
- Removed golint suppressions in .travis.yml.

## v2.1.0

- Added `StatusCode` to `Error` for more easily obtaining the HTTP Reponse StatusCode (if any)

## v2.0.0

- Changed `to.StringMapPtr` method signature to return a pointer
- Changed `ServicePrincipalCertificateSecret` and `NewServicePrincipalTokenFromCertificate` to support generic certificate and private keys

## v1.0.0

- Added Logging inspectors to trace http.Request / Response
- Added support for User-Agent header
- Changed WithHeader PrepareDecorator to use set vs. add
- Added JSON to error when unmarshalling fails
- Added Client#Send method
- Corrected case of "Azure" in package paths
- Added "to" helpers, Azure helpers, and improved ease-of-use
- Corrected golint issues

## v1.0.1

- Added CHANGELOG.md

## v1.1.0

- Added mechanism to retrieve a ServicePrincipalToken using a certificate-signed JWT
- Added an example of creating a certificate-based ServicePrincipal and retrieving an OAuth token using the certificate

## v1.1.1

- Introduce godeps and vendor dependencies introduced in v1.1.1
