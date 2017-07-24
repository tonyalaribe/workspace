package file

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// FilePersister is an implementation of the File Persister
type Persister struct {
	RootDirectory string
}

// Save persists data to a loader.
func (fp Persister) Save(workspaceID string, formID string, submissionName string, b64Data string) (string, error) {
	pathToSubmission := filepath.Join(fp.RootDirectory, workspaceID, formID, submissionName)
	err := os.MkdirAll(pathToSubmission, os.ModePerm)
	if err != nil {
		return "", err
	}
	fullPath := Base64ToFileSystem(b64Data, pathToSubmission)
	return filepath.Join("/", fullPath), nil
}

func Base64ToFileSystem(b64 string, location string) string {
	if len(strings.Split(b64, "base64,")) < 2 {
		return b64
	}
	byt, err := base64.StdEncoding.DecodeString(strings.Split(b64, "base64,")[1])
	if err != nil {
		log.Println(err)
		//Cant decode string, so its probably an already processed url.
		return b64
	}

	meta := strings.Split(b64, "base64,")[0]
	filename := strings.Replace(strings.Split(meta, "name=")[1], ";", "", -1)

	fullPath := filepath.Join(location, filename)
	err = ioutil.WriteFile(fullPath, byt, 0644)
	if err != nil {
		log.Println(err)
	}
	return fullPath
}
