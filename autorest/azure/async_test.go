package azure

import (
	"fmt"
	"io/ioutil"
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

	if getAsyncOperation(r) != mocks.TestAzureAsyncURL {
		t.Errorf("azure: getAsyncOperation failed to extract the Azure-AsyncOperation header -- expected %v, received %v", mocks.TestURL, getAsyncOperation(r))
	}
}

func TestGetAsyncOperation_ReturnsEmptyStringIfHeaderIsAbsent(t *testing.T) {
	r := mocks.NewResponse()

	if len(getAsyncOperation(r)) != 0 {
		t.Errorf("azure: getAsyncOperation failed to return empty string when the Azure-AsyncOperation header is absent -- received %v", getAsyncOperation(r))
	}
}

func TestHasSucceeded_ReturnsTrueForSuccess(t *testing.T) {
	if !hasSucceeded(operationSucceeded) {
		t.Error("azure: hasSucceeded failed to return true for success")
	}
}

func TestHasSucceeded_ReturnsFalseOtherwise(t *testing.T) {
	if hasSucceeded("not a success string") {
		t.Error("azure: hasSucceeded returned true for a non-success")
	}
}

func TestHasTerminated_ReturnsTrueForValidTerminationStates(t *testing.T) {
	for _, state := range []string{operationSucceeded, operationCanceled, operationFailed} {
		if !hasTerminated(state) {
			t.Errorf("azure: hasTerminated failed to return true for the '%s' state", state)
		}
	}
}

func TestHasTerminated_ReturnsFalseForUnknownStates(t *testing.T) {
	if hasTerminated("not a known state") {
		t.Error("azure: hasTerminated returned true for an unknown state")
	}
}

func TestOperationError_ErrorReturnsAString(t *testing.T) {
	s := (operationError{Code: "server code", Message: "server error"}).Error()
	if s == "" {
		t.Errorf("azure: operationError#Error failed to return an error")
	}
	if !strings.Contains(s, "server code") || !strings.Contains(s, "server error") {
		t.Errorf("azure: operationError#Error returned a malformed error -- error='%v'", s)
	}
}

func TestOperationResource_StateReturnsState(t *testing.T) {
	if (operationResource{Status: "state"}).state() != "state" {
		t.Errorf("azure: operationResource#state failed to return the correct state")
	}
}

func TestOperationResource_HasSucceededReturnsFalseIfNotSuccess(t *testing.T) {
	if (operationResource{Status: "not a success string"}).hasSucceeded() {
		t.Errorf("azure: operationResource#hasSucceeded failed to return false for a canceled operation")
	}
}

func TestOperationResource_HasSucceededReturnsTrueIfSuccessful(t *testing.T) {
	if !(operationResource{Status: operationSucceeded}).hasSucceeded() {
		t.Errorf("azure: operationResource#hasSucceeded failed to return true for a successful operation")
	}
}

func TestOperationResource_HasTerminatedReturnsTrueForKnownStates(t *testing.T) {
	for _, state := range []string{operationSucceeded, operationCanceled, operationFailed} {
		if !(operationResource{Status: state}).hasTerminated() {
			t.Errorf("azure: operationResource#hasTerminated failed to return true for the '%s' state", state)
		}
	}
}

func TestOperationResource_HasTerminatedReturnsFalseForUnknownStates(t *testing.T) {
	if (operationResource{Status: "not a known state"}).hasTerminated() {
		t.Errorf("azure: operationResource#hasTerminated returned true for a non-terminal operation")
	}
}

func TestProvisioningStatus_StateReturnsState(t *testing.T) {
	if (provisioningStatus{provisioningProperties{"state"}}).state() != "state" {
		t.Errorf("azure: provisioningStatus#state failed to return the correct state")
	}
}

func TestProvisioningStatus_HasSucceededReturnsFalseIfNotSuccess(t *testing.T) {
	if (provisioningStatus{provisioningProperties{"not a success string"}}).hasSucceeded() {
		t.Errorf("azure: provisioningStatus#hasSucceeded failed to return false for a canceled operation")
	}
}

