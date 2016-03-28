package azure

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
)

const (
	headerAsyncOperation = "Azure-AsyncOperation"
)

const (
	methodDelete = "DELETE"
	methodPatch  = "PATCH"
	methodPost   = "POST"
	methodPut    = "PUT"
	methodGet    = "GET"

	operationInProgress string = "InProgress"
	operationCanceled   string = "Canceled"
	operationFailed     string = "Failed"
	operationSucceeded  string = "Succeeded"
)

// DoPollForAsynchronous returns a SendDecorator that polls if the http.Response is for an Azure
// long-running operation. It will delay between requests for the duration specified in the
// RetryAfter header or, if the header is absent, the passed delay. Polling may be canceled by
// closing the optional channel on the http.Request.
func DoPollForAsynchronous(delay time.Duration) autorest.SendDecorator {
	return func(s autorest.Sender) autorest.Sender {
		return autorest.SenderFunc(func(r *http.Request) (resp *http.Response, err error) {
			resp, err = s.Do(r)
			if err != nil {
				return resp, err
			}

			// Note:
			// 	newPollingState maps the operation status -- retrieved from either a provisioningState
			// 	field, the status field of an OperationResource, or inferred from the HTTP status code --
			// 	into a well-known states. Since the process begins from the initial request, the state
			//	always comes from either a the provisioningState returned or is inferred from the HTTP
			//	status code. Subsequent requests will return an Azure OperationResource object if the
			//	service returns the Azure-AsyncOperation header. The boolean, isAzureAsyncOperation,
			//	conveys to newPollingState if it should expect a provisioningState field or an
			//	OperationResource object.
			//

			// Ensure a cancel channel exists since the loop requires it and it is optional.
			if resp.Request.Cancel == nil {
				resp.Request.Cancel = make(chan struct{})
			}

			isAzureAsyncOperation := false
			ps := pollingState{state: operationFailed}
			for err == nil {
				ps, err = newPollingState(resp, isAzureAsyncOperation)
				if err != nil {
					break
				}
				if ps.hasTerminated() {
					if !ps.hasSucceeded() {
						err = autorest.NewError("azure", "DoPollForAsynchronous", "Polling terminated with status '%s'", ps.state)
					}
					break
				}

				r, err = newPollingRequest(resp, &isAzureAsyncOperation)
				if err != nil {
					return resp, err
				}

				resp, err = autorest.SendWithSender(s, r,
					autorest.AfterDelay(autorest.GetRetryAfter(resp, delay)))
			}

			return resp, err
		})
	}
}

func getAsyncOperation(resp *http.Response) string {
	return resp.Header.Get(http.CanonicalHeaderKey(headerAsyncOperation))
}

func hasSucceeded(state string) bool {
	return state == operationSucceeded
}

func hasTerminated(state string) bool {
	switch state {
	case operationCanceled, operationFailed, operationSucceeded:
		return true
	default:
		return false
	}
}

type provisioningTracker interface {
	state() string
	hasSucceeded() bool
	hasTerminated() bool
}

type operationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (oe operationError) Error() string {
	return fmt.Sprintf("Azure Operation Error: Code=%q Message=%q", oe.Code, oe.Message)
}

type operationResource struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Status          string                 `json:"status"`
	Properties      map[string]interface{} `json:"properties"`
	OperationError  operationError         `json:"error"`
	StartTime       date.Time              `json:"startTime"`
	EndTime         date.Time              `json:"endTime"`
	PercentComplete float64                `json:"percentComplete"`
}

func (or operationResource) state() string {
	return or.Status
}

func (or operationResource) hasSucceeded() bool {
	return hasSucceeded(or.state())
}

func (or operationResource) hasTerminated() bool {
	return hasTerminated(or.state())
}

type provisioningProperties struct {
	ProvisioningState string `json:"provisioningState"`
}

type provisioningStatus struct {
	Properties provisioningProperties `json:"properties"`
}

func (ps provisioningStatus) state() string {
	return ps.Properties.ProvisioningState
}

func (ps provisioningStatus) hasSucceeded() bool {
	return hasSucceeded(ps.state())
}

func (ps provisioningStatus) hasTerminated() bool {
	return hasTerminated(ps.state())
}

type pollingState struct {
	state   string
	code    string
	message string
}

func (ps pollingState) hasSucceeded() bool {
	return hasSucceeded(ps.state)
}

func (ps pollingState) hasTerminated() bool {
	return hasTerminated(ps.state)
}

func newPollingState(resp *http.Response, isAzureAsyncOperation bool) (pollingState, error) {
	pollState := pollingState{state: operationFailed}

	var pt provisioningTracker
	if isAzureAsyncOperation {
		pt = &operationResource{}
	} else {
		pt = &provisioningStatus{}
	}

	b := &bytes.Buffer{}
	err := autorest.Respond(resp,
		autorest.ByCopying(b),
		autorest.ByUnmarshallingJSON(pt),
		autorest.ByClosing())
	resp.Body = ioutil.NopCloser(b)
	if err != nil {
		return pollState, err
	}

	// -- Terminal states apply regardless
	// -- Unknown states are per-service inprogress states
	// -- Otherwise, infer state from HTTP status code
	if pt.hasTerminated() {
		pollState.state = pt.state()
	} else if pt.state() != "" {
		pollState.state = operationInProgress
	} else {
		switch resp.StatusCode {
		case http.StatusAccepted:
			pollState.state = operationInProgress

		case http.StatusNoContent, http.StatusCreated, http.StatusOK:
			pollState.state = operationSucceeded

		default:
			pollState.state = operationFailed
		}
	}

	return pollState, nil
}

func newPollingRequest(resp *http.Response, isAzureAsyncOperation *bool) (*http.Request, error) {
	req := resp.Request
	if req == nil {
		return nil, autorest.NewError("azure", "newPollingRequest", "Azure Polling Error - Original HTTP request is missing")
	}

	// Prefer the Azure-AsyncOperation header
	uri := getAsyncOperation(resp)
	*isAzureAsyncOperation = (uri != "")

	// Else, use the Location header
	if !*isAzureAsyncOperation {
		uri = autorest.GetLocation(resp)
	}

	// Lastly, requests against an existing resource, use the last request URI
	if uri == "" {
		req.Method = strings.ToUpper(req.Method)
		if req.Method == methodPatch || req.Method == methodPut || req.Method == methodGet {
			uri = req.URL.String()
		}
	}

	if uri == "" {
		return nil, autorest.NewError("azure", "newPollingRequest", "Azure Polling Error - Unable to obtain polling URI for %s %s", req.Method, req.RequestURI)
	}

	reqPoll, err := autorest.Prepare(&http.Request{Cancel: req.Cancel},
		autorest.AsGet(),
		autorest.WithBaseURL(uri))
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "azure", "newPollingRequest", nil, "Failure creating poll request to %s", uri)
	}

	return reqPoll, nil
}
