package tube

import (
	"testing"
)

func TestNewChannelFetcherById(t *testing.T) {
	type args struct {
		channelID string
	}
	tests := []struct {
		name    string
		args    args
		want    func(fetcher *ChannelFetcher) bool
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				channelID: "UCBKDRq35-L8xev4O7ZqBeLg",
				//channelID: "UC_x5XG1OV2P6uZZ5FSM9Ttw",
			},
			want: func(fetcher *ChannelFetcher) bool {
				if len(fetcher.videos) == 0 {
					return false
				}
				if fetcher.innerContextStr == "" {
					return false
				}
				if fetcher.initialDataStr == "" {
					return false
				}
				return true
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChannelFetcherById(tt.args.channelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChannelFetcherById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && !tt.want(got) {
				t.Errorf("NewChannelFetcherById() miss got")
				return
			}
		})
	}
}
