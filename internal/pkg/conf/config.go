package conf

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	WebConfig        *WebConfig        `yaml:"web"`
	DBConfig         *DBConfig         `yaml:"db"`
	SubscriberConfig *SubscriberConfig `yaml:"subscriber"`
}

type WebConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type DBConfig struct {
	DSN    string     `yaml:"dsn"`
	Driver DriverType `yaml:"driver"`
}

type DriverType string

const (
	MySQLDriver  DriverType = "mysql"
	SQLiteDriver DriverType = "sqlite"
)

type SubscriberConfig struct {
	FetcherConfig    *FetcherConfig    `yaml:"fetcher"`
	DownloaderConfig *DownloaderConfig `yaml:"downloader"`
}

type FetcherConfig struct {
	Enable                bool         `yaml:"enable"`
	ProxyConfig           *ProxyConfig `yaml:"proxy"`
	FetcheIntervalSeconds int          `yaml:"fetch_interval_seconds"`
}

type DownloaderConfig struct {
	Enable                  bool         `yaml:"enable"`
	ProxyConfig             *ProxyConfig `yaml:"proxy"`
	BasePath                string       `yaml:"base_path"`
	YtDlpLink               string       `yaml:"yt_dlp_link"`
	DownloadIntervalSeconds int          `yaml:"download_interval_seconds"`
}

type ProxyConfig struct {
	Proxies []string `yaml:"proxies"`
}

func ReadConfig(content []byte) (*Config, error) {
	var config Config
	err := yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
