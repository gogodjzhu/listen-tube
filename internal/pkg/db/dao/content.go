package dao

import (
	"time"

	"github.com/gogodjzhu/listen-tube/internal/pkg/db"
)

type Content struct {
	ID            uint         `gorm:"id;primaryKey;autoIncrement"`
	Platform      string       `gorm:"platform"`
	ChannelCredit string       `gorm:"channel_credit"`
	Title         string       `gorm:"title"`
	Thumbnail     string       `gorm:"thumbnail"`
	ContentCredit string       `gorm:"content_credit"`
	State         ContentState `gorm:"state"`
	Path          string       `gorm:"path"`
	CreateAt      time.Time    `gorm:"create_at"`
	UpdateAt      time.Time    `gorm:"update_at"`
}

type ContentState int

const (
	ContentStateInited ContentState = iota
	ContentStateDownloading
	ContentStateDownloaded
)

func (Content) TableName() string {
	return "t_content"
}

type ContentMapper struct {
	*db.BasicMapper[Content]
}

func NewContentMapper(ds *db.DatabaseSource) (*ContentMapper, error) {
	bm, err := db.NewBasicMapper[Content](ds)
	if err != nil {
		return nil, err
	}
	return &ContentMapper{
		bm,
	}, nil
}
