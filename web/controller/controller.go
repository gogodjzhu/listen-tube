package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gogodjzhu/listen-tube/internal/app/auth"
	"github.com/gogodjzhu/listen-tube/internal/app/subscribe"
	"github.com/gogodjzhu/listen-tube/internal/pkg/db"
	"github.com/gogodjzhu/listen-tube/internal/pkg/db/dao"
	"github.com/gogodjzhu/listen-tube/internal/pkg/tube/downloader"
	"github.com/gogodjzhu/listen-tube/internal/pkg/tube/fetcher"
	"github.com/gogodjzhu/listen-tube/web/controller/buzz"
	"github.com/gogodjzhu/listen-tube/web/controller/middleware/jwt"
	"github.com/gin-contrib/cors"
)

type Controller struct {
	Router          *gin.Engine
	Conf            *Config
	subcribeService *subscribe.SubscribeService
	authService     *auth.AuthService
}

type Config struct {
	Port             int
	DBConfig         *db.Config
	DownloaderConfig *downloader.Config
}

func NewController(ctx context.Context, conf *Config) (*Controller, error) {
	r := gin.Default()
	// TODO: only allow specific origin
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))
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

	subscribeService, err := subscribe.NewSubscribeService(unionMapper, downloaderInstance, fetcherInstance)
	if err != nil {
		return nil, err
	}
	err = subscribeService.Start(ctx)
	if err != nil {
		return nil, err
	}
	authService, err := auth.NewAuthService(unionMapper)
	if err != nil {
		return nil, err
	}

	return &Controller{
		Conf:            conf,
		Router:          r,
		subcribeService: subscribeService,
		authService:     authService,
	}, nil
}

func (c *Controller) Start() error {
	buzzController, err := buzz.NewBuzzController(c.subcribeService)
	if err != nil {
		return err
	}
	jwtMiddleware, err := jwt.NewJWTMiddleware(c.authService)
	if err != nil {
		return err
	}

	authGroup := c.Router.Group("/auth")
	authGroup.POST("/login", jwtMiddleware.LoginHandler)
	authGroup.POST("/refresh_token", jwtMiddleware.RefreshHandler)
	authGroup.POST("/logout", jwtMiddleware.LogoutHandler)
	authGroup.POST("/register", jwtMiddleware.RegisterHandler)
	authGroup.Use(jwtMiddleware.MiddlewareFunc())
	{
		authGroup.GET("/current_user", jwtMiddleware.UserInfoHandler)
	}

	buzzGroup := c.Router.Group("/buzz")
	buzzGroup.Use(jwtMiddleware.MiddlewareFunc())
	{
		buzzController.AddHandler(buzzGroup)
	}

	return c.Router.Run(fmt.Sprintf("0.0.0.0:%d", c.Conf.Port)) // listen and serve on 0.0.0.0:8080
}
