package autorest

import (
	"fmt"
	"net/http"
	"time"
)

// PollingMode sets how, if at all, clients composed with Client will poll.
type PollingMode string

const (
	// Poll until reaching a maximum number of attempts
	PollUntilAttempts PollingMode = "poll-until-attempts"

	// Poll until a specified time.Duration has passed
	PollUntilDuration PollingMode = "poll-until-duration"

	// Do not poll at all
	DoNotPoll PollingMode = "not-at-all"
)

// RequestInspector defines a single method that returns a PrepareDecorator used to inspect the
// http.Request prior to sending.
type RequestInspector interface {
	WithInspection() PrepareDecorator
}

// ResponseInspector defines a single method that returns a ResponseDecorator used to inspect the
// http.Response prior to responding.
type ResponseInspector interface {
	ByInspecting() RespondDecorator
}

// Client is a convenience base for autorest generated clients. It provides default, "do nothing"
// implementations of an Authorizer, RequestInspector, and ResponseInspector. It also returns the
// standard, undecorated http.Client as a default Sender. Lastly, it supports basic request polling,
// limited to a maximum number of attempts or a specified duration.
//
// Most uses of generated clients can customize behavior or inspect requests through by supplying a
// custom Authorizer, custom RequestInspector, and / or custom ResponseInspector. Users may log
// requests, implement circuit breakers (see https://msdn.microsoft.com/en-us/library/dn589784.aspx)
// or otherwise influence sending the request by providing a decorated Sender.
type Client struct {
	Authorizer        Authorizer
	Sender            Sender
	RequestInspector  RequestInspector
	ResponseInspector ResponseInspector

	PollingMode     PollingMode
	PollingAttempts int
	PollingDuration time.Duration
}

// PollIfNeeded is a convenience method that will poll if the passed http.Response requires it.
func (c *Client) PollIfNeeded(resp *http.Response, codes ...int) (*http.Response, error) {
	req, delay, err := CreatePollingRequest(resp, c, codes...)
	if err != nil {
		return resp, fmt.Errorf("autorest: Failed to create Poll request (%v)", err)
	}
	if req == nil {
		return resp, nil
	}

	if c.PollForAttempts() {
		return PollForAttempts(c, req, delay, c.PollingAttempts, codes...)
	} else if c.PollForDuration() {
		return PollForDuration(c, req, delay, c.PollingDuration, codes...)
	} else {
		return resp, fmt.Errorf("autorest: Polling for %s is required, but polling is disabled", req.URL)
	}
}

// DoNotPoll returns true if the client should not poll, false otherwise.
func (c Client) DoNotPoll() bool {
	return len(c.PollingMode) == 0 || c.PollingMode == DoNotPoll
}

// PollForAttempts returns true if the PollingMode is set to ForAttempts, false otherwise.
func (c Client) PollForAttempts() bool {
	return c.PollingMode == PollUntilAttempts
}

// PollForDuration return true if the PollingMode is set to ForDuration, false otherwise.
func (c Client) PollForDuration() bool {
	return c.PollingMode == PollUntilDuration
}

// Do is a convenience method that invokes the Sender of the Client. It will use the default
// http.Client if no Sender is set.
func (c *Client) Do(r *http.Request) (*http.Response, error) {
	if c.Sender == nil {
		c.Sender = &http.Client{}
	}
	return c.Sender.Do(r)
}

// WithAuthorization is a convenience method that returns the WithAuthorization PrepareDecorator
// from the current Authorizer. If not Authorizer is set, it returns the WithAuthorization
// PrepareDecorator from the NullAuthorizer.
func (c *Client) WithAuthorization() PrepareDecorator {
	if c.Authorizer == nil {
		c.Authorizer = NullAuthorizer{}
	}
	return c.Authorizer.WithAuthorization()
}

// WithInspection is a convenience method that passes the request to the supplied RequestInspector,
// if present, or returns the WithNothing PrepareDecorator otherwise.
func (c *Client) WithInspection() PrepareDecorator {
	if c.RequestInspector == nil {
		return WithNothing()
	} else {
		return c.RequestInspector.WithInspection()
	}
}

// ByInspecting is a convenience method that passes the response to the supplied ResponseInspector,
// if present, or returns the ByIgnoring RespondDecorator otherwise.
func (c *Client) ByInspecting() RespondDecorator {
	if c.ResponseInspector == nil {
		return ByIgnoring()
	} else {
		return c.ResponseInspector.ByInspecting()
	}
}
