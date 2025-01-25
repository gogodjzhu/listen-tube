package subscribe

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/gogodjzhu/listen-tube/internal/pkg/conf"
	"github.com/gogodjzhu/listen-tube/internal/pkg/db/dao"
	"github.com/gogodjzhu/listen-tube/internal/pkg/tube/downloader"
	"github.com/gogodjzhu/listen-tube/internal/pkg/tube/fetcher"
	"github.com/gogodjzhu/listen-tube/internal/pkg/util/str"
)

// Subscribe is an interface that defines the methods that must be implemented by a subscription service.
type Subscribe interface {
	// Subscribe subscribes a user to a channel
	Subscribe(userID, channelID string) error
}

type SubscribeService struct {
	subscriptionMapper *dao.SubscriptionMapper
	channelMapper      *dao.ChannelMapper
	contentMapper      *dao.ContentMapper
	userMapper         *dao.UserMapper
	downloader         *downloader.Downloader
	fetcher            *fetcher.Fetcher
}

func NewSubscribeService(mapper *dao.UnionMapper, config *conf.SubscriberConfig) (*SubscribeService, error) {
	downloader, err := downloader.NewDownloader(config.DownloaderConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create downloader")
	}
	fetcher := fetcher.NewFetcher(config.FetcherConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create fetcher")
	}
	ss := &SubscribeService{
		subscriptionMapper: mapper.SubscriptionMapper,
		userMapper:         mapper.UserMapper,
		channelMapper:      mapper.ChannelMapper,
		contentMapper:      mapper.ContentMapper,
		downloader:         downloader,
		fetcher:            fetcher,
	}

	return ss, nil
}

// Start the background tasks to fetch and download content periodically
func (s *SubscribeService) Start(ctx context.Context) error {
	go s.fetcher.TryStart(ctx, s.takeNextFetcher, s.updateFetchResult)
	go s.downloader.TryStart(ctx, s.takeNextDownload, s.updateDownloadResult)
	return nil
}

// TODO: test this method
func (s *SubscribeService) takeNextDownload() *dao.Content {
	sql := "SELECT * FROM t_content WHERE state = ? ORDER BY published_time DESC LIMIT 1"
	contents, err := s.contentMapper.SelectBySQL(sql, dao.ContentStatePrepared)
	if err != nil {
		log.Errorf("failed to list content: %v", err)
		return nil
	}
	if len(contents) == 0 {
		log.Warn("no content to download...")
		return nil
	}
	return contents[0]
}

// TODO: test this method
func (s *SubscribeService) updateDownloadResult(c dao.Content, r *downloader.Result) {
	state := dao.ContentStateDownloaded
	info := "finished"
	if !r.Finished {
		state = dao.ContentStateFailed
		log.Warnf("failed to download content %s", c.ContentCredit)
	}
	_, err := s.contentMapper.Update(&dao.Content{ID: c.ID}, &dao.Content{
		State:    state,
		Path:     r.Output,
		Info:     info,
		UpdateAt: time.Now(),
	})
	if err != nil {
		log.Errorf("failed to update content %s, err%v", c.ContentCredit, err)
	}
}

// TODO: test this method
func (s *SubscribeService) takeNextFetcher() *dao.Channel {
	sql := "SELECT * FROM t_channel ORDER BY update_at ASC LIMIT 1"
	channels, err := s.channelMapper.SelectBySQL(sql)
	if err != nil {
		log.Errorf("failed to list channel, err:%v", err)
		return nil
	}
	if len(channels) == 0 {
		log.Warn("no channel to fetch...")
		return nil
	}
	return channels[0]
}

// TODO: test this method
func (s *SubscribeService) updateFetchResult(c *dao.Channel, r *fetcher.Result) {
	for _, content := range r.Contents {
		state := dao.ContentStatePrepared
		info := "prepared"
		if content.MembersOnly {
			state = dao.ContentStateFailed
			info = "skip for members only"
		}
		newContent := &dao.Content{
			Platform:      "YouTube",
			ChannelCredit: c.ChannelCredit,
			Title:         content.Title,
			Thumbnail:     content.Thumbnail,
			ContentCredit: content.Credit,
			State:         state,
			Info:          info,
			PublishedTime: content.PublishedTime,
			Length:        content.Length,
			CreateAt:      time.Now(),
			UpdateAt:      time.Now(),
		}
		if oldContents, err := s.contentMapper.Select(&dao.Content{ContentCredit: content.Credit}); err != nil {
			log.Errorf("failed to list content %s, err:%v", content.Credit, err)
			continue
		} else if len(oldContents) > 0 {
			log.Debugf("content %s already exists", content.Credit)
			continue
		}
		if _, err := s.contentMapper.Insert(newContent); err != nil {
			log.Errorf("failed to create content %s, err:%v", content.Credit, err)
		}
	}
	if _, err := s.channelMapper.Update(&dao.Channel{ID: c.ID}, &dao.Channel{
		UpdateAt: time.Now(),
	}); err != nil {
		log.Errorf("failed to update channel %s, err:%v", c.ChannelCredit, err)
	}
}

