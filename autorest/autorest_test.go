package autorest

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/azure/go-autorest/autorest/mocks"
)

const (
	testAuthorizationHeader = "BEARER SECRETTOKEN"
	testUrl                 = "https://microsoft.com/a/b/c/"
)

func TestCreatePollingRequestIgnoresSuccess(t *testing.T) {
	resp := mocks.NewResponse()

	req, _, _ := CreatePollingRequest(resp, NullAuthorizer{})
	if req != nil {
		t.Error("autorest: CreatePollingRequest did not ignore a successful response")
	}
}

func TestCreatePollingRequestIgnoresSuccessWithoutError(t *testing.T) {
	resp := mocks.NewResponse()

	_, _, err := CreatePollingRequest(resp, NullAuthorizer{})
	if err != nil {
		t.Errorf("autorest: CreatePollingRequest returned an error when ignoring success (%v)", err)
	}
}

func TestCreatePollingRequestLeavesBodyOpenOfIgnoredSuccess(t *testing.T) {
	resp := mocks.NewResponse()

	req, _, _ := CreatePollingRequest(resp, NullAuthorizer{})
	if req != nil {
		t.Error("autorest: CreatePollingRequest closed the responise body while ignoring a successful response")
	}
}

func TestCreatePollingRequestDefaultsToAcceptedStatusCode(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _, _ := CreatePollingRequest(resp, NullAuthorizer{})
	if req == nil {
		t.Error("autorest: CreatePollingRequest failed to create a request for default 202 Accepted status code")
	}
}

func TestCreatePollingRequestReturnsErrorsForUnexpectedStatusCodes(t *testing.T) {
	resp := mocks.NewResponseWithStatus("500 ServerError", 500)
	addAcceptedHeaders(resp)

	_, _, err := CreatePollingRequest(resp, NullAuthorizer{})
	if err == nil {
		t.Errorf("autorest: CreatePollingRequest did not return an error when ignoring a status code (%v)", err)
	}
}

func TestCreatePollingRequestDoesNotCreateARequestForUnexpectedStatusCodes(t *testing.T) {
	resp := mocks.NewResponseWithStatus("500 ServerError", 500)
	addAcceptedHeaders(resp)

	req, _, _ := CreatePollingRequest(resp, NullAuthorizer{})
	if req != nil {
		t.Error("autorest: CreatePollingRequest created a new request when ignoring a status code")
	}
}

func TestCreatePollingRequestLeavesBodyOpenOfIgnoredStatusCode(t *testing.T) {
	resp := mocks.NewResponseWithStatus("500 ServerError", 500)
	addAcceptedHeaders(resp)

	CreatePollingRequest(resp, NullAuthorizer{})
	if !resp.Body.(*mocks.Body).IsOpen() {
		t.Error("autorest: CreatePollingRequest closed Body of an ignored status code")
	}
}

func TestCreatePollingRequestFailsWithoutLocationHeader(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	mocks.AddResponseHeader(resp, http.CanonicalHeaderKey(headerRetryAfter), "0")

	_, _, err := CreatePollingRequest(resp, NullAuthorizer{})
	if err == nil {
		t.Error("autorest: CreatePollingRequest failed to detect missing Location header")
	} else {
		f, _ := regexp.MatchString(".*Missing Location header.*", fmt.Sprintf("%v", err))
		if !f {
			t.Errorf("autorest: CreatePollingRequest returned an unexpected error (%v) for a missing Location", err)
		}
	}
}

func TestCreatePollingRequestFailsWithoutRetryAfterHeader(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	mocks.AddResponseHeader(resp, http.CanonicalHeaderKey(headerLocation), "https://microsoft.com/a/b/c/")

	_, _, err := CreatePollingRequest(resp, NullAuthorizer{})
	if err == nil {
		t.Error("autorest: CreatePollingRequest failed to detect missing Retry-After header")
	} else {
		f, _ := regexp.MatchString(".*Missing Retry-After header.*", fmt.Sprintf("%v", err))
		if !f {
			t.Errorf("autorest: CreatePollingRequest returned an unexpected error (%v) for a missing Retry-After", err)
		}
	}
}

func TestCreatePollingRequestClosesTheResponseBody(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	CreatePollingRequest(resp, NullAuthorizer{})
	if resp.Body.(*mocks.Body).IsOpen() {
		t.Error("autorest: CreatePollingRequest failed to close the response body when creating a new request")
	}
}

func TestCreatePollingRequestReturnsAGetRequest(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _, _ := CreatePollingRequest(resp, NullAuthorizer{})
	if req.Method != "GET" {
		t.Errorf("autorest: CreatePollingRequest did not create an HTTP GET request -- actual method %v", req.Method)
	}
}

func TestCreatePollingRequestProvidesTheURL(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _, _ := CreatePollingRequest(resp, NullAuthorizer{})
	if req.URL.String() != testUrl {
		t.Errorf("autorest: CreatePollingRequest did not create an HTTP with the expected URL -- received %v, expected %v", req.URL, testUrl)
	}
}

