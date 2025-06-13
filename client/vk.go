package client

import (
	"fmt"
	"vk-cli/auth"
	"vk-cli/config"

	"github.com/SevereCloud/vksdk/v2/api"
)

func BuildVkClient() (*api.VK, error) {
	config := config.GetConfig()
	token, err := auth.CheckAuth(config.TokenFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	vk := api.NewVK(token)
	return vk, nil
}
