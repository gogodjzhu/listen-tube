package fetcher

import (
	"testing"

	"github.com/gogodjzhu/listen-tube/internal/pkg/conf"
)

func TestChannelFetcher_Fetch(t *testing.T) {
	type fields struct {
		proxies []string
	}
	type args struct {
		opt FetchOption
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		validate func(*Result) bool
	}{
		{
			name: "Valid fetch",
			fields: fields{
				proxies: []string{},
			},
			args: args{
				opt: FetchOption{
					ChannelCredit: "UC_x5XG1OV2P6uZZ5FSM9Ttw",
				},
			},
			wantErr: false,
			validate: func(result *Result) bool {
				return len(result.Contents) > 0
			},
		},
		// Add more test cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf := NewFetcher(&conf.FetcherConfig{
				Enable:                true,
			})
			got, err := cf.Fetch(tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChannelFetcher.Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.validate(got) {
				t.Errorf("ChannelFetcher.Fetch() validation failed for result: %v", got)
			}
		})
	}
}

func TestFetcher_ParseChannelCredit(t *testing.T) {
	type fields struct {
		proxies []string
	}
	type args struct {
		channelCredit string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid URL",
			fields: fields{
				proxies: []string{},
			},
			args: args{
				channelCredit: "https://www.youtube.com/channel/UC_x5XG1OV2P6uZZ5FSM9Ttw",
			},
			want:    "UC_x5XG1OV2P6uZZ5FSM9Ttw",
			wantErr: false,
		},
		{
			name: "Valid URL",
			fields: fields{
				proxies: []string{},
			},
			args: args{
				channelCredit: "https://www.youtube.com/@GoogleDevelopers",
			},
			want:    "UC_x5XG1OV2P6uZZ5FSM9Ttw",
			wantErr: false,
		},
		{
			name: "Valid URL",
			fields: fields{
				proxies: []string{},
			},
			args: args{
				channelCredit: "@GoogleDevelopers",
			},
			want:    "UC_x5XG1OV2P6uZZ5FSM9Ttw",
			wantErr: false,
		},
		{
			name: "Invalid URL",
			fields: fields{
				proxies: []string{},
			},
			args: args{
				channelCredit: "https://www.youtube.com/channel/invalid_channel",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Invalid URL",
			fields: fields{
				proxies: []string{},
			},
			args: args{
				channelCredit: "invalid_channel",
			},
			want:    "",
			wantErr: true,
		},
		// Add more test cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf := &Fetcher{
				proxies: tt.fields.proxies,
			}
			got, err := cf.ParseChannelCredit(tt.args.channelCredit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetcher.ParseChannelCredit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Fetcher.ParseChannelCredit() = %v, want %v", got, tt.want)
			}
		})
	}
}
