package autorest

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/azure/go-autorest/autorest/mocks"
)

func TestClientIsPollingAllowed(t *testing.T) {
	c := Client{PollingMode: PollUntilAttempts}
	r := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)

	err := c.IsPollingAllowed(r)
	if err != nil {
		t.Errorf("autorest: Client#IsPollingAllowed returned an error for an http.Response that requires polling (%v)", err)
	}
}

func TestClientIsPollingAllowedIgnoresOk(t *testing.T) {
	c := Client{PollingMode: PollUntilAttempts}
	r := mocks.NewResponse()

	err := c.IsPollingAllowed(r)
	if err != nil {
		t.Error("autorest: Client#IsPollingAllowed returned an error for an http.Response that does not require polling")
	}
}

func TestClientIsPollingAllowedIgnoresDisabledForIgnoredStatusCode(t *testing.T) {
	c := Client{PollingMode: PollUntilAttempts}
	r := mocks.NewResponseWithStatus("400 BadRequest", http.StatusBadRequest)

	err := c.IsPollingAllowed(r)
	if err != nil {
		t.Errorf("autorest: Client#IsPollingAllowed returned an error for an http.Response that requires polling (%v)", err)
	}
}

func TestClientIsPollingAllowedIgnoredPollingMode(t *testing.T) {
	c := Client{PollingMode: DoNotPoll}
	r := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)

	err := c.IsPollingAllowed(r)
	if err == nil {
		t.Error("autorest: Client#IsPollingAllowed failed to return an error when polling is disabled")
	}
}

func TestClientPollAsNeededIgnoresOk(t *testing.T) {
	c := Client{}
	s := mocks.NewSender()
	c.Sender = s
	r := mocks.NewResponse()

	resp, err := c.PollAsNeeded(r)
	if err != nil {
		t.Errorf("autorest: Client#PollAsNeeded failed when given a successful HTTP request (%v)", err)
	}
	if s.Attempts() > 0 {
		t.Error("autorest: Client#PollAsNeeded attempted to poll a successful HTTP request")
	}

	Respond(resp,
		ByClosing())
}

func TestClientPollAsNeededLeavesBodyOpen(t *testing.T) {
	c := Client{}
	c.Sender = mocks.NewSender()
	r := mocks.NewResponse()

	resp, err := c.PollAsNeeded(r)
	if err != nil {
		t.Errorf("autorest: Client#PollAsNeeded failed when given a successful HTTP request (%v)", err)
	}
	if !resp.Body.(*mocks.Body).IsOpen() {
		t.Error("autorest: Client#PollAsNeeded unexpectedly closed the response body")
	}

	Respond(resp,
		ByClosing())
}

func TestClientPollAsNeededReturnsErrorWhenPollingDisabled(t *testing.T) {
	c := Client{}
	c.Sender = mocks.NewSender()
	c.PollingMode = DoNotPoll

	r := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	setAcceptedHeaders(r)

	_, err := c.PollAsNeeded(r)
	if err == nil {
		t.Error("autorest: Client#PollAsNeeded failed to return an error when polling was disabled but required")
	}

	Respond(r,
		ByClosing())
}

func TestClientPollAsNeededAllowsInspectionOfRequest(t *testing.T) {
	c := Client{}
	c.Sender = mocks.NewSender()
	c.PollingMode = PollUntilAttempts
	c.PollingAttempts = 1

	mi := &mockInspector{}
	c.RequestInspector = mi

	r := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	setAcceptedHeaders(r)

	c.PollAsNeeded(r)
	if !mi.wasInvoked {
		t.Error("autorest: Client#PollAsNeeded failed to allow inspection of polling request")
	}

	Respond(r,
		ByClosing())
}

func TestClientPollAsNeededReturnsErrorIfUnableToCreateRequest(t *testing.T) {
	c := Client{}
	c.Authorizer = mockFailingAuthorizer{}
	c.Sender = mocks.NewSender()
	c.PollingMode = PollUntilAttempts
	c.PollingAttempts = 1

	r := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	setAcceptedHeaders(r)

	_, err := c.PollAsNeeded(r)
	if err == nil {
		t.Error("autorest: Client#PollAsNeeded failed to return error when unable to create request")
	}

	Respond(r,
		ByClosing())
}

func TestClientPollAsNeededPollsForAttempts(t *testing.T) {
	c := Client{}
	c.PollingMode = PollUntilAttempts
	c.PollingAttempts = 5

	s := mocks.NewSender()
	s.EmitStatus("202 Accepted", http.StatusAccepted)
	c.Sender = s

	r := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	setAcceptedHeaders(r)
	s.SetResponse(r)
	s.ReuseResponse(true)

	resp, _ := c.PollAsNeeded(r)
	if s.Attempts() != 5 {
		t.Errorf("autorest: Client#PollAsNeeded did not poll the expected number of attempts -- expected %v, actual %v",
			c.PollingAttempts, s.Attempts())
	}

	Respond(resp,
		ByClosing())
}

func TestClientPollAsNeededPollsForDuration(t *testing.T) {
	c := Client{}
	c.PollingMode = PollUntilDuration
	c.PollingDuration = 10 * time.Millisecond

	s := mocks.NewSender()
	s.EmitStatus("202 Accepted", http.StatusAccepted)
	c.Sender = s

	r := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	setAcceptedHeaders(r)
	s.SetResponse(r)
	s.ReuseResponse(true)

	d1 := 10 * time.Millisecond
	start := time.Now()
	resp, _ := c.PollAsNeeded(r)
	d2 := time.Now().Sub(start)
	if d2 < d1 {
		t.Errorf("autorest: Client#PollAsNeeded did not poll for the expected duration -- expected %v, actual %v",
			d1.Seconds(), d2.Seconds())
	}

	Respond(resp,
		ByClosing())
}

