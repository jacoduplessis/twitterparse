package main

import (
	"encoding/json"
	"github.com/jacoduplessis/twitterparse"
	"log"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go USER_ID\n")
	}
	username := os.Args[1]

	tc, err := twitterparse.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	tweets, err := tc.GetProfileTweets(username)

	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(tweets)
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(b)

}
