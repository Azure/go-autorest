package azure

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/mocks"
)

const (
	headerAuthorization = "Authorization"
	longDelay           = 5 * time.Second
	retryDelay          = 10 * time.Millisecond
	testLogPrefix       = "azure:"
)

// Use a Client Inspector to set the request identifier.
func ExampleWithClientID() {
	uuid := "71FDB9F4-5E49-4C12-B266-DE7B4FD999A6"
	req, _ := autorest.Prepare(&http.Request{},
		autorest.AsGet(),
		autorest.WithBaseURL("https://microsoft.com/a/b/c/"))

	c := autorest.Client{Sender: mocks.NewSender()}
	c.RequestInspector = WithReturningClientID(uuid)

	c.Send(req)
	fmt.Printf("Inspector added the %s header with the value %s\n",
		HeaderClientID, req.Header.Get(HeaderClientID))
	fmt.Printf("Inspector added the %s header with the value %s\n",
		HeaderReturnClientID, req.Header.Get(HeaderReturnClientID))
	// Output:
	// Inspector added the x-ms-client-request-id header with the value 71FDB9F4-5E49-4C12-B266-DE7B4FD999A6
	// Inspector added the x-ms-return-client-request-id header with the value true
}

func TestWithReturningClientIDReturnsError(t *testing.T) {
	var errIn error
	uuid := "71FDB9F4-5E49-4C12-B266-DE7B4FD999A6"
	_, errOut := autorest.Prepare(&http.Request{},
		withErrorPrepareDecorator(&errIn),
		WithReturningClientID(uuid))

	if errOut == nil || errIn != errOut {
		t.Errorf("azure: WithReturningClientID failed to exit early when receiving an error -- expected (%v), received (%v)",
			errIn, errOut)
	}
}

func TestWithClientID(t *testing.T) {
	uuid := "71FDB9F4-5E49-4C12-B266-DE7B4FD999A6"
	req, _ := autorest.Prepare(&http.Request{},
		WithClientID(uuid))

	if req.Header.Get(HeaderClientID) != uuid {
		t.Errorf("azure: WithClientID failed to set %s -- expected %s, received %s",
			HeaderClientID, uuid, req.Header.Get(HeaderClientID))
	}
}

func TestWithReturnClientID(t *testing.T) {
	b := false
	req, _ := autorest.Prepare(&http.Request{},
		WithReturnClientID(b))

	if req.Header.Get(HeaderReturnClientID) != strconv.FormatBool(b) {
		t.Errorf("azure: WithReturnClientID failed to set %s -- expected %s, received %s",
			HeaderClientID, strconv.FormatBool(b), req.Header.Get(HeaderClientID))
	}
}

func TestExtractClientID(t *testing.T) {
	uuid := "71FDB9F4-5E49-4C12-B266-DE7B4FD999A6"
	resp := mocks.NewResponse()
	mocks.SetResponseHeader(resp, HeaderClientID, uuid)

	if ExtractClientID(resp) != uuid {
		t.Errorf("azure: ExtractClientID failed to extract the %s -- expected %s, received %s",
			HeaderClientID, uuid, ExtractClientID(resp))
	}
}

func TestExtractRequestID(t *testing.T) {
	uuid := "71FDB9F4-5E49-4C12-B266-DE7B4FD999A6"
	resp := mocks.NewResponse()
	mocks.SetResponseHeader(resp, HeaderRequestID, uuid)

	if ExtractRequestID(resp) != uuid {
		t.Errorf("azure: ExtractRequestID failed to extract the %s -- expected %s, received %s",
			HeaderRequestID, uuid, ExtractRequestID(resp))
	}
}

func TestIsAzureError_ReturnsTrueForAzureError(t *testing.T) {
	if !IsAzureError(&RequestError{}) {
		t.Errorf("azure: IsAzureError failed to return true for an Azure Service error")
	}
}

func TestIsAzureError_ReturnsFalseForNonAzureError(t *testing.T) {
	if IsAzureError(fmt.Errorf("An Error")) {
		t.Errorf("azure: IsAzureError return true for an non-Azure Service error")
	}
}

func TestNewErrorWithError_UsesReponseStatusCode(t *testing.T) {
	e := NewErrorWithError(fmt.Errorf("Error"), "packageType", "method", mocks.NewResponseWithStatus("Forbidden", http.StatusForbidden), "message")
	if e.StatusCode != http.StatusForbidden {
		t.Errorf("azure: NewErrorWithError failed to use the Status Code of the passed Response -- expected %v, received %v", http.StatusForbidden, e.StatusCode)
	}
}

