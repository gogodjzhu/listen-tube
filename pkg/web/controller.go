package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	Router *gin.Engine
	Conf   *Config
	API    *APIController
}

type Config struct {
	Port int
}

func NewController(conf *Config) (*Controller, error) {
	r := gin.Default()
	r.StaticFS("/page", http.Dir("/IdeaProjects/listen-tube/pkg/web/page"))

	apiController, err := NewAPIController(r)
	if err != nil {
		return nil, err
	}
	return &Controller{
		Conf:   conf,
		Router: r,
		API:    apiController,
	}, nil
}

func (c *Controller) Start() error {
	return c.Router.Run(fmt.Sprintf("0.0.0.0:%d", c.Conf.Port)) // listen and serve on 0.0.0.0:8080
}

type APIController struct {
}

func NewAPIController(r *gin.Engine) (*APIController, error) {
	apiController := &APIController{}
	r.Handle(apiController.Ping())
	return apiController, nil
}

func (a *APIController) Ping() (string, string, func(context *gin.Context)) {
	return http.MethodGet, "ping", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "pong",
		})
	}
}

func main() {
	c, err := NewController(&Config{Port: 8080})
	if err != nil {
		panic(err)
	}
	c.Start()
}
