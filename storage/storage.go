/*
Package storage makes it possible to perform file related operation from the rest of workspace without worrying about the underlying storage implementation.

The storage mechanism to use, and the relevant connection details should be specified in the workspace.yaml file.

.workspace.yaml

The .workspcae.yaml file coul have the following storage related files:

- persistence-type: This represents the storage backend to be used. eg local, s3, openstack,

*/
package storage

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/graymeta/stow"
	// support Azure storage
	_ "github.com/graymeta/stow/azure"
	// support Google storage
	stowgcs "github.com/graymeta/stow/google"
	// support local storage
	stowlocal "github.com/graymeta/stow/local"
	// support swift storage
	stowswift "github.com/graymeta/stow/swift"
	// support s3 storage
	stows3 "github.com/graymeta/stow/s3"
	// support oracle storage

	_ "github.com/graymeta/stow/oracle"
	"gitlab.com/middlefront/workspace/config"
)

type StowFile struct {
	ID   string
	Name string
	URL  string
	Size int64
}

var stowBucket stow.Container
var stowLoc stow.Location

//StorageInit initilizes stow with the storage details from the .workspace.yaml configuration. StorageInit creates a stowBucket and stowLoc instance, which can be shared during the lifetime of the application
func StorageInit() error {
	conf := config.Get()
	var err error
	switch conf.PersistenceType {
	case "gcs": //google cloud storage
		b, err := ioutil.ReadFile(conf.GCS.ConfigJSON)
		if err != nil {
			return err
		}

		stowLoc, err = stow.Dial(stowgcs.Kind, stow.ConfigMap{
			stowgcs.ConfigJSON: string(b),
			// stowgcs.ConfigProjectId: "past-3",
		})
		if err != nil {
			return err
		}

		stowBucket, err = stowLoc.Container(conf.GCS.BucketName)
		if err != nil {
			return err
		}
		break

	case "s3":
		stowLoc, err = stow.Dial(stows3.Kind, stow.ConfigMap{
			stows3.ConfigAccessKeyID: conf.S3.AccessKeyID,
			stows3.ConfigSecretKey:   conf.S3.SecretKey,
			stows3.ConfigEndpoint:    conf.S3.Endpoint,
			stows3.ConfigRegion:      conf.S3.Region,
			stows3.ConfigDisableSSL:  conf.S3.DisableSSL,
		})
		if err != nil {
			return err
		}

		stowBucket, err = stowLoc.Container(conf.S3.BucketName)
		if err != nil {
			return err
		}
		break

	case "openstack":
		stowLoc, err = stow.Dial(stowswift.Kind, stow.ConfigMap{
			stowswift.ConfigKey:           conf.Openstack.TenantID,
			stowswift.ConfigTenantName:    conf.Openstack.TenantName,
			stowswift.ConfigUsername:      conf.Openstack.Username,
			stowswift.ConfigTenantAuthURL: conf.Openstack.IdentityEndpoint,
		})
		if err != nil {
			return err
		}

		stowBucket, err = stowLoc.Container(conf.Openstack.BucketName)
		if err != nil {
			return err
		}
		break

	case "local":
		stowLoc, err = stow.Dial(stowlocal.Kind, stow.ConfigMap{
			stowlocal.ConfigKeyPath: conf.FileSystem.Path,
		})
		if err != nil {
			return err
		}

		stowBucket, err = stowLoc.Container(conf.FileSystem.BucketName)
		if err != nil {
			return err
		}
		break

	default:
		return errors.New("unknown storage Type: " + conf.PersistenceType)
	}

	return nil
}

//UploadStream recieves an io.reader which is used stream the file to the repective storage backend. This is ideal, as it prevents sitting the file(bytes) in memory during uploads.
//Path represents they path to the file, which might be a single filename, or follow a backslack delimited structure.
func UploadStream(path string, r io.Reader, size int64) (StowFile, error) {
	file := StowFile{}
	resp, err := stowBucket.Put(path, r, size, nil)
	if err != nil {
		return file, err
	}

	file.Name = resp.Name()
	file.URL = resp.URL().String()
	file.ID = resp.ID()
	file.Size, err = resp.Size()
	if err != nil {
		return file, err
	}

	return file, nil
}

//UploadBase64 recieves a base64 string, which it converts to the underlying file and uploads via UploadStream
func UploadBase64(path string, b64 string) (StowFile, error) {
	file := StowFile{}
	if len(strings.Split(b64, "base64,")) < 2 {
		file.URL = b64
		return file, nil
	}
	byt, err := base64.StdEncoding.DecodeString(strings.Split(b64, "base64,")[1])
	if err != nil {
		file.URL = b64
		return file, err
	}

	meta := strings.Split(b64, "base64,")[0]
	filename := strings.Replace(strings.Split(meta, "name=")[1], ";", "", -1)
	pathToFile := filepath.Join(path, filename)
	fileReader := bytes.NewReader(byt)
	return UploadStream(pathToFile, fileReader, fileReader.Size())
}

//GetByPath is able to load a flle from stowBucket given the file ID(which is usually a path to the file on the storage bucket)
func GetByPath(path string) (stow.Item, error) {
	var item stow.Item
	item, err := stowBucket.Item(path)
	if err != nil {
		return item, err
	}
	return item, nil
}

//GetByURL is able to load a file from stowLoc (stow location) which is bucket independent, given a url which usually encodes information like the storage backend kind, and the bucket the files exist in.
func GetByURL(urlstr string) (stow.Item, error) {
	var item stow.Item
	stowURL, err := url.Parse(urlstr)
	if err != nil {
		return item, err
	}
	item, err = stowLoc.ItemByURL(stowURL)
	if err != nil {
		return item, err
	}
	return item, nil
}
