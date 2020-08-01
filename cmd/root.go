package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/BeryJu/gopyazo/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
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
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is config.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log.SetLevel(log.DebugLevel)
	log.AddHook(&apmlogrus.Hook{})
	t := apm.DefaultTracer
	t.SetLogger(log.WithField("component", "apm"))

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile("config.yaml")
	}
	viper.SetEnvPrefix("gopyazo")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.WithField("config-file", viper.ConfigFileUsed()).Info("Using config file")
	}

	if viper.GetString(config.ConfigLogFormat) == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}
}
