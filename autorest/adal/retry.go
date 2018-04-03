// +build go1.8

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

package adal

import (
	"bytes"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"
)

var (
	// statusCodesForRetry are a defined group of status code for which the client will retry
	statusCodesForRetry = []int{
		http.StatusRequestTimeout,      // 408
		http.StatusTooManyRequests,     // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout,      // 504
	}
)

// retriableRequest provides facilities for retrying an HTTP request.
type retriableRequest struct {
	req *http.Request
	rc  io.ReadCloser
	br  *bytes.Reader
}

// prepare signals that the request is about to be sent.
func (rr *retriableRequest) prepare() (err error) {
	// preserve the request body; this is to support retry logic as
	// the underlying transport will always close the reqeust body
	if rr.req.Body != nil {
		if rr.rc != nil {
			rr.req.Body = rr.rc
		} else if rr.br != nil {
			_, err = rr.br.Seek(0, io.SeekStart)
			rr.req.Body = ioutil.NopCloser(rr.br)
		}
		if err != nil {
			return err
		}
		if rr.req.GetBody != nil {
			// this will allow us to preserve the body without having to
			// make a copy.  note we need to do this on each iteration
			rr.rc, err = rr.req.GetBody()
			if err != nil {
				return err
			}
		} else if rr.br == nil {
			// fall back to making a copy (only do this once)
			err = rr.prepareFromByteReader()
		}
	}
	return err
}

// newRetriableRequest returns a wrapper around an HTTP request that support retry logic.
func newRetriableRequest(req *http.Request) *retriableRequest {
	return &retriableRequest{req: req}
}

// request returns the wrapped HTTP request.
func (rr *retriableRequest) request() *http.Request {
	return rr.req
}

func (rr *retriableRequest) prepareFromByteReader() (err error) {
	// fall back to making a copy (only do this once)
	b := []byte{}
	if rr.req.ContentLength > 0 {
		b = make([]byte, rr.req.ContentLength)
		_, err = io.ReadFull(rr.req.Body, b)
		if err != nil {
			return err
		}
	} else {
		b, err = ioutil.ReadAll(rr.req.Body)
		if err != nil {
			return err
		}
	}
	rr.br = bytes.NewReader(b)
	rr.req.Body = ioutil.NopCloser(rr.br)
	return err
}

func retry(req *http.Request, backoff time.Duration, attempts int) (resp *http.Response, err error) {
	s := http.Client{}
	rr := newRetriableRequest(req)
	// Increment to add the first call (attempts denotes number of retries)
	for attempt := 0; attempt < attempts; {
		err = rr.prepare()
		if err != nil {
			return
		}
		resp, err = s.Do(rr.request())
		// we want to retry if err is not nil (e.g. transient network failure).
		if err == nil && !responseHasStatusCode(resp, statusCodesForRetry...) {
			return
		}
		delayed := delayWithRetryAfter(resp)
		if !delayed {
			delayForBackoff(backoff, attempt)
		}
		// don't count a 429 against the number of attempts
		// so that we continue to retry until it succeeds
		if resp == nil || resp.StatusCode != http.StatusTooManyRequests {
			attempt++
		}
	}
	return
}

// delayWithRetryAfter invokes time.After for the duration specified in the "Retry-After" header in
// responses with status code 429
func delayWithRetryAfter(resp *http.Response) bool {
	if resp == nil {
		return false
	}
	retryAfter, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
	if resp.StatusCode == http.StatusTooManyRequests && retryAfter > 0 {
		time.Sleep(time.Duration(retryAfter) * time.Second)
		return true
	}
	return false
}

// delayForBackoff invokes time.After for the supplied backoff duration raised to the power of
// passed attempt (i.e., an exponential backoff delay). Backoff duration is in seconds and can set
// to zero for no delay. The delay may be canceled by closing the passed channel. If terminated early,
// returns false.
// Note: Passing attempt 1 will result in doubling "backoff" duration. Treat this as a zero-based attempt
// count.
func delayForBackoff(backoff time.Duration, attempt int) {
	time.Sleep(time.Duration(backoff.Seconds()*math.Pow(2, float64(attempt))) * time.Second)
}

// responseHasStatusCode returns true if the status code in the HTTP Response is in the passed set
// and false otherwise.
func responseHasStatusCode(resp *http.Response, codes ...int) bool {
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
