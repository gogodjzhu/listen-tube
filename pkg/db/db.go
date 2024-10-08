package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseSource struct {
	DB *gorm.DB
}

type Config struct {
	DSN string
}

func NewDatabaseSource(conf *Config) (*DatabaseSource, error) {
	db, err := gorm.Open(sqlite.Open(conf.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	dm := &DatabaseSource{
		DB: db,
	}
	return dm, nil
}
