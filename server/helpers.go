package server

import (
	u "dojer/utils"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var indexPage = "index.html"

func loggerValuesFunc(c echo.Context, v middleware.RequestLoggerValues) error {
	var status, method string
	route := color.New(color.Bold, color.FgWhite).Sprint(v.URI)

	switch {
	case v.Status >= 200 && v.Status < 300:
		status = u.Green(v.Status)
	case v.Status >= 300 && v.Status < 400:
		status = u.Magenta(v.Status)
	case v.Status >= 400 && v.Status < 500:
		status = u.Yellow(v.Status)
	default:
		status = u.Red(status)
	}

	switch v.Method {
	case "GET":
		method = u.Cyan(v.Method)
	case "POST":
		method = u.Blue(v.Method)
	case "PUT":
		method = u.Yellow(v.Method)
	case "DELETE":
		method = u.Red(v.Method)
	default:
		method = u.White(v.Method)
	}

	logFormat := fmt.Sprintf("%s %s %s %s %s", v.RemoteIP, status, method, v.Host, route)
	u.Log(logFormat)

	return nil
}

func fileFS(c echo.Context, file string, filesystem fs.FS) error {
	f, err := filesystem.Open(file)
	if err != nil {
		return echo.ErrNotFound
	}
	defer f.Close()

	fi, _ := f.Stat()
	if fi.IsDir() {
		file = filepath.Join(file, indexPage)
		f, err = filesystem.Open(file)
		if err != nil {
			return echo.ErrNotFound
		}
		defer f.Close()
		if fi, err = f.Stat(); err != nil {
			return err
		}
	}
	ff, ok := f.(io.ReadSeeker)
	if !ok {
		return errors.New("file does not implement io.ReadSeeker")
	}
	http.ServeContent(c.Response(), c.Request(), fi.Name(), fi.ModTime(), ff)
	return nil
}
