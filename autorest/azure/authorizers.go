package azure

import "github.com/Azure/go-autorest/autorest"

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

const (
	apiKeyAuthorizerHeader      = "Ocp-Apim-Subscription-Key"
	bingAPISdkHeader            = "X-BingApis-SDK-Client"
	golangBingAPISdkHeaderValue = "Go-SDK"
)

// APIKeyAuthorizer implements API Key authorization.
type APIKeyAuthorizer struct {
	headers         map[string]interface{}
	queryParameters map[string]interface{}
}

// NewAPIKeyAuthorizerWithHeaders creates an ApiKeyAuthorizer with headers.
func NewAPIKeyAuthorizerWithHeaders(headers map[string]interface{}) *APIKeyAuthorizer {
	return NewAPIKeyAuthorizer(headers, nil)
}

// NewAPIKeyAuthorizerWithQueryParameters creates an ApiKeyAuthorizer with query parameters.
func NewAPIKeyAuthorizerWithQueryParameters(queryParameters map[string]interface{}) *APIKeyAuthorizer {
	return NewAPIKeyAuthorizer(nil, queryParameters)
}

// NewAPIKeyAuthorizer creates an ApiKeyAuthorizer with headers.
func NewAPIKeyAuthorizer(headers map[string]interface{}, queryParameters map[string]interface{}) *APIKeyAuthorizer {
	return &APIKeyAuthorizer{headers: headers, queryParameters: queryParameters}
}

// WithAuthorization returns a PrepareDecorator that adds an HTTP headers and Query Paramaters
func (aka *APIKeyAuthorizer) WithAuthorization() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.DecoratePreparer(p, autorest.WithHeaders(aka.headers), autorest.WithQueryParameters(aka.queryParameters))
	}
}

// CognitiveServicesAuthorizer implements authorization for Cognitive Services.
type CognitiveServicesAuthorizer struct {
	subscriptionKey string
}

// NewCognitiveServicesAuthorizer is
func NewCognitiveServicesAuthorizer(subscriptionKey string) *CognitiveServicesAuthorizer {
	return &CognitiveServicesAuthorizer{subscriptionKey: subscriptionKey}
}

// WithAuthorization is
func (csa *CognitiveServicesAuthorizer) WithAuthorization() autorest.PrepareDecorator {
	headers := make(map[string]interface{})
	headers[apiKeyAuthorizerHeader] = csa.subscriptionKey
	headers[bingAPISdkHeader] = golangBingAPISdkHeaderValue

	return NewAPIKeyAuthorizerWithHeaders(headers).WithAuthorization()
}

// EventGridKeyAuthorizer implements authorization for event grid using key authentication.
type EventGridKeyAuthorizer struct {
	topicKey string
}

// NewEventGridKeyAuthorizer creates a new EventGridKeyAuthorizer
// with the specified topic key.
func NewEventGridKeyAuthorizer(topicKey string) EventGridKeyAuthorizer {
	return EventGridKeyAuthorizer{topicKey: topicKey}
}

// WithAuthorization returns a PrepareDecorator that adds the aeg-sas-key authentication header.
func (egta EventGridKeyAuthorizer) WithAuthorization() autorest.PrepareDecorator {
	headers := map[string]interface{}{
		"aeg-sas-key": egta.topicKey,
	}
	return NewAPIKeyAuthorizerWithHeaders(headers).WithAuthorization()
}
