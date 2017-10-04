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
	utf16leBOM := []byte{255, 254}
	utf16beBOM := []byte{254, 255}
	utf8BOM := []byte{239, 187, 191}

	switch {
	case bytes.HasPrefix(b, utf16leBOM):
		b = bytes.TrimPrefix(b, utf16leBOM)
		u16 := make([]uint16, (len(b) / 2))
		buf := bytes.NewReader(b)
		err := binary.Read(buf, binary.LittleEndian, &u16)
		if err != nil {
			return nil, err
		}
		return []byte(string(utf16.Decode(u16))), nil
	case bytes.HasPrefix(b, utf16beBOM):
		b = bytes.TrimPrefix(b, utf16beBOM)
		u16 := make([]uint16, (len(b) / 2))
		buf := bytes.NewReader(b)
		err := binary.Read(buf, binary.BigEndian, &u16)
		if err != nil {
			return nil, err
		}
		return []byte(string(utf16.Decode(u16))), nil
	case bytes.HasPrefix(b, utf8BOM):
		return bytes.TrimPrefix(b, utf8BOM), nil
	}
	return b, nil
}

func getResourceKey(baseURI string) string {
	for k, v := range pcEndpoints {
		if baseURI == v {
			return k
		}
	}
	return ""
}
