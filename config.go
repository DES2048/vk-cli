package main

import "github.com/ilyakaznacheev/cleanenv"

// Config ...
type Config struct {
	DefaultGroupID int    `yaml:"defaultGroupId" env:"DEFAULT_GROUP_ID, env-required"`
	TokenFile      string `yaml:"tokenFile" env:"TOKEN_FILE" env-default:"token.txt"`
}

// ReadConfig ...
func ReadConfig(file string) (*Config, error) {
	var config Config
	err := cleanenv.ReadConfig(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
