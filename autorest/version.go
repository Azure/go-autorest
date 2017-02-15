package autorest

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	major = 7
	minor = 3
	patch = 1
	tag   = ""
)

var version string

// Version returns the semantic version (see http://semver.org).
func Version() string {
	if version == "" {
		verBuilder := bytes.NewBufferString(fmt.Sprintf("%d.%d.%d", major, minor, patch))
		if tag != "" && tag != "-" {
			updated := strings.TrimPrefix(tag, "-")
			verBuilder.WriteRune('-')
			verBuilder.WriteString(updated)
		}
		version = verBuilder.String()
	}
	return version
}
