package dao

import "github.com/gogodjzhu/listen-tube/internal/pkg/db"

type UnionMapper struct {
	ChannelMapper      *ChannelMapper
	SubscriptionMapper *SubscriptionMapper
	ContentMapper      *ContentMapper
	UserMapper         *UserMapper
}

func NewUnionMapper(ds *db.DatabaseSource) (*UnionMapper, error) {
	cm, err := NewChannelMapper(ds)
	if err != nil {
		return nil, err
	}
	sm, err := NewSubscriptionMapper(ds)
	if err != nil {
		return nil, err
	}
	co, err := NewContentMapper(ds)
	if err != nil {
		return nil, err
	}
	um, err := NewUserMapper(ds)
	if err != nil {
		return nil, err
	}
	return &UnionMapper{
		ChannelMapper:      cm,
		SubscriptionMapper: sm,
		ContentMapper:      co,
		UserMapper:         um,
	}, nil
}