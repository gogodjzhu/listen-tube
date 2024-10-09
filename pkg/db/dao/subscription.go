package dao

import (
	"example.com/m/v2/pkg/db"
	"time"
)

type Subscription struct {
	ID        uint      `gorm:"id;primaryKey;autoIncrement"`
	UserId    string    `gorm:"user_id"`
	ChannelId string    `gorm:"channel_id"`
	CreateAt  time.Time `gorm:"create_at"`
	UpdateAt  time.Time `gorm:"update_at"`
}

func (Subscription) TableName() string {
	return "t_subscription"
}

type SubscriptionMapper struct {
	*db.BasicMapper[Subscription]
}

func NewSubscriptionMapper(ds *db.DatabaseSource) (*SubscriptionMapper, error) {
	bm, err := db.NewBasicMapper[Subscription](ds)
	if err != nil {
		return nil, err
	}
	return &SubscriptionMapper{
		bm,
	}, nil
}
