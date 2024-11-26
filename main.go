package main

import (
	"github.com/gogodjzhu/listen-tube/web/controller"
	"github.com/gogodjzhu/listen-tube/internal/pkg/db"
	"github.com/gogodjzhu/listen-tube/internal/pkg/tube/downloader"
)

func main() {
	conf := &controller.Config{
		Port: 8080,
		DBConfig: &db.Config{
			DSN: "/tmp/listen-tube.db",
		},
		DownloaderConfig: &downloader.Config{
			BinUri:   "/tmp/listen-tube/.bin/yt-dlp",
			BinURL:   "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux",
			BasePath: "/tmp/listen-tube/",
		},
	}

	c, err := controller.NewController(conf)
	if err != nil {
		panic(err)
	}
	c.Start()
}