package downloader

import (
	"dojer/utils"
	"errors"
	"fmt"
	"path/filepath"
)

const DEFAULT_PERM = 0755

var CurrentDownloads = map[string]DownloadItem{}

type DownloadItem struct {
	Source        string
	ID            string `json:"id"`
	Type          string
	MangaName     string `name`
	ChapterNumber string
	Pages         []string
	Progress      float32 `json:"progress"`
	AltProg       string  `json:"prog"`
}

type ExtractedData struct {
	Type          string
	Identifier    string
	Images        []string
	ChapterNumber string
	Source        string
}

func (d *DownloadItem) AddToQueue() {
	DownloadQueue.AddItem(*d)
}

func (d *DownloadItem) GetPath() string {
	return filepath.Join(utils.GetDownloadsFolder(), d.ID)
}

func Download(item DownloadItem) error {
	if len(item.Pages) == 0 {
		return errors.New("Item pages is empty")
	}
	if isProgressListEmpty() {
		fmt.Fprint(pw, "\033[s")
	}
	CurrentDownloads[item.ID] = item
	createNewProgress(item.ID, "Starting")
	update()
	err := downloadPages(&item)
	if err != nil {
		return err
	}

	err = createThumbnail(&item)
	if err != nil {
		return err
	}

	delete(progressList, item.ID)
	return nil
}
