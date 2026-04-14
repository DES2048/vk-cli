package cmd

import (
	"fmt"
	"log"
	"vk-cli/config"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

type SizeFlagValue struct {
	Value uint64
	IsGt  bool
}

func (v *SizeFlagValue) String() string {
	if v.Value != 0 {
		sign := "<"
		if v.IsGt {
			sign = ">"
		}
		return fmt.Sprintf("%s%d", sign, v.Value)
	}
	return ""
}

func (v *SizeFlagValue) Set(value string) error {
	sign := value[0]

	if sign == byte('<') {
		v.IsGt = false
	} else if sign == byte('>') {
		v.IsGt = true
	} else {
		return fmt.Errorf("invalid cmp op, should be <10mb, or >10mb")
	}
	var err error
	v.Value, err = humanize.ParseBytes(value[1:])
	return err
}

func (v *SizeFlagValue) Type() string {
	return "size"
}

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
