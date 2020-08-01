package cmd

import (
	"github.com/BeryJu/gopyazo/pkg/config"
	"github.com/BeryJu/gopyazo/pkg/hash"
	"github.com/BeryJu/gopyazo/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run gopyazo Server",
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetDefault(config.ConfigListen, "localhost:8000")
		viper.SetDefault(config.ConfigAuthenticationDriver, "null")
		viper.SetDefault(config.ConfigRootDir, "webroot")

		server := server.New()
		server.HashMap = hash.New()
		go server.HashMap.RunIndexer()

		server.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
