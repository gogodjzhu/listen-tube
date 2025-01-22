package conf

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
	BinUri   string `yaml:"bin_uri"`
	BinURL   string `yaml:"bin_url"`
	BasePath string `yaml:"base_path"`
}
