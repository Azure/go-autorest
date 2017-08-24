package adal

import (
	"github.com/mitchellh/go-homedir"
	"strconv"
)

// AzureCLIToken represents an AccessToken from the Azure CLI
type AzureCLIToken struct {
	AccessToken      string `json:"accessToken"`
	Authority        string `json:"_authority"`
	ClientID         string `json:"_clientId"`
	ExpiresIn        int    `json:"expiresIn"`
	ExpiresOn        string `json:"expiresOn"`
	IdentityProvider string `json:"identityProvider"`
	IsMRRT           bool   `json:"isMRRT"`
	RefreshToken     string `json:"refreshToken"`
	Resource         string `json:"resource"`
	TokenType        string `json:"tokenType"`
	UserID           string `json:"userId"`
}

// AzureCLIAccessTokensPath returns the path where access tokens are stored from the Azure CLI
func AzureCLIAccessTokensPath() (string, error) {
	return homedir.Expand("~/.azure/accessTokens.json")
}

// ToToken converts an AzureCLIToken to a Token
func (t AzureCLIToken) ToToken() Token {
	return Token{
		AccessToken:  t.AccessToken,
		Type:         t.TokenType,
		ExpiresIn:    strconv.Itoa(t.ExpiresIn),
		ExpiresOn:    t.ExpiresOn,
		RefreshToken: t.RefreshToken,
		Resource:     t.Resource,
	}
}
