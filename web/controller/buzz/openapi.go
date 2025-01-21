package buzz

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gogodjzhu/listen-tube/internal/app/subscribe"
)

type OpenAPIController struct {
	subscribeService *subscribe.SubscribeService
}

func NewOpenAPIController(subscribeService *subscribe.SubscribeService) (*OpenAPIController, error) {
	return &OpenAPIController{
		subscribeService: subscribeService,
	}, nil
}

func (c *OpenAPIController) AddHandler(r gin.IRoutes) error {
	// ...existing code...

	r.GET("/content/stream/:contentCredit", func(ctx *gin.Context) {
		contentCredit := ctx.Param("contentCredit")
		content, err := c.subscribeService.GetContent(contentCredit)
		if err != nil || content == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
			return
		}

		file, err := os.Open(content.Path)
		if (err != nil) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		defer file.Close()

		ctx.Header("Content-Type", "audio/mp3")
		http.ServeContent(ctx.Writer, ctx.Request, filepath.Base(content.Path), content.UpdateAt, file)
	})

	return nil
}
