package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Please provide username")
	}

	username := os.Args[1]

	r, err := http.Get("https://twitter.com/" + username)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	io.Copy(os.Stdout, r.Body)

}
