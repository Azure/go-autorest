package azure

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/mocks"
)

func TestGetAsyncOperation_ReturnsAzureAsyncOperationHeader(t *testing.T) {
	r := newAsynchronousResponse()

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

func TestIsAsynchronousResponse_ReturnsTrueForLongRunningResponse(t *testing.T) {
	r := newAsynchronousResponse()

	if !IsAsynchronousResponse(r) {
		t.Errorf("azure: IsAsynchronousResponse returned false for a long-running response")
	}
}

func TestIsAsynchronousResponse_ReturnsFalseIfStatusCodeIsNotCreated(t *testing.T) {
	r := newAsynchronousResponse()
	r.StatusCode = http.StatusOK

	if IsAsynchronousResponse(r) {
		t.Errorf("azure: IsAsynchronousResponse returned true for a response without a 201 status code")
	}
}

func TestIsAsynchronousResponse_ReturnsFalseIfAsyncHeaderIsAbsent(t *testing.T) {
	r := newAsynchronousResponse()
	r.Header.Del(HeaderAsyncOperation)

	if IsAsynchronousResponse(r) {
		t.Errorf("azure: IsAsynchronousResponse returned true for a response without an Azure-AsyncOperation header")
	}
}

func TestNewOperationResourceRequest_LeavesBodyOpen(t *testing.T) {
	r := newAsynchronousResponse()

	NewOperationResourceRequest(r, nil)
	if !r.Body.(*mocks.Body).IsOpen() {
		t.Error("azure: NewOperationResourceRequest closed the http.Request Body when the Azure-AsyncOperation header was missing")
	}
}

func TestNewOperationResourceRequest_DoesNotReturnARequestWhenAzureAsyncOperationHeaderIsMissing(t *testing.T) {
	r := newAsynchronousResponse()
	r.Header.Del(HeaderAsyncOperation)

	req, _ := NewOperationResourceRequest(r, nil)
	if req != nil {
		t.Error("azure: NewOperationResourceRequest returned an http.Request when the Azure-AsyncOperation header was missing")
	}
}

func TestNewOperationResourceRequest_ReturnsAnErrorWhenPrepareFails(t *testing.T) {
	r := newAsynchronousResponse()
	r.Header.Set(http.CanonicalHeaderKey(HeaderAsyncOperation), mocks.TestBadURL)

	_, err := NewOperationResourceRequest(r, nil)
	if err == nil {
		t.Error("azure: NewOperationResourceRequest failed to return an error when Prepare fails")
	}
}

func TestNewOperationResourceRequest_DoesNotReturnARequestWhenPrepareFails(t *testing.T) {
	r := newAsynchronousResponse()
	r.Header.Set(http.CanonicalHeaderKey(HeaderAsyncOperation), mocks.TestBadURL)

	req, _ := NewOperationResourceRequest(r, nil)
	if req != nil {
		t.Error("azure: NewOperationResourceRequest returned an http.Request when Prepare failed")
	}
}

func TestNewOperationResourceRequest_ReturnsAGetRequest(t *testing.T) {
	r := newAsynchronousResponse()

	req, _ := NewOperationResourceRequest(r, nil)
	if req.Method != "GET" {
		t.Errorf("azure: NewOperationResourceRequest did not create an HTTP GET request -- actual method %v", req.Method)
	}
}

func TestNewOperationResourceRequest_ProvidesTheURL(t *testing.T) {
	r := newAsynchronousResponse()

	req, _ := NewOperationResourceRequest(r, nil)
	if req.URL.String() != mocks.TestURL {
		t.Errorf("azure: NewOperationResourceRequest did not create an HTTP with the expected URL -- received %v, expected %v", req.URL, mocks.TestURL)
	}
}

func TestOperationError_ErrorReturnsAString(t *testing.T) {
	s := (OperationError{Code: "server code", Message: "server error"}).Error()
	if s == "" {
		t.Errorf("azure: OperationError#Error failed to return an error")
	}
	if !strings.Contains(s, "server code") || !strings.Contains(s, "server error") {
		t.Errorf("azure: OperationError#ToError returned a malformed error -- error(%v)", s)
	}
}

func TestOperationResource_HasSucceededReturnsFalseIfCanceled(t *testing.T) {
	if (OperationResource{Status: OperationCanceled}).HasSucceeded() {
		t.Errorf("azure: OperationResource#HasSucceeded failed to return false for a canceled operation")
	}
}

func TestOperationResource_HasSucceededReturnsFalseIfFailed(t *testing.T) {
	if (OperationResource{Status: OperationFailed}).HasSucceeded() {
		t.Errorf("azure: OperationResource#HasSucceeded failed to return false for a failed operation")
	}
}

func TestOperationResource_HasSucceededReturnsTrueIfSuccessful(t *testing.T) {
	if !(OperationResource{Status: OperationSucceeded}).HasSucceeded() {
		t.Errorf("azure: OperationResource#HasSucceeded failed to return true for a successful operation")
	}
}

func TestOperationResource_HasTerminatedReturnsTrueIfCanceled(t *testing.T) {
	if !(OperationResource{Status: OperationCanceled}).HasTerminated() {
		t.Errorf("azure: OperationResource#HasTerminated failed to return true for a canceled operation")
	}
}

func TestOperationResource_HasTerminatedReturnsTrueIfFailed(t *testing.T) {
	if !(OperationResource{Status: OperationFailed}).HasTerminated() {
		t.Errorf("azure: OperationResource#HasTerminated failed to return true for a failed operation")
	}
}

func TestOperationResource_HasTerminatedReturnsTrueIfSuccessful(t *testing.T) {
	if !(OperationResource{Status: OperationSucceeded}).HasTerminated() {
		t.Errorf("azure: OperationResource#HasTerminated failed to return true for a successful operation")
	}
}

func TestOperationResource_HasTerminatedReturnsFalseNotTerminated(t *testing.T) {
	if (OperationResource{Status: "UnknownStatus"}).HasTerminated() {
		t.Errorf("azure: OperationResource#HasTerminated returned true for a non-terminal operation")
	}
}

func TestOperationResource_GetErrorReturnsErrorIfCanceled(t *testing.T) {
	if (OperationResource{Status: OperationCanceled}).GetError() == nil {
		t.Errorf("azure: OperationResource#GetError failed to return an error for canceled operations")
	}
}

func TestOperationResource_GetErrorReturnsErrorIfFailed(t *testing.T) {
	if (OperationResource{Status: OperationFailed}).GetError() == nil {
		t.Errorf("azure: OperationResource#GetError failed to return an error for failed operations")
	}
}

func TestOperationResource_GetErrorDoesReturnUnlessFailedOrCanceled(t *testing.T) {
	if (OperationResource{Status: OperationSucceeded}).GetError() != nil {
		t.Errorf("azure: OperationResource#GetError return an error for operation that was not canceled or failed")
	}
}

func TestDoPollForAsynchronous_IgnoresUnspecifiedStatusCodes(t *testing.T) {
	client := mocks.NewSender()

	r, _ := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Duration(0), time.Duration(0)))

	if client.Attempts() != 1 {
		t.Errorf("azure: DoPollForAsynchronous polled for unspecified status code")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_PollsForSpecifiedStatusCodes(t *testing.T) {
	client := mocks.NewSender()
	client.AppendResponse(newAsynchronousResponse())

	r, _ := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond, time.Millisecond))

	if client.Attempts() != 2 {
		t.Errorf("azure: DoPollForAsynchronous failed to poll for specified status code")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_CanBeCanceled(t *testing.T) {
	cancel := make(chan struct{})
	delay := 5 * time.Second

	client := mocks.NewSender()
	client.AppendResponse(newAsynchronousResponse())
	client.AppendAndRepeatResponse(newOperationResourceResponse("Busy"), -1)

	var wg sync.WaitGroup
	wg.Add(1)
	start := time.Now()
	go func() {
		wg.Done()
		r, _ := autorest.SendWithSender(client, mocks.NewRequest(),
			DoPollForAsynchronous(time.Millisecond, time.Millisecond))
		autorest.Respond(r,
			autorest.ByClosing())
	}()
	wg.Wait()
	close(cancel)
	time.Sleep(5 * time.Millisecond)
	if time.Since(start) >= delay {
		t.Errorf("azure: DoPollForAsynchronous failed to cancel")
	}
}

func TestDoPollForAsynchronous_ClosesAllNonreturnedResponseBodiesWhenPolling(t *testing.T) {
	r1 := newAsynchronousResponse()
	r2 := newOperationResourceResponse("busy")
	r3 := newOperationResourceResponse(OperationSucceeded)

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)
	client.AppendResponse(r3)

	r, _ := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond, time.Millisecond))

	if r1.Body.(*mocks.Body).IsOpen() || r2.Body.(*mocks.Body).CloseAttempts() < 2 {
		t.Errorf("azure: DoPollForAsynchronous did not close unreturned response bodies")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_LeavesLastResponseBodyOpen(t *testing.T) {
	r1 := newAsynchronousResponse()
	r2 := newOperationResourceResponse("busy")
	r3 := newOperationResourceResponse(OperationSucceeded)

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)
	client.AppendResponse(r3)

	r, _ := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond, time.Millisecond))

	if !r.Body.(*mocks.Body).IsOpen() {
		t.Errorf("azure: DoPollForAsynchronous did not leave open the body of the last response")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_DoesNotPollIfCreatingOperationRequestFails(t *testing.T) {
	r1 := newAsynchronousResponse()
	mocks.SetResponseHeader(r1, http.CanonicalHeaderKey(HeaderAsyncOperation), mocks.TestBadURL)
	r2 := newOperationResourceResponse("busy")

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)

	r, _ := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond, time.Millisecond))

	if client.Attempts() > 1 {
		t.Errorf("azure: DoPollForAsynchronous polled with an invalidly formed operation request")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_StopsPollingAfterAnError(t *testing.T) {
	r1 := newAsynchronousResponse()
	r2 := newOperationResourceResponse("busy")

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)
	client.SetError(fmt.Errorf("Faux Error"))
	client.SetEmitErrorAfter(2)

	r, _ := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond, time.Millisecond))

	if client.Attempts() > 3 {
		t.Errorf("azure: DoPollForAsynchronous failed to stop polling after receiving an error")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_ReturnsPollingError(t *testing.T) {
	client := mocks.NewSender()
	client.AppendAndRepeatResponse(newAsynchronousResponse(), 5)
	client.SetError(fmt.Errorf("Faux Error"))
	client.SetEmitErrorAfter(1)

	r, err := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond, time.Millisecond))

	if err == nil {
		t.Errorf("azure: DoPollForAsynchronous failed to return error from polling")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_PollsUntilOperationResourceHasTerminated(t *testing.T) {
	r1 := newAsynchronousResponse()
	r2 := newOperationResourceResponse("busy")
	r3 := newOperationResourceResponse(OperationCanceled)

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)
	client.AppendAndRepeatResponse(r3, 1)

	r, _ := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond, time.Millisecond))

	if client.Attempts() < 4 {
		t.Errorf("azure: DoPollForAsynchronous stopped polling before receiving a terminated OperationResource")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_StopsPollingWhenOperationResourceHasTerminated(t *testing.T) {
	r1 := newAsynchronousResponse()
	r2 := newOperationResourceResponse("busy")
	r3 := newOperationResourceResponse(OperationCanceled)

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)
	client.AppendAndRepeatResponse(r3, 2)

	r, _ := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond, time.Millisecond))

	if client.Attempts() > 4 {
		t.Errorf("azure: DoPollForAsynchronous failed to stop after receiving a terminated OperationResource")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_ReturnsAnErrorForCanceledOperations(t *testing.T) {
	r1 := newAsynchronousResponse()
	r2 := newOperationResourceResponse("busy")
	r3 := newOperationResourceErrorResponse(OperationCanceled)

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)
	client.AppendAndRepeatResponse(r3, 1)

	r, err := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond, time.Millisecond))

	if err == nil || !strings.Contains(fmt.Sprintf("%v", err), "BadArgument") {
		t.Errorf("azure: DoPollForAsynchronous failed to return an appropriate error for a canceled OperationResource")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_ReturnsAnErrorForFailedOperations(t *testing.T) {
	r1 := newAsynchronousResponse()
	r2 := newOperationResourceResponse("busy")
	r3 := newOperationResourceErrorResponse(OperationFailed)

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)
	client.AppendAndRepeatResponse(r3, 1)

	r, err := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond, time.Millisecond))

	if err == nil || !strings.Contains(fmt.Sprintf("%v", err), "BadArgument") {
		t.Errorf("azure: DoPollForAsynchronous failed to return an appropriate error for a canceled OperationResource")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_ReturnsNoErrorForSuccessfulOperations(t *testing.T) {
	r1 := newAsynchronousResponse()
	r2 := newOperationResourceResponse("busy")
	r3 := newOperationResourceErrorResponse(OperationSucceeded)

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)
	client.AppendAndRepeatResponse(r3, 1)

	r, err := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond, time.Millisecond))

	if err != nil {
		t.Errorf("azure: DoPollForAsynchronous returned an error for a successful OperationResource")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_StopsPollingIfItReceivesAnInvalidOperationResource(t *testing.T) {
	r1 := newAsynchronousResponse()
	r2 := newOperationResourceResponse("busy")
	r3 := mocks.NewResponseWithContent(operationResourceIllegal)
	r4 := newOperationResourceResponse(OperationSucceeded)

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)
	client.AppendAndRepeatResponse(r3, 1)
	client.AppendAndRepeatResponse(r4, 1)

	r, err := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond, time.Millisecond))

	if client.Attempts() > 4 {
		t.Errorf("azure: DoPollForAsynchronous failed to polling after receiving an invalid OperationResource")
	}
	if err == nil {
		t.Errorf("azure: DoPollForAsynchronous failed to return an error after receving an invalid OperationResource")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func newAsynchronousResponse() *http.Response {
	r := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetResponseHeader(r, http.CanonicalHeaderKey(HeaderAsyncOperation), mocks.TestURL)
	mocks.SetRetryHeader(r, retryDelay)
	return r
}

const (
	operationResourceIllegal = `
	This is not JSON and should fail...badly.
`

	operationResourceFormat = `
	{
		"id": "/subscriptions/id/locations/westus/operationsStatus/sameguid",
		"name": "sameguid",
		"status" : "%s",
		"startTime" : "2006-01-02T15:04:05Z",
		"endTime" : "2006-01-02T16:04:05Z",
		"percentComplete" : 50.00,

		"properties" : {}
	}
`

	operationResourceErrorFormat = `
	{
		"id": "/subscriptions/id/locations/westus/operationsStatus/sameguid",
		"name": "sameguid",
		"status" : "%s",
		"startTime" : "2006-01-02T15:04:05Z",
		"endTime" : "2006-01-02T16:04:05Z",
		"percentComplete" : 50.00,

		"properties" : {},
		"error" : {
			"code" : "BadArgument",
			"message" : "The provided database 'foo' has an invalid username."
		}
	}
`
)

func newOperationResourceResponse(status string) *http.Response {
	return mocks.NewResponseWithContent(fmt.Sprintf(operationResourceFormat, status))
}

func newOperationResourceErrorResponse(status string) *http.Response {
	return mocks.NewResponseWithContent(fmt.Sprintf(operationResourceErrorFormat, status))
}
