package cmd

import (
	"log"
	"log/slog"
	"strconv"
	"vk-cli/auth"
	"vk-cli/config"
	"vk-cli/group"
	"vk-cli/photos"

	"github.com/SevereCloud/vksdk/v2/api"
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
			if len(groupPhotoVar) > 0 {
				groupId, err = strconv.Atoi(groupPhotoVar)
				if err != nil {
					// try to get group id by name
					group, err := group.GetGroupByName(vk, groupPhotoVar)
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
			if addPhotoAlbumTitle != "" {
				var err error

				album, err := photos.GetPhotoAlbumByTitle(vk, -groupId, addPhotoAlbumTitle)
				if err != nil {
					log.Fatalf("failed to get albums: %s\n", err)
				}

				if album == nil {
					photoAlbumID, err = photos.AddPhotoAlbum(vk, groupId, addPhotoAlbumTitle)
					if err != nil {
						log.Fatalf("Failed to create album: %s\n", err)
					}

					slog.Info("Added new album", "id", photoAlbumID, "title", addPhotoAlbumTitle)
				} else {
					photoAlbumID = album.ID
				}
			}

			// upload photos
			filenames := args

			for idx, filename := range filenames {
				slog.Info("upload photo", "index", idx+1, "of", len(filenames), "file", filename)

				photos.UploadPhoto(vk, filename, groupId, photoAlbumID)
			}
		},
	}
)

func init() {
	uploadPhotoCmd.Flags().StringVar(&groupPhotoVar, "group", "", "group name or group id")
	uploadPhotoCmd.Flags().IntVar(&photoAlbumID, "album", 0, "album id")
	uploadPhotoCmd.Flags().StringVar(&addPhotoAlbumTitle, "add-album", "", "title for new album")
	RootCmd.AddCommand(uploadPhotoCmd)
}
