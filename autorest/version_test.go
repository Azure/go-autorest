package autorest

import (
	"io"
	"os/exec"
	"testing"

	"fmt"

	"sync"

	"github.com/Masterminds/semver"
)

func TestVersion(t *testing.T) {
	var declaredVersion *semver.Version
	if temp, err := semver.NewVersion(Version()); nil == err {
		declaredVersion = temp
		t.Logf("Declared Version: %s", declaredVersion.String())
	} else {
		t.Error(err)
	}

	var currentVersion *semver.Version
	if temp, err := getMaxReleasedVersion(); nil == err {
		currentVersion = temp
		t.Logf("Current Release Version: %s", currentVersion.String())
	} else {
		t.Error(err)
	}

	if !declaredVersion.GreaterThan(currentVersion) {
		t.Log("autorest: Assertion that the Declared version is greater than Current Release Version failed", currentVersion.String(), declaredVersion.String())
		t.Fail()
	}
}

// Variables required by getMaxReleasedVersion. None of these variables should be used outside of that
// function.
var (
	maxReleasedVersion *semver.Version
)

func getMaxReleasedVersion() (*semver.Version, error) {
	if nil == maxReleasedVersion {
		var wg sync.WaitGroup
		var currentTag string
		var emptyVersion semver.Version
		reader, writer := io.Pipe()
		tagLister := exec.Command("git", "tag")
		tagLister.Stdout = writer

		if err := tagLister.Start(); nil != err {
			return nil, err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			maxReleasedVersion = &emptyVersion
			for {
				if parity, err := fmt.Fscanln(reader, &currentTag); err != nil || parity != 1 {
					break
				}

				if currentVersion, err := semver.NewVersion(currentTag); err == nil && currentVersion.GreaterThan(maxReleasedVersion) {
					maxReleasedVersion = currentVersion
				}
			}
		}()

		if err := tagLister.Wait(); nil != err {
			return nil, err
		}
		writer.Close()
		wg.Wait()
	}
	return maxReleasedVersion, nil
}
