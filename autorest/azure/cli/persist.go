package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dimchansky/utfbom"
)

// LoadCLIProfile restores an AzureCLIProfile object from a file located at 'path'.
func LoadCLIProfile(path string) (AzureCLIProfile, error) {
	var profile AzureCLIProfile

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return profile, fmt.Errorf("failed to open file (%s) while loading token: %v", path, err)
	}
	reader := utfbom.SkipOnly(bytes.NewReader(contents))

	dec := json.NewDecoder(reader)
	if err = dec.Decode(&profile); err != nil {
		return profile, fmt.Errorf("failed to decode contents of file (%s) into a AzureCLIProfile representation: %v", path, err)
	}

	return profile, nil
}

// LoadCLITokens restores a set of Token objects from a file located at 'path'.
func LoadCLITokens(path string) ([]Token, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file (%s) while loading token: %v", path, err)
	}
	defer file.Close()

	var tokens []Token

	dec := json.NewDecoder(file)
	if err = dec.Decode(&tokens); err != nil {
		return nil, fmt.Errorf("failed to decode contents of file (%s) into a `cli.Token` representation: %v", path, err)
	}

	return tokens, nil
}
