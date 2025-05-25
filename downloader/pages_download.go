package downloader

import (
	u "dojer/utils"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func downloadPages(downloadItem *DownloadItem) error {
	var fail error
	var wg sync.WaitGroup

	sem := make(chan string, 3) // download 3 pages at a time
	count := make(chan int, 1)
	breakLoop := false

	p := progressList[downloadItem.ID]
	p.setMessage(u.Blue("Downloading"))

	for i, imageUrl := range downloadItem.Pages {
		if breakLoop {
			break
		}

		wg.Add(1)
		sem <- imageUrl
		count <- i + 1
		go func(imageUrl string, i int, count chan int) {
			defer func() {
				<-sem
				wg.Done()
			}()

			if err := pageWorker(imageUrl, downloadItem.GetPath(), i); err != nil {
				fail = err
				breakLoop = true
			}
			c := <-count
			p.setProgress(fmt.Sprintf("%d/%d", c, len(downloadItem.Pages)))
			downloadItem.Progress = float32((c / len(downloadItem.Pages)) * 100)
			downloadItem.AltProg = fmt.Sprintf("%d/%d", c, len(downloadItem.Pages))
			CurrentDownloads[downloadItem.ID] = *downloadItem
		}(imageUrl, i, count)
	}

	wg.Wait()

	defer func() {
		close(sem)
		close(count)
		delete(CurrentDownloads, downloadItem.ID)
	}()

	if fail != nil {
		os.RemoveAll(downloadItem.GetPath())
		return fail
	}

	p.setMessage(u.Green("Completed \u2713"))

	return nil
}

func pageWorker(url string, path string, i int) error {

	imagePath := filepath.Join(path, fmt.Sprintf("%d.jpg", i+1))

	if u.FileExists(imagePath) {
		return nil
	}
	err := os.MkdirAll(path, DEFAULT_PERM)
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		f := fmt.Sprintf("Try to get %s but got status %s\n", url, resp.Status)
		return errors.New(f)
	}

	file, err := os.Create(imagePath)
	if err != nil {
		resp.Body.Close()
		return err
	}

	if _, err := io.Copy(file, resp.Body); err != nil {
		resp.Body.Close()
		return err
	}
	resp.Body.Close()
	update()
	return nil
}
