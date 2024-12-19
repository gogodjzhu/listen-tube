package buzz

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gogodjzhu/listen-tube/internal/app/subscribe"
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
		result := c.AddSubscription(&req)
		ctx.JSON(http.StatusOK, result)
	})

	r.POST("/subscription/delete", func(ctx *gin.Context) {
		var req DeleteSubscriptionRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := c.DeleteSubscription(&req)
		ctx.JSON(http.StatusOK, result)
	})

	r.GET("/subscription/list", func(ctx *gin.Context) {
		var req ListSubscriptionRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := c.ListSubscription(&req)
		ctx.JSON(http.StatusOK, result)
	})

	r.GET("/content/list", func(ctx *gin.Context) {
		var req ListContentRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := c.ListContent(&req)
		ctx.JSON(http.StatusOK, result)
	})

	return nil
}

// AddSubscription adds a new subscription for a user to a channel.
func (c *BuzzController) AddSubscription(req *AddSubscriptionRequest) *AddSubscriptionResult {
	if err := c.subscribeService.AddSubscription(req.UserID, req.ChannelID); err != nil {
		return &AddSubscriptionResult{Code: 1, Msg: err.Error()}
	} else {
		return &AddSubscriptionResult{Code: 0, Msg: "ok"}
	}
}

// DeleteSubscription deletes a subscription for a user to a channel.
func (c *BuzzController) DeleteSubscription(req *DeleteSubscriptionRequest) *DeleteSubscriptionResult {
	if err := c.subscribeService.DeleteSubscription(req.UserID, req.ChannelID); err != nil {
		return &DeleteSubscriptionResult{Code: 1, Msg: err.Error()}
	} else {
		return &DeleteSubscriptionResult{Code: 0, Msg: "ok"}
	}
}

// ListSubscription lists all subscriptions for a user.
func (c *BuzzController) ListSubscription(req *ListSubscriptionRequest) *ListSubscriptionResult {
	subscriptions, err := c.subscribeService.ListSubscription(req.UserID)
	if err != nil {
		return &ListSubscriptionResult{Code: 1, Msg: err.Error()}
	}
	result := make([]*Subscription, len(subscriptions))
	for i, sub := range subscriptions {
		result[i] = &Subscription{
			Platform:      "youtube", // TODO: get platform
			ChannelCredit: sub.ChannelCredit,
			ChannelName:   "", // TODO: get channel name
			CreateAt:      sub.CreateAt,
			UpdateAt:      sub.UpdateAt,
		}
	}
	return &ListSubscriptionResult{Subscriptions: result, Code: 0, Msg: "ok"}
}

// ListContent lists all contents for a user.
func (c *BuzzController) ListContent(req *ListContentRequest) *ListContentResult {
	contents, err := c.subscribeService.ListContent(req.UserID)
	if err != nil {
		return &ListContentResult{Code: 1, Msg: err.Error()}
	}
	result := make([]*Content, len(contents))
	for i, content := range contents {
		result[i] = &Content{
			Platform:      content.Platform,
			ChannelCredit: content.ChannelCredit,
			ContentCredit: content.ContentCredit,
			Title:         content.Title,
			Thumbnail:     content.Thumbnail,
			State:         int(content.State),
			CreateAt:      content.CreateAt,
			UpdateAt:      content.UpdateAt,
		}
	}
	return &ListContentResult{Contents: result, Code: 0, Msg: "ok"}
}

type AddSubscriptionRequest struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
}

type AddSubscriptionResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type DeleteSubscriptionRequest struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
}

type DeleteSubscriptionResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type ListSubscriptionRequest struct {
	UserID string `json:"user_id"`
}

type ListSubscriptionResult struct {
	Subscriptions []*Subscription `json:"subscriptions"`
	Code          int             `json:"code"`
	Msg           string          `json:"msg"`
}

type ListContentRequest struct {
	UserID string `json:"user_id"`
}

type ListContentResult struct {
	Contents []*Content `json:"contents"`
	Code     int        `json:"code"`
	Msg      string     `json:"msg"`
}

type Subscription struct {
	Platform      string    `json:"platform"`
	ChannelCredit string    `json:"channel_credit"`
	ChannelName   string    `json:"channel_name"`
	CreateAt      time.Time `json:"create_at"`
	UpdateAt      time.Time `json:"update_at"`
}

type Content struct {
	Platform      string    `json:"platform"`
	ChannelCredit string    `json:"channel_credit"`
	ContentCredit string    `json:"content_credit"`
	Title         string    `json:"title"`
	Thumbnail     string    `json:"thumbnail"`
	State         int       `json:"state"`
	CreateAt      time.Time `json:"create_at"`
	UpdateAt      time.Time `json:"update_at"`
}
