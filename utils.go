package twitterparse

import (
	"fmt"
	"strings"
)

func URLFromUsername(username string) string {
	return fmt.Sprintf("https://twitter.com/%s", strings.TrimSpace(username))
}
