package dao

import (
	"time"

	"github.com/gogodjzhu/listen-tube/internal/pkg/db"
)

type User struct {
	ID       uint      `gorm:"id;primaryKey;autoIncrement"`
	Credit   string    `gorm:"credit"`
	Name     string    `gorm:"name"`
	CreateAt time.Time `gorm:"create_at"`
	UpdateAt time.Time `gorm:"update_at"`
}

func (User) TableName() string {
	return "t_user"
}

type UserMapper struct {
	*db.BasicMapper[User]
}

func NewUserMapper(ds *db.DatabaseSource) (*UserMapper, error) {
	bm, err := db.NewBasicMapper[User](ds)
	if err != nil {
		return nil, err
	}
	return &UserMapper{
		bm,
	}, nil
}
