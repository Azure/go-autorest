// +build windows

package adal

import (
	"os"
	"path/filepath"
)

// msiPath is the path to the MSI Extension settings file (to discover the endpoint)
var msiPath = filepath.Join(os.Getenv("SystemDrive"), "WindowsAzure/Config/ManagedIdentity-Settings")
