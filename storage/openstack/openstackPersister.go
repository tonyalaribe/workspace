package openstack

import (
	"bytes"
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
)

// Persister is an implementation of the File Persister
type Persister struct {
	OpenstackSession *gophercloud.ServiceClient
	BucketName       string
	ResourceBaseURL  string
}

// Save persists data to a loader.
func (fp Persister) Save(workspaceID string, formID string, submissionName string, b64Data string) (string, error) {
	pathToSubmission := filepath.Join(workspaceID, formID, submissionName)
	err := os.MkdirAll(pathToSubmission, os.ModePerm)
	if err != nil {
		return "", err
	}
	fullPath, err := Base64ToOpenstack(b64Data, fp.ResourceBaseURL, pathToSubmission, fp.OpenstackSession, fp.BucketName)
	return fullPath, err
}

//Base64ToOpenstack puts file into s3 bucket
func Base64ToOpenstack(b64 string, resourcePath, imagename string, openstackClient *gophercloud.ServiceClient, bucketName string) (string, error) {
	log.Println("in base64ToOpenstack")
	if len(strings.Split(b64, "base64,")) < 2 {
		return b64, nil
	}
	byt, err := base64.StdEncoding.DecodeString(strings.Split(b64, "base64,")[1])
	if err != nil {
		log.Println(err)
		return "", err
	}
	meta := strings.Split(b64, "base64,")[0]
	filename := strings.Replace(strings.Split(meta, "name=")[1], ";", "", -1)
	pathToFile := filepath.Join(imagename, filename)

	// You have the option of specifying additional configuration options.
	opts := objects.CreateOpts{}

	// Now execute the upload
	res := objects.Create(openstackClient, bucketName, pathToFile, bytes.NewReader(byt), opts)
	if res.Err != nil {
		log.Println(res)
	}

	// We have the option of extracting the resulting
	return resourcePath + "/" + pathToFile, nil
}
