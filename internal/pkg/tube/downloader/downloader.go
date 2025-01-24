package downloader

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gogodjzhu/listen-tube/internal/pkg/conf"
	"github.com/gogodjzhu/listen-tube/internal/pkg/db/dao"
	"github.com/gogodjzhu/listen-tube/internal/pkg/util/errors"
	"github.com/gogodjzhu/listen-tube/internal/pkg/util/ioutil"
	log "github.com/sirupsen/logrus"
)

// Downloader is responsible for downloading contents using youtube-dl. @see https://github.com/yt-dlp/yt-dlp
type Downloader struct {
	conf   *conf.DownloaderConfig
	binUri string
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
	d := &Downloader{
		conf:   conf,
		binUri: filepath.Join(conf.BasePath, ".bin", "yt-dlp"),
	}
	if err := d.prepare(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Downloader) TryStart(ctx context.Context, next func() *dao.Content, update func(dao.Content, *Result)) {
	if !d.conf.Enable {
		log.Info("downloader disabled")
		return
	}
	// periodically (from d.conf.DownloadIntervalSeconds) fetch the content
	timer := time.NewTicker(time.Duration(d.conf.DownloadIntervalSeconds) * time.Second)
	for {
		select {
		case <-ctx.Done():
			log.Info("downloader stopped")
			return
		case <-timer.C:
			content := next()
			if content == nil {
				continue
			}
			result, err := d.Download(ctx, &DownloadOption{
				ContentCredit: content.ContentCredit,
				Format:        "mp3",
				Force:         false,
			})
			result.Err = err
			update(*content, result)
		}
	}
}

func (d *Downloader) prepare() error {
	// execute ``yt-dlp --version`` to check if the binary is working
	cmd := exec.Command(d.binUri, "--version")
	if err := cmd.Run(); err != nil {
		log.Info("yt-dlp binary not found, downloading...")
		// download yt-dlp binary if not exists
		if err := ioutil.DownloadFile(d.conf.YtDlpLink, d.binUri, true, 0755); err != nil {
			log.Errorf("failed to download yt-dlp binary: %v", err)
			return errors.ErrFailedOS
		}
	}
	cmd = exec.Command(d.binUri, "--version")
	versionOutput, err := cmd.Output()
	if err != nil {
		log.Errorf("yt-dlp binary is not working: %v", err)
		return errors.ErrFailedOS
	}
	// log the version of yt-dlp binary and its path
	log.Infof("yt-dlp version: %s, path: %s", strings.TrimSpace(string(versionOutput)), d.binUri)

	// create base path if not exists
	if err := os.MkdirAll(d.conf.BasePath, os.ModePerm); err != nil {
		log.Errorf("failed to create base path: %v", err)
		return errors.ErrFailedOS
	}
	// check if ffmpeg exists
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		log.Errorf("ffmpeg not found: %v", err)
		return errors.ErrFailedOS
	}
	return nil
}

// Download downloads a content based on the provided DownloadOption and returns the Result.
func (d *Downloader) Download(ctx context.Context, opt *DownloadOption) (*Result, error) {
	// check if option is valid
	if err := opt.Validate(); err != nil {
		return nil, errors.ErrInvalidParams
	}

	// prepare the output file path
	outPath := filepath.Join(d.conf.BasePath, opt.ContentCredit)
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
					log.Debugf("downloading %s: %s", opt.ContentCredit, m)
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
		return nil, err
	}
	result.Finished = true
	log.Debugf("downloaded content: %s", opt.ContentCredit)
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