func TestCreatePollingRequestAppliesAuthorization(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _, _ := CreatePollingRequest(resp, TestAuthorizer{})
	if req.Header.Get(http.CanonicalHeaderKey(headerAuthorization)) != testAuthorizationHeader {
		t.Errorf("autorest: CreatePollingRequest did not apply authorization -- received %v, expected %v",
			req.Header.Get(http.CanonicalHeaderKey(headerAuthorization)), testAuthorizationHeader)
	}
}

func TestPollForAttemptsStops(t *testing.T) {
	client := mocks.NewClient()
	client.EmitErrors(-1)

	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _, _ := CreatePollingRequest(resp, NullAuthorizer{})

	PollForAttempts(client, req, time.Duration(0), 5)
	if client.Attempts() < 5 || client.Attempts() > 5 {
		t.Errorf("autorest: PollForAttempts stopped incorrectly -- expected %v attempts, actual attempts were %v", 5, client.Attempts())
	}
}

func TestPollForDurationsStops(t *testing.T) {
	client := mocks.NewClient()
	client.EmitErrors(-1)

	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _, _ := CreatePollingRequest(resp, NullAuthorizer{})

	d := 10 * time.Millisecond
	start := time.Now()
	PollForDuration(client, req, time.Duration(0), d)
	if time.Now().Sub(start) < d {
		t.Error("autorest: PollForDuration stopped too soon")
	}
}

func TestPollForDurationsStopsWithinReason(t *testing.T) {
	client := mocks.NewClient()
	client.EmitErrors(-1)

	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _, _ := CreatePollingRequest(resp, NullAuthorizer{})

	d := 10 * time.Millisecond
	start := time.Now()
	PollForDuration(client, req, time.Duration(0), d)
	if time.Now().Sub(start) > (time.Duration(5.0) * d) {
		t.Error("autorest: PollForDuration took too long to stop -- exceeded 5 times expected duration")
	}
}

func TestPollingHonorsDelay(t *testing.T) {
	client := mocks.NewClient()
	client.EmitErrors(-1)

	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _, _ := CreatePollingRequest(resp, NullAuthorizer{})

	d1 := 10 * time.Millisecond
	start := time.Now()
	PollForAttempts(client, req, d1, 1)
	d2 := time.Now().Sub(start)
	if d2 < d1 {
		t.Errorf("autorest: Polling failed to honor delay -- expected %v, actual %v", d1.Seconds(), d2.Seconds())
	}
}

func TestPollingReturnsErrorForExpectedStatusCode(t *testing.T) {
	client := mocks.NewClient()
	client.EmitStatus("202 Accepted", 202)

	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _, _ := CreatePollingRequest(resp, NullAuthorizer{})

	resp, err := PollForAttempts(client, req, time.Duration(0), 1, 202)
	if err == nil {
		t.Error("autorest: doPoll failed to emit error for known status code")
	}
}

func TestPollingReturnsNoErrorForUnexpectedStatusCode(t *testing.T) {
	client := mocks.NewClient()
	client.EmitStatus("500 ServerError", 500)

	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _, _ := CreatePollingRequest(resp, NullAuthorizer{})

	resp, err := PollForAttempts(client, req, time.Duration(0), 1, 202)
	if err != nil {
		t.Error("autorest: doPoll emitted error for unknown status code")
	}
}

func TestPollingReturnsDefaultsToAcceptedStatusCode(t *testing.T) {
	client := mocks.NewClient()
	client.EmitStatus("202 Accepted", 202)

	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _, _ := CreatePollingRequest(resp, NullAuthorizer{})

	resp, err := PollForAttempts(client, req, time.Duration(0), 1)
	if err == nil {
		t.Error("autorest: doPoll failed to default to HTTP 202")
	}
}

func TestdoPollLeavesFinalBodyOpen(t *testing.T) {
	client := mocks.NewClient()
	client.EmitStatus("500 ServerError", 500)

	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _, _ := CreatePollingRequest(resp, NullAuthorizer{})

	resp, _ = PollForAttempts(client, req, time.Duration(0), 1)
	if !resp.Body.(*mocks.Body).IsOpen() {
		t.Error("autorest: doPoll unexpectedly closed the response body")
	}
}

func TestResponseHasStatusCode(t *testing.T) {
	codes := []int{200, 202}
	resp := &http.Response{StatusCode: 202}
	if !ResponseHasStatusCode(resp, codes...) {
		t.Errorf("autorest: ResponseHasStatusCode failed to find %v in %v", resp.StatusCode, codes)
	}
}

func TestResponseHasStatusCodeNotPresent(t *testing.T) {
	codes := []int{200, 202}
	resp := &http.Response{StatusCode: 500}
	if ResponseHasStatusCode(resp, codes...) {
		t.Errorf("autorest: ResponseHasStatusCode unexpectedly found %v in %v", resp.StatusCode, codes)
	}
}

func addAcceptedHeaders(resp *http.Response) {
	mocks.AddResponseHeader(resp, http.CanonicalHeaderKey(headerLocation), testUrl)
	mocks.AddResponseHeader(resp, http.CanonicalHeaderKey(headerRetryAfter), "0")
}

type TestAuthorizer struct{}

func (ta TestAuthorizer) WithAuthorization() PrepareDecorator {
	return WithHeader(headerAuthorization, testAuthorizationHeader)
}
