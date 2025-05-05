package cmd

import (
	"log"
	"vk-cli/auth"
	"vk-cli/config"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use: "auth",
	Run: func(cmd *cobra.Command, args []string) {
		// get config
		config, err := config.ReadConfig(ConfigFile)
		if err != nil {
			log.Fatalf("Failed to load config: %s\n", err)
		}

		token, err := auth.Auth()
		if err != nil {
			log.Fatalf("failed to obtain token: %s", err)
		}

		err = auth.WriteTokenToFile(token, config.TokenFile)
		if err != nil {
			log.Fatalf("failed to save token to file: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(authCmd)
}
