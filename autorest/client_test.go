package autorest

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/azure/go-autorest/autorest/mocks"
)

const (
	headerMockAuthorizer = "x-mock-authorizer"
)

func TestClientShouldPoll(t *testing.T) {
	c := &Client{PollingMode: PollUntilAttempts}
	r := mocks.NewResponseWithStatus("202 Accepted", 202)

	if !c.ShouldPoll(r) {
		t.Error("autorest: Client#ShouldPoll failed to return true for an http.Response that requires polling")
	}
}

func TestClientShouldPollIgnoresOk(t *testing.T) {
	c := &Client{PollingMode: PollUntilAttempts}
	r := mocks.NewResponse()

	if c.ShouldPoll(r) {
		t.Error("autorest: Client#ShouldPoll failed to return false for an http.Response that does not require polling")
	}
}

func TestClientShouldPollIgnoredPollingMode(t *testing.T) {
	c := &Client{PollingMode: DoNotPoll}
	r := mocks.NewResponse()

	if c.ShouldPoll(r) {
		t.Error("autorest: Client#ShouldPoll failed to return false when polling is disabled")
	}
}

func TestClientPollAsNeededIgnoresOk(t *testing.T) {
	c := &Client{}
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
	c := &Client{}
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

func TestClientPollAsNeededPollsForAttempts(t *testing.T) {
	c := &Client{}
	c.PollingMode = PollUntilAttempts
	c.PollingAttempts = 5

	s := mocks.NewSender()
	s.EmitStatus("202 Accepted", 202)
	c.Sender = s

	r := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(r)

	resp, _ := c.PollAsNeeded(r)
	if s.Attempts() != 5 {
		t.Errorf("autorest: Client#PollAsNeeded did not poll the expected number of attempts -- expected %v, actual %v",
			c.PollingAttempts, s.Attempts())
	}

	Respond(resp,
		ByClosing())
}

func TestClientPollAsNeededPollsForDuration(t *testing.T) {
	c := &Client{}
	c.PollingMode = PollUntilDuration
	c.PollingDuration = 10 * time.Millisecond

	s := mocks.NewSender()
	s.EmitStatus("202 Accepted", 202)
	c.Sender = s

	r := mocks.NewResponseWithStatus("202 Accepted", 202)
	addAcceptedHeaders(r)

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
	c := &Client{}

	if !c.DoNotPoll() {
		t.Errorf("autorest: Client requested polling by default, expected no polling (%v)", c.PollingMode)
	}
}

func TestClientDoNotPollForAttempts(t *testing.T) {
	c := &Client{}
	c.PollingMode = PollUntilAttempts

	if c.DoNotPoll() {
		t.Errorf("autorest: Client failed to request polling after polling mode set to %s", c.PollingMode)
	}
}

func TestClientDoNotPollForDuration(t *testing.T) {
	c := &Client{}
	c.PollingMode = PollUntilDuration

	if c.DoNotPoll() {
		t.Errorf("autorest: Client failed to request polling after polling mode set to %s", c.PollingMode)
	}
}

func TestClientPollForAttempts(t *testing.T) {
	c := &Client{}
	c.PollingMode = PollUntilAttempts

	if !c.PollForAttempts() {
		t.Errorf("autorest: Client#SetPollingMode failed to set polling by attempts -- polling mode set to %s", c.PollingMode)
	}
}

func TestClientPollForDuration(t *testing.T) {
	c := &Client{}
	c.PollingMode = PollUntilDuration

	if !c.PollForDuration() {
		t.Errorf("autorest: Client#SetPollingMode failed to set polling for duration -- polling mode set to %s", c.PollingMode)
	}
}

func TestClientDoSetsDefaultSender(t *testing.T) {
	c := &Client{}

	c.Do(&http.Request{})
	if !reflect.DeepEqual(c.Sender, &http.Client{}) {
		t.Error("autorest: Client#Do failed to set the default Sender to the http.Client")
	}
}

func TestClientWithAuthorizer(t *testing.T) {
	c := &Client{}
	c.Authorizer = mockAuthorizer{}

	req, _ := Prepare(&http.Request{},
		c.WithAuthorization())

	if req.Header.Get(headerMockAuthorizer) == "" {
		t.Error("autorest: Client#WithAuthorizer failed to return the WithAuthorizer from the active Authorizer")
	}
}

func TestClientWithAuthorizerSetsDefaultAuthorizer(t *testing.T) {
	c := &Client{}

	Prepare(&http.Request{},
		c.WithAuthorization())

	if !reflect.DeepEqual(c.Authorizer, NullAuthorizer{}) {
		t.Error("autorest: Client#WithAuthorizer failed to set the default Authorizer to the NullAuthorizer")
	}
}

func TestClientWithInspection(t *testing.T) {
	c := &Client{}
	r := &mockInspector{}
	c.RequestInspector = r

	Prepare(&http.Request{},
		c.WithInspection())

	if !r.wasInvoked {
		t.Error("autorest: Client#WithInspection failed to invoke RequestInspector")
	}
}

func TestClientWithInspectionSetsDefault(t *testing.T) {
	c := &Client{}

	r1 := &http.Request{}
	r2, _ := Prepare(r1,
		c.WithInspection())

	if !reflect.DeepEqual(r1, r2) {
		t.Error("autorest: Client#WithInspection failed to provide a default RequestInspector")
	}
}

func TestClientByInspecting(t *testing.T) {
	c := &Client{}
	r := &mockInspector{}
	c.ResponseInspector = r

	Respond(&http.Response{},
		c.ByInspecting())

	if !r.wasInvoked {
		t.Error("autorest: Client#ByInspecting failed to invoke ResponseInspector")
	}
}

func TestClientByInspectingSetsDefault(t *testing.T) {
	c := &Client{}

	r := &http.Response{}
	Respond(r,
		c.ByInspecting())

	if !reflect.DeepEqual(r, &http.Response{}) {
		t.Error("autorest: Client#ByInspecting failed to provide a default ResponseInspector")
	}
}

type mockAuthorizer struct{}

func (ma mockAuthorizer) WithAuthorization() PrepareDecorator {
	return WithHeader(headerMockAuthorizer, "authorized")
}

type mockInspector struct {
	wasInvoked bool
}

func (mi *mockInspector) WithInspection() PrepareDecorator {
	return func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			mi.wasInvoked = true
			return p.Prepare(r)
		})
	}
}

func (mi *mockInspector) ByInspecting() RespondDecorator {
	return func(r Responder) Responder {
		return ResponderFunc(func(resp *http.Response) error {
			mi.wasInvoked = true
			return r.Respond(resp)
		})
	}
}
