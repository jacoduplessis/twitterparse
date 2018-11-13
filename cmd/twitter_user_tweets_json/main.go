package main

import (
	"encoding/json"
	"github.com/jacoduplessis/twitterparse"
	"log"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go USERNAME\n")
	}
	username := os.Args[1]

	url := twitterparse.URLFromUsername(username)

	bts, err := twitterparse.Fetch(url)
	if err != nil {
		log.Fatal(err)
	}

	tweets, err := twitterparse.ParseTweetsBytes(bts)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(tweets)
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(b)

}
