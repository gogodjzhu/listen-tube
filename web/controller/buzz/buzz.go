package buzz

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gogodjzhu/listen-tube/internal/app/subscribe"
	utiltime "github.com/gogodjzhu/listen-tube/internal/pkg/util/time"
	"github.com/gogodjzhu/listen-tube/web/controller/middleware/jwt"
)

type BuzzController struct {
	subscribeService *subscribe.SubscribeService
}

func NewBuzzController(subscribeService *subscribe.SubscribeService) (*BuzzController, error) {
	return &BuzzController{
		subscribeService: subscribeService,
	}, nil
}

func (c *BuzzController) AddHandler(r gin.IRoutes) error {
	r.POST("/subscription/add", func(ctx *gin.Context) {
		var req AddSubscriptionRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userinfo := jwt.GetCurrentUser(ctx)
		result := c.AddSubscription(userinfo, &req)
		ctx.JSON(http.StatusOK, result)
	})

	r.POST("/subscription/delete", func(ctx *gin.Context) {
		var req DeleteSubscriptionRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userinfo := jwt.GetCurrentUser(ctx)
		result := c.DeleteSubscription(userinfo, &req)
		ctx.JSON(http.StatusOK, result)
	})

	r.GET("/subscription/list", func(ctx *gin.Context) {
		var req ListSubscriptionRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userinfo := jwt.GetCurrentUser(ctx)
		result := c.ListSubscription(userinfo, &req)
		ctx.JSON(http.StatusOK, result)
	})

	r.POST("/content/list", func(ctx *gin.Context) {
		var req ListContentRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userinfo := jwt.GetCurrentUser(ctx)
		result := c.ListContent(userinfo, &req)
		ctx.JSON(http.StatusOK, result)
	})

	// 新增的端点，用于流式传输音频文件
	r.GET("/content/stream/:contentCredit", func(ctx *gin.Context) {
		contentCredit := ctx.Param("contentCredit")
		content, err := c.subscribeService.GetContent(contentCredit)
		if err != nil || content == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
			return
		}

		file, err := os.Open(content.Path)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		defer file.Close()

		ctx.Header("Content-Type", "audio/mp3")
		http.ServeContent(ctx.Writer, ctx.Request, filepath.Base(content.Path), content.UpdateAt, file)
	})

	return nil
}

// AddSubscription adds a new subscription for a user to a channel.
func (c *BuzzController) AddSubscription(userInfo *jwt.UserInfo, req *AddSubscriptionRequest) *AddSubscriptionResult {
	if err := c.subscribeService.AddSubscription(userInfo.UserCredit, req.ChannelID); err != nil {
		return &AddSubscriptionResult{Code: 1, Msg: err.Error()}
	} else {
		return &AddSubscriptionResult{Code: 0, Msg: "ok"}
	}
}

// DeleteSubscription deletes a subscription for a user to a channel.
func (c *BuzzController) DeleteSubscription(userInfo *jwt.UserInfo, req *DeleteSubscriptionRequest) *DeleteSubscriptionResult {
	if err := c.subscribeService.DeleteSubscription(userInfo.UserCredit, req.ChannelID); err != nil {
		return &DeleteSubscriptionResult{Code: 1, Msg: err.Error()}
	} else {
		return &DeleteSubscriptionResult{Code: 0, Msg: "ok"}
	}
}

// ListSubscription lists all subscriptions for a user.
func (c *BuzzController) ListSubscription(userInfo *jwt.UserInfo, req *ListSubscriptionRequest) *ListSubscriptionResult {
	subscriptions, err := c.subscribeService.ListSubscription(userInfo.UserCredit)
	if err != nil {
		return &ListSubscriptionResult{Code: 1, Msg: err.Error()}
	}
	result := make([]*Subscription, len(subscriptions))
	for i, sub := range subscriptions {
		result[i] = &Subscription{
			Platform:      "youtube", // TODO: get platform
			ChannelCredit: sub.ChannelCredit,
			ChannelName:   "", // TODO: get channel name
			CreateAt:      sub.CreateAt.Unix(),
			UpdateAt:      sub.UpdateAt.Unix(),
		}
	}
	return &ListSubscriptionResult{Subscriptions: result, Code: 0, Msg: "ok"}
}

// ListContent lists all contents for a user.
func (c *BuzzController) ListContent(userInfo *jwt.UserInfo, req *ListContentRequest) *ListContentResult {
	contents, err := c.subscribeService.ListContent(userInfo.UserCredit, req.PageIndex, req.PageSize)
	if err != nil {
		return &ListContentResult{Code: 1, Msg: err.Error()}
	}
	channelCredits := make(map[string]string)
	for _, content := range contents {
		if _, ok := channelCredits[content.ChannelCredit]; !ok {
			channel, err := c.subscribeService.GetChannel(content.ChannelCredit)
			if err != nil {
				continue
			}
			channelCredits[content.ChannelCredit] = channel.Name
		}
	}

	result := make([]*Content, len(contents))
	for i, content := range contents {
		result[i] = &Content{
			Platform:      content.Platform,
			Credit:        content.ContentCredit,
			Name:          content.Title,
			ChannelCredit: content.ChannelCredit,
			ChannelName:   channelCredits[content.ChannelCredit],
			Thumbnail:     content.Thumbnail,
			PublishedTime: utiltime.TranslateDuration2Accessibility(time.Now(), content.PublishedTime),
			// format duration, format: 01:00:10, 10:10, 00:10
			Length:        utiltime.FormatDuration(content.Length),
			State:         int(content.State),
			CreateAt:      content.CreateAt.Unix(),
			UpdateAt:      content.UpdateAt.Unix(),
		}
	}
	return &ListContentResult{Contents: result, Code: 0, Msg: "ok"}
}



type AddSubscriptionRequest struct {
	ChannelID string `json:"channel_id"`
}

type AddSubscriptionResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type DeleteSubscriptionRequest struct {
	ChannelID string `json:"channel_id"`
}

type DeleteSubscriptionResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type ListSubscriptionRequest struct {
}

type ListSubscriptionResult struct {
	Subscriptions []*Subscription `json:"subscriptions"`
	Code          int             `json:"code"`
	Msg           string          `json:"msg"`
}

type ListContentRequest struct {
	PageIndex int `form:"page_index"`
	PageSize  int `form:"page_size"`
}

type ListContentResult struct {
	Contents []*Content `json:"contents"`
	Code     int        `json:"code"`
	Msg      string     `json:"msg"`
}

type Subscription struct {
	Platform      string `json:"platform"`
	ChannelCredit string `json:"channel_credit"`
	ChannelName   string `json:"channel_name"`
	CreateAt      int64  `json:"create_at"`
	UpdateAt      int64  `json:"update_at"`
}

type Content struct {
	Platform      string `json:"platform"`
	Name          string `json:"name"`
	Credit        string `json:"credit"`
	ChannelName   string `json:"channel_name"`
	ChannelCredit string `json:"channel_credit"`
	Thumbnail     string `json:"thumbnail"`
	PublishedTime string `json:"published_time"`
	Length        string `json:"length"`
	State         int    `json:"state"`
	CreateAt      int64  `json:"create_at"`
	UpdateAt      int64  `json:"update_at"`
}
