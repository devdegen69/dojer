package server

import (
	"dojer/ui"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func routes(e *echo.Echo) {
	bindApiRoutes(e)
	bindSpaRoutes(e)
}

func bindSpaRoutes(e *echo.Echo) {
	e.GET("/*", func(c echo.Context) error {
		p := c.Param("*")

		tmpPath, err := url.PathUnescape(p)
		if err != nil {
			return fmt.Errorf("failed to unescape path variable: %w", err)
		}
		p = tmpPath

		name := filepath.ToSlash(filepath.Clean(strings.TrimPrefix(p, "/")))

		fileErr := fileFS(c, name, ui.DistDirFs)

		if fileErr != nil && errors.Is(fileErr, echo.ErrNotFound) {
			return fileFS(c, "index.html", ui.DistDirFs)
		}

		return fileErr
	},
		middleware.Gzip(),
	)
}

func bindApiRoutes(e *echo.Echo) {
	g := e.Group("/api")

	g.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept},
	}))

	g.GET("", listDojs).Name = "listDojs"
	g.GET("/search/:query", searchDoj).Name = "serachDojs"
	g.GET("/:id", getDoj).Name = "getDoj"
	g.POST("/download", downloadDojs).Name = "getDoj"
	g.GET("/downloadStatus", getDownloadStatus).Name = "getDoj"
}
