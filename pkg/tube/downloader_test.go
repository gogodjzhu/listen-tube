package tube

import "testing"

func TestDownloader_Download(t *testing.T) {
	type args struct {
		video Video
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Downloader{}
			d.Download(tt.args.video)
		})
	}
}
