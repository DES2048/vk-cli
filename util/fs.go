package util

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type ExtSet map[string]bool

type File struct {
	Path string
	Info fs.FileInfo
}

var VideoFileExtSet = ExtSet{
	"mp4": true,
	"m4v": true,
	"avi": true,
	"mkv": true,
	"wmv": true,
	"mov": true,
}

var ImageFileExtSet = ExtSet{
	"jpg":  true,
	"jpeg": true,
	"bmp":  true,
	"png":  true,
	"gif":  true,
}

func TestFileNameByExtSet(filename string, extSet ExtSet) bool {
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	_, ok := extSet[strings.ToLower(ext)]
	return ok
}

func GetFilesByExtSet(root string, extSet ExtSet, recursive bool) ([]*File, error) {
	files := make([]*File, 0)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !recursive && d.IsDir() && path != root {
			return filepath.SkipDir
		}
		if TestFileNameByExtSet(path, extSet) {
			fi, err := d.Info()
			if err != nil {
				return err
			}
			files = append(files, &File{
				Path: path,
				Info: fi,
			})
		}
		return nil
	})
	return files, err
}

func GetVideoFilesFromDir(path string, recursive bool) ([]*File, error) {
	return GetFilesByExtSet(path, VideoFileExtSet, recursive)
}

func GetImageFilesFromDir(path string, recursive bool) ([]*File, error) {
	return GetFilesByExtSet(path, ImageFileExtSet, recursive)
}

func FilterFilesByExtSet(files []string, extSet ExtSet) []string {
	return slices.DeleteFunc(files, func(e string) bool {
		return !TestFileNameByExtSet(e, extSet)
	})
}

func GetFilesFromArgs(args []string, extSet ExtSet) ([]*File, error) {
	filenames := args
	videofiles := make([]*File, 0)

	// work with dirs if any
	for idx, filename := range filenames {
		stat, err := os.Stat(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to stat %s :%w", filename, err)
		}

		if stat.IsDir() {
			videos, err := GetFilesByExtSet(filename, extSet, false)
			if err != nil {
				return nil, fmt.Errorf("failed to get videofiles from dir: %w", err)
			}

			videofiles = append(videofiles, videos...)
			filenames = append(filenames[:idx], filenames[:idx+1]...)
		} else {
			if TestFileNameByExtSet(filename, VideoFileExtSet) {
				videofiles = append(videofiles, &File{
					Path: filename,
					Info: stat,
				})
			}
		}

	}
	return videofiles, nil
}
