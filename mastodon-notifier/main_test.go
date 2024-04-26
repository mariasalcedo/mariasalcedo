package main

import (
	"context"
	"fmt"
	"github.com/mariasalcedo/mariasalcedo/mastodon-notifier/client"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPostStatus(t *testing.T) {
	tests := []struct {
		name           string
		toot           client.TootRequest
		mockResponse   string
		mockStatusCode int
		wantErr        bool
	}{
		{
			name: "successful post",
			toot: client.TootRequest{
				Status:     "Test status",
				Visibility: "public",
			},
			mockResponse:   `{"id":"123","url":"http://mastodon.example.com/status/123","content":"Test status","created_at":"2024-04-25T00:00:00Z","visibility":"public"}`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "error from server",
			toot: client.TootRequest{
				Status:     "Test status",
				Visibility: "public",
			},
			mockResponse:   `{"error":"Internal server error"}`,
			mockStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.mockStatusCode)
				fmt.Fprintln(w, tc.mockResponse)
			}))
			defer mockServer.Close()

			config := &client.Config{
				Server:       mockServer.URL,
				ClientID:     "clientID",
				ClientSecret: "clientSecret",
				AccessToken:  "access-token",
			}

			c := &client.Client{
				Client: http.Client{Timeout: 10 * time.Second},
				Config: config,
			}

			resp, err := c.PostStatus(context.Background(), tc.toot)

			if (err != nil) != tc.wantErr {
				t.Errorf("postStatus() error = %v, wantErr %v", err, tc.wantErr)
			}
			log.Print("ID: " + resp.ID + " URL: " + resp.URL + " status: " + resp.Content + " visibility: " + resp.Visibility)
		})
	}
}
