package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/boltdb/bolt"
	"github.com/mikespook/gorbac"
	"gitlab.com/middlefront/workspace/storage"
	"gitlab.com/middlefront/workspace/storage/file"
	"gitlab.com/middlefront/workspace/storage/s3"
)

type Config struct {
	FormsMetadata       string
	WorkspacesContainer string
	WorkspacesMetadata  string
	UsersBucket         string
	Auth0ApiToken       string
	Auth0ClientSecret   string

	RootDirectory     string
	BoltFile          string
	SubmissionsBucket []byte
	DB                *bolt.DB
	FileManager       storage.FileManager
	RolesManager      *gorbac.RBAC

	PersistenceType    string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	AWSEndpoint        string
	AWSS3BucketName    string
}

var (
	config Config
)

//Using Init not init, so i can manually determine when the content of config are initalized, as opposed to initializing whenever the package is imported (initialization should happen at app startup, which is only when imported by the main.go file).
func Init(c Config) {
	config = c

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
	case "local":
		config.FileManager = file.Persister{
			RootDirectory: config.RootDirectory,
		}
		break
	default:
		log.Fatal("unknown storage Type: " + config.PersistenceType)
	}

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
