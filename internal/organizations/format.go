package organizations

import (
	"fmt"
	"strings"
)

// formatURL takes a string with a user-inputted URL and returns a standardized, formatted URL which contains a protocol if one does not already exist.
func formatURL(url string) string {
	if len(url) < 1 {
		return ""
	}

	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}

	return fmt.Sprintf("http://%s", url)
}
