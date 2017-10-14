// +build !go1.8

package autorest

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// RetriableRequest provides facilities for retrying an HTTP request.
type RetriableRequest struct {
	req *http.Request
	br  *bytes.Reader
}

// Prepare signals that the request is about to be sent.
func (rr *RetriableRequest) Prepare() (err error) {
	// preserve the request body; this is to support retry logic as
	// the underlying transport will always close the reqeust body
	if rr.req.Body != nil {
		if rr.br != nil {
			_, err = rr.br.Seek(0, 0 /*io.SeekStart*/)
			rr.req.Body = ioutil.NopCloser(rr.br)
		}
		if err != nil {
			return err
		}
		if rr.br == nil {
			// fall back to making a copy (only do this once)
			err = rr.prepareFromByteReader()
		}
	}
	return err
}

func removeRequestBody(req *http.Request) {
	req.Body = nil
	req.ContentLength = 0
}
