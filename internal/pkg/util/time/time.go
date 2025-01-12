package time

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

/**
* TranslateAccessibility2Duration formats time text to time.Duration
**/
func TranslateAccessibility2Duration(str string) (time.Duration, error) {
	str = strings.ToLower(str)
	str = strings.ReplaceAll(str, "ago", "")
	str = strings.TrimSpace(str)
	parts := strings.Split(str, " ")
	if len(parts) != 2 {
		return 10 * time.Minute, errors.New("invalid duration: " + str)
	}
	interval, err := strconv.Atoi(parts[0])
	if err != nil {
		return 10 * time.Minute, errors.New("invalid duration: " + str)
	}
	unit := parts[1]
	switch {
	case strings.Contains(unit, "second"):
		return time.Duration(interval) * time.Second, nil
	case strings.Contains(unit, "minute"):
		return time.Duration(interval) * time.Minute, nil
	case strings.Contains(unit, "hour"):
		return time.Duration(interval) * time.Hour, nil
	case strings.Contains(unit, "day"):
		return time.Duration(interval) * 24 * time.Hour, nil
	case strings.Contains(unit, "week"):
		return time.Duration(interval) * 7 * 24 * time.Hour, nil
	case strings.Contains(unit, "month"):
		return time.Duration(interval) * 30 * 24 * time.Hour, nil
	case strings.Contains(unit, "year"):
		return time.Duration(interval) * 365 * 24 * time.Hour, nil
	default:
		return time.Duration(rand.Int()) * time.Minute, errors.New("invalid duration: " + str)
	}
}

// TranslateDuration2Accessibility calculate related - now, and return the result in human readable format
func TranslateDuration2Accessibility(now time.Time, relatedTime time.Time) string {
	duration := now.Sub(relatedTime)
	if duration < 0 {
		return "in the future"
	}
	if duration < time.Minute {
		return "just now"
	}
	if duration < time.Hour {
		return strconv.Itoa(int(duration/time.Minute)) + " minutes ago"
	}
	if duration < 24*time.Hour {
		return strconv.Itoa(int(duration/time.Hour)) + " hours ago"
	}
	if duration < 7*24*time.Hour {
		return strconv.Itoa(int(duration/(24*time.Hour))) + " days ago"
	}
	if duration < 30*24*time.Hour {
		return strconv.Itoa(int(duration/(7*24*time.Hour))) + " weeks ago"
	}
	if duration < 365*24*time.Hour {
		return strconv.Itoa(int(duration/(30*24*time.Hour))) + " months ago"
	}
	return strconv.Itoa(int(duration/(365*24*time.Hour))) + " years ago"
}

/**
* TranslateDuration formats time text to time.Duration
**/
func TranslateDuration(str string) (time.Duration, error) {
	parts := strings.Split(str, ":")
	switch len(parts) {
	case 1:
		return time.ParseDuration(parts[0] + "s")
	case 2:
		return time.ParseDuration(parts[0] + "m" + parts[1] + "s")
	case 3:
		return time.ParseDuration(parts[0] + "h" + parts[1] + "m" + parts[2] + "s")
	default:
		return 0, errors.New("invalid duration format: " + str)
	}
}

// FormatDuration formats the duration into a string with the format: 01:00:10, 10:10, 00:10
func FormatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}