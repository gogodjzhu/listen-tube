package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gogodjzhu/listen-tube/internal/app/auth"
	"github.com/gogodjzhu/listen-tube/internal/app/subscribe"
	"github.com/gogodjzhu/listen-tube/internal/pkg/conf"
	"github.com/gogodjzhu/listen-tube/internal/pkg/db"
	"github.com/gogodjzhu/listen-tube/internal/pkg/db/dao"
	"github.com/gogodjzhu/listen-tube/web/controller/buzz"
	"github.com/gogodjzhu/listen-tube/web/controller/middleware/jwt"
)

type Controller struct {
	Router          *gin.Engine
	Conf            *conf.Config
	subcribeService *subscribe.SubscribeService
	authService     *auth.AuthService
}

func NewController(ctx context.Context, conf *conf.Config) (*Controller, error) {
	r := gin.Default()
	// TODO: only allow specific origin
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8090"}
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
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

	subscribeService, err := subscribe.NewSubscribeService(unionMapper, conf.SubscriberConfig)
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
	openApiController, err := buzz.NewOpenAPIController(c.subcribeService)
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
	openapiGroup := c.Router.Group("/openapi")
	openApiController.AddHandler(openapiGroup)

	return c.Router.Run(fmt.Sprintf("0.0.0.0:%d", c.Conf.WebConfig.Port)) // listen and serve on 0.0.0.0:8080
}
