package main

import (
	"context"

	"github.com/gogodjzhu/listen-tube/internal/pkg/db"
	"github.com/gogodjzhu/listen-tube/internal/pkg/tube/downloader"
	"github.com/gogodjzhu/listen-tube/web/controller"
)

func main() {
	conf := &controller.Config{
		Port: 8080,
		DBConfig: &db.Config{
			DSN: "/tmp/listen-demo/listen-tube.db",
		},
		DownloaderConfig: &downloader.Config{
			BinUri:   "/tmp/listen-demo/listen-tube/.bin/yt-dlp",
			BinURL:   "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux",
			BasePath: "/tmp/listen-demo/listen-tube/",
		},
	}

	ctx := context.Background()
	c, err := controller.NewController(ctx, conf)
	if err != nil {
		panic(err)
	}
	c.Start()
}
