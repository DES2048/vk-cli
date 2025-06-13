package group

import (
	"errors"
	"slices"
	"strconv"
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

func GetGroupId(client *api.VK, idOrName string) (int, error) {
	// try convert to int
	groupId, err := strconv.Atoi(idOrName)
	if err != nil {
		// try to get group id by name
		group, err := GetGroupByName(client, idOrName)
		if err != nil {
			return -1, err
		}

		if group == nil {
			return -1, errors.New("group with given name not found")
		}
		return group.ID, nil
	} else {
		return groupId, nil
	}
}
