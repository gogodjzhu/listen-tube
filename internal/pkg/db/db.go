package db

import "C"
import (
	"errors"
	"os"
	"path/filepath"

	"github.com/gogodjzhu/listen-tube/internal/pkg/conf"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	logger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DatabaseSource struct {
	DB *gorm.DB
}

type Pagenation struct {
	PageIndex int
	PageSize  int
	OrderBy   []string
	Order     string
}

func NewDatabaseSource(c *conf.DBConfig) (*DatabaseSource, error) {
	var db *gorm.DB
	var err error

	switch c.Driver {
	case conf.MySQLDriver:
		db, err = gorm.Open(mysql.Open(c.DSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	case conf.SQLiteDriver:
		if _, err := os.Stat(c.DSN); os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(c.DSN), os.ModePerm)
			if err != nil {
				return nil, err
			}
			file, err := os.Create(c.DSN)
			if err != nil {
				return nil, err
			}
			file.Close()
			log.Infof("Created sqlite db file: %s", c.DSN)
		}
		db, err = gorm.Open(sqlite.Open(c.DSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	default:
		return nil, errors.New("unsupported database driver")
	}

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

func (d *BasicMapper[T]) SelectWithPage(where *T, p *Pagenation) ([]*T, error) {
	var tArr []*T
	result := d.DB.Where(where).Limit(p.PageSize).Offset(p.PageIndex * p.PageSize).Find(&tArr)
	if result.Error != nil {
		return nil, result.Error
	}
	return tArr, nil
}

func (d *BasicMapper[T]) SelectBySQL(sql string, args ...interface{}) ([]*T, error) {
	var tArr []*T
	result := d.DB.Raw(sql, args...).Find(&tArr)
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