func TestProvisioningStatus_HasSucceededReturnsTrueIfSuccessful(t *testing.T) {
	if !(provisioningStatus{provisioningProperties{operationSucceeded}}).hasSucceeded() {
		t.Errorf("azure: provisioningStatus#hasSucceeded failed to return true for a successful operation")
	}
}

func TestProvisioningStatus_HasTerminatedReturnsTrueForKnownStates(t *testing.T) {
	for _, state := range []string{operationSucceeded, operationCanceled, operationFailed} {
		if !(provisioningStatus{provisioningProperties{state}}).hasTerminated() {
			t.Errorf("azure: provisioningStatus#hasTerminated failed to return true for the '%s' state", state)
		}
	}
}

func TestProvisioningStatus_HasTerminatedReturnsFalseForUnknownStates(t *testing.T) {
	if (provisioningStatus{provisioningProperties{"not a known state"}}).hasTerminated() {
		t.Errorf("azure: provisioningStatus#hasTerminated returned true for a non-terminal operation")
	}
}

func TestPollingState_HasSucceededReturnsFalseIfNotSuccess(t *testing.T) {
	if (pollingState{state: "not a success string"}).hasSucceeded() {
		t.Errorf("azure: pollingState#hasSucceeded failed to return false for a canceled operation")
	}
}

func TestPollingState_HasSucceededReturnsTrueIfSuccessful(t *testing.T) {
	if !(pollingState{state: operationSucceeded}).hasSucceeded() {
		t.Errorf("azure: pollingState#hasSucceeded failed to return true for a successful operation")
	}
}

func TestPollingState_HasTerminatedReturnsTrueForKnownStates(t *testing.T) {
	for _, state := range []string{operationSucceeded, operationCanceled, operationFailed} {
		if !(pollingState{state: state}).hasTerminated() {
			t.Errorf("azure: pollingState#hasTerminated failed to return true for the '%s' state", state)
		}
	}
}

func TestPollingState_HasTerminatedReturnsFalseForUnknownStates(t *testing.T) {
	if (pollingState{state: "not a known state"}).hasTerminated() {
		t.Errorf("azure: pollingState#hasTerminated returned true for a non-terminal operation")
	}
}

func TestNewPollingState_ReturnsAnErrorIfOneOccurs(t *testing.T) {
	resp := mocks.NewResponseWithContent(operationResourceIllegal)
	_, err := newPollingState(resp, false)
	if err == nil {
		t.Errorf("azure: newPollingState failed to return an error after a JSON parsing error")
	}
}

func TestNewPollingState_ReturnsTerminatedForKnownProvisioningStates(t *testing.T) {
	for _, state := range []string{operationSucceeded, operationCanceled, operationFailed} {
		resp := mocks.NewResponseWithContent(fmt.Sprintf(pollingStateFormat, state))
		resp.StatusCode = 42
		ps, _ := newPollingState(resp, false)
		if !ps.hasTerminated() {
			t.Errorf("azure: newPollingState failed to return a terminating pollingState for the '%s' state", state)
		}
	}
}

func TestNewPollingState_ReturnsSuccessForSuccessfulProvisioningState(t *testing.T) {
	resp := mocks.NewResponseWithContent(fmt.Sprintf(pollingStateFormat, operationSucceeded))
	resp.StatusCode = 42
	ps, _ := newPollingState(resp, false)
	if !ps.hasSucceeded() {
		t.Errorf("azure: newPollingState failed to return a successful pollingState for the '%s' state", operationSucceeded)
	}
}

func TestNewPollingState_ReturnsInProgressForAllOtherProvisioningStates(t *testing.T) {
	s := "not a recognized state"
	resp := mocks.NewResponseWithContent(fmt.Sprintf(pollingStateFormat, s))
	resp.StatusCode = 42
	ps, _ := newPollingState(resp, false)
	if ps.hasTerminated() {
		t.Errorf("azure: newPollingState returned terminated for unknown state '%s'", s)
	}
}

