package cmd

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"vk-cli/client"
	"vk-cli/config"
	"vk-cli/group"
	"vk-cli/util"
	"vk-cli/video"

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
		RunE: func(cmd *cobra.Command, args []string) error {
			groupId := 0

			// get config
			config := config.GetConfig()
			// build client
			vk, err := client.BuildVkClient()
			if err != nil {
				return err
			}
			// if group is a group title get group id
			if len(groupVar) > 0 {
				groupId, err = group.GetGroupId(vk, groupVar)
				if err != nil {
					return fmt.Errorf("failed to get group id: %w", err)
				}
			} else {
				groupId = config.DefaultGroupID
			}
			// add album if any
			if addAlbumTitle != "" {
				var err error

				album, err := video.GetVideoAlbumByTitle(vk, -groupId, addAlbumTitle)
				if err != nil {
					return fmt.Errorf("failed to get albums: %w", err)
				}

				if album == nil {
					albumID, err = video.AddVideoAlbum(vk, groupId, addAlbumTitle)
					if err != nil {
						return fmt.Errorf("failed to create album: %w", err)
					}

					slog.Info("Added new album", "id", albumID, "title", addAlbumTitle)
				} else {
					albumID = album.ID
				}
			}

			videofiles, err := util.GetFilenamesFromArgs(args, util.VideoFileExtSet)
			if err != nil {
				return fmt.Errorf("failed to get videos from args: %w", err)
			}

			videos, err := video.GetVideos(vk, -groupId, albumID)
			if err != nil {
				return fmt.Errorf("failed to get videos list: %w", err)
			}
			// create videos map
			videosTitleMap := make(map[string]bool, len(videos))

			for _, video := range videos {
				videosTitleMap[video.Title] = true
			}

			for idx, filename := range videofiles {
				videoName := filepath.Base(filename)
				videoTitle, _ := strings.CutSuffix(videoName, filepath.Ext(videoName))

				if _, ok := videosTitleMap[videoTitle]; ok {
					slog.Info("skipped video", "index", idx+1, "of", len(videofiles), "file", filename, "title", videoTitle)
					continue
				}

				slog.Info("upload video", "index", idx+1, "of", len(videofiles), "file", filename, "title", videoTitle)

				video.UploadVideo(vk, filename, groupId, albumID, videoTitle)
			}
			return nil
		},
	}
)

func init() {
	uploadVideoCmd.Flags().StringVar(&groupVar, "group", "", "group name or group id")
	uploadVideoCmd.Flags().IntVar(&albumID, "album", 0, "album id")
	uploadVideoCmd.Flags().StringVar(&addAlbumTitle, "add-album", "", "title for new album")

	RootCmd.AddCommand(uploadVideoCmd)
}
