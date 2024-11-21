package main

import (
	"flag"
	"log"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api"
)

func main() {
	albumID := 0
	addAlbumTitle := ""
	configFilename := ""

	// flags
	flag.IntVar(&albumID, "album", 0, "album id")
	flag.StringVar(&addAlbumTitle, "add-album", "", "title for new album")
	flag.StringVar(&configFilename, "config", "config.toml", "path to config file")

	flag.Parse()

	// get config
	config, err := ReadConfig(configFilename)
	if err != nil {
		log.Fatalf("Failed to load config: %s\n", err)
	}

	// get token
	token, err := CheckAuth(config.TokenFile)
	if err != nil {
		log.Fatalf("Failed to get token:%s\n", err)
	}

	vk := api.NewVK(token)

	// add album if any
	if addAlbumTitle != "" {
		var err error

		album, err := getVideoAlbumByTitle(vk, -config.DefaultGroupID, addAlbumTitle)
		if err != nil {
			log.Fatalf("failed to get albums: %s\n", err)
		}

		if album == nil {
			albumID, err = addVideoAlbum(vk, config.DefaultGroupID, addAlbumTitle)
			if err != nil {
				log.Fatalf("Failed to create album: %s\n", err)
			}

			slog.Info("Added new album", "id", albumID, "title", addAlbumTitle)
		} else {
			albumID = album.ID
		}
	}

	// upload videos
	filenames := flag.Args()

	if len(filenames) == 0 {
		slog.Info("No files to upload")
		return
	}

	for idx, filename := range filenames {
		videoName := filepath.Base(filename)
		videoTitle, _ := strings.CutSuffix(videoName, filepath.Ext(videoName))

		slog.Info("upload video", "index", idx+1, "of", len(filenames), "file", filename, "title", videoTitle)

		uploadVideo(vk, filename, config.DefaultGroupID, albumID, videoTitle)
	}
}
