package autorest

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	mimeTypeJson = "application/json"

	headerAuthorization = "Authorization"
	headerContentType   = "Content-Type"
)

// Preparer is the interface that wraps the Prepare method.
//
// Prepare accepts and possibly modifies an http.Request (e.g., adding Headers).
// Implementations must ensure to not share or hold state since Preparers may be shared and re-used.
type Preparer interface {
	Prepare(*http.Request) (*http.Request, error)
}

// PreparerFunc is a method that implements the Preparer interface.
type PreparerFunc func(*http.Request) (*http.Request, error)

// Prepare implements the Preparer interface on PreparerFunc.
func (pf PreparerFunc) Prepare(r *http.Request) (*http.Request, error) {
	return pf(r)
}

// PrepareDecorator takes and possibly decorates, by wrapping, a Preparer.
// Decorators may affect the http.Request and pass it along or, first, pass the http.Request along then
// affect the result.
// By convention, the names of PrepareDecorators should begin with "With."
type PrepareDecorator func(Preparer) Preparer

// CreatePreparer creates, decorates, and returns a Preparer.
// Without decorators, the returned Preparer returns the passed http.Request unmodified.
// Preparers are safe to share and re-use.
func CreatePreparer(decorators ...PrepareDecorator) Preparer {
	return DecoratePreparer(Preparer(PreparerFunc(func(r *http.Request) (*http.Request, error) { return r, nil })), decorators...)
}

// DecoratePreparer accepts a Preparer and a, possibly empty, set of PrepareDecorators, which it applies to the Preparer.
// Decorators are applied in the order received, but their affect upon the request depends on whether they are a
// pre-decorator (change the http.Request and then pass it along) or a post-decorator (pass the http.Request along and
// alter it on return).
func DecoratePreparer(p Preparer, decorators ...PrepareDecorator) Preparer {
	for i := len(decorators) - 1; i >= 0; i-- {
		p = decorators[i](p)
	}
	return p
}

// Prepare accepts an http.Request and a, possibly empty, set of PrepareDecorators.
// It creates a Preparer from the decorators it then applies to the passed http.Request.
func Prepare(r *http.Request, decorators ...PrepareDecorator) (*http.Request, error) {
	return CreatePreparer(decorators...).Prepare(r)
}

// Header returns a PrepareDecorator that adds the specified HTTP header and value to the http.Request.
// It will canonicalize the passed header name (via http.CanonicalHeaderKey) before adding the header.
func Header(header string, value string) PrepareDecorator {
	return func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r.Header.Add(http.CanonicalHeaderKey(header), value)
			return p.Prepare(r)
		})
	}
}

// WithBearerAuthorization returns a PrepareDecorator that adds an HTTP Authorization header whose value is
// "Bearer " followed by the supplied token.
func WithBearerAuthorization(token string) PrepareDecorator {
	return Header(headerAuthorization, fmt.Sprintf("Bearer %s", token))
}

// WithJsonContentType returns a PrepareDecorator that adds an HTTP ContentType header whose value is
// "application/json".
func WithJsonContentType() PrepareDecorator {
	return Header(headerContentType, mimeTypeJson)
}

// WithURL returns a PrepareDecorator that populates the http.Request with a url.URL constructed from the
// supplied baseUrl and path.
func WithURL(baseUrl string, path string) PrepareDecorator {
	return func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			u, err := url.Parse(baseUrl)
			if err != nil {
				return r, err
			}
			u.Path = strings.TrimRight(u.Path, "/")
			if strings.HasPrefix(path, "/") {
				u.Path = path
			} else {
				u.Path += "/" + path
			}
			r.URL = u
			return p.Prepare(r)
		})
	}
}

// WithPathParameters returns a PrepareDecorator that replaces brace-enclosed keys within the request path
// (i.e., http.Request.URL.Path) with the corresponding values from the passed map.
func WithPathParameters(pathParameters map[string]interface{}) PrepareDecorator {
	return func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			for key, value := range ensureValueStrings(pathParameters) {
				r.URL.Path = strings.Replace(r.URL.Path, "{"+key+"}", value, -1)
			}
			return p.Prepare(r)
		})
	}
}

// WithQueryParameters returns a PrepareDecorators that encodes and applies the query parameters given
// in the supplied map (i.e., key=value).
func WithQueryParameters(queryParameters map[string]interface{}) PrepareDecorator {
	return func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			v := r.URL.Query()
			for key, value := range ensureValueStrings(queryParameters) {
				v.Add(key, value)
			}
			r.URL.RawQuery = v.Encode()
			return p.Prepare(r)
		})
	}
}
