package dao

import (
	"time"

	"github.com/gogodjzhu/listen-tube/internal/pkg/db"
)

type Subscription struct {
	ID            uint      `gorm:"id;primaryKey;autoIncrement"`
	UserCredit    string    `gorm:"user_credit"`
	ChannelCredit string    `gorm:"channel_credit"`
	CreateAt      time.Time `gorm:"create_at"`
	UpdateAt      time.Time `gorm:"update_at"`
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
