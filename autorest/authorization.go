package autorest

import (
	"fmt"
	"net/http"
)

// Authorizer is the interface that provides a PrepareDecorator used to supply request
// authorization. Most often, the Authorizer decorator runs last so it has access to the full
// state of the formed HTTP request.
type Authorizer interface {
	WithAuthorization() PrepareDecorator
}

// NullAuthorizer implements a default, "do nothing" Authorizer.
type NullAuthorizer struct{}

// WithAuthorization returns a PrepareDecorator that does nothing.
func (na NullAuthorizer) WithAuthorization() PrepareDecorator {
	return WithNothing()
}

// BearerAuthorizer implements the bearer authorization
type BearerAuthorizer struct {
	token string
}

func withBearerAuthorization(token string) PrepareDecorator {
	return WithHeader(headerAuthorization, fmt.Sprintf("Bearer %s", token))
}

// WithAuthorization returns a PrepareDecorator that adds an HTTP Authorization header whose
// value is "Bearer " followed by the AccessToken of the ServicePrincipalToken.
//
// By default, the token will automatically refresh if nearly expired (as determined by the
// RefreshWithin interval). Use the AutoRefresh method to enable or disable automatically refreshing
// tokens.
func (ba *BearerAuthorizer) WithAuthorization() PrepareDecorator {
	return func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			return (withBearerAuthorization(ba.token)(p)).Prepare(r)
		})
	}
}
