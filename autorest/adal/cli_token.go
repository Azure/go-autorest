package adal

import "strconv"

type AzureCLIToken struct {
	AccessToken      string `json:"accessToken"`
	Authority        string `json:"_authority"`
	ClientID         string `json:"_clientId"`
	ExpiresIn        int     `json:"expiresIn"`
	ExpiresOn        string `json:"expiresOn"`
	IdentityProvider string `json:"identityProvider"`
	IsMRRT           bool    `json:"isMRRT"`
	RefreshToken     string `json:"refreshToken"`
	Resource         string `json:"resource"`
	TokenType        string `json:"tokenType"`
	UserId           string `json:"userId"`
}

func AzureCLIAccessTokensPath() string {
	return "~/.azure/accessTokens.json"
}

func (t AzureCLIToken) ToToken() Token {
	return Token{
		AccessToken: t.AccessToken,
		Type: t.TokenType,
		ExpiresIn: strconv.Itoa(t.ExpiresIn),
		ExpiresOn: t.ExpiresOn,
		RefreshToken: t.RefreshToken,
		Resource: t.Resource,
	}
}