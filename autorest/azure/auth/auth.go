package auth

// Copyright 2017 Microsoft Corporation
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"golang.org/x/crypto/pkcs12"
)

// GetAuthorizerWithDefaults tries to get an authorizer using the following methods:
// 1. Client credentials
// 2. Client certificate
// 3. MSI
func GetAuthorizerWithDefaults() (*autorest.BearerAuthorizer, error) {

	options := NewAuthorizerOptionsFromEnvironment()

	//1.Client Credentials
	authorizer, err := GetAuthorizerFromClientCredentials(options)

	//2. Client Certificate
	if err != nil {
		authorizer, err = GetAuthorizerFromClientCertificate(options)
	}

	//3. MSI
	if err != nil {
		authorizer, err = GetAuthorizerFromMSI(options.Environment)
	}

	return authorizer, nil
}

// GetAuthorizerFromMSI gets an authorizer from MSI. Note that this will only work when running in an Azure environment
func GetAuthorizerFromMSI(env azure.Environment) (*autorest.BearerAuthorizer, error) {
	msiEndpoint, err := adal.GetMSIVMEndpoint()

	if err != nil {
		return nil, err
	}

	spToken, err := adal.NewServicePrincipalTokenFromMSI(msiEndpoint, env.ResourceManagerEndpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to get oauth token from MSI: %v", err)
	}

	return autorest.NewBearerAuthorizer(spToken), nil
}

// GetAuthorizerFromClientCredentials gets an authorizer from client credentials.
func GetAuthorizerFromClientCredentials(options AuthorizerOptions) (*autorest.BearerAuthorizer, error) {
	oauthConfig, err := adal.NewOAuthConfig(options.Environment.ActiveDirectoryEndpoint, options.TenantID)
	if err != nil {
		return nil, err
	}

	spToken, err := adal.NewServicePrincipalToken(*oauthConfig, options.ClientID, options.ClientSecret, options.Environment.ResourceManagerEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth token from client credentials: %v", err)
	}

	return autorest.NewBearerAuthorizer(spToken), nil
}

// GetAuthorizerFromClientCertificate gets an authorizer from a client certificate.
func GetAuthorizerFromClientCertificate(options AuthorizerOptions) (*autorest.BearerAuthorizer, error) {
	oauthConfig, err := adal.NewOAuthConfig(options.Environment.ActiveDirectoryEndpoint, options.TenantID)

	certData, err := ioutil.ReadFile(options.CertificatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read the certificate file (%s): %v", options.CertificatePath, err)
	}

	certificate, rsaPrivateKey, err := decodePkcs12(certData, options.CertificatePassword)
	if err != nil {
		return nil, fmt.Errorf("failed to decode pkcs12 certificate while creating spt: %v", err)
	}

	spToken, err := adal.NewServicePrincipalTokenFromCertificate(*oauthConfig, options.ClientID, certificate, rsaPrivateKey, options.Environment.ResourceManagerEndpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to get oauth token from certificate auth: %v", err)
	}

	return autorest.NewBearerAuthorizer(spToken), nil
}

// GetAuthorizerFromDeviceFlow gets an authorizer from device flow.
func GetAuthorizerFromDeviceFlow(options AuthorizerOptions) (*autorest.BearerAuthorizer, error) {
	oauthClient := &autorest.Client{}
	oauthConfig, err := adal.NewOAuthConfig(options.Environment.ActiveDirectoryEndpoint, options.TenantID)
	deviceCode, err := adal.InitiateDeviceAuth(oauthClient, *oauthConfig, options.ClientID, options.Environment.ActiveDirectoryEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to start device auth flow: %s", err)
	}

	fmt.Println(*deviceCode.Message)

	token, err := adal.WaitForUserCompletion(oauthClient, deviceCode)
	if err != nil {
		return nil, fmt.Errorf("failed to finish device auth flow: %s", err)
	}

	spToken, err := adal.NewServicePrincipalTokenFromManualToken(*oauthConfig, options.ClientID, options.Environment.ResourceManagerEndpoint, *token)
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth token from device flow: %v", err)
	}

	return autorest.NewBearerAuthorizer(spToken), nil
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

// AuthorizerOptions provides the options to get an authorizer
type AuthorizerOptions struct {
	TenantID            string
	ClientID            string
	ClientSecret        string
	CertificatePath     string
	CertificatePassword string
	Environment         azure.Environment
}

// NewAuthorizerOptionsForClientCredentials return an AuthorizerOptions object configured to obtain an Authorizer through Client Credentials.
func NewAuthorizerOptionsForClientCredentials(tenantID string, clientID string, clientSecret string, env azure.Environment) AuthorizerOptions {
	return AuthorizerOptions{
		TenantID:     tenantID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Environment:  env,
	}
}

// NewAuthorizerOptionsForClientCertificate return an AuthorizerOptions object configured to obtain an Authorizer through client certificate.
func NewAuthorizerOptionsForClientCertificate(tenantID string, clientID string, certificatePath string, certificatePassword string, env azure.Environment) AuthorizerOptions {
	return AuthorizerOptions{
		TenantID:            tenantID,
		ClientID:            clientID,
		CertificatePath:     certificatePath,
		CertificatePassword: certificatePassword,
		Environment:         env,
	}
}

// NewAuthorizerOptionsFromEnvironment return an AuthorizerOptions configured from environment variables.
func NewAuthorizerOptionsFromEnvironment() AuthorizerOptions {
	options := AuthorizerOptions{}
	options.TenantID = os.Getenv("AZURE_TENANT_ID")
	options.ClientID = os.Getenv("AZURE_CLIENT_ID")
	options.ClientSecret = os.Getenv("AZURE_CLIENT_SECRET")
	options.CertificatePath = os.Getenv("AZURE_CERTIFICATE_PATH")
	options.CertificatePassword = os.Getenv("AZURE_CERTIFICATE_PASSWORD")

	envName := os.Getenv("AZURE_ENVIRONMENT")
	if envName == "" {
		options.Environment = azure.PublicCloud
	} else {
		env, err := azure.EnvironmentFromName(envName)
		if err != nil {
			panic(err)
		}
		options.Environment = env
	}

	return options
}
