package tube

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

type ChannelFetcher struct {
	proxies         []string
	innerContextStr string
	initialDataStr  string
	videos          []Video
}

type Video struct {
	ID        string
	Title     string
	Thumbnail string
}

func NewChannelFetcherById(channelID string) (*ChannelFetcher, error) {
	baseUrl := fmt.Sprintf("https://www.youtube.com/channel/%s/videos?view=0&flow=grid", channelID)
	client := getSession(nil) // TODO support proxies

	//var err error
	html := httpGet(client, baseUrl)
	innerContextStr := getTextFromHtml(html, "INNERTUBE_CONTEXT", 2, "\"}},") + "\"}}"
	initialDataStr := getTextFromHtml(html, "var ytInitialData = ", 0, "};") + "}"
	var videos []Video

	selectedSection := gjson.Get(initialDataStr, "contents.twoColumnBrowseResultsRenderer.tabs.#(tabRenderer.selected=true)")
	contents := gjson.Get(selectedSection.Raw, "tabRenderer.content.richGridRenderer.contents")
	for _, content := range contents.Array() {
		videoRenderer := gjson.Get(content.Raw, "richItemRenderer.content.videoRenderer")
		videoId := gjson.Get(videoRenderer.Raw, "videoId")
		title := gjson.Get(videoRenderer.Raw, "title.runs.0.text")
		thumbnail := gjson.Get(videoRenderer.Raw, "thumbnail.thumbnails.@reverse.url") // @reverse to get the last one (the largest one)
		videos = append(videos, Video{
			ID:        videoId.Raw,
			Title:     title.Raw,
			Thumbnail: thumbnail.Raw,
		})
	}

	fmt.Println(fmt.Sprintf("%+v", videos))
	return &ChannelFetcher{
		innerContextStr: innerContextStr,
		initialDataStr:  initialDataStr,
		videos:          videos,
	}, nil
}

func getSession(proxies []string) *http.Client {
	client := &http.Client{}
	// Set up proxies if provided
	if len(proxies) > 0 {
		randomIdx := rand.Intn(len(proxies))
		proxyURL, _ := url.Parse(proxies[randomIdx])
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}
	return client
}

func httpGet(client *http.Client, url string) string {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en")
	req.AddCookie(&http.Cookie{Name: "CONSENT", Value: "YES+cb", Domain: ".youtube.com"})

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func getTextFromHtml(html, key string, numChars int, stop string) string {
	posBegin := strings.Index(html, key) + len(key) + numChars
	posEnd := strings.Index(html[posBegin:], stop) + posBegin
	return html[posBegin:posEnd]
}
