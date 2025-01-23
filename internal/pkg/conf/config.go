package conf

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	WebConfig        *WebConfig        `yaml:"web"`
	DBConfig         *DBConfig         `yaml:"db"`
	DownloaderConfig *DownloaderConfig `yaml:"downloader"`
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

type DownloaderConfig struct {
	Enable   bool   `yaml:"enable"`
	BinUri   string `yaml:"bin_uri"`
	BinURL   string `yaml:"bin_url"`
	BasePath string `yaml:"base_path"`
}

func ReadConfig(content []byte) (*Config, error) {
	var config Config
	err := yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
