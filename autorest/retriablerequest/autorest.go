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

package retriablerequest

import (
	"math"
	"net/http"
	"strconv"
	"time"
)

var (
	// StatusCodesForRetry are a defined group of status code for which the client will retry
	StatusCodesForRetry = []int{
		http.StatusRequestTimeout,      // 408
		http.StatusTooManyRequests,     // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout,      // 504
	}

	// DefaultRetryAttempts is number of attempts for retry status codes (5xx).
	DefaultRetryAttempts = 5

	// DefaultRetryDuration is the duration to wait between retries.
	DefaultRetryDuration = 1 * time.Second
)

// DelayWithRetryAfter invokes time.After for the duration specified in the "Retry-After" header in
// responses with status code 429
func DelayWithRetryAfter(resp *http.Response, cancel <-chan struct{}) bool {
	if resp == nil {
		return false
	}
	retryAfter, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
	if resp.StatusCode == http.StatusTooManyRequests && retryAfter > 0 {
		select {
		case <-time.After(time.Duration(retryAfter) * time.Second):
			return true
		case <-cancel:
			return false
		}
	}
	return false
}

// DelayForBackoff invokes time.After for the supplied backoff duration raised to the power of
// passed attempt (i.e., an exponential backoff delay). Backoff duration is in seconds and can set
// to zero for no delay. The delay may be canceled by closing the passed channel. If terminated early,
// returns false.
// Note: Passing attempt 1 will result in doubling "backoff" duration. Treat this as a zero-based attempt
// count.
func DelayForBackoff(backoff time.Duration, attempt int, cancel <-chan struct{}) bool {
	select {
	case <-time.After(time.Duration(backoff.Seconds()*math.Pow(2, float64(attempt))) * time.Second):
		return true
	case <-cancel:
		return false
	}
}

// ResponseHasStatusCode returns true if the status code in the HTTP Response is in the passed set
// and false otherwise.
func ResponseHasStatusCode(resp *http.Response, codes ...int) bool {
	if resp == nil {
		return false
	}
	return containsInt(codes, resp.StatusCode)
}

func containsInt(ints []int, n int) bool {
	for _, i := range ints {
		if i == n {
			return true
		}
	}
	return false
}

func Retry(req *http.Request, backoff time.Duration, attempts int) (resp *http.Response, err error) {
	s := http.Client{}
	rr := NewRetriableRequest(req)
	// Increment to add the first call (attempts denotes number of retries)
	for attempt := 0; attempt < attempts; {
		err = rr.Prepare()
		if err != nil {
			return
		}
		resp, err = s.Do(rr.Request())
		// we want to retry if err is not nil (e.g. transient network failure).
		if err == nil && !ResponseHasStatusCode(resp, StatusCodesForRetry...) {
			return
		}
		delayed := DelayWithRetryAfter(resp, req.Cancel)
		if !delayed {
			DelayForBackoff(backoff, attempt, req.Cancel)
		}
		// don't count a 429 against the number of attempts
		// so that we continue to retry until it succeeds
		if resp == nil || resp.StatusCode != http.StatusTooManyRequests {
			attempt++
		}
	}
	return
}
