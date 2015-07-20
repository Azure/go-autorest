package autorest

import (
	"log"
	"net/http"
	"time"
)

// Sender is the interface that wraps the Do method to send HTTP requests.
//
// The standard http.Client conforms to this interface.
type Sender interface {
	Do(*http.Request) (*http.Response, error)
}

// SenderFunc is a method that implements the Sender interface.
type SenderFunc func(*http.Request) (*http.Response, error)

// Do implements the Sender interface on SenderFunc.
func (sf SenderFunc) Do(r *http.Request) (*http.Response, error) {
	return sf(r)
}

// SendDecorator takes and possibily decorates, by wrapping, a Sender.
// Decorators may affect the http.Request and pass it along or, first, pass the http.Request along then
// react to the http.Response result.
// By convention, the names of SendDecorators should begin with "With."
type SendDecorator func(Sender) Sender

// CreateSender creates, decorates, and returns, as a Sender, the default http.Client.
func CreateSender(decorators ...SendDecorator) Sender {
	return DecorateSender(&http.Client{}, decorators...)
}

// DecorateSender accepts a Sender and a, possibly empty, set of SendDecorators, which is applies to the Sender.
// Decorators are applied in the order received, but their affect upon the request depends on whether they are a
// pre-decorator (change the http.Request and then pass it along) or a post-decorator (pass the http.Request along and
// react to the results in http.Response).
func DecorateSender(s Sender, decorators ...SendDecorator) Sender {
	for i := len(decorators) - 1; i >= 0; i-- {
		s = decorators[i](s)
	}
	return s
}

// Send sends, by means of the default http.Client, the passed http.Request, returning the http.Response and possible error.
// It also accepts a, possibly empty, set of SendDecorators which it will apply the http.Client before invoking the Do method.
//
// Send is a convenience method and not recommended for production. Advanced users should create and share their own Sender
// (e.g., instance of http.Client), applying SendDecorators (by means of DecorateSender) as needed.
func Send(r *http.Request, decorators ...SendDecorator) (*http.Response, error) {
	return CreateSender(decorators...).Do(r)
}

// WithLogging returns a SendDecorator that implements simple before and after logging of the request.
func WithLogging(logger log.Logger) SendDecorator {
	return func(s Sender) Sender {
		return SenderFunc(func(r *http.Request) (*http.Response, error) {
			log.Printf("autorest: Sending %s %s\n", r.Method, r.URL)
			resp, err := s.Do(r)
			log.Printf("autorest: %s %s received %s\n", r.Method, r.URL, resp.Status)
			return resp, err
		})
	}
}

// WithRetry returns a SendDecorator that implements simple retry logic (i.e., retry for any received error, up to a maximum
// number of attempts, exponentially backing off using the supplied backoff time.Duration).
func WithRetry(attempts int, backoff time.Duration) SendDecorator {
	return func(s Sender) Sender {
		return SenderFunc(func(r *http.Request) (resp *http.Response, err error) {
			for attempt := 0; attempt < attempts; attempt++ {
				if resp, err = s.Do(r); err == nil {
					break
				}
				time.Sleep(backoff * time.Duration(attempt))
			}
			return resp, err
		})
	}
}
