package twitterparse

import (
	"fmt"
	"net/http"
	"strings"
)

func URLFromUsername(username string) string {
	return fmt.Sprintf("https://twitter.com/%s", strings.TrimSpace(username))
}

func FetchUserWithClientAndParse(client http.Client, username string) ([]Tweet, error) {

	resp, err := client.Get(URLFromUsername(username))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ParseTweets(resp.Body)

}