func TestNewErrorWithError_ReturnsUnwrappedError(t *testing.T) {
	e1 := RequestError{}
	e1.ServiceError = &ServiceError{Code: "42", Message: "A Message"}
	e1.StatusCode = 200
	e1.RequestID = "A RequestID"
	e2 := NewErrorWithError(&e1, "packageType", "method", nil, "message")

	if !reflect.DeepEqual(e1, e2) {
		t.Errorf("azure: NewErrorWithError wrapped an RequestError -- expected %T, received %T", e1, e2)
	}
}

func TestNewErrorWithError_WrapsAnError(t *testing.T) {
	e1 := fmt.Errorf("Inner Error")
	var e2 interface{} = NewErrorWithError(e1, "packageType", "method", nil, "message")

	if _, ok := e2.(RequestError); !ok {
		t.Errorf("azure: NewErrorWithError failed to wrap a standard error -- received %T", e2)
	}
}

func TestWithErrorUnlessStatusCode_NotAnAzureError(t *testing.T) {
	body := `<html>
		<head>
			<title>IIS Error page</title>
		</head>
		<body>Some non-JSON error page</body>
	</html>`
	r := mocks.NewResponseWithContent(body)
	r.Request = mocks.NewRequest()
	r.StatusCode = http.StatusBadRequest
	r.Status = http.StatusText(r.StatusCode)

	err := autorest.Respond(r,
		WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	ok, _ := err.(*RequestError)
	if ok != nil {
		t.Fatalf("azure: azure.RequestError returned from malformed response: %v", err)
	}

	// the error body should still be there
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != body {
		t.Fatalf("response body is wrong. got=%q exptected=%q", string(b), body)
	}
}

func TestWithErrorUnlessStatusCode_FoundAzureError(t *testing.T) {
	j := `{
		"error": {
			"code": "InternalError",
			"message": "Azure is having trouble right now."
		}
	}`
	uuid := "71FDB9F4-5E49-4C12-B266-DE7B4FD999A6"
	r := mocks.NewResponseWithContent(j)
	mocks.SetResponseHeader(r, HeaderRequestID, uuid)
	r.Request = mocks.NewRequest()
	r.StatusCode = http.StatusInternalServerError
	r.Status = http.StatusText(r.StatusCode)

	err := autorest.Respond(r,
		WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())

	if err == nil {
		t.Fatalf("azure: returned nil error for proper error response")
	}
	azErr, ok := err.(*RequestError)
	if !ok {
		t.Fatalf("azure: returned error is not azure.RequestError: %T", err)
	}
	if expected := "InternalError"; azErr.ServiceError.Code != expected {
		t.Fatalf("azure: wrong error code. expected=%q; got=%q", expected, azErr.ServiceError.Code)
	}
	if azErr.ServiceError.Message == "" {
		t.Fatalf("azure: error message is not unmarshaled properly")
	}
	if expected := http.StatusInternalServerError; azErr.StatusCode != expected {
		t.Fatalf("azure: got wrong StatusCode=%d Expected=%d", azErr.StatusCode, expected)
	}
	if expected := uuid; azErr.RequestID != expected {
		t.Fatalf("azure: wrong request ID in error. expected=%q; got=%q", expected, azErr.RequestID)
	}

	_ = azErr.Error()

	// the error body should still be there
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Error(err)
	}
	if string(b) != j {
		t.Fatalf("response body is wrong. got=%q expected=%q", string(b), j)
	}

}

func TestGetAsyncOperation_ReturnsAzureAsyncOperationHeader(t *testing.T) {
	r := newLongRunningResponse()

	if GetAsyncOperation(r) != mocks.TestURL {
		t.Errorf("azure: GetAsyncOperation failed to extract the Azure-AsyncOperation header -- expected %v, received %v", mocks.TestURL, GetAsyncOperation(r))
	}
}

func TestGetAsyncOperation_ReturnsEmptyStringIfHeaderIsAbsent(t *testing.T) {
	r := mocks.NewResponse()

	if len(GetAsyncOperation(r)) != 0 {
		t.Errorf("azure: GetAsyncOperation failed to return empty string when the Azure-AsyncOperation header is absent -- received %v", GetAsyncOperation(r))
	}
}

func TestReponseIsLongRunning_ReturnsTrueForLongRunningResponse(t *testing.T) {
	r := newLongRunningResponse()

	if !ResponseIsLongRunning(r) {
		t.Errorf("azure: ResponseIsLongRunning returned false for a long-running response")
	}
}

