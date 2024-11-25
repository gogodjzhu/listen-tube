package buzz

import (
	"time"

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

// AddSubscription adds a new subscription for a user to a channel.
func (c *BuzzController) AddSubscription(req *AddSubscriptionRequest) *AddSubscriptionResult {
	err := c.subscribeService.AddSubscription(req.UserID, req.ChannelID)
	return &AddSubscriptionResult{Err: err}
}

// DeleteSubscription deletes a subscription for a user to a channel.
func (c *BuzzController) DeleteSubscription(req *DeleteSubscriptionRequest) *DeleteSubscriptionResult {
	err := c.subscribeService.DeleteSubscription(req.UserID, req.ChannelID)
	return &DeleteSubscriptionResult{Err: err}
}

// ListSubscription lists all subscriptions for a user.
func (c *BuzzController) ListSubscription(req *ListSubscriptionRequest) *ListSubscriptionResult {
	subscriptions, err := c.subscribeService.ListSubscription(req.UserID)
	if err != nil {
		return &ListSubscriptionResult{Err: err}
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
	return &ListSubscriptionResult{Subscriptions: result}
}

// ListContent lists all contents for a user.
func (c *BuzzController) ListContent(req *ListContentRequest) *ListContentResult {
	contents, err := c.subscribeService.ListContent(req.UserID)
	if err != nil {
		return &ListContentResult{Err: err}
	}
	result := make([]*Content, len(contents))
	for i, content := range contents {
		result[i] = &Content{
			Platform:      content.Platform,
			ChannelCredit: content.ChannelCredit,
			ContentCredit: content.ContentCredit,
			Title:         content.Title,
			Thumbnail:     content.Thumbnail,
			State:         string(content.State),
			CreateAt:      content.CreateAt,
			UpdateAt:      content.UpdateAt,
		}
	}
	return &ListContentResult{Contents: result}
}

type AddSubscriptionRequest struct {
	UserID    string
	ChannelID string
}

type AddSubscriptionResult struct {
	Err error
}

type DeleteSubscriptionRequest struct {
	UserID    string
	ChannelID string
}

type DeleteSubscriptionResult struct {
	Err error
}

type ListSubscriptionRequest struct {
	UserID string
}

type ListSubscriptionResult struct {
	Subscriptions []*Subscription
	Err           error
}

type ListContentRequest struct {
	UserID string
}

type ListContentResult struct {
	Contents []*Content
	Err      error
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
	State         string    `json:"state"`
	CreateAt      time.Time `json:"create_at"`
	UpdateAt      time.Time `json:"update_at"`
}
