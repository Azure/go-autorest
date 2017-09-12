package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/mitchellh/go-homedir"
)

// AzureCLIToken represents an AccessToken from the Azure CLI
type AzureCLIToken struct {
	AccessToken      string `json:"accessToken"`
	Authority        string `json:"_authority"`
	ClientID         string `json:"_clientId"`
	ExpiresOn        string `json:"expiresOn"`
	IdentityProvider string `json:"identityProvider"`
	IsMRRT           bool   `json:"isMRRT"`
	RefreshToken     string `json:"refreshToken"`
	Resource         string `json:"resource"`
	TokenType        string `json:"tokenType"`
	UserID           string `json:"userId"`
}

// ToToken converts an AzureCLIToken to a Token
func (t AzureCLIToken) ToToken() (*adal.Token, error) {
	tokenExpirationDate, err := ParseAzureCLIExpirationDate(t.ExpiresOn)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Token Expiration Date %q: %+v", t.ExpiresOn, err)
	}

	difference := tokenExpirationDate.Sub(adal.ExpirationBase)
	token := adal.Token{
		AccessToken:  t.AccessToken,
		Type:         t.TokenType,
		ExpiresIn:    "3600",
		ExpiresOn:    strconv.Itoa(int(difference.Seconds())),
		RefreshToken: t.RefreshToken,
		Resource:     t.Resource,
	}
	return &token, nil
}

// AzureCLIAccessTokensPath returns the path where access tokens are stored from the Azure CLI
func AzureCLIAccessTokensPath() (string, error) {
	return homedir.Expand("~/.azure/accessTokens.json")
}

// ParseAzureCLIExpirationDate parses either a Azure CLI or CloudShell date into a time object
func ParseAzureCLIExpirationDate(input string) (*time.Time, error) {
	// CloudShell (and potentially the Azure CLI in future)
	expirationDate, cloudShellErr := time.Parse(time.RFC3339, input)
	if cloudShellErr != nil {
		// Azure CLI (Python) e.g. 2017-08-31 19:48:57.998857 (plus the local timezone)
		cliFormat := "2006-01-02 15:04:05.999999"
		expirationDate, cliErr := time.ParseInLocation(cliFormat, input, time.Local)
		if cliErr == nil {
			return &expirationDate, nil
		}

		return nil, fmt.Errorf("Error parsing expiration date %q.\n\nCloudShell Error: \n%+v\n\nCLI Error:\n%+v", input, cloudShellErr, cliErr)
	}

	return &expirationDate, nil
}
