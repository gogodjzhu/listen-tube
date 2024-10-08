package dao

import (
	"example.com/m/v2/pkg/db"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Channel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Platform  string    `gorm:"platform"`
	ChannelID string    `gorm:"channel_id"`
	Name      string    `gorm:"name"`
	CreateAt  time.Time `gorm:"create_at"`
	UpdateAt  time.Time `gorm:"update_at"`
}

func (Channel) TableName() string {
	return "t_channel"
}

type ChannelMapper struct {
	*db.DatabaseSource
}

func NewChannelMapper(db *db.DatabaseSource) (*ChannelMapper, error) {
	err := db.DB.AutoMigrate(&Channel{})
	if err != nil {
		return nil, err
	}
	return &ChannelMapper{
		db,
	}, nil
}

func (d *ChannelMapper) Insert(channel *Channel) (int64, error) {
	result := d.DB.Create(channel)
	return result.RowsAffected, result.Error
}

func (d *ChannelMapper) ListAll() ([]Channel, error) {
	var channels []Channel
	sql := fmt.Sprintf("SELECT * FROM %s", Channel{}.TableName())
	result := d.DB.Raw(sql).Scan(&channels)
	if result.Error != nil {
		return nil, result.Error
	}
	return channels, nil
}
