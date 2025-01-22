package downloader

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gogodjzhu/listen-tube/internal/pkg/conf"
	"github.com/gogodjzhu/listen-tube/internal/pkg/util/errors"
	"github.com/gogodjzhu/listen-tube/internal/pkg/util/ioutil"
	log "github.com/sirupsen/logrus"
)

// Downloader is responsible for downloading contents using youtube-dl. @see https://github.com/yt-dlp/yt-dlp
type Downloader struct {
	basePath string // base output directory
	binUri   string // path to youtube-dl binary
}

type Config struct {
	BinUri   string // path to store youtube-dl binary
	BinURL   string // url to download the youtube-dl binary
	BasePath string // base output directory
}

var DefaultConfig = &Config{
	BinUri:   "/tmp/listen-tube/.bin/yt-dlp",
	BinURL:   "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux",
	BasePath: "/tmp/listen-tube/",
}

func (opt *DownloadOption) Validate() error {
	if opt == nil {
		return errors.ErrInvalidParams
	}
	if len(opt.ContentCredit) == 0 {
		return errors.ErrInvalidParams
	}
	return nil
}

// NewDownloader creates a new Downloader instance and ensures the necessary binaries and directories are set up.
func NewDownloader(conf *conf.DownloaderConfig) (*Downloader, error) {
	// download the youtube-dl binary if not exists
	if _, err := os.Stat(conf.BinUri); err != nil {
		if !os.IsNotExist(err) {
			log.Errorf("failed to check youtube-dl binary: %v", err)
			return nil, errors.ErrFailedOS
		}
		if err := ioutil.DownloadFile(conf.BinURL, conf.BinUri); err != nil {
			log.Errorf("failed to download youtube-dl binary: %v", err)
			return nil, errors.ErrFailedOS
		}
		// make the binary executable
		if err := os.Chmod(conf.BinUri, 0755); err != nil {
			log.Errorf("failed to make the binary executable: %v", err)
			return nil, errors.ErrFailedOS
		}
	}
	// create base path if not exists
	if err := os.MkdirAll(conf.BasePath, os.ModePerm); err != nil {
		log.Errorf("failed to create base path: %v", err)
		return nil, errors.ErrFailedOS
	}
	// check if ffmpeg exists
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		log.Errorf("ffmpeg not found: %v", err)
		return nil, errors.ErrFailedOS
	}
	return &Downloader{
		basePath: conf.BasePath,
		binUri:   conf.BinUri,
	}, nil
}

// Download downloads a content based on the provided DownloadOption and returns the Result.
func (d *Downloader) Download(ctx context.Context, opt *DownloadOption) (*Result, error) {
	// check if option is valid
	if err := opt.Validate(); err != nil {
		return nil, errors.ErrInvalidParams
	}

	// prepare the output file path
	outPath := filepath.Join(d.basePath, opt.ContentCredit)
	if err := os.Mkdir(outPath, os.ModePerm); err != nil && !os.IsExist(err) {
		log.Errorf("failed to create output directory: %v", err)
		return nil, errors.ErrFailedOS
	}

	// clean the output directory if force download
	if opt.Force {
		if err := os.RemoveAll(outPath); err != nil {
			log.Errorf("failed to remove output directory: %v", err)
			return nil, errors.ErrFailedOS
		}
	}

	// prepare the Result struct
	result := &Result{
		Finished:   false,
		Err:        nil,
		Progress:   0,
		ContentURL: "https://www.youtube.com/watch?v=" + opt.ContentCredit,
		Output:     filepath.Join(outPath, "worstaudio."+opt.Format),
	}

	messageChan := make(ioutil.ChanWriter)
	defer close(messageChan)
	go func() {
		percentagePattern := regexp.MustCompile(`\[download\]\s+(\d+\.\d+|\d+)%\s+of\s+.+\s+in`)
		for msg := range messageChan {
			for _, m := range strings.Split(msg, "\n") {
				if len(strings.TrimSpace(m)) != 0 {
					match := percentagePattern.FindStringSubmatch(m)
					if len(match) > 1 {
						progress, _ := strconv.ParseFloat(match[1], 64)
						result.Progress = progress
					}
					log.Infof("downloading %s: %s", opt.ContentCredit, m)
				}
			}
		}
	}()

	// prepare the download command
	args := make([]string, 0)
	// if opt.Format != "" {
	// 	args = append(args, "--merge-output-format", opt.Format)
	// }
	args = append(args, "-f", "worstaudio")
	args = append(args, "-o", result.Output)
	args = append(args, result.ContentURL)
	cmd := exec.Command(d.binUri, args...)
	cmd.Stdout = messageChan
	cmd.Stderr = messageChan
	err := cmd.Run()
	if err != nil {
		log.Errorf("failed to download content: %v", err)
		result.Err = errors.ErrFailedOS
	}
	result.Finished = true
	return result, nil
}

type DownloadOption struct {
	ContentCredit string // content credit
	Format        string // download format
	Force         bool   // force download, delete the existing file
}

type Result struct {
	Finished   bool    // download finished
	Err        error   // error message
	Progress   float64 // download progress
	ContentURL string  // content url
	Output     string  // absolute output file path
}
