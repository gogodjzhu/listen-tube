package db

import "C"
import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	logger "gorm.io/gorm/logger"
)

type DatabaseSource struct {
	DB *gorm.DB
}

type Config struct {
	DSN string
}

func NewDatabaseSource(conf *Config) (*DatabaseSource, error) {
	db, err := gorm.Open(sqlite.Open(conf.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	dm := &DatabaseSource{
		DB: db,
	}
	return dm, nil
}

type BasicMapper[T schema.Tabler] struct {
	*DatabaseSource
}

func NewBasicMapper[T schema.Tabler](ds *DatabaseSource) (*BasicMapper[T], error) {
	var t T
	err := ds.DB.AutoMigrate(&t)
	if err != nil {
		return nil, err
	}
	return &BasicMapper[T]{
		ds,
	}, nil
}

func (d *BasicMapper[T]) Insert(t *T) (int64, error) {
	result := d.DB.Create(t)
	return result.RowsAffected, result.Error
}

func (d *BasicMapper[T]) Select(where *T) ([]*T, error) {
	var tArr []*T
	result := d.DB.Where(where).Find(&tArr)
	if result.Error != nil {
		return nil, result.Error
	}
	return tArr, nil
}

func (d *BasicMapper[T]) Delete(t *T) (int64, error) {
	result := d.DB.Delete(t)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return result.RowsAffected, nil
}

func (d *BasicMapper[T]) Update(old *T, new *T) (int64, error) {
	result := d.DB.Model(old).Updates(new)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return result.RowsAffected, nil
}