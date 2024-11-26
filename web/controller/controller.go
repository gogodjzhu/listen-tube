package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gogodjzhu/listen-tube/web/controller/buzz"
	"github.com/gogodjzhu/listen-tube/internal/app/subscribe"
	"github.com/gogodjzhu/listen-tube/internal/pkg/db"
	"github.com/gogodjzhu/listen-tube/internal/pkg/db/dao"
	"github.com/gogodjzhu/listen-tube/internal/pkg/tube/downloader"
	"github.com/gogodjzhu/listen-tube/internal/pkg/tube/fetcher"
)

type Controller struct {
	Router *gin.Engine
	Conf   *Config
}

type Config struct {
	Port            int
	DBConfig        *db.Config
	DownloaderConfig *downloader.Config
}

func NewController(conf *Config) (*Controller, error) {
	r := gin.Default()
	r.StaticFS("/static", http.Dir("web/static"))

	// Create SubscribeService instance
	ds, err := db.NewDatabaseSource(conf.DBConfig)
	if err != nil {
		return nil, err
	}

	unionMapper, err := dao.NewUnionMapper(ds)
	if err != nil {
		return nil, err
	}

	downloaderInstance, err := downloader.NewDownloader(conf.DownloaderConfig)
	if err != nil {
		return nil, err
	}

	fetcherInstance := fetcher.NewFetcher(fetcher.Config{})

	ctx := context.Background()
	subscribeService, err := subscribe.NewSubscribeService(ctx, unionMapper, downloaderInstance, fetcherInstance)
	if err != nil {
		return nil, err
	}

	buzzController, err := buzz.NewBuzzController(subscribeService)
	if err != nil {
		return nil, err
	}

	r.POST("/subscription/add", func(ctx *gin.Context) {
		var req buzz.AddSubscriptionRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := buzzController.AddSubscription(&req)
		if result.Err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	})

	r.POST("/subscription/delete", func(ctx *gin.Context) {
		var req buzz.DeleteSubscriptionRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := buzzController.DeleteSubscription(&req)
		if result.Err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	})

	r.GET("/subscription/list", func(ctx *gin.Context) {
		var req buzz.ListSubscriptionRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := buzzController.ListSubscription(&req)
		if result.Err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	})

	r.GET("/content/list", func(ctx *gin.Context) {
		var req buzz.ListContentRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := buzzController.ListContent(&req)
		if result.Err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	})

	return &Controller{
		Conf:   conf,
		Router: r,
	}, nil
}

func (c *Controller) Start() error {
	return c.Router.Run(fmt.Sprintf("0.0.0.0:%d", c.Conf.Port)) // listen and serve on 0.0.0.0:8080
}