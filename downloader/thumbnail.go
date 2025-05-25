package downloader

import (
	"bytes"
	"dojer/utils"
	"image/jpeg"
	"image/png"
	"golang.org/x/image/webp"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

func createThumbnail(d *DownloadItem) error {
	thumbnailPath := utils.GetThumbnailPathOf(d.ID)
	if utils.FileExists(thumbnailPath) {
		return nil
	}
	_ = utils.EnsureExists(thumbnailPath)

	imagePath := filepath.Join(d.GetPath(), "1.jpg")
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	imageData, err := io.ReadAll(imageFile)
	if err != nil {
		return err
	}

	img, err := jpeg.Decode(bytes.NewReader(imageData))
	if err != nil {
		img, err = png.Decode(bytes.NewReader(imageData))
		if err != nil {
			img, err = webp.Decode(bytes.NewReader(imageData))
			if err != nil {
				return err
			}
		}
	}

	previewW := 500
	previewH := math.Ceil(float64(previewW) * 1.5)
	resizedImg := resize.Thumbnail(uint(previewW), uint(previewH), img, resize.Bilinear)

	var buffer bytes.Buffer
	err = jpeg.Encode(&buffer, resizedImg, &jpeg.Options{Quality: 75})
	if err != nil {
		return err
	}

	err = os.WriteFile(thumbnailPath, buffer.Bytes(), 0755)
	if err != nil {
		return err
	}
	return nil
}
