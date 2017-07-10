package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/middlefront/workspace/web"
)

// serveCmd represents the serve command, which periodically polls the database for updates
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the process to serve workspace",
	Long:  `This command starts a server that serves the workspace project.`,
	Run: func(cmd *cobra.Command, args []string) {
		web.App()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
