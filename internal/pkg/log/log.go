package log

import (
	"github.com/sirupsen/logrus"
  "github.com/lestrrat-go/file-rotatelogs"
  "io"
  "os"
  "time"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:           "2006-01-02 15:04:05",
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		DisableLevelTruncation:    false,
	})
	rotateWriter, err := rotatelogs.New(
		"logs/rotate/logfile.%Y%m%d.log",
		rotatelogs.WithLinkName("logs/logfile"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithRotationSize(1024*1024*1024),
	)
	if err != nil {
		panic(err)
	}
	writers := []io.Writer{rotateWriter, os.Stdout}
	logrus.SetOutput(io.MultiWriter(writers...))
}
