package twitterparse

import (
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

const TOKEN = "AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"

var cookieRegex = regexp.MustCompile(`document\.cookie.*?gt=(.*?);`)
var indexURL = "https://mobile.twitter.com"

type GuestToken string

type TwitterClient struct {
	Client http.Client
	Token  GuestToken
}

func NewClientWithHTTPClient(c http.Client) (*TwitterClient, error) {

	guestToken, err := initialRequest(c)
	if err != nil {
		return nil, err
	}

	if guestToken == "" {
		return nil, errors.New("Twitter Guest Token is missing")
	} else {
		fmt.Printf("Twitter Guest Token is %s\n", guestToken)
	}

	tc := &TwitterClient{
		Client: c,
		Token:  guestToken,
	}

	return tc, nil
}

func NewClient() (*TwitterClient, error) {

	c := http.Client{
		Timeout: time.Second * 10,
	}

	return NewClientWithHTTPClient(c)
}

func (tc *TwitterClient) GetProfile(userID string) (*Profile, error) {

	// also see https://api.twitter.com/2/timeline/media/%s.json

	u := fmt.Sprintf("https://api.twitter.com/2/timeline/profile/%s.json?include_profile_interstitial_type=1&include_blocking=1&include_blocked_by=1&include_followed_by=1&include_want_retweets=1&include_mute_edge=1&include_can_dm=1&include_can_media_tag=1&skip_status=1&cards_platform=Web-12&include_cards=1&include_composer_source=true&include_ext_alt_text=true&include_reply_count=1&tweet_mode=extended&include_entities=true&include_user_entities=true&include_ext_media_color=true&send_error_codes=true&include_tweet_replies=false&userId=%s&count=20&ext=mediaStats,highlightedLabel",
		userID, userID)

	r, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", TOKEN))
	r.Header.Add("Origin", "https://mobile.twitter.com")
	r.Header.Add("Referer", "https://mobile.twitter.com")
	r.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/11.0 Safari/605.1.15 Epiphany/605.1.15")
	r.Header.Add("x-twitter-active-user", "yes")
	r.Header.Add("x-twitter-client-language", "en")
	r.Header.Add("x-guest-token", string(tc.Token))

	resp, err := tc.Client.Do(r)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ParseProfile(resp.Body)

}

func (tc *TwitterClient) GetProfileTweets(userID string) ([]*Tweet, error) {

	var tweets []*Tweet

	profile, err := tc.GetProfile(userID)
	if err != nil {
		return nil, err
	}

	for _, twt := range profile.GlobalObjects.Tweets {

		t := &Tweet{}
		t.RelativeTime = humanize.Time(twt.Time)
		t.Timestamp = twt.Time.Unix()
		user, ok := profile.GlobalObjects.Users[twt.UserID]
		if !ok {
			fmt.Println("Could not find user in profile struct")
			continue
		}

		t.UserName = user.Name
		t.UserHandle = user.ScreenName
		t.UserAvatar = user.ProfileImageURL
		t.UserID = twt.UserID

		t.Content = twt.FullText

		for _, media := range twt.Entities.Media {
			if media.Type == "photo" {
				t.ImageURL = media.MediaURL
				break
			}
		}

		for _, media := range twt.ExtendedEntities.Media {
			if media.Type == "video" {
				t.Video = true
				t.VideoThumbnail = media.MediaURL
				for _, variant := range media.VideoInfo.Variants {

					if variant.Bitrate != 0 && variant.Bitrate > 400000 && variant.Bitrate < 1200000 {
						t.VideoSource = variant.URL
						break
					}
				}
				break
			}
		}

		tweets = append(tweets, t)
	}

	return tweets, nil

}

func initialRequest(client http.Client) (GuestToken, error) {

	guestToken := ""

	req1, err := http.NewRequest("GET", indexURL, nil)
	if err != nil {
		return "", err
	}

	req1.Header.Add("Origin", "https://mobile.twitter.com")
	req1.Header.Add("Referer", "https://mobile.twitter.com")
	req1.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/11.0 Safari/605.1.15 Epiphany/605.1.15")

	r1, err := client.Do(req1)
	if err != nil {
		return "", err
	}
	defer r1.Body.Close()
	body, err := ioutil.ReadAll(r1.Body)
	if err != nil {
		return "", err
	}

	matches := cookieRegex.FindSubmatch(body)

	if len(matches) > 1 {
		guestToken = string(matches[1])
	}

	return GuestToken(guestToken), nil

}
