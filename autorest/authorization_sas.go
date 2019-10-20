package autorest

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
	"fmt"
	"net/http"
	"strings"
)

// SASTokenAuthorizer implements an authorization for SAS Token Authentication
// this can be used for interaction with Blob Storage Endpoints
type SASTokenAuthorizer struct {
	sasToken string
}

// NewSASTokenAuthorizer creates a SASTokenAuthorizer using the given credentials
func NewSASTokenAuthorizer(sasToken string) (*SASTokenAuthorizer, error) {
	if strings.TrimSpace(sasToken) == "" {
		return nil, fmt.Errorf("`sasToken` cannot be empty!")
	}

	token := sasToken
	if strings.HasPrefix(sasToken, "?") {
		token = strings.TrimPrefix(sasToken, "?")
	}

	return &SASTokenAuthorizer{
		sasToken: token,
	}, nil
}

// WithAuthorization returns a PrepareDecorator that adds an HTTP Authorization header whose
// value is "SharedKey " followed by the computed key.
// This can be used for the Blob, Queue, and File Services
//
// from: https://docs.microsoft.com/en-us/rest/api/storageservices/authorize-with-shared-key
// You may use Shared Key Lite authorization to authorize a request made against the
// 2009-09-19 version and later of the Blob and Queue services,
// and version 2014-02-14 and later of the File services.
func (sas *SASTokenAuthorizer) WithAuthorization() PrepareDecorator {
	return func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)
			if err != nil {
				return r, err
			}

			queryString := r.URL.RawQuery
			if queryString != "" {
				queryString = fmt.Sprintf("%s&%s", queryString, sas.sasToken)
			} else {
				queryString = sas.sasToken
			}

			r.URL.RawQuery = queryString
			r.RequestURI = r.URL.String()
			return Prepare(r)
		})
	}
}
