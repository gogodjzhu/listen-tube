package conf

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	validYAML := `
web:
  port: 8080
  host: "localhost"
db:
  dsn: "user:password@/dbname"
  driver: "mysql"
downloader:
  enable: true
  bin_uri: "/usr/local/bin"
  bin_url: "http://example.com/bin"
  base_path: "/downloads"
`

	t.Run("valid config", func(t *testing.T) {
		config, err := ReadConfig([]byte(validYAML))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if config.WebConfig.Port != 8080 {
			t.Errorf("expected port 8080, got %d", config.WebConfig.Port)
		}
		if config.WebConfig.Host != "localhost" {
			t.Errorf("expected host localhost, got %s", config.WebConfig.Host)
		}
		if config.DBConfig.DSN != "user:password@/dbname" {
			t.Errorf("expected DSN user:password@/dbname, got %s", config.DBConfig.DSN)
		}
		if config.DBConfig.Driver != MySQLDriver {
			t.Errorf("expected driver mysql, got %s", config.DBConfig.Driver)
		}
		if !config.DownloaderConfig.Enable {
			t.Errorf("expected downloader to be enabled")
		}
		if config.DownloaderConfig.BinUri != "/usr/local/bin" {
			t.Errorf("expected bin_uri /usr/local/bin, got %s", config.DownloaderConfig.BinUri)
		}
		if config.DownloaderConfig.BinURL != "http://example.com/bin" {
			t.Errorf("expected bin_url http://example.com/bin, got %s", config.DownloaderConfig.BinURL)
		}
		if config.DownloaderConfig.BasePath != "/downloads" {
			t.Errorf("expected base_path /downloads, got %s", config.DownloaderConfig.BasePath)
		}
	})
}