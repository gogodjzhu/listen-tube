package fetcher

import (
	"fmt"
	"strings"

	"github.com/gogodjzhu/listen-tube/internal/pkg/util/http"
	"github.com/tidwall/gjson"
)

type Fetcher struct {
	proxies []string
}

type Config struct {
	Proxies []string
}

func NewFetcher(config Config) *Fetcher {
	return &Fetcher{proxies: config.Proxies}
}

func (cf *Fetcher) Fetch(opt FetchOption) (*Result, error) {
	baseUrl := fmt.Sprintf("https://www.youtube.com/channel/%s/videos?view=0&flow=grid", opt.ChannelCredit)
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
	selectedSection := gjson.Get(initialDataStr, "contents.twoColumnBrowseResultsRenderer.tabs.#(tabRenderer.selected=true)")
	contentsRaws := gjson.Get(selectedSection.Raw, "tabRenderer.content.richGridRenderer.contents")
	for _, content := range contentsRaws.Array() {
		videoRenderer := gjson.Get(content.Raw, "richItemRenderer.content.videoRenderer")
		videoId := gjson.Get(videoRenderer.Raw, "videoId")
		title := gjson.Get(videoRenderer.Raw, "title.runs.0.text")
		thumbnail := gjson.Get(videoRenderer.Raw, "thumbnail.thumbnails.@reverse.0.url")
		contents = append(contents, Content{
			Credit:    videoId.Raw,
			Title:     title.Raw,
			Thumbnail: thumbnail.Raw,
		})
	}

	return &Result{
		Platform:      "youtube",
		ChannelCredit: opt.ChannelCredit,
		Contents:      contents,
		Err:           nil,
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
	Platform      string
	ChannelCredit string
	Contents      []Content
	Err           error
}

type Content struct {
	Credit    string
	Title     string
	Thumbnail string
}
