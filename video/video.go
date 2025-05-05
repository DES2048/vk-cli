package video

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

func UploadVideo(client *api.VK, filename string, groupID int, albumID int, videoTitle string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening videofile: %s\n", err)
	}

	resp, err := client.UploadVideo(api.Params{
		"name":     videoTitle,
		"group_id": groupID,
		"album_id": albumID,
	}, file)
	if err != nil {
		slog.Error("Failed to upload video", "file", filename, "error", err)
	}

	m, _ := json.MarshalIndent(resp, "", "  ")
	slog.Debug("upload response", "response", string(m))
}

func GetVideoAlbumByTitle(client *api.VK, ownerID int, title string) (*object.VideoVideoAlbum, error) {
	resp, err := client.VideoGetAlbums(api.Params{
		"owner_id": ownerID,
	})
	if err != nil {
		return nil, err
	}

	idx := slices.IndexFunc(resp.Items, func(e object.VideoVideoAlbum) bool {
		return strings.EqualFold(e.Title, title)
	})

	if idx < 0 {
		return nil, nil
	}

	return &resp.Items[idx], nil
}

func AddVideoAlbum(client *api.VK, groupID int, title string) (int, error) {
	resp, err := client.VideoAddAlbum(api.Params{
		"group_id": groupID,
		"title":    title,
	})
	if err != nil {
		return -1, err
	}

	return resp.AlbumID, nil
}

func GetVideos(client *api.VK, ownerID int, albumID int) ([]object.VideoVideo, error) {
	resp, err := client.VideoGet(api.Params{
		"owner_id": ownerID,
		"album_id": albumID,
	})
	return resp.Items, err
}