func TestNewPollingState_ReturnsSuccessWhenProvisioningStateFieldIsAbsentForSuccessStatusCodes(t *testing.T) {
	for _, sc := range []int{http.StatusOK, http.StatusCreated, http.StatusNoContent} {
		resp := mocks.NewResponseWithContent(pollingStateEmpty)
		resp.StatusCode = sc
		ps, _ := newPollingState(resp, false)
		if !ps.hasSucceeded() {
			t.Errorf("azure: newPollingState failed to return success when the provisionState field is absent for Status Code %d", sc)
		}
	}
}

func TestNewPollingState_ReturnsInProgressWhenProvisioningStateFieldIsAbsentForAccepted(t *testing.T) {
	resp := mocks.NewResponseWithContent(pollingStateEmpty)
	resp.StatusCode = http.StatusAccepted
	ps, _ := newPollingState(resp, false)
	if ps.hasTerminated() {
		t.Errorf("azure: newPollingState returned terminated when the provisionState field is absent for Status Code Accepted")
	}
}

func TestNewPollingState_ReturnsFailedWhenProvisioningStateFieldIsAbsentForUnknownStatusCodes(t *testing.T) {
	resp := mocks.NewResponseWithContent(pollingStateEmpty)
	resp.StatusCode = 42
	ps, _ := newPollingState(resp, false)
	if !ps.hasTerminated() || ps.hasSucceeded() {
		t.Errorf("azure: newPollingState did not return failed when the provisionState field is absent for an unknown Status Code")
	}
}

func TestNewPollingState_ReturnsTerminatedForKnownOperationResourceStates(t *testing.T) {
	for _, state := range []string{operationSucceeded, operationCanceled, operationFailed} {
		resp := mocks.NewResponseWithContent(fmt.Sprintf(operationResourceFormat, state))
		resp.StatusCode = 42
		ps, _ := newPollingState(resp, true)
		if !ps.hasTerminated() {
			t.Errorf("azure: newPollingState failed to return a terminating pollingState for the '%s' state", state)
		}
	}
}

func TestNewPollingState_ReturnsSuccessForSuccessfulOperationResourceState(t *testing.T) {
	resp := mocks.NewResponseWithContent(fmt.Sprintf(operationResourceFormat, operationSucceeded))
	resp.StatusCode = 42
	ps, _ := newPollingState(resp, true)
	if !ps.hasSucceeded() {
		t.Errorf("azure: newPollingState failed to return a successful pollingState for the '%s' state", operationSucceeded)
	}
}

func TestNewPollingState_ReturnsInProgressForAllOtherOperationResourceStates(t *testing.T) {
	s := "not a recognized state"
	resp := mocks.NewResponseWithContent(fmt.Sprintf(operationResourceFormat, s))
	resp.StatusCode = 42
	ps, _ := newPollingState(resp, true)
	if ps.hasTerminated() {
		t.Errorf("azure: newPollingState returned terminated for unknown state '%s'", s)
	}
}

func TestNewPollingState_CopiesTheResponseBody(t *testing.T) {
	s := fmt.Sprintf(pollingStateFormat, operationSucceeded)
	resp := mocks.NewResponseWithContent(s)
	resp.StatusCode = 42
	newPollingState(resp, true)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("azure: newPollingState failed to replace the http.Response Body -- Error='%v'", err)
	}
	if string(b) != s {
		t.Errorf("azure: newPollingState failed to copy the http.Response Body -- Expected='%s' Received='%s'", s, string(b))
	}
}

func TestNewPollingState_ClosesTheOriginalResponseBody(t *testing.T) {
	resp := mocks.NewResponse()
	b := resp.Body.(*mocks.Body)
	newPollingState(resp, false)
	if b.IsOpen() {
		t.Error("azure: newPollingState failed to close the original http.Response Body")
	}
}

func TestNewPollingRequest_FailsWhenResponseLacksRequest(t *testing.T) {
	isOperationResource := false
	resp := newAsynchronousResponse()
	resp.Request = nil

	_, err := newPollingRequest(resp, &isOperationResource)
	if err == nil {
		t.Error("azure: newPollingRequest failed to return an error when the http.Response lacked the original http.Request")
	}
}

func TestNewPollingRequest_PrefersTheAzureAsyncOperationHeader(t *testing.T) {
	isOperationResource := false
	resp := newAsynchronousResponse()

	req, _ := newPollingRequest(resp, &isOperationResource)
	if req.URL.String() != mocks.TestAzureAsyncURL {
		t.Error("azure: newPollingRequest failed to prefer the Azure-AsyncOperation header")
	}
}

