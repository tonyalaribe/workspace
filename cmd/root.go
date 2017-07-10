package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	conf "gitlab.com/middlefront/workspace/config"
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
	config.FormsMetadata = viper.GetString("forms-metadata")
	config.WorkspacesContainer = viper.GetString("workspaces-container")
	config.WorkspacesMetadata = viper.GetString("workspaces-metadata")
	config.UsersBucket = viper.GetString("users-bucket")
	config.RootDirectory = viper.GetString("root-directory")
	config.PersistenceType = viper.GetString("persistence-type")

	config.Auth0ApiToken = viper.GetString("auth0-api-token")
	config.Auth0ClientSecret = viper.GetString("auth0-client-secret")

	config.AWSAccessKeyID = viper.GetString("AWS_ACCESS_KEY_ID")
	config.AWSSecretAccessKey = viper.GetString("AWS_SECRET_ACCESS_KEY")
	config.AWSRegion = viper.GetString("AWS_REGION")
	config.AWSEndpoint = viper.GetString("AWS_ENDPOINT")
	config.AWSS3BucketName = viper.GetString("AWS_S3_BUCKET")

	conf.Init(config)

}
