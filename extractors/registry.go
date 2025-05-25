package extractors

import (
	"fmt"
	"dojer/downloader"
	"net/url"
	"sync"
	"testing"
)

type Tester interface {
	Test(t *testing.T)
}

type ExtractedData struct {
	Type          string
	Identifier    string
	Images        []string
	ChapterNumber string
	Source        string
}

type Extractor interface {
	Match(u *url.URL) bool
	Extract(u *url.URL) (ExtractedData, error)
}

type ExtractorRegistry struct {
	extractors []Extractor
}

func (er *ExtractorRegistry) Register(e Extractor) {
	er.extractors = append(er.extractors, e)
}

func (er *ExtractorRegistry) FindExtractor(u *url.URL) Extractor {
	for _, extractor := range er.extractors {
		if extractor.Match(u) {
			return extractor
		}
	}
	return nil
}

var registry = &ExtractorRegistry{
	extractors: []Extractor{},
}

func Run(urls []string, useQueue bool) {

	var wg sync.WaitGroup
	ch := make(chan *url.URL, 2)

	for _, urlString := range urls {
		u, err := url.Parse(urlString)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			continue
		}

		ch <- u
		extractor := registry.FindExtractor(u)
		if extractor != nil {
			wg.Add(1)
			go func() {
				defer func() {
					<-ch
					wg.Done()
				}()

				data, err := extractor.Extract(u)
				if err != nil {
					fmt.Printf("Extractor error with the url %s: %s\n", u, err.Error())
					return
				}
				downloadItem := downloader.DownloadItem{
					Type:          data.Type,
					ID:            data.Identifier,
					MangaName:     data.Identifier,
					Source:        data.Source,
					Pages:         data.Images,
					ChapterNumber: data.ChapterNumber,
				}

				if useQueue {
					downloader.DownloadQueue.AddItem(downloadItem)
				} else {
					downloadItem := downloadItem
					err := downloader.Download(downloadItem)
					if err != nil {
						fmt.Printf("err: %v\n", err)
					}
				}
			}()
		} else {
			fmt.Printf("No extractor found for URL %s\n", u)
		}
	}
	wg.Wait()
}
