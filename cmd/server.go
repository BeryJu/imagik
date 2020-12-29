package cmd

import (
	"github.com/BeryJu/imagik/pkg/hash"
	"github.com/BeryJu/imagik/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run imagik Server",
	Run: func(cmd *cobra.Command, args []string) {
		server := server.New()
		server.HashMap = hash.New()
		go server.HashMap.RunIndexer()

		err := server.Run()
		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
