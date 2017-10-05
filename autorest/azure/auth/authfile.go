package auth

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"unicode/utf16"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/dimchansky/utfbom"
)

var pcEndpoints = map[string]string{
	"managementEndpointUrl":          "https://management.core.windows.net",
	"resourceManagerEndpointUrl":     "https://management.azure.com",
	"activeDirectoryEndpointUrl":     "https://login.microsoftonline.com",
	"galleryEndpointUrl":             "https://gallery.azure.com",
	"activeDirectoryGraphResourceId": "https://graph.windows.net",
}

// Auth includes authentication details for ARM clients
type Auth struct {
	Authorizer *autorest.BearerAuthorizer
	File       map[string]string
	BaseURI    string
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

	var af map[string]string
	err = json.Unmarshal(decoded, &af)
	if err != nil {
		return
	}
	auth.File = af

	k := getResourceKey(baseURI)
	auth.BaseURI = af[k]

	config, err := adal.NewOAuthConfig(af["activeDirectoryEndpointUrl"], af["tenantId"])
	if err != nil {
		return
	}

	spToken, err := adal.NewServicePrincipalToken(*config, af["clientId"], af["clientSecret"], af[k])
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

func getResourceKey(baseURI string) string {
	for k, v := range pcEndpoints {
		if baseURI == v {
			return k
		}
	}
	return ""
}
