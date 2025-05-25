package server

import (
	d "dojer/downloader"
	u "dojer/utils"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func Init(port string) {
	e := echo.New()
	e.HTTPErrorHandler = errorHandler
	e.HideBanner = true

	e.Static("/downs", u.GetDownloadsFolder())

	routes(e)
	setMiddlewares(e)
	go startBackupInterval()
	go d.ListenDownloadQueue()

	serverPort := fmt.Sprintf(":%s", port)
	e.Logger.Fatal(e.Start(serverPort))
}

func startBackupInterval() {
	interval := time.Hour * time.Duration(viper.GetInt("backup.interval"))
	for tick := range time.Tick(interval) {
		err := backup(tick)
		if err != nil {
			u.LogWarn(err.Error())
		} else {
			u.LogSuccess("Backup", "Done backuping the database.")
		}
	}
}

func setMiddlewares(e *echo.Echo) {
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:     true,
		LogURI:        true,
		LogMethod:     true,
		LogHost:       true,
		LogRemoteIP:   true,
		LogValuesFunc: loggerValuesFunc,
	}))
}
