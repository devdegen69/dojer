package server

import (
	d "dojer/downloader"
	"dojer/extractors"
	"dojer/store"
	u "dojer/utils"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
)

func listDojs(c echo.Context) error {

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 40
	}
	seed := c.QueryParam("seed")

	ds, pagination := store.List(store.ListRequest{
		Page:  page,
		Seed:  seed,
		Limit: limit,
	})

	data := map[string]interface{}{
		"dojs":       ds,
		"pagination": pagination,
	}

	return c.JSON(200, data)
}

func getDoj(c echo.Context) error {
	var d store.Doujinshi
	id := c.Param("id")
	if id == "surprise" {
		d = store.PickRandom()
	} else {
		d = store.Get(id)
		if d.ID == "" {
			return echo.ErrNotFound
		}
	}
	pages := d.Pages

	var images []string
	var sizes []string
	per := 27

	for i := 1; i <= pages; i++ {
		page := fmt.Sprintf("%d.jpg", i)
		images = append(images, fmt.Sprintf("/downs/%s/%s", d.ID, page))
		imgPath := filepath.Join(u.GetDownloadsFolder(), d.ID, page)
		file, err := os.Open(imgPath)
		if err != nil {
			return err
		}
		img, _, err := image.DecodeConfig(file)
		if err != nil {
			return err
		}
		w := img.Width * per / 100
		h := img.Height * per / 100

		sizes = append(sizes, fmt.Sprintf("%dx%d", w, h))
	}

	data := map[string]interface{}{
		"title":    d.Title,
		"doj":      d,
		"counters": d.Counters(),
		"images":   images,
		"sizes":    sizes,
	}

	return c.JSON(200, data)
}

func searchDoj(c echo.Context) error {

	search := c.Param("query")
	if search == "" {
		return echo.ErrNotFound
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	d, pagination := store.Search(search, page)
	// if len(d) == 0 {
	// 	return c.JSON(410, echo.Map{"dojs": []string{}})
	// }

	data := map[string]interface{}{
		"dojs":       d,
		"pagination": pagination,
		"search":     search,
	}

	return c.JSON(200, data)
}

func downloadDojs(c echo.Context) error {
	downloadList := new(struct {
		Ids []string `json:"ids"`
	})

	if err := c.Bind(downloadList); err != nil {
		return err
	}

	for i, v := range downloadList.Ids {
		downloadList.Ids[i] = fmt.Sprintf("https://nhentai.net/g/%s/", v)
	}

	go extractors.Run(downloadList.Ids, true)

	return c.JSON(200, downloadList)
}

func getDownloadStatus(c echo.Context) error {
	return c.JSON(200, echo.Map{"items": d.CurrentDownloads})
}

func errorHandler(err error, c echo.Context) {
	switch v := err.(type) {
	case *echo.HTTPError:
		_ = c.JSON(v.Code, echo.Map{"message": v.Message})
	default:
		_ = c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error, try again later."})
		u.LogError(v.Error())
	}
}
