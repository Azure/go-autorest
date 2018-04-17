package azure

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

import (
	"net/http"
	"testing"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/mocks"
)

const (
	TestTenantID                = "TestTenantID"
	TestActiveDirectoryEndpoint = "https://login/test.com/"
)

func TestApiKeyAuthorization(t *testing.T) {

	headers := make(map[string]interface{})
	queryParameters := make(map[string]interface{})

	dummyAuthHeader := "dummyAuthHeader"
	dummyAuthHeaderValue := "dummyAuthHeaderValue"

	dummyAuthQueryParameter := "dummyAuthQueryParameter"
	dummyAuthQueryParameterValue := "dummyAuthQueryParameterValue"

	headers[dummyAuthHeader] = dummyAuthHeaderValue
	queryParameters[dummyAuthQueryParameter] = dummyAuthQueryParameterValue

	aka := NewAPIKeyAuthorizer(headers, queryParameters)

	req, err := autorest.Prepare(mocks.NewRequest(), aka.WithAuthorization())

	if err != nil {
		t.Fatalf("azure: APIKeyAuthorizer#WithAuthorization returned an error (%v)", err)
	} else if req.Header.Get(http.CanonicalHeaderKey(dummyAuthHeader)) != dummyAuthHeaderValue {
		t.Fatalf("azure: APIKeyAuthorizer#WithAuthorization failed to set %s header", dummyAuthHeader)

	} else if req.URL.Query().Get(dummyAuthQueryParameter) != dummyAuthQueryParameterValue {
		t.Fatalf("azure: APIKeyAuthorizer#WithAuthorization failed to set %s query parameter", dummyAuthQueryParameterValue)
	}
}

func TestCognitivesServicesAuthorization(t *testing.T) {
	subscriptionKey := "dummyKey"
	csa := NewCognitiveServicesAuthorizer(subscriptionKey)
	req, err := autorest.Prepare(mocks.NewRequest(), csa.WithAuthorization())

	if err != nil {
		t.Fatalf("azure: CognitiveServicesAuthorizer#WithAuthorization returned an error (%v)", err)
	} else if req.Header.Get(http.CanonicalHeaderKey(bingAPISdkHeader)) != golangBingAPISdkHeaderValue {
		t.Fatalf("azure: CognitiveServicesAuthorizer#WithAuthorization failed to set %s header", bingAPISdkHeader)
	} else if req.Header.Get(http.CanonicalHeaderKey(apiKeyAuthorizerHeader)) != subscriptionKey {
		t.Fatalf("azure: CognitiveServicesAuthorizer#WithAuthorization failed to set %s header", apiKeyAuthorizerHeader)
	}
}
