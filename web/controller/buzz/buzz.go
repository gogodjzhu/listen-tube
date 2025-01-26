package buzz

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gogodjzhu/listen-tube/internal/app/subscribe"
	utiltime "github.com/gogodjzhu/listen-tube/internal/pkg/util/time"
	"github.com/gogodjzhu/listen-tube/web/controller/middleware/interceptor"
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

	r.GET("/content/list", func(ctx *gin.Context) {
		var req ListContentRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userinfo := jwt.GetCurrentUser(ctx)
		result := c.ListContent(userinfo, &req)
		ctx.JSON(http.StatusOK, result)
	})

	return nil
}

// AddSubscription adds a new subscription for a user to a channel.
func (c *BuzzController) AddSubscription(userInfo *jwt.UserInfo, req *AddSubscriptionRequest) *interceptor.APIResponseDTO[bool] {
	if err := c.subscribeService.AddSubscription(userInfo.UserCredit, req.ChannelID); err != nil {
		return interceptor.NewDefaultErrorResponse[bool](err.Error())
	} else {
		return interceptor.NewDefaultSuccessResponse(true)
	}
}

// DeleteSubscription deletes a subscription for a user to a channel.
func (c *BuzzController) DeleteSubscription(userInfo *jwt.UserInfo, req *DeleteSubscriptionRequest) *interceptor.APIResponseDTO[bool] {
	if err := c.subscribeService.DeleteSubscription(userInfo.UserCredit, req.ChannelID); err != nil {
		return interceptor.NewDefaultErrorResponse[bool](err.Error())
	} else {
		return interceptor.NewDefaultSuccessResponse(true)
	}
}

// ListSubscription lists all subscriptions for a user.
func (c *BuzzController) ListSubscription(userInfo *jwt.UserInfo, req *ListSubscriptionRequest) *interceptor.APIResponseDTO[[]*Subscription] {
	subscriptions, err := c.subscribeService.ListSubscription(userInfo.UserCredit)
	if err != nil {
		return interceptor.NewDefaultErrorResponse[[]*Subscription](err.Error())
	}
	result := make([]*Subscription, len(subscriptions))
	for i, sub := range subscriptions {
		channel, err := c.subscribeService.GetChannel(sub.ChannelCredit)
		if err != nil {
			continue
		}
		result[i] = &Subscription{
			Platform:         "youtube", // TODO: get platform
			ChannelName:      channel.Name,
			ChannelThubmnail: channel.Thumbnails,
			CreateAt:         sub.CreateAt.Unix(),
			UpdateAt:         sub.UpdateAt.Unix(),
		}
	}
	return interceptor.NewDefaultSuccessResponse(result)
}

// ListContent lists all contents for a user.
func (c *BuzzController) ListContent(userInfo *jwt.UserInfo, req *ListContentRequest) *interceptor.APIResponseDTO[[]*Content] {
	contents, err := c.subscribeService.ListContent(userInfo.UserCredit, req.PageIndex, req.PageSize)
	if err != nil {
		return interceptor.NewDefaultErrorResponse[[]*Content](err.Error())
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
			Length:   utiltime.FormatDuration(content.Length),
			State:    int(content.State),
			CreateAt: content.CreateAt.Unix(),
			UpdateAt: content.UpdateAt.Unix(),
		}
	}
	return interceptor.NewDefaultSuccessResponse(result)
}

type AddSubscriptionRequest struct {
	ChannelID string `json:"channel_id"`
}

type DeleteSubscriptionRequest struct {
	ChannelID string `json:"channel_id"`
}

type ListSubscriptionRequest struct {
}

type ListContentRequest struct {
	PageIndex int `form:"page_index"`
	PageSize  int `form:"page_size"`
}

type Subscription struct {
	Platform         string `json:"platform"`
	ChannelName      string `json:"channel_name"`
	ChannelThubmnail string `json:"channel_thumbnail"`
	CreateAt         int64  `json:"create_at"`
	UpdateAt         int64  `json:"update_at"`
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
