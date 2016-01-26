package azure

/*
  This file is largely based on rjw57/oauth2device's code, with the follow differences:
   * scope -> resource, and only allow a single one
   * receive "Message" in the DeviceCode struct and show it to users as the prompt
   * azure-xplat-cli has the following behavior that this emulates:
     - does not send client_secret during the token exchange
     - sends resource again in the token exchange request
*/

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Azure/go-autorest/autorest"
)

const (
	// OAuthDeviceEndpoint is Azure's OAuth2 Device Flow Endpoint
	OAuthDeviceEndpoint = "https://login.microsoftonline.com/common/oauth2/devicecode"
	// OAuthTokenEndpoint is Azure's OAuth2 Token Endpoint
	OAuthTokenEndpoint = "https://login.microsoftonline.com/common/oauth2/token"

	// authAPIVersionQueryParamName is the name
	authAPIVersionQueryParamName  = "api-version"
	authAPIVersionQueryParamValue = "1.0"
)

var (
	// ErrDeviceGeneric represents an unknown error from the token endpoint when using device flow
	ErrDeviceGeneric = fmt.Errorf("Error while retrieving OAuth token: Unknown Error")

	// ErrDeviceAccessDenied represents an access denied error from the token endpoint when using device flow
	ErrDeviceAccessDenied = fmt.Errorf("Error while retrieving OAuth token: Access Denied")

	// ErrDeviceAuthorizationPending represents the server waiting on the user to complete the device flow
	ErrDeviceAuthorizationPending = fmt.Errorf("Error while retrieving OAuth token: Authorization Pending")

	// ErrDeviceCodeExpired represents the server timing out and expiring the code during device flow
	ErrDeviceCodeExpired = fmt.Errorf("Error while retrieving OAuth token: Code Expired")

	// ErrDeviceSlowDown represents the service telling us we're polling too often during device flow
	ErrDeviceSlowDown = fmt.Errorf("Error while retrieving OAuth token: Slow Down")
)

// DeviceCode is the object returned by the device auth endpoint
// It contains information to instruct the user to complete the auth flow
type DeviceCode struct {
	DeviceCode      *string `json:"device_code"`
	UserCode        *string `json:"user_code"`
	VerificationURL *string `json:"verification_url"`
	ExpiresIn       *int64  `json:"expires_in,string"`
	Interval        *int64  `json:"interval,string"`

	Message  *string `json:"message"` // Azure specific
	Resource *string // store this as well, needed during token exchange for azure
}

// TokenError is the object returned by the token exchange endpoint
// when something is amiss
type TokenError struct {
	Error            *string `json:"error"`
	ErrorCodes       []int   `json:"error_codes"`
	ErrorDescription *string `json:"error_description"`
	Timestamp        *string `json:"timestamp"`
	TraceID          *string `json:"trace_id"`
}

// DeviceToken is the object return by the token exchange endpoint
// It can either look like a Token or an ErrorToken, so put both here
// and check for presence of "Error" to know if we are in error state
type DeviceToken struct {
	Token
	TokenError
}

// InitiateDeviceAuth initiates a device auth flow. It returns a DeviceCode
// that can be used with CheckForUserCompletion or WaitForUserCompletion.
func InitiateDeviceAuth(client *autorest.Client, clientID, resource string) (*DeviceCode, error) {
	req, _ := autorest.Prepare(
		&http.Request{},
		autorest.AsPost(),
		autorest.WithBaseURL(OAuthDeviceEndpoint),
		autorest.WithFormData(url.Values{
			"client_id": []string{clientID},
			"resource":  []string{resource},
		}),
		autorest.WithQueryParameters(map[string]interface{}{
			authAPIVersionQueryParamName: authAPIVersionQueryParamValue,
		}))

	resp, err := client.Send(req)
	if err != nil {
		return nil, fmt.Errorf("Error occurred while requesting Device Authorization Code. %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"Request for Device Authorization Code resulted in Status: %d: %s",
			resp.StatusCode,
			http.StatusText(resp.StatusCode))
	}

	var code DeviceCode
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&code)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode Device Authorization Code json")
	}

	code.Resource = &resource

	return &code, nil
}

// CheckForUserCompletion takes a DeviceCode and checks with the Azure AD OAuth endpoint
// to see if the device flow has: been completed, timed out, or otherwise failed
func CheckForUserCompletion(client *autorest.Client, clientID string, code *DeviceCode) (*Token, error) {
	req, _ := autorest.Prepare(
		&http.Request{},
		autorest.AsPost(),
		autorest.WithBaseURL(OAuthTokenEndpoint),
		autorest.WithFormData(url.Values{
			"client_id":  []string{clientID},
			"code":       []string{*code.DeviceCode},
			"grant_type": []string{OAuthGrantTypeDeviceCode},
			"resource":   []string{*code.Resource},
		}),
		autorest.WithQueryParameters(map[string]interface{}{
			authAPIVersionQueryParamName: authAPIVersionQueryParamValue,
		}))

	resp, err := client.Send(req)
	if err != nil {
		return nil, fmt.Errorf("Error occurred while exchanging device code for token. %s", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
		return nil, fmt.Errorf(
			"Exchanging Device Code for OAuth token resulted in Status: %d: %s",
			resp.StatusCode,
			http.StatusText(resp.StatusCode))
	}

	var token DeviceToken
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&token)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode Device Token json")
	}

	if token.Error == nil {
		return &token.Token, nil
	}

	switch *token.Error {
	case "authorization_pending":
		return nil, ErrDeviceAuthorizationPending
	case "slow_down":
		return nil, ErrDeviceSlowDown
	case "access_denied":
		return nil, ErrDeviceAccessDenied
	case "code_expired":
		return nil, ErrDeviceCodeExpired
	default:
		return nil, ErrDeviceGeneric
	}
}

// WaitForUserCompletion calls CheckForUserCompletion repeatedly until a token is granted or an error state occurs.
// This prevents the user from looping and checking against 'ErrDeviceAuthorizationPending'.
func WaitForUserCompletion(client *autorest.Client, clientID string, code *DeviceCode) (*Token, error) {
	intervalDuration := time.Duration(*code.Interval) * time.Second
	waitDuration := intervalDuration

	for {
		token, err := CheckForUserCompletion(client, clientID, code)

		if err == nil {
			return token, nil
		}

		switch err {
		case ErrDeviceSlowDown:
			waitDuration = waitDuration + (intervalDuration / 2)
		case ErrDeviceAuthorizationPending:
			// noop
		default: // everything else is "fatal" to us
			return nil, err
		}

		time.Sleep(waitDuration)
	}
}
