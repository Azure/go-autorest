package azure

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/Azure/go-autorest/autorest"
)

var (
	statusCodesForRetry = append(autorest.StatusCodesForRetry, http.StatusConflict)
)

// Client is the Azure specific implementation for autorest.Sender
type Client struct {
	autorest.Client
}

// Do implements the Sender interface by invoking the active Sender after applying authorization.
// If Sender is not set, it uses a new instance of http.Client. In both cases it will, if UserAgent
// is set, apply set the User-Agent header. This is the Azure specific implementation
func (azc Client) Do(r *http.Request) (*http.Response, error) {
	if r.UserAgent() == "" {
		r, _ = autorest.Prepare(r,
			autorest.WithUserAgent(azc.UserAgent))
	}
	r, err := autorest.Prepare(r,
		azc.WithInspection(),
		azc.WithAuthorization())
	if err != nil {
		return nil, NewErrorWithError(err, "autorest/Client", "Do", nil, "Preparing request failed")
	}
	resp, err := autorest.SendWithSender(azc.sender(), r,
		DoRetryForStatusCodes(azc, statusCodesForRetry...))
	autorest.Respond(resp,
		azc.ByInspecting())
	return resp, err
}

// sender returns the Sender to which to send requests.
func (azc Client) sender() autorest.Sender {
	if azc.Sender == nil {
		j, _ := cookiejar.New(nil)
		return &http.Client{Jar: j}
	}
	return azc.Sender
}
