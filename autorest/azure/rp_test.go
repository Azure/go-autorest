package azure

import (
	"net/http"
	"testing"
	"time"

	"github.com/Azure/go-autorest/autorest"

	"github.com/Azure/go-autorest/autorest/mocks"
)

func TestRegisterResourceProvider(t *testing.T) {
	client := mocks.NewSender()
	client.AppendResponse(mocks.NewResponseWithStatus("MissingSubscriptionRegistration", http.StatusTooManyRequests))
	client.AppendResponse(mocks.NewResponseWithStatus("200 OK", http.StatusOK))

	r, err := autorest.SendWithSender(client, mocks.NewRequest(),
		RegisterResourceProvider(),
		autorest.DoRetryForStatusCodes(5, time.Duration(2*time.Second), autorest.StatusCodesForRetry...),
	)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	autorest.Respond(r,
		autorest.ByDiscardingBody(),
		autorest.ByClosing(),
	)

	if r.StatusCode != http.StatusOK {
		t.Fatalf("azure: Sender#RegisterResourceProvider -- Got: StatusCode %v; Want: StatusCode 200 OK", r.StatusCode)
	}
}
