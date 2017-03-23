package adal

import (
	"fmt"
	"net/url"
)

const (
	activeDirectoryAPIVersion       = "1.0"
	activeDirectoryEndpointTemplate = "%s/oauth2/%s?api-version=%s"
)

// OAuthConfig represents the endpoints needed
// in OAuth operations
type OAuthConfig struct {
	AuthorizeEndpoint  url.URL
	TokenEndpoint      url.URL
	DeviceCodeEndpoint url.URL
}

// NewOAuthConfig returns an OAuthConfig with tenant specific urls
func NewOAuthConfig(activeDirectoryEndpoint, tenantID string) (*OAuthConfig, error) {
	u, err := url.Parse(activeDirectoryEndpoint)
	if err != nil {
		return nil, err
	}
	authorizeURL, err := u.Parse(fmt.Sprintf(activeDirectoryEndpointTemplate, tenantID, "authorize", activeDirectoryAPIVersion))
	if err != nil {
		return nil, err
	}
	tokenURL, err := u.Parse(fmt.Sprintf(activeDirectoryEndpointTemplate, tenantID, "token", activeDirectoryAPIVersion))
	if err != nil {
		return nil, err
	}
	deviceCodeURL, err := u.Parse(fmt.Sprintf(activeDirectoryEndpointTemplate, tenantID, "devicecode", activeDirectoryAPIVersion))
	if err != nil {
		return nil, err
	}

	return &OAuthConfig{
		AuthorizeEndpoint:  *authorizeURL,
		TokenEndpoint:      *tokenURL,
		DeviceCodeEndpoint: *deviceCodeURL,
	}, nil
}
