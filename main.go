package main

import (
	"context"

	"github.com/gogodjzhu/listen-tube/internal/pkg/conf"
	_ "github.com/gogodjzhu/listen-tube/internal/pkg/log"
	"github.com/gogodjzhu/listen-tube/web/controller"
)

func main() {
	config := &conf.Config{
		WebConfig: &conf.WebConfig{
			Port: 8080,
			Host: "0.0.0.0",
		},
		DBConfig: &conf.DBConfig{
			DSN:    "/tmp/listen-demo/listen-tube.db",
			Driver: conf.SQLiteDriver,
		},
		DownloaderConfig: &conf.DownloaderConfig{
			BinUri:   "/tmp/listen-demo/listen-tube/.bin/yt-dlp",
			BinURL:   "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux",
			BasePath: "/tmp/listen-demo/listen-tube/",
		},
	}

	ctx := context.Background()
	c, err := controller.NewController(ctx, config)
	if err != nil {
		panic(err)
	}
	c.Start()
}
