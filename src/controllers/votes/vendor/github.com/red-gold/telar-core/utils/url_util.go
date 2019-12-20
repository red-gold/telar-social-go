package utils

import (
	"fmt"

	coreConfig "github.com/red-gold/telar-core/config"
)

const baseFunctionURL = "/function"

// GetPrettyURL return
func GetPrettyURL() string {
	if *coreConfig.AppConfig.QueryPrettyURL {
		return ""
	}
	return baseFunctionURL
}

// GetPrettyURL formats according to pretty URL from (baseFunctionURL+url) and returns the resulting string.
func GetPrettyURLf(url string) string {
	if *coreConfig.AppConfig.QueryPrettyURL {
		return url
	}
	return fmt.Sprintf("%s%s", baseFunctionURL, url)
}
