package s3

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Persister is an implementation of the File Persister
type Persister struct {
	AWSSession *session.Session
	BucketName string
}

// Save persists data to a loader.
func (fp Persister) Save(name string, workspace string, b64Data string) (string, error) {
	pathToSubmission := filepath.Join(workspace, name)
	err := os.MkdirAll(pathToSubmission, os.ModePerm)
	if err != nil {
		return "", err
	}
	fullPath, err := Base64ToS3(b64Data, pathToSubmission, fp.AWSSession, fp.BucketName)
	return fullPath, err
}

//Base64ToS3 puts file into s3 bucket
func Base64ToS3(b64 string, imagename string, awsSession *session.Session, bucketName string) (string, error) {
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

	uploader := s3manager.NewUploader(awsSession)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(pathToFile),
		ACL:         aws.String("public-read"),
		Body:        bytes.NewReader(byt),
		ContentType: aws.String(meta),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", result.Location)

	return result.Location, nil
}
