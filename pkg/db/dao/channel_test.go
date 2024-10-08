package dao

import (
	"example.com/m/v2/pkg/db"
	"github.com/google/uuid"
	"os"
	"testing"
	"time"
)

func TestDao(t *testing.T) {
	dsn := "/tmp/listen-tube-test-db.db"
	defer os.RemoveAll(dsn)
	databaseSource, err := db.NewDatabaseSource(&db.Config{
		DSN: dsn,
	})
	if err != nil {
		t.Errorf("NewDatabaseSource() failed, err:%v", err)
	}
	channelMapper, err := NewChannelMapper(databaseSource)
	if err != nil {
		t.Errorf("NewChannelMapper() failed, err:%v", err)
	}
	affected, err := channelMapper.Insert(&Channel{
		ID:        uuid.New(),
		Platform:  "testPlatform",
		ChannelID: "testChannelId",
		Name:      "testName",
		CreateAt:  time.Now(),
		UpdateAt:  time.Now(),
	})
	if err != nil {
		t.Errorf("channelMapper.Insert() failed, err:%v", err)
	}
	if affected != 1 {
		t.Errorf("channelMapper.Insert() affected != 1 but %d", affected)
	}
	channels, err := channelMapper.ListAll()
	if err != nil {
		t.Errorf("channelMapper.ListAll() failed, err:%v", err)
	}
	if len(channels) != 1 {
		t.Errorf("channelMapper.ListAll() len(result) != 1 but %d", affected)
	}

}
