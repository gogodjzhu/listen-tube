package str

import (
	"regexp"
	"strings"
)

// ArrayToString converts an array of strings to a single string.
func ArrayToStringWithSplit(arr []string, split string) string {
	if len(arr) == 0 {
		return ""
	}

	result := arr[0]
	for i := 1; i < len(arr); i++ {
		result += split + arr[i]
	}
	return result
}

// StringToArray converts a string to an array of strings.
func StringToArrayWithSplit(str, split string) []string {
	if str == "" {
		return []string{}
	}
	return strings.Split(str, split)
}

// ReplaceAllRegex replaces the matched regex in a string with another string.
func ReplaceAllRegex(str, regex, replace string) string {
	re := regexp.MustCompile(regex)
	return re.ReplaceAllString(str, replace)
}