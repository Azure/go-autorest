package azure

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/mocks"
	"github.com/Azure/go-autorest/autorest/to"
)

const (
	TestResource = "SomeResource"
	TestClientID = "SomeClientID"
)

const MockDeviceCodeResponse = `
{
	"device_code": "10000-40-1234567890",
	"user_code": "ABCDEF",
	"verification_url": "http://aka.ms/deviceauth",
	"expires_in": "900",
	"interval": "0"
}
`

const MockDeviceTokenResponse = `{
	"access_token": "accessToken",
	"refresh_token": "refreshToken",
	"expires_in": "1000",
	"expires_on": "2000",
	"not_before": "3000",
	"resource": "resource",
	"token_type": "type"
}
`

func TestDeviceCodeIncludesResource(t *testing.T) {
	sender := mocks.NewSender()
	sender.EmitContent(MockDeviceCodeResponse)
	sender.EmitStatus("OK", 200)
	client := &autorest.Client{Sender: sender}

	code, err := InitiateDeviceAuth(client, TestClientID, TestResource)
	if err != nil {
		t.Errorf("unexpected error initiating device auth")
	}

	if *code.Resource != TestResource {
		t.Errorf("InitiateDeviceAuth failed to stash the resource in the DeviceCode struct")
	}
}

func TestDeviceCodeReturnsErrorIfSendingFails(t *testing.T) {
	sender := mocks.NewSender()
	sender.EmitErrors(1)
	sender.SetError(fmt.Errorf("this is an error"))
	client := &autorest.Client{Sender: sender}

	_, err := InitiateDeviceAuth(client, TestClientID, TestResource)
	if err == nil {
		t.Errorf("failed to get err") // (expecting Error occurred while requesting Device Authorization Code)
	}
}

func TestDeviceCodeReturnsErrorIfBadRequest(t *testing.T) {
	sender := mocks.NewSender()
	sender.EmitStatus("Bad Request", 400)
	client := &autorest.Client{Sender: sender}

	_, err := InitiateDeviceAuth(client, TestClientID, TestResource)
	if err == nil {
		t.Errorf("failed to get err") // (expecting Request for Device Authorization Code resulted in Status)
	}
}

func TestDeviceCodeReturnsErrorIfCannotDeserializeDeviceCode(t *testing.T) {
	gibberishJSON := strings.Replace(MockDeviceCodeResponse, "expires_in", "\":, :gibberish", -1)
	sender := mocks.NewSender()
	sender.EmitContent(gibberishJSON)
	client := &autorest.Client{Sender: sender}

	_, err := InitiateDeviceAuth(client, TestClientID, TestResource)
	if err == nil {
		t.Errorf("failed to get error") // (expecting Failed to decode Device Authorization Code json)
	}
}

func deviceCode() *DeviceCode {
	var deviceCode DeviceCode
	json.Unmarshal([]byte(MockDeviceCodeResponse), &deviceCode)
	deviceCode.Resource = to.StringPtr("resource")
	return &deviceCode
}

func TestDeviceTokenReturns(t *testing.T) {
	sender := mocks.NewSender()
	sender.EmitContent(MockDeviceTokenResponse)
	sender.EmitStatus("OK", 200)
	client := &autorest.Client{Sender: sender}

	_, err := WaitForUserCompletion(client, TestClientID, deviceCode())
	if err != nil {
		t.Errorf("got error unexpectedly")
	}
}

func TestDeviceTokenReturnsErrorIfSendingFails(t *testing.T) {
	sender := mocks.NewSender()
	sender.EmitErrors(1)
	sender.SetError(fmt.Errorf("this is an error"))
	client := &autorest.Client{Sender: sender}

	_, err := WaitForUserCompletion(client, TestClientID, deviceCode())
	if err == nil {
		t.Errorf("failed to get err") // (expecting Error occurred while exchanging device code for token)
	}
}

func TestDeviceTokenReturnsErrorIfServerError(t *testing.T) {
	sender := mocks.NewSender()
	sender.EmitStatus("Internal Server Error", 500)
	client := &autorest.Client{Sender: sender}

	_, err := WaitForUserCompletion(client, TestClientID, deviceCode())
	if err == nil {
		t.Errorf("failed to get error") // (expecting Exchanging Device Code for OAuth token resulted in Status)
	}
}

func TestDeviceTokenReturnsErrorIfCannotDeserializeDeviceToken(t *testing.T) {
	gibberishJSON := strings.Replace(MockDeviceTokenResponse, "expires_in", ";:\"gibberish", -1)
	sender := mocks.NewSender()
	sender.EmitContent(gibberishJSON)
	client := &autorest.Client{Sender: sender}

	_, err := WaitForUserCompletion(client, TestClientID, deviceCode())
	if err == nil {
		t.Errorf("failed to get err") // (expecting Error occurred while exchanging device code for token)
	}
}

func errorDeviceTokenResponse(message string) string {
	return `{ "error": "` + message + `" }`
}

func TestDeviceTokenReturnsErrorIfAuthorizationPending(t *testing.T) {
	sender := mocks.NewSender()
	sender.EmitContent(errorDeviceTokenResponse("authorization_pending"))
	sender.EmitStatus("Bad Request", 400)
	client := &autorest.Client{Sender: sender}

	go func() {
		time.Sleep(1 * time.Second)
		sender.EmitContent(errorDeviceTokenResponse("access_denied"))
	}()

	_, _ = WaitForUserCompletion(client, TestClientID, deviceCode())
}

func TestDeviceTokenReturnsErrorIfSlowDown(t *testing.T) {
	sender := mocks.NewSender()
	sender.EmitContent(errorDeviceTokenResponse("slow_down"))
	sender.EmitStatus("Bad Request", 400)
	client := &autorest.Client{Sender: sender}

	go func() {
		time.Sleep(1 * time.Second)
		sender.EmitContent(errorDeviceTokenResponse("access_denied"))
	}()

	_, _ = WaitForUserCompletion(client, TestClientID, deviceCode())
}

func TestDeviceTokenReturnsErrorIfAccessDenied(t *testing.T) {
	sender := mocks.NewSender()
	sender.EmitContent(errorDeviceTokenResponse("access_denied"))
	sender.EmitStatus("Bad Request", 400)
	client := &autorest.Client{Sender: sender}

	_, err := WaitForUserCompletion(client, TestClientID, deviceCode())
	if err == nil {
		t.Errorf("failed to get error")
	}
	if err != ErrDeviceAccessDenied {
		t.Errorf("got wrong error")
	}
}

func TestDeviceTokenReturnsErrorIfCodeExpired(t *testing.T) {
	sender := mocks.NewSender()
	sender.EmitContent(errorDeviceTokenResponse("code_expired"))
	sender.EmitStatus("Bad Request", 400)
	client := &autorest.Client{Sender: sender}

	_, err := WaitForUserCompletion(client, TestClientID, deviceCode())
	if err == nil {
		t.Errorf("failed to get error")
	}
	if err != ErrDeviceCodeExpired {
		t.Errorf("got wrong error")
	}
}

func TestDeviceTokenReturnsErrorForUnknownError(t *testing.T) {
	sender := mocks.NewSender()
	sender.EmitContent(errorDeviceTokenResponse("unknown_error"))
	sender.EmitStatus("Bad Request", 400)
	client := &autorest.Client{Sender: sender}

	_, err := WaitForUserCompletion(client, TestClientID, deviceCode())
	if err == nil {
		t.Errorf("failed to get error")
	}
	if err != ErrDeviceGeneric {
		fmt.Println(err)
		t.Errorf("got wrong error")
	}
}
