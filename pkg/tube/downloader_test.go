package tube

import (
	"testing"
)

func TestDownloader_Download(t *testing.T) {
	type fields struct {
		videoId string
	}
	type args struct {
		force bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				videoId: "PagHbCgW7Uo",
			},
			args: args{
				force: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDownloader(tt.fields.videoId)
			if err != nil {
				t.Errorf("NewDownloader() error = %v", err)
			}
			if err := d.Download(tt.args.force); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
