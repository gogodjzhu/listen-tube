package conf

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	content := []byte(`
web:
  port: 8080
  host: "localhost"
db:
  dsn: "user:password@/dbname"
  driver: "mysql"
subscriber:
  fetcher:
    enable: true
    proxy:
      proxies: ["http://proxy1", "http://proxy2"]
    fetch_interval_seconds: 60
  downloader:
    enable: true
    proxy:
      proxies: ["http://proxy1", "http://proxy2"]
    base_path: "/downloads"
    yt_dlp_link: "http://yt-dlp"
    download_interval_seconds: 120
`)

	config, err := ReadConfig(content)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.WebConfig.Port != 8080 {
		t.Errorf("Expected WebConfig.Port to be 8080, got %d", config.WebConfig.Port)
	}
	if config.WebConfig.Host != "localhost" {
		t.Errorf("Expected WebConfig.Host to be 'localhost', got %s", config.WebConfig.Host)
	}
	if config.DBConfig.DSN != "user:password@/dbname" {
		t.Errorf("Expected DBConfig.DSN to be 'user:password@/dbname', got %s", config.DBConfig.DSN)
	}
	if config.DBConfig.Driver != MySQLDriver {
		t.Errorf("Expected DBConfig.Driver to be 'mysql', got %s", config.DBConfig.Driver)
	}
	if !config.SubscriberConfig.FetcherConfig.Enable {
		t.Errorf("Expected SubscriberConfig.FetcherConfig.Enable to be true, got %v", config.SubscriberConfig.FetcherConfig.Enable)
	}
	if len(config.SubscriberConfig.FetcherConfig.ProxyConfig.Proxies) != 2 {
		t.Errorf("Expected 2 proxies, got %d", len(config.SubscriberConfig.FetcherConfig.ProxyConfig.Proxies))
	}
	if config.SubscriberConfig.FetcherConfig.FetcheIntervalSeconds != 60 {
		t.Errorf("Expected FetcherConfig.FetcheIntervalSeconds to be 60, got %d", config.SubscriberConfig.FetcherConfig.FetcheIntervalSeconds)
	}
	if !config.SubscriberConfig.DownloaderConfig.Enable {
		t.Errorf("Expected SubscriberConfig.DownloaderConfig.Enable to be true, got %v", config.SubscriberConfig.DownloaderConfig.Enable)
	}
	if config.SubscriberConfig.DownloaderConfig.BasePath != "/downloads" {
		t.Errorf("Expected DownloaderConfig.BasePath to be '/downloads', got %s", config.SubscriberConfig.DownloaderConfig.BasePath)
	}
	if config.SubscriberConfig.DownloaderConfig.YtDlpLink != "http://yt-dlp" {
		t.Errorf("Expected DownloaderConfig.YtDlpLink to be 'http://yt-dlp', got %s", config.SubscriberConfig.DownloaderConfig.YtDlpLink)
	}
	if config.SubscriberConfig.DownloaderConfig.DownloadIntervalSeconds != 120 {
		t.Errorf("Expected DownloaderConfig.DownloadIntervalSeconds to be 120, got %d", config.SubscriberConfig.DownloaderConfig.DownloadIntervalSeconds)
	}
}