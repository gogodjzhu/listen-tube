package subscribe

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/gogodjzhu/listen-tube/internal/pkg/db"
	"github.com/gogodjzhu/listen-tube/internal/pkg/db/dao"
	"github.com/gogodjzhu/listen-tube/internal/pkg/tube/downloader"
	"github.com/gogodjzhu/listen-tube/internal/pkg/tube/fetcher"
)

var fixedTime = time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)

func setupSuite(t *testing.T) func(t *testing.T) {
	// Teardown: delete test db files
	files, err := filepath.Glob("/tmp/listen-tube-unit-test-*.db")
	if err != nil {
		t.Fatalf("Failed to list test db files: %v", err)
	}
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			t.Fatalf("Failed to delete test db file %s: %v", file, err)
		}
	}
	// Return a function to teardown the test suite
	return func(t *testing.T) {
		// Teardown: delete test db files
		files, err := filepath.Glob("/tmp/listen-tube-unit-test-*.db")
		if err != nil {
			t.Fatalf("Failed to list test db files: %v", err)
		}
		for _, file := range files {
			if err := os.Remove(file); err != nil {
				t.Fatalf("Failed to delete test db file %s: %v", file, err)
			}
		}
	}
}

func setupTest(t *testing.T, s *SubscribeService) func(t *testing.T) {
	// Insert initial entities
	user1 := &dao.User{
		Credit:   "validUser1",
		Name:     "Test User1",
		CreateAt: fixedTime,
		UpdateAt: fixedTime,
	}
	_, err := s.userMapper.Insert(user1)
	if err != nil {
		t.Fatalf("Failed to insert user1: %v", err)
	}
	user2 := &dao.User{
		Credit:   "validUser2",
		Name:     "Test User2",
		CreateAt: fixedTime,
		UpdateAt: fixedTime,
	}
	_, err = s.userMapper.Insert(user2)
	if err != nil {
		t.Fatalf("Failed to insert user2: %v", err)
	}

	channel := &dao.Channel{
		Platform:      "YouTube",
		Name:          "Test Channel",
		ChannelCredit: "UC_x5XG1OV2P6uZZ5FSM9Ttw",
		CreateAt:      fixedTime,
		UpdateAt:      fixedTime,
	}
	_, err = s.channelMapper.Insert(channel)
	if err != nil {
		t.Fatalf("Failed to insert channel: %v", err)
	}

	subscription := &dao.Subscription{
		UserCredit:    "validUser1",
		ChannelCredit: "UC_x5XG1OV2P6uZZ5FSM9Ttw",
		CreateAt:      fixedTime,
		UpdateAt:      fixedTime,
	}
	_, err = s.subscriptionMapper.Insert(subscription)
	if err != nil {
		t.Fatalf("Failed to insert subscription: %v", err)
	}

	content := &dao.Content{
		Platform:      "YouTube",
		ChannelCredit: "UC_x5XG1OV2P6uZZ5FSM9Ttw",
		Title:         "Test Content",
		Thumbnail:     "http://example.com/thumbnail.jpg",
		ContentCredit: "dQw4w9WgXcQ",
		State:         dao.ContentStatePrepared,
		CreateAt:      fixedTime,
		UpdateAt:      fixedTime,
	}
	_, err = s.contentMapper.Insert(content)
	if err != nil {
		t.Fatalf("Failed to insert content: %v", err)
	}

	// Return a function to teardown the test
	return func(t *testing.T) {
		// Teardown: clean up the database
		s.userMapper.DB.Exec("DELETE FROM " + user1.TableName())
		s.channelMapper.DB.Exec("DELETE FROM " + channel.TableName())
		s.subscriptionMapper.DB.Exec("DELETE FROM " + subscription.TableName())
		s.contentMapper.DB.Exec("DELETE FROM " + content.TableName())
	}
}

func MockSubscribeService() *SubscribeService {
	conf := &db.Config{
		DSN: fmt.Sprintf("/tmp/listen-tube-unit-test-%d.db", time.Now().UnixNano()),
	}
	ds, err := db.NewDatabaseSource(conf)
	if err != nil {
		panic(err)
	}

	unionMapper, err := dao.NewUnionMapper(ds)
	if err != nil {
		panic(err)
	}

	downloaderConfig := &downloader.Config{
		BinUri:   "/tmp/listen-tube/.bin/yt-dlp",
		BinURL:   "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux",
		BasePath: "/tmp/listen-tube/",
	}
	downloaderInstance, err := downloader.NewDownloader(downloaderConfig)
	if err != nil {
		panic(err)
	}

	fetcherInstance := fetcher.NewFetcher(fetcher.Config{})

	subscribeService, err := NewSubscribeService(unionMapper, downloaderInstance, fetcherInstance)
	if err != nil {
		panic(err)
	}

	return subscribeService
}

