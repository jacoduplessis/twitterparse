package twitterparse

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dustin/go-humanize"
	"io"
	"strconv"
	"time"
)

type Tweet struct {
	ID             string `json:"id"`
	Permalink      string `json:"permalink"`
	Content        string `json:"content"`
	Timestamp      int64  `json:"timestamp"`
	ISOTime        string `json:"iso_time"`
	RelativeTime   string `json:"relative_time"`
	UserName       string `json:"user_name"`
	UserHandle     string `json:"user_handle"`
	UserID         string `json:"user_id"`
	UserAvatar     string `json:"user_avatar"`
	ImageURL       string `json:"image_url"`
	Context        string `json:"context"`
	Quoted         *Tweet `json:"quoted"`
	Video          bool   `json:"video"`
	VideoThumbnail string `json:"video_thumbnail"`
	VideoSource    string `json:"video_source"`
}

type ProfileTweet struct {
	CreatedAt string    `json:"created_at"`
	ID        string    `json:"id_string"`
	FullText  string    `json:"full_text"`
	Time      time.Time `json:"time"`
	UserID    string    `json:"user_id_str"`

	Entities struct {
		Media []struct {
			MediaURL string `json:"media_url_https"`
			Type     string `json:"type"`
		} `json:"media"`
	} `json:"entities"`
	ExtendedEntities struct {
		Media []struct {
			MediaURL  string `json:"media_url_https"`
			Type      string `json:"type"`
			VideoInfo struct {
				AspectRatio [2]int `json:"aspect_ratio"`
				Duration    int    `json:"duration_millis"`
				Variants    []struct {
					Bitrate     int    `json:"bitrate"`
					ContentType string `json:"content_type"`
					URL         string `json:"url"`
				} `json:"variants"`
			} `json:"video_info"`
		} `json:"media"`
	} `json:"extended_entities"`
}

type ProfileUser struct {
	ID              string `json:"id_str"`
	Name            string `json:"name"`
	ScreenName      string `json:"screen_name"`
	Location        string `json:"location"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url_https"`
	BannerImageURL  string `json:"profile_banner_url"`
}

type Profile struct {
	GlobalObjects struct {
		Tweets map[string]*ProfileTweet `json:"tweets"`
		Users  map[string]*ProfileUser  `json:"users"`
	} `json:"globalObjects"`
}

func ParseProfile(r io.Reader) (*Profile, error) {

	profile := &Profile{}
	err := json.NewDecoder(r).Decode(profile)
	if err != nil {
		return nil, err
	}

	for _, tweet := range profile.GlobalObjects.Tweets {

		tm, err := time.Parse(time.RubyDate, tweet.CreatedAt)
		if err != nil {
			fmt.Printf("error parsing date %s: %v", tweet.CreatedAt, err)
		}

		tweet.Time = tm

	}

	return profile, nil

}

func ParseTweets(r io.Reader) ([]*Tweet, error) {

	var tweets []*Tweet
	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return nil, err
	}

	doc.Find(".tweet").Each(func(i int, s *goquery.Selection) {

		timestampStr := s.Find(".js-short-timestamp").AttrOr("data-time", "")
		timestamp, _ := strconv.Atoi(timestampStr)
		timeValue := time.Unix(int64(timestamp), 0)
		relativeTime := humanize.Time(timeValue)

		media := s.Find(".AdaptiveMedia")

		video := media.HasClass("is-video")

		tweet := &Tweet{

			ID:           s.AttrOr("data-tweet-id", ""),
			Permalink:    s.AttrOr("data-permalink-path", ""),
			Content:      s.Find(".tweet-text").Text(),
			UserName:     s.AttrOr("data-name", ""),
			UserHandle:   s.AttrOr("data-screen-name", ""),
			UserID:       s.AttrOr("data-user-id", ""),
			UserAvatar:   s.Find(".avatar").AttrOr("src", ""),
			Timestamp:    int64(timestamp),
			RelativeTime: relativeTime,
			ISOTime:      timeValue.Format("2006-01-02T15:04:05Z"),
			ImageURL:     s.Find(".js-adaptive-photo").AttrOr("data-image-url", ""),
			Video:        video,
		}

		tweets = append(tweets, tweet)
	})

	return tweets, nil
}
