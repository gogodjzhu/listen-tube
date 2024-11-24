package dao

import (
	"time"

	"github.com/gogodjzhu/listen-tube/internal/pkg/db"
)

type Channel struct {
	ID            uint      `gorm:"id;primaryKey;autoIncrement"`
	Platform      string    `gorm:"platform"`
	Name          string    `gorm:"name"`
	ChannelCredit string    `gorm:"channel_credit"`
	CreateAt      time.Time `gorm:"create_at"`
	UpdateAt      time.Time `gorm:"update_at"`
}

func (Channel) TableName() string {
	return "t_channel"
}

type ChannelMapper struct {
	*db.BasicMapper[Channel]
}

func NewChannelMapper(ds *db.DatabaseSource) (*ChannelMapper, error) {
	bm, err := db.NewBasicMapper[Channel](ds)
	if err != nil {
		return nil, err
	}
	return &ChannelMapper{
		bm,
	}, nil
}
