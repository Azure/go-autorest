package autorest

import (
	"fmt"
	"net/http"
	"time"
)

const (
	// The default delay between polling requests (only used if the http.Request lacks a well-formed
	// Retry-After header).
	DefaultPollingDelay = 60 * time.Second

	// The default total polling duration.
	DefaultPollingDuration = 10 * time.Minute
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

var (
	// Generated clients should compose using the DefaultClient instead of allocating a new Client
	// instance. Users can then established widely used Client defaults by replacing or modifying the
	// DefaultClient before instantiating a generated client.
	DefaultClient = &Client{PollingMode: PollUntilDuration, PollingDuration: DefaultPollingDuration}
)

// Client is the base for autorest generated clients. It provides default, "do nothing"
// implementations of an Authorizer, RequestInspector, and ResponseInspector. It also returns the
// standard, undecorated http.Client as a default Sender. Lastly, it supports basic request polling,
// limited to a maximum number of attempts or a specified duration.
//
// Most customization of generated clients is best achieved by supplying a custom Authorizer, custom
// RequestInspector, and / or custom ResponseInspector. Users may log requests, implement circuit
// breakers (see https://msdn.microsoft.com/en-us/library/dn589784.aspx) or otherwise influence
// sending the request by providing a decorated Sender.
type Client struct {
	Authorizer        Authorizer
	Sender            Sender
	RequestInspector  RequestInspector
	ResponseInspector ResponseInspector

	PollingMode     PollingMode
	PollingAttempts int
	PollingDuration time.Duration
}

// ShouldPoll returns true if the client allows polling and the passed http.Response requires it,
// otherwise it returns false.
func (c *Client) ShouldPoll(resp *http.Response, codes ...int) bool {
	return !c.DoNotPoll() && ResponseRequiresPolling(resp, codes...)
}

// PollAsNeeded is a convenience method that will poll if the passed http.Response requires it.
func (c *Client) PollAsNeeded(resp *http.Response, codes ...int) (*http.Response, error) {
	if !ResponseRequiresPolling(resp, codes...) {
		return resp, nil
	}

	req, err := CreatePollingRequest(resp, c)
	if err != nil {
		return resp, fmt.Errorf("autorest: Unable to create polling request for response to %s (%v)",
			resp.Request.URL, err)
	}

	delay := GetRetryDelay(resp, DefaultPollingDelay)

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

// Do is a convenience method that invokes the Sender of the Client. If no Sender is set, it will
// be set to the default http.Client.
func (c *Client) Do(r *http.Request) (*http.Response, error) {
	if c.Sender == nil {
		c.Sender = &http.Client{}
	}
	return c.Sender.Do(r)
}

// WithAuthorization is a convenience method that returns the WithAuthorization PrepareDecorator
// from the current Authorizer. If not Authorizer is set, it sets it to the NullAuthorizer.
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
