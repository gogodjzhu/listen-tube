package tube

import (
	"example.com/m/v2/pkg/util"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Downloader struct {
	videoURL    string
	path        string
	filename    string
	messageChan util.ChanWriter
}

type DownloadResult struct {
}

func NewDownloader(videoId string) (*Downloader, error) {
	return &Downloader{
		videoURL:    fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoId),
		path:        fmt.Sprintf("/tmp/listen-tube/%s", videoId),
		filename:    "video",
		messageChan: make(util.ChanWriter),
	}, nil
}

func (d *Downloader) Download(force bool) error {
	defer close(d.messageChan)
	go func() {
		for msg := range d.messageChan {
			for _, m := range strings.Split(msg, "\n") {
				if len(strings.TrimSpace(m)) != 0 {
					fmt.Println(m)
				}
			}
		}
	}()

	if force {
		if err := os.RemoveAll(d.path); err != nil {
			return err
		}
	}
	cmd := exec.Command("/IdeaProjects/listen-tube/pkg/tmp/yt-dlp_linux", "--merge-output-format", "mp4", d.videoURL, "-o", d.path+"/"+d.filename)
	str := cmd.String()
	fmt.Println(str)
	cmd.Stdout = d.messageChan
	cmd.Stderr = d.messageChan
	return cmd.Run()
}
