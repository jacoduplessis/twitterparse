package twitterparse

import (
	"os"
	"testing"
)

func TestParseTweets(t *testing.T) {

	f, _ := os.Open("test/user.html")
	defer f.Close()
	tweets, _ := ParseTweets(f)

	twt := tweets[4]

	if twt.ID != "1062095461148966919" {
		t.Fail()
	}

	if twt.UserID != "14063426" {
		t.Fail()
	}

	if twt.Timestamp != "1542058265" {
		t.Fail()
	}

	if twt.UserHandle != "PGATOUR" {
		t.Fail()
	}

	if twt.UserName != "PGA TOUR" {
		t.Fail()
	}

	if twt.Content != "A win got him in. \n\nFrom Mexico to Hawaii. \nFrom @MayakobaGolf to @Sentry_TOC.\n\nMatt Kuchar is living the good life.pic.twitter.com/q4gR5NvtEp" {
		t.Fail()
	}

	if twt.Permalink != "/PGATOUR/status/1062095461148966919" {
		t.Fail()
	}

	if twt.UserAvatar != "https://pbs.twimg.com/profile_images/985892552439189505/VP47Cu-F_bigger.jpg" {
		t.Fail()
	}

	if twt.ImageURL != "https://pbs.twimg.com/media/Dr1SMVKWwAArMlw.jpg" {
		t.Fail()
	}
}
