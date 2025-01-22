package downloader

import (
	"context"
	"reflect"
	"testing"

	"github.com/gogodjzhu/listen-tube/internal/pkg/conf"
)

func TestDownloader_Download(t *testing.T) {
	type fields struct {
		binUri   string
		binURL   string
		basePath string
	}
	type args struct {
		ctx context.Context
		opt *DownloadOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Result
		wantErr bool
	}{
		{
			name: "Valid download",
			fields: fields{
				binUri:   "/tmp/listen-tube/.bin/yt-dlp",
				binURL:   "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux",
				basePath: "/tmp/listen-tube/",
			},
			args: args{
				ctx: context.Background(),
				opt: &DownloadOption{
					ContentCredit: "dQw4w9WgXcQ",
					Format:        "mp4",
					Force:         true,
				},
			},
			want: &Result{
				Finished:   true,
				Err:        nil,
				Progress:   100,
				ContentURL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
				Output:     "/tmp/listen-tube/dQw4w9WgXcQ/worstaudio.mp4",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDownloader(&conf.DownloaderConfig{
				BinUri:   tt.fields.binUri,
				BinURL:   tt.fields.binURL,
				BasePath: tt.fields.basePath,
			})
			if err != nil {
				t.Errorf("NewDownloader() error = %v", err)
				return
			}
			got, err := d.Download(tt.args.ctx, tt.args.opt)

			if (err != nil) != tt.wantErr {
				t.Errorf("Downloader.Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Downloader.Download() = %v, want %v", got, tt.want)
			}
		})
	}
}
