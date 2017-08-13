package config

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mikespook/gorbac"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"gitlab.com/middlefront/workspace/database"
	"gitlab.com/middlefront/workspace/database/boltdb"
	"gitlab.com/middlefront/workspace/storage"
	"gitlab.com/middlefront/workspace/storage/file"
	"gitlab.com/middlefront/workspace/storage/s3"
)

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
	FileManager       storage.FileManager
	RolesManager      *gorbac.RBAC

	PersistenceType    string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	AWSEndpoint        string
	AWSS3BucketName    string

	DatabaseType string
	Database     database.Database
}

var (
	config Config
)

//Using Init not init, so i can manually determine when the content of config are initalized, as opposed to initializing whenever the package is imported (initialization should happen at app startup, which is only when imported by the main.go file).
func Init(c Config) {
	config = c
	config.SubmissionsBucket = []byte("submissions")
	switch config.PersistenceType {
	case "s3":
		// creds := credentials.NewEnvCredentials()
		creds := credentials.NewStaticCredentials(config.AWSAccessKeyID, config.AWSSecretAccessKey, "")
		credValue, err := creds.Get()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%#v", credValue)
		awsConfig := &aws.Config{
			Credentials: creds,
			Region:      aws.String(config.AWSRegion),
		}
		endpoint := config.AWSEndpoint
		if endpoint != "" {
			awsConfig.Endpoint = aws.String(endpoint)
			awsConfig.DisableSSL = aws.Bool(true)
			awsConfig.S3ForcePathStyle = aws.Bool(true)
		}
		sess := session.New(awsConfig)
		config.FileManager = s3.Persister{
			AWSSession: sess,
			BucketName: config.AWSS3BucketName,
		}
		break
	case "openstack":
		// creds := credentials.NewEnvCredentials()
		opts := gophercloud.AuthOptions{
			IdentityEndpoint: "https://storage.bhs1.cloud.ovh.net/v2.0/",
			Username:         "nGaW4ryDxkz3",
			Password:         "GmEhvEHWMzwheeJB6hRz8q283HNubYC2",
			TenantID:         "AUTH_74439eb9176b44c78f1a279cb21f554d",
		}
		provider, err := openstack.AuthenticatedClient(opts)
		if err != nil {
			log.Println(err)
		}
		client, err := openstack.NewObjectStorageV1(provider, gophercloud.EndpointOpts{})
		if err != nil {
			log.Println(err)
		}
		log.Print(client.ResourceBaseURL())
		break
	case "local":
		config.FileManager = file.Persister{
			RootDirectory: config.RootDirectory,
		}
		break
	default:
		log.Fatal("unknown storage Type: " + config.PersistenceType)
	}

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