func TestResponseIsLongRunning_ReturnsFalseIfStatusCodeIsNotCreated(t *testing.T) {
	r := newLongRunningResponse()
	r.StatusCode = http.StatusOK

	if ResponseIsLongRunning(r) {
		t.Errorf("azure: ResponseIsLongRunning returned true for a response without a 201 status code")
	}
}

func TestResponseIsLongRunning_ReturnsFalseIfAsyncHeaderIsAbsent(t *testing.T) {
	r := newLongRunningResponse()
	r.Header.Del(HeaderAsyncOperation)

	if ResponseIsLongRunning(r) {
		t.Errorf("azure: ResponseIsLongRunning returned true for a response without an Azure-AsyncOperation header")
	}
}

func TestNewAsyncPollingRequest_LeavesBodyOpenWhenAzureAsyncOperationHeaderIsMissing(t *testing.T) {
	r := newLongRunningResponse()
	r.Header.Del(HeaderAsyncOperation)

	NewAsyncPollingRequest(r, autorest.Client{Authorizer: autorest.NullAuthorizer{}})
	if !r.Body.(*mocks.Body).IsOpen() {
		t.Error("azure: NewAsyncPollingRequest closed the http.Request Body when the Azure-AsyncOperation header was missing")
	}
}

func TestNewAsyncPollingRequest_DoesNotReturnARequestWhenAzureAsyncOperationHeaderIsMissing(t *testing.T) {
	r := newLongRunningResponse()
	r.Header.Del(HeaderAsyncOperation)

	req, _ := NewAsyncPollingRequest(r, autorest.Client{Authorizer: autorest.NullAuthorizer{}})
	if req != nil {
		t.Error("azure: NewAsyncPollingRequest returned an http.Request when the Azure-AsyncOperation header was missing")
	}
}

func TestNewAsyncPollingRequest_ReturnsAnErrorWhenPrepareFails(t *testing.T) {
	r := newLongRunningResponse()

	_, err := NewAsyncPollingRequest(r, autorest.Client{Authorizer: mockFailingAuthorizer{}})
	if err == nil {
		t.Error("azure: NewAsyncPollingRequest failed to return an error when Prepare fails")
	}
}

func TestNewAsyncPollingRequest_LeavesBodyOpenWhenPrepareFails(t *testing.T) {
	r := newLongRunningResponse()
	r.Header.Set(http.CanonicalHeaderKey(HeaderAsyncOperation), "")

	_, err := NewAsyncPollingRequest(r, autorest.Client{Authorizer: autorest.NullAuthorizer{}})
	if !r.Body.(*mocks.Body).IsOpen() {
		t.Errorf("azure: NewAsyncPollingRequest closed the http.Request Body when Prepare returned an error (%v)", err)
	}
}

func TestNewAsyncPollingRequest_DoesNotReturnARequestWhenPrepareFails(t *testing.T) {
	r := newLongRunningResponse()
	r.Header.Set(http.CanonicalHeaderKey(HeaderAsyncOperation), mocks.TestBadURL)

	req, _ := NewAsyncPollingRequest(r, autorest.Client{Authorizer: autorest.NullAuthorizer{}})
	if req != nil {
		t.Error("azure: NewAsyncPollingRequest returned an http.Request when Prepare failed")
	}
}

func TestNewAsyncPollingRequest_ClosesTheResponseBody(t *testing.T) {
	r := newLongRunningResponse()

	NewAsyncPollingRequest(r, autorest.Client{Authorizer: autorest.NullAuthorizer{}})
	if r.Body.(*mocks.Body).IsOpen() {
		t.Error("azure: NewAsyncPollingRequest failed to close the response body when creating a new request")
	}
}

func TestNewAsyncPollingRequest_ReturnsAGetRequest(t *testing.T) {
	r := newLongRunningResponse()

	req, _ := NewAsyncPollingRequest(r, autorest.Client{Authorizer: autorest.NullAuthorizer{}})
	if req.Method != "GET" {
		t.Errorf("azure: NewAsyncPollingRequest did not create an HTTP GET request -- actual method %v", req.Method)
	}
}

func TestNewAsyncPollingRequest_ProvidesTheURL(t *testing.T) {
	r := newLongRunningResponse()

	req, _ := NewAsyncPollingRequest(r, autorest.Client{Authorizer: autorest.NullAuthorizer{}})
	if req.URL.String() != mocks.TestURL {
		t.Errorf("azure: NewAsyncPollingRequest did not create an HTTP with the expected URL -- received %v, expected %v", req.URL, mocks.TestURL)
	}
}

