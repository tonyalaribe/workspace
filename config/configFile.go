package config

import (
	"fmt"
	"log"

	"gitlab.com/middlefront/workspace/filePersistence"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/viper"
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetConfigName(".workspace") // name of config file (without extension)
	viper.AddConfigPath(".")          // The apps root root directory as first search path
	viper.AddConfigPath("$HOME")      // adding home directory as second search path
	// viper.AutomaticEnv()              // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	//This initialization is placed here, because initConfig is a callback that is called after cobra has parsed the config file and other variables. The other ideal location would have been the init function, but the init function is called before the config has been parsed, and hence the absense of the needed variables.

	config.FormsMetadata = viper.GetString("forms-metadata")
	config.WorkspacesContainer = viper.GetString("workspaces-container")
	config.WorkspacesMetadata = viper.GetString("workspaces-metadata")
	config.UsersBucket = viper.GetString("users-bucket")
	config.Auth0ApiToken = viper.GetString("auth0-api-token")
	config.RootDirectory = viper.GetString("root-directory")
	persistenceType := viper.GetString("persistence-type")

	switch persistenceType {
	case "s3":
		// creds := credentials.NewEnvCredentials()
		creds := credentials.NewStaticCredentials(viper.GetString("AWS_ACCESS_KEY_ID"), viper.GetString("AWS_SECRET_ACCESS_KEY"), "")
		credValue, err := creds.Get()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%#v", credValue)
		awsConfig := &aws.Config{
			Credentials: creds,
			Region:      aws.String(viper.GetString("AWS_REGION")),
		}
		endpoint := viper.GetString("AWS_ENDPOINT")
		if endpoint != "" {
			awsConfig.Endpoint = aws.String(endpoint)
			awsConfig.DisableSSL = aws.Bool(true)
			awsConfig.S3ForcePathStyle = aws.Bool(true)
		}
		sess := session.New(awsConfig)
		config.FileManager = filePersistence.S3Persister{
			AWSSession: sess,
			BucketName: viper.GetString("AWS_S3_BUCKET"),
		}
		break
	case "local":
		config.FileManager = filePersistence.FilePersister{RootDirectory: config.RootDirectory}
		break
	default:
		log.Fatal("unknown storage Type: " + persistenceType)
	}
}
