package api

import (
	"io"
	"net/http"
)

type Config struct {
	Server       string
	ClientID     string
	ClientSecret string
	AccessToken  string
}

type Client struct {
	http.Client
	Config *Config
}

type Request struct {
	RequestBody io.Reader
	Headers     http.Header
	Prefix      string
	Path        string
	Method      string
}

type TootRequest struct {
	Status      string `json:"status"`
	Visibility  string `json:"visibility"`
	Sensitive   bool   `json:"sensitive,omitempty"`
	SpoilerText string `json:"spoiler_text,omitempty"`
	Language    string `json:"language,omitempty"`
}

type TootResponse struct {
	ID         string `json:"id"`
	URL        string `json:"url"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	Visibility string `json:"visibility"`
}
