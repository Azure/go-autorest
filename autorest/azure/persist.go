package azure

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

// LoadToken restores a Token object from a file located at 'path'.
func LoadToken(path string) (*Token, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file (%s) while loading token: %v", path, err)
	}

	var token Token

	dec := json.NewDecoder(file)
	err = dec.Decode(&token)
	if err != nil {
		return nil, fmt.Errorf("failed to decode contents of file (%s) into Token representation: %v", path, err)
	}

	return &token, nil
}

// SaveToken persists an oauth token at the given location on disk.
// It moves the new file into place so it can safely be used to replace an existing file
// that maybe accessed by multiple processes.
func SaveToken(path string, token Token) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory (%s) to store token in: %v", dir, err)
	}

	tempPath := path + fmt.Sprintf("%d", rand.Int31())

	newFile, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("failed to write token to temp file (%s) while saving token: %v", tempPath, err)
	}

	enc := json.NewEncoder(newFile)
	err = enc.Encode(token)
	if err != nil {
		return fmt.Errorf("failed to encode token to file (%s) while saving token: %v", tempPath, err)
	}

	err = os.Rename(tempPath, path)
	if err != nil {
		return fmt.Errorf("failed to move temporary token to desired output location. source=(%s). destination=(%s). error = %v", tempPath, path, err)
	}

	return nil
}
