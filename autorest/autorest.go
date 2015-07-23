/*

  This package implements an HTTP request pipeline suitable for use across multiple go-routines
  and provides the shared routines relied on by AutoRest (see https://github.com/azure/autorest/)
  generated Go code.

  The package breaks sending and responding to HTTP requests into three phases: Preparing, Sending,
  and Responding. A typical pattern is:

    req, err := Prepare(&http.Request{},
      WithAuthorization())

    resp, err := Send(req,
      WithLogging(logger),
      DoErrorIfStatusCode(500),
      DoCloseIfError(),
      DoRetryForAttempts(5, time.Second))

    err = Respond(resp,
      ByClosing())

  Each phase relies on decorators to modify and / or manage processing. Decorators may first modify
  and then pass the data along, pass the data first and then modify the result, or wrap themselves
  around passing the data (such as a logger might do). Decorators run in the order provided. For
  example, the following:

    req, err := Prepare(&http.Request{},
      WithBaseURL("https://microsoft.com/"),
      WithPath("a"),
      WithPath("b"),
      WithPath("c"))

  will set the URL to:

    https://microsoft.com/a/b/c

  Preparers and Responders may be shared and re-used (assuming the underlying decorators support
  sharing and re-use). Performant use is obtained by creating one or more Preparers and Responders
  shared among multiple go-routines, and a single Sender shared among multiple sending go-routines,
  all bound together by means of input / output channels.

  Decorators hold their passed state within a closure (such as the path components in the example
  above). Be careful to share Preparers and Responders only in a context where such held state
  applies. For example, it may not make sense to share a Preparer that applies a query string from a
  fixed set of values. Similarly, sharing a Responder that reads the response body into a passed
  struct (e.g., ByUnmarshallingJson) is likely incorrect.

  Lastly, the Swagger specification (https://swagger.io) that drives AutoRest
  (https://github.com/azure/autorest/) precisely defines two date forms (i.e., date and date-time).
  The two sub-packages -- github.com/azure/go-autorest/autorest/date and
  github.com/azure/go-autorest/autorest/datetime -- provide time.Time derivations to ensure correct
  parsing and formatting.

  See the included examples for more detail.

*/
package autorest

import (
	"fmt"
	"net/http"
	"time"
)

const (
	headerLocation   = "Location"
	headerRetryAfter = "Retry-After"
)

// CreatePollingRequest allocates and returns a new http.Request, along with suggested a retry
// delay, if the supplied http.Response code is with the passed set (the set defaults to an
// HTTP 202). The http.Response must include both Location and Retry-After headers. If the passed
// http.Response is to be polled, this method will close the http.Response body. A passed http.Response
// with an HTTP 200 status code is ignored.
func CreatePollingRequest(resp *http.Response, authorizer Authorizer, codes ...int) (*http.Request, time.Duration, error) {
	d := time.Duration(0)

	if resp.StatusCode == 200 {
		return nil, d, nil
	}

	if len(codes) == 0 {
		codes = []int{202}
	}

	if !ResponseHasStatusCode(resp, codes...) {
		return nil, d, fmt.Errorf("autorest: Poll requested for unexpected status code -- %v not in %v", resp.StatusCode, codes)
	}

	location := resp.Header.Get(headerLocation)
	if location == "" {
		return nil, d, fmt.Errorf("autorest: Missing Location header in poll request to %s", resp.Request.URL)
	}

	retry := resp.Header.Get(headerRetryAfter)
	if retry == "" {
		return nil, d, fmt.Errorf("autorest: Missing Retry-After header in poll request to %s", resp.Request.URL)
	}

	d, err := time.ParseDuration(retry + "s")
	if err != nil {
		return nil, d, fmt.Errorf("autorest: Failed to parse retry duration (%s) in poll request to %s -- (%v)", retry, location, err)
	}

	req, err := Prepare(&http.Request{},
		AsGet(),
		WithBaseURL(location),
		authorizer.WithAuthorization())
	if err != nil {
		return nil, d, fmt.Errorf("autorest: Failure creating poll request to %s (%v)", location, err)
	}

	Respond(resp,
		ByClosing())

	return req, d, nil
}

// PollForAttempts will retry the passed http.Request until it receives an HTTP status code outside
// the passed set or has made the specified number of attempts. The set of status codes defaults to
// HTTP 202.
func PollForAttempts(s Sender, req *http.Request, delay time.Duration, attempts int, codes ...int) (*http.Response, error) {
	return SendWithSender(
		decorateForPolling(s, delay, codes...),
		req,
		DoRetryForAttempts(attempts, time.Duration(0)))
}

// PollForDuration will retry the passed http.Request until it receives an HTTP status code outside
// the passed set or the total time meets or exceeds the specified duration. The set of status codes
// defaults to HTTP 202.
func PollForDuration(s Sender, req *http.Request, delay time.Duration, total time.Duration, codes ...int) (*http.Response, error) {
	return SendWithSender(
		decorateForPolling(s, delay, codes...),
		req,
		DoRetryForDuration(total, time.Duration(0)))
}

func decorateForPolling(s Sender, delay time.Duration, codes ...int) Sender {
	if len(codes) == 0 {
		codes = []int{202}
	}

	return DecorateSender(s,
		AfterDelay(delay),
		DoCloseIfError(),
		DoErrorIfStatusCode(codes...))
}

// ResponseHasStatusCode returns true if the status code in the HTTP Response is in the passed set
// and false otherwise.
func ResponseHasStatusCode(resp *http.Response, codes ...int) bool {
	return containsInt(codes, resp.StatusCode)
}
