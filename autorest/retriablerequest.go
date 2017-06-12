package autorest

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

// RetriableRequest provides facilities for retrying an HTTP request.
type RetriableRequest struct {
	req *http.Request
	rc  io.ReadCloser
	br  *bytes.Reader
}

// NewRetriableRequest returns a wrapper around an HTTP request that support retry logic.
func NewRetriableRequest(req *http.Request) *RetriableRequest {
	return &RetriableRequest{req: req}
}

// Request returns the wrapped HTTP request.
func (rr *RetriableRequest) Request() *http.Request {
	return rr.req
}

// Prepare signals that the request is about to be sent.
func (rr *RetriableRequest) Prepare() (err error) {
	// preserve the request body; this is to support retry logic as
	// the underlying transport will always close the reqeust body
	if rr.req.Body != nil {
		if rr.req.GetBody != nil {
			// this will allow us to preserve the body without having to
			// make a copy.  note we need to do this on each iteration
			rr.rc, err = rr.req.GetBody()
			if err != nil {
				return err
			}
		} else if rr.br == nil {
			// fall back to making a copy (only do this once)
			b := []byte{}
			if rr.req.ContentLength > 0 {
				b = make([]byte, rr.req.ContentLength)
				_, err = io.ReadFull(rr.req.Body, b)
				if err != nil {
					return err
				}
			} else {
				b, err = ioutil.ReadAll(rr.req.Body)
				if err != nil {
					return err
				}
			}
			rr.br = bytes.NewReader(b)
			rr.req.Body = ioutil.NopCloser(rr.br)
		}
	}
	return err
}

// Reset indicates that the request will be sent again.
func (rr *RetriableRequest) Reset() {
	if rr.rc != nil {
		rr.req.Body = rr.rc
	} else if rr.br != nil {
		rr.br.Seek(0, io.SeekStart)
	}
}
