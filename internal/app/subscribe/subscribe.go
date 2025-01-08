package subscribe

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

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

func NewSubscribeService(mapper *dao.UnionMapper, downloader *downloader.Downloader, fetcher *fetcher.Fetcher) (*SubscribeService, error) {
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
	// Schedule to execute `FetchContent` and `DownloadContent` periodically in the background
	go s.scheduleFetchContent(ctx)
	go s.scheduleDownloadContent(ctx)
	return nil
}

func (s *SubscribeService) scheduleFetchContent(ctx context.Context) {
	// Execute FetchContent immediately
	log.Info("Starting initial FetchContent task")
	fetchCount, err := s.FetchContent()
	if err != nil {
		log.Errorf("Error during initial FetchContent task: %v", err)
	} else {
		log.Infof("Initial FetchContent task completed, fetched %d contents", fetchCount)
	}

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("Stopping scheduled FetchContent task")
			return
		case <-ticker.C:
			log.Info("Starting scheduled FetchContent task")
			fetchCount, err := s.FetchContent()
			if err != nil {
				log.Errorf("Error during scheduled FetchContent task: %v", err)
			} else {
				log.Infof("Scheduled FetchContent task completed, fetched %d contents", fetchCount)
			}
		}
	}
}

func (s *SubscribeService) scheduleDownloadContent(ctx context.Context) {
	// Execute DownloadContent immediately
	log.Info("Starting initial DownloadContent task")
	err := s.DownloadContent()
	if err != nil {
		log.Errorf("Error during initial DownloadContent task: %v", err)
	} else {
		log.Info("Initial DownloadContent task completed")
	}

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("Stopping scheduled DownloadContent task")
			return
		case <-ticker.C:
			log.Info("Starting scheduled DownloadContent task")
			err := s.DownloadContent()
			if err != nil {
				log.Errorf("Error during scheduled DownloadContent task: %v", err)
			} else {
				log.Info("Scheduled DownloadContent task completed")
			}
		}
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

	// check if the channel exists. if not, check by `ChannelFetcher.Fetch()` and create a new channel if it exists
	channel, err := s.channelMapper.Select(&dao.Channel{ChannelCredit: channelCredit})
	if err != nil || len(channel) == 0 {
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
		_, err = s.channelMapper.Insert(newChannel)
		if err != nil {
			return fmt.Errorf("failed to create new channel")
		}
		for _, content := range result.Contents {
			_, err := s.contentMapper.Insert(&dao.Content{
				Platform:      string(result.Platform),
				ChannelCredit: result.ChannelID,
				Title:         content.Title,
				Thumbnail:     content.Thumbnail,
				ContentCredit: content.Credit,
				State:         dao.ContentStatePrepared,
				CreateAt:      time.Now(),
				UpdateAt:      time.Now(),
			})
			if err != nil {
				return fmt.Errorf("failed to create new content")
			}
		}
	}

	// create a new subscription
	subscription := &dao.Subscription{
		UserCredit:    userCredit,
		ChannelCredit: channelCredit,
		CreateAt:      time.Now(),
		UpdateAt:      time.Now(),
	}
	_, err = s.subscriptionMapper.Insert(subscription)
	if err != nil {
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
	pageSql := "SELECT * FROM t_content WHERE channel_credit IN (?) ORDER BY create_at DESC LIMIT ? OFFSET ?"
	return s.contentMapper.SelectBySQL(pageSql, channelCredits, pageSize, (pageIndex - 1) * pageSize)
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

// FetchContent fetches content for all subscriptions.
func (s *SubscribeService) FetchContent() (int, error) {
	// list all subscriptions of all users and related channels
	subscriptions, err := s.subscriptionMapper.Select(&dao.Subscription{})
	if err != nil {
		return 0, errors.Wrap(err, "failed to list subscriptions")
	}

	fetchCount := 0
	// loop through the channels and fetch the contents by `ChannelFetcher.Fetch()`
	for _, subscription := range subscriptions {
		result, err := s.fetcher.Fetch(fetcher.FetchOption{ChannelCredit: subscription.ChannelCredit})
		if err != nil {
			log.Errorf("failed to fetch content for channelCredit %s: %v", subscription.ChannelCredit, err)
			continue
		}

		// save the fetched contents to the database, with the state of `ContentStateInited`
		for _, content := range result.Contents {
			newContent := &dao.Content{
				Platform:      "YouTube",
				ChannelCredit: subscription.ChannelCredit,
				Title:         content.Title,
				Thumbnail:     content.Thumbnail,
				ContentCredit: content.Credit,
				State:         dao.ContentStatePrepared,
				CreateAt:      time.Now(),
				UpdateAt:      time.Now(),
			}
			oldContent, err := s.contentMapper.Select(&dao.Content{ContentCredit: content.Credit})
			if err != nil {
				return fetchCount, errors.Wrap(err, "failed to list content")
			}
			if len(oldContent) > 0 {
				log.Debugf("content %s already exists", content.Credit)
				continue
			}
			_, err = s.contentMapper.Insert(newContent)
			if err != nil {
				return fetchCount, errors.Wrap(err, "failed to insert content")
			}
			fetchCount++
		}
	}
	return fetchCount, nil
}

// DownloadContent downloads content for all subscriptions.
func (s *SubscribeService) DownloadContent() error {
	// list all the content with the state of `ContentStateInited`
	contents, err := s.contentMapper.Select(&dao.Content{State: dao.ContentStatePrepared})
	if err != nil {
		return errors.Wrap(err, "failed to list content")
	}
	log.Infof("going to download %d contents...", len(contents))

	// loop through the content and download the content by `Downloader.Download()`, and update the state of the content to `ContentStateDownloading`
	for _, content := range contents {
		content.State = dao.ContentStateDownloading
		_, err = s.contentMapper.Update(&dao.Content{ID: content.ID}, content)
		if err != nil {
			return errors.Wrap(err, "failed to update content state to downloading")
		}

		// download the content
		log.Infof("downloading content %s...", content.ContentCredit)
		result, err := s.downloader.Download(context.Background(), &downloader.DownloadOption{
			ContentCredit: content.ContentCredit,
			Rename:        content.Title,
			Format:        "mp3",
			Force:         false,
		})
		if err != nil {
			return errors.Wrap(err, "failed to download content")
		}

		// parse the result of the download and update the state of the content to `ContentStateDownloaded` if the download is successful
		if result.Finished {
			log.Infof("downloaded content %s: %s", content.ContentCredit, result.Output)
			content.State = dao.ContentStateDownloaded
			content.Path = result.Output
			_, err = s.contentMapper.Update(&dao.Content{ID: content.ID}, content)
			if err != nil {
				return errors.Wrap(err, "failed to update content state to downloaded")
			}
		} else {
			log.Errorf("failed to download content %s: %v", content.ContentCredit, result.Err)
		}

		log.Info("sleep 5 min before next download")
		time.Sleep(5 * time.Minute)
	}
	return nil
}
