package auth

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode/utf16"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/dimchansky/utfbom"
)

// Auth includes authentication details for ARM clients
type Auth struct {
	Authorizer *autorest.BearerAuthorizer
	File
	BaseURI string
}

// File represents the authentication file
type File struct {
	ClientID                string `json:"clientId,ompitempty"`
	ClientSecret            string `json:"clientSecret,ompitempty"`
	SubscriptionID          string `json:"subscriptionId,ompitempty"`
	TenantID                string `json:"tenantId,ompitempty"`
	ActiveDirectoryEndpoint string `json:"activeDirectoryEndpointUrl,ompitempty"`
	ResourceManagerEndpoint string `json:"resourceManagerEndpointUrl,ompitempty"`
	GraphResourceID         string `json:"activeDirectoryGraphResourceId,ompitempty"`
	SQLManagementEndpoint   string `json:"sqlManagementEndpointUrl,ompitempty"`
	GalleryEndpoint         string `json:"galleryEndpointUrl,ompitempty"`
	ManagementEndpoint      string `json:"managementEndpointUrl,ompitempty"`
}

// GetTokenFromAuthFile creates an authorizer from an Azure CLI auth file
func GetTokenFromAuthFile(baseURI string) (auth Auth, err error) {
	fileLocation := os.Getenv("AZURE_AUTH_LOCATION")
	if fileLocation == "" {
		return auth, errors.New("auth file not found. Environment variable AZURE_AUTH_LOCATION is not set")
	}

	contents, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		return
	}

	// Auth file might be encoded
	decoded, err := decode(contents)
	if err != nil {
		return
	}

	err = json.Unmarshal(decoded, &auth.File)
	if err != nil {
		return
	}

	resource, err := getResource(auth.File, baseURI)
	if err != nil {
		return
	}
	auth.BaseURI = resource

	config, err := adal.NewOAuthConfig(auth.ActiveDirectoryEndpoint, auth.TenantID)
	if err != nil {
		return
	}

	spToken, err := adal.NewServicePrincipalToken(*config, auth.ClientID, auth.ClientSecret, resource)
	if err != nil {
		return
	}

	auth.Authorizer = autorest.NewBearerAuthorizer(spToken)
	return
}

func decode(b []byte) ([]byte, error) {
	reader, enc := utfbom.Skip(bytes.NewReader(b))

	switch enc {
	case utfbom.UTF16LittleEndian:
		u16 := make([]uint16, (len(b)/2)-1)
		err := binary.Read(reader, binary.LittleEndian, &u16)
		if err != nil {
			return nil, err
		}
		return []byte(string(utf16.Decode(u16))), nil
	case utfbom.UTF16BigEndian:
		u16 := make([]uint16, (len(b)/2)-1)
		err := binary.Read(reader, binary.BigEndian, &u16)
		if err != nil {
			return nil, err
		}
		return []byte(string(utf16.Decode(u16))), nil
	}
	return ioutil.ReadAll(reader)
}

func getResource(f File, baseURI string) (string, error) {
	// Compare dafault base URI from the SDK to the endpounts from the public cloud
	if !strings.HasSuffix(baseURI, "/") {
		baseURI += "/"
	}
	switch baseURI {
	case azure.PublicCloud.ServiceManagementEndpoint:
		return f.ManagementEndpoint, nil
	case azure.PublicCloud.ResourceManagerEndpoint:
		return f.ResourceManagerEndpoint, nil
	case azure.PublicCloud.ActiveDirectoryEndpoint:
		return f.ActiveDirectoryEndpoint, nil
	case azure.PublicCloud.GalleryEndpoint:
		return f.GalleryEndpoint, nil
	case azure.PublicCloud.GraphEndpoint:
		return f.GraphResourceID, nil
	}
	return "", fmt.Errorf("auth: base URI not found in endpoints")
}
