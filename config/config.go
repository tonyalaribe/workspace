package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

type Config struct {
	RootDirectory     string
	BoltFile          string
	SubmissionsBucket []byte
	DB                *bolt.DB
}

var (
	config Config
)

const (
	FORMS_METADATA       = "forms_metadata"
	WORKSPACES_METADATA  = "workspaces_metadata"
	WORKSPACES_CONTAINER = "workspaces_container"
)

//Using Init not init, so i can manually determine when the content of config are initalized, as opposed to initializing whenever the package is imported (initialization should happen at app startup, which is only when imported by the main.go file).
func Init() {
	config.RootDirectory = filepath.Join(".", "data")
	config.BoltFile = filepath.Join(config.RootDirectory, "workspace.db")
	config.SubmissionsBucket = []byte("submissions")

	os.MkdirAll(config.RootDirectory, os.ModePerm)
	db, err := bolt.Open(config.BoltFile, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		log.Println(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(WORKSPACES_METADATA))
		tx.CreateBucketIfNotExists([]byte(WORKSPACES_CONTAINER))
		return nil
	})

	log.Println(db.GoString())
	config.DB = db
}

func Get() *Config {
	return &config
}
