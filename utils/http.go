package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
)


func Get(url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"cookie": viper.GetStringSlice("nhentai.cookies"),
		"user-agent": []string{viper.GetString("nhentai.user_agent")},
	}

	// panic(fmt.Sprintf("Cookies: %s, User agent: %s", cookies, userAgent))
	maxRetries := 5
	retryDelay := 2

	urlParts := strings.Split(url, "/")
	code := urlParts[len(urlParts)-1]

	var doc *goquery.Document
	for i := 0; i < maxRetries; i++ {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}

		switch resp.StatusCode {
		case 429:
			time.Sleep(time.Duration(retryDelay) * time.Second)
			continue
		case 403:
			resp.Body.Close()

			newErr := errors.New(fmt.Sprintf("[%s] get catch by cloudflare, try update the cookies or user-agent", code))
			cacheUrl := "https://webcache.googleusercontent.com/search?q=cache"
			fmt.Printf("Trying to use cached route of nhentai.net/g/%s\n", code)
			req, err := http.NewRequest("GET", fmt.Sprintf("%s:%s/", cacheUrl, url), nil)
			if err != nil {
				return nil, newErr
			}

			resp, err = http.DefaultClient.Do(req)
			if err != nil || resp.StatusCode == 404 {
				return nil, newErr
			}

		case 404:
			resp.Body.Close()
			return nil, errors.New("this doj don't exist.")
		case 200:
			break
		}

		document, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, err
		}

		doc = document
	}

	return doc, err
}
