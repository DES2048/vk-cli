package photos

import (
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
)

func UploadPhoto(client *api.VK, filename string, groupID int, albumID int) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %s\n", err)
	}

	resp, err := client.UploadPhotoGroup(groupID, albumID, file)
	if err != nil {
		slog.Error("Failed to upload photo", "file", filename, "error", err)
	}

	m, _ := json.MarshalIndent(resp, "", "  ")
	slog.Debug("upload response", "response", string(m))
}

func GetPhotoAlbumByTitle(client *api.VK, ownerID int, title string) (*object.PhotosPhotoAlbumFull, error) {
	resp, err := client.PhotosGetAlbums(api.Params{
		"owner_id": ownerID,
	})
	if err != nil {
		return nil, err
	}

	idx := slices.IndexFunc(resp.Items, func(e object.PhotosPhotoAlbumFull) bool {
		return strings.EqualFold(e.Title, title)
	})

	if idx < 0 {
		return nil, nil
	}

	return &resp.Items[idx], nil
}

func AddPhotoAlbum(client *api.VK, groupID int, title string) (int, error) {
	resp, err := client.PhotosCreateAlbum(api.Params{
		"group_id": groupID,
		"title":    title,
	})
	if err != nil {
		return -1, err
	}

	return resp.ID, nil
}
