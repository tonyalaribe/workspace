package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	conf "gitlab.com/middlefront/workspace/config"
	"gitlab.com/middlefront/workspace/storage"
)

var (
	cfgFile string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Workspace ---",
	Long: `--------
	`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.workspace.yaml)")
}

func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".workspace") // name of config file (without extension)
	viper.AddConfigPath(".")          // The apps root root directory as first search path
	viper.AddConfigPath("$HOME")      // adding home directory as second search path
	viper.AutomaticEnv()              // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	//This initialization is placed here, because initConfig is a callback that is called after cobra has parsed the config file and other variables. The other ideal location would have been the init function, but the init function is called before the config has been parsed, and hence the absense of the needed variables.

	config := conf.Config{}
	config.AppMetadata = viper.GetString("app-metadata")
	config.FormsMetadata = viper.GetString("forms-metadata")
	config.WorkspacesContainer = viper.GetString("workspaces-container")
	config.WorkspacesMetadata = viper.GetString("workspaces-metadata")
	config.UsersBucket = viper.GetString("users-bucket")
	config.RootDirectory = viper.GetString("root-directory")
	config.PersistenceType = viper.GetString("persistence-type")

	config.Auth0ApiToken = viper.GetString("auth0-api-token")
	config.Auth0ClientSecret = viper.GetString("auth0-client-secret")

	config.DatabaseType = viper.GetString("database-type")

	s3 := viper.GetStringMapString("s3")
	config.S3 = conf.S3{
		Endpoint:    s3["endpoint"],
		AccessKeyID: s3["access-key-id"],
		SecretKey:   s3["secret-key"],
		Region:      s3["region"],
		DisableSSL:  s3["disable-ssl"],
		BucketName:  s3["bucket-name"],
	}

	gcs := viper.GetStringMapString("gcs")
	config.GCS = conf.GCS{
		ConfigJSON: gcs["config-json"],
		// ProjectID:  gcs["project-id"],
	}

	openstack := viper.GetStringMapString("openstack")
	config.Openstack = conf.Openstack{
		IdentityEndpoint: openstack["identity-endpoint"],
		Username:         openstack["username"],
		Password:         openstack["password"],
		TenantID:         openstack["tenant-id"],
		TenantName:       openstack["tenant-name"],
		BucketName:       openstack["bucket-name"],
		ApiKey:           openstack["api-key"],
	}

	filesystem := viper.GetStringMapString("filesystem")
	config.FileSystem = conf.FileSystem{
		Path: filesystem["path"],
	}

	conf.Init(config)
	err := storage.StorageInit()
	if err != nil {
		log.Fatal(err)
	}
}
