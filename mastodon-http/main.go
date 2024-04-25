package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/mariasalcedo/mariasalcedo/mastodon/api"
	"log"
	"os"
)

func main() {
	var msgArg, visibleArg string
	dryRunArg := flag.Bool("dryrun", false, "to validate credentials, default=false")
	flag.StringVar(&msgArg, "message", "my test message", "message to toot, defaults to my test message")
	flag.StringVar(&visibleArg, "visibility", "public", "visibility flag, default=public")

	flag.Parse()

	clientID := retrieveEnvVar("MASTODON_CLIENT_KEY")
	clientSecret := retrieveEnvVar("MASTODON_CLIENT_SECRET")

	c := api.AuthenticatedClient(context.Background(), &api.Config{
		Server:       "https://mastodon.green",
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})

	c.ValidateCredentials(context.Background())

	if !*dryRunArg {
		if !isFlagPassed("message") {
			log.Fatalf("--message is required")
		}
		toot := api.TootRequest{
			Status:     msgArg,
			Visibility: visibleArg,
		}
		tootResponse := c.PostStatus(context.Background(), toot)
		fmt.Printf("url=%s id=%s", tootResponse.URL, tootResponse.ID)
	}
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func retrieveEnvVar(key string) string {
	if value, ok := os.LookupEnv(key); !ok {
		log.Fatalf("%s not set\n", key)
	} else {
		if value == "" {
			log.Fatalf("%s is empty\n", key)
		}
	}
	return os.Getenv(key)
}