func TestSubscribeService_AddSubscription(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name string
		args struct {
			userCredit    string
			channelCredit string
		}
		wantErr bool
	}{
		{
			name: "Valid subscription",
			args: struct {
				userCredit    string
				channelCredit string
			}{
				userCredit:    "validUser2",
				channelCredit: "UC_x5XG1OV2P6uZZ5FSM9Ttw",
			},
			wantErr: false,
		},
		{
			name: "Invalid subscription",
			args: struct {
				userCredit    string
				channelCredit string
			}{
				userCredit:    "invalidUser",
				channelCredit: "inUC_x5XG1OV2P6uZZ5FSM9Ttw",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MockSubscribeService()
			teardownTest := setupTest(t, s)
			defer teardownTest(t)

			if err := s.AddSubscription(tt.args.userCredit, tt.args.channelCredit); (err != nil) != tt.wantErr {
				t.Errorf("SubscribeService.AddSubscription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubscribeService_DeleteSubscription(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name string
		args struct {
			userCredit    string
			channelCredit string
		}
		wantErr bool
	}{
		{
			name: "Valid deletion",
			args: struct {
				userCredit    string
				channelCredit string
			}{
				userCredit:    "validUser1",
				channelCredit: "UC_x5XG1OV2P6uZZ5FSM9Ttw",
			},
			wantErr: false,
		},
		{
			name: "Invalid deletion",
			args: struct {
				userCredit    string
				channelCredit string
			}{
				userCredit:    "invalidUser",
				channelCredit: "inUC_x5XG1OV2P6uZZ5FSM9Ttw",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MockSubscribeService()
			teardownTest := setupTest(t, s)
			defer teardownTest(t)

			if err := s.DeleteSubscription(tt.args.userCredit, tt.args.channelCredit); (err != nil) != tt.wantErr {
				t.Errorf("SubscribeService.DeleteSubscription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubscribeService_ListSubscription(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name string
		args struct {
			userCredit string
		}
		want    []*dao.Subscription
		wantErr bool
	}{
		{
			name: "Valid list",
			args: struct {
				userCredit string
			}{
				userCredit: "validUser1",
			},
			want:    []*dao.Subscription{{ID: 1, UserCredit: "validUser1", ChannelCredit: "UC_x5XG1OV2P6uZZ5FSM9Ttw", CreateAt: fixedTime, UpdateAt: fixedTime}},
			wantErr: false,
		},
		{
			name: "Invalid list",
			args: struct {
				userCredit string
			}{
				userCredit: "invalidUser",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MockSubscribeService()
			teardownTest := setupTest(t, s)
			defer teardownTest(t)

			got, err := s.ListSubscription(tt.args.userCredit)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubscribeService.ListSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SubscribeService.ListSubscription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscribeService_FetchContent(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name    string
		wantMin int
		wantErr bool
	}{
		{
			name:    "Fetch content",
			wantMin: 1, // Adjust this based on expected minimum fetch count
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MockSubscribeService()
			teardownTest := setupTest(t, s)
			defer teardownTest(t)

			got, err := s.FetchContent()
			if (err != nil) != tt.wantErr {
				t.Errorf("SubscribeService.FetchContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got <= tt.wantMin {
				t.Errorf("SubscribeService.FetchContent() = %v, want > %v", got, tt.wantMin)
			}
		})
	}
}

func TestSubscribeService_DownloadContent(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name string
	}{
		{
			name: "Download content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MockSubscribeService()
			teardownTest := setupTest(t, s)
			defer teardownTest(t)

			s.DownloadContent()
		})
	}
}

func TestSubscribeService_ListContent(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name string
		args struct {
			userCredit string
		}
		want    []*dao.Content
		wantErr bool
	}{
		{
			name: "Valid list content",
			args: struct {
				userCredit string
			}{
				userCredit: "validUser1",
			},
			want: []*dao.Content{
				{
					ID:            1,
					Platform:      "YouTube",
					ChannelCredit: "UC_x5XG1OV2P6uZZ5FSM9Ttw",
					Title:         "Test Content",
					Thumbnail:     "http://example.com/thumbnail.jpg",
					ContentCredit: "dQw4w9WgXcQ",
					State:         dao.ContentStatePrepared,
					CreateAt:      fixedTime,
					UpdateAt:      fixedTime,
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid list content",
			args: struct {
				userCredit string
			}{
				userCredit: "invalidUser",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MockSubscribeService()
			teardownTest := setupTest(t, s)
			defer teardownTest(t)

			got, err := s.ListContent(tt.args.userCredit, 0, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubscribeService.ListContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SubscribeService.ListContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

	func TestConstant_list(t *testing.T) {
		if dao.ContentStateFailed != -1 {
			t.Errorf("ContentStateFailed = %v, want -1", dao.ContentStateFailed)
		}
		if dao.ContentStateInited != 0 {
			t.Errorf("ContentStateInited = %v, want 0", dao.ContentStateInited)
		}
		if dao.ContentStatePrepared != 1 {
			t.Errorf("ContentStatePrepared = %v, want 1", dao.ContentStatePrepared)
		}
		if dao.ContentStateDownloading != 2 {
			t.Errorf("ContentStateDownloading = %v, want 2", dao.ContentStateDownloading)
		}
		if dao.ContentStateDownloaded != 3 {
			t.Errorf("ContentStateDownloaded = %v, want 3", dao.ContentStateDownloaded)
		}
	}
