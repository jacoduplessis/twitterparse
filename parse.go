package twitterparse

import (
	"github.com/PuerkitoBio/goquery"
	"io"
)

type Tweet struct {
	ID         string `json:"id"`
	Permalink  string `json:"permalink"`
	Content    string `json:"content"`
	Timestamp  string `json:"timestamp"`
	UserName   string `json:"user_name"`
	UserHandle string `json:"user_handle"`
	UserID     string `json:"user_id"`
	UserAvatar string `json:"user_avatar"`
	ImageURL   string `json:"image_url"`
}

func ParseTweets(r io.Reader) ([]Tweet, error) {

	var tweets []Tweet
	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return nil, err
	}

	doc.Find(".tweet").Each(func(i int, s *goquery.Selection) {

		tweet := Tweet{

			ID:         s.AttrOr("data-tweet-id", ""),
			Permalink:  s.AttrOr("data-permalink-path", ""),
			Content:    s.Find(".tweet-text").Text(),
			UserName:   s.AttrOr("data-name", ""),
			UserHandle: s.AttrOr("data-screen-name", ""),
			UserID:     s.AttrOr("data-user-id", ""),
			UserAvatar: s.Find(".avatar").AttrOr("src", ""),
			Timestamp:  s.Find(".js-short-timestamp").AttrOr("data-time", ""),
			ImageURL:   s.Find(".js-adaptive-photo").AttrOr("data-image-url", ""),
		}

		tweets = append(tweets, tweet)
	})

	return tweets, nil
}
