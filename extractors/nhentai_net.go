package extractors

import (
	"dojer/store"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
)

type NhentaiGalleryInfo struct {
	Code       string
	Name       string
	Title      string
	Images     []string
	Parodies   string
	Characters string
	Tags       string
	Artists    string
	Groups     string
	Languages  string
	Categories string
	Pages      string
	Uploaded   string
}

type Nhentai struct {
	pattern *regexp.Regexp
}

func init() {
	registry.extractors = append(registry.extractors, NhentaiExtractor())
	viper.SetDefault("nhentai.cookies", []string{})
	viper.SetDefault("nhentai.user_agent", "")
}

func NhentaiExtractor() *Nhentai {
	pattern := regexp.MustCompile(`^https?://(www\.)?nhentai\.net/g/\d+`)
	return &Nhentai{pattern: pattern}
}

func (e *Nhentai) Match(u *url.URL) bool {
	return e.pattern.MatchString(u.String())
}

func (e *Nhentai) Test(t *testing.T) {
	testURL, _ := url.Parse("http://www.nhentai.net/g/1234")
	result := e.Match(testURL)
	if !result {
		t.Errorf("Expected true, got false")
	} else {
		t.Log("Test Paseed")
	}
}

func (e *Nhentai) Extract(url *url.URL) (ExtractedData, error) {
	var data ExtractedData
	userAgent := viper.GetString("nhentai.user_agent")
	cookies := viper.GetStringSlice("nhentai.cookies")
	if len(cookies) == 0 || userAgent == "" {
		return data, fmt.Errorf("Cookies or User-agent missing for the extractor %s.", url.Host)
	}
	doc, err := fetchDocument(url, http.Header{
		"Accept":          []string{"text/html", "application/xhtml+xml", "application/xml;q=0.9", "image/avif", "image/jxl", "image/webp", "*/*;q=0.8"},
		"Accept-Language": []string{"en-US", "en;q=0.5"},
		"Connection":      []string{"keep-alive"},
		"Cookie":          cookies,
		"Host":            []string{"nhentai.net"},
		"Sec-Fetch-Dest":  []string{"document"},
		"Sec-Fetch-Mode":  []string{"navigate"},
		"Sec-Fetch-Site":  []string{"none"},
		"Sec-Fetch-User":  []string{"?1"},
		"Sec-GPC":         []string{"1"},
		"User-Agent":      []string{userAgent},
		// "Accept-Encoding": []string{"gzip", "deflate", "br"},
	})

	if err != nil {
		return data, err
	}

	g := getNhentaiGalleryInfo(doc, url.String())
	gPages, err := strconv.Atoi(g.Pages)
	if err != nil {
		gPages = 0
	}

	data.Identifier = g.Code
	data.Source = "nhentai.net"
	data.Type = "doujinshi"
	data.Images = g.Images

	nhentaiDoujinshi := &store.Doujinshi{
		ID:         g.Code,
		Name:       g.Name,
		Title:      g.Title,
		Parodies:   g.Parodies,
		Characters: g.Characters,
		Tags:       g.Tags,
		Artists:    g.Artists,
		Groups:     g.Groups,
		Languages:  g.Languages,
		Categories: g.Categories,
		Pages:      gPages,
		Uploaded:   g.Uploaded,
		CreatedAt:  time.Now(),
	}

	err = store.Insert(*nhentaiDoujinshi)
	if err != nil {
		return data, err
	}

	return data, nil
}

func getNhentaiGalleryInfo(doc *goquery.Document, url string) *NhentaiGalleryInfo {

	urlParts := strings.Split(url, "/")
	code := urlParts[len(urlParts)-2]
	title := doc.Find(".title").Text()
	name := doc.Find(".title>.pretty").First().Text()

	images := doc.Find("a.gallerythumb>img").Map(func(i int, img *goquery.Selection) string {
		val, exists := img.Attr("data-src")
		if exists {
			// original: https://2t.nhentai.net/gallleries/1233/1.webp
			re := regexp.MustCompile(`t(\d)`) // replace t1 with 1
			val = re.ReplaceAllString(val, "i$1") // replace t1 with i1
			re = regexp.MustCompile(`(\d)t`)
			val = re.ReplaceAllString(val, "$1")
			val := regexp.MustCompile(`\.(jpg|png|webp)\.(jpg|png|webp)`).ReplaceAllString(val, ".$1")
			return val
		}
		return ""
	})

	tags := make(map[string][]string)
	keys := []string{
		"parodies",
		"characters",
		"tags",
		"artists",
		"groups",
		"languages",
		"categories",
		"pages",
		"uploaded",
	}

	count := 0

	doc.Find(".tag-container").Each(func(i int, s *goquery.Selection) {
		key := keys[count]
		s.Find(".tags").Each(func(i int, s *goquery.Selection) {
			if s.Find("a").Size() == 0 {
				tags[key] = append(tags[key], "null")
			}
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				tagName := s.Find(".name").First().Text()
				if tagName != "" {
					tags[key] = append(tags[key], tagName)
				}
			})
			count = count + 1
		})
	})

	doujinshi := NhentaiGalleryInfo{
		Code:       code,
		Name:       name,
		Title:      title,
		Images:     images,
		Parodies:   strings.Join(tags["parodies"], ","),
		Characters: strings.Join(tags["characters"], ","),
		Tags:       strings.Join(tags["tags"], ","),
		Artists:    strings.Join(tags["artists"], ","),
		Groups:     strings.Join(tags["groups"], ","),
		Languages:  strings.Join(tags["languages"], ","),
		Categories: strings.Join(tags["categories"], ","),
		Pages:      strings.Join(tags["pages"], ","),
		Uploaded:   strings.Join(tags["uploaded"], ","),
	}

	return &doujinshi
}
