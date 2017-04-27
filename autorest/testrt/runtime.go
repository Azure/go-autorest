package testrt

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"io/ioutil"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Credentials contains data used for authentication.
type Credentials struct {
	TenantID string `json:"tenantId"`
	ClientID string `json:"clientId"`
	Secret   string `json:"secret"`
}

// ReservedParams contains the contents of the reserved block sent from the client.
type ReservedParams struct {
	Headers     http.Header `json:"headers"`
	Credentials Credentials `json:"credentials"`
}

// RawResponse is used to marshal a HTTP response to the client.
type RawResponse struct {
	StatusCode int             `json:"statusCode"`
	Headers    http.Header     `json:"headers"`
	Response   json.RawMessage `json:"response"`
}

func (rr RawResponse) String() string {
	return fmt.Sprintf("%v %v %v", rr.StatusCode, rr.Headers, string(rr.Response))
}

// ToRawResponse converts an autorest response to a RawResponse for marshaling.
func ToRawResponse(r *autorest.Response, v interface{}) (*RawResponse, error) {
	var b []byte
	var err error

	if v != nil {
		b, err = json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal raw response %v", err)
		}
	}
	rr := RawResponse{
		StatusCode: r.StatusCode,
		Headers:    r.Header,
		Response:   b,
	}
	return &rr, nil
}

type syncFile struct {
	file *os.File
}

func (fw syncFile) Write(b []byte) (n int, err error) {
	n, err = fw.file.Write(b)
	// immediately sync the file so the content show up in real time
	fw.file.Sync()
	return
}

// CreateSyncronousLogfile creates a logger that synchronously writes to a log file in the temp directory.
func CreateSyncronousLogfile(prefix string) (*log.Logger, error) {
	logFile, err := ioutil.TempFile("", prefix)
	if err != nil {
		return nil, err
	}
	sf := syncFile{
		file: logFile,
	}
	return log.New(sf, ">>> ", log.LstdFlags|log.Lshortfile), nil
}

// CreateBearerAuthorizer creates a bearer authorizer based on the current credentials.
func (c Credentials) CreateBearerAuthorizer() (autorest.Authorizer, error) {
	oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, c.TenantID)
	if err != nil {
		return nil, err
	}
	spToken, err := adal.NewServicePrincipalToken(*oauthConfig, c.ClientID, c.Secret, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return nil, err
	}
	return autorest.NewBearerAuthorizer(spToken), nil
}
