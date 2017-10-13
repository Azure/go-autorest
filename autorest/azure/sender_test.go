package azure

import (
	"net/http"
	"testing"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/mocks"
)

func TestRegisterResourceProvider(t *testing.T) {
	client := mocks.NewSender()
	client.AppendResponse(mocks.NewResponseWithBodyAndStatus(mocks.NewBody(`{
	"error":{
		"code":"MissingSubscriptionRegistration",
		"message":"The subscription registration is in 'Unregistered' state. The subscription must be registered to use namespace 'Microsoft.EventGrid'. See https://aka.ms/rps-not-found for how to register subscriptions.",
		"details":[
			{
				"code":"MissingSubscriptionRegistration",
				"target":"Microsoft.EventGrid",
				"message":"The subscription registration is in 'Unregistered' state. The subscription must be registered to use namespace 'Microsoft.EventGrid'. See https://aka.ms/rps-not-found for how to register subscriptions."
			}
		]
	}
}
`), http.StatusConflict, "MissingSubscriptionRegistration"))

	client.AppendResponse(mocks.NewResponseWithBodyAndStatus(mocks.NewBody(`{
	"registrationState": "Registering"
}
`), http.StatusOK, "200 OK"))

	client.AppendResponse(mocks.NewResponseWithBodyAndStatus(mocks.NewBody(`{
	"registrationState": "Registered"
}
`), http.StatusOK, "200 OK"))

	client.AppendResponse(mocks.NewResponseWithStatus("200 OK", http.StatusOK))

	r, err := autorest.SendWithSender(client, mocks.NewRequestForURL("https://lol/subscriptions/rofl"),
		DoRetryForStatusCodes(Client{}, statusCodesForRetry...),
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
