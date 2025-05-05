package cmd

import (
	"log"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"
	"vk-cli/auth"
	"vk-cli/config"
	"vk-cli/group"
	"vk-cli/video"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/spf13/cobra"
)

var (
	albumID       = 0
	addAlbumTitle = ""
	groupVar      = ""

	uploadVideoCmd = &cobra.Command{
		Use:   "upload-video",
		Short: " upload video(s) to group",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			groupId := 0

			// get config
			config, err := config.ReadConfig(ConfigFile)
			if err != nil {
				log.Fatalf("Failed to load config: %s\n", err)
			}

			// get token
			token, err := auth.CheckAuth(config.TokenFile)
			if err != nil {
				log.Fatalf("Failed to get token:%s\n", err)
			}

			vk := api.NewVK(token)

			// if group is a group title get group id
			if len(groupVar) > 0 {
				groupId, err = strconv.Atoi(groupVar)
				if err != nil {
					// try to get group id by name
					group, err := group.GetGroupByName(vk, groupVar)
					if err != nil {
						log.Fatalf("failed to get group by name: %s", err)
					}

					if group == nil {
						log.Fatal("group with given name not found")
					}
					groupId = group.ID
				}
			} else {
				groupId = config.DefaultGroupID
			}
			// add album if any
			if addAlbumTitle != "" {
				var err error

				album, err := video.GetVideoAlbumByTitle(vk, -groupId, addAlbumTitle)
				if err != nil {
					log.Fatalf("failed to get albums: %s\n", err)
				}

				if album == nil {
					albumID, err = video.AddVideoAlbum(vk, groupId, addAlbumTitle)
					if err != nil {
						log.Fatalf("Failed to create album: %s\n", err)
					}

					slog.Info("Added new album", "id", albumID, "title", addAlbumTitle)
				} else {
					albumID = album.ID
				}
			}

			// upload videos
			filenames := args

			videos, err := video.GetVideos(vk, -groupId, albumID)
			if err != nil {
				log.Fatalf("failed to get videos list: %s", err)
			}
			// create videos map
			videosTitleMap := make(map[string]bool)

			for _, video := range videos {
				videosTitleMap[video.Title] = true
			}

			for idx, filename := range filenames {
				videoName := filepath.Base(filename)
				videoTitle, _ := strings.CutSuffix(videoName, filepath.Ext(videoName))

				if _, ok := videosTitleMap[videoTitle]; ok {
					slog.Info("skipped video", "index", idx+1, "of", len(filenames), "file", filename, "title", videoTitle)
					continue
				}

				slog.Info("upload video", "index", idx+1, "of", len(filenames), "file", filename, "title", videoTitle)

				video.UploadVideo(vk, filename, groupId, albumID, videoTitle)
			}
		},
	}
)

func init() {
	uploadVideoCmd.Flags().StringVar(&groupVar, "group", "", "group name or group id")
	uploadVideoCmd.Flags().IntVar(&albumID, "album", 0, "album id")
	uploadVideoCmd.Flags().StringVar(&addAlbumTitle, "add-album", "", "title for new album")

	RootCmd.AddCommand(uploadVideoCmd)
}
