package fetcher

import (
	"testing"
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
				return len(result.Contents) > 0 && result.Err == nil
			},
		},
		// Add more test cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf := NewFetcher(Config{
				Proxies: tt.fields.proxies,
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