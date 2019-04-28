package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var configFile string
var verbose bool

// TODO Change me
const cmdName = "http-server"

var (
	Version  string = ""
	Revision string = ""
	Time     string = ""
)

var rootCmd = &cobra.Command{
	Use:   cmdName,
	Short: "",
	Long:  ``,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	//rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", fmt.Sprintf("config file (default \"./.%s.yml\")", cmdName))
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Set Default Viper configs.
	viper.SetDefault("build.version", Version)
	viper.SetDefault("build.revision", Revision)
	viper.SetDefault("build.time", Time)
}

func initConfig() {}
