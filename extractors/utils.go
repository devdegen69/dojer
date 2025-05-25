package extractors

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func fetchDocument(u *url.URL, headers http.Header) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	if len(headers) > 0 {
		req.Header = headers
	}

	maxRetries := 5
	var doc *goquery.Document
	for i := 1; i < maxRetries; i++ {
		time.Sleep(time.Millisecond * 400)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making HTTP request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			switch resp.StatusCode {
			case http.StatusTooManyRequests:
				time.Sleep(time.Second * 3)
				continue

			default:
				return nil, fmt.Errorf("HTTP error: status code %d", resp.StatusCode)
			}
		}

		doc, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error parsing HTML content: %w", err)
		}

		break
	}

	return doc, nil
}

func filterNotEmpty(strings []string) []string {
	var result []string
	for _, s := range strings {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}
