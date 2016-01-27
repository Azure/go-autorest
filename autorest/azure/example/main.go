package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"golang.org/x/crypto/pkcs12"
)

const resourceGroupURLTemplate = "https://management.azure.com/subscriptions/{subscription-id}/resourcegroups"
const apiVersion = "2015-01-01"
const xplatClientID = "04b07795-8ddb-461a-bbee-02f9e1bf7b46"

var (
	mode           string
	tenantID       string
	subscriptionID string
	applicationID  string

	tokenCachePath string

	certificatePath string
)

func init() {
	flag.StringVar(&mode, "mode", "device", "mode of operation for SPT creation")
	flag.StringVar(&certificatePath, "certificatePath", "", "path to pk12/pfx certificate")
	flag.StringVar(&applicationID, "applicationId", "", "application id")
	flag.StringVar(&tenantID, "tenantId", "", "tenant id")
	flag.StringVar(&subscriptionID, "subscriptionId", "", "subscription id")
	flag.StringVar(&tokenCachePath, "tokenCachePath", "", "location of oauth token cache")

	flag.Parse()

	log.Printf("mode(%s) certPath(%s) appID(%s) tenantID(%s), subID(%s)\n",
		mode, certificatePath, applicationID, tenantID, subscriptionID)

	if strings.Trim(tenantID, " ") == "" ||
		strings.Trim(subscriptionID, " ") == "" {
		log.Fatalln("Bad usage. Please specify applicationID, tenantID, subscriptionID")
	}

	if mode != "certificate" && mode != "cached" && mode != "device" {
		log.Fatalln("Bad usage. Mode must be one of 'certificate', 'cached' or 'device'.")
	}

	if mode == "device" && applicationID == "" {
		applicationID = xplatClientID
	}

	if mode == "certificate" && strings.Trim(certificatePath, " ") == "" {
		log.Fatalln("Bad usage. Mode 'certificate' requires the 'certificatePath' argument.")
	}
}

func getSptFromCachedToken(clientID, tenantID, resource string, callbacks ...azure.TokenRefreshCallback) (*azure.ServicePrincipalToken, error) {
	token, err := azure.LoadToken(tokenCachePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load token from cache: %v", err)
	}

	spt := azure.NewServicePrincipalTokenFromManualToken(
		clientID,
		tenantID,
		resource,
		*token,
		callbacks...)

	return spt, nil
}

func decodePkcs12(pkcs []byte, password string) (*x509.Certificate, *rsa.PrivateKey, error) {
	privateKey, certificate, err := pkcs12.Decode(pkcs, password)
	if err != nil {
		return nil, nil, err
	}

	rsaPrivateKey, isRsaKey := privateKey.(*rsa.PrivateKey)
	if !isRsaKey {
		return nil, nil, fmt.Errorf("PKCS#12 certificate must contain an RSA private key")
	}

	return certificate, rsaPrivateKey, nil
}

func getSptFromCertificate(clientID, tenantID, resource, certicatePath string, callbacks ...azure.TokenRefreshCallback) (*azure.ServicePrincipalToken, error) {
	certData, err := ioutil.ReadFile(certificatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read the certificate file (%s): %v", certificatePath, err)
	}

	certificate, rsaPrivateKey, err := decodePkcs12(certData, "")
	if err != nil {
		return nil, fmt.Errorf("failed to decode pkcs12 certificate while creating spt")
	}

	spt := azure.NewServicePrincipalTokenFromCertificate(
		clientID,
		certificate,
		rsaPrivateKey,
		tenantID,
		azure.AzureResourceManagerScope,
		callbacks...)

	return spt, nil
}

func getSptFromDeviceFlow(clientID, tenantID, resource string, callbacks ...azure.TokenRefreshCallback) (*azure.ServicePrincipalToken, error) {
	oauthClient := &autorest.Client{}
	deviceCode, err := azure.InitiateDeviceAuth(oauthClient, clientID, azure.AzureResourceManagerScope)
	if err != nil {
		return nil, fmt.Errorf("failed to start device auth flow: %s", err)
	}

	fmt.Println(*deviceCode.Message)

	token, err := azure.WaitForUserCompletion(oauthClient, clientID, deviceCode)
	if err != nil {
		return nil, fmt.Errorf("failed to finish device auth flow: %s", err)
	}

	spt := azure.NewServicePrincipalTokenFromManualToken(
		clientID,
		tenantID,
		resource,
		*token,
		callbacks...)

	return spt, nil
}

func getResourceGroups(client *autorest.Client) ([]string, error) {
	p := map[string]interface{}{"subscription-id": subscriptionID}
	q := map[string]interface{}{"api-version": apiVersion}

	req, _ := autorest.Prepare(&http.Request{},
		autorest.AsGet(),
		autorest.WithBaseURL(resourceGroupURLTemplate),
		autorest.WithPathParameters(p),
		autorest.WithQueryParameters(q))

	resp, err := client.Send(req)
	if err != nil {
		return nil, err
	}

	value := struct {
		ResourceGroups []struct {
			Name string `json:"name"`
		} `json:"value"`
	}{}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&value)
	if err != nil {
		return nil, err
	}

	var names = make([]string, len(value.ResourceGroups))
	for i, name := range value.ResourceGroups {
		names[i] = name.Name
	}
	return names, nil
}

func saveToken(spt azure.Token) {
	err := azure.SaveToken(tokenCachePath, spt)
	if err != nil {
		log.Println("error saving token", err)
	} else {
		log.Println("saved token to", tokenCachePath)
	}
}

func main() {
	var spt *azure.ServicePrincipalToken
	var err error

	resource := azure.AzureResourceManagerScope

	callback := func(t azure.Token) error {
		log.Println("refresh callback was called because the cached oauth token was stale")
		saveToken(spt.Token)
		return nil
	}

	if tokenCachePath != "" {
		log.Println("tokenCachePath specified; attempting to load from", tokenCachePath)
		spt, err = getSptFromCachedToken(applicationID, tenantID, resource, callback)
		if err != nil {
			spt = nil // just in case, this is the condition below
			log.Println("loading from cache failed:", err)
		}
	}

	if spt == nil {
		log.Println("authenticating via 'mode'", mode)
		switch mode {
		case "device":
			spt, err = getSptFromDeviceFlow(applicationID, tenantID, resource, callback)
		case "certificate":
			spt, err = getSptFromCertificate(applicationID, tenantID, resource, certificatePath, callback)
		}
		if err != nil {
			log.Fatalln("failed to retrieve token:", err)
		}

		// should save it as soon as you get it since Refresh won't be called for some time
		if tokenCachePath != "" {
			saveToken(spt.Token)
		}
	}

	client := &autorest.Client{}
	client.Authorizer = spt

	groupNames, err := getResourceGroups(client)
	if err != nil {
		log.Fatalln("failed to retrieve groups:", err)
	}

	log.Println("Groups:", strings.Join(groupNames, ","))
}
