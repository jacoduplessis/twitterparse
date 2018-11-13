package twitterparse

import "net/http"

func FetchUserWithClientAndParse(c http.Client, username string) ([]Tweet, error) {

	b, err := FetchWithClient(c, URLFromUsername(username))
	if err != nil {
		return nil, err
	}
	return ParseTweetsBytes(b)

}
