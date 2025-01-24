package fetcher

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gogodjzhu/listen-tube/internal/pkg/conf"
	"github.com/gogodjzhu/listen-tube/internal/pkg/db/dao"
	"github.com/gogodjzhu/listen-tube/internal/pkg/util/http"
	utiltime "github.com/gogodjzhu/listen-tube/internal/pkg/util/time"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Fetcher struct {
	proxies []string
	conf 	*conf.FetcherConfig
}

func NewFetcher(config *conf.FetcherConfig) *Fetcher {
	var proxies []string
	if config.ProxyConfig != nil {
		proxies = config.ProxyConfig.Proxies
	}
	return &Fetcher{
		proxies: proxies,
		conf:    config,
	}
}

func (cf *Fetcher) TryStart(ctx context.Context, next func() *dao.Channel, update func(*dao.Channel, *Result)) {
	if !cf.conf.Enable {
		log.Info("fetcher disabled")
		return
	}
	// periodically (from d.conf.DownloadIntervalSeconds) fetch the content
	timer := time.NewTicker(time.Duration(cf.conf.FetcheIntervalSeconds) * time.Second)
	for {
		select {
		case <-ctx.Done():
			log.Info("fetcher stopped")
			return
		case <-timer.C:
			channel := next()
			if channel == nil {
				continue
			}
			result, err := cf.Fetch(FetchOption{
				ChannelCredit: channel.ChannelCredit,
			})
			result.Err = err
			update(channel, result)
		}
	}
}

// ParseChannelCredit parse channel credit from channel credit string
func (cf *Fetcher) ParseChannelCredit(channelCredit string) (string, error) {
	channelCredit = strings.TrimSpace(channelCredit)
	if !strings.HasPrefix(channelCredit, "https://") && !strings.HasPrefix(channelCredit, "http://") {
		if strings.HasPrefix(channelCredit, "@") {
			channelCredit = "https://www.youtube.com/" + channelCredit
		} else {
			channelCredit = "https://www.youtube.com/channel/" + channelCredit
		}
	}
	html, err := http.HttpGet(cf.proxies, channelCredit)
	if err != nil {
		return "", err
	}
	initialDataStr, err := getTextFromHtml(html, "var ytInitialData = ", 0, "};")
	if err != nil {
		return "", err
	}
	initialDataStr += "}"
	externalId := gjson.Get(initialDataStr, "metadata.channelMetadataRenderer.externalId")
	return externalId.Str, nil
}

func (cf *Fetcher) Fetch(opt FetchOption) (*Result, error) {
	channelID, err := cf.ParseChannelCredit(opt.ChannelCredit)
	if err != nil {
		return nil, err
	}
	baseUrl := fmt.Sprintf("https://www.youtube.com/channel/%s/videos?view=0&flow=grid", channelID)
	html, err := http.HttpGet(cf.proxies, baseUrl)
	if err != nil {
		return nil, err
	}

	// innerContextStr := getTextFromHtml(html, "INNERTUBE_CONTEXT", 2, "\"}},") + "\"}}"
	initialDataStr, err := getTextFromHtml(html, "var ytInitialData = ", 0, "};")
	if err != nil {
		return nil, err
	}
	initialDataStr += "}"

	var contents []Content
	selectedSection := gjson.Get(initialDataStr, "contents.twoColumnBrowseResultsRenderer.tabs.#(tabRenderer.title=Videos)")
	contentsRaws := gjson.Get(selectedSection.Raw, "tabRenderer.content.richGridRenderer.contents")
	for _, content := range contentsRaws.Array() {
		continuationItemRenderer := gjson.Get(content.Raw, "continuationItemRenderer")
		if continuationItemRenderer.Exists() {
			continue
		}
		videoRenderer := gjson.Get(content.Raw, "richItemRenderer.content.videoRenderer")
		videoId := gjson.Get(videoRenderer.Raw, "videoId")
		title := gjson.Get(videoRenderer.Raw, "title.runs.0.text")
		if videoId.Str == "" || title.Str == "" {
			log.Warnf("videoId or title is empty, skip. channel: %s", opt.ChannelCredit)
			continue
		}
		thumbnail := gjson.Get(videoRenderer.Raw, "thumbnail.thumbnails.@reverse.0.url")
		thumbnailStr := strings.Split(thumbnail.Str, "?")[0]
		thumbnailStr = regexp.MustCompile(`hqdefault_custom_[0-9]+\.jpg`).ReplaceAllString(thumbnailStr, "hqdefault.jpg")
		publishedTimeText := gjson.Get(videoRenderer.Raw, "publishedTimeText.simpleText")
		publishedTime, err := utiltime.TranslateAccessibility2Duration(publishedTimeText.Str)
		if err != nil {
			log.Warnf("failed to parse published time: %s", publishedTimeText.Str)
		}
		lengthText := gjson.Get(videoRenderer.Raw, "lengthText.simpleText")
		length, err := utiltime.TranslateDuration(lengthText.Str)
		if err != nil {
			log.Warnf("failed to parse length: %s", lengthText.Str)
		}
		membersOnly := false
		badges := gjson.Get(videoRenderer.Raw, "badges")
		if badges.Exists() {
			for _, badge := range badges.Array() {
				if gjson.Get(badge.Raw, "metadataBadgeRenderer.label").Str == "Members only" {
					membersOnly = true
					break
				}
			}
		}
		contents = append(contents, Content{
			Credit:        videoId.Str,
			Title:         title.Str,
			Thumbnail:     thumbnailStr,
			PublishedTime: time.Now().Add(-publishedTime),
			Length:        length,
			MembersOnly:   membersOnly,
		})
	}
	metadata := gjson.Get(initialDataStr, "metadata.channelMetadataRenderer")
	title := gjson.Get(metadata.Raw, "title")
	description := gjson.Get(metadata.Raw, "description")
	var ownerUrls []string
	for _, ownerUrl := range gjson.Get(metadata.Raw, "ownerUrls").Array() {
		ownerUrls = append(ownerUrls, ownerUrl.Str)
	}
	var thumbnails []string
	for _, thumbnail := range gjson.Get(metadata.Raw, "avatar.thumbnails").Array() {
		t := gjson.Get(thumbnail.Raw, "url")
		thumbnails = append(thumbnails, t.Str)
	}
	// TODO: ignore member

	return &Result{
		Platform:    "youtube",
		ChannelID:   opt.ChannelCredit,
		Title:       title.Str,
		Description: description.Str,
		Thumbnails:  thumbnails,
		OwnerUrls:   ownerUrls,
		Contents:    contents,
		Err:         nil,
	}, nil
}

func getTextFromHtml(html, key string, numChars int, stop string) (string, error) {
	posBegin := strings.Index(html, key)
	if posBegin == -1 {
		return "", fmt.Errorf("key %s not found in html", key)
	}
	posBegin += len(key) + numChars

	posEnd := strings.Index(html[posBegin:], stop)
	if posEnd == -1 {
		return "", fmt.Errorf("stop string %s not found in html after key %s", stop, key)
	}
	posEnd += posBegin

	return html[posBegin:posEnd], nil
}

type FetchOption struct {
	ChannelCredit string
}

type Result struct {
	Platform    dao.Platform
	ChannelID   string
	Title       string
	Description string
	Thumbnails  []string
	OwnerUrls   []string
	Contents    []Content
	Err         error
}

type Content struct {
	Credit        string
	Title         string
	Thumbnail     string
	PublishedTime time.Time
	Length        time.Duration
	MembersOnly   bool
}
