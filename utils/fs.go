package utils

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// GetDataPath concatenates the provided paths with the configured data directory and returns the resulting path.
// If no path is provided, it returns the configured data directory.
//
// Example:
//
//	// Assuming data_dir is configured as "{data_path}"
//	path := GetDataPath("folder1", "file.txt")
//	// Returns: "{data_path}/folder1/file.txt"
func GetDataPath(path ...string) string {
	if flag.Lookup("test.v") != nil {
		dir, err := os.UserConfigDir()
		if err != nil {
			panic(err)
		}
		joinedPaths := filepath.Join(path...)
		dataPath := filepath.Join(dir, "testingdojer", joinedPaths)
		return dataPath
	}
	joinedPaths := filepath.Join(path...)
	dataPath := filepath.Join(viper.GetString("data_dir"), joinedPaths)
	return dataPath
}

// GetDownloadsFolder returns the path to the downloads folder within the configured data directory.
//
// Example:
//
//	// Assuming data_dir is configured as "{data_path}"
//	folder := GetDownloadsFolder()
//	// Returns: "{data_path}/downloads"
func GetDownloadsFolder() string {
	return GetDataPath("downloads")
}

// GetThumbnailsFolder returns the path to the thumbnails folder within the downloads folder.
//
// Example:
//
//	// Assuming data_dir is configured as "{data_path}"
//	folder := GetThumbnailsFolder()
//	// Returns: "{data_path}/downloads/thumbs"
func GetThumbnailsFolder() string {
	return filepath.Join(GetDownloadsFolder(), "thumbs")
}

// GetThumbnailPathOf constructs the path for a thumbnail file based on the provided ID.
//
// Example:
//
//	// Assuming thumbnails are stored in "{data_path}/downloads/thumbs"
//	thumbnailPath := GetThumbnailPathOf("123")
//	// Returns: "{data_path}/downloads/thumbs/123.jpg"
func GetThumbnailPathOf(id string) string {
	thumbsFolder := GetThumbnailsFolder()

	return filepath.Join(thumbsFolder, fmt.Sprintf("%s.jpg", id))
}

// EnsureExists creates directories along the provided path if they don't already exist.
// It returns an error if any error occurs during directory creation.
//
// Example:
//
//	// Create the directory "{data_path}/downloads/thumbs" if it doesn't exist
//	err := EnsureExists("{data_path}/downloads/thumbs")
//	if err != nil {
//	    // Handle error
//	}
func EnsureExists(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return nil
}

// FileExists checks if a file exists at the specified path.
//
// Example:
//
//	// Check if the file "{data_path}/downloads/thumbs/123.jpg" exists
//	exists := FileExists("{data_path}/downloads/thumbs/123.jpg")
//	// Returns: true if the file exists, false otherwise
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
