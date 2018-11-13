package twitterparse

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func URLFromUsername(username string) string {
	return fmt.Sprintf("https://twitter.com/%s", strings.TrimSpace(username))
}

func FetchWithClient(client http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func Fetch(url string) ([]byte, error) {

	client := http.Client{
		Timeout: time.Second * 5,
	}
	return FetchWithClient(client, url)
}