func TestNewAsyncPollingRequest_AppliesAuthorization(t *testing.T) {
	r := newLongRunningResponse()

	req, _ := NewAsyncPollingRequest(r, autorest.Client{Authorizer: mockAuthorizer{}})
	if req.Header.Get(http.CanonicalHeaderKey(headerAuthorization)) != mocks.TestAuthorizationHeader {
		t.Errorf("azure: NewAsyncPollingRequest did not apply authorization -- received %v, expected %v",
			req.Header.Get(http.CanonicalHeaderKey(headerAuthorization)), mocks.TestAuthorizationHeader)
	}
}

func TestWithAsyncPolling_Waits(t *testing.T) {
	resp := newLongRunningResponse()

	client := mocks.NewSender()
	client.SetResponse(resp)

	req, _ := NewAsyncPollingRequest(resp, autorest.DefaultClient)

	tt := time.Now()
	r, _ := autorest.SendWithSender(client, req,
		withAsyncResponseDecorator(2),
		WithAsyncPolling(retryDelay))
	s := time.Since(tt)
	if s < retryDelay {
		t.Error("azure: WithAsyncPolling failed to wait for at least the specified duration")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestWithAsyncPolling_Polls(t *testing.T) {
	resp := newLongRunningResponse()

	client := mocks.NewSender()
	client.SetResponse(resp)

	req, _ := NewAsyncPollingRequest(resp, autorest.DefaultClient)

	r, _ := autorest.SendWithSender(client, req,
		withAsyncResponseDecorator(2),
		WithAsyncPolling(retryDelay))
	if client.Attempts() != 3 {
		t.Errorf("azure: WithAsyncPolling failed to poll -- expected %v, actual %v", 2, client.Attempts())
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestWithAsyncPolling_DoesNotPollForNormalResponses(t *testing.T) {
	client := mocks.NewSender()

	r, _ := autorest.SendWithSender(client, mocks.NewRequest(),
		WithAsyncPolling(retryDelay))
	if client.Attempts() != 1 {
		t.Error("azure: WithAsyncPolling polled for a non-long-running Response")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestWithAsyncPolling_CancelsWhenSignaled(t *testing.T) {
	resp := newLongRunningResponse()
	mocks.SetRetryHeader(resp, longDelay)

	client := mocks.NewSender()
	client.SetResponse(resp)

	req, _ := NewAsyncPollingRequest(resp, autorest.DefaultClient)
	cancel := make(chan struct{})
	req.Cancel = cancel

	tt := time.Now()
	var r *http.Response
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()
		r, _ = autorest.SendWithSender(client, req,
			WithAsyncPolling(longDelay))
	}()
	wg.Wait()
	close(cancel)
	time.Sleep(10 * time.Millisecond)
	s := time.Since(tt)
	if s >= longDelay {
		t.Errorf("azure: WithAsyncPolling failed to cancel polling")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func withErrorPrepareDecorator(e *error) autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			*e = fmt.Errorf("azure: Faux Prepare Error")
			return r, *e
		})
	}
}

func withAsyncResponseDecorator(n int) autorest.SendDecorator {
	i := 0
	return func(s autorest.Sender) autorest.Sender {
		return autorest.SenderFunc(func(r *http.Request) (*http.Response, error) {
			resp, err := s.Do(r)
			if err == nil {
				if i < n {
					resp.StatusCode = http.StatusCreated
					resp.Header = http.Header{}
					resp.Header.Add(http.CanonicalHeaderKey(HeaderAsyncOperation), mocks.TestURL)
					i++
				} else {
					resp.StatusCode = http.StatusOK
					resp.Header.Del(http.CanonicalHeaderKey(HeaderAsyncOperation))
				}
			}
			return resp, err
		})
	}
}

func newLongRunningResponse() *http.Response {
	r := mocks.NewResponse()
	r.Header = http.Header{}
	mocks.SetRetryHeader(r, retryDelay)
	r.Header.Add(http.CanonicalHeaderKey(HeaderAsyncOperation), mocks.TestURL)
	r.StatusCode = http.StatusCreated
	return r
}

type mockAuthorizer struct{}

func (ma mockAuthorizer) WithAuthorization() autorest.PrepareDecorator {
	return autorest.WithHeader(headerAuthorization, mocks.TestAuthorizationHeader)
}

type mockFailingAuthorizer struct{}

func (mfa mockFailingAuthorizer) WithAuthorization() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			return r, fmt.Errorf("ERROR: mockFailingAuthorizer returned expected error")
		})
	}
}
