package client

import (
	"io"
	"net/http"
)

type Config struct {
	Server       string
	ClientID     string
	ClientSecret string
	ClientCode   string
	AccessToken  string
}

type Client struct {
	http.Client
	Config *Config
}

type MastodonClientRequest struct {
	RequestBody io.Reader
	Headers     http.Header
	Prefix      string
	Path        string
	Method      string
}

type TootRequest struct {
	Status     string `json:"status"`
	Visibility string `json:"visibility"`
}

type TootResponse struct {
	ID         string `json:"id"`
	URL        string `json:"url"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	Visibility string `json:"visibility"`
}

type ValidCreds struct {
	Name       string `json:"name"`
	WebsiteURL string `json:"website"`
	VapidKey   string `json:"vapid_key"`
}
