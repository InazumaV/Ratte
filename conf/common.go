package conf

import "strings"

func IsHttpUrl(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
