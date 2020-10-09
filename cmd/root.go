package cmd

import (
	"fmt"
	"os"

	"github.com/BeryJu/gopyazo/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopyazo",
	Short: "A small Fileserver",
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
	log.SetLevel(log.DebugLevel)
	config.DefaultConfig()

	configPath := "./config.yml"

	if err := config.LoadConfig(configPath); err == nil {
		log.WithField("config-file", configPath).Info("Using config file")
	}

	if viper.GetString(config.C.LogFormat) == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}
}
