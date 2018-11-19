package twitterparse

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/dustin/go-humanize"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Tweet struct {
	ID             string `json:"id"`
	Permalink      string `json:"permalink"`
	Content        string `json:"content"`
	Timestamp      int    `json:"timestamp"`
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
			Timestamp:    timestamp,
			RelativeTime: relativeTime,
			ISOTime:      timeValue.Format("2006-01-02T15:04:05Z"),
			ImageURL:     s.Find(".js-adaptive-photo").AttrOr("data-image-url", ""),
			Video:        video,
		}

		if video {

			playerStyle := s.Find(".PlayableMedia-player").AttrOr("style", "")

			tweet.VideoThumbnail = getVideoThumbnail(playerStyle)

			// videoId := getVideoID(playerStyle)
			// tweet.VideoSource = fmt.Sprintf("https://video.twimg.com/tweet_video/%s.mp4", videoId)
			// https://api.twitter.com/2/timeline/profile/55246492.json?include_profile_interstitial_type=1&include_blocking=1&include_blocked_by=1&include_followed_by=1&include_want_retweets=1&include_mute_edge=1&include_can_dm=1&include_can_media_tag=1&skip_status=1&cards_platform=Web-12&include_cards=1&include_composer_source=true&include_ext_alt_text=true&include_reply_count=1&tweet_mode=extended&include_entities=true&include_user_entities=true&include_ext_media_color=true&send_error_codes=true&include_tweet_replies=false&userId=55246492&count=20&ext=mediaStats%2ChighlightedLabel
		}

		tweets = append(tweets, tweet)
	})

	return tweets, nil
}

func getVideoID(style string) string {

	reID := regexp.MustCompile(`/(.*?)\.jpg`)
	idMatches := reID.FindStringSubmatch(style)
	if len(idMatches) < 2 {
		return ""
	}
	return idMatches[1]
}

func getVideoThumbnail(style string) string {
	reID := regexp.MustCompile(`url\((.*)\)`)
	idMatches := reID.FindStringSubmatch(style)
	if len(idMatches) < 2 {
		return ""
	}
	return strings.Trim(idMatches[1], "\"'")
}
