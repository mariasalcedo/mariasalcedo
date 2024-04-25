package main

import (
	"context"
	"github.com/mattn/go-mastodon"
	"log"
	"os"
)

func PostAToot(Text string, Visibility string, DryRun bool) {
	c := mastodon.NewClient(&mastodon.Config{
		Server:       "https://mastodon.green",
		ClientID:     os.Getenv("MASTODON_CLIENT_KEY"),
		ClientSecret: os.Getenv("MASTODON_CLIENT_SECRET"),
	})

	err := c.AuthenticateApp(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if DryRun {
		_, e := c.VerifyAppCredentials(context.Background())

		if e != nil {
			log.Fatalf("%#v\n", e)
		}
	} else {
		toot := mastodon.Toot{
			Status:     Text,
			Visibility: Visibility,
		}
		_, err = c.PostStatus(context.Background(), &toot)

		if err != nil {
			log.Fatalf("%#v\n", err)
		}
	}
}

func main() {
	PostAToot("This is the content of my new post!", "public", true)
}