func TestNewPollingRequest_ReturnsTrueWhenUsingTheAzureAsyncOperationHeader(t *testing.T) {
	isOperationResource := false
	resp := newAsynchronousResponse()

	newPollingRequest(resp, &isOperationResource)
	if !isOperationResource {
		t.Error("azure: newPollingRequest failed to return true when using the Azure-AsyncOperation header")
	}
}

func TestNewPollingRequest_PrefersLocationWhenTheAzureAsyncOperationHeaderMissing(t *testing.T) {
	isOperationResource := false
	resp := newAsynchronousResponse()
	resp.Header.Del(http.CanonicalHeaderKey(headerAsyncOperation))

	req, _ := newPollingRequest(resp, &isOperationResource)
	if req.URL.String() != mocks.TestLocationURL {
		t.Error("azure: newPollingRequest failed to prefer the Location header when the Azure-AsyncOperation header is missing")
	}
}

func TestNewPollingRequest_UsesTheObjectLocationIfAsyncHeadersAreMissing(t *testing.T) {
	isOperationResource := false
	resp := newAsynchronousResponse()
	resp.Header.Del(http.CanonicalHeaderKey(headerAsyncOperation))
	resp.Header.Del(http.CanonicalHeaderKey(autorest.HeaderLocation))
	resp.Request.Method = methodPatch

	req, _ := newPollingRequest(resp, &isOperationResource)
	if req.URL.String() != mocks.TestURL {
		t.Error("azure: newPollingRequest failed to use the Object URL when the asynchronous headers are missing")
	}
}

func TestNewPollingRequest_RecognizesLowerCaseHTTPVerbs(t *testing.T) {
	for _, m := range []string{"patch", "put"} {
		isOperationResource := false
		resp := newAsynchronousResponse()
		resp.Header.Del(http.CanonicalHeaderKey(headerAsyncOperation))
		resp.Header.Del(http.CanonicalHeaderKey(autorest.HeaderLocation))
		resp.Request.Method = m

		req, _ := newPollingRequest(resp, &isOperationResource)
		if req.URL.String() != mocks.TestURL {
			t.Errorf("azure: newPollingRequest failed to recognize the lower-case HTTP verb '%s'", m)
		}
	}
}

func TestNewPollingRequest_ReturnsAnErrorIfAsyncHeadersAreMissingForANewOrDeletedObject(t *testing.T) {
	isOperationResource := false
	resp := newAsynchronousResponse()
	resp.Header.Del(http.CanonicalHeaderKey(headerAsyncOperation))
	resp.Header.Del(http.CanonicalHeaderKey(autorest.HeaderLocation))

	for _, m := range []string{methodDelete, methodPost} {
		resp.Request.Method = m
		req, err := newPollingRequest(resp, &isOperationResource)
		if req != nil {
			t.Errorf("azure: newPollingRequest returned an http.Request even though it could not determine the polling URL for Method '%s'", m)
		}
		if err == nil {
			t.Errorf("azure: newPollingRequest failed to return an error even though it could not determine the polling URL for Method '%s'", m)
		}
	}
}

func TestNewPollingRequest_ReturnsAnErrorWhenPrepareFails(t *testing.T) {
	isOperationResource := false
	resp := newAsynchronousResponse()
	resp.Header.Set(http.CanonicalHeaderKey(headerAsyncOperation), mocks.TestBadURL)

	_, err := newPollingRequest(resp, &isOperationResource)
	if err == nil {
		t.Error("azure: newPollingRequest failed to return an error when Prepare fails")
	}
}

func TestNewPollingRequest_DoesNotReturnARequestWhenPrepareFails(t *testing.T) {
	isOperationResource := false
	resp := newAsynchronousResponse()
	resp.Header.Set(http.CanonicalHeaderKey(headerAsyncOperation), mocks.TestBadURL)

	req, _ := newPollingRequest(resp, &isOperationResource)
	if req != nil {
		t.Error("azure: newPollingRequest returned an http.Request when Prepare failed")
	}
}

