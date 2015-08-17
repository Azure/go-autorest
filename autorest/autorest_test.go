package autorest

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/azure/go-autorest/autorest/mocks"
)

const (
	testAuthorizationHeader = "BEARER SECRETTOKEN"
	testDelay               = 0 * time.Second
	testBadUrl              = ""
	testUrl                 = "https://microsoft.com/a/b/c/"
)

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

func TestResponseRequiresPollingIgnoresSuccess(t *testing.T) {
	resp := mocks.NewResponse()

	if ResponseRequiresPolling(resp) {
		t.Error("autorest: ResponseRequiresPolling did not ignore a successful response")
	}
}

func TestResponseRequiresPollingLeavesBodyOpen(t *testing.T) {
	resp := mocks.NewResponse()

	ResponseRequiresPolling(resp)
	if !resp.Body.(*mocks.Body).IsOpen() {
		t.Error("autorest: ResponseRequiresPolling closed the responise body while ignoring a successful response")
	}
}

func TestResponseRequiresPollingDefaultsToAcceptedStatusCode(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	if !ResponseRequiresPolling(resp) {
		t.Error("autorest: ResponseRequiresPolling failed to create a request for default 202 Accepted status code")
	}
}

func TestResponseRequiresPollingReturnsFalseForUnexpectedStatusCodes(t *testing.T) {
	resp := mocks.NewResponseWithStatus("500 ServerError", 500)
	addAcceptedHeaders(resp)

	if ResponseRequiresPolling(resp) {
		t.Error("autorest: ResponseRequiresPolling did not return false when ignoring a status code")
	}
}

func TestCreatePollingRequestLeavesBodyOpenWhenLocationHeaderIsMissing(t *testing.T) {
	resp := mocks.NewResponseWithStatus("500 ServerError", 500)

	CreatePollingRequest(resp, NullAuthorizer{})
	if !resp.Body.(*mocks.Body).IsOpen() {
		t.Error("autorest: CreatePollingRequest closed the http.Request Body when the Location header was missing")
	}
}

func TestCreatePollingRequestDoesNotReturnARequestWhenLocationHeaderIsMissing(t *testing.T) {
	resp := mocks.NewResponseWithStatus("500 ServerError", 500)

	req, _ := CreatePollingRequest(resp, NullAuthorizer{})
	if req != nil {
		t.Error("autorest: CreatePollingRequest returned an http.Request when the Location header was missing")
	}
}

func TestCreatePollingRequestReturnsAnErrorWhenPrepareFails(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)
	resp.Header.Set(http.CanonicalHeaderKey(headerLocation), testBadUrl)

	_, err := CreatePollingRequest(resp, NullAuthorizer{})
	if err == nil {
		t.Error("autorest: CreatePollingRequest failed to return an error when Prepare fails")
	}
}

func TestCreatePollingRequestLeavesBodyOpenWhenPrepareFails(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)
	resp.Header.Set(http.CanonicalHeaderKey(headerLocation), testBadUrl)

	_, err := CreatePollingRequest(resp, NullAuthorizer{})
	if !resp.Body.(*mocks.Body).IsOpen() {
		t.Errorf("autorest: CreatePollingRequest closed the http.Request Body when Prepare returned an error (%v)", err)
	}
}

func TestCreatePollingRequestDoesNotReturnARequestWhenPrepareFails(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)
	resp.Header.Set(http.CanonicalHeaderKey(headerLocation), testBadUrl)

	req, _ := CreatePollingRequest(resp, NullAuthorizer{})
	if req != nil {
		t.Error("autorest: CreatePollingRequest returned an http.Request when Prepare failed")
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

	req, _ := CreatePollingRequest(resp, NullAuthorizer{})
	if req.Method != "GET" {
		t.Errorf("autorest: CreatePollingRequest did not create an HTTP GET request -- actual method %v", req.Method)
	}
}

func TestCreatePollingRequestProvidesTheURL(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _ := CreatePollingRequest(resp, NullAuthorizer{})
	if req.URL.String() != testUrl {
		t.Errorf("autorest: CreatePollingRequest did not create an HTTP with the expected URL -- received %v, expected %v", req.URL, testUrl)
	}
}

func TestCreatePollingRequestAppliesAuthorization(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _ := CreatePollingRequest(resp, TestAuthorizer{})
	if req.Header.Get(http.CanonicalHeaderKey(headerAuthorization)) != testAuthorizationHeader {
		t.Errorf("autorest: CreatePollingRequest did not apply authorization -- received %v, expected %v",
			req.Header.Get(http.CanonicalHeaderKey(headerAuthorization)), testAuthorizationHeader)
	}
}

func TestGetRetryDelay(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	d := GetRetryDelay(resp, DefaultPollingDelay)
	if d != testDelay {
		t.Errorf("autorest: GetRetryDelay failed to returned the expected delay -- expected %v, received %v", testDelay, d)
	}
}

func TestGetRetryDelayReturnsDefaultDelayIfRetryHeaderIsMissing(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", 202)

	d := GetRetryDelay(resp, DefaultPollingDelay)
	if d != DefaultPollingDelay {
		t.Errorf("autorest: GetRetryDelay failed to returned the default delay for a missing Retry-After header -- expected %v, received %v",
			DefaultPollingDelay, d)
	}
}

