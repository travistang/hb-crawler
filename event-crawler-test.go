package main

import (
	"fmt"
	hb "hb-crawler/rating-gain/hiking-buddies"
	"log"
)

func failIfError(err error, msg string) {
	if err != nil {
		log.Fatal(msg)
	}
}

func main() {
	fmt.Printf("Starting to crawl hiking buddies...\n")
	credentials := hb.Credential{
		Email:    "opulent_umpires0w@icloud.com",
		Password: "fygveq-5ruqJa-gusgap",
	}

	fmt.Printf("Logging in with email %s...\n", credentials.Email)
	err, cookies := hb.Login(&credentials)

	failIfError(err, fmt.Sprintf("Error when logging in, error %+v\n", err))

	fmt.Printf("Fetching event list with cookies %s\n", *&cookies.SessionId)
	err, events := hb.CrawlEventList(cookies)

	failIfError(err, fmt.Sprintf("Failed to fetch events, error %+v\n", err))
	fmt.Printf("Events: %+v\n", events)
}
