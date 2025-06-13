package cmd

import (
	"log"
	"vk-cli/config"

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
	cfg, err := config.ReadConfig(ConfigFile)
	if err != nil {
		log.Fatalf("Failed to load config: %s\n", err)
	}
	config.SetConfig(cfg)
}
