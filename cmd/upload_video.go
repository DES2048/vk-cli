package cmd

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"vk-cli/client"
	"vk-cli/config"
	"vk-cli/group"
	"vk-cli/util"
	"vk-cli/video"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

var (
	albumVar     = ""
	addAlbum     = false
	groupVar     = ""
	modTimeSince time.Time
	sizeFVar     = SizeFlagValue{}

	uploadVideoCmd = &cobra.Command{
		Use:   "upload-video",
		Short: " upload video(s) to group",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupId := 0
			albumId := 0

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

			// get album by title
			album, err := video.GetVideoAlbumByTitle(vk, -groupId, albumVar)
			if err != nil {
				return fmt.Errorf("failed to get albums: %w", err)
			}

			// if album not found
			if album == nil {
				if addAlbum {
					albumId, err = video.AddVideoAlbum(vk, groupId, albumVar)
					if err != nil {
						return fmt.Errorf("failed to create album: %w", err)
					}
					slog.Info("Added new album", "id", albumVar, "title", addAlbum)
				} else {
					return fmt.Errorf("failed to find album with title: %s", albumVar)
				}
			} else {
				albumId = album.ID
			}

			// get videofiles from args
			videofiles, err := util.GetFilesFromArgs(args, util.VideoFileExtSet)
			if err != nil {
				return fmt.Errorf("failed to get videos from args: %w", err)
			}

			// filter videofiles if any
			if !modTimeSince.IsZero() {
				slog.Debug("used mt since", "value", modTimeSince.Format(time.DateTime))
				videofiles = slices.DeleteFunc(videofiles, func(f *util.File) bool {
					return f.Info.ModTime().Before(modTimeSince)
				})
			}

			if sizeFVar.Value > 0 {
				slog.Debug("used size filter", "value", sizeFVar.Value, "gt", sizeFVar.IsGt)

				videofiles = slices.DeleteFunc(videofiles, func(f *util.File) bool {
					filterGt := sizeFVar.Value >= uint64(f.Info.Size())
					if sizeFVar.IsGt {
						return filterGt
					} else {
						return !filterGt
					}
					//	return sizeFVar.IsGt && cmp <= 0 || !sizeFVar.IsGt && cmp > 0
				})

			}

			// get videos list from vk
			videos, err := video.GetVideos(vk, -groupId, albumId)
			if err != nil {
				return fmt.Errorf("failed to get videos list: %w", err)
			}
			// create videos map
			videosTitleMap := make(map[string]bool, len(videos))

			for _, video := range videos {
				videosTitleMap[video.Title] = true
			}

			// upload videos loop
			// counters
			successCount, errCount, skipCount := 0, 0, 0
			for idx, file := range videofiles {
				videoName := file.Info.Name()
				videoTitle, _ := strings.CutSuffix(videoName, filepath.Ext(videoName))

				uplLogger := slog.With(
					slog.Int("index", idx),
					slog.Int("of", len(videofiles)),
					slog.String("file", file.Path),
					slog.String("title", videoTitle),
					slog.String("size_h", humanize.Bytes(uint64(file.Info.Size()))),
				)

				if _, ok := videosTitleMap[videoTitle]; ok {
					skipCount++
					uplLogger.Info("skipped video")
					continue
				}

				uplLogger.Info("upload video")

				err := video.UploadVideo(vk, file.Path, groupId, albumId, videoTitle)
				if err != nil {
					errCount++
					uplLogger.Error("Failed to upload video", slog.String("error", err.Error()))
				}

				successCount++
			}

			slog.Info("Summary", "success", successCount, "skipped", skipCount, "failed", errCount, "of", len(videofiles))
			return nil
		},
	}
)

func init() {
	uploadVideoCmd.Flags().StringVarP(&groupVar, "group", "g", "", "group name or group id")
	uploadVideoCmd.Flags().StringVarP(&albumVar, "album", "a", "", "album name or album id")
	uploadVideoCmd.Flags().BoolVar(&addAlbum, "add-album", false, "allow to create new album if it didn't find")
	uploadVideoCmd.Flags().TimeVar(&modTimeSince, "mt-since", time.Time{}, []string{time.DateTime, time.DateOnly}, "filter files by modtime since")
	uploadVideoCmd.Flags().Var(&sizeFVar, "size", " size filter in format <100mb or >1mb etc")
	RootCmd.AddCommand(uploadVideoCmd)
}