func TestNewPollingRequest_ReturnsAGetRequest(t *testing.T) {
	isOperationResource := false
	resp := newAsynchronousResponse()

	req, _ := newPollingRequest(resp, &isOperationResource)
	if req.Method != "GET" {
		t.Errorf("azure: newPollingRequest did not create an HTTP GET request -- actual method %v", req.Method)
	}
}

func TestDoPollForAsynchronous_IgnoresUnspecifiedStatusCodes(t *testing.T) {
	client := mocks.NewSender()

	r, _ := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Duration(0)))

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
		DoPollForAsynchronous(time.Millisecond))

	if client.Attempts() != 2 {
		t.Errorf("azure: DoPollForAsynchronous failed to poll for specified status code")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_CanBeCanceled(t *testing.T) {
	cancel := make(chan struct{})
	delay := 5 * time.Second

	r1 := newAsynchronousResponse()

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(newOperationResourceResponse("Busy"), -1)

	var wg sync.WaitGroup
	wg.Add(1)
	start := time.Now()
	go func() {
		req := mocks.NewRequest()
		req.Cancel = cancel

		wg.Done()

		r, _ := autorest.SendWithSender(client, req,
			DoPollForAsynchronous(10*time.Second))
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
	b1 := r1.Body.(*mocks.Body)
	r2 := newOperationResourceResponse("busy")
	b2 := r2.Body.(*mocks.Body)
	r3 := newOperationResourceResponse(operationSucceeded)
	b3 := r3.Body.(*mocks.Body)

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)
	client.AppendResponse(r3)

	r, _ := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond))

	if b1.IsOpen() || b2.IsOpen() || b3.IsOpen() {
		t.Errorf("azure: DoPollForAsynchronous did not close unreturned response bodies")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_LeavesLastResponseBodyOpen(t *testing.T) {
	r1 := newAsynchronousResponse()
	r2 := newOperationResourceResponse("busy")
	r3 := newOperationResourceResponse(operationSucceeded)

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)
	client.AppendResponse(r3)

	r, _ := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond))

	b, err := ioutil.ReadAll(r.Body)
	if len(b) <= 0 || err != nil {
		t.Errorf("azure: DoPollForAsynchronous did not leave open the body of the last response - Error='%v'", err)
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_DoesNotPollIfOriginalRequestReturnedAnError(t *testing.T) {
	r1 := newAsynchronousResponse()
	r2 := newOperationResourceResponse("busy")

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendResponse(r2)
	client.SetError(fmt.Errorf("Faux Error"))

	r, _ := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond))

	if client.Attempts() != 1 {
		t.Errorf("azure: DoPollForAsynchronous tried to poll after receiving an error")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_DoesNotPollIfCreatingOperationRequestFails(t *testing.T) {
	r1 := newAsynchronousResponse()
	mocks.SetResponseHeader(r1, http.CanonicalHeaderKey(headerAsyncOperation), mocks.TestBadURL)
	r2 := newOperationResourceResponse("busy")

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)

	r, _ := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond))

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
		DoPollForAsynchronous(time.Millisecond))

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
		DoPollForAsynchronous(time.Millisecond))

	if err == nil {
		t.Errorf("azure: DoPollForAsynchronous failed to return error from polling")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_PollsUntilOperationResourceHasTerminated(t *testing.T) {
	r1 := newAsynchronousResponse()
	r2 := newOperationResourceResponse("busy")
	r3 := newOperationResourceResponse(operationCanceled)

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)
	client.AppendAndRepeatResponse(r3, 1)

	r, _ := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond))

	if client.Attempts() < 4 {
		t.Errorf("azure: DoPollForAsynchronous stopped polling before receiving a terminated OperationResource")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_StopsPollingWhenOperationResourceHasTerminated(t *testing.T) {
	r1 := newAsynchronousResponse()
	r2 := newOperationResourceResponse("busy")
	r3 := newOperationResourceResponse(operationCanceled)

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)
	client.AppendAndRepeatResponse(r3, 2)

	r, _ := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond))

	if client.Attempts() > 4 {
		t.Errorf("azure: DoPollForAsynchronous failed to stop after receiving a terminated OperationResource")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_ReturnsAnErrorForCanceledOperations(t *testing.T) {
	r1 := newAsynchronousResponse()
	r2 := newOperationResourceResponse("busy")
	r3 := newOperationResourceErrorResponse(operationCanceled)

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)
	client.AppendAndRepeatResponse(r3, 1)

	r, err := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond))

	if err == nil || !strings.Contains(fmt.Sprintf("%v", err), "Canceled") {
		t.Errorf("azure: DoPollForAsynchronous failed to return an appropriate error for a canceled OperationResource")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_ReturnsAnErrorForFailedOperations(t *testing.T) {
	r1 := newAsynchronousResponse()
	r2 := newOperationResourceResponse("busy")
	r3 := newOperationResourceErrorResponse(operationFailed)

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)
	client.AppendAndRepeatResponse(r3, 1)

	r, err := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond))

	if err == nil || !strings.Contains(fmt.Sprintf("%v", err), "Failed") {
		t.Errorf("azure: DoPollForAsynchronous failed to return an appropriate error for a canceled OperationResource")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_ReturnsNoErrorForSuccessfulOperations(t *testing.T) {
	r1 := newAsynchronousResponse()
	r2 := newOperationResourceResponse("busy")
	r3 := newOperationResourceErrorResponse(operationSucceeded)

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)
	client.AppendAndRepeatResponse(r3, 1)

	r, err := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond))

	if err != nil {
		t.Errorf("azure: DoPollForAsynchronous returned an error for a successful OperationResource")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

func TestDoPollForAsynchronous_StopsPollingIfItReceivesAnInvalidOperationResource(t *testing.T) {
	r1 := newAsynchronousResponse()
	r2 := newOperationResourceResponse("busy")
	r3 := newOperationResourceResponse("busy")
	r3.Body = mocks.NewBody(operationResourceIllegal)
	r4 := newOperationResourceResponse(operationSucceeded)

	client := mocks.NewSender()
	client.AppendResponse(r1)
	client.AppendAndRepeatResponse(r2, 2)
	client.AppendAndRepeatResponse(r3, 1)
	client.AppendAndRepeatResponse(r4, 1)

	r, err := autorest.SendWithSender(client, mocks.NewRequest(),
		DoPollForAsynchronous(time.Millisecond))

	if client.Attempts() > 4 {
		t.Errorf("azure: DoPollForAsynchronous failed to stop polling after receiving an invalid OperationResource")
	}
	if err == nil {
		t.Errorf("azure: DoPollForAsynchronous failed to return an error after receving an invalid OperationResource")
	}

	autorest.Respond(r,
		autorest.ByClosing())
}

const (
	operationResourceIllegal = `
	This is not JSON and should fail...badly.
	`

	pollingStateFormat = `
	{
		"unused" : {
			"somefield" : 42
		},
		"properties" : {
			"provisioningState": "%s"
		}
	}
	`

	pollingStateEmpty = `
	{
		"unused" : {
			"somefield" : 42
		},
		"properties" : {
		}
	}
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

func newAsynchronousResponse() *http.Response {
	r := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	r.Body = mocks.NewBody(fmt.Sprintf(pollingStateFormat, operationInProgress))
	mocks.SetResponseHeader(r, http.CanonicalHeaderKey(headerAsyncOperation), mocks.TestAzureAsyncURL)
	mocks.SetResponseHeader(r, http.CanonicalHeaderKey(autorest.HeaderLocation), mocks.TestLocationURL)
	mocks.SetRetryHeader(r, retryDelay)
	r.Request = mocks.NewRequestForURL(mocks.TestURL)
	return r
}

func newOperationResourceResponse(status string) *http.Response {
	r := newAsynchronousResponse()
	r.Body = mocks.NewBody(fmt.Sprintf(operationResourceFormat, status))
	return r
}

func newOperationResourceErrorResponse(status string) *http.Response {
	r := newAsynchronousResponse()
	r.Body = mocks.NewBody(fmt.Sprintf(operationResourceErrorFormat, status))
	return r
}
