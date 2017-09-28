package config

import (
	"log"
	"time"

	"github.com/mikespook/gorbac"
	"gitlab.com/middlefront/workspace/database"
	"gitlab.com/middlefront/workspace/database/boltdb"
)

type Openstack struct {
	IdentityEndpoint string
	Username         string
	Password         string
	TenantID         string
	TenantName       string
	BucketName       string
	ApiKey           string
}
type S3 struct {
	Endpoint    string
	AccessKeyID string
	SecretKey   string
	Region      string
	DisableSSL  string
	BucketName  string
}
type GCS struct {
	ConfigJSON string
	BucketName string
}
type FileSystem struct {
	Path       string
	BucketName string
}

type Config struct {
	AppMetadata         string
	FormsMetadata       string
	WorkspacesContainer string
	WorkspacesMetadata  string
	UsersBucket         string
	Auth0ApiToken       string
	Auth0ClientSecret   string

	RootDirectory     string
	BoltFile          string
	SubmissionsBucket []byte
	// FileManager       storage.FileManager
	RolesManager *gorbac.RBAC

	DatabaseType string
	Database     database.Database

	PersistenceType string

	Openstack  Openstack
	S3         S3
	GCS        GCS
	FileSystem FileSystem
}

var (
	config Config
)

//Using Init not init, so i can manually determine when the content of config are initalized, as opposed to initializing whenever the package is imported (initialization should happen at app startup, which is only when imported by the main.go file).
func Init(c Config) {
	config = c
	config.SubmissionsBucket = []byte("submissions")

	switch config.DatabaseType {
	case "boltdb":
		db, err := boltdb.New(config.RootDirectory,
			config.AppMetadata,
			config.WorkspacesMetadata, config.WorkspacesContainer, config.UsersBucket, config.FormsMetadata)
		if err != nil {
			log.Println(err)
		}
		config.Database = database.Database(db)
		break
	default:
		log.Fatal("unknown database type: " + config.DatabaseType)
	}

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
