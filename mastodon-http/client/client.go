package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

func AuthenticatedClient(ctx context.Context, config *Config) *Client {
	c := &Client{
		Client: http.Client{Timeout: 10 * time.Second},
		Config: config,
	}

	params := url.Values{
		"client_id":     {config.ClientID},
		"client_secret": {config.ClientSecret},
		"grant_type":    {"client_credentials"},
		"redirect_uri":  {"urn:ietf:wg:oauth:2.0:oob"},
	}

	c.Authenticate(ctx, params)
	return c
}

func (c *Client) Authenticate(ctx context.Context, params url.Values) {
	request := Request{
		RequestBody: strings.NewReader(params.Encode()),
		Headers: http.Header{
			"Content-Type": {"application/x-www-form-urlencoded"},
		},
		Prefix: "Authenticate to mastodon",
		Path:   "/oauth/token",
		Method: http.MethodPost,
	}

	body := c.DoRequest(ctx, request)

	var res struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.Unmarshal(body, &res); err != nil {
		log.Fatalf("Failed to unmarshal authentication response: %v", err)
	}

	c.Config.AccessToken = res.AccessToken
}

func (c *Client) DoRequest(ctx context.Context, r Request) []byte {
	requestUrl, err := url.Parse(c.Config.Server)
	if err != nil {
		log.Fatalf("parse server URL on %s failed: %v", r.Prefix, err)
	}
	requestUrl.Path = path.Join(requestUrl.Path, r.Path)

	req, err := http.NewRequest(r.Method, requestUrl.String(), r.RequestBody)
	if err != nil {
		log.Fatalf("wrap request to %s failed: %v", r.Prefix, err)
	}

	req = req.WithContext(ctx)
	req.Header = r.Headers

	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf("%s failed: %v", r.Prefix, err)
	}
	defer resp.Body.Close()

	responseBody := handleResponseBody(r.Prefix, resp)
	return responseBody
}

func handleResponseBody(origin string, resp *http.Response) []byte {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("%s failed, HTTP response code: %s, reading response body failed: %w", origin, resp.Status, err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("%s failed, HTTP response code: %s, Body: %s", origin, resp.Status, string(body))
	}

	return body
}

func (c *Client) ValidateCredentials(ctx context.Context) {
	request := Request{
		Headers: http.Header{
			"Authorization": {"Bearer " + c.Config.AccessToken},
		},
		Prefix: "validateCredentials",
		Path:   "/api/v1/apps/verify_credentials",
		Method: http.MethodPost,
	}

	c.DoRequest(ctx, request)
}

func (c *Client) PostStatus(ctx context.Context, toot TootRequest) TootResponse {
	payload, err := json.Marshal(toot)
	if err != nil {
		log.Fatalf("postStatus json marshaling failed: %w", err)
	}

	request := Request{
		RequestBody: bytes.NewBuffer(payload),
		Headers: http.Header{
			"Authorization": {"Bearer " + c.Config.AccessToken},
			"Content-Type":  {"application/json"},
		},
		Prefix: "postStatus",
		Path:   "/api/v1/statuses",
		Method: http.MethodPost,
	}

	body := c.DoRequest(ctx, request)

	var statusResponse TootResponse
	if err := json.Unmarshal(body, &statusResponse); err != nil {
		log.Fatalf("postStatus unmarshaling status response failed: %w", err)
	}

	log.Printf("Status Posted: %s, URL: %s", statusResponse.ID, statusResponse.URL)

	return statusResponse
}
