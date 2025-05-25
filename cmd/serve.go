package cmd

import (
	"dojer/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var port = 8033
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start a small webserver",
	Long:  `start a local webserver to read the downloaded doujinshis`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Init(viper.GetString("server.port"))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
