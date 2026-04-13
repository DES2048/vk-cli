package video

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/SevereCloud/vksdk/v3/object"
)

func UploadVideo(client *api.VK, filename string, groupID int, albumID int, videoTitle string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening videofile: %w", err)
	}

	resp, err := client.UploadVideo(api.Params{
		"name":     videoTitle,
		"group_id": groupID,
		"album_id": albumID,
	}, file)
	if err != nil {
		return fmt.Errorf("failed to upload video: %w", err)
	}

	m, _ := json.MarshalIndent(resp, "", "  ")
	slog.Debug("upload response", "response", string(m))
	return nil
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

	allCount := resp.Count
	fetched := len(resp.Items)
	items := resp.Items

	for allCount > fetched {
		resp, err := client.VideoGet(api.Params{
			"owner_id": ownerID,
			"album_id": albumID,
			"offset":   fetched,
		})
		if err != nil {
			return nil, err
		}

		fetched += len(resp.Items)
		items = append(items, resp.Items...)

	}
	slog.Debug("videos list", "count", resp.Count, "items", len(items))
	return items, err
}
