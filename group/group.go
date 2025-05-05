package group

import (
	"slices"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
)

func GetGroupByName(client *api.VK, name string) (*object.GroupsGroup, error) {
	resp, err := client.GroupsGetExtended(api.Params{
		"filter": "admin",
	})
	if err != nil {
		return nil, err
	}

	idx := slices.IndexFunc(resp.Items, func(e object.GroupsGroup) bool {
		return strings.EqualFold(e.Name, name)
	})

	if idx < 0 {
		return nil, nil
	}

	return &resp.Items[idx], nil
}
