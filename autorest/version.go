package autorest

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
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
	var once sync.Once
	if version == "" {
		once.Do(func() {
			semver := fmt.Sprintf("%d.%d.%d", major, minor, patch)
			verBuilder := bytes.NewBufferString(semver)
			if tag != "" && tag != "-" {
				updated := strings.TrimPrefix(tag, "-")
				_, err := verBuilder.WriteString("-" + updated)
				if err == nil {
					verBuilder = bytes.NewBufferString(semver)
				}
			}
			version = verBuilder.String()
		})
	}
	return version
}
