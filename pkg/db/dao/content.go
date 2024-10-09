package dao

import (
	"example.com/m/v2/pkg/db"
	"time"
)

type Content struct {
	ID        uint      `gorm:"id;primaryKey;autoIncrement"`
	Platform  string    `gorm:"platform"`
	ChannelID string    `gorm:"channel_id"`
	Title     string    `gorm:"title"`
	CreateAt  time.Time `gorm:"create_at"`
	UpdateAt  time.Time `gorm:"update_at"`
}

func (Content) TableName() string {
	return "content"
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
