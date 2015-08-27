package autorest

import (
	"net/http"
	"time"
)

const (
	// DefaultPollingDelay is the default delay between polling requests (only used if the
	// http.Request lacks a well-formed Retry-After header).
	DefaultPollingDelay = 60 * time.Second

	// DefaultPollingDuration is the default total polling duration.
	DefaultPollingDuration = 15 * time.Minute
)

// PollingMode sets how, if at all, clients composed with Client will poll.
type PollingMode string

const (
	// PollUntilAttempts polling mode polls until reaching a maximum number of attempts.
	PollUntilAttempts PollingMode = "poll-until-attempts"

	// PollUntilDuration polling mode polls until a specified time.Duration has passed.
	PollUntilDuration PollingMode = "poll-until-duration"

	// DoNotPoll disables polling.
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
	// DefaultClient is the base from which generated clients should create a Client instance. Users
	// can then established widely used Client defaults by replacing or modifying the DefaultClient
	// before instantiating a generated client.
	DefaultClient = Client{PollingMode: PollUntilDuration, PollingDuration: DefaultPollingDuration}
)

// Client is the base for autorest generated clients. It provides default, "do nothing"
// implementations of an Authorizer, RequestInspector, and ResponseInspector. It also returns the
// standard, undecorated http.Client as a default Sender. Lastly, it supports basic request polling,
// limited to a maximum number of attempts or a specified duration.
//
// Generated clients should also use Error (see NewError and NewErrorWithError) for errors and
// return responses that compose with Response.
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

// IsPollingAllowed returns an error if the client allows polling and the passed http.Response
// requires it, otherwise it returns nil.
func (c Client) IsPollingAllowed(resp *http.Response, codes ...int) error {
	if c.DoNotPoll() && ResponseRequiresPolling(resp, codes...) {
		return NewError("autorest.Client", "IsPollingAllowed", "Response to %s requires polling but polling is disabled",
			resp.Request.URL)
	}
	return nil
}

// PollAsNeeded is a convenience method that will poll if the passed http.Response requires it.
func (c Client) PollAsNeeded(resp *http.Response, codes ...int) (*http.Response, error) {
	if !ResponseRequiresPolling(resp, codes...) {
		return resp, nil
	}

	if c.DoNotPoll() {
		return resp, NewError("autorest.Client", "PollAsNeeded", "Polling for %s is required, but polling is disabled",
			resp.Request.URL)
	}

	req, err := NewPollingRequest(resp, c)
	if err != nil {
		return resp, NewErrorWithError(err, "autorest.Client", "PollAsNeeded", "Unable to create polling request for response to %s",
			resp.Request.URL)
	}

	Prepare(req,
		c.WithInspection())

	if c.PollForAttempts() {
		return PollForAttempts(c, req, DefaultPollingDelay, c.PollingAttempts, codes...)
	}
	return PollForDuration(c, req, DefaultPollingDelay, c.PollingDuration, codes...)
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

// Do is a convenience method that invokes the Sender of the Client. If no Sender is set, it uses
// a new instance of http.Client.
func (c Client) Do(r *http.Request) (*http.Response, error) {
	return c.sender().Do(r)
}

// sender returns the Sender to which to send requests.
func (c Client) sender() Sender {
	if c.Sender == nil {
		return &http.Client{}
	}
	return c.Sender
}

// WithAuthorization is a convenience method that returns the WithAuthorization PrepareDecorator
// from the current Authorizer. If not Authorizer is set, it uses the NullAuthorizer.
func (c Client) WithAuthorization() PrepareDecorator {
	return c.authorizer().WithAuthorization()
}

// authorizer returns the Authorizer to use.
func (c Client) authorizer() Authorizer {
	if c.Authorizer == nil {
		return NullAuthorizer{}
	}
	return c.Authorizer
}

// WithInspection is a convenience method that passes the request to the supplied RequestInspector,
// if present, or returns the WithNothing PrepareDecorator otherwise.
func (c Client) WithInspection() PrepareDecorator {
	if c.RequestInspector == nil {
		return WithNothing()
	}
	return c.RequestInspector.WithInspection()
}

// ByInspecting is a convenience method that passes the response to the supplied ResponseInspector,
// if present, or returns the ByIgnoring RespondDecorator otherwise.
func (c Client) ByInspecting() RespondDecorator {
	if c.ResponseInspector == nil {
		return ByIgnoring()
	}
	return c.ResponseInspector.ByInspecting()
}

// Response serves as the base for all responses from generated clients. It provides access to the
// last http.Response.
type Response struct {
	*http.Response `json:"-"`
}

// GetPollingDelay extracts the polling delay from the Retry-After header of the response. If
// the header is absent or is malformed, it will return the supplied default delay time.Duration.
func (r Response) GetPollingDelay(defaultDelay time.Duration) time.Duration {
	return GetPollingDelay(r.Response, defaultDelay)
}

// GetPollingLocation retrieves the polling URL from the Location header of the response.
func (r Response) GetPollingLocation() string {
	return GetPollingLocation(r.Response)
}