func TestGetRetryDelayReturnsDefaultDelayIfRetryHeaderIsMalformed(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)
	resp.Header.Set(http.CanonicalHeaderKey(headerRetryAfter), "a very bad non-integer value")

	d := GetRetryDelay(resp, DefaultPollingDelay)
	if d != DefaultPollingDelay {
		t.Errorf("autorest: GetRetryDelay failed to returned the default delay for a malformed Retry-After header -- expected %v, received %v",
			DefaultPollingDelay, d)
	}
}

func TestPollForAttemptsStops(t *testing.T) {
	client := mocks.NewSender()
	client.EmitErrors(-1)

	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _ := CreatePollingRequest(resp, NullAuthorizer{})

	PollForAttempts(client, req, time.Duration(0), 5)
	if client.Attempts() < 5 || client.Attempts() > 5 {
		t.Errorf("autorest: PollForAttempts stopped incorrectly -- expected %v attempts, actual attempts were %v", 5, client.Attempts())
	}
}

func TestPollForDurationsStops(t *testing.T) {
	client := mocks.NewSender()
	client.EmitErrors(-1)

	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _ := CreatePollingRequest(resp, NullAuthorizer{})

	d := 10 * time.Millisecond
	start := time.Now()
	PollForDuration(client, req, time.Duration(0), d)
	if time.Now().Sub(start) < d {
		t.Error("autorest: PollForDuration stopped too soon")
	}
}

func TestPollForDurationsStopsWithinReason(t *testing.T) {
	client := mocks.NewSender()
	client.EmitErrors(-1)

	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _ := CreatePollingRequest(resp, NullAuthorizer{})

	d := 10 * time.Millisecond
	start := time.Now()
	PollForDuration(client, req, time.Duration(0), d)
	if time.Now().Sub(start) > (time.Duration(5.0) * d) {
		t.Error("autorest: PollForDuration took too long to stop -- exceeded 5 times expected duration")
	}
}

func TestPollingHonorsDelay(t *testing.T) {
	client := mocks.NewSender()
	client.EmitErrors(-1)

	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _ := CreatePollingRequest(resp, NullAuthorizer{})

	d1 := 10 * time.Millisecond
	start := time.Now()
	PollForAttempts(client, req, d1, 1)
	d2 := time.Now().Sub(start)
	if d2 < d1 {
		t.Errorf("autorest: Polling failed to honor delay -- expected %v, actual %v", d1.Seconds(), d2.Seconds())
	}
}

func TestPollingReturnsErrorForExpectedStatusCode(t *testing.T) {
	client := mocks.NewSender()
	client.EmitStatus("202 Accepted", 202)

	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _ := CreatePollingRequest(resp, NullAuthorizer{})

	resp, err := PollForAttempts(client, req, time.Duration(0), 1, 202)
	if err == nil {
		t.Error("autorest: Polling failed to emit error for known status code")
	}
}

func TestPollingReturnsNoErrorForUnexpectedStatusCode(t *testing.T) {
	client := mocks.NewSender()
	client.EmitStatus("500 ServerError", 500)

	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _ := CreatePollingRequest(resp, NullAuthorizer{})

	resp, err := PollForAttempts(client, req, time.Duration(0), 1, 202)
	if err != nil {
		t.Error("autorest: Polling emitted error for unknown status code")
	}
}

func TestPollingReturnsDefaultsToAcceptedStatusCode(t *testing.T) {
	client := mocks.NewSender()
	client.EmitStatus("202 Accepted", 202)

	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _ := CreatePollingRequest(resp, NullAuthorizer{})

	resp, err := PollForAttempts(client, req, time.Duration(0), 1)
	if err == nil {
		t.Error("autorest: Polling failed to default to HTTP 202")
	}
}

func TestPollingLeavesFinalBodyOpen(t *testing.T) {
	client := mocks.NewSender()
	client.EmitStatus("500 ServerError", 500)

	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _ := CreatePollingRequest(resp, NullAuthorizer{})

	resp, _ = PollForAttempts(client, req, time.Duration(0), 1)
	if !resp.Body.(*mocks.Body).IsOpen() {
		t.Error("autorest: Polling unexpectedly closed the response body")
	}
}

func TestDecorateForPollingCloseBodyOnEachAttempt(t *testing.T) {
	client := mocks.NewSender()
	client.EmitStatus("202 Accepted", 202)
	client.ReuseResponse(true)

	resp := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(resp)

	req, _ := CreatePollingRequest(resp, NullAuthorizer{})
	resp, _ = PollForAttempts(client, req, time.Duration(0), 5)
	if resp.Body.(*mocks.Body).CloseAttempts() != 5 {
		t.Errorf("autorest: decorateForPolling failed to close the response Body between requests -- expected %v, received %v",
			5, resp.Body.(*mocks.Body).CloseAttempts())
	}
}

func addAcceptedHeaders(resp *http.Response) {
	mocks.AddResponseHeader(resp, http.CanonicalHeaderKey(headerLocation), testUrl)
	mocks.AddResponseHeader(resp, http.CanonicalHeaderKey(headerRetryAfter), fmt.Sprintf("%v", int(testDelay.Seconds())))
}

type TestAuthorizer struct{}

func (ta TestAuthorizer) WithAuthorization() PrepareDecorator {
	return WithHeader(headerAuthorization, testAuthorizationHeader)
}
