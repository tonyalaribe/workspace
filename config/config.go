package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/mikespook/gorbac"
	"gitlab.com/middlefront/workspace/storage"
)

type Config struct {
	FormsMetadata       string
	WorkspacesContainer string
	WorkspacesMetadata  string
	UsersBucket         string
	Auth0ApiToken       string

	RootDirectory     string
	BoltFile          string
	SubmissionsBucket []byte
	DB                *bolt.DB
	FileManager       storage.FileManager
	RolesManager      *gorbac.RBAC
}

var (
	config Config
)

//Using Init not init, so i can manually determine when the content of config are initalized, as opposed to initializing whenever the package is imported (initialization should happen at app startup, which is only when imported by the main.go file).
func Init() {
	initConfig()

	config.BoltFile = filepath.Join(config.RootDirectory, "workspace.db")
	config.SubmissionsBucket = []byte("submissions")

	os.MkdirAll(config.RootDirectory, os.ModePerm)
	db, err := bolt.Open(config.BoltFile, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		log.Println(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(config.WorkspacesMetadata))
		tx.CreateBucketIfNotExists([]byte(config.WorkspacesContainer))
		tx.CreateBucketIfNotExists([]byte(config.UsersBucket))
		return nil
	})

	log.Println(db.GoString())
	config.DB = db

	config.RolesManager = GenerateRolesInstance()
	defer SavePermissions()
	go func() {
		for range time.Tick(time.Second * 10) {
			SavePermissions()
		}
	}()
}

func Get() *Config {
	return &config
}
