package cmd

import (
	"fmt"
	"log/slog"
	"vk-cli/client"
	"vk-cli/config"
	"vk-cli/group"
	"vk-cli/photos"
	"vk-cli/util"

	"github.com/spf13/cobra"
)

var (
	photoAlbumID       = 0
	addPhotoAlbumTitle = ""
	groupPhotoVar      = ""

	uploadPhotoCmd = &cobra.Command{
		Use:   "upload-photo",
		Short: " upload photos(s) to group",
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
			if len(groupPhotoVar) > 0 {
				groupId, err = group.GetGroupId(vk, groupPhotoVar)
				if err != nil {
					return fmt.Errorf("failed to get group id: %w", err)
				}
				// try to get group id by name
			} else {
				groupId = config.DefaultGroupID
			}
			// add album if any
			if addPhotoAlbumTitle != "" {
				var err error

				album, err := photos.GetPhotoAlbumByTitle(vk, -groupId, addPhotoAlbumTitle)
				if err != nil {
					return fmt.Errorf("failed to get albums: %w", err)
				}

				if album == nil {
					photoAlbumID, err = photos.AddPhotoAlbum(vk, groupId, addPhotoAlbumTitle)
					if err != nil {
						return fmt.Errorf("failed to create album: %w", err)
					}

					slog.Info("Added new album", "id", photoAlbumID, "title", addPhotoAlbumTitle)
				} else {
					photoAlbumID = album.ID
				}
			}

			// upload photos
			filenames, err := util.GetFilenamesFromArgs(args, util.ImageFileExtSet)
			if err != nil {
				return fmt.Errorf("failed to get photos from args: %w", err)
			}

			for idx, filename := range filenames {
				slog.Info("upload photo", "index", idx+1, "of", len(filenames), "file", filename)

				photos.UploadPhoto(vk, filename, groupId, photoAlbumID)
			}
			return nil
		},
	}
)

func init() {
	uploadPhotoCmd.Flags().StringVar(&groupPhotoVar, "group", "", "group name or group id")
	uploadPhotoCmd.Flags().IntVar(&photoAlbumID, "album", 0, "album id")
	uploadPhotoCmd.Flags().StringVar(&addPhotoAlbumTitle, "add-album", "", "title for new album")
	RootCmd.AddCommand(uploadPhotoCmd)
}
