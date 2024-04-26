package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/mariasalcedo/mariasalcedo/mastodon-notifier/client"
	"log"
	"os"
)

func main() {
	var msgArg, visibleArg string
	dryRunArg := flag.Bool("dry-run", false, "to validate credentials, default=false")
	oauthArg := flag.Bool("oauth", false, "to validate credentials, default=false")
	flag.StringVar(&msgArg, "message", "my test message", "message to toot, defaults to my test message")
	flag.StringVar(&visibleArg, "visibility", "unlisted", "visibility flag, default=unlisted")

	flag.Parse()

	c := func() *client.Client {
		switch *oauthArg {
		case true:
			return GetOauthClient()
		default:
			return GetTokenizedClient()
		}
	}()

	if *dryRunArg {
		if err := c.ValidateCredentials(context.Background()); err != nil {
			log.Fatal(err)
		}
		return
	}

	if !IsFlagPassed("message") {
		log.Fatalf("--message is required")
	}
	toot := client.TootRequest{
		Status:     msgArg,
		Visibility: visibleArg,
	}
	tootResponse, err := c.PostStatus(context.Background(), toot)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("url=%s id=%s", tootResponse.URL, tootResponse.ID)
}

func GetTokenizedClient() *client.Client {
	clientToken := RetrieveEnvVar("MASTODON_CLIENT_TOKEN")
	return client.NewClient(&client.Config{
		Server:      "https://mastodon.green",
		AccessToken: clientToken,
	})
}

func GetOauthClient() *client.Client {
	clientID := RetrieveEnvVar("MASTODON_CLIENT_KEY")
	clientSecret := RetrieveEnvVar("MASTODON_CLIENT_SECRET")
	clientCode := RetrieveEnvVar("MASTODON_CLIENT_ACCESS_CODE")

	c := client.NewClient(&client.Config{
		Server:       "https://mastodon.green",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		ClientCode:   clientCode,
	})
	if err := c.Authenticate(context.Background()); err != nil {
		log.Fatal(err)
	}
	return c
}

func IsFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func RetrieveEnvVar(key string) string {
	if value, ok := os.LookupEnv(key); !ok {
		log.Fatalf("%s not set\n", key)
	} else {
		if value == "" {
			log.Fatalf("%s is empty\n", key)
		}
	}
	return os.Getenv(key)
}
