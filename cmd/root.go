package cmd

import (
	"fmt"
	"os"
	"path"

	"beryju.org/imagik/pkg/config"
	"beryju.org/imagik/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "imagik",
	Short: "A small Fileserver",
	Run: func(cmd *cobra.Command, args []string) {
		os.MkdirAll(path.Join(os.TempDir(), "imagik/"), 0750)
		server := server.New()
		go server.HashMap.RunIndexer()

		err := server.Run()
		if err != nil {
			log.Error(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(onInit)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is config.yaml)")
}

// onInit reads in config file and ENV variables if set.
func onInit() {
	log.SetFormatter(&log.JSONFormatter{})
	config.DefaultConfig()
	err := config.LoadConfig(cfgFile)

	if err == nil {
		log.WithField("config-file", cfgFile).Info("Using config file")
	} else {
		log.WithField("config-file", cfgFile).WithError(err).Warning("Failed to read config")
	}

	if config.C.LogFormat != "json" {
		log.SetFormatter(&log.TextFormatter{})
	}
	if config.C.Debug {
		log.SetLevel(log.TraceLevel)
	}
}
