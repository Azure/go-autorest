/*
Package azure provides Azure-specific implementations used with AutoRest.

See the included examples for more detail.
*/
package azure

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Azure/go-autorest/autorest"
)

const (
	// HeaderClientID is the Azure extension header to set a user-specified request ID.
	HeaderClientID = "x-ms-client-request-id"

	// HeaderReturnClientID is the Azure extension header to set if the user-specified request ID
	// should be included in the response.
	HeaderReturnClientID = "x-ms-return-client-request-id"

	// HeaderRequestID is the Azure extension header of the service generated request ID returned
	// in the response.
	HeaderRequestID = "x-ms-request-id"
)

// Error describes an error response returned by Azure service.
type Error struct {
	ServiceError *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"` // JSON object returned in service response body
	StatusCode int    // HTTP response status code
	RequestID  string // x-ms-request-id-header
}

// Error returns a human-friendly error message from service error.
func (e *Error) Error() string {
	return fmt.Sprintf("azure: Service returned an error. Code=%q Message=%q Status=%d",
		e.ServiceError.Code, e.ServiceError.Message, e.StatusCode)
}

// WithReturningClientID returns a PrepareDecorator that adds an HTTP extension header of
// x-ms-client-request-id whose value is the passed, undecorated UUID (e.g.,
// "0F39878C-5F76-4DB8-A25D-61D2C193C3CA"). It also sets the x-ms-return-client-request-id
// header to true such that UUID accompanies the http.Response.
func WithReturningClientID(uuid string) autorest.PrepareDecorator {
	preparer := autorest.CreatePreparer(
		WithClientID(uuid),
		WithReturnClientID(true))

	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)
			if err != nil {
				return r, err
			}
			return preparer.Prepare(r)
		})
	}
}

// WithClientID returns a PrepareDecorator that adds an HTTP extension header of
// x-ms-client-request-id whose value is passed, undecorated UUID (e.g.,
// "0F39878C-5F76-4DB8-A25D-61D2C193C3CA").
func WithClientID(uuid string) autorest.PrepareDecorator {
	return autorest.WithHeader(HeaderClientID, uuid)
}

// WithReturnClientID returns a PrepareDecorator that adds an HTTP extension header of
// x-ms-return-client-request-id whose boolean value indicates if the value of the
// x-ms-client-request-id header should be included in the http.Response.
func WithReturnClientID(b bool) autorest.PrepareDecorator {
	return autorest.WithHeader(HeaderReturnClientID, strconv.FormatBool(b))
}

// ExtractClientID extracts the client identifier from the x-ms-client-request-id header set on the
// http.Request sent to the service (and returned in the http.Response)
func ExtractClientID(resp *http.Response) string {
	return autorest.ExtractHeaderValue(HeaderClientID, resp)
}

// ExtractRequestID extracts the Azure server generated request identifier from the
// x-ms-request-id header.
func ExtractRequestID(resp *http.Response) string {
	return autorest.ExtractHeaderValue(HeaderRequestID, resp)
}

// WithErrorUnlessStatusCode returns a RespondDecorator that emits an
// azure.Error by reading the response body unless the response HTTP status code
// is among the set passed.
//
// If there is a chance service may return responses other than the Azure error
// format and the response cannot be parsed into an error, a decoding error will
// be returned containing the response body. In any case, the Responder will
// return an error if the status code is not satisfied.
//
// If this Responder returns an error, the response body will be replaced with
// an in-memory reader, which needs no further closing.
func WithErrorUnlessStatusCode(codes ...int) autorest.RespondDecorator {
	return func(r autorest.Responder) autorest.Responder {
		return autorest.ResponderFunc(func(resp *http.Response) error {
			err := r.Respond(resp)
			if err == nil && !autorest.ResponseHasStatusCode(resp, codes...) {
				var e Error
				defer resp.Body.Close()

				b, decodeErr := autorest.CopyAndDecode(autorest.EncodedAsJSON, resp.Body, &e)
				resp.Body = ioutil.NopCloser(&b) // replace body with in-memory reader
				if decodeErr != nil || e.ServiceError == nil {
					return fmt.Errorf("autorest/azure: error response cannot be parsed: %q error: %v", b.String(), err)
				}

				e.RequestID = ExtractRequestID(resp)
				e.StatusCode = resp.StatusCode
				err = &e
			}
			return err
		})
	}
}
