package autorest

import "net/http"

// Copyright 2017 Microsoft Corporation
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

// AuthenticationError is an interface used by errors returned by authentication responses.
// For example, token referesh errors
type AuthenticationError interface {
	error
	Response() *http.Response
}

// IsAuthenticationError returns true if the specified error implements the TokenRefreshError
// interface.  If err is a DetailedError it will walk the chain of Original errors.
func IsAuthenticationError(err error) bool {
	if _, ok := err.(AuthenticationError); ok {
		return true
	}
	if de, ok := err.(DetailedError); ok {
		return IsAuthenticationError(de.Original)
	}
	return false
}

// WithAuthorization returns a PrepareDecorator that adds the specified value to the Authorization header
func WithAuthorization(authorization string) PrepareDecorator {
	return WithHeader(HeaderAuthorization, authorization)
}