func TestClientDoNotPoll(t *testing.T) {
	c := Client{}

	if !c.DoNotPoll() {
		t.Errorf("autorest: Client requested polling by default, expected no polling (%v)", c.PollingMode)
	}
}

func TestClientDoNotPollForAttempts(t *testing.T) {
	c := Client{}
	c.PollingMode = PollUntilAttempts

	if c.DoNotPoll() {
		t.Errorf("autorest: Client failed to request polling after polling mode set to %s", c.PollingMode)
	}
}

func TestClientDoNotPollForDuration(t *testing.T) {
	c := Client{}
	c.PollingMode = PollUntilDuration

	if c.DoNotPoll() {
		t.Errorf("autorest: Client failed to request polling after polling mode set to %s", c.PollingMode)
	}
}

func TestClientPollForAttempts(t *testing.T) {
	c := Client{}
	c.PollingMode = PollUntilAttempts

	if !c.PollForAttempts() {
		t.Errorf("autorest: Client#SetPollingMode failed to set polling by attempts -- polling mode set to %s", c.PollingMode)
	}
}

func TestClientPollForDuration(t *testing.T) {
	c := Client{}
	c.PollingMode = PollUntilDuration

	if !c.PollForDuration() {
		t.Errorf("autorest: Client#SetPollingMode failed to set polling for duration -- polling mode set to %s", c.PollingMode)
	}
}

func TestClientSenderReturnsHttpClientByDefault(t *testing.T) {
	c := Client{}

	if fmt.Sprintf("%T", c.sender()) != "*http.Client" {
		t.Error("autorest: Client#sender failed to return http.Client by default")
	}
}

func TestClientSenderReturnsSetSender(t *testing.T) {
	c := Client{}

	s := mocks.NewSender()
	c.Sender = s

	if c.sender() != s {
		t.Error("autorest: Client#sender failed to return set Sender")
	}
}

func TestClientDoInvokesSender(t *testing.T) {
	c := Client{}

	s := mocks.NewSender()
	c.Sender = s

	c.Do(&http.Request{})
	if s.Attempts() != 1 {
		t.Error("autorest: Client#Do failed to invoke the Sender")
	}
}

func TestClientAuthorizerReturnsNullAuthorizerByDefault(t *testing.T) {
	c := Client{}

	if fmt.Sprintf("%T", c.authorizer()) != "autorest.NullAuthorizer" {
		t.Error("autorest: Client#authorizer failed to return the NullAuthorizer by default")
	}
}

func TestClientAuthorizerReturnsSetAuthorizer(t *testing.T) {
	c := Client{}
	c.Authorizer = mockAuthorizer{}

	if fmt.Sprintf("%T", c.authorizer()) != "autorest.mockAuthorizer" {
		t.Error("autorest: Client#authorizer failed to return the set Authorizer")
	}
}

func TestClientWithAuthorizer(t *testing.T) {
	c := Client{}
	c.Authorizer = mockAuthorizer{}

	req, _ := Prepare(&http.Request{},
		c.WithAuthorization())

	if req.Header.Get(headerAuthorization) == "" {
		t.Error("autorest: Client#WithAuthorizer failed to return the WithAuthorizer from the active Authorizer")
	}
}

func TestClientWithInspection(t *testing.T) {
	c := Client{}
	r := &mockInspector{}
	c.RequestInspector = r

	Prepare(&http.Request{},
		c.WithInspection())

	if !r.wasInvoked {
		t.Error("autorest: Client#WithInspection failed to invoke RequestInspector")
	}
}

func TestClientWithInspectionSetsDefault(t *testing.T) {
	c := Client{}

	r1 := &http.Request{}
	r2, _ := Prepare(r1,
		c.WithInspection())

	if !reflect.DeepEqual(r1, r2) {
		t.Error("autorest: Client#WithInspection failed to provide a default RequestInspector")
	}
}

func TestClientByInspecting(t *testing.T) {
	c := Client{}
	r := &mockInspector{}
	c.ResponseInspector = r

	Respond(&http.Response{},
		c.ByInspecting())

	if !r.wasInvoked {
		t.Error("autorest: Client#ByInspecting failed to invoke ResponseInspector")
	}
}

func TestClientByInspectingSetsDefault(t *testing.T) {
	c := Client{}

	r := &http.Response{}
	Respond(r,
		c.ByInspecting())

	if !reflect.DeepEqual(r, &http.Response{}) {
		t.Error("autorest: Client#ByInspecting failed to provide a default ResponseInspector")
	}
}

func TestResponseGetPollingDelay(t *testing.T) {
	d1 := 10 * time.Second

	r := mocks.NewResponse()
	setRetryHeader(r, d1)
	ar := Response{Response: r}

	d2 := ar.GetPollingDelay(time.Duration(0))
	if d1 != d2 {
		t.Errorf("autorest: Response#GetPollingDelay failed to return the correct delay -- expected %v, received %v",
			d1, d2)
	}
}

func TestResponseGetPollingDelayReturnsDefault(t *testing.T) {
	ar := Response{Response: mocks.NewResponse()}

	d1 := 10 * time.Second
	d2 := ar.GetPollingDelay(d1)
	if d1 != d2 {
		t.Errorf("autorest: Response#GetPollingDelay failed to return the default delay -- expected %v, received %v",
			d1, d2)
	}
}

func TestResponseGetPollingLocation(t *testing.T) {
	r := mocks.NewResponse()
	setLocationHeader(r, testURL)
	ar := Response{Response: r}

	if ar.GetPollingLocation() != testURL {
		t.Errorf("autorest: Response#GetPollingLocation failed to return correct URL -- expected %v, received %v",
			testURL, ar.GetPollingLocation())
	}
}
