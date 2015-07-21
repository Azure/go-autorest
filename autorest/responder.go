package autorest

import (
	"encoding/json"
	"net/http"
)

// Responder is the interface that wraps the Respond method.
//
// Respond accepts and reacts to an http.Response.
// Implementations must ensure to not share or hold state since Responders may be shared and re-used.
type Responder interface {
	Respond(*http.Response) error
}

// ResponderFunc is a method that implements the Responder interface.
type ResponderFunc func(*http.Response) error

// Respond implements the Responder interface on ResponderFunc.
func (rf ResponderFunc) Respond(r *http.Response) error {
	return rf(r)
}

// RespondDecorator takes and possibly decorates, by wrapping, a Responder.
// Decorators may react to the http.Response and pass it along or, first, pass the http.Response along then react.
// By convention, the names of RespondDecorators should begin with "By."
type RespondDecorator func(Responder) Responder

// CreateResponder creates, decorates, and returns a Responder.
// Without decorators, the returned Responder returns the passed http.Response unmodified.
// Responders may or may not be safe to share and re-used: It depends on the applied decorators.
// For example, a standard decorator that closes the response body is fine to share whereas a decorator that
// reads the body into a passed struct is not.
//
// To prevent memory leaks, ensure that at least one Responder closes the response body. By convention, decorators
// that close the response body end in "Closing" (e.g., ByUnmarshallingJsonAndClosing).
func CreateResponder(decorators ...RespondDecorator) Responder {
	return DecorateResponder(Responder(ResponderFunc(func(r *http.Response) error { return nil })), decorators...)
}

// DecorateResponder accepts a Responder and a, possibly empty, set of RespondDecorators, which it applies to the Responder.
// Decorators are applied in the order received, but their affect upon the request depends on whether they are a
// pre-decorator (react to the http.Response and then pass it along) or a post-decorator (pass the http.Response along and
// then react).
func DecorateResponder(r Responder, decorators ...RespondDecorator) Responder {
	for i := len(decorators) - 1; i >= 0; i-- {
		r = decorators[i](r)
	}
	return r
}

// Respond accepts an http.Response and a, possibly empty, set of RespondDecorators.
// It creates a Responder from the decorators it then applies to the passed http.Response.
func Respond(r *http.Response, decorators ...RespondDecorator) error {
	return CreateResponder(decorators...).Respond(r)
}

// ByClosing returns a RespondDecorator that first invokes the passed Responder, passing the http.Response, after
// which it closes the http.Response.Body.
func ByClosing() RespondDecorator {
	return func(r Responder) Responder {
		return ResponderFunc(func(resp *http.Response) error {
			err := r.Respond(resp)
			if err == nil {
				resp.Body.Close()
			}
			return nil
		})
	}
}

// ByUnmarshallingJsonAndClosing returns a RespondDecorator that decodes a JSON document returned in the response Body
// into the value pointed to by v, after which is closes the response body.
func ByUnmarshallingJsonAndClosing(v interface{}) RespondDecorator {
	return func(r Responder) Responder {
		return ResponderFunc(func(resp *http.Response) error {
			defer resp.Body.Close()
			d := json.NewDecoder(resp.Body)
			err := d.Decode(v)
			if err != nil {
				return err
			}
			return r.Respond(resp)
		})
	}
}