// AddSubscription adds a new subscription for a user to a channel.
func (s *SubscribeService) AddSubscription(userCredit, channelCredit string) error {
	// check if the user exists
	user, err := s.userMapper.Select(&dao.User{Credit: userCredit})
	if err != nil || len(user) == 0 {
		return fmt.Errorf("user does not exist")
	}

	// check if the user has already subscribed to the channel, if so, return an error
	subscriptions, err := s.subscriptionMapper.Select(&dao.Subscription{UserCredit: userCredit, ChannelCredit: channelCredit})
	if err == nil && len(subscriptions) > 0 {
		return fmt.Errorf("already subscribed to the channel")
	}

	// create a new channel if not exists, and update the channel with the latest fetch result
	existedChannel, err := s.channelMapper.Select(&dao.Channel{ChannelCredit: channelCredit})
	if err != nil {
		return fmt.Errorf("failed to list channel")
	}
	if len(existedChannel) == 0 {
		result, err := s.fetcher.Fetch(fetcher.FetchOption{ChannelCredit: channelCredit})
		if err != nil || len(result.Contents) == 0 {
			return fmt.Errorf("channel does not exist")
		}
		newChannel := &dao.Channel{
			Platform:      dao.Platform(result.Platform),
			Name:          result.Title,
			Description:   result.Description,
			OwnerUrls:     str.ArrayToStringWithSplit(result.OwnerUrls, ","),
			Thumbnails:    str.ArrayToStringWithSplit(result.Thumbnails, ","),
			ChannelCredit: result.ChannelID,
			CreateAt:      time.Now(),
			UpdateAt:      time.Now(),
		}
		// create a new channel
		if _, err = s.channelMapper.Insert(newChannel); err != nil {
			return fmt.Errorf("failed to create new channel")
		}
		// update the channel with the latest fetch result
		s.updateFetchResult(newChannel, result)
	}

	if _, err = s.subscriptionMapper.Insert(&dao.Subscription{
		UserCredit:    userCredit,
		ChannelCredit: channelCredit,
		CreateAt:      time.Now(),
		UpdateAt:      time.Now(),
	}); err != nil {
		return fmt.Errorf("failed to create subscription")
	}
	return nil
}

// DeleteSubscription deletes an existing subscription for a user from a channel.
func (s *SubscribeService) DeleteSubscription(userCredit, channelCredit string) error {
	// check if the user exists
	user, err := s.userMapper.Select(&dao.User{Credit: userCredit})
	if err != nil || len(user) != 1 {
		return fmt.Errorf("user does not exist, err: %v", err)
	}

	// check if the user has already subscribed to the channel, if not, return an error
	subscriptions, err := s.subscriptionMapper.Select(&dao.Subscription{UserCredit: userCredit, ChannelCredit: channelCredit})
	if err != nil || len(subscriptions) != 1 {
		return fmt.Errorf("not subscribed to the channel, err: %v", err)
	}

	// delete the subscription record
	_, err = s.subscriptionMapper.Delete(subscriptions[0])
	if err != nil {
		return fmt.Errorf("failed to delete subscription, err: %v", err)
	}
	return nil
}

// ListSubscription lists all subscriptions for a user.
func (s *SubscribeService) ListSubscription(userCredit string) ([]*dao.Subscription, error) {
	// check if the user exists
	user, err := s.userMapper.Select(&dao.User{Credit: userCredit})
	if err != nil || len(user) == 0 {
		return nil, fmt.Errorf("user does not exist")
	}

	// list the subscriptions and related channels of the user
	return s.subscriptionMapper.Select(&dao.Subscription{UserCredit: userCredit})
}

// ListContent lists all contents for a user.
func (s *SubscribeService) ListContent(userCredit string, pageIndex, pageSize int) ([]*dao.Content, error) {
	// check if the user exists
	user, err := s.userMapper.Select(&dao.User{Credit: userCredit})
	if err != nil || len(user) == 0 {
		return nil, fmt.Errorf("user does not exist")
	}

	// list the subscribed channels of the user
	subscriptions, err := s.subscriptionMapper.Select(&dao.Subscription{UserCredit: userCredit})
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions")
	}
	channelCredits := make([]interface{}, 0)
	for _, subscription := range subscriptions {
		channelCredits = append(channelCredits, subscription.ChannelCredit)
	}

	// list the contents of the subscribed channels
	pageSql := "SELECT * FROM t_content WHERE state = 3 AND channel_credit IN (?) ORDER BY published_time DESC LIMIT ? OFFSET ?"
	return s.contentMapper.SelectBySQL(pageSql, channelCredits, pageSize, (pageIndex-1)*pageSize)
}

// GetContent gets a content by its credit.
func (s *SubscribeService) GetContent(contentCredit string) (*dao.Content, error) {
	// list the content by its credit
	contents, err := s.contentMapper.Select(&dao.Content{ContentCredit: contentCredit})
	if err != nil || len(contents) == 0 {
		return nil, fmt.Errorf("content does not exist")
	}
	return contents[0], nil
}

// GetChannel gets a channel by its credit.
func (s *SubscribeService) GetChannel(channelCredit string) (*dao.Channel, error) {
	// list the channel by its credit
	channels, err := s.channelMapper.Select(&dao.Channel{ChannelCredit: channelCredit})
	if err != nil || len(channels) == 0 {
		return nil, fmt.Errorf("channel does not exist")
	}
	return channels[0], nil
}
