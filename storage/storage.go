package storage

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
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

var stowBucket stow.Container

func StorageInit() {
	conf := config.Get()
	switch conf.PersistenceType {
	case "gcs": //google cloud storage
		b, err := ioutil.ReadFile(conf.GCS.ConfigJSON)
		if err != nil {
			log.Fatal(err)
		}

		stowLoc, err := stow.Dial(stowgcs.Kind, stow.ConfigMap{
			stowgcs.ConfigJSON: string(b),
			// stowgcs.ConfigProjectId: "past-3",
		})
		if err != nil {
			log.Fatal(err)
		}

		stowBucket, err = stowLoc.Container(conf.GCS.BucketName)
		if err != nil {
			log.Fatal(err)
		}

		break
	case "s3":
		stowLoc, err := stow.Dial(stows3.Kind, stow.ConfigMap{
			stows3.ConfigAccessKeyID: conf.S3.AccessKeyID,
			stows3.ConfigSecretKey:   conf.S3.SecretKey,
			stows3.ConfigEndpoint:    conf.S3.Endpoint,
			stows3.ConfigRegion:      conf.S3.Region,
			stows3.ConfigDisableSSL:  conf.S3.DisableSSL,
		})
		if err != nil {
			log.Fatal(err)
		}

		stowBucket, err = stowLoc.Container(conf.S3.BucketName)
		if err != nil {
			log.Fatal(err)
		}

		break
	case "openstack":
		stowLoc, err := stow.Dial(stowswift.Kind, stow.ConfigMap{
			stowswift.ConfigKey:           conf.Openstack.TenantID,
			stowswift.ConfigTenantName:    conf.Openstack.TenantName,
			stowswift.ConfigUsername:      conf.Openstack.Username,
			stowswift.ConfigTenantAuthURL: conf.Openstack.IdentityEndpoint,
		})
		if err != nil {
			log.Fatal(err)
		}

		stowBucket, err = stowLoc.Container(conf.Openstack.BucketName)
		if err != nil {
			log.Fatal(err)
		}

		break
	case "local":
		stowLoc, err := stow.Dial(stowlocal.Kind, stow.ConfigMap{
			stowlocal.ConfigKeyPath: conf.FileSystem.Path,
		})
		if err != nil {
			log.Fatal(err)
		}

		stowBucket, err = stowLoc.Container(conf.FileSystem.BucketName)
		if err != nil {
			log.Fatal(err)
		}

		break
	default:
		log.Fatal("unknown storage Type: " + conf.PersistenceType)
	}

}

type StowFile struct {
	ID   string
	Name string
	URL  string
	Size int64
}

func UploadStream(name string, r io.Reader, size int64) (StowFile, error) {
	file := StowFile{}
	log.Println(stowBucket.Name())
	log.Println(stowBucket.ID())
	resp, err := stowBucket.Put(name, r, size, nil)
	if err != nil {
		return file, err
	}
	log.Printf("%#v", resp)

	file.Name = resp.Name()
	file.URL = resp.URL().String()
	file.ID = resp.ID()
	file.Size, err = resp.Size()
	if err != nil {
		log.Println(err)
	}

	log.Println(file)

	return file, nil
}

func UploadBase64(path string, b64 string) (StowFile, error) {
	file := StowFile{}
	if len(strings.Split(b64, "base64,")) < 2 {
		file.URL = b64
		return file, nil
	}
	byt, err := base64.StdEncoding.DecodeString(strings.Split(b64, "base64,")[1])
	if err != nil {
		log.Println(err)
		file.URL = b64
		return file, err
	}
	meta := strings.Split(b64, "base64,")[0]
	filename := strings.Replace(strings.Split(meta, "name=")[1], ";", "", -1)
	pathToFile := filepath.Join(path, filename)

	fileReader := bytes.NewReader(byt)
	return UploadStream(pathToFile, fileReader, fileReader.Size())

}

func GetFromStorage(file StowFile) (stow.Item, error) {
	// log.Println(stowBucket.Name())
	log.Println(stowBucket.ID())
	item, err := stowBucket.Item(file.ID)
	if err != nil {
		log.Println(err)
	}

	// log.Printf("%#v", resp)
	return item, nil
}
