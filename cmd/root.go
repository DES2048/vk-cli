package cmd

import (
	"github.com/spf13/cobra"
)

var (
	ConfigFile string

	RootCmd = &cobra.Command{
		Use: "vk-cli",
	}
)

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&ConfigFile, "config", "config.toml", "path to config file")
}

func initConfig() {
}
